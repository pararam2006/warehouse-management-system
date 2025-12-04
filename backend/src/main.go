package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"warehouse-management-system/src/config"
	"warehouse-management-system/src/controllers"
	"warehouse-management-system/src/middleware"
	"warehouse-management-system/src/repositories"
	"warehouse-management-system/src/services"

	"github.com/gorilla/mux"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	dsn := "file:" + cfg.DBPath + "?_pragma=foreign_keys(ON)"

	log.Println("SQLite DSN:", dsn)
	db, err := config.OpenSQLite(dsn)
	if err != nil {
		log.Fatalf("failed to open SQLite database: %v", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	// Выполняем миграцию схемы (idempotent).
	if err := config.MigrateSQLite(db); err != nil {
		log.Fatalf("failed to migrate SQLite schema: %v", err)
	}

	// Инициализация репозиториев (слой хранения данных)
	userRepo := repositories.NewUserRepository(db)
	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	supplierRepo := repositories.NewSupplierRepository(db)
	warehouseRepo := repositories.NewWarehouseRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	// Инициализация сервисов
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	productService := services.NewProductService(productRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	supplierService := services.NewSupplierService(supplierRepo)
	warehouseService := services.NewWarehouseService(warehouseRepo, productRepo)
	orderService := services.NewOrderService(orderRepo, warehouseRepo, productRepo)

	// Инициализация контроллеров
	authController := controllers.NewAuthController(authService)
	productController := controllers.NewProductController(productService)
	categoryController := controllers.NewCategoryController(categoryService)
	supplierController := controllers.NewSupplierController(supplierService)
	warehouseController := controllers.NewWarehouseController(warehouseService)
	orderController := controllers.NewOrderController(orderService)

	// Инициализация роутера
	router := mux.NewRouter()

	// Middleware
	router.Use(middleware.CORS)
	router.Use(middleware.LoggingMiddleware)

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// Auth routes
	api.HandleFunc("/auth/register", authController.Register).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/login", authController.Login).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/me", middleware.AuthMiddleware(authController.GetMe, cfg.JWTSecret)).Methods("GET", "OPTIONS")

	// Products routes
	api.HandleFunc("/products", middleware.AuthMiddleware(productController.GetProducts, cfg.JWTSecret)).Methods("GET", "OPTIONS")
	api.HandleFunc("/products", middleware.AuthMiddleware(middleware.RoleMiddleware(productController.CreateProduct, "admin", "manager"), cfg.JWTSecret)).Methods("POST", "OPTIONS")
	api.HandleFunc("/products/{id}", middleware.AuthMiddleware(productController.GetProduct, cfg.JWTSecret)).Methods("GET", "OPTIONS")
	api.HandleFunc("/products/{id}", middleware.AuthMiddleware(middleware.RoleMiddleware(productController.UpdateProduct, "admin", "manager"), cfg.JWTSecret)).Methods("PUT", "OPTIONS")
	api.HandleFunc("/products/{id}", middleware.AuthMiddleware(middleware.RoleMiddleware(productController.DeleteProduct, "admin"), cfg.JWTSecret)).Methods("DELETE", "OPTIONS")

	// Categories routes
	api.HandleFunc("/categories", middleware.AuthMiddleware(categoryController.GetCategories, cfg.JWTSecret)).Methods("GET", "OPTIONS")
	api.HandleFunc("/categories", middleware.AuthMiddleware(middleware.RoleMiddleware(categoryController.CreateCategory, "admin", "manager"), cfg.JWTSecret)).Methods("POST", "OPTIONS")
	api.HandleFunc("/categories/{id}", middleware.AuthMiddleware(categoryController.GetCategory, cfg.JWTSecret)).Methods("GET", "OPTIONS")
	api.HandleFunc("/categories/{id}", middleware.AuthMiddleware(middleware.RoleMiddleware(categoryController.UpdateCategory, "admin", "manager"), cfg.JWTSecret)).Methods("PUT", "OPTIONS")
	api.HandleFunc("/categories/{id}", middleware.AuthMiddleware(middleware.RoleMiddleware(categoryController.DeleteCategory, "admin"), cfg.JWTSecret)).Methods("DELETE", "OPTIONS")

	// Suppliers routes
	api.HandleFunc("/suppliers", middleware.AuthMiddleware(supplierController.GetSuppliers, cfg.JWTSecret)).Methods("GET", "OPTIONS")
	api.HandleFunc("/suppliers", middleware.AuthMiddleware(middleware.RoleMiddleware(supplierController.CreateSupplier, "admin", "manager"), cfg.JWTSecret)).Methods("POST", "OPTIONS")
	api.HandleFunc("/suppliers/{id}", middleware.AuthMiddleware(supplierController.GetSupplier, cfg.JWTSecret)).Methods("GET", "OPTIONS")
	api.HandleFunc("/suppliers/{id}", middleware.AuthMiddleware(middleware.RoleMiddleware(supplierController.UpdateSupplier, "admin", "manager"), cfg.JWTSecret)).Methods("PUT", "OPTIONS")
	api.HandleFunc("/suppliers/{id}", middleware.AuthMiddleware(middleware.RoleMiddleware(supplierController.DeleteSupplier, "admin"), cfg.JWTSecret)).Methods("DELETE", "OPTIONS")

	// Warehouse operations routes
	api.HandleFunc("/warehouse/receipt", middleware.AuthMiddleware(middleware.RoleMiddleware(warehouseController.Receipt, "admin", "manager", "storekeeper"), cfg.JWTSecret)).Methods("POST", "OPTIONS")
	api.HandleFunc("/warehouse/write-off", middleware.AuthMiddleware(middleware.RoleMiddleware(warehouseController.WriteOff, "admin", "manager", "storekeeper"), cfg.JWTSecret)).Methods("POST", "OPTIONS")
	api.HandleFunc("/warehouse/reserve", middleware.AuthMiddleware(middleware.RoleMiddleware(warehouseController.Reserve, "admin", "manager"), cfg.JWTSecret)).Methods("POST", "OPTIONS")
	api.HandleFunc("/warehouse/inventory", middleware.AuthMiddleware(warehouseController.GetInventory, cfg.JWTSecret)).Methods("GET", "OPTIONS")

	// Orders routes
	api.HandleFunc("/orders", middleware.AuthMiddleware(orderController.GetOrders, cfg.JWTSecret)).Methods("GET", "OPTIONS")
	api.HandleFunc("/orders", middleware.AuthMiddleware(middleware.RoleMiddleware(orderController.CreateOrder, "admin", "manager"), cfg.JWTSecret)).Methods("POST", "OPTIONS")
	api.HandleFunc("/orders/{id}", middleware.AuthMiddleware(orderController.GetOrder, cfg.JWTSecret)).Methods("GET", "OPTIONS")
	api.HandleFunc("/orders/{id}/status", middleware.AuthMiddleware(middleware.RoleMiddleware(orderController.UpdateOrderStatus, "admin", "manager"), cfg.JWTSecret)).Methods("PUT", "OPTIONS")

	// Отдача страниц фронтенда (пути относительно корня проекта).
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/public/index.html")
	}).Methods("GET")
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/public/register.html")
	}).Methods("GET")
	router.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/public/dashboard.html")
	}).Methods("GET")
	router.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/public/products.html")
	}).Methods("GET")
	router.HandleFunc("/warehouse", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/public/warehouse.html")
	}).Methods("GET")
	router.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/public/orders.html")
	}).Methods("GET")

	// Отдача статических JS файлов из frontend/src/
	router.PathPrefix("/src/").Handler(http.StripPrefix("/src/", http.FileServer(http.Dir("../frontend/src/"))))

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
