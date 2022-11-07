package downloader

import (
	"GiveMeAnOffer/parse"
	"time"
)

type DownloadTaskUIItem struct {
	*parse.ParserTask
	CreatedAt time.Time `json:"time"`
	UpdatedAt time.Time
	Status    string            `json:"status"`
	IsDone    bool              `json:"isDone"`
	VideoPath string            `json:"videoPath"`
	State     DownloadTaskState `json:"state"`
}
