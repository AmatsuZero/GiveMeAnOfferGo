package bilibili_api

type VideoInfo struct {
	BaseResponse
	Bvid        string        `json:"bvid"`
	Aid         string        `json:"aid"`
	SliceCount  int           `json:"videos"`
	Tid         string        `json:"tid"`
	Tnamae      string        `json:"tnamae"`
	Copyright   string        `json:"copyright"`
	Pic         string        `json:"pic"`
	Title       string        `json:"title"`
	Pubdate     int           `json:"pubdate"`
	CTime       int           `json:"c_time"`
	Desc        string        `json:"desc"`
	State       int           `json:"state"`
	Attribute   int           `json:"attribute"`
	Duration    int           `json:"duration"`
	MissionId   int           `json:"mission_id"`
	RedirectURL string        `json:"redirect_url"`
	Rights      VideoRights   `json:"rights"`
	Owner       VideoOwner    `json:"owner"`
	Stat        VideoState    `json:"stat"`
	Dynamic     string        `json:"dynamic"`
	Cid         int           `json:"cid"`
	Dimension   int           `json:"dimension"`
	NoCache     bool          `json:"no_cache"`
	Pages       []VideoPages  `json:"pages"`
	Subtitle    VideoSubtitle `json:"subtitle"`
	Staff       VideoStaff    `json:"staff, omitempty"`
}

type VideoOwner struct {
}

type VideoRights struct {
	Bp       int `json:"bp"`
	Elec     int `json:"elec"`
	download int `json:"download"`
	movide   int `json:"movide"`
}

func (rights VideoRights) CanDownload() bool {
	return rights.download == 1
}

func (rights VideoRights) IsMovie() bool {
	return rights.movide == 1
}

type VideoState struct {
}

type VideoPages struct {
}

type VideoDimension struct {
}

type VideoStaff struct {
}

type VideoSubtitle struct {
}
