package bilibili_api

import "testing"

func TestGetBaseURL(t *testing.T) {
	t.Log(getVideoDetailURL("85440373", ""))
	t.Log(kBaseURL.String())
}
