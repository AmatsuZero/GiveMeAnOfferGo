package main

import (
	"context"
	"golang.org/x/exp/slices"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"time"
)

type DownloadTaskUIItem struct {
	*ParserTask
	CreatedAt time.Time `json:"time"`
	UpdatedAt time.Time
	Status    string            `json:"status"`
	IsDone    bool              `json:"isDone"`
	VideoPath string            `json:"videoPath"`
	State     DownloadTaskState `json:"state"`
}

// 通知前端任务列表添加任务
func (a *App) addTaskNotifyItem(task *ParserTask) *DownloadTaskUIItem {
	// 先从记录里面查找
	for _, t := range a.tasks {
		if t.Url == task.Url {
			return t
		}
	}

	// 通知前端任务列表添加任务
	item := &DownloadTaskUIItem{}
	a.db.FirstOrCreate(item, &DownloadTaskUIItem{
		ParserTask: task,
		Status:     "初始化...",
		State:      DownloadTaskIdle,
	})
	a.tasks = append(a.tasks, item)
	a.eventsEmit(TaskNotifyCreate, item)
	return item
}

func (a *App) RemoveTaskNotifyItem(item *DownloadTaskUIItem) (err error) {
	err = a.removeLocalNotifyItem(item)
	if err != nil {
		return err
	}
	idx := slices.IndexFunc[*DownloadTaskUIItem](a.tasks, func(ele *DownloadTaskUIItem) bool {
		return ele.Url == item.Url
	})
	if idx == -1 {
		return
	}
	a.tasks = append(a.tasks[:idx], a.tasks[idx+1:]...)
	a.db.Where("url = ?", item.Url).Delete(item)
	return err
}

func (a *App) removeLocalNotifyItem(item *DownloadTaskUIItem) (err error) {
	// 先看看是否已经下载过是否存在
	if _, err = os.Stat(item.VideoPath); os.IsNotExist(err) {
		return nil
	}
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
	return
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

	for _, task := range a.tasks {
		a.eventsEmit(TaskNotifyCreate, task)
	}
}

func (a *App) saveTracks() {
	a.db.Save(a.tasks)
}

func (a *App) ClearTasks() {
	for _, task := range a.tasks {
		a.RemoveTaskNotifyItem(task)
	}
}
