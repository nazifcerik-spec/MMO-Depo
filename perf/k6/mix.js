// Load mix (DD-238, docs/performance/ccu-scenarios.md). Run via
//   scripts/loadtest.sh [smoke|t1] [duration]
//
// Active players (30%) play over HTTP at ~40 rps per 300 players with the
// doc's action weights; AFK players (70%) hold the WebSocket, send a
// heartbeat every minute and measure global-chat fanout latency. Every VU
// is its own player and its own client IP (X-Forwarded-For, trusted only
// in the perf run), so the per-IP and per-account rate limits apply as in
// production. Thresholds are the SLOs in docs/performance/slo-slis.md; a
// miss fails the run.
import http from 'k6/http';
import ws from 'k6/ws';
import { check, sleep } from 'k6';
import { Counter, Trend } from 'k6/metrics';
import exec from 'k6/execution';

const BASE = __ENV.BASE || 'http://api:8080';
const TIER = __ENV.TIER || 'smoke';
const TIERS = { smoke: { active: 30, afk: 70, ramp: 20 }, t1: { active: 300, afk: 700, ramp: 90 } };
const tier = TIERS[TIER];
const DURATION = __ENV.DURATION || (TIER === 't1' ? '10m' : '3m');
const PACE = 7.5; // seconds between an active player's actions → 40 rps per 300 players

const fanout = new Trend('ws_fanout_ms', true);
const rateLimited = new Counter('rate_limited');

export const options = {
  scenarios: {
    active: { executor: 'constant-vus', exec: 'active', vus: tier.active, duration: DURATION },
    afk: { executor: 'constant-vus', exec: 'afk', vus: tier.afk, duration: DURATION, gracefulStop: '10s' },
  },
  thresholds: {
    http_req_failed: ['rate<0.005'],
    checks: ['rate>0.99'],
    rate_limited: ['count==0'],
    'http_req_duration{name:sheet}': ['p(95)<150', 'p(99)<500'],
    'http_req_duration{name:inventory}': ['p(95)<250', 'p(99)<700'],
    'http_req_duration{name:reveal}': ['p(95)<250', 'p(99)<700'],
    'http_req_duration{name:vendor_buy}': ['p(95)<200', 'p(99)<600'],
    'http_req_duration{name:market_search}': ['p(95)<200', 'p(99)<500'],
    ws_fanout_ms: ['p(95)<100', 'p(99)<400'],
  },
};

// Domain conflicts (e.g. an activity already running) are answers, not failures.
http.setResponseCallback(http.expectedStatuses({ min: 200, max: 299 }, 409));

export function setup() {
  return { run: Math.random().toString(36).slice(2, 8) };
}

let me = null; // per VU (every VU has its own JS runtime)

function ip() {
  const id = exec.vu.idInTest;
  return `10.${(id >> 16) & 255}.${(id >> 8) & 255}.${id & 255}`;
}

function uuid() {
  return 'xxxxxxxx-xxxx-4xxx-8xxx-xxxxxxxxxxxx'.replace(/x/g, () => ((Math.random() * 16) | 0).toString(16));
}

function req(method, path, body, name) {
  const headers = { 'Content-Type': 'application/json', 'X-Forwarded-For': ip() };
  if (me && me.token) headers.Authorization = `Bearer ${me.token}`;
  if (me && me.char) headers['X-Character-Id'] = me.char;
  if (method !== 'GET') headers['Idempotency-Key'] = uuid();
  const res = http.request(method, BASE + path, body === undefined ? null : JSON.stringify(body), {
    headers,
    tags: { name },
  });
  if (res.status === 429) rateLimited.add(1);
  check(res, { [`${name} ok`]: (r) => (r.status >= 200 && r.status < 300) || r.status === 409 });
  return res;
}

function json(res) {
  try {
    return res.json();
  } catch (_) {
    return {};
  }
}

