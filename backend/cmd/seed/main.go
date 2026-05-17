// Seed: populates the database with realistic Turkish test data.
// Usage: go run ./cmd/seed/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"kare-rehber/backend/internal/config"
	"kare-rehber/backend/internal/database"
	"kare-rehber/backend/internal/models"
	"kare-rehber/backend/internal/services"

	"gorm.io/datatypes"
)

func ptr[T any](v T) *T { return &v }

func hash(pwd string) string {
	h, err := services.HashPassword(pwd)
	if err != nil {
		log.Fatal(err)
	}
	return h
}

func dob(year, month, day int) *time.Time {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return &t
}

func ago(days int) time.Time {
	return time.Now().AddDate(0, 0, -days)
}

func catJSON(ratings [6]int, notes [6]string) datatypes.JSON {
	names := []string{
		"Çalışma Stili",
		"Motivasyon Durumu",
		"Hedefler ve Başarı Durumu",
		"Gelişim Serüveni",
		"Ders Bazlı Gelişim",
		"Hizmet Şuuru (vakıf ile irtibat, muhabbet, gidiş geliş)",
	}
	cats := make([]map[string]interface{}, 6)
	for i := range cats {
		cats[i] = map[string]interface{}{"name": names[i], "rating": ratings[i], "notes": notes[i]}
	}
	b, _ := json.Marshal(map[string]interface{}{"categories": cats})
	return datatypes.JSON(b)
}

