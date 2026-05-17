// ─── API ──────────────────────────────────────────────────────────────────────
const API_BASE = 'http://localhost:8080/api/v1';

async function kareApi(method, path, body) {
  const token = localStorage.getItem('kare_token');
  const opts = { method, headers: { 'Content-Type': 'application/json' } };
  if (token) opts.headers['Authorization'] = 'Bearer ' + token;
  if (body) opts.body = JSON.stringify(body);
  try {
    const res = await fetch(API_BASE + path, opts);
    if (res.status === 401) { kareLogout(); return null; }
    const data = await res.json();
    if (!res.ok) { kareToast(data.error || 'Bir hata oluştu', 'danger'); return null; }
    return data;
  } catch (e) {
    kareToast('Sunucuya bağlanılamıyor', 'danger');
    return null;
  }
}

// ─── AUTH ─────────────────────────────────────────────────────────────────────
async function kareLogin(username, password) {
  const data = await kareApi('POST', '/auth/login', { username, password });
  if (!data || !data.token) return null;
  localStorage.setItem('kare_token', data.token);
  localStorage.setItem('kare_user', JSON.stringify(data.user));
  return data.user;
}

function kareLogout() {
  localStorage.removeItem('kare_token');
  localStorage.removeItem('kare_user');
  const inSub = /\/(admin|koc|koordinator|veli|ogrenci)\//.test(location.pathname);
  location.href = inSub ? '../index.html' : 'index.html';
}

// Checks token + optional role whitelist. Returns false and redirects if invalid.
function kareGuard(roles) {
  const token = localStorage.getItem('kare_token');
  const user = JSON.parse(localStorage.getItem('kare_user') || 'null');
  const inSub = /\/(admin|koc|koordinator|veli|ogrenci)\//.test(location.pathname);
  const loginHref = inSub ? '../index.html' : 'index.html';
  if (!token || !user) { location.href = loginHref; return false; }
  if (roles && !roles.includes(user.role)) { location.href = loginHref; return false; }
  return true;
}

function kareCurrentUser() {
  return JSON.parse(localStorage.getItem('kare_user') || 'null');
}

// ─── MENU ─────────────────────────────────────────────────────────────────────
const KARE_MENUS = {
  admin: [
    { section: 'Genel', items: [
      { href: 'dashboard.html', icon: 'tabler-smart-home', label: 'Panel' },
      { href: 'kullanicilar.html', icon: 'tabler-users', label: 'Kullanıcılar' },
    ]},
    { section: 'Eşleştirme', items: [
      { href: 'eslestirme-koc.html', icon: 'tabler-link', label: 'Öğrenci–Koç' },
      { href: 'eslestirme-koordinator.html', icon: 'tabler-building', label: 'Öğrenci–Koordinatör' },
    ]},
    { section: 'Takip', items: [
      { href: 'gorusmeler.html', icon: 'tabler-notes', label: 'Görüşme Kayıtları' },
      { href: 'koc-uyarilar.html', icon: 'tabler-alert-triangle', label: 'Koç Uyarıları' },
      { href: 'haftalar.html', icon: 'tabler-calendar', label: 'Hafta Tanımları' },
    ]},
    { section: 'İletişim', items: [
      { href: 'mesajlar.html', icon: 'tabler-mail', label: 'Mesajlar' },
      { href: 'sms.html', icon: 'tabler-message', label: 'SMS' },
    ]},
    { section: 'Sistem', items: [
      { href: 'raporlar.html', icon: 'tabler-chart-bar', label: 'Raporlar' },
      { href: 'loglar.html', icon: 'tabler-history', label: 'Log Kayıtları' },
    ]},
  ],
  koc: [
    { section: 'Koç', items: [
      { href: 'dashboard.html', icon: 'tabler-smart-home', label: 'Panel' },
      { href: 'ogrenciler.html', icon: 'tabler-school', label: 'Öğrencilerim' },
      { href: 'gorusme.html', icon: 'tabler-forms', label: 'Görüşme Formu' },
      { href: 'mesajlar.html', icon: 'tabler-mail', label: 'Mesajlar' },
    ]},
  ],
  koordinator: [
    { section: 'Koordinatör', items: [
      { href: 'dashboard.html', icon: 'tabler-smart-home', label: 'Panel' },
      { href: 'ogrenciler.html', icon: 'tabler-users', label: 'Öğrenciler' },
      { href: 'gorusmeler.html', icon: 'tabler-notes', label: 'Görüşmeler' },
    ]},
  ],
  veli: [
    { section: 'Veli', items: [
      { href: 'dashboard.html', icon: 'tabler-smart-home', label: 'Panel' },
      { href: 'gorusmeler.html', icon: 'tabler-notes', label: 'Onaylı Görüşmeler' },
      { href: 'mesajlar.html', icon: 'tabler-mail', label: 'Mesajlar' },
    ]},
  ],
  ogrenci: [
    { section: 'Öğrenci', items: [
      { href: 'dashboard.html', icon: 'tabler-smart-home', label: 'Panel' },
      { href: 'mesajlar.html', icon: 'tabler-mail', label: 'Mesajlar' },
    ]},
  ],
};

