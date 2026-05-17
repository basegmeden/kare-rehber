. "$PSScriptRoot\generate-pages.ps1"

$kul = @"
<$d class="card">
  <$d class="card-header d-flex justify-content-between flex-wrap gap-2">
    <h5 class="mb-0">Kullanıcı listesi</h5>
    <button class="btn btn-primary btn-sm" data-toast="Kullanıcı eklendi"><i class="ti tabler-plus me-1"></i> Ekle</button>
  </$d>
  <$d class="card-body">
    <$d class="row mb-4 g-3">
      <$d class="col-md-4"><input class="form-control" placeholder="Ara..."></$d>
      <$d class="col-md-3"><select class="form-select"><option>Tüm roller</option><option>Öğrenci</option><option>Koç</option><option>Koordinatör</option></select></$d>
    </$d>
    <$d class="table-responsive"><table class="table table-hover"><thead><tr><th>Ad Soyad</th><th>Rol</th><th>Telefon</th><th>İl</th><th>Durum</th><th></th></tr></thead><tbody>
    <tr><td>Ali Kaya</td><td><span class="badge bg-label-info">Öğrenci</span></td><td>0532 ***</td><td>İstanbul</td><td><span class="badge bg-label-success">Aktif</span></td><td><button class="btn btn-sm btn-icon btn-text-secondary"><i class="ti tabler-edit"></i></button></td></tr>
    <tr><td>Ayşe Yılmaz</td><td><span class="badge bg-label-primary">Koç</span></td><td>0544 ***</td><td>Ankara</td><td><span class="badge bg-label-success">Aktif</span></td><td><button class="btn btn-sm btn-icon btn-text-secondary"><i class="ti tabler-edit"></i></button></td></tr>
    <tr><td>Fatma Vakfı</td><td><span class="badge bg-label-warning">Koordinatör</span></td><td>0555 ***</td><td>İzmir</td><td><span class="badge bg-label-secondary">Beklemede</span></td><td><button class="btn btn-sm btn-primary" data-toast="Giriş bilgileri SMS ile gönderildi">Giriş Bilgileri Gönder</button></td></tr>
    </tbody></table></$d>
  </$d>
</$d>
"@

$ek = @"
<$d class="alert alert-primary" role="alert">Kesin kaydı tamamlanan öğrenciler ilgili koçlarla eşleştirilir.</$d>
<$d class="row g-6">
  <$d class="col-lg-6"><$d class="card"><$d class="card-header"><h5 class="mb-0">Eşleştirilmemiş öğrenciler</h5></$d><$d class="card-body">
    <$d class="row g-2 mb-3"><$d class="col-6"><select class="form-select form-select-sm"><option>Tüm iller</option><option>İstanbul</option></select></$d><$d class="col-6"><input class="form-control form-control-sm" placeholder="Ara"></$d></$d>
    <table class="table table-sm"><thead><tr><th><input type="checkbox" class="form-check-input"></th><th>Öğrenci</th><th>İl</th></tr></thead><tbody>
    <tr><td><input type="checkbox" class="form-check-input"></td><td>Ali Kaya</td><td>İstanbul</td></tr>
  </tbody></table></$d></$d></$d>
  <$d class="col-lg-6"><$d class="card"><$d class="card-header"><h5 class="mb-0">Koç havuzu</h5></$d><$d class="card-body">
    <select class="form-select mb-3"><option>Ayşe Yılmaz</option><option>Mehmet Koç</option></select>
    <button class="btn btn-primary w-100" data-toast="Toplu eşleştirme tamamlandı">Toplu Eşleştir</button>
  </$d></$d></$d>
</$d>
"@

$eko = $ek -replace 'koçlarla','koordinatörlerle' -replace 'Koç havuzu','Koordinatör / Vakıf' -replace 'Koç havuzu','Koordinatör seçin'

$gor = @"
<$d class="card"><$d class="card-header d-flex flex-wrap gap-2 justify-content-between">
  <h5 class="mb-0">Görüşme kayıtları</h5>
  <$d class="d-flex gap-2"><select class="form-select form-select-sm" style="width:auto"><option>3. Hafta (11-17 Mayıs)</option></select><select class="form-select form-select-sm" style="width:auto"><option>Onay bekliyor</option></select></$d>
