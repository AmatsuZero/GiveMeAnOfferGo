package app

import (
	"GiveMeAnOffer/downloader"
	"GiveMeAnOffer/eventbus"
	"GiveMeAnOffer/parse"
	"context"
	"golang.org/x/exp/slices"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

// 通知前端任务列表添加任务
func (a *App) addTaskNotifyItem(task *parse.ParserTask) *downloader.DownloadTaskUIItem {
	item := &downloader.DownloadTaskUIItem{
		ParserTask: task,
		Status:     "初始化...",
		State:      downloader.DownloadTaskIdle,
	}

	if a.isCliMode() {
		a.tasks = append(a.tasks, item)
		a.EventsEmit(eventbus.TaskNotifyCreate, item)
		return item
	}

	// 先从记录里面查找
	for _, t := range a.tasks {
		if t.Url == task.Url {
			a.EventsEmit(eventbus.TaskAddEvent, parse.EventMessage{
				Code:    -1,
				Message: "任务已经存在",
			})
			return t
		}
	}

	// 通知前端任务列表添加任务
	a.db.Create(item)
	a.tasks = append(a.tasks, item)
	a.EventsEmit(eventbus.TaskNotifyCreate, item)
	return item
}

func (a *App) RemoveTaskNotifyItem(item *downloader.DownloadTaskUIItem) (err error) {
	err = a.removeLocalNotifyItem(item)
	if err != nil {
		return err
	}
	idx := slices.IndexFunc(a.tasks, func(ele *downloader.DownloadTaskUIItem) bool {
		return ele.Url == item.Url
	})
	if idx == -1 {
		return
	}
	a.tasks = append(a.tasks[:idx], a.tasks[idx+1:]...)
	a.db.Where("url = ?", item.Url).Delete(item)
	return err
}

func (a *App) removeLocalNotifyItem(item *downloader.DownloadTaskUIItem) (err error) {
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
		a.LogErrorf("初始化数据库失败: %v", err)
		return
	}

	a.db = db
	err = a.db.AutoMigrate(&downloader.DownloadTaskUIItem{})
	if err != nil {
		a.LogErrorf("数据库自动迁移失败：%v", err)
	}
	a.tasks = make([]*downloader.DownloadTaskUIItem, 0)
	a.queues = map[string]downloader.StoppableTask{}
}

func (a *App) DomReady(_ context.Context) {
	result := a.db.Find(&a.tasks)
	a.LogInfof("读取数据 %v 条", result.RowsAffected)

	for _, task := range a.tasks {
		a.EventsEmit(eventbus.TaskNotifyCreate, task)
	}
}

func (a *App) saveTracks() {
	if len(a.tasks) > 0 && !a.isCliMode() {
		a.db.Save(a.tasks)
	}
}

func (a *App) ClearTasks() {
	for _, task := range a.tasks {
		if e := a.RemoveTaskNotifyItem(task); e != nil {
			a.LogErrorf("移除任务失败: %v", e)
		}
	}
}
