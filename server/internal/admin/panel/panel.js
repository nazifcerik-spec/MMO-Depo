'use strict';
// MeNa World admin panel (DD-360). Talks only to /internal/admin and /internal/ops.
// Everything from the server is written with textContent, never as HTML.

const app = document.getElementById('app');
const TOKEN_KEY = 'mena.admin.token';
let token = sessionStorage.getItem(TOKEN_KEY);
let me = null;
let perms = new Set();

function h(tag, attrs, ...children) {
  const el = document.createElement(tag);
  for (const [k, v] of Object.entries(attrs || {})) {
    if (k === 'on') {
      for (const [ev, fn] of Object.entries(v)) el.addEventListener(ev, fn);
    } else if (k === 'class') {
      el.className = v;
    } else if (v === true) {
      el.setAttribute(k, '');
    } else if (v !== false && v != null) {
      el.setAttribute(k, v);
    }
  }
  for (const c of children.flat(Infinity)) {
    if (c != null && c !== false) el.append(c.nodeType ? c : String(c));
  }
  return el;
}

class ApiError extends Error {
  constructor(status, code, message) {
    super(message || code);
    this.status = status;
    this.code = code;
  }
}

const MESSAGES = {
  E_BAD_CREDENTIALS: 'E-posta, şifre ya da kod hatalı.',
  E_BAD_CODE: 'Kod hatalı.',
  E_FORBIDDEN: 'Bu işlem için iznin yok.',
  E_REASON_REQUIRED: 'Gerekçe yazmalısın.',
  E_NOT_FOUND: 'Bulunamadı.',
  E_INVALID: 'Girilen değer geçersiz.',
  E_INVALID_ID: 'Kimlik biçimi geçersiz.',
  E_LISTING_UNAVAILABLE: 'İlan artık etkin değil.',
  E_REPORT_RESOLVED: 'Bu şikâyet zaten sonuçlanmış.',
  E_TICKET_CLOSED: 'Talep kapalı.',
  E_VALIDATION: 'Girilen değer geçersiz.',
  E_INSUFFICIENT_BALANCE: 'Bakiye yetersiz.',
};
const say = (e) => MESSAGES[e.code] || e.message || 'Bir şeyler ters gitti.';

async function raw(method, path, body) {
  const headers = {};
  if (token) headers.Authorization = 'Bearer ' + token;
  if (body !== undefined) headers['Content-Type'] = 'application/json';
  const res = await fetch(path, { method, headers, body: body === undefined ? undefined : JSON.stringify(body) });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    if (res.status === 401 && token && path !== '/internal/admin/login') signedOut();
    throw new ApiError(res.status, data.error?.code || 'E_' + res.status, data.error?.message);
  }
  return data;
}

// A critical action asks for the code once, then retries.
async function api(method, path, body) {
  try {
    return await raw(method, path, body);
  } catch (e) {
    if (e.code !== 'E_STEPUP_REQUIRED') throw e;
    await stepUp();
    return raw(method, path, body);
  }
}

function stepUp() {
  const dlg = document.getElementById('stepup');
  const form = document.getElementById('stepup-form');
  const input = document.getElementById('stepup-code');
  const err = document.getElementById('stepup-error');
  const cancel = document.getElementById('stepup-cancel');
  input.value = '';
  err.textContent = '';
  dlg.showModal();
  return new Promise((resolve, reject) => {
    const done = () => {
      form.removeEventListener('submit', onSubmit);
      cancel.removeEventListener('click', onCancel);
      dlg.close();
    };
    const onCancel = () => {
      done();
      reject(new ApiError(0, 'E_CANCELLED', 'Vazgeçildi.'));
    };
    const onSubmit = async (ev) => {
      ev.preventDefault();
      try {
        await raw('POST', '/internal/admin/stepup', { code: input.value.trim() });
        done();
        resolve();
      } catch (e) {
        err.textContent = say(e);
      }
    };
    form.addEventListener('submit', onSubmit);
    cancel.addEventListener('click', onCancel);
  });
}