</$d><$d class="table-responsive"><table class="table"><thead><tr><th>Öğrenci</th><th>Koç</th><th>Hafta</th><th>Durum</th><th></th></tr></thead><tbody>
<tr><td>Ali Kaya</td><td>Ayşe Yılmaz</td><td>3. Hafta</td><td><span class="badge bg-label-warning">Onay bekliyor</span></td><td><button class="btn btn-sm btn-success" data-toast="Onaylandı">Onayla</button> <button class="btn btn-sm btn-label-secondary">Düzenle</button></td></tr>
<tr><td>Zeynep A.</td><td>Ayşe Yılmaz</td><td>3. Hafta</td><td><span class="badge bg-label-success">Onaylandı</span></td><td><button class="btn btn-sm btn-label-secondary">Log</button></td></tr>
</tbody></table></$d><$d class="card-footer text-muted small">Admin düzenlemeleri LOG kayıtlarında tutulur.</$d></$d>
"@

$uyar = @"
<$d class="card"><$d class="card-header d-flex justify-content-between"><h5 class="mb-0">Görüşmesini geciktiren koçlar</h5><button class="btn btn-primary btn-sm" data-toast="Toplu SMS gönderildi"><i class="ti tabler-send"></i> Toplu SMS</button></$d>
<$d class="table-responsive"><table class="table table-hover"><thead><tr><th>Koç</th><th>Eksik öğrenci</th><th>Durum</th><th></th></tr></thead><tbody>
<tr><td>Ayşe Yılmaz</td><td>3</td><td><span class="badge bg-label-danger">Pazartesi gecikmesi</span></td><td><button class="btn btn-sm btn-outline-primary" data-toast="SMS gönderildi">SMS</button></td></tr>
</tbody></table></$d></$d>
"@

$hafta = @"
<$d class="card"><$d class="card-header d-flex justify-content-between"><h5 class="mb-0">Hafta tanımları</h5><button class="btn btn-primary btn-sm" data-toast="Hafta eklendi"><i class="ti tabler-plus"></i> Hafta Ekle</button></$d>
<table class="table"><thead><tr><th>Hafta</th><th>Başlangıç</th><th>Bitiş</th><th>Durum</th><th></th></tr></thead><tbody>
<tr><td>3. Hafta</td><td>11 May 2026</td><td>17 May 2026</td><td><span class="badge bg-label-success">Aktif</span></td><td>—</td></tr>
<tr><td>2. Hafta</td><td>4 May</td><td>10 May</td><td><span class="badge bg-label-info">Önceki</span></td><td>—</td></tr>
<tr><td>1. Hafta</td><td>27 Nis</td><td>3 May</td><td><span class="badge bg-label-secondary">Kapalı</span></td><td><button class="btn btn-sm btn-label-primary" data-toast="Hafta açıldı">Admin için aç</button></td></tr>
</tbody></table></$d>
"@

$mesaj = @"
<$d class="row g-6"><$d class="col-md-5"><$d class="card h-100"><$d class="card-header"><h5 class="mb-0">Gelen mesajlar</h5></$d><$d class="list-group list-group-flush">
  <a href="#" class="list-group-item list-group-item-action active"><$d class="d-flex w-100 justify-content-between"><strong>Veli — Mehmet D.</strong><small>Bugün</small></$d><p class="mb-0 small">Oğlumun ders durumu hakkında bilgi alabilir miyim?</p></a>
  <a href="#" class="list-group-item list-group-item-action"><strong>Öğrenci — Ali K.</strong><p class="mb-0 small text-muted">Geri bildirim mesajı</p></a>
</$d></$d><$d class="col-md-7"><$d class="card h-100"><$d class="card-header"><h5 class="mb-0">Yanıt</h5></$d><$d class="card-body"><textarea class="form-control mb-3" rows="6" placeholder="Cevap yazın..."></textarea><button class="btn btn-primary" data-toast="Mesaj gönderildi">Gönder</button></$d></$d></$d>
"@

$sms = @"
<$d class="row g-6"><$d class="col-md-6"><$d class="card"><$d class="card-header"><h5 class="mb-0">SMS gönder</h5></$d><$d class="card-body">
  <select class="form-select mb-3"><option>Tüm geciken koçlar</option><option>Seçili kullanıcılar</option></select>
  <textarea class="form-control mb-3" rows="4">Lütfen haftalık görüşme notlarınızı giriniz.</textarea>
  <button class="btn btn-primary" data-toast="SMS gönderildi">Gönder</button>
</$d></$d><$d class="col-md-6"><$d class="card"><$d class="card-header"><h5 class="mb-0">Giriş bilgileri SMS</h5></$d><$d class="card-body"><p class="text-muted">Onaylanan öğrenci/koçlara link, kullanıcı adı ve şifre gönderilir.</p><button class="btn btn-label-primary" data-toast="Giriş bilgileri gönderildi">Giriş Bilgileri Gönder</button></$d></$d></$d>
"@

