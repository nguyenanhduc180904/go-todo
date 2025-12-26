package main

import (
	"database/sql"
	"fmt"
	"go-todo/internal/handler"
	"go-todo/internal/store"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import driver (quan trọng!)

	_ "go-todo/docs" // Import docs

	httpSwagger "github.com/swaggo/http-swagger"

	_ "go-todo/docs" // Import side-effect để load cấu hình Swagger
)

// @title Todo List API
// @version 1.0
// @description API quản lý danh sách công việc (Todo List) đơn giản bằng Golang.
// @host localhost:8080
// @BasePath /
func main() {
	_ = godotenv.Load()
	// 1. Lấy connection string từ biến môi trường (Chuẩn Cloud)
	connStr := os.Getenv("DB_CONN_STR")
	if connStr == "" {
		// Giá trị mặc định cho local dev
		connStr = ""
	}

	// 2. Kết nối Database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Kiểm tra kết nối
	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}
	fmt.Println("Connected to Database successfully!")

	// 3. Khởi tạo Store và Handler
	pgStore := store.NewPostgresStore(db)

	// Tự động tạo bảng nếu chưa có
	if err := pgStore.InitSchema(); err != nil {
		log.Fatal("Failed to init schema:", err)
	}

	todoHandler := handler.NewTodoHandler(pgStore)

	// 4. Setup Router (ServeMux)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /todos", todoHandler.CreateTodo)
	mux.HandleFunc("GET /todos", todoHandler.ListTodos)
	mux.HandleFunc("GET /todos/{id}", todoHandler.GetTodo)
	mux.HandleFunc("PUT /todos/{id}", todoHandler.UpdateTodo)
	mux.HandleFunc("DELETE /todos/{id}", todoHandler.DeleteTodo)

	// THÊM ROUTE CHO SWAGGER
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// 5. Chạy Server
	port := os.Getenv("PORT") // Cloud thường cấp PORT qua biến môi trường
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running on port " + port)
	http.ListenAndServe(":"+port, mux)
}