function signedOut() {
  token = null;
  sessionStorage.removeItem(TOKEN_KEY);
  renderLogin();
}

// Runs an action from a button, showing its outcome in the given element.
function action(btn, out, fn) {
  btn.addEventListener('click', async () => {
    btn.disabled = true;
    out.textContent = '';
    out.className = '';
    try {
      const msg = await fn();
      if (msg) {
        out.textContent = msg;
        out.className = 'success';
      }
    } catch (e) {
      if (e.code !== 'E_CANCELLED') {
        out.textContent = say(e);
        out.className = 'error';
      }
    } finally {
      btn.disabled = false;
    }
  });
}

const fmtDate = (s) => (s ? new Date(s).toLocaleString('tr-TR') : '—');
const fmtGold = (c) => {
  const n = Math.abs(Number(c || 0));
  const g = Math.floor(n / 10000), s = Math.floor((n % 10000) / 100), r = n % 100;
  return (c < 0 ? '-' : '') + `${g}a ${s}g ${r}b`;
};

function table(columns, rows, onRow) {
  if (!rows.length) return h('p', { class: 'muted' }, 'Kayıt yok.');
  return h('div', { class: 'tablewrap' }, h('table', null,
    h('thead', null, h('tr', null, columns.map((c) => h('th', { class: c.num ? 'num' : null }, c.title)))),
    h('tbody', null, rows.map((r) => h('tr', { class: onRow ? 'clickable' : null, on: onRow ? { click: () => onRow(r) } : {} },
      columns.map((c) => {
        const v = c.get(r);
        return h('td', { class: [c.num ? 'num' : '', c.mono ? 'mono' : ''].join(' ').trim() || null }, v == null ? '—' : v);
      })))),
  ));
}

function details(obj, keys) {
  return h('dl', null, Object.entries(keys).map(([k, title]) => [h('dt', null, title), h('dd', null, obj[k] == null ? '—' : /_at$/.test(k) ? fmtDate(obj[k]) : String(obj[k]))]));
}

// --- login --------------------------------------------------------------------

function renderLogin() {
  const email = h('input', { type: 'email', autocomplete: 'username', required: true, id: 'login-email' });
  const password = h('input', { type: 'password', autocomplete: 'current-password', required: true, id: 'login-password' });
  const code = h('input', { inputmode: 'numeric', autocomplete: 'one-time-code', maxlength: 6, required: true, id: 'login-code' });
  const err = h('p', { class: 'error' });
  const form = h('form', { class: 'login', on: { submit: async (ev) => {
    ev.preventDefault();
    err.textContent = '';
    try {
      const r = await raw('POST', '/internal/admin/login', { email: email.value.trim(), password: password.value, code: code.value.trim() });
      token = r.token;
      sessionStorage.setItem(TOKEN_KEY, token);
      await start();
    } catch (e) {
      err.textContent = say(e);
    }
  } } },
    h('h1', null, 'MeNa World Yönetim'),
    h('label', null, 'E-posta', email),
    h('label', null, 'Şifre', password),
    h('label', null, 'Doğrulayıcı kodu', code),
    err,
    h('button', { type: 'submit' }, 'Giriş yap'));
  app.replaceChildren(form);
  email.focus();
}

// --- shell --------------------------------------------------------------------

const PAGES = [
  { id: 'overview', title: 'Özet', need: 'player.read', render: pageOverview },
  { id: 'players', title: 'Oyuncular', need: 'player.read', render: pagePlayers },
  { id: 'reports', title: 'Şikâyetler', need: 'ops.reports', render: pageReports },
  { id: 'support', title: 'Destek', need: 'support.tickets', render: pageSupport },
  { id: 'ledger', title: 'Defter', need: 'ledger.read', render: pageLedger },
  { id: 'market', title: 'Pazar', need: 'market.listing_cancel', render: pageMarket },
  { id: 'world', title: 'Lonca ve etkinlik', need: null, render: pageWorld },
  { id: 'errors', title: 'İstemci hataları', need: 'ops.client_errors', render: pageErrors },
  { id: 'audit', title: 'Denetim kaydı', need: 'audit.read', render: pageAudit },
  { id: 'admins', title: 'Yöneticiler', need: 'admin.role_grant', render: pageAdmins },
];