$rapor = @"
<$d class="row g-6 mb-6">
  <$d class="col-md-4"><$d class="card"><$d class="card-body"><span class="text-muted">Görüşme (bu ay)</span><h3 class="mb-0 mt-1">4.820</h3></$d></$d>
  <$d class="col-md-4"><$d class="card"><$d class="card-body"><span class="text-muted">Aktif öğrenci</span><h3 class="mb-0 mt-1">1.102</h3></$d></$d>
  <$d class="col-md-4"><$d class="card"><$d class="card-body"><span class="text-muted">Pasif öğrenci</span><h3 class="mb-0 mt-1 text-warning">146</h3></$d></$d>
</$d>
<$d class="card"><$d class="card-header"><h5 class="mb-0">İl bazlı dağılım</h5></$d><table class="table"><thead><tr><th>İl</th><th>Öğrenci</th><th>Görüşme oranı</th></tr></thead><tbody>
<tr><td>İstanbul</td><td>412</td><td><$d class="progress" style="height:6px"><$d class="progress-bar" style="width:94%"></$d></$d> 94%</td></tr>
<tr><td>Ankara</td><td>298</td><td>91%</td></tr></tbody></table></$d>
"@

$log = @"
<$d class="card"><table class="table"><thead><tr><th>Tarih</th><th>Kullanıcı</th><th>İşlem</th><th>Detay</th></tr></thead><tbody>
<tr><td>16.05.2026 10:32</td><td>Admin</td><td><span class="badge bg-label-warning">Düzenleme</span></td><td>Ali Kaya — 3. hafta notu güncellendi</td></tr>
<tr><td>15.05.2026 18:01</td><td>Admin</td><td><span class="badge bg-label-success">Onay</span></td><td>Veli paneline açıldı</td></tr>
</tbody></table></$d>
"@

$adminPages = @{
  'kullanicilar.html' = @('Kullanıcılar','Tüm kullanıcıları yönetin',$kul)
  'eslestirme-koc.html' = @('Öğrenci–Koç Eşleştirme','İl bazlı filtreleme ve toplu eşleştirme',$ek)
  'eslestirme-koordinator.html' = @('Öğrenci–Koordinatör','Vakıf eşleştirmesi',$eko)
  'gorusmeler.html' = @('Görüşme Kayıtları','İncele, onayla, güncelle',$gor)
  'koc-uyarilar.html' = @('Koç Uyarıları','Zamanında görüşmeyen koçlar',$uyar)
  'haftalar.html' = @('Hafta Tanımları','Admin tarih listesi',$hafta)
  'mesajlar.html' = @('Mesajlar','Veli ve öğrenci mesajları',$mesaj)
  'sms.html' = @('SMS Gönderimi','Toplu ve bireysel',$sms)
  'raporlar.html' = @('Raporlar','İstatistik ve performans',$rapor)
  'loglar.html' = @('Log Kayıtları','Tüm değişiklikler',$log)
}

foreach ($kv in $adminPages.GetEnumerator()) {
  $path = Join-Path $base "admin\$($kv.Key)"
  [System.IO.File]::WriteAllText($path, (Shell 'admin' $kv.Value[0] $kv.Value[1] $kv.Value[2] 1), $enc)
}

# KOC
$kdash = @"
<$d class="alert alert-info">Aktif hafta: <strong>3. Hafta (11-17 Mayıs 2026)</strong>. En fazla bir önceki haftada işlem yapabilirsiniz.</$d>
<$d class="row g-6"><$d class="col-md-4"><$d class="card"><$d class="card-body"><span class="text-muted">Öğrencilerim</span><h3 class="mb-0">12</h3></$d></$d>
<$d class="col-md-4"><$d class="card"><$d class="card-body"><span class="text-muted">Bu hafta eksik</span><h3 class="mb-0 text-warning">2</h3></$d></$d>
<$d class="col-md-4"><$d class="card"><$d class="card-body"><span class="text-muted">Tamamlanan</span><h3 class="mb-0 text-success">10</h3></$d></$d></$d>
"@

$kogr = @"
<$d class="card"><table class="table table-hover"><thead><tr><th>Öğrenci</th><th>İl</th><th>Bu hafta</th><th></th></tr></thead><tbody>
<tr><td>Ali Kaya</td><td>İstanbul</td><td><span class="badge bg-label-warning">Girilmedi</span></td><td><a href="gorusme.html" class="btn btn-sm btn-primary">Görüşme gir</a></td></tr>
<tr><td>Zeynep A.</td><td>Ankara</td><td><span class="badge bg-label-success">Tamam</span></td><td><a href="gorusme.html" class="btn btn-sm btn-label-secondary">Görüntüle</a></td></tr>
</tbody></table></$d>
"@

