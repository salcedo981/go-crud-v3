package routers

import (
	"go_template_v3/pkg/middleware"
	ctrFeatureOne "go_template_v3/pkg/services/featureOne/controller"
	svcHealthcheck "go_template_v3/pkg/services/healthcheck"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func APIRoute(app *fiber.App) {
    app.Use("/uploads", static.New("./uploads", static.Config{
        MaxAge: 3600, // 1 hour cache
    }))

	publicV1 := app.Group("/api/public/v1")
	privateV1 := app.Group("/api/private/v1")

	// HealthCheck
	publicV1.Get("/", svcHealthcheck.HealthCheck)
	privateV1.Get("/", svcHealthcheck.HealthCheck)

	// Expense Category
	expenseCategoryEndpoint := publicV1.Group("/expense-categories")
	expenseCategoryEndpoint.Post("/", ctrFeatureOne.AddExpenseCategory)
	expenseCategoryEndpoint.Get("/", ctrFeatureOne.GetExpenseCategories)

	// Auth Routes
	authGroup := publicV1.Group("/auth")
	authGroup.Post("/register", ctrFeatureOne.Register)
	authGroup.Post("/login", ctrFeatureOne.Login)
	authGroup.Post("/forgot-password", ctrFeatureOne.ForgotPassword)
	authGroup.Post("/verify-reset-token", ctrFeatureOne.VerifyResetToken)
	authGroup.Post("/reset-password", ctrFeatureOne.ResetPassword)

	// Protected Auth Routes
	authGroupProtected := publicV1.Group("/auth", middleware.AuthMiddleware)
	authGroupProtected.Put("/update-user", ctrFeatureOne.UpdateUser)
	authGroupProtected.Post("/logout", ctrFeatureOne.Logout)

	// Protected expense routes
	expenseGroup := publicV1.Group("/expenses", middleware.AuthMiddleware)
	expenseGroup.Put("/batch", ctrFeatureOne.BatchUpdateExpenses)
	expenseGroup.Put("/batch-async", ctrFeatureOne.BatchUpdateExpensesAsync)
	expenseGroup.Get("/batch-async/:jobId", ctrFeatureOne.GetBatchJobStatus)

	expenseGroup.Post("/", ctrFeatureOne.AddExpense)
	expenseGroup.Post("/v2", ctrFeatureOne.AddExpenseV2)
	expenseGroup.Get("/", ctrFeatureOne.GetExpenses)
	expenseGroup.Get("/:id", ctrFeatureOne.GetExpense)
	expenseGroup.Delete("/:id", ctrFeatureOne.DeleteExpense)
	expenseGroup.Put("/:id", ctrFeatureOne.UpdateExpense)

	expenseGroup.Post("test-internal-send-request", ctrFeatureOne.TestInternalSendRequest)

}
