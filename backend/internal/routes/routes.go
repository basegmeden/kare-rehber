package routes

import (
	"kare-rehber/backend/internal/handlers"
	"kare-rehber/backend/internal/middleware"
	"kare-rehber/backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api/v1")

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", handlers.Login)
	auth.Get("/me", middleware.Auth(), handlers.Me)

	// Public registration
	reg := api.Group("/register")
	reg.Post("/student", handlers.RegisterStudent)
	reg.Post("/coach", handlers.RegisterCoach)

	// Admin
	admin := api.Group("/admin",
		middleware.Auth(),
		middleware.RequireRole(models.RoleAdmin),
	)
	admin.Get("/users", handlers.AdminListUsers)
	admin.Post("/users", handlers.AdminCreateUser)
	admin.Put("/users/:id", handlers.AdminUpdateUser)

	admin.Get("/students", handlers.AdminListStudents)
	admin.Put("/students/:id", handlers.AdminUpdateStudent)

	admin.Get("/coaches", handlers.AdminListCoaches)
	admin.Put("/coaches/:id", handlers.AdminUpdateCoach)
	admin.Get("/coaches/alerts", handlers.AdminCoachAlerts)

	admin.Get("/coordinators", handlers.AdminListCoordinators)

	admin.Get("/weeks", handlers.AdminListWeeks)
	admin.Post("/weeks", handlers.AdminCreateWeek)
	admin.Put("/weeks/:id", handlers.AdminUpdateWeek)

	admin.Get("/meetings", handlers.AdminListMeetings)
	admin.Put("/meetings/:id/approve", handlers.AdminApproveMeeting)
	admin.Put("/meetings/:id/reject", handlers.AdminRejectMeeting)
	admin.Put("/meetings/:id", handlers.AdminEditMeeting)

	admin.Get("/matching/students", handlers.AdminStudentsByCity)
	admin.Post("/matching/coach", handlers.AdminMatchCoach)
	admin.Post("/matching/coordinator", handlers.AdminMatchCoordinator)

	admin.Post("/sms/send", handlers.AdminSendSMS)
	admin.Post("/sms/send-credentials/:id", handlers.AdminSendCredentialsSMS)
	admin.Get("/sms/logs", handlers.AdminSMSLogs)

	admin.Get("/reports/overview", handlers.AdminReportOverview)
	admin.Get("/reports/by-city", handlers.AdminReportByCity)
	admin.Get("/reports/missing-meetings", handlers.AdminReportMissingMeetings)
	admin.Get("/reports/mentor-performance", handlers.AdminReportMentorPerformance)

	admin.Get("/logs", handlers.AdminAuditLogs)
	admin.Get("/messages", handlers.AdminListMessages)
	admin.Post("/messages", handlers.AdminSendMessage)

	// Coach
	coach := api.Group("/coach",
		middleware.Auth(),
		middleware.RequireRole(models.RoleCoach),
	)
	coach.Get("/students", handlers.CoachListStudents)
	coach.Get("/weeks", handlers.CoachListWeeks)
	coach.Get("/meetings", handlers.CoachListMeetings)
	coach.Post("/meetings", handlers.CoachSubmitMeeting)
	coach.Put("/meetings/:id", handlers.CoachUpdateMeeting)
	coach.Get("/messages", handlers.CoachListMessages)
	coach.Post("/messages", handlers.CoachSendMessage)

	// Coordinator
	coord := api.Group("/coordinator",
		middleware.Auth(),
		middleware.RequireRole(models.RoleCoordinator),
	)
	coord.Get("/students", handlers.CoordinatorListStudents)
	coord.Get("/meetings", handlers.CoordinatorListMeetings)

	// Parent
	parent := api.Group("/parent",
		middleware.Auth(),
		middleware.RequireRole(models.RoleParent),
	)
	parent.Get("/meetings", handlers.ParentListMeetings)
	parent.Get("/messages", handlers.ParentListMessages)
	parent.Post("/messages", handlers.ParentSendMessage)

	// Student
	student := api.Group("/student",
		middleware.Auth(),
		middleware.RequireRole(models.RoleStudent),
	)
	student.Get("/profile", handlers.StudentGetProfile)
	student.Get("/messages", handlers.StudentListMessages)
	student.Post("/messages", handlers.StudentSendMessage)
}
