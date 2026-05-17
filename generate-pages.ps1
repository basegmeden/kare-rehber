$base = 'c:\Users\mozek\Desktop\kare-rehber'
$d = -join ([char]100,[char]105,[char]118)
$A = 'https://demos.pixinvent.com/vuexy-html-admin-template/assets'
$enc = New-Object System.Text.UTF8Encoding $false

foreach ($folder in @('admin','koc','koordinator','veli','ogrenci')) {
  New-Item -ItemType Directory -Force -Path "$base\$folder" | Out-Null
}

function Head($title, $depth) {
  $p = if ($depth -eq 1) { '../' } else { '' }
  @"
<!DOCTYPE html>
<html lang="tr" class="layout-navbar-fixed layout-menu-fixed layout-compact" dir="ltr" data-bs-theme="dark" data-assets-path="$A/">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>$title | KARE Rehber</title>
<link rel="icon" type="image/x-icon" href="$A/img/favicon/favicon.ico">
<link rel="stylesheet" href="$A/vendor/fonts/iconify-icons.css">
<link rel="stylesheet" href="$A/vendor/libs/node-waves/node-waves.css">
<link rel="stylesheet" href="$A/vendor/css/core.css">
<link rel="stylesheet" href="$A/css/demo.css">
<link rel="stylesheet" href="$A/vendor/libs/perfect-scrollbar/perfect-scrollbar.css">
<link rel="stylesheet" href="${p}css/kare.css">
<script src="$A/vendor/js/helpers.js"></script>
<script src="$A/js/config.js"></script>
</head>
"@
}

function Shell($role, $title, $subtitle, $body, $depth) {
  $p = if ($depth -eq 1) { '../' } else { '' }
  $exit = "${p}index.html"
  $h = Head $title $depth
  $logo = @'
<svg width="32" height="32" viewBox="0 0 32 32" fill="none"><rect width="32" height="32" rx="8" fill="#4F46E5"/><rect x="15" y="9" width="2" height="16" fill="#C7D2FE"/><rect x="7" y="10" width="8" height="15" rx="2" fill="#FFFFFF"/><rect x="17" y="10" width="8" height="15" rx="2" fill="#FFFFFF"/><rect x="13" y="6" width="5" height="7" rx="1" fill="#34D399"/></svg>
'@
  @"

$h
<body>
<$d class="layout-wrapper layout-content-navbar">
  <$d class="layout-container">
    <aside id="layout-menu" class="layout-menu menu-vertical menu bg-menu-theme">
      <$d class="app-brand demo">
        <a href="${p}index.html" class="app-brand-link">
          <span class="app-brand-logo demo">$logo</span>
          <span class="app-brand-text demo menu-text fw-bold ms-2">KARE</span>
        </a>
        <a href="javascript:void(0)" class="layout-menu-toggle menu-link text-large ms-auto"><i class="ti tabler-x d-block d-xl-none"></i></a>
      </$d>
      <$d class="menu-inner-shadow"></$d>
      <ul class="menu-inner py-1" id="kare-menu-inner" data-role="$role" data-prefix="" data-exit="$exit"></ul>
    </aside>
    <$d class="layout-page">
      <nav class="layout-navbar container-xxl navbar navbar-expand-xl navbar-detached align-items-center bg-navbar-theme" id="layout-navbar">
        <$d class="layout-menu-toggle navbar-nav align-items-xl-center me-3 me-xl-0 d-xl-none">
          <a class="nav-item nav-link px-0 me-xl-4" href="javascript:void(0)"><i class="ti tabler-menu-2"></i></a>
        </$d>
        <$d class="navbar-nav-right d-flex align-items-center ms-auto">
          <ul class="navbar-nav flex-row align-items-center">
            <li class="nav-item navbar-dropdown dropdown-user dropdown">
              <a class="nav-link dropdown-toggle hide-arrow" href="javascript:void(0);" data-bs-toggle="dropdown">
                <span class="avatar avatar-online"><span class="avatar-initial rounded-circle bg-label-primary">K</span></span>
              </a>
              <ul class="dropdown-menu dropdown-menu-end">
                <li><a class="dropdown-item" href="#"><$d class="d-flex"><$d class="flex-grow-1"><h6 class="mb-0" id="kare-user-name">Kullanıcı</h6><small class="text-muted">$role</small></$d></$d></a></li>
                <li><hr class="dropdown-divider"></li>
                <li><a class="dropdown-item" href="$exit"><i class="ti tabler-logout me-2"></i>Çıkış</a></li>
              </ul>
            </li>
          </ul>
        </$d>
      </nav>
      <$d class="content-wrapper">
        <$d class="container-xxl flex-grow-1 container-p-y">
          <$d class="row"><$d class="col-12 mb-4">
            <h4 class="mb-1">$title</h4>
            <p class="text-muted mb-0">$subtitle</p>
          </$d></$d>
          $body
        </$d>
      </$d>
    </$d>
  </$d>
</$d>
<script src="$A/vendor/libs/jquery/jquery.js"></script>
<script src="$A/vendor/libs/popper/popper.js"></script>
<script src="$A/vendor/js/bootstrap.js"></script>
<script src="$A/vendor/libs/node-waves/node-waves.js"></script>
<script src="$A/vendor/libs/perfect-scrollbar/perfect-scrollbar.js"></script>
<script src="$A/vendor/libs/hammer/hammer.js"></script>
<script src="$A/vendor/js/menu.js"></script>
<script src="$A/js/main.js"></script>
<script src="${p}js/kare.js"></script>
</body>
</html>
"@
}

