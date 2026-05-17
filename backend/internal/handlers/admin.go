package handlers

import (
	"fmt"
	"strconv"
	"time"

	"encoding/json"

	"kare-rehber/backend/internal/database"
	"kare-rehber/backend/internal/middleware"
	"kare-rehber/backend/internal/models"
	"kare-rehber/backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

// --- Users ---

func AdminListUsers(c *fiber.Ctx) error {
	role := c.Query("role")
	var users []models.User
	q := database.DB.Model(&models.User{})
	if role != "" {
		q = q.Where("role = ?", role)
	}
	q.Find(&users)
	return c.JSON(users)
}

func AdminCreateUser(c *fiber.Ctx) error {
	var body struct {
		Name        string `json:"name"`
		Surname     string `json:"surname"`
		Phone       string `json:"phone"`
		City        string `json:"city"`
		Role        string `json:"role"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		DateOfBirth string `json:"date_of_birth"`
		Organization string `json:"organization"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	hash, err := services.HashPassword(body.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error"})
	}

	user := models.User{
		Name:         body.Name,
		Surname:      body.Surname,
		Phone:        body.Phone,
		City:         body.City,
		Role:         models.Role(body.Role),
		Username:     body.Username,
		PasswordHash: hash,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username or phone already exists"})
	}

	switch models.Role(body.Role) {
	case models.RoleCoach:
		database.DB.Create(&models.Coach{UserID: user.ID, PoolStatus: models.PoolStatusApproved})
	case models.RoleCoordinator:
		database.DB.Create(&models.Coordinator{UserID: user.ID, Organization: body.Organization})
	case models.RoleStudent:
		database.DB.Create(&models.Student{UserID: user.ID, RegistrationStatus: models.RegistrationConfirmed})
	}

	logAudit(c, "create_user", "user", &user.ID, map[string]interface{}{"username": user.Username, "role": string(user.Role)})
	return c.Status(fiber.StatusCreated).JSON(user)
}

func AdminUpdateUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	var body struct {
		Name    string `json:"name"`
		Surname string `json:"surname"`
		Phone   string `json:"phone"`
		City    string `json:"city"`
		Status  string `json:"status"`
	}
	c.BodyParser(&body)

	updates := map[string]interface{}{}
	if body.Name != "" {
		updates["name"] = body.Name
	}
	if body.Surname != "" {
		updates["surname"] = body.Surname
	}
	if body.Phone != "" {
		updates["phone"] = body.Phone
	}
	if body.City != "" {
		updates["city"] = body.City
	}
	if body.Status != "" {
		updates["status"] = body.Status
	}

	database.DB.Model(&user).Updates(updates)
	logAudit(c, "update_user", "user", &user.ID, updates)
	return c.JSON(user)
}

// --- Students ---

func AdminListStudents(c *fiber.Ctx) error {
	city := c.Query("city")
	status := c.Query("registration_status")
	var students []models.Student
	q := database.DB.Preload("User")
	if city != "" {
		q = q.Joins("JOIN users ON users.id = students.user_id").Where("users.city = ?", city)
	}
	if status != "" {
		q = q.Where("students.registration_status = ?", status)
	}
	q.Find(&students)
	return c.JSON(students)
}

func AdminUpdateStudent(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var student models.Student
	if err := database.DB.First(&student, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
	}
	var body struct {
		RegistrationStatus string `json:"registration_status"`
		School             string `json:"school"`
		Grade              string `json:"grade"`
		Notes              string `json:"notes"`
	}
	c.BodyParser(&body)
	updates := map[string]interface{}{}
	if body.RegistrationStatus != "" {
		updates["registration_status"] = body.RegistrationStatus
	}
	if body.School != "" {
		updates["school"] = body.School
	}
	if body.Grade != "" {
		updates["grade"] = body.Grade
	}
	if body.Notes != "" {
		updates["notes"] = body.Notes
	}
	database.DB.Model(&student).Updates(updates)
	return c.JSON(student)
}

// --- Coaches ---

func AdminListCoaches(c *fiber.Ctx) error {
	poolStatus := c.Query("pool_status")
	var coaches []models.Coach
	q := database.DB.Preload("User")
	if poolStatus != "" {
		q = q.Where("pool_status = ?", poolStatus)
	}
	q.Find(&coaches)
	return c.JSON(coaches)
}

func AdminUpdateCoach(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var coach models.Coach
	if err := database.DB.First(&coach, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Coach not found"})
	}
	var body struct {
		PoolStatus string `json:"pool_status"`
	}
	c.BodyParser(&body)
	if body.PoolStatus != "" {
		database.DB.Model(&coach).Update("pool_status", body.PoolStatus)

		if models.PoolStatus(body.PoolStatus) == models.PoolStatusApproved {
			database.DB.Model(&models.User{}).Where("id = ?", coach.UserID).Update("status", models.UserStatusActive)
		}
	}
	return c.JSON(coach)
}

// --- Coordinators ---

func AdminListCoordinators(c *fiber.Ctx) error {
	var coords []models.Coordinator
	database.DB.Preload("User").Find(&coords)
	return c.JSON(coords)
}

// --- Weeks ---

func AdminListWeeks(c *fiber.Ctx) error {
	var weeks []models.Week
	database.DB.Order("week_number desc").Find(&weeks)
	return c.JSON(weeks)
}

func AdminCreateWeek(c *fiber.Ctx) error {
	var body struct {
		WeekNumber int    `json:"week_number"`
		Label      string `json:"label"`
		StartDate  string `json:"start_date"`
		EndDate    string `json:"end_date"`
		IsActive   bool   `json:"is_active"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	start, _ := time.Parse("2006-01-02", body.StartDate)
	end, _ := time.Parse("2006-01-02", body.EndDate)

	if body.IsActive {
		database.DB.Model(&models.Week{}).Where("is_active = true").Update("is_active", false)
	}

	week := models.Week{
		WeekNumber: body.WeekNumber,
		Label:      body.Label,
		StartDate:  start,
		EndDate:    end,
		IsActive:   body.IsActive,
		CreatedBy:  middleware.CurrentUserID(c),
	}
	database.DB.Create(&week)
	return c.Status(fiber.StatusCreated).JSON(week)
}

func AdminUpdateWeek(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var week models.Week
	if err := database.DB.First(&week, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Week not found"})
	}
	var body struct {
		IsActive bool `json:"is_active"`
		IsLocked bool `json:"is_locked"`
	}
	c.BodyParser(&body)

	if body.IsActive {
		database.DB.Model(&models.Week{}).Where("is_active = true AND id != ?", id).Update("is_active", false)
	}

	database.DB.Model(&week).Updates(map[string]interface{}{
		"is_active": body.IsActive,
		"is_locked": body.IsLocked,
	})
	return c.JSON(week)
}

// --- Meetings ---

func AdminListMeetings(c *fiber.Ctx) error {
	weekID := c.Query("week_id")
	status := c.Query("status")
	var meetings []models.Meeting
	q := database.DB.Preload("Student.User").Preload("Coach.User").Preload("Week")
	if weekID != "" {
		q = q.Where("week_id = ?", weekID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	q.Find(&meetings)
	return c.JSON(meetings)
}

func AdminApproveMeeting(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	adminID := middleware.CurrentUserID(c)
	now := time.Now()
	result := database.DB.Model(&models.Meeting{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      models.MeetingStatusApproved,
		"approved_by": adminID,
		"approved_at": now,
	})
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Meeting not found"})
	}
	logAudit(c, "approve_meeting", "meeting", uintPtr(uint(id)), nil)
	return c.JSON(fiber.Map{"message": "Görüşme onaylandı"})
}

func AdminRejectMeeting(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	result := database.DB.Model(&models.Meeting{}).Where("id = ?", id).Update("status", models.MeetingStatusRejected)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Meeting not found"})
	}
	logAudit(c, "reject_meeting", "meeting", uintPtr(uint(id)), nil)
	return c.JSON(fiber.Map{"message": "Görüşme reddedildi"})
}

func AdminEditMeeting(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	adminID := middleware.CurrentUserID(c)
	var body struct {
		Rating int    `json:"rating"`
		Notes  string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := services.AdminEditMeeting(uint(id), adminID, body.Rating, body.Notes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	logAudit(c, "edit_meeting", "meeting", uintPtr(uint(id)), map[string]interface{}{"rating": body.Rating})
	return c.JSON(fiber.Map{"message": "Görüşme güncellendi"})
}

// --- Coach Alerts ---

func AdminCoachAlerts(c *fiber.Ctx) error {
	coaches, err := services.GetCoachAlerts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(coaches)
}

// --- Matching ---

func AdminMatchCoach(c *fiber.Ctx) error {
	var body struct {
		CoachID    uint   `json:"coach_id"`
		StudentIDs []uint `json:"student_ids"`
	}
	if err := c.BodyParser(&body); err != nil || len(body.StudentIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "coach_id and student_ids required"})
	}
	if err := services.AssignCoachToStudents(body.CoachID, body.StudentIDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	logAudit(c, "match_coach", "coach_student", &body.CoachID, map[string]interface{}{"student_count": len(body.StudentIDs)})
	return c.JSON(fiber.Map{"message": fmt.Sprintf("%d öğrenci eşleştirildi", len(body.StudentIDs))})
}

func AdminMatchCoordinator(c *fiber.Ctx) error {
	var body struct {
		CoordinatorID uint   `json:"coordinator_id"`
		StudentIDs    []uint `json:"student_ids"`
	}
	if err := c.BodyParser(&body); err != nil || len(body.StudentIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "coordinator_id and student_ids required"})
	}
	if err := services.AssignCoordinatorToStudents(body.CoordinatorID, body.StudentIDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": fmt.Sprintf("%d öğrenci eşleştirildi", len(body.StudentIDs))})
}

func AdminStudentsByCity(c *fiber.Ctx) error {
	city := c.Query("city")
	students, err := services.GetStudentsByCity(city)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(students)
}

// --- SMS ---

func AdminSendSMS(c *fiber.Ctx) error {
	var body struct {
		Phone   string `json:"phone"`
		Message string `json:"message"`
		Type    string `json:"type"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	smsType := models.SMSType(body.Type)
	if smsType == "" {
		smsType = models.SMSTypeCustom
	}
	if err := services.SendSMS(nil, body.Phone, body.Message, smsType); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SMS gönderilemedi"})
	}
	return c.JSON(fiber.Map{"message": "SMS gönderildi"})
}

func AdminSendCredentialsSMS(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	user, err := services.FindUserByID(uint(id))
	if err != nil || user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	msg := fmt.Sprintf("KARE-REHBER giriş bilgileriniz: Kullanıcı adı: %s | Şifre: %s | Link: kare.ulued.org",
		user.Username, user.Phone[len(user.Phone)-4:])
	uid := user.ID
	if err := services.SendSMS(&uid, user.Phone, msg, models.SMSTypeCredentials); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SMS gönderilemedi"})
	}
	return c.JSON(fiber.Map{"message": "Giriş bilgileri SMS ile gönderildi"})
}