function kareMenuHtml(role, prefix, exitHref) {
  const menu = KARE_MENUS[role] || KARE_MENUS.admin;
  const current = location.pathname.split('/').pop();
  let html = '';
  menu.forEach((g) => {
    html += `<li class="menu-header small"><span class="menu-header-text">${g.section}</span></li>`;
    g.items.forEach((item) => {
      const active = current === item.href ? ' active' : '';
      html += `<li class="menu-item${active}"><a href="${prefix}${item.href}" class="menu-link"><i class="menu-icon icon-base ti ${item.icon}"></i><div>${item.label}</div></a></li>`;
    });
  });
  html += `<li class="menu-item"><a href="javascript:void(0)" class="menu-link" onclick="kareLogout()"><i class="menu-icon icon-base ti tabler-logout"></i><div>Çıkış</div></a></li>`;
  return html;
}

// ─── UI HELPERS ───────────────────────────────────────────────────────────────
function kareToast(msg, type) {
  type = type || 'success';
  const el = document.createElement('div');
  el.className = `alert alert-${type} alert-dismissible fade show position-fixed bottom-0 end-0 m-4`;
  el.style.zIndex = 9999;
  el.innerHTML = msg + '<button type="button" class="btn-close" data-bs-dismiss="alert"></button>';
  document.body.appendChild(el);
  setTimeout(() => el.remove(), 3500);
}

function kareSpinner() {
  return '<div class="text-center py-4"><div class="spinner-border text-primary"></div></div>';
}

function kareSkeleton(rows, cols) {
  const widths = [55, 80, 65, 90, 45, 70];
  let html = '';
  for (let r = 0; r < rows; r++) {
    html += '<tr class="kare-skel-row">';
    for (let c = 0; c < cols; c++) {
      const w = widths[(r * cols + c) % widths.length];
      html += `<td><span class="kare-skel" style="width:${w}%"></span></td>`;
    }
    html += '</tr>';
  }
  return html;
}

function kareEmptyState(icon, title, desc, btnHtml) {
  return `<tr><td colspan="99"><div class="kare-empty">
    <i class="ti ${icon} kare-empty-icon"></i>
    <h6>${title}</h6>
    <p>${desc}</p>
    ${btnHtml || ''}
  </div></td></tr>`;
}

function kareMarkInvalid(el, msg) {
  el.classList.add('is-invalid');
  let fb = el.nextElementSibling;
  if (!fb || !fb.classList.contains('invalid-feedback')) {
    fb = document.createElement('div');
    fb.className = 'invalid-feedback';
    el.parentNode.insertBefore(fb, el.nextSibling);
  }
  fb.textContent = msg || 'Bu alan gerekli';
}

function kareClearInvalid(el) {
  el.classList.remove('is-invalid');
  const fb = el.nextElementSibling;
  if (fb && fb.classList.contains('invalid-feedback')) fb.remove();
}

function kareBadge(status) {
  const map = {
    pending:   ['warning', 'Onay Bekliyor'],
    approved:  ['success', 'Onaylandı'],
    rejected:  ['danger',  'Reddedildi'],
    pre:       ['warning', 'Ön Kayıt'],
    confirmed: ['success', 'Kesin Kayıt'],
    active:    ['success', 'Aktif'],
    passive:   ['secondary', 'Pasif'],
    admin:     ['danger',  'Admin'],
    koc:       ['primary', 'Koç'],
    koordinator: ['warning', 'Koordinatör'],
    veli:      ['info',    'Veli'],
    ogrenci:   ['success', 'Öğrenci'],
  };
  const [cls, label] = map[status] || ['secondary', status];
  return `<span class="badge bg-label-${cls}">${label}</span>`;
}

// ─── INIT ─────────────────────────────────────────────────────────────────────
document.addEventListener('DOMContentLoaded', () => {
  // Rating buttons
  document.querySelectorAll('.rating-btn').forEach((btn) => {
    btn.addEventListener('click', () => {
      btn.parentElement.querySelectorAll('.rating-btn').forEach((b) => {
        b.classList.remove('active', 'btn-primary');
        b.classList.add('btn-outline-primary');
      });
      btn.classList.add('active', 'btn-primary');
      btn.classList.remove('btn-outline-primary');
    });
  });

  // Render sidebar menu
  const menuEl = document.getElementById('kare-menu-inner');
  if (menuEl && menuEl.dataset.role) {
    menuEl.innerHTML = kareMenuHtml(menuEl.dataset.role, menuEl.dataset.prefix || '', menuEl.dataset.exit || '../index.html');

    // Navbar breadcrumb: aktif menü başlığını navbar'a ekle
    const activeItem = menuEl.querySelector('.menu-item.active .menu-link div');
    const navbar = document.getElementById('layout-navbar');
    if (activeItem && navbar) {
      const titleEl = document.createElement('div');
      titleEl.className = 'kare-page-title d-none d-md-flex';
      titleEl.textContent = activeItem.textContent;
      // ms-auto'yu navbar-nav-right'tan kaldır; flex:1 breadcrumb'da
      const navRight = navbar.querySelector('.navbar-nav-right');
      if (navRight) navRight.classList.remove('ms-auto');
      const toggleBtn = navbar.querySelector('.layout-menu-toggle');
      if (toggleBtn) toggleBtn.after(titleEl);
    }
  }

  // User display name in navbar
  const user = kareCurrentUser();
  const un = document.getElementById('kare-user-name');
  if (un && user) un.textContent = user.name + ' ' + user.surname;
});
