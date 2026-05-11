package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// ファイル全体の定義
type define struct {
	MaxID int    `json:"max_id"`
	Tasks []task `json:"tasks"`
}

// タスクの定義
type task struct {
	ID     int       `json:"id"`
	Name   string    `json:"name"`
	Limit  time.Time `json:"limit"`
	Status string    `json:"status"`
}

const (
	condID = iota + 1
	condName
	condStatus
)

func load() (*define, error) {
	var def define
	file, err := os.OpenFile("def.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("ファイルを開けませんでした: %w", err)

	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&def)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("ファイルの内容を読み込めませんでした。: %w", err)
	}

	return &def, nil
}

func save(def *define) error {
	file, err := os.OpenFile("def.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("ファイルを開けませんでした: %w", err)

	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(def)
	if err != nil {
		return fmt.Errorf("ファイルに保存できませんでした。: %w", err)
	}

	return nil
}

// タスクを追加する関数
// define構造体に対するメソッドを定義

func (d *define) addTask(name string) {

	d.MaxID++

	newTask := task{
		ID:     d.MaxID,
		Name:   name,
		Limit:  time.Now().AddDate(0, 0, 7),
		Status: "NotStarted",
	}
	d.Tasks = append(d.Tasks, newTask)
	fmt.Println("タスクが追加されました。")
}

func (d *define) Search(cond string) {
	// condをfield=xxxの形式で受け取る
	parts := strings.SplitN(cond, "=", 2)
	if len(parts) != 2 || (parts[0] != "id" && parts[0] != "name" && parts[0] != "status") {
		fmt.Println("検索条件はfield=xxxの形式で指定してください")
		return
	}

	condType := 0
	iVal := 0
	sVal := ""
	var err error

	switch parts[0] {
	case "id":
		condType = condID
		iVal, err = strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("idは整数で指定してください")
			return
		}
	case "name":
		condType = condName
		sVal = parts[1]
	case "status":
		condType = condStatus
		sVal = parts[1]
	}

	for _, task := range d.Tasks {
		switch condType {
		case condID:
			if task.ID == iVal {
				fmt.Printf("id=%d name=%s limit=%s status=%s \n",
					task.ID, task.Name, task.Limit.Format("2006-01-02"), task.Status)
			}

		case condName:
			if task.Name == sVal {
				fmt.Printf("id=%d name=%s limit=%s status=%s \n",
					task.ID, task.Name, task.Limit.Format("2006-01-02"), task.Status)
			}
		case condStatus:
			if task.Status == sVal {
				fmt.Printf("id=%d name=%s limit=%s status=%s \n",
					task.ID, task.Name, task.Limit.Format("2006-01-02"), task.Status)
			}
		}
	}
}

func (d *define) Done(id int) bool {
	for i, task := range d.Tasks {
		if id == task.ID {
			d.Tasks[i].Status = "Done"
			return true
		}
	}
	return false
}

func (d *define) Delete(id int) bool {
	for i, task := range d.Tasks {
		if id == task.ID {
			d.Tasks = append(d.Tasks[:i], d.Tasks[i+1:]...)
			return true
		}
	}
	return false
}

func (d *define) Deadline(id int, limitStr string) bool {
	for i, task := range d.Tasks {
		if id == task.ID {
			limit, err := time.Parse("20060102", limitStr)
			if err != nil {
				fmt.Println(err)
				return false
			}
			d.Tasks[i].Limit = limit
			return true
		}
	}
	return false
}
