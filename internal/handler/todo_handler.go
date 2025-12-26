package handler

import (
	"encoding/json"
	"go-todo/internal/model"
	"go-todo/internal/store"
	"net/http"
	"strconv"
	"time"
)

type TodoHandler struct {
	Store *store.PostgresStore
}

func NewTodoHandler(s *store.PostgresStore) *TodoHandler {
	return &TodoHandler{Store: s}
}

// Helper để gửi JSON response nhanh
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// CreateTodo godoc
// @Summary Tạo công việc mới
// @Description Tạo một Todo với title và status mặc định
// @Tags todos
// @Accept  json
// @Produce  json
// @Param todo body model.Todo true "Dữ liệu Todo"
// @Success 201 {object} model.Todo
// @Failure 400 {string} string "Invalid input"
// @Router /todos [post]
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var t model.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Tạo ID
	t.ID = strconv.FormatInt(time.Now().UnixNano(), 10)

	// Cập nhật: Xử lý lỗi từ DB
	if err := h.Store.Create(t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, t)
}

// ListTodos godoc
// @Summary Lấy danh sách công việc
// @Description Trả về danh sách tất cả các công việc hiện có
// @Tags todos
// @Produce  json
// @Success 200 {array} model.Todo
// @Router /todos [get]
func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
	// Cập nhật: Nhận data và error
	todos, err := h.Store.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Nếu mảng rỗng thì trả về mảng rỗng thay vì null
	if todos == nil {
		todos = []model.Todo{}
	}
	jsonResponse(w, http.StatusOK, todos)
}

// GetTodo godoc
// @Summary Lấy chi tiết công việc
// @Description Lấy thông tin công việc dựa trên ID
// @Tags todos
// @Produce  json
// @Param id path string true "Todo ID"
// @Success 200 {object} model.Todo
// @Failure 404 {string} string "Not Found"
// @Router /todos/{id} [get]
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") // Go 1.22 feature: Lấy tham số từ URL cực dễ

	t, err := h.Store.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, t)
}

// UpdateTodo godoc
// @Summary Cập nhật công việc
// @Description Cập nhật thông tin dựa trên ID
// @Tags todos
// @Accept  json
// @Produce  json
// @Param id path string true "Todo ID"
// @Param todo body model.Todo true "Dữ liệu cập nhật"
// @Success 200 {object} model.Todo
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Not Found"
// @Router /todos/{id} [put]
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var t model.Todo

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.Store.Update(id, t); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Trả về object đã update
	t.ID = id
	jsonResponse(w, http.StatusOK, t)
}

// DeleteTodo godoc
// @Summary Xóa công việc
// @Description Xóa vĩnh viễn công việc dựa trên ID
// @Tags todos
// @Param id path string true "Todo ID"
// @Success 204 "No Content"
// @Failure 404 {string} string "Not Found"
// @Router /todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.Store.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
