package main

import (
	"fmt"
	"net/http"
)

func main() {
	// ここではサーバー起動のみを実行
	http.HandleFunc("GET /tasks", handleList)
	http.HandleFunc("POST /task", handleAdd)
	http.HandleFunc("PATCH /tasks/{id}/done", handleDone)
	http.HandleFunc("DELETE /tasks/{id}", handleDelete)
	http.HandleFunc("PATCH /tasks/{id}/deadline", handleDeadline)

	fmt.Println("サーバーを起動します...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("サーバーの起動に失敗しました: %v\n", err)
	}
}