async function start() {
  try {
    const r = await raw('GET', '/internal/admin/session/me');
    me = r.admin;
    perms = new Set(r.permissions);
  } catch (e) {
    signedOut();
    return;
  }
  renderShell(location.hash.slice(1) || 'overview');
}

function renderShell(current) {
  const pages = PAGES.filter((p) => !p.need || perms.has(p.need));
  const page = h('section', { class: 'page' });
  const buttons = {};
  const open = (p, arg) => {
    for (const [id, b] of Object.entries(buttons)) b.classList.toggle('active', id === p.id);
    location.hash = p.id;
    page.replaceChildren();
    p.render(page, arg, (id, a) => open(pages.find((x) => x.id === id), a));
  };
  const nav = h('nav', null,
    h('div', { class: 'brand' }, h('h1', null, 'MeNa World')),
    pages.map((p) => (buttons[p.id] = h('button', { type: 'button', id: 'nav-' + p.id, on: { click: () => open(p) } }, p.title))),
    h('div', { class: 'who' }, me.email, h('br'), (me.roles || []).join(', '), h('br'),
      h('button', { type: 'button', class: 'ghost', id: 'logout', on: { click: async () => {
        try { await raw('POST', '/internal/admin/logout', {}); } catch (_) { /* the session ends here anyway */ }
        signedOut();
      } } }, 'Çıkış')));
  app.replaceChildren(h('div', { class: 'layout' }, nav, page));
  open(pages.find((p) => p.id === current) || pages[0]);
}

function loading(el, fn) {
  el.replaceChildren(h('p', { class: 'muted' }, 'Yükleniyor…'));
  fn().catch((e) => el.replaceChildren(h('p', { class: 'error' }, say(e))));
}

// --- pages --------------------------------------------------------------------

function pageOverview(page) {
  page.append(h('h2', null, 'Özet'));
  const body = h('div');
  page.append(body);
  loading(body, async () => {
    const o = await raw('GET', '/internal/ops/overview');
    const stat = (label, v) => h('div', { class: 'stat' }, h('b', null, Number(v).toLocaleString('tr-TR')), h('span', { class: 'muted' }, label));
    body.replaceChildren(h('div', { class: 'stats' },
      stat('Hesap', o.accounts), stat('Son 24 saat yeni hesap', o.accounts_new_24h), stat('Son 24 saat giriş', o.logins_24h),
      stat('Etkin karakter', o.characters), stat('Çevrimiçi bağlantı', o.online),
      stat('Açık şikâyet', o.open_reports), stat('Açık destek talebi', o.open_tickets),
      stat('Etkin ilan', o.active_listings), stat('Etkin açık artırma', o.active_auctions),
      stat('Lonca', o.guilds), stat('Son 24 saat istemci hatası', o.client_errors_24h)));
  });
}

