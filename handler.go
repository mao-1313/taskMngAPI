package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type addRequest struct {
	Name string `json:"name"`
}

type deadlineRequest struct {
	Deadline string `json:"deadline"`
}

func handleList(w http.ResponseWriter, r *http.Request) {
	d, err := load()
	if err != nil {
		http.Error(w, "タスクの読み込みに失敗しました", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d.Tasks)
	fmt.Printf("タスクの一覧を返しました: %v\n", d.Tasks)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	//  タスクの追加処理
	// ファイルを読み込んでタスクの一覧を取得
	d, err := load()
	if err != nil {
		http.Error(w, "タスクの読み込みに失敗しました", http.StatusInternalServerError)
		return
	}

	// リクエスト内容をデコードして、構造体に格納
	var newReq addRequest
	err = json.NewDecoder(r.Body).Decode(&newReq)
	if err != nil {
		http.Error(w, "リクエストの内容が不正です", http.StatusBadRequest)
		return
	}

	// タスクを追加して、ファイルに保存
	d.addTask(newReq.Name)
	err = save(d)
	if err != nil {
		http.Error(w, "タスクの保存に失敗しました", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func handleDone(w http.ResponseWriter, r *http.Request) {
	// リクエストパスからIDを取得
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "IDの形式が不正です", http.StatusBadRequest)
		return
	}
	// タスクの完了処理
	// ファイルを読み込んでタスクの一覧を取得
	d, err := load()
	if err != nil {
		http.Error(w, "タスクの読み込みに失敗しました", http.StatusInternalServerError)
		return
	}

	// タスクを追加して、ファイルに保存
	if !d.Done(id) {
		http.Error(w, "タスクが見つかりませんでした", http.StatusNotFound)
		return
	}
	err = save(d)
	if err != nil {
		http.Error(w, "タスクの保存に失敗しました", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	// リクエストパスからIDを取得
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "IDの形式が不正です", http.StatusBadRequest)
		return
	}
	// タスクの完了処理
	// ファイルを読み込んでタスクの一覧を取得
	d, err := load()
	if err != nil {
		http.Error(w, "タスクの読み込みに失敗しました", http.StatusInternalServerError)
		return
	}

	// タスクを追加して、ファイルに保存
	if !d.Delete(id) {
		http.Error(w, "タスクが見つかりませんでした", http.StatusNotFound)
		return
	}
	err = save(d)
	if err != nil {
		http.Error(w, "タスクの保存に失敗しました", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleDeadline(w http.ResponseWriter, r *http.Request) {
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

	// タスクの完了処理
	// ファイルを読み込んでタスクの一覧を取得
	d, err := load()
	if err != nil {
		http.Error(w, "タスクの読み込みに失敗しました", http.StatusInternalServerError)
		return
	}

	// タスクを追加して、ファイルに保存
	if !d.Deadline(id, newReq.Deadline) {
		http.Error(w, "タスクが見つかりませんでした", http.StatusNotFound)
		return
	}
	err = save(d)
	if err != nil {
		http.Error(w, "タスクの保存に失敗しました", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
