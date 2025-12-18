package integration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"example.com/pz16/internal/db"
	"example.com/pz16/internal/httpapi"
	"example.com/pz16/internal/repo"
	"example.com/pz16/internal/service"
)

func newServer(t *testing.T, dsn string) *httptest.Server {
	t.Helper()

	dbx, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	// Применяем миграции
	db.MustApplyMigrations(dbx)

	// Очистка БД и закрытие соединения после теста
	t.Cleanup(func() {
		_, _ = dbx.Exec(`TRUNCATE TABLE notes RESTART IDENTITY CASCADE`)
		_ = dbx.Close()
	})

	r := gin.Default()
	svc := service.Service{
		Notes: repo.NoteRepo{DB: dbx},
	}
	httpapi.Router{Svc: &svc}.Register(r)

	return httptest.NewServer(r)
}

func TestCreateAndGetNote(t *testing.T) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		t.Skip("DB_DSN not set (use `make up` and `make test`)")
	}

	srv := newServer(t, dsn)
	defer srv.Close()

	// ---------- 1) CREATE ----------
	resp, err := http.Post(
		srv.URL+"/notes",
		"application/json",
		strings.NewReader(`{"title":"Hello","content":"World"}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()

	var created map[string]any
	if err := json.Unmarshal(body, &created); err != nil {
		t.Fatal(err)
	}

	id, ok := created["id"].(float64)
	if !ok {
		t.Fatalf("id not found in response: %v", created)
	}

	// ---------- 2) GET ----------
	url := fmt.Sprintf("%s/notes/%d", srv.URL, int64(id))
	resp2, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp2.StatusCode)
	}

	body2, _ := io.ReadAll(resp2.Body)
	_ = resp2.Body.Close()

	var got map[string]any
	if err := json.Unmarshal(body2, &got); err != nil {
		t.Fatal(err)
	}

	if got["title"] != "Hello" {
		t.Fatalf("unexpected title: %v", got["title"])
	}
	if got["content"] != "World" {
		t.Fatalf("unexpected content: %v", got["content"])
	}
}

func TestGetNoteNotFound(t *testing.T) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		t.Skip("DB_DSN not set (use `make up` and `make test`)")
	}

	srv := newServer(t, dsn)
	defer srv.Close()

	resp, err := http.Get(fmt.Sprintf("%s/notes/999999", srv.URL))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}
