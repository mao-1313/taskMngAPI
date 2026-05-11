package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("データベースの初期化に失敗しました: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			deadlimit DATE NOT NULL DEFAULT (date('now', '+7 day')),
			status TEXT NOT NULL
		)
	`)

	if err != nil {
		t.Fatalf("テーブルの作成に失敗しました: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})
	return db
}

func TestHandleList_0rec(t *testing.T) {
	h := handler{db: setupDB(t)}

	req := httptest.NewRequest("GET", "/tasks", nil)
	rec := httptest.NewRecorder()
	h.handleList(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("期待されるステータスコードは %d ですが、実際は %d です", http.StatusOK, rec.Code)
	}
	body := strings.TrimSpace(rec.Body.String())
	if body != "[]" {
		t.Errorf("期待されるボディは[]ですが、実際は%sです", body)
	}

}

func TestHandleList_1rec(t *testing.T) {
	// DB初期化
	h := handler{db: setupDB(t)}

	// データ追加
	_, err := h.db.Exec(`INSERT INTO tasks (name, deadlimit, status)
		VALUES ('テストタスク', '2024-06-30', '未完了')`)
	if err != nil {
		t.Fatalf("データの追加に失敗しました: %v", err)
	}

	// データ取得
	req := httptest.NewRequest("GET", "/tasks", nil)
	rec := httptest.NewRecorder()
	h.handleList(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("期待されるステータスコードは %d ですが、実際は %d です", http.StatusOK, rec.Code)
	}

	// データ内容確認
	body := strings.TrimSpace(rec.Body.String())
	expected := `[{"id":1,"name":"テストタスク","deadlimit":"2024-06-30T00:00:00Z","status":"未完了"}]`
	if body != expected {
		t.Errorf("期待されるボディは%sですが、実際は%sです", expected, body)
	}
}

func TestHandleDone(t *testing.T) {
	// DB初期化
	h := handler{db: setupDB(t)}

	// データ追加
	_, err := h.db.Exec(`INSERT INTO tasks (name, deadlimit, status)
		VALUES ('テストタスク', '2024-06-30', 'Not Started')`)
	if err != nil {
		t.Fatalf("データの追加に失敗しました: %v", err)
	}

	// データ更新
	req := httptest.NewRequest("PATCH", "/tasks/1/done", nil)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()
	h.handleDone(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("期待されるステータスコードは %d ですが、実際は %d です", http.StatusOK, rec.Code)
	}

	// データ内容確認
	var status string
	err = h.db.QueryRow(`SELECT status FROM tasks WHERE id = 1`).Scan(&status)
	if err != nil {
		t.Fatalf("データの取得に失敗しました: %v", err)
	}
	if status != "done" {
		t.Errorf("期待されるステータスは done ですが、実際は %s です", status)
	}
}

func TestHandleDelete(t *testing.T) {
	// DB初期化
	h := handler{db: setupDB(t)}

	// データ追加
	_, err := h.db.Exec(`INSERT INTO tasks (name, deadlimit, status)
		VALUES ('テストタスク', '2024-06-30', 'Not Started')`)
	if err != nil {
		t.Fatalf("データの追加に失敗しました: %v", err)
	}

	// データ更新
	req := httptest.NewRequest("DELETE", "/tasks/1", nil)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()
	h.handleDelete(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("期待されるステータスコードは %d ですが、実際は %d です", http.StatusOK, rec.Code)
	}

	// データ件数確認
	count := 0
	h.db.QueryRow(`SELECT COUNT(*) FROM tasks`).Scan(&count)
	if count != 0 {
		t.Fatalf("データの削除に失敗しました")
	}
}