function pagePlayers(page, arg) {
  const q = h('input', { class: 'grow', id: 'player-q', placeholder: 'Kullanıcı adı, e-posta ya da karakter adı' });
  const list = h('div');
  const detail = h('div');
  const search = () => {
    detail.replaceChildren();
    loading(list, async () => {
      const r = await raw('GET', '/internal/ops/players?q=' + encodeURIComponent(q.value.trim()));
      list.replaceChildren(table([
        { title: 'Kullanıcı', get: (p) => p.username },
        { title: 'E-posta', get: (p) => p.email },
        { title: 'Durum', get: (p) => p.state },
        { title: 'Son giriş', get: (p) => fmtDate(p.last_login_at) },
      ], r.players || [], (p) => openPlayer(p.id)));
    });
  };
  const openPlayer = (id) => loading(detail, async () => {
    const d = await raw('GET', '/internal/ops/players/' + id);
    const out = h('span');
    const reason = h('input', { class: 'grow', placeholder: 'Gerekçe (zorunlu)', id: 'sanction-reason' });
    const hours = h('input', { type: 'number', min: 1, value: 24, id: 'mute-hours' });
    const days = h('input', { type: 'number', min: 0, value: 7, id: 'ban-days' });
    const act = (label, cls, path, payload, done) => {
      const b = h('button', { type: 'button', class: cls });
      b.textContent = label;
      action(b, out, async () => {
        await api('POST', '/internal/ops/players/' + id + path, { reason: reason.value.trim(), ...payload() });
        openPlayer(id);
        return done;
      });
      return b;
    };
    detail.replaceChildren(h('div', { class: 'card' },
      h('h2', null, d.username || id),
      details(d, { id: 'Hesap', username: 'Kullanıcı', email: 'E-posta', state: 'Durum', created_at: 'Kayıt', last_login_at: 'Son giriş' }),
      h('h3', null, 'Karakterler'),
      table([
        { title: 'Ad', get: (c) => c.name },
        { title: 'Sınıf', get: (c) => c.class_code },
        { title: 'Sv', num: true, get: (c) => c.level },
        { title: 'Durum', get: (c) => c.state },
      ], d.characters || [], (c) => openCharacter(c.id)),
      h('div', { id: 'character-detail' }),
      h('h3', null, 'Yaptırımlar'),
      table([
        { title: 'Tür', get: (s) => s.kind },
        { title: 'Bitiş', get: (s) => (s.until ? fmtDate(s.until) : 'süresiz') },
        { title: 'Gerekçe', get: (s) => s.reason },
        { title: 'Kaldırıldı', get: (s) => fmtDate(s.revoked_at) },
      ], d.sanctions || []),
      h('h3', null, 'İşlem'),
      h('div', { class: 'row' }, reason),
      h('div', { class: 'row' },
        h('label', null, 'Saat', hours), act('Sustur', 'ghost', '/mute', () => ({ hours: Number(hours.value) }), 'Susturuldu.'),
        act('Susturmayı kaldır', 'ghost', '/unmute', () => ({}), 'Kaldırıldı.'),
        h('label', null, 'Gün (0 = süresiz)', days), act('Yasakla', 'danger', '/ban', () => ({ days: Number(days.value) }), 'Yasaklandı.'),
        act('Yasağı kaldır', 'ghost', '/unban', () => ({}), 'Kaldırıldı.')),
      out));
  });
  const openCharacter = (id) => {
    const box = document.getElementById('character-detail');
    loading(box, async () => {
      const c = await raw('GET', '/internal/ops/characters/' + id);
      const out = h('span');
      const reason = h('input', { class: 'grow', placeholder: 'Gerekçe (zorunlu)' });
      const freeze = h('button', { type: 'button', class: 'danger' }, c.state === 'frozen' ? 'Çözdür' : 'Dondur');
      action(freeze, out, async () => {
        await api('POST', '/internal/ops/characters/' + id + (c.state === 'frozen' ? '/unfreeze' : '/freeze'), { reason: reason.value.trim() });
        openCharacter(id);
        return 'Tamam.';
      });
      box.replaceChildren(h('div', { class: 'card' },
        details(c, { id: 'Kimlik', name: 'Ad', class_code: 'Sınıf', level: 'Seviye', state: 'Durum', location_region_code: 'Bölge', gold_copper: 'Altın (bakır)', item_count: 'Eşya' }),
        h('div', { class: 'row' }, reason, freeze), out));
    });
  };
  page.append(h('h2', null, 'Oyuncular'),
    h('form', { class: 'row', on: { submit: (ev) => { ev.preventDefault(); search(); } } }, q, h('button', { type: 'submit', id: 'player-search' }, 'Ara')),
    h('div', { class: 'card' }, list), detail);
  if (arg?.account) openPlayer(arg.account);
}

const REPORT_REASONS = { harassment: 'Taciz / hakaret', cheating: 'Hile / bot', name: 'Uygunsuz isim', spam: 'Spam / reklam', scam: 'Dolandırıcılık', other: 'Diğer' };

