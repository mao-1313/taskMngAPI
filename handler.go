package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type handler struct {
	db *sql.DB
}

type addRequest struct {
	Name string `json:"name"`
}

type deadlineRequest struct {
	Deadline string `json:"deadline"`
}

// タスクの定義
type task struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Limit  string `json:"deadlimit"`
	Status string `json:"status"`
}

func (h *handler) handleList(w http.ResponseWriter, r *http.Request) {
	// タスク一覧をJSON形式で返す
	// sqlite3のデータベースからタスクの一覧を取得
	recs, err := h.db.Query("SELECT id, name, deadlimit, status FROM tasks")
	if err != nil {
		http.Error(w, "タスクの取得に失敗しました", http.StatusInternalServerError)
		return
	}
	defer recs.Close()

	var tasks []task
	for recs.Next() {
		var t task
		err := recs.Scan(&t.ID, &t.Name, &t.Limit, &t.Status)
		if err != nil {
			http.Error(w, "タスクの取得に失敗しました", http.StatusInternalServerError)
			fmt.Printf("タスクの取得に失敗しました: %v\n", err)
			return
		}
		tasks = append(tasks, t)
	}

	json.NewEncoder(w).Encode(tasks)
	fmt.Printf("タスクの一覧を返しました: %v\n", tasks)
}

func (h *handler) handleAdd(w http.ResponseWriter, r *http.Request) {

	// リクエスト内容をデコードして、構造体に格納
	var newReq addRequest
	err := json.NewDecoder(r.Body).Decode(&newReq)
	if err != nil {
		http.Error(w, "リクエストの内容が不正です", http.StatusBadRequest)
		return
	}

	// タスクを追加して、sqlite3のdbに保存
	_, err = h.db.Exec("INSERT INTO tasks (name, status) VALUES (?, 'pending')", newReq.Name)
	if err != nil {
		http.Error(w, "タスクの保存に失敗しました", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *handler) handleDone(w http.ResponseWriter, r *http.Request) {
	// リクエストパスからIDを取得
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "IDの形式が不正です", http.StatusBadRequest)
		return
	}

	// IDに対応するタスクのステータスを"done"に更新
	result, err := h.db.Exec("UPDATE tasks SET status = 'done' WHERE id = ?", id)
	if err != nil {
		http.Error(w, "タスクの更新に失敗しました", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "タスクの更新に失敗しました", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "タスクが見つかりませんでした", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	// リクエストパスからIDを取得
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "IDの形式が不正です", http.StatusBadRequest)
		return
	}

	// IDに対応するタスクを削除
	result, err := h.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		http.Error(w, "タスクの削除に失敗しました", http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "タスクが見つかりませんでした", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) handleDeadline(w http.ResponseWriter, r *http.Request) {
	// リクエストパスからIDを取得
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "IDの形式が不正です", http.StatusBadRequest)
		return
	}

	// リクエスト内容をデコードして、構造体に格納
	var newReq deadlineRequest
	err = json.NewDecoder(r.Body).Decode(&newReq)
	if err != nil {
		http.Error(w, "リクエストの内容が不正です", http.StatusBadRequest)
		return
	}

	// newReq.Deadlineの形式チェック
	if _, err := time.Parse("2006-01-02", newReq.Deadline); err != nil {
		http.Error(w, "期限の形式が不正です。YYYY-MM-DD形式で指定してください", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec("UPDATE tasks SET deadlimit = ? WHERE id = ?", newReq.Deadline, id)
	if err != nil {
		http.Error(w, "タスクの更新に失敗しました", http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "タスクが見つかりませんでした", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	// タスク一覧をJSON形式で返す
	// sqlite3のデータベースからクエリの内容で検索してタスクの一覧を取得
	cond := r.URL.Query().Get("cond")
	if cond == "" {
		http.Error(w, "検索条件が指定されていません", http.StatusBadRequest)
		return
	}

	// クエリで指定された条件を解析して=で分割
	parts := strings.SplitN(cond, "=", 2)
	if len(parts) != 2 || (parts[0] != "id" && parts[0] != "name" && parts[0] != "status") {
		http.Error(w, "検索結果はfield=xxxの形式で指定してください", http.StatusBadRequest)
		return
	}

	field, val := parts[0], parts[1]

	tasks := []task{}

	var rows *sql.Rows
	var err error

	switch field {
	case "id":
		rows, err = h.db.Query("SELECT id, name, deadlimit, status FROM tasks WHERE id = ?", val)
	case "name":
		rows, err = h.db.Query("SELECT id, name, deadlimit, status FROM tasks WHERE name = ?", val)
	case "status":
		rows, err = h.db.Query("SELECT id, name, deadlimit, status FROM tasks WHERE status = ?", val)
	}
	if err != nil {
		http.Error(w, "タスクの取得に失敗しました", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var t task
		err := rows.Scan(&t.ID, &t.Name, &t.Limit, &t.Status)
		if err != nil {
			http.Error(w, "タスクの取得に失敗しました", http.StatusInternalServerError)
			fmt.Printf("タスクの取得に失敗しました: %v\n", err)
			return
		}
		tasks = append(tasks, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
	fmt.Printf("タスクの一覧を返しました: %v\n", tasks)
}
