package eventbus

import "github.com/wailsapp/wails/v2/pkg/runtime"

const (
	TaskAddEvent      = "task-add-reply"
	SelectVariant     = "select-variant"
	OnVariantSelected = "variant-selected"
	TaskNotifyCreate  = "task-notify-create"
	TaskStop          = "task-stop"
	TaskStatusUpdate  = "task-notify-update"
	TaskFinish        = "task-notify-end"
)

type RuntimeHandler interface {
	EventsEmit(eventName string, optionalData ...interface{})
	EventsOnce(eventName string, callback func(optionalData ...interface{}))
	MessageDialog(dialogOptions runtime.MessageDialogOptions) (string, error)
	EventsOn(eventName string, callback func(optionalData ...interface{}))
}