$kgor = @"
<$d class="btn-group mb-4" role="group">
  <button type="button" class="btn btn-primary week-chip active">3. Hafta (11-17 Mayıs)</button>
  <button type="button" class="btn btn-label-secondary week-chip">2. Hafta (4-10 Mayıs)</button>
  <button type="button" class="btn btn-label-secondary disabled">1. Hafta (kapalı)</button>
</$d>
<$d class="card"><$d class="card-body">
  <$d class="mb-4"><label class="form-label">Öğrenci</label><select class="form-select"><option>Ali Kaya</option></select></$d>
  <$d class="mb-4"><label class="form-label">Genel değerlendirme (1-5)</label><$d class="btn-group">
    <button type="button" class="btn btn-outline-primary rating-btn">1</button><button type="button" class="btn btn-outline-primary rating-btn">2</button>
    <button type="button" class="btn btn-outline-primary rating-btn active">3</button><button type="button" class="btn btn-outline-primary rating-btn">4</button><button type="button" class="btn btn-outline-primary rating-btn">5</button>
  </$d></$d>
  <$d class="mb-4"><label class="form-label">Görüşme notu</label><textarea class="form-control" rows="5"></textarea></$d>
  <button class="btn btn-primary" data-toast="Görüşme kaydedildi">Kaydet</button>
</$d></$d>
"@

$kmes = @"
<$d class="card"><$d class="card-body"><textarea class="form-control mb-3" rows="4" placeholder="Admin veya koordinatöre mesaj..."></textarea><button class="btn btn-primary" data-toast="Mesaj gönderildi">Gönder</button></$d></$d>
"@

@('dashboard.html','Yönetim Paneli','Özet',$kdash),
@('ogrenciler.html','Öğrencilerim','Tüm öğrenciler',$kogr),
@('gorusme.html','Haftalık Görüşme','Haftada 1 — hafta sonu beklenir',$kgor),
@('mesajlar.html','Mesajlar','İletişim',$kmes) | ForEach-Object {
  [System.IO.File]::WriteAllText((Join-Path $base "koc\$($_[0])"), (Shell 'koc' $_[1] $_[2] $_[3] 1), $enc)
}

# KOORDINATOR
@('dashboard.html','Koordinatör Paneli','Vakıf öğrenci takibi',$kdash),
@('ogrenciler.html','Öğrenciler','Eşleşen öğrenciler',$kogr),
@('gorusmeler.html','Rehber Görüşmeleri','Admin onayı olmadan görüntüleme',$gor) | ForEach-Object {
  [System.IO.File]::WriteAllText((Join-Path $base "koordinator\$($_[0])"), (Shell 'koordinator' $_[1] $_[2] $_[3] 1), $enc)
}

# VELI
$vd = @"
<$d class="alert alert-primary">Sadece admin tarafından onaylanan haftanın tamamı görüntülenir.</$d>
<$d class="row g-6 mb-4"><$d class="col-md-6"><$d class="card"><$d class="card-body"><span class="text-muted">Öğrenci</span><h5 class="mb-0">Ali Kaya</h5></$d></$d>
<$d class="col-md-6"><$d class="card"><$d class="card-body"><span class="text-muted">Koç</span><h5 class="mb-0">Ayşe Yılmaz</h5></$d></$d></$d>
"@
$vg = @"
<$d class="card mb-4"><$d class="card-header"><h5 class="mb-0">3. Hafta (11-17 Mayıs) <span class="badge bg-label-success">Onaylandı</span></h5></$d><$d class="card-body"><p><strong>Değerlendirme:</strong> 4/5</p><p class="mb-0">Ali bu hafta hedeflerine uygun ilerledi.</p></$d></$d>
"@
@('dashboard.html','Veli Paneli','Onaylı bilgiler',$vd),
@('gorusmeler.html','Onaylı Görüşmeler','Haftalık raporlar',$vg),
@('mesajlar.html','Mesajlar','Admin / koordinatör',$kmes) | ForEach-Object {
  [System.IO.File]::WriteAllText((Join-Path $base "veli\$($_[0])"), (Shell 'veli' $_[1] $_[2] $_[3] 1), $enc)
}

