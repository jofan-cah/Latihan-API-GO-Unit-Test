package main

import (
	"api-go/database"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func setupMockDB() (sqlmock.Sqlmock, *sql.DB) {
	// Membuat mock database
	db, mock, err := sqlmock.New() // mock DB
	if err != nil {
		panic(err)
	}
	database.DB = db // Mengatur DB global untuk mock
	return mock, db
}

func TestAPIEndpoints(t *testing.T) {
	// Setup mock DB
	mock, db := setupMockDB()
	defer db.Close()

	t.Run("Test POST /users/create", func(t *testing.T) {
		// Data payload yang dikirim melalui API
		payload := `{"name":"Test User", "email":"test@example.com"}`
		req, _ := http.NewRequest("POST", "/users/create", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")

		// Mock query untuk INSERT ke database
		mock.ExpectExec("INSERT INTO users").
			WithArgs("Test User", "test@example.com").
			WillReturnResult(sqlmock.NewResult(1, 1)) // Mock hasil insert dengan ID 1

		// Membuat recorder untuk mencatat respons HTTP
		rr := httptest.NewRecorder()

		// Memanggil handler dengan request yang sudah dibuat
		handler := http.HandlerFunc(createUser)
		handler.ServeHTTP(rr, req)

		// Memastikan status code adalah 200 OK
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("POST /users/create failed: got %v want %v", status, http.StatusOK)
		} else {
			t.Log("API POST /users/create Success") // Menambahkan keterangan sukses
		}

		// Memeriksa apakah nama yang dikirim dalam respons sesuai dengan yang diharapkan
		expected := `"name":"Test User"`
		if !strings.Contains(rr.Body.String(), expected) {
			t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}

		// Memastikan bahwa semua expectation yang ditentukan pada mock DB terpenuhi
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unmet expectations: %v", err)
		}
	})

	t.Run("Test GET /users", func(t *testing.T) {
		// Mock query untuk SELECT dari database
		rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at"}).
			AddRow(1, "Test User", "test@example.com", "2024-01-01")

		mock.ExpectQuery("SELECT id, name, email, created_at FROM users").
			WillReturnRows(rows)

		// Membuat request GET
		req, _ := http.NewRequest("GET", "/users", nil)

		// Membuat recorder untuk mencatat respons HTTP
		rr := httptest.NewRecorder()

		// Memanggil handler dengan request GET
		handler := http.HandlerFunc(getUsers)
		handler.ServeHTTP(rr, req)

		// Memastikan status code adalah 200 OK
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("GET /users failed: got %v want %v", status, http.StatusOK)
		} else {
			t.Log("API GET /users Success") // Menambahkan keterangan sukses
		}

		// Memastikan body respons mengandung nama pengguna yang diharapkan
		expected := `"name":"Test User"`
		if !strings.Contains(rr.Body.String(), expected) {
			t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}

		// Memastikan bahwa semua expectation yang ditentukan pada mock DB terpenuhi
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unmet expectations: %v", err)
		}
	})

	t.Run("Test DELETE /users/{id}", func(t *testing.T) {
		// Mock query untuk DELETE dari database
		mock.ExpectExec("DELETE FROM users WHERE id = ?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1)) // Mock hasil delete dengan 1 baris terpengaruh

		// Membuat request DELETE
		req, _ := http.NewRequest("DELETE", "/users/1", nil)

		// Membuat recorder untuk mencatat respons HTTP
		rr := httptest.NewRecorder()

		// Memanggil handler dengan request DELETE
		handler := http.HandlerFunc(deleteUser)
		handler.ServeHTTP(rr, req)

		// Memastikan status code adalah 200 OK
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("DELETE /users/{id} failed: got %v want %v", status, http.StatusOK)
		} else {
			t.Log("API DELETE /users/{id} Success") // Menambahkan keterangan sukses
		}

		// Memastikan bahwa semua expectation yang ditentukan pada mock DB terpenuhi
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unmet expectations: %v", err)
		}
	})
}