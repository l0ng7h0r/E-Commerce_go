package main

import (
	"log"
	"time"
	_ "time/tzdata"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	adminHandler "github.com/l0ng7h0r/ecommerce/internal/handler/admin"
	sellerHandler "github.com/l0ng7h0r/ecommerce/internal/handler/seller"
	userHandler "github.com/l0ng7h0r/ecommerce/internal/handler/user"

	"github.com/l0ng7h0r/ecommerce/docs/swagger"
	"github.com/l0ng7h0r/ecommerce/internal/middleware"
	"github.com/l0ng7h0r/ecommerce/internal/repository"
	"github.com/l0ng7h0r/ecommerce/internal/usecase"
	"github.com/l0ng7h0r/ecommerce/pkg/config"
	"github.com/l0ng7h0r/ecommerce/pkg/database"
	"github.com/l0ng7h0r/ecommerce/pkg/phajay"
	"github.com/l0ng7h0r/ecommerce/pkg/supabase"

	"github.com/gofiber/contrib/v3/swaggo"
	_ "github.com/l0ng7h0r/ecommerce/docs/admin"
	_ "github.com/l0ng7h0r/ecommerce/docs/seller"
	_ "github.com/l0ng7h0r/ecommerce/docs/user"
)

func init() {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err == nil {
		time.Local = loc
	}
}

