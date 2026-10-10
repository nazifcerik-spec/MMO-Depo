// Logins at a fixed arrival rate from many client IPs against pre-registered
// accounts (DD-238, docs/performance/load-login.md).
//   scripts/loadtest.sh login [duration]          RATE=logins/s (default 10)
//
// RATE <= 20 is the background load (§2.1): p95 < 300 ms, errors < 0.5%.
// Above that it is a storm (§2.2): the server is expected to shed with a fast
// 503 + Retry-After (bounded argon2id slots) — what must hold is that the
// logins it accepts stay inside p95 < 800 ms and nothing else fails.
import http from 'k6/http';
import { check } from 'k6';
import { Counter, Trend } from 'k6/metrics';

const BASE = __ENV.BASE || 'http://api:8080';
const RATE = Number(__ENV.RATE || 10);
const ACCOUNTS = Number(__ENV.ACCOUNTS || 200);
const DURATION = __ENV.DURATION || '1m';
const STORM = RATE > 20;

const okLatency = new Trend('login_ok_ms', true);
const shed = new Counter('login_shed_503');
const failed = new Counter('login_failed_other');

export const options = {
  setupTimeout: '10m',
  scenarios: {
    logins: {
      executor: 'constant-arrival-rate', exec: 'login', rate: RATE, timeUnit: '1s', duration: DURATION,
      preAllocatedVUs: Math.min(RATE * 4, 2000), maxVUs: 3000,
    },
  },
  thresholds: STORM
    ? { login_ok_ms: ['p(95)<800'], login_failed_other: ['count==0'] }
    : { login_ok_ms: ['p(95)<300'], login_shed_503: ['count==0'], 'http_req_failed{name:login}': ['rate<0.005'] },
};

const rnd = () => Math.floor(Math.random() * 250) + 1;
const ip = () => `10.${rnd()}.${rnd()}.${rnd()}`;
const headers = () => ({ 'Content-Type': 'application/json', 'X-Forwarded-For': ip(), 'Idempotency-Key': `${Date.now()}-${Math.random()}` });

export function setup() {
  const run = Math.random().toString(36).slice(2, 8);
  const users = [];
  for (let i = 0; i < ACCOUNTS; i++) {
    const u = `l${run}${i}`;
    const res = http.post(`${BASE}/v1/auth/register`, JSON.stringify({
      email: `${u}@example.com`, username: u, password: `perf-pass-${u}`, accept_terms: true,
    }), { headers: headers(), tags: { name: 'register' } });
    if (res.status === 201 || res.status === 200) users.push(u);
  }
  if (users.length < ACCOUNTS) throw new Error(`registered ${users.length}/${ACCOUNTS}`);
  return { users };
}

export function login(data) {
  const u = data.users[Math.floor(Math.random() * data.users.length)];
  const res = http.post(`${BASE}/v1/auth/login`, JSON.stringify({ email_or_username: u, password: `perf-pass-${u}` }),
    { headers: headers(), tags: { name: 'login' } });
  if (res.status === 200) okLatency.add(res.timings.duration);
  else if (res.status === 503 && res.headers['Retry-After']) shed.add(1);
  else failed.add(1);
  check(res, { 'login 200 or shed': (r) => r.status === 200 || r.status === 503 });
}
