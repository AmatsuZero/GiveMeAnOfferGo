package main

import (
	"context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"time"
)

type DownloadTaskUIItem struct {
	TaskName  string    `json:"taskName"`
	CreatedAt time.Time `json:"time"`
	UpdatedAt time.Time
	Status    string            `json:"status"`
	Url       string            `gorm:"primaryKey" json:"url"`
	IsDone    bool              `json:"isDone"`
	VideoPath string            `json:"videoPath"`
	State     DownloadTaskState `json:"state"`
}

// 通知前端任务列表添加任务
func (a *App) addTaskNotifyItem(task *ParserTask) *DownloadTaskUIItem {
	// 先从记录里面查找
	idx, ok := a.tasksIdx[task.Url]
	if ok {
		return a.tasks[idx]
	}
	// 通知前端任务列表添加任务
	item := &DownloadTaskUIItem{
		TaskName: task.TaskName,
		Status:   "初始化...",
		Url:      task.Url,
		State:    DownloadTaskIdle,
	}
	a.tasks = append(a.tasks, item)
	a.tasksIdx[task.Url] = len(a.tasks) - 1
	a.eventsEmit(TaskNotifyCreate, item)
	return item
}

func (a *App) RemoveTaskNotifyItem(item *DownloadTaskUIItem) (err error) {
	if !item.IsDone {
		task, ok := a.queues[item.Url]
		if ok {
			task.Stop()
		}
	}
	if item.IsDone {
		err = os.Remove(item.VideoPath)
	} else {
		err = os.RemoveAll(item.VideoPath)
	}
	a.db.Delete(item)
	return err
}

func (a *App) initDB() {
	dbPath := filepath.Join(appFolder, "records.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		QueryFields: true,
	})
	if err != nil {
		a.logErrorf("初始化数据库失败: %v", err)
		return
	}

	a.db = db
	err = a.db.AutoMigrate(&DownloadTaskUIItem{})
	if err != nil {
		a.logErrorf("数据库自动迁移失败：%v", err)
	}
	a.tasks = make([]*DownloadTaskUIItem, 0)
	a.queues = map[string]StoppableTask{}
}

func (a *App) domReady(ctx context.Context) {
	result := a.db.Find(&a.tasks)
	a.logInfof("读取数据 %v 条", result.RowsAffected)
	a.tasksIdx = make(map[string]int, result.RowsAffected)

	for i, task := range a.tasks {
		a.tasksIdx[task.Url] = i
		a.eventsEmit(TaskNotifyCreate, task)
	}
}

func (a *App) saveTracks() {
	a.db.Save(a.tasks)
}