func AdminSMSLogs(c *fiber.Ctx) error {
	var logs []models.SMSLog
	database.DB.Order("created_at desc").Find(&logs)
	return c.JSON(logs)
}

// --- Reports ---

func AdminReportOverview(c *fiber.Ctx) error {
	report, err := services.GetOverviewReport()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(report)
}

func AdminReportByCity(c *fiber.Ctx) error {
	report, err := services.GetCityReport()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(report)
}

func AdminReportMissingMeetings(c *fiber.Ctx) error {
	coaches, err := services.GetMissingMeetingsReport()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(coaches)
}

func AdminReportMentorPerformance(c *fiber.Ctx) error {
	report, err := services.GetMentorPerformanceReport()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(report)
}

// --- Audit Logs ---

func AdminAuditLogs(c *fiber.Ctx) error {
	var logs []models.AuditLog
	database.DB.Preload("User").Order("created_at desc").Limit(500).Find(&logs)
	return c.JSON(logs)
}

func AdminListMessages(c *fiber.Ctx) error {
	var messages []models.Message
	database.DB.Preload("Sender").Order("created_at desc").Limit(200).Find(&messages)
	return c.JSON(messages)
}

func AdminSendMessage(c *fiber.Ctx) error {
	userID := middleware.CurrentUserID(c)
	var body struct {
		RecipientID uint   `json:"recipient_id"`
		Body        string `json:"body"`
	}
	if err := c.BodyParser(&body); err != nil || body.Body == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "recipient_id and body required"})
	}
	msg := models.Message{
		SenderID:    userID,
		RecipientID: body.RecipientID,
		Body:        body.Body,
	}
	database.DB.Create(&msg)
	return c.Status(fiber.StatusCreated).JSON(msg)
}

// --- Helpers ---

func logAudit(c *fiber.Ctx, action, entity string, entityID *uint, details map[string]interface{}) {
	var detailsJSON datatypes.JSON
	if details != nil {
		b, _ := json.Marshal(details)
		detailsJSON = datatypes.JSON(b)
	}
	entry := models.AuditLog{
		UserID:     middleware.CurrentUserID(c),
		Action:     action,
		EntityType: entity,
		EntityID:   entityID,
		Details:    detailsJSON,
	}
	database.DB.Create(&entry)
}

func uintPtr(v uint) *uint { return &v }