# INDEX - Login cover
$index = @"
<!DOCTYPE html>
<html lang="tr" class="layout-wide customizer-hide" dir="ltr" data-bs-theme="dark" data-assets-path="$A/">
<head>
<meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Giriş | KARE Rehber</title>
<link rel="stylesheet" href="$A/vendor/fonts/iconify-icons.css">
<link rel="stylesheet" href="$A/vendor/css/core.css">
<link rel="stylesheet" href="$A/css/demo.css">
<link rel="stylesheet" href="$A/vendor/css/pages/page-auth.css">
<script src="$A/vendor/js/helpers.js"></script>
<script src="$A/js/config.js"></script>
</head>
<body>
<$d class="authentication-wrapper authentication-cover">
  <a href="index.html" class="app-brand auth-cover-brand">
    <span class="app-brand-logo demo"><span class="text-primary fw-bold fs-2">K</span></span>
    <span class="app-brand-text demo text-heading fw-bold">KARE Rehber</span>
  </a>
  <$d class="authentication-inner row m-0">
    <$d class="d-none d-xl-flex col-xl-8 p-0">
      <$d class="auth-cover-bg d-flex justify-content-center align-items-center p-5">
        <$d class="text-white">
          <h2 class="text-white mb-3">Öğrenci gelişimini şeffaf takip edin</h2>
          <p class="mb-0 opacity-75">Koç, veli, koordinatör ve admin — haftalık görüşmeler, SMS ve raporlama tek platformda.</p>
        </$d>
      </$d>
    </$d>
    <$d class="d-flex col-12 col-xl-4 align-items-center authentication-bg p-sm-12 p-6">
      <$d class="w-px-400 mx-auto">
        <h4 class="mb-1">Hoş geldiniz 👋</h4>
        <p class="mb-6">Rolünüzü seçerek demo panele giriş yapın</p>
        <form id="login-form">
          <input type="hidden" id="selected-role" value="admin">
          <$d class="mb-4">
            <label class="form-label">Rol</label>
            <$d class="d-grid gap-2">
              <button type="button" class="btn btn-outline-primary role-btn selected" data-role="admin">Admin</button>
              <button type="button" class="btn btn-outline-primary role-btn" data-role="koc">Koç</button>
              <button type="button" class="btn btn-outline-primary role-btn" data-role="koordinator">Koordinatör</button>
              <button type="button" class="btn btn-outline-primary role-btn" data-role="veli">Veli</button>
              <button type="button" class="btn btn-outline-primary role-btn" data-role="ogrenci">Öğrenci</button>
            </$d>
          </$d>
          <$d class="mb-4"><label class="form-label">Kullanıcı adı</label><input type="text" class="form-control" value="demo" readonly></$d>
          <$d class="mb-4"><label class="form-label">Şifre</label><input type="password" class="form-control" value="demo" readonly></$d>
          <button class="btn btn-primary d-grid w-100 mb-4" type="submit">Giriş Yap</button>
        </form>
        <p class="text-center mb-2"><span class="text-muted">Kayıt formları:</span></p>
        <p class="text-center"><a href="ogrenci-kayit.html">Öğrenci ön kayıt</a> · <a href="koc-basvuru.html">Koç başvuru</a></p>
      </$d>
    </$d>
  </$d>