function pageReports(page, _arg, go) {
  const state = h('select', { id: 'report-state' }, h('option', { value: 'open' }, 'Açık'), h('option', { value: 'actioned' }, 'İşlem yapıldı'), h('option', { value: 'dismissed' }, 'Reddedildi'));
  const list = h('div');
  const load = () => loading(list, async () => {
    const r = await raw('GET', '/internal/ops/reports?limit=100&state=' + state.value);
    const rows = r.reports || [];
    if (!rows.length) {
      list.replaceChildren(h('p', { class: 'muted' }, 'Kayıt yok.'));
      return;
    }
    list.replaceChildren(...rows.map((x) => {
      const note = h('input', { class: 'grow', placeholder: 'Karar notu' });
      const out = h('span');
      const resolve = (label, st, cls) => {
        const b = h('button', { type: 'button', class: cls }, label);
        action(b, out, async () => {
          await api('POST', '/internal/ops/reports/' + x.id + '/resolve', { state: st, note: note.value.trim() });
          load();
          return 'Tamam.';
        });
        return b;
      };
      return h('div', { class: 'card report' },
        h('p', null, h('b', null, x.target_name), ' · ', REPORT_REASONS[x.category] || x.category, ' ',
          x.open_against > 1 ? h('span', { class: 'badge warn' }, x.open_against + ' açık şikâyet') : null),
        h('p', { class: 'muted' }, 'Bildiren: ', x.reporter_name || 'silinmiş karakter', ' · ', fmtDate(x.created_at)),
        x.note ? h('p', null, x.note) : null,
        x.chat_excerpt ? h('div', { class: 'quote' }, x.chat_excerpt) : null,
        x.state === 'open'
          ? [h('div', { class: 'row' }, note, resolve('İşlem yapıldı', 'actioned', 'danger'), resolve('Reddet', 'dismissed', 'ghost'),
              x.target_account && perms.has('player.read')
                ? h('button', { type: 'button', class: 'ghost', on: { click: () => go('players', { account: x.target_account }) } }, 'Oyuncuyu aç') : null), out]
          : h('p', { class: 'muted' }, 'Karar: ', x.resolution_note || '—'));
    }));
  });
  state.addEventListener('change', load);
  page.append(h('h2', null, 'Şikâyetler'), h('div', { class: 'row' }, h('label', null, 'Durum', state)), h('div', { class: 'spacer' }), list);
  load();
}

function pageSupport(page) {
  const state = h('select', null, h('option', { value: 'open' }, 'Açık'), h('option', { value: 'answered' }, 'Yanıtlandı'), h('option', { value: 'closed' }, 'Kapalı'));
  const list = h('div');
  const load = () => loading(list, async () => {
    const r = await raw('GET', '/internal/ops/support/tickets?limit=100&state=' + state.value);
    const rows = r.tickets || [];
    if (!rows.length) {
      list.replaceChildren(h('p', { class: 'muted' }, 'Kayıt yok.'));
      return;
    }
    list.replaceChildren(...rows.map((t) => {
      const reply = h('textarea', { placeholder: 'Yanıt' });
      const close = h('input', { type: 'checkbox' });
      const out = h('span');
      const send = h('button', { type: 'button' }, 'Yanıtla');
      action(send, out, async () => {
        await api('POST', '/internal/ops/support/tickets/' + t.id + '/reply', { reply: reply.value.trim(), close: close.checked });
        load();
        return 'Gönderildi.';
      });
      return h('div', { class: 'card' },
        h('p', null, h('b', null, t.subject), ' · ', h('span', { class: 'badge' }, t.category)),
        h('p', { class: 'muted' }, fmtDate(t.created_at)),
        h('p', null, t.body),
        t.reply ? h('div', { class: 'quote' }, t.reply) : null,
        t.state !== 'closed' ? [reply, h('div', { class: 'row' }, h('label', null, 'Kapat', close), send), out] : null);
    }));
  });
  state.addEventListener('change', load);
  page.append(h('h2', null, 'Destek'), h('div', { class: 'row' }, h('label', null, 'Durum', state)), list);
  load();
}