// Registers, creates and enters a character, funds it. Starts are spread
// over the ramp so sign-up (argon2id) is not one thundering herd.
function bootstrap(data) {
  sleep(Math.random() * tier.ramp);
  const id = exec.vu.idInTest;
  const handle = `p${data.run}${id}`;
  me = {};
  const reg = json(req('POST', '/v1/auth/register', {
    email: `${handle}@example.com`, username: handle, password: `perf-pass-${handle}`, accept_terms: true,
  }, 'register'));
  me.token = reg.access_token;
  const created = json(req('POST', '/v1/character', { name: `P${data.run}${id}`.slice(0, 16), race_code: 'human', class_code: 'warrior' }, 'char_create'));
  me.char = created.id;
  req('POST', `/v1/character/${me.char}/enter`, {}, 'char_enter');
  req('POST', '/v1/dev/mint', { character_id: me.char, amount: 1000000, currency: 'gold' }, 'mint');
  const vendor = json(req('GET', '/v1/vendor/general_ironhold', undefined, 'vendor'));
  const potion = (vendor.sells || []).find((o) => o.template_code === 'potion_hp_small');
  me.potionPrice = potion ? potion.buy_price_copper : 0;
}

const c = () => `/v1/character/${me.char}`;

const actions = [
  [6, () => req('GET', `${c()}/sheet`, undefined, 'sheet')],
  [4, () => req('GET', `${c()}/inventory`, undefined, 'inventory')],
  [4, () => {
    const active = json(req('GET', `${c()}/activity/active`, undefined, 'activity_active'));
    if (active.contract) req('POST', `/v1/activity/${active.contract.id}/reveal`, {}, 'reveal');
    else req('POST', `${c()}/activity`, { kind: 'mining', duration_seconds: 1800 }, 'activity_start');
  }],
  [1.5, () => req('POST', '/v1/vendor/general_ironhold/buy', {
    template_code: 'potion_hp_small', expected_price_copper: me.potionPrice, character_id: me.char,
  }, 'vendor_buy')],
  [3.5, () => req('GET', '/v1/market/listings?sort=price_asc&limit=25&offset=0', undefined, 'market_search')],
  [2, () => req('GET', `${c()}/wallet`, undefined, 'wallet')],
  [2, () => req('GET', '/v1/notifications?unread=1', undefined, 'notifications')],
  [1, () => req('GET', `${c()}/quests`, undefined, 'quests')],
  // ≈4% of actions: 1.7 chat/s at T1 (ccu-scenarios §4).
  [1, () => req('POST', `${c()}/chat`, { channel: 'global', body: `perf ${Date.now()}` }, 'chat_send')],
];
const totalWeight = actions.reduce((s, [w]) => s + w, 0);

function pick() {
  let r = Math.random() * totalWeight;
  for (const [w, fn] of actions) {
    if ((r -= w) < 0) return fn;
  }
  return actions[0][1];
}

export function active(data) {
  if (!me) bootstrap(data);
  pick()();
  sleep(PACE * (0.5 + Math.random()));
}

export function afk(data) {
  if (!me) bootstrap(data);
  const url = `${BASE.replace(/^http/, 'ws')}/v1/ws?token=${me.token}`;
  const res = ws.connect(url, { headers: { 'X-Forwarded-For': ip() }, tags: { name: 'ws' } }, (socket) => {
    socket.on('message', (raw) => {
      let ev;
      try {
        ev = JSON.parse(raw);
      } catch (_) {
        return;
      }
      if (ev.kind === 'chat.message' && ev.payload && ev.payload.sent_at) {
        fanout.add(Date.now() - Date.parse(ev.payload.sent_at));
      }
    });
    socket.setInterval(() => req('GET', '/v1/notifications?unread=1', undefined, 'heartbeat'), 60000 * (0.8 + Math.random() * 0.4));
    socket.setTimeout(() => socket.close(), 30 * 60 * 1000); // the scenario's end interrupts first
  });
  check(res, { 'ws upgraded': (r) => r && r.status === 101 });
  sleep(5); // reconnect pause if the server dropped us
}