func main() {
	config.Load()
	database.Connect()

	// ── Weeks ─────────────────────────────────────────────────────────────────
	weeks := []models.Week{
		{
			WeekNumber: 1,
			Label:      "1. Hafta (28 Nisan – 4 Mayıs)",
			StartDate:  time.Date(2026, 4, 28, 0, 0, 0, 0, time.UTC),
			EndDate:    time.Date(2026, 5, 4, 23, 59, 59, 0, time.UTC),
			IsActive:   false, IsLocked: true, CreatedBy: 1,
		},
		{
			WeekNumber: 2,
			Label:      "2. Hafta (5 – 11 Mayıs)",
			StartDate:  time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC),
			EndDate:    time.Date(2026, 5, 11, 23, 59, 59, 0, time.UTC),
			IsActive:   false, IsLocked: true, CreatedBy: 1,
		},
		{
			WeekNumber: 3,
			Label:      "3. Hafta (12 – 18 Mayıs)",
			StartDate:  time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC),
			EndDate:    time.Date(2026, 5, 18, 23, 59, 59, 0, time.UTC),
			IsActive:   true, IsLocked: false, CreatedBy: 1,
		},
		{
			WeekNumber: 4,
			Label:      "4. Hafta (19 – 25 Mayıs)",
			StartDate:  time.Date(2026, 5, 19, 0, 0, 0, 0, time.UTC),
			EndDate:    time.Date(2026, 5, 25, 23, 59, 59, 0, time.UTC),
			IsActive:   false, IsLocked: false, CreatedBy: 1,
		},
	}
	for i := range weeks {
		if err := database.DB.Create(&weeks[i]).Error; err != nil {
			fmt.Printf("Hafta zaten mevcut (hafta %d)\n", i+1)
		}
	}
	fmt.Println("✓ Haftalar oluşturuldu")

	// ── Coaches ───────────────────────────────────────────────────────────────
	coachDefs := []struct {
		name, surname, phone, city, exp, username string
	}{
		{"Ahmet", "Yılmaz", "05321000001", "İstanbul", "3 yıl yüz yüze öğrenci takibi, 2 yıl online koçluk", "koc1"},
		{"Zeynep", "Kaya", "05321000002", "Ankara", "5 yıl öğretmenlik, 2 yıl akademik koçluk", "koc2"},
		{"Mehmet", "Demir", "05321000003", "İzmir", "Eğitim psikoloğu, 4 yıl sınav koçluğu deneyimi", "koc3"},
		{"Fatma", "Çelik", "05321000004", "Bursa", "TYT/AYT uzmanı, 6 yıl bireysel ders deneyimi", "koc4"},
		{"Hasan", "Şahin", "05321000005", "Konya", "YKS koçu, 2 yıl deneyim", "koc5"},
	}
	var coaches []models.Coach
	var coachUserIDs []uint
	for _, d := range coachDefs {
		u := models.User{
			Name: d.name, Surname: d.surname, Phone: d.phone,
			City: d.city, Role: models.RoleCoach,
			Username: d.username, PasswordHash: hash("koc123"),
			Status: models.UserStatusActive,
		}
		if err := database.DB.Create(&u).Error; err != nil {
			fmt.Printf("Koç zaten mevcut: %s\n", d.username)
			continue
		}
		c := models.Coach{UserID: u.ID, Experience: d.exp, PoolStatus: models.PoolStatusApproved}
		database.DB.Create(&c)
		coaches = append(coaches, c)
		coachUserIDs = append(coachUserIDs, u.ID)
	}
	fmt.Printf("✓ %d koç oluşturuldu\n", len(coaches))

	// ── Coordinators ──────────────────────────────────────────────────────────
	coordDefs := []struct {
		name, surname, phone, city, org, username string
	}{
		{"Mustafa", "Arslan", "05331000001", "İstanbul", "ULUED İstanbul Şubesi", "koordinator1"},
		{"Esra", "Öztürk", "05331000002", "Ankara", "ULUED Ankara Şubesi", "koordinator2"},
	}
	var coords []models.Coordinator
	for _, d := range coordDefs {
		u := models.User{
			Name: d.name, Surname: d.surname, Phone: d.phone,
			City: d.city, Role: models.RoleCoordinator,
			Username: d.username, PasswordHash: hash("koor123"),
			Status: models.UserStatusActive,
		}
		if err := database.DB.Create(&u).Error; err != nil {
			fmt.Printf("Koordinatör zaten mevcut: %s\n", d.username)
			continue
		}
		coord := models.Coordinator{UserID: u.ID, Organization: d.org}
		database.DB.Create(&coord)
		coords = append(coords, coord)
	}
	fmt.Printf("✓ %d koordinatör oluşturuldu\n", len(coords))

	// ── Students ──────────────────────────────────────────────────────────────
	type studentDef struct {
		name, surname, phone, city, school, grade, username string
		dob                                                   *time.Time
		status                                                models.RegistrationStatus
	}
	sDefs := []studentDef{
		// Coach 1 – İstanbul (3 öğrenci)
		{"Ali", "Arslan", "05411000001", "İstanbul", "Kadıköy Anadolu Lisesi", "11. Sınıf", "ogrenci1", dob(2008, 3, 15), models.RegistrationConfirmed},
		{"Büşra", "Doğan", "05411000002", "İstanbul", "Ümraniye Fen Lisesi", "10. Sınıf", "ogrenci2", dob(2009, 7, 22), models.RegistrationConfirmed},
		{"Yusuf", "Öztürk", "05411000003", "İstanbul", "Beşiktaş Anadolu Lisesi", "12. Sınıf", "ogrenci3", dob(2007, 11, 8), models.RegistrationConfirmed},
		// Coach 2 – Ankara (3 öğrenci)
		{"Merve", "Yıldız", "05411000004", "Ankara", "Çankaya Lisesi", "10. Sınıf", "ogrenci4", dob(2009, 5, 30), models.RegistrationConfirmed},
		{"İbrahim", "Koç", "05411000005", "Ankara", "Ankara Fen Lisesi", "11. Sınıf", "ogrenci5", dob(2008, 1, 12), models.RegistrationConfirmed},
		{"Seda", "Erdoğan", "05411000006", "Ankara", "Keçiören Lisesi", "9. Sınıf", "ogrenci6", dob(2010, 9, 3), models.RegistrationPre},
		// Coach 3 – İzmir (3 öğrenci)
		{"Ömer", "Çetin", "05411000007", "İzmir", "Karşıyaka Anadolu Lisesi", "12. Sınıf", "ogrenci7", dob(2007, 4, 18), models.RegistrationConfirmed},
		{"Elif", "Şimşek", "05411000008", "İzmir", "Bornova Lisesi", "11. Sınıf", "ogrenci8", dob(2008, 6, 25), models.RegistrationConfirmed},
		{"Furkan", "Polat", "05411000009", "İzmir", "Konak Fen Lisesi", "10. Sınıf", "ogrenci9", dob(2009, 2, 14), models.RegistrationPre},
		// Coach 4 – Bursa (2 öğrenci)
		{"Hatice", "Güneş", "05411000010", "Bursa", "Nilüfer Anadolu Lisesi", "11. Sınıf", "ogrenci10", dob(2008, 8, 7), models.RegistrationConfirmed},
		{"Tarık", "Aktaş", "05411000011", "Bursa", "Osmangazi Lisesi", "10. Sınıf", "ogrenci11", dob(2009, 12, 20), models.RegistrationConfirmed},
		// Coach 5 – Konya (2 öğrenci) — bu koç 3. haftayı girmedi (uyarı çıkacak)
		{"Rabia", "Kılıç", "05411000012", "Konya", "Meram Fen Lisesi", "11. Sınıf", "ogrenci12", dob(2008, 10, 5), models.RegistrationConfirmed},
		{"Emre", "Boz", "05411000013", "Konya", "Selçuklu Anadolu Lisesi", "9. Sınıf", "ogrenci13", dob(2010, 3, 28), models.RegistrationPre},
		// Eşleştirilmemiş ön kayıtlar
		{"Selin", "Duman", "05411000014", "Antalya", "Muratpaşa Lisesi", "10. Sınıf", "ogrenci14", dob(2009, 6, 11), models.RegistrationPre},
		{"Kerem", "Aydın", "05411000015", "Gaziantep", "Şahinbey Lisesi", "11. Sınıf", "ogrenci15", dob(2008, 8, 19), models.RegistrationPre},
	}

	var students []models.Student
	var studentUserIDs []uint
	for _, d := range sDefs {
		u := models.User{
			Name: d.name, Surname: d.surname, Phone: d.phone,
			City: d.city, Role: models.RoleStudent,
			Username: d.username, PasswordHash: hash("ogrenci123"),
			Status: models.UserStatusActive, DateOfBirth: d.dob,
		}
		if err := database.DB.Create(&u).Error; err != nil {
			fmt.Printf("Öğrenci zaten mevcut: %s\n", d.username)
			continue
		}
		s := models.Student{
			UserID: u.ID, School: d.school, Grade: d.grade,
			RegistrationStatus: d.status,
		}
		database.DB.Create(&s)
		students = append(students, s)
		studentUserIDs = append(studentUserIDs, u.ID)
	}
	fmt.Printf("✓ %d öğrenci oluşturuldu\n", len(students))

	// ── Parents ───────────────────────────────────────────────────────────────
	// 3 veli -> öğrenciler 0, 3, 9 ile eşleştirildi
	parentDefs := []struct {
		name, surname, phone, username string
		studentIdx                      int
	}{
		{"Fatma", "Arslan", "05511000001", "veli1", 0},  // Ali Arslan'ın annesi
		{"Kemal", "Yıldız", "05511000002", "veli2", 3},  // Merve Yıldız'ın babası
		{"Ayşe", "Güneş", "05511000003", "veli3", 9},    // Hatice Güneş'in annesi
	}
	for _, d := range parentDefs {
		u := models.User{
			Name: d.name, Surname: d.surname, Phone: d.phone,
			City: students[d.studentIdx].User.City,
			Role: models.RoleParent, Username: d.username,
			PasswordHash: hash("veli123"), Status: models.UserStatusActive,
		}
		if err := database.DB.Create(&u).Error; err != nil {
			fmt.Printf("Veli zaten mevcut: %s\n", d.username)
			continue
		}
		// Link student to parent
		if d.studentIdx < len(students) {
			database.DB.Model(&students[d.studentIdx]).Update("parent_user_id", u.ID)
		}
	}
	fmt.Println("✓ Veliler oluşturuldu")

	// ── Coach–Student assignments ──────────────────────────────────────────────
	// coach[0]->öğrenci[0,1,2]  coach[1]->öğrenci[3,4,5]
	// coach[2]->öğrenci[6,7,8]  coach[3]->öğrenci[9,10]
	// coach[4]->öğrenci[11,12]
	assignments := [][2]int{
		{0, 0}, {0, 1}, {0, 2},
		{1, 3}, {1, 4}, {1, 5},
		{2, 6}, {2, 7}, {2, 8},
		{3, 9}, {3, 10},
		{4, 11}, {4, 12},
	}
	for _, a := range assignments {
		ci, si := a[0], a[1]
		if ci >= len(coaches) || si >= len(students) {
			continue
		}
		database.DB.Create(&models.CoachStudent{
			CoachID: coaches[ci].ID, StudentID: students[si].ID,
			AssignedAt: ago(21), IsActive: true,
		})
	}
	fmt.Println("✓ Koç–öğrenci atamaları yapıldı")

	// ── Coordinator–Student assignments ───────────────────────────────────────
	// coord[0] (İstanbul) -> öğrenci[0,1,2]
	// coord[1] (Ankara)   -> öğrenci[3,4,5]
	coordAssignments := [][2]int{
		{0, 0}, {0, 1}, {0, 2},
		{1, 3}, {1, 4}, {1, 5},
	}
	for _, a := range coordAssignments {
		ci, si := a[0], a[1]
		if ci >= len(coords) || si >= len(students) {
			continue
		}
		database.DB.Create(&models.CoordinatorStudent{
			CoordinatorID: coords[ci].ID, StudentID: students[si].ID,
			AssignedAt: ago(21), IsActive: true,
		})
	}
	fmt.Println("✓ Koordinatör–öğrenci atamaları yapıldı")

	// ── Meetings ──────────────────────────────────────────────────────────────
	type meetingDef struct {
		coachIdx, studentIdx, weekIdx int
		rating                         int
		notes                          string
		ratings                        [6]int
		catNotes                       [6]string
		status                         models.MeetingStatus
	}

	approvedAt := ago(5)

	meetings := []meetingDef{
		// ── Hafta 1 (onaylandı) ─────────────────────────────────────────
		{0, 0, 0, 4,
			"Ali bu hafta düzenli çalışma programı oluşturdu. Sabah erken kalkıp ders yapma alışkanlığı geliştiriyor. TYT matematik testlerinde net sayısı artmaya başladı.",
			[6]int{4, 5, 4, 3, 4, 5},
			[6]string{
				"Pomodoro tekniğini uyguluyor, günde 6 saat verimli çalışıyor.",
				"Hedefe odaklı ve motivasyonu yüksek. Üniversite hayalinden bahsetti.",
				"TYT'de 340+ hedefliyor, şu an 310 civarında. Mantıklı bir yol haritası var.",
				"Matematik ve Türkçe dengesini iyi kurdu.",
				"Matematik netini 32'den 35'e çıkardı. Türkçe sabit.",
				"Geçen hafta vakıf toplantısına katıldı. Muhabbetleri iyi.",
			},
			models.MeetingStatusApproved},
		{0, 1, 0, 3,
			"Büşra derslere konsantrasyonunu artırmaya çalışıyor. Zaman zaman dikkat dağınıklığı yaşıyor; telefon kullanımını sınırlandırması önerildi.",
			[6]int{3, 3, 3, 3, 3, 4},
			[6]string{
				"Çalışma süresi yeterli ama kalitesini artırması gerekiyor.",
				"Motivasyonu orta düzeyde. Küçük başarılarla cesaretlendiriliyor.",
				"Hedefleri var ama net değil. Somutlaştırma çalışması yapıldı.",
				"Yavaş ama istikrarlı bir ilerleme var.",
				"Biyoloji konularında güçlü, matematik biraz zayıf.",
				"Vakıf etkinliklerine düzenli katılıyor.",
			},
			models.MeetingStatusApproved},
		{0, 2, 0, 5,
			"Yusuf TYT hazırlıklarında harika bir hafta geçirdi. Matematik netini 35'e çıkardı, Türkçe'de de çok başarılı.",
			[6]int{5, 5, 5, 5, 5, 5},
			[6]string{
				"Son derece disiplinli. Günlük hedeflerini eksiksiz tamamlıyor.",
				"Motivasyonu zirveде. Sınavı kazanacağına inanıyor.",
				"İTÜ Elektrik-Elektronik Mühendisliği hedefi var. Puanı oraya gidebilir.",
				"Her hafta belirgin gelişim görülüyor.",
				"TYT toplam net 120 oldu. Hedefte ilerliyor.",
				"Vakıf gidişleri düzenli, arkadaşlarıyla muhabbeti çok iyi.",
			},
			models.MeetingStatusApproved},
		{1, 3, 0, 4,
			"Merve AYT Türkçe alanında güçlü bir hafta geçirdi. Makale analizleri ve edebiyat sorularında başarılı.",
			[6]int{4, 4, 4, 4, 4, 3},
			[6]string{
				"Sistematik çalışıyor; konu listesi yapıp takip ediyor.",
				"Motivasyonu iyi. Arkadaşlarından etkilenmiyor, odağını koruyor.",
				"Boğaziçi Türk Dili hedefi var. Gerçekçi ve ulaşılabilir.",
				"AYT Türkçe alanında somut gelişme var.",
				"Edebiyat ve Dil bilgisi çok güçlü. Matematik daha zayıf.",
				"Vakıf toplantısına bu hafta katılamadı, mazereti vardı.",
			},
			models.MeetingStatusApproved},
		{1, 4, 0, 5,
			"İbrahim Fen Bilimleri alanında harika performans gösterdi. Fizik ve Kimya problemlerini çok hızlı çözüyor.",
			[6]int{5, 5, 5, 5, 5, 4},
			[6]string{
				"Günlük 8 saat düzenli çalışıyor. Çok disiplinli.",
				"Olimpiyat hedefi var; motivasyonu tavan.",
				"ODTÜ Fizik Mühendisliği veya Boğaziçi EEE hedefliyor.",
				"Her geçen hafta yeni bir konuyu bitiriyor.",
				"Fizik netini 38'e çıkardı. Kimya da çok güçlü.",
				"Vakıf muhabbetleri çok iyi; liderlik vasfı var.",
			},
			models.MeetingStatusApproved},
		{1, 5, 0, 2,
			"Seda sisteme yeni alıştı. Motivasyon ve düzen konusunda destek gerekiyor. Ailevi bazı sorunları paylaştı.",
			[6]int{2, 2, 2, 2, 2, 3},
			[6]string{
				"Henüz düzenli çalışma alışkanlığı oturmadı.",
				"Motivasyonu düşük; teşvik edildi.",
				"Hedefleri çok belirsiz. Bir sonraki görüşmede netleştirilecek.",
				"İlk hafta olduğu için gelişim değerlendirilemedi.",
				"Türkçe biraz daha güçlü; diğer alanlarda eksikler var.",
				"Vakıf hakkında bilgi verildi, merak etti.",
			},
			models.MeetingStatusApproved},
		{2, 6, 0, 4,
			"Ömer hem TYT hem AYT'ye eş zamanlı hazırlanıyor. Bu hafta karma çalışma programı oluşturduk.",
			[6]int{4, 4, 4, 4, 3, 4},
			[6]string{
				"TYT ve AYT dengesini iyi kuruyor.",
				"12. sınıf baskısını hissediyor ama yönetiyor.",
				"Dokuz Eylül Tıp hedefliyor. Puan aralığında ilerliyor.",
				"Tarih konularında belirgin gelişim.",
				"Fen bilimleri güçlü, Türkçe biraz daha çalışılmalı.",
				"Vakıf gidişleri düzenli.",
			},
			models.MeetingStatusApproved},
		{2, 7, 0, 3,
			"Elif Türkçe ve matematik dengesini çalışıyoruz. Sosyal bilimler sorularını daha hızlı çözüyor.",
			[6]int{3, 4, 3, 3, 3, 4},
			[6]string{
				"Orta düzey disiplin; bazen program dışına çıkıyor.",
				"İyimser ve istekli; cesaretlendirilmeye ihtiyacı yok.",
				"Ege Üniversitesi Mimarlık hedefliyor.",
				"Bu hafta geometri konularını bitirdi.",
				"Türkçe güçlü, matematik orta.",
				"Arkadaşlarıyla birlikte vakfa gidiyor.",
			},
			models.MeetingStatusApproved},
		{2, 8, 0, 2,
			"Furkan bu hafta düşük motivasyonla geldi. Aile baskısından bahsetti. Destekleyici bir yaklaşım benimsedik.",
			[6]int{2, 2, 3, 2, 2, 3},
			[6]string{
				"Düzensiz çalışma; hedeflerini erteliyor.",
				"Motivasyon çok düşük. Aile baskısı var.",
				"Hedefleri yok denecek kadar belirsiz.",
				"Bu hafta somut bir gelişme gözlemlenmedi.",
				"Matematikte ciddi eksikler var.",
				"Vakıf etkinliklerine bazen katılıyor.",
			},
			models.MeetingStatusApproved},
		{3, 9, 0, 5,
			"Hatice bu hafta tüm deneme sınavlarında hedef puanının üzerinde çıktı. Çok başarılı bir hafta.",
			[6]int{5, 5, 5, 5, 5, 5},
			[6]string{
				"Günlük 9 saat çalışıyor. Öz disiplin mükemmel.",
				"Motivasyonu tavan. Kazanacağından emin.",
				"Hacettepe Tıp birinci tercihinde. Puanı oraya gidebilir.",
				"Her konuyu eksiksiz tamamlıyor.",
				"Biyoloji ve Kimya netlerini 40'a çıkardı.",
				"Vakıf muhabbetleri çok verimli. Gruba liderlik yapıyor.",
			},
			models.MeetingStatusApproved},
		{3, 10, 0, 3,
			"Tarık fizik konularında desteğe ihtiyaç duyuyor. Online kaynak tavsiye ettim. Genel tempoo yeterli.",
			[6]int{3, 3, 3, 3, 2, 4},
			[6]string{
				"Günde 5 saat çalışıyor ama odak zaman zaman dağılıyor.",
				"Motivasyonu orta. Örneklerle teşvik edildi.",
				"Uludağ Mühendislik hedefliyor. Gerçekçi.",
				"Kimya konuları tamamlandı, fizik kaldı.",
				"Fizik netinde sorun var: 12. Çözüm kaynakları verildi.",
				"Vakıf gidişleri düzenli.",
			},
			models.MeetingStatusApproved},
		{4, 11, 0, 4,
			"Rabia bu hafta güzel bir çalışma temposu yakaladı. Matematik net sayısını 28'den 32'ye çıkardı.",
			[6]int{4, 4, 4, 4, 4, 3},
			[6]string{
				"Pomodoro ile çalışıyor. 6 saat verimli ders.",
				"Motivasyonu iyi; üniversite hayali somut.",
				"Selçuk Üniversitesi Tıp hedefliyor.",
				"Bu hafta matematik hızlanması dikkat çekici.",
				"Matematik ve Fen dengesi iyi kuruldu.",
				"Vakıf ile irtibatı var ama gidişler seyrek.",
			},
			models.MeetingStatusApproved},
		{4, 12, 0, 3,
			"Emre ders programı oluşturma konusunda yardım istedi. Birlikte haftalık plan hazırladık.",
			[6]int{3, 3, 3, 3, 3, 4},
			[6]string{
				"Program oluşturuldu ama uygulaması henüz oturmadı.",
				"Motivasyonu var ama yönlendirilmeye ihtiyacı var.",
				"Uzun vadeli hedefleri belirsiz.",
				"İlk hafta çalışması; somut gelişim sonraki haftalarda görülecek.",
				"Türkçe biraz daha güçlü.",
				"Vakıf etkinliklerine meraklı.",
			},
			models.MeetingStatusApproved},

		// ── Hafta 2 (onaylandı) ─────────────────────────────────────────
		{0, 0, 1, 4,
			"Ali bu hafta Türkçe paragraf testlerinde büyük ilerleme kaydetti. Okuduğunu anlama hızı arttı.",
			[6]int{4, 5, 4, 4, 4, 5},
			[6]string{
				"Çalışma programını aksatmadan uyguladı.",
				"Motivasyonu devam ediyor. Deneme sonu iyi hissettirdi.",
				"TYT net ortalaması 315'e çıktı.",
				"Türkçe hız kazanma dikkat çekici.",
				"Türkçe net +5 oldu. Matematik sabit.",
				"Vakıf muhabbet grubuna katıldı.",
			},
			models.MeetingStatusApproved},
		{0, 1, 1, 4,
			"Büşra konsantrasyon sorunlarını büyük ölçüde aşmaya başladı. Pomodoro tekniğini benimsedi.",
			[6]int{4, 4, 3, 4, 3, 4},
			[6]string{
				"Pomodoro ile çalışma süresi 5 saate çıktı.",
				"Bu hafta çok daha pozitif geldi görüşmeye.",
				"Hedeflerini not ederek takip ediyor.",
				"Konsantrasyon belirgin şekilde arttı.",
				"Biyoloji netinde artış var.",
				"Etkinliklere devam ediyor.",
			},
			models.MeetingStatusApproved},
		{0, 2, 1, 5,
			"Yusuf son TYT denemesinde 385 net yaptı. Hedefine çok yaklaştı. Motivasyonu devam ediyor.",
			[6]int{5, 5, 5, 5, 5, 5},
			[6]string{
				"Ders programını sıfır hatayla uyguluyor.",
				"Sınavı geçeceğinden emin; çok sakin.",
				"İTÜ hedefine puanı yetecek seviyeye geldi.",
				"Her hafta 5+ net artışı kaydediyor.",
				"Matematik: 38 net. Türkçe: 37 net.",
				"Vakıf etkinliklerine düzenli; arkadaşları motive ediyor.",
			},
			models.MeetingStatusApproved},
		{1, 3, 1, 4,
			"Merve bu hafta kelime bilgisini geliştirmek için roman okumaya başladı. Edebi metin analizleri çok iyi.",
			[6]int{4, 4, 4, 4, 4, 4},
			[6]string{
				"Roman okuma alışkanlığı kazanıyor.",
				"Pozitif ve kararlı.",
				"AYT Türkçe netini 38'e çıkardı.",
				"Edebiyat alanında kayda değer ilerleme.",
				"Türkçe net artışı devam ediyor.",
				"Vakıf gidişlerine başladı.",
			},
			models.MeetingStatusApproved},
		{1, 4, 1, 5,
			"İbrahim Kimya testlerinde de üst düzey performans gösteriyor. Fizik olimpiyat ön elemeye katılacak.",
			[6]int{5, 5, 5, 5, 5, 5},
			[6]string{
				"Günde 9 saat çalışıyor. Çok disiplinli.",
				"Olimpiyat motivasyonu çok yüksek.",
				"Boğaziçi EEE puanına ulaştı.",
				"Her geçen gün yeni zirve.",
				"Fizik+Kimya net toplamı 78.",
				"Vakıf muhabbet grubuna liderlik yapıyor.",
			},
			models.MeetingStatusApproved},
		{1, 5, 1, 3,
			"Seda bu hafta çok daha motive görünüyordu. Hedeflerini netleştirdik. Biraz daha zaman lazım.",
			[6]int{3, 3, 3, 3, 3, 3},
			[6]string{
				"Günde 4 saat çalışmaya başladı.",
				"İkinci hafta motivasyonu yükseldi.",
				"Eczacılık veya Hemşirelik düşünüyor.",
				"Yavaş ama istikrarlı ilerleme.",
				"Biyoloji biraz daha güçlü.",
				"Vakıf toplantısına bu hafta katıldı.",
			},
			models.MeetingStatusApproved},
		{2, 6, 1, 4,
			"Ömer tarih konularını tamamladı. Yükseköğretim hedefini netleştirdi. Deneme puanları yüksek.",
			[6]int{4, 4, 4, 4, 4, 4},
			[6]string{
				"TYT+AYT dengesi iyi kuruluyor.",
				"Sınav heyecanını yönetiyor.",
				"DEÜ Tıp kesin hedef.",
				"Tarih konuları tamamen bitti.",
				"Tarih netini 17'ye çıkardı.",
				"Vakıf gidişleri sürekli.",
			},
			models.MeetingStatusApproved},
		{2, 7, 1, 4,
			"Elif matematik günlük 3 saat çalışıyor. Geometri konularında net sayısı arttı.",
			[6]int{4, 4, 3, 4, 4, 4},
			[6]string{
				"Çalışma süresi artmış.",
				"Pozitif ve iddialı.",
				"EGE Mimarlık hedefinde ilerliyor.",
				"Geometri konularında belirgin atılım.",
				"Geometri net +4 oldu.",
				"Vakıf arkadaşlarıyla çalışma grupları yapıyor.",
			},
			models.MeetingStatusApproved},
		{2, 8, 1, 3,
			"Furkan motivasyonu toparladı. Arkadaşlarıyla birlikte ders çalışmaya başladı. İyiye gidiyor.",
			[6]int{3, 3, 3, 3, 3, 4},
			[6]string{
				"Grup çalışması motivasyonu artırdı.",
				"İkinci haftada belirgin iyileşme.",
				"Hedefler netleşiyor.",
				"Aile baskısını biraz daha yönetiyor.",
				"Matematik netinde küçük artış.",
				"Vakıf etkinliklerine katılmaya başladı.",
			},
			models.MeetingStatusApproved},
		{3, 9, 1, 5,
			"Hatice tüm alanlarda dengeli ve güçlü bir hafta geçirdi. Burs başvurusunu tamamladı.",
			[6]int{5, 5, 5, 5, 5, 5},
			[6]string{
				"Çalışma temposu eksiksiz.",
				"Motivasyonu zirvede.",
				"Hacettepe Tıp hedefine bu hafta da adım attı.",
				"Burs başvurusu yapıldı: güçlü bir portfolyo.",
				"Biyoloji: 39 net. Kimya: 38 net.",
				"Vakıf etkinliği düzenliyor: çok aktif.",
			},
			models.MeetingStatusApproved},
		{3, 10, 1, 4,
			"Tarık bu hafta fizik sorularında belirgin ilerleme kaydetti. Formülleri ezberlemek yerine anlamayı tercih ediyor.",
			[6]int{4, 4, 4, 4, 4, 4},
			[6]string{
				"Çalışma kalitesi arttı.",
				"Geçen haftaya göre çok daha motive.",
				"Uludağ Mühendislik hedefine yaklaşıyor.",
				"Fizik konularında anlayışla ilerleme.",
				"Fizik neti 12'den 17'ye çıktı.",
				"Vakıf gidişleri devam ediyor.",
			},
			models.MeetingStatusApproved},
		{4, 11, 1, 4,
			"Rabia sınav kaygısını yönetmede ilerleme kaydetti. Nefes egzersizleri yapıyor. Başarılı bir hafta.",
			[6]int{4, 4, 4, 4, 4, 4},
			[6]string{
				"Sınav öncesi rutin oluşturdu.",
				"Kaygısını yönetmeyi öğreniyor.",
				"Selçuk Tıp hedefinde güvenle ilerliyor.",
				"Kaygı yönetiminde somut ilerleme.",
				"Matematik net 32'de sabit; hedef 38.",
				"Vakıf muhabbetleri daha düzenli.",
			},
			models.MeetingStatusApproved},
		{4, 12, 1, 4,
			"Emre bu hafta çok daha disiplinli çalıştı. Günlük hedeflerini takip ediyor. Güzel ilerleme.",
			[6]int{4, 4, 4, 4, 3, 4},
			[6]string{
				"Haftalık planı eksiksiz uyguladı.",
				"Motivasyonu artmış; planlamayı seviyor.",
				"Kısa vadeli hedeflerle ilerleme sağlandı.",
				"İkinci haftada belirgin ilerleme.",
				"Türkçe güçlü; matematikte destek lazım.",
				"Vakıf etkinliklerine katılım arttı.",
			},
			models.MeetingStatusApproved},

		// ── Hafta 3 – aktif hafta (koç 4 / Hasan girmedi → uyarı) ───────
		{0, 0, 2, 5,
			"Ali son deneme sınavında matematik netini 40'a çıkardı. Mükemmel bir hafta. Sınava hazır.",
			[6]int{5, 5, 5, 5, 5, 5},
			[6]string{
				"Çalışma programı mükemmel işliyor.",
				"Hedefine olan inancı çok güçlü.",
				"TYT ortalaması 350'yi geçti.",
				"Her haftaki gelişim ivmesi artarak devam ediyor.",
				"Matematik: 40 net. Türkçe: 38 net.",
				"Vakıf etkinliğinde sunum yaptı. Harika!",
			},
			models.MeetingStatusPending},
		{0, 1, 2, 4,
			"Büşra Türkçe alanında güçlendi. Paragraf sorularında neredeyse hiç hata yapmıyor.",
			[6]int{4, 4, 4, 4, 4, 4},
			[6]string{
				"Çalışma kalitesi üst düzey.",
				"Pozitif enerjiyle geliyor görüşmeye.",
				"Hedefine ulaşma olasılığı yüksek.",
				"Türkçe alanında ani yükseliş.",
				"Türkçe: 38 net. Matematik: 25 net.",
				"Düzenli vakıf katılımı var.",
			},
			models.MeetingStatusPending},
		{0, 2, 2, 5,
			"Yusuf hedef üniversitesini netleştirdi. Motivasyonu çok yüksek. İTÜ için puanı yeterli.",
			[6]int{5, 5, 5, 5, 5, 5},
			[6]string{
				"Çalışma disiplini eksiksiz.",
				"Zirve motivasyonunda.",
				"İTÜ EEE için puanı yeterli.",
				"Son haftaların en başarılısı.",
				"TYT toplam: 390+",
				"Vakıf grubunda aktif rol üstlendi.",
			},
			models.MeetingStatusPending},
		{1, 3, 2, 4,
			"Merve bu hafta AYT Türkçe denemesinde en yüksek puanını aldı. Çok güzel bir gelişim.",
			[6]int{4, 5, 4, 4, 4, 4},
			[6]string{
				"Çalışma rutini çok oturmuş.",
				"Bu hafta çok mutlu geldi; sınavdan umutlu.",
				"Boğaziçi için puan aralığına girmeye başladı.",
				"AYT Türkçe patlaması.",
				"AYT Türkçe net: 38.",
				"Vakıf toplantısında konuşma yaptı.",
			},
			models.MeetingStatusPending},
		{1, 4, 2, 5,
			"İbrahim Fizik olimpiyatı ön eleme turuna geçti. Tebrik edildi. Performansı inanılmaz.",
			[6]int{5, 5, 5, 5, 5, 5},
			[6]string{
				"Olimpiyat çalışması + YKS dengesi mükemmel.",
				"Olimpiyat başarısı motivasyonu ikiye katladı.",
				"Boğaziçi EEE kesinleşti.",
				"Tarihi bir başarı: olimpiyat eleme turuna geçti.",
				"Fizik neti 40'a ulaştı.",
				"Vakıf topluluğuna ilham veriyor.",
			},
			models.MeetingStatusPending},
		{1, 5, 2, 3,
			"Seda bu hafta vakıf toplantısına katıldı. Güzel bir deneyim yaşadı. Derslerinde de iyileşme var.",
			[6]int{3, 3, 3, 3, 3, 5},
			[6]string{
				"Çalışma süresi artmaya devam ediyor.",
				"Motivasyonu yükseliş trendinde.",
				"Hemşirelik veya Eczacılık kararını verdi.",
				"Sosyal ortam derse motivasyonu artırdı.",
				"Biyoloji netinde artış var.",
				"Vakıf toplantısı çok olumlu etki bıraktı. Devam etmek istiyor.",
			},
			models.MeetingStatusPending},
		{2, 6, 2, 4,
			"Ömer sınav heyecanını yönetmeyi öğrendi. Sınav stratejisi çalışmaları çok verimli.",
			[6]int{4, 4, 4, 4, 4, 4},
			[6]string{
				"Sınav stratejisi oturdu.",
				"Sakin ve kararlı.",
				"DEÜ Tıp için sınırda; son haftalar kritik.",
				"Zaman yönetimi çok iyileşti.",
				"TYT+AYT dengesi korunuyor.",
				"Vakıf gidişleri devam ediyor.",
			},
			models.MeetingStatusPending},
		{2, 7, 2, 4,
			"Elif koçluktan çok memnun olduğunu ifade etti. İlerleme kaydetmesi çok sevindirici.",
			[6]int{4, 5, 4, 4, 4, 4},
			[6]string{
				"Çalışma temposu çok iyi.",
				"Motive, mutlu ve hevesli.",
				"EGE Mimarlık hedefi sağlam.",
				"Son 3 haftanın en başarılısı.",
				"Geometri netinde sürekli artış.",
				"Vakıf arkadaşlarıyla çok uyumlu.",
			},
			models.MeetingStatusPending},
		{2, 8, 2, 3,
			"Furkan matematik konularında daha sistematik çalışmaya başladı. Motivasyonu korunuyor.",
			[6]int{3, 3, 3, 3, 3, 4},
			[6]string{
				"Sistematik çalışma alışkanlığı oturmuş.",
				"Motivasyonu orta ama stabil.",
				"Hedefleri biraz netleşti.",
				"Yavaş ama sürekli ilerleme.",
				"Matematik netinde küçük artış.",
				"Vakıf etkinliklerine devam ediyor.",
			},
			models.MeetingStatusPending},
		{3, 9, 2, 5,
			"Hatice bu hafta mock sınavda 455 puan aldı. Tebrikler! Hacettepe Tıp kesinleşti.",
			[6]int{5, 5, 5, 5, 5, 5},
			[6]string{
				"Çalışma sistemi mükemmel.",
				"Başarı motivasyonu daha da artırdı.",
				"Hacettepe Tıp kesin hedef. Puan oraya gidiyor.",
				"Mock sınav zirvesi: 455 puan.",
				"Biyoloji+Kimya net toplamı 78.",
				"Vakıf etkinliği organize ediyor.",
			},
			models.MeetingStatusPending},
		{3, 10, 2, 4,
			"Tarık kimya konusunda da iyileşme gösteriyor. Deney videoları izlemeye başladı. Başarılı hafta.",
			[6]int{4, 4, 4, 4, 4, 4},
			[6]string{
				"Çalışma kalitesi artmış.",
				"Bu hafta çok daha motivasyonlu geldi.",
				"Uludağ Mühendislik için neti yeterliye yaklaşıyor.",
				"Kimya konularında da ilerleme.",
				"Kimya neti 14'ten 18'e çıktı.",
				"Vakıf gidişleri düzenli.",
			},
			models.MeetingStatusPending},
		// Coach 4 (Hasan/Konya) → HAFTA 3 GÖRÜŞMESİ YOK (uyarı listesinde görünecek)
	}

	approvedUID := uint(1)
	for _, md := range meetings {
		if md.coachIdx >= len(coaches) || md.studentIdx >= len(students) || md.weekIdx >= len(weeks) {
			continue
		}
		det := catJSON(md.ratings, md.catNotes)
		m := models.Meeting{
			CoachID:   coaches[md.coachIdx].ID,
			StudentID: students[md.studentIdx].ID,
			WeekID:    weeks[md.weekIdx].ID,
			Rating:    md.rating,
			Notes:     md.notes,
			Details:   det,
			Status:    md.status,
		}
		if md.status == models.MeetingStatusApproved {
			m.ApprovedBy = &approvedUID
			m.ApprovedAt = &approvedAt
		}
		database.DB.Create(&m)
	}
	fmt.Printf("✓ %d görüşme kaydı oluşturuldu\n", len(meetings))

	// ── Messages ──────────────────────────────────────────────────────────────
	msgs := []struct {
		senderIdx, recipientIdx int // index into userIDs: 0=admin, 1-5=coaches, 6-7=coords
		senderID, recipientID   uint
		body                    string
		daysAgo                 int
	}{}

	adminID := uint(1)
	var coach1ID, coach2ID, coach3ID, coord1ID, parent1ID uint
	if len(coachUserIDs) > 0 {
		coach1ID = coachUserIDs[0]
	}
	if len(coachUserIDs) > 1 {
		coach2ID = coachUserIDs[1]
	}
	if len(coachUserIDs) > 2 {
		coach3ID = coachUserIDs[2]
	}
	// Get coordinator user IDs
	var coordU []models.User
	database.DB.Where("role = ?", models.RoleCoordinator).Find(&coordU)
	if len(coordU) > 0 {
		coord1ID = coordU[0].ID
	}
	// Get first parent user ID
	var parentU []models.User
	database.DB.Where("role = ?", models.RoleParent).Find(&parentU)
	if len(parentU) > 0 {
		parent1ID = parentU[0].ID
	}

	type msgDef struct {
		senderID, recipientID uint
		body                  string
		daysAgo               int
	}

	_ = msgs
	messageDefs := []msgDef{}

	if coach1ID > 0 {
		messageDefs = append(messageDefs,
			msgDef{coach1ID, adminID, "Merhaba, bu hafta Ali Arslan ile çok verimli bir görüşme yaptık. Matematik netini 40'a çıkarmayı başardı. Devam eden motivasyonu için ödüllendirme önerilebilir.", 5},
			msgDef{adminID, coach1ID, "Harika haber Ahmet Bey! Ali'nin gelişimini yakından takip ediyoruz. Teşekkürler.", 5},
			msgDef{coach1ID, adminID, "Yusuf için hedef üniversite bilgisi gerekiyor; veli ile görüşme ayarlanabilir mi?", 3},
			msgDef{adminID, coach1ID, "Tabii, bu hafta içinde ayarlayacağız.", 3},
		)
	}
	if coach2ID > 0 {
		messageDefs = append(messageDefs,
			msgDef{coach2ID, adminID, "İbrahim'in fizik olimpiyatına katılacağını bildirmek istedim. Destek verebilir miyiz?", 7},
			msgDef{adminID, coach2ID, "Kesinlikle! Önce okul müdürlüğüyle iletişime geçelim. Detayları paylaşır mısınız?", 7},
			msgDef{coach2ID, adminID, "Olimpiyat 25 Mayıs'ta. İbrahim çok heyecanlı. Ulaşım desteği sağlanabilir mi?", 6},
		)
	}
	if coach3ID > 0 {
		messageDefs = append(messageDefs,
			msgDef{coach3ID, adminID, "Furkan bu hafta motivasyon konusunda ilerledi. Aile durumu hakkında bilgi almak istiyorum.", 10},
			msgDef{adminID, coach3ID, "Furkan'ın ailesiyle geçen hafta görüştük. Bilgi paylaşacağız.", 10},
		)
	}
	if coord1ID > 0 {
		messageDefs = append(messageDefs,
			msgDef{coord1ID, adminID, "İstanbul öğrencilerinin bu haftaki durumunu inceledim. Genel seyir çok iyi. Teşekkürler.", 4},
			msgDef{adminID, coord1ID, "Güzel haber Mustafa Bey. Yusuf özellikle çok başarılı.", 4},
		)
	}
	if parent1ID > 0 {
		messageDefs = append(messageDefs,
			msgDef{parent1ID, adminID, "Merhaba, Ali'nin bu haftaki görüşmesinden haberdardım. Çok mutlu oldum. Teşekkürler.", 6},
			msgDef{adminID, parent1ID, "Fatma Hanım merhaba. Ali gerçekten çok başarılı bir öğrenci. Koçuyla uyum çok iyi.", 6},
		)
	}

	for _, md := range messageDefs {
		if md.senderID == 0 || md.recipientID == 0 {
			continue
		}
		database.DB.Create(&models.Message{
			SenderID:    md.senderID,
			RecipientID: md.recipientID,
			Body:        md.body,
			CreatedAt:   ago(md.daysAgo),
			IsRead:      md.daysAgo > 2,
		})
	}
	fmt.Printf("✓ %d mesaj oluşturuldu\n", len(messageDefs))

	// ── SMS Logs ──────────────────────────────────────────────────────────────
	smsLogs := []models.SMSLog{
		{Phone: "05411000001", Message: "KARE Rehber: Giriş bilgileriniz — Kullanıcı adı: ogrenci1, Şifre: ogrenci123. kare.ulued.org", Type: models.SMSTypeCredentials, Status: models.SMSStatusSent, CreatedAt: ago(20)},
		{Phone: "05411000004", Message: "KARE Rehber: Giriş bilgileriniz — Kullanıcı adı: ogrenci4, Şifre: ogrenci123. kare.ulued.org", Type: models.SMSTypeCredentials, Status: models.SMSStatusSent, CreatedAt: ago(20)},
		{Phone: "05321000001", Message: "KARE Rehber: Koç paneline giriş bilgileriniz — Kullanıcı adı: koc1, Şifre: koc123. kare.ulued.org", Type: models.SMSTypeCredentials, Status: models.SMSStatusSent, CreatedAt: ago(22)},
		{Phone: "05321000002", Message: "KARE Rehber: Koç paneline giriş bilgileriniz — Kullanıcı adı: koc2, Şifre: koc123. kare.ulued.org", Type: models.SMSTypeCredentials, Status: models.SMSStatusSent, CreatedAt: ago(22)},
		{Phone: "05321000005", Message: "KARE: Bu hafta görüşme notunuzu girmediniz. Lütfen sisteme giriniz.", Type: models.SMSTypeAlert, Status: models.SMSStatusSent, CreatedAt: ago(1)},
		{Phone: "05411000006", Message: "KARE Rehber: Ön kaydınız alınmıştır. En kısa sürede sizinle iletişime geçeceğiz.", Type: models.SMSTypeCustom, Status: models.SMSStatusSent, CreatedAt: ago(14)},
	}
	for _, s := range smsLogs {
		database.DB.Create(&s)
	}
	fmt.Printf("✓ %d SMS logu oluşturuldu\n", len(smsLogs))

	fmt.Println("\n✅ Seed tamamlandı!")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("Admin    → admin / admin123")
	fmt.Println("Koçlar   → koc1..koc5 / koc123")
	fmt.Println("Koordinatör → koordinator1, koordinator2 / koor123")
	fmt.Println("Veliler  → veli1, veli2, veli3 / veli123")
	fmt.Println("Öğrenciler → ogrenci1..ogrenci15 / ogrenci123")
	fmt.Println("─────────────────────────────────────────")
}