function pageLedger(page) {
  const account = h('input', { class: 'grow', placeholder: 'Karakter ya da hesap kimliği' });
  const txn = h('input', { class: 'grow', placeholder: 'İşlem kimliği' });
  const list = h('div');
  const search = () => loading(list, async () => {
    const q = account.value.trim() ? 'account=' + encodeURIComponent(account.value.trim()) : 'txn=' + encodeURIComponent(txn.value.trim());
    const r = await raw('GET', '/internal/ops/ledger?limit=200&' + q);
    const reversal = perms.has('ledger.reversal<=100k_gold') || perms.has('ledger.reversal_any');
    list.replaceChildren(table([
      { title: 'Zaman', get: (e) => fmtDate(e.created_at) },
      { title: 'Tür', get: (e) => e.kind },
      { title: 'Hesap', get: (e) => e.account_kind },
      { title: 'Tutar', num: true, get: (e) => fmtGold(e.delta) },
      { title: 'Sonra', num: true, get: (e) => fmtGold(e.balance_after) },
      { title: 'İşlem', mono: true, get: (e) => e.txn_id },
    ], r.entries || [], reversal ? (e) => openReverse(e) : null), detail);
  });
  const detail = h('div');
  const openReverse = (e) => {
    const reason = h('input', { class: 'grow', placeholder: 'Gerekçe (zorunlu)' });
    const out = h('span');
    const b = h('button', { type: 'button', class: 'danger' }, 'Bu işlemi ters çevir');
    action(b, out, async () => {
      await api('POST', '/internal/ops/ledger/' + e.txn_id + '/reverse', { reason: reason.value.trim() });
      return 'Ters kayıt yazıldı.';
    });
    detail.replaceChildren(h('div', { class: 'card' }, h('p', { class: 'mono' }, e.txn_id), h('div', { class: 'row' }, reason, b), out));
  };
  page.append(h('h2', null, 'Defter'),
    h('form', { class: 'row', on: { submit: (ev) => { ev.preventDefault(); search(); } } }, account, txn, h('button', { type: 'submit' }, 'Ara')),
    h('div', { class: 'card' }, list));
}

function pageMarket(page) {
  const list = h('div');
  const reason = h('input', { class: 'grow', placeholder: 'Gerekçe (zorunlu)' });
  const refund = h('input', { type: 'checkbox' });
  const out = h('p');
  const load = () => loading(list, async () => {
    const r = await raw('GET', '/internal/ops/market/listings');
    const rows = r.listings || [];
    list.replaceChildren(table([
      { title: 'Eşya', get: (l) => l.template_code },
      { title: 'Adet', num: true, get: (l) => l.stack_quantity },
      { title: 'Tür', get: (l) => (l.kind === 'auction' ? 'Açık artırma' : 'Sabit') },
      { title: 'Fiyat', num: true, get: (l) => fmtGold(l.current_bid_copper || l.total_price_copper) },
      { title: 'Satıcı', get: (l) => l.seller_name },
      { title: 'Bitiş', get: (l) => fmtDate(l.expires_at) },
      { title: '', get: (l) => {
        const b = h('button', { type: 'button', class: 'danger' }, 'İptal');
        action(b, out, async () => {
          await api('POST', '/internal/ops/market/' + l.id + '/cancel', { reason: reason.value.trim(), refund_fee: refund.checked });
          load();
          return 'İlan iptal edildi.';
        });
        return b;
      } },
    ], rows));
  });
  page.append(h('h2', null, 'Pazar'), h('div', { class: 'row' }, reason, h('label', null, 'Ücreti iade et', refund)), out, h('div', { class: 'card' }, list));
  load();
}

