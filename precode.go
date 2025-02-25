package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта.

// Обработчик для получения всех задач.
func getTasks(w http.ResponseWriter, r *http.Request) {

	// устанавливаем заголовок типа содержимого.
	w.Header().Set("Content-Type", "application/json")

	// сериализируем данные в json.
	jsonTasks, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // возвращение ошибки в случае неудачи сериализации.
		return
	}
	w.WriteHeader(http.StatusOK) // возвращем статус ответа.
	w.Write(jsonTasks)           // возвращаем ответ в формате json.
}

// Обработчик для отправки задачи на сервер.
func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// проверка наличия задачи по id
	if task.ID == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка на существование задачи с таким же ID
	if _, exists := tasks[task.ID]; exists {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// Обработчик задач по заданному ID.
func getTask(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id") // извлечение id.
	w.Header().Set("Content-Type", "application/json")

	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Задачи с таким id нет!", http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Ошибка при кодировании задачи", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

// Обработчик удаления задач по заданному id.
func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // извлечение задачи по id.
	_, exists := tasks[id]
	if !exists {
		http.Error(w, "Задачи с таким id нет!", http.StatusBadRequest)
		return
	}

	delete(tasks, id) // удаляем задачу по id.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики.
	// определяем эндроинты и методы.
	r.Get("/tasks", getTasks)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
