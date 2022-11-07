package custom_error

import "fmt"

type NetworkError struct {
	Code int
	URL  string
}

func (e NetworkError) Error() string {
	return fmt.Sprintf("下载失败：Received HTTP %v for %v", e.Code, e.URL)
}