function pageWorld(page) {
  const blocks = [h('h2', null, 'Lonca ve etkinlik')];
  if (perms.has('guild.disband_override')) {
    const id = h('input', { class: 'grow', placeholder: 'Lonca kimliği' });
    const reason = h('input', { class: 'grow', placeholder: 'Gerekçe (zorunlu)' });
    const out = h('span');
    const b = h('button', { type: 'button', class: 'danger' }, 'Loncayı dağıt');
    action(b, out, async () => {
      await api('POST', '/internal/ops/guilds/' + encodeURIComponent(id.value.trim()) + '/disband', { reason: reason.value.trim() });
      return 'Lonca dağıtıldı.';
    });
    blocks.push(h('div', { class: 'card' }, h('h3', null, 'Lonca dağıtma'),
      h('p', { class: 'muted' }, 'Kasasındaki altın ve bankasındaki eşyalar erişilemez kalır. Önce lideriyle konuş.'),
      h('div', { class: 'row' }, id, reason, b), out));
  }
  if (perms.has('content.event_trigger')) {
    const code = h('input', { class: 'grow', placeholder: 'Etkinlik şablonu kodu' });
    const mins = h('input', { type: 'number', min: 0, value: 0 });
    const reason = h('input', { class: 'grow', placeholder: 'Gerekçe (zorunlu)' });
    const out = h('span');
    const b = h('button', { type: 'button' }, 'Başlat');
    action(b, out, async () => {
      const r = await api('POST', '/internal/ops/world-events/' + encodeURIComponent(code.value.trim()) + '/trigger',
        { starts_in_minutes: Number(mins.value), reason: reason.value.trim() });
      return 'Başlatıldı: ' + r.spawn_id;
    });
    blocks.push(h('div', { class: 'card' }, h('h3', null, 'Dünya etkinliği'), h('div', { class: 'row' }, code, h('label', null, 'Kaç dakika sonra', mins), reason, b), out));
  }
  if (blocks.length === 1) blocks.push(h('p', { class: 'muted' }, 'Bu bölüm için iznin yok.'));
  page.append(...blocks);
}

function pageErrors(page) {
  const list = h('div');
  page.append(h('h2', null, 'İstemci hataları'), h('div', { class: 'card' }, list));
  loading(list, async () => {
    const r = await raw('GET', '/internal/ops/client-errors?limit=100');
    list.replaceChildren(table([
      { title: 'Sayı', num: true, get: (e) => e.count },
      { title: 'Tür', get: (e) => e.kind },
      { title: 'Mesaj', get: (e) => e.message },
      { title: 'Platform', get: (e) => e.platform },
      { title: 'Sürüm', get: (e) => e.app_version },
      { title: 'Son', get: (e) => fmtDate(e.last_seen) },
    ], r.errors || []));
  });
}

function pageAudit(page) {
  const list = h('div');
  page.append(h('h2', null, 'Denetim kaydı'), h('div', { class: 'card' }, list));
  loading(list, async () => {
    const r = await raw('GET', '/internal/ops/audit?limit=200');
    list.replaceChildren(table([
      { title: 'Zaman', get: (e) => fmtDate(e.created_at || e.at) },
      { title: 'Yönetici', mono: true, get: (e) => e.admin_user_id },
      { title: 'İşlem', get: (e) => e.action },
      { title: 'Hedef', mono: true, get: (e) => e.target_id },
      { title: 'Gerekçe', get: (e) => e.reason },
    ], r.actions || []));
  });
}

const ROLES = ['viewer', 'support_l1', 'support_l2', 'game_designer', 'ops_engineer', 'admin'];

function pageAdmins(page) {
  const list = h('div');
  const load = () => loading(list, async () => {
    const r = await raw('GET', '/internal/admin/users');
    list.replaceChildren(...(r.users || []).map((u) => {
      const role = h('select', null, ROLES.map((x) => h('option', { value: x }, x)));
      const out = h('span');
      const set = (grant, label, cls) => {
        const b = h('button', { type: 'button', class: cls }, label);
        action(b, out, async () => {
          await api('POST', '/internal/admin/users/' + u.id + '/role', { role: role.value, grant });
          load();
          return 'Tamam.';
        });
        return b;
      };
      return h('div', { class: 'card' }, h('p', null, h('b', null, u.email), ' · ', (u.roles || []).join(', ')),
        h('div', { class: 'row' }, role, set(true, 'Rol ver', ''), set(false, 'Rolü al', 'ghost')), out);
    }));
  });
  page.append(h('h2', null, 'Yöneticiler'), list);
  load();
}

if (token) start();
else renderLogin();
