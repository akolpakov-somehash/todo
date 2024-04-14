package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	filename = "tasks.json"
	dirname  = ".todo"
)

func init() {
	homedir, _ := os.UserHomeDir()
	path := filepath.Join(homedir, dirname, filename)
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		storage := StorageJson{
			filepath: path,
		}
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			panic(err)
		}
		_, err = os.Create(path)
		if err != nil {
			panic(err)
		}
		initData := StorageData{
			Increment: 1,
			Tasks:     []*Task{},
		}
		err = storage.saveData(&initData)
		if err != nil {
			panic(err)
		}
	}
}

func NewStorageJson() (Storage, error) {
	homedir, _ := os.UserHomeDir()
	path := filepath.Join(homedir, dirname, filename)
	storage := StorageJson{
		filepath: path,
	}
	return &storage, nil
}

type StorageData struct {
	Increment uint    `json:"increment"`
	Tasks     []*Task `json:"tasks"`
}

func (t *Task) MarshalJSON() (jsn []byte, err error) {
	m := map[string]interface{}{
		"id":          t.ID,
		"description": t.Description,
		"status":      t.Status,
		"dateCreated": t.DateCreated.Format(time.RFC3339),
		"dateUpdated": t.DateUpdated.Format(time.RFC3339),
	}
	jsn, err = json.Marshal(m)
	return
}

type StorageJson struct {
	filepath string
}

func (s *StorageJson) AddTask(description string) error {
	data, error := s.loadData()
	if error != nil {
		return error
	}
	task := Task{
		ID:          data.Increment,
		Description: description,
		Status:      StatusPending,
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
	}
	data.Tasks = append(data.Tasks, &task)
	data.Increment++
	return s.saveData(data)
}

func (s *StorageJson) CompleteTask(id uint) error {
	data, error := s.loadData()
	if error != nil {
		return error
	}
	for _, task := range data.Tasks {
		if task.ID == id {
			task.Status = StatusCompleted
			task.DateUpdated = time.Now()
			break
		}
	}
	return s.saveData(data)
}

func (s *StorageJson) ListTasks(filter ListFilter) ([]*Task, error) {
	data, error := s.loadData()
	if error != nil {
		return nil, error
	}
	if filter == FilterAll {
		return data.Tasks, nil
	}
	var tasks []*Task
	for _, task := range data.Tasks {
		if filter == FilterCompleted && task.Status == StatusCompleted {
			tasks = append(tasks, task)
		} else if filter == FilterPending && task.Status == StatusPending {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (s *StorageJson) DeleteTask(id uint) error {
	data, error := s.loadData()
	if error != nil {
		return error
	}
	for i, task := range data.Tasks {
		if task.ID == id {
			data.Tasks = append((data.Tasks)[:i], (data.Tasks)[i+1:]...)
			break
		}
	}
	return s.saveData(data)
}

func (s *StorageJson) loadData() (*StorageData, error) {
	data, error := os.ReadFile(s.filepath)
	if error != nil {
		return nil, error
	}
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	var storageData StorageData
	error = decoder.Decode(&storageData)
	return &storageData, error
}

func (s *StorageJson) saveData(data *StorageData) error {
	file, err := os.OpenFile(s.filepath, os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer file.Close()
		encoder := json.NewEncoder(file)
		err = encoder.Encode(data)
	}
	return err
}