# OGRENCI
$od = @"
<$d class="alert alert-warning">Görüşme notlarınızı göremezsiniz.</$d>
<$d class="card"><$d class="card-header"><h5 class="mb-0">Geri bildirim</h5></$d><$d class="card-body"><textarea class="form-control mb-3" rows="4"></textarea><button class="btn btn-primary" data-toast="Geri bildirim admin listesine eklendi">Gönder</button></$d></$d>
"@
@('dashboard.html','Öğrenci Paneli','Sınırlı erişim',$od),
@('mesajlar.html','Mesajlar','Admin / koordinatör',$kmes) | ForEach-Object {
  [System.IO.File]::WriteAllText((Join-Path $base "ogrenci\$($_[0])"), (Shell 'ogrenci' $_[1] $_[2] $_[3] 1), $enc)
}

# Public forms
$pubHead = @"
<!DOCTYPE html>
<html lang="tr" data-bs-theme="dark" data-assets-path="$A/">
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Öğrenci Kayıt | KARE</title>
<link rel="stylesheet" href="$A/vendor/css/core.css"><link rel="stylesheet" href="$A/css/demo.css">
</head><body class="bg-body">
<nav class="navbar navbar-expand-lg bg-navbar-theme"><$d class="container-xxl">
  <a class="navbar-brand fw-bold text-primary" href="index.html">KARE Rehber</a>
  <a href="index.html" class="btn btn-sm btn-outline-primary">Giriş</a>
</$d></nav>
<$d class="container-xxl py-6"><$d class="row justify-content-center"><$d class="col-lg-8">
"@

$pubForm = @"
<$d class="card"><$d class="card-header"><h4 class="mb-0">Öğrenci Ön Kayıt</h4><small class="text-muted">kare.ulued.org/ogrenci-kayit-formu</small></$d>
<$d class="card-body">
<$d class="alert alert-info">Eksiksiz form → SMS bilgilendirme → telefon görüşmesi → kesin kayıt</$d>
<form>
<$d class="row g-4"><$d class="col-md-6"><label class="form-label">Ad *</label><input class="form-control" required></$d>
<$d class="col-md-6"><label class="form-label">Soyad *</label><input class="form-control" required></$d>
<$d class="col-md-6"><label class="form-label">Telefon *</label><input type="tel" class="form-control" required></$d>
<$d class="col-md-6"><label class="form-label">Doğum tarihi *</label><input type="date" class="form-control" required></$d>
<$d class="col-12"><label class="form-label">İl *</label><select class="form-select"><option>İstanbul</option><option>Ankara</option></select></$d>
<$d class="col-12"><button type="button" class="btn btn-primary" onclick="kareToast('Ön kayıt alındı')">Başvuruyu Gönder</button></$d>
</$d></form></$d></$d>
"@

$pubEnd = "</$d></$d></$d><script src=`"$A/vendor/js/bootstrap.js`"></script><script src=`"js/kare.js`"></script></body></html>"
[System.IO.File]::WriteAllText("$base\ogrenci-kayit.html", $pubHead + $pubForm + $pubEnd, $enc)
$kocForm = $pubForm -replace 'Öğrenci Ön Kayıt','Koç Ön Başvuru' -replace 'ogrenci-kayit','koc-on-basvuru' -replace 'Öğrenci','Koç' -replace 'Doğum tarihi','Deneyim (yıl)'
[System.IO.File]::WriteAllText("$base\koc-basvuru.html", ($pubHead -replace 'Öğrenci Kayıt','Koç Başvuru') + $kocForm + $pubEnd, $enc)

# kare.css
$css = ".role-btn.selected { pointer-events: none; } .week-chip.active { box-shadow: 0 0 0 2px var(--bs-primary); }"
[System.IO.File]::WriteAllText("$base\css\kare.css", $css, $enc)

# Fix index login script
$idx = Get-Content "$base\index.html" -Raw -Encoding UTF8
$idx = $idx -replace "document\.querySelectorAll\('\.role-btn'\)\.forEach\(b=>b\.onclick=\(\)=>\{document\.querySelectorAll\('\.role-btn'\)\.forEach\(x=>x\.classList\.remove\('selected','btn-primary'\)\);x\.classList\.add\('btn-outline-primary'\);b\.classList\.add\('selected','btn-primary'\);b\.classList\.remove\('btn-outline-primary'\);document\.getElementById\('selected-role'\)\.value=b\.dataset\.role\}\);","document.querySelectorAll('.role-btn').forEach(b=>b.onclick=()=>{document.querySelectorAll('.role-btn').forEach(x=>{x.classList.remove('selected','btn-primary');x.classList.add('btn-outline-primary')});b.classList.add('selected','btn-primary');b.classList.remove('btn-outline-primary');document.getElementById('selected-role').value=b.dataset.role});"
[System.IO.File]::WriteAllText("$base\index.html", $idx, $enc)

Write-Host "All pages generated:" (Get-ChildItem $base -Recurse -Filter *.html).Count