// @title           E-Commerce API Platform
// @version         2.0
// @description     Role-Based Microservices API Monorepo (User, Seller, Admin)
// @host            localhost:3000
// @BasePath        /api/v2

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg.DBDsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// --- Phajay Payment Gateway Client ---
	phajayClient := phajay.NewClient(cfg.PhajaySecretKey)

	// --- Supabase Storage Client ---
	supabaseClient := supabase.NewClient(cfg.SupabaseURL, cfg.SupabaseServiceRoleKey)
	_ = supabaseClient

	// --- Repositories ---
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// --- Usecases ---
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg)
	productUsecase := usecase.NewProductUsecase(productRepo)
	cartUsecase := usecase.NewCartUsecase(cartRepo, productRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, cartRepo, productRepo)
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepo, orderRepo, phajayClient)

	// --- Handlers ---
	// Customer Handlers
	uAuthH := userHandler.NewUserAuthHandler(authUsecase)
	uProdH := userHandler.NewUserProductHandler(productUsecase)
	uCartH := userHandler.NewUserCartHandler(cartUsecase)
	uOrderH := userHandler.NewUserOrderHandler(orderUsecase)
	uPayH := userHandler.NewUserPaymentHandler(paymentUsecase)

	// Seller Handlers
	sAuthH := sellerHandler.NewSellerAuthHandler(authUsecase)
	sProdH := sellerHandler.NewSellerProductHandler(productUsecase, supabaseClient)

	// Admin Handlers
	aAuthH := adminHandler.NewAdminAuthHandler(authUsecase)
	aUserH := adminHandler.NewAdminUserMgmtHandler(authUsecase)
	aCatH := adminHandler.NewAdminCategoryHandler(productUsecase)
	aOrderH := adminHandler.NewAdminOrderMgmtHandler(orderUsecase)

	// --- Middleware ---
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	app := fiber.New()
	app.Use(cors.New())

	// --- 3 Swagger Portals with Portal Navigation Buttons ---
	app.Get("/swagger/user/*", swaggo.New(swagger.GetUserConfig()))
	app.Get("/swagger/seller/*", swaggo.New(swagger.GetSellerConfig()))
	app.Get("/swagger/admin/*", swaggo.New(swagger.GetAdminConfig()))

	api := app.Group("/api/v2")

	// ── Public Webhooks ────────────────────────────────────────────────────────
	api.Post("/webhooks/phajay", uPayH.PhajayWebhook)

	// ═══════════════════════════════════════════════════════════════════════════
	// 1. CUSTOMER ROUTER PORTAL (/api/v2/user)
	// ═══════════════════════════════════════════════════════════════════════════
	uPortal := api.Group("/user")

	// Public Auth & Products
	uPortal.Post("/register", uAuthH.Register)
	uPortal.Post("/login", uAuthH.Login)
	uPortal.Post("/refresh", uAuthH.Refresh)
	uPortal.Post("/logout", uAuthH.Logout)

	uPortal.Get("/products", uProdH.GetAllProducts)
	uPortal.Get("/products/:id", uProdH.GetProductByID)
	uPortal.Get("/categories", uProdH.GetAllCategories)
	uPortal.Get("/products/seller/:sellerId", uProdH.GetProductsBySeller)

	// Customer Authenticated Routes
	uAuth := uPortal.Group("", authMiddleware.Auth())

	// Cart
	uAuth.Get("/cart", uCartH.GetCart)
	uAuth.Post("/cart/items", uCartH.AddItem)
	uAuth.Put("/cart/items/:productId", uCartH.UpdateItem)
	uAuth.Delete("/cart/items/:productId", uCartH.RemoveItem)
	uAuth.Delete("/cart", uCartH.ClearCart)

	// Orders & Order History
	uAuth.Post("/orders", uOrderH.CreateOrder)
	uAuth.Get("/orders", uOrderH.GetMyOrders)
	uAuth.Get("/orders/:id", uOrderH.GetOrderByID)

	// Payments
	uAuth.Post("/payments", uPayH.CreatePayment)
	uAuth.Get("/payments/order/:orderId", uPayH.GetPaymentByOrder)

	// ═══════════════════════════════════════════════════════════════════════════
	// 2. SELLER ROUTER PORTAL (/api/v2/seller)
	// ═══════════════════════════════════════════════════════════════════════════
	sPortal := api.Group("/seller")

	sPortal.Post("/login", sAuthH.Login)
	sPortal.Post("/refresh", sAuthH.Refresh)
	sPortal.Post("/logout", sAuthH.Logout)

	sAuth := sPortal.Group("", authMiddleware.Auth(), authMiddleware.RequireRole("seller"))

	sAuth.Post("/products", sProdH.CreateProduct)
	sAuth.Post("/products/upload-image", sProdH.UploadImage)
	sAuth.Get("/products", sProdH.GetMyProducts)
	sAuth.Put("/products/:id", sProdH.UpdateProduct)
	sAuth.Delete("/products/:id", sProdH.DeleteProduct)
	sAuth.Post("/categories", sProdH.CreateCategory)

	// ═══════════════════════════════════════════════════════════════════════════
	// 3. ADMIN ROUTER PORTAL (/api/v2/admin)
	// ═══════════════════════════════════════════════════════════════════════════
	aPortal := api.Group("/admin")

	aPortal.Post("/login", aAuthH.Login)
	aPortal.Post("/refresh", aAuthH.Refresh)
	aPortal.Post("/logout", aAuthH.Logout)

	aAuth := aPortal.Group("", authMiddleware.Auth(), authMiddleware.RequireRole("admin"))

	// User Management
	aAuth.Post("/users", aUserH.CreateUser)
	aAuth.Get("/users", aUserH.GetAllUsers)
	aAuth.Get("/users/:id", aUserH.GetUserByID)
	aAuth.Delete("/users/:id", aUserH.DeleteUser)

	// Category Management
	aAuth.Post("/categories", aCatH.CreateCategory)

	// Order Management
	aAuth.Get("/orders", aOrderH.GetAllOrders)
	aAuth.Patch("/orders/:id/status", aOrderH.UpdateOrderStatus)

	log.Printf("Server listening on port %s...", cfg.AppPort)
	log.Printf("User Swagger:   http://localhost:%s/swagger/user/index.html", cfg.AppPort)
	log.Printf("Seller Swagger: http://localhost:%s/swagger/seller/index.html", cfg.AppPort)
	log.Printf("Admin Swagger:  http://localhost:%s/swagger/admin/index.html", cfg.AppPort)

	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
