package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// データベースの初期化
	db, err := sql.Open("sqlite3", "tasks.db")
	if err != nil {
		fmt.Printf("データベースの初期化に失敗しました: %v\n", err)
		return
	}
	defer db.Close()

	h := handler{db: db}

	// テーブルの作成
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			deadlimit DATE NOT NULL DEFAULT (date('now', '+7 day')),
			status TEXT NOT NULL
		)
	`)

	if err != nil {
		fmt.Printf("テーブルの作成に失敗しました: %v\n", err)
		return
	}

	// ここではサーバー起動のみを実行
	http.HandleFunc("GET /tasks", h.handleList)
	http.HandleFunc("POST /task", h.handleAdd)
	http.HandleFunc("PATCH /tasks/{id}/done", h.handleDone)
	http.HandleFunc("DELETE /tasks/{id}", h.handleDelete)
	http.HandleFunc("PATCH /tasks/{id}/deadline", h.handleDeadline)
	http.HandleFunc("GET /tasks/search", h.handleSearch)

	fmt.Println("サーバーを起動します...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("サーバーの起動に失敗しました: %v\n", err)
	}
}