</$d>
<script src="$A/vendor/libs/jquery/jquery.js"></script>
<script src="$A/vendor/js/bootstrap.js"></script>
<script src="js/kare.js"></script>
<script>
document.querySelectorAll('.role-btn').forEach(b=>b.onclick=()=>{document.querySelectorAll('.role-btn').forEach(x=>x.classList.remove('selected','btn-primary'));x.classList.add('btn-outline-primary');b.classList.add('selected','btn-primary');b.classList.remove('btn-outline-primary');document.getElementById('selected-role').value=b.dataset.role});
document.getElementById('login-form').onsubmit=e=>{e.preventDefault();kareLogin(document.getElementById('selected-role').value)};
</script>
</body>
</html>
"@

[System.IO.File]::WriteAllText("$base\index.html", $index, $enc)

# Admin dashboard content
$dashBody = @"
<$d class="row g-6 mb-6">
  <$d class="col-sm-6 col-xl-3"><$d class="card"><$d class="card-body"><$d class="d-flex align-items-start justify-content-between"><$d><span class="text-heading">Toplam Öğrenci</span><h4 class="mb-0 mt-2">1.248</h4></$d><span class="badge bg-label-primary rounded p-2"><i class="ti tabler-users"></i></span></$d></$d></$d></$d>
  <$d class="col-sm-6 col-xl-3"><$d class="card"><$d class="card-body"><span class="text-heading">Aktif Koç</span><h4 class="mb-0 mt-2">86</h4></$d></$d></$d>
  <$d class="col-sm-6 col-xl-3"><$d class="card"><$d class="card-body"><span class="text-heading">Eksik Görüşme</span><h4 class="mb-0 mt-2 text-warning">23</h4></$d></$d></$d>
  <$d class="col-sm-6 col-xl-3"><$d class="card"><$d class="card-body"><span class="text-heading">Bekleyen Mesaj</span><h4 class="mb-0 mt-2">7</h4></$d></$d></$d>
</$d>
<$d class="row g-6">
  <$d class="col-xl-6"><$d class="card h-100"><$d class="card-header"><h5 class="mb-0">Son Kayıtlar</h5></$d><$d class="table-responsive text-nowrap"><table class="table table-hover"><thead><tr><th>Ad</th><th>İl</th><th>Durum</th></tr></thead><tbody>
  <tr><td>Ali Kaya</td><td>İstanbul</td><td><span class="badge bg-label-warning">Ön kayıt</span></td></tr>
  <tr><td>Zeynep Arslan</td><td>Ankara</td><td><span class="badge bg-label-success">Kesin kayıt</span></td></tr>
  </tbody></table></$d></$d></$d>
  <$d class="col-xl-6"><$d class="card h-100"><$d class="card-header"><h5 class="mb-0">Geciken Koçlar</h5></$d><$d class="table-responsive"><table class="table"><thead><tr><th>Koç</th><th>Eksik</th><th></th></tr></thead><tbody>
  <tr><td>Ayşe Yılmaz</td><td>3 öğrenci</td><td><button class="btn btn-sm btn-outline-primary" data-toast="SMS gönderildi">SMS</button></td></tr>
  </tbody></table></$d></$d></$d>
</$d>
"@

$pages = @{
  "admin\dashboard.html" = @{ role='admin'; title='Yönetim Paneli'; sub='Sistem özeti'; body=$dashBody; depth=1 }
}

[System.IO.File]::WriteAllText("$base\admin\dashboard.html", (Shell 'admin' 'Yönetim Paneli' 'Sistem özeti' $dashBody 1), $enc)
Write-Host "Generated dashboard"
