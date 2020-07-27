package bilibili_api

import (
	"context"
	"encoding/json"
	"github.com/reactivex/rxgo/v2"
	"net/http"
	"net/url"
)

/*
获取视频详情
*/
type VideoInfoRequest struct {
	baseRequest
	Aid  string // 视频avID
	Bvid string // 视频avID
}

func (request VideoInfoRequest) IsParamsValid() bool {
	return len(request.Aid) > 0 || len(request.Bvid) > 0
}

func (request VideoInfoRequest) Request() (*http.Request, error) {
	return request.requestByPath("x/web-interface/view")
}

func (request VideoInfoRequest) requestByPath(path string) (*http.Request, error) {
	if !request.IsParamsValid() || len(path) == 0 {
		return nil, kInvalidParamError
	}
	base := *kBaseURL
	base.Path = path
	base.RawQuery = request.queryItems(base.Query()).Encode()
	return http.NewRequest("GET", base.String(), nil)
}

func (request VideoInfoRequest) queryItems(query url.Values) url.Values {
	if len(request.Aid) > 0 {
		query.Add("aid", request.Aid)
	} else {
		query.Add("bvid", request.Bvid)
	}
	return query
}

/*
获取完整信息
*/
func (request VideoInfoRequest) Fetch(client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	req, err := request.Request()
	if err != nil {
		return rxgo.Thrown(err)
	}
	return request.baseRequest.fetch(client, req, opts...).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		var info VideoInfo
		data := i.([]byte)
		err := json.Unmarshal(data, &info)
		if err != nil {
			return nil, err
		}
		return info, nil
	})
}

type VideoDescRequest struct {
	VideoInfoRequest
}

func (request VideoDescRequest) Request() (*http.Request, error) {
	return request.requestByPath("x/web-interface/archive/desc")
}

func (request VideoDescRequest) Fetch(client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	req, err := request.Request()
	if err != nil {
		return rxgo.Thrown(err)
	}
	return request.fetch(client, req, opts...).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		var info VideoDescription
		data := i.([]byte)
		err := json.Unmarshal(data, &info)
		if err != nil {
			return nil, err
		}
		return info, nil
	})
}

type VideoDescription struct {
	BaseResponse
	Data string
}

type VideoPageListRequest struct {
	VideoInfoRequest
}

func (request VideoPageListRequest) Request() (*http.Request, error) {
	return request.requestByPath("x/player/pagelist")
}

func (request VideoPageListRequest) Fetch(client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	req, err := request.Request()
	if err != nil {
		return rxgo.Thrown(err)
	}
	return request.fetch(client, req, opts...).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		var info VideoPageList
		data := i.([]byte)
		err := json.Unmarshal(data, &info)
		if err != nil {
			return nil, err
		}
		return info, nil
	})
}

type VideoPageList struct {
	BaseResponse
	Data []VideoPages
}

type VideoInfo struct {
	BaseResponse
	Data struct {
		Bvid        string // 视频 bvID
		Aid         int    // 视频 avID
		SliceCount  int    // 视频分P总数
		Tid         int    // 分区ID
		Tname       string // 子分区名称
		Copyright   int    // 版权标志
		Pic         string // 视频封面图片 url
		Title       string // 视频标题
		Pubdate     int
		CTime       int
		Desc        string
		State       int
		Attribute   int
		Duration    int
		MissionId   int
		RedirectURL string
		Rights      VideoRights
		Owner       VideoOwner
		Stat        VideoState
		Dynamic     string // 视频同步发布的的动态的文字内容
		Cid         int
		Dimension   VideoDimension
		NoCache     bool
		Pages       []VideoPages
		Subtitle    VideoSubtitle
		Staff       []VideoStaff `json:"staff, omitempty"`
	}
}

func (info VideoInfo) IsReprint() bool {
	return info.Data.Copyright == 2
}

type VideoOwner struct {
	Mid  int    // Up主 uid
	Name string // Up主昵称
	Face string // Up主头像
}

type VideoRights struct {
	Bp            int
	Elec          int
	Download      int
	Movie         int
	Pay           int
	Hd5           int
	NoReprint     int
	Autoplay      int
	UgcPay        int
	IsCooperation int
	UgcPayPreview int
	NoBackground  int
}

func (rights VideoRights) IsAutoPlay() bool {
	return rights.Autoplay == 1
}

func (rights VideoRights) CanDownload() bool {
	return rights.Download == 1
}

func (rights VideoRights) IsMovie() bool {
	return rights.Movie == 1
}

func (rights VideoRights) IsVIPOnly() bool {
	return rights.Pay == 1
}

func (rights VideoRights) IsHD() bool {
	return rights.Hd5 == 1
}

func (rights VideoRights) CanReprint() bool {
	return rights.NoReprint == 0
}

func (rights VideoRights) IsMadeByCooperation() bool {
	return rights.IsCooperation == 1
}

type VideoState struct {
	Aid                int
	ViewCount          int `json:"view"`
	DanmukuCount       int `json:"danmuku"`
	ReplyCount         int `json:"reply"`
	FavoriteCount      int `json:"favorite"`
	CoinCount          int `json:"coin"`
	ShareCount         int `json:"share"`
	NowRank            int
	HistoryHighestRank int `json:"his_rank"`
	LikeCount          int `json:"like"`
	Dislike            int
	Evaluation         string `json:"evaluation, omitempty"` // 默认为空
}

func (state VideoState) IsBlocked() bool {
	return state.ViewCount == -1
}

type VideoPages struct {
	Cid       int
	Page      int
	From      string
	PartTitle string
	Duration  int
	Vid       string
	WebLink   string `json:"weblink"`
	Dimension VideoDimension
}

func (pages VideoPages) IsFromMangoTV() bool {
	return pages.From == "hunan"
}

type VideoDimension struct {
	Width  int // 当前分P 宽度
	Height int // 当前分P 高度
	Rotate int // 是否将宽高对换
}

func (dimension VideoDimension) IsPortrait() bool {
	return dimension.Rotate == 1
}

type VideoStaff struct {
	VideoOwner
	Title         string
	FollowerCount int          `json:"follower"`
	VIPInfo       VIPInfo      `json:"vip"`
	OfficialInfo  OfficialInfo `json:"official"`
}

type OfficialInfo struct {
	RoleLevel int    `json:"role"`
	Title     string `json:"title, omitempty"`
	Desc      string `json:"desc, omitempty"`
	Type      int
}

func (info OfficialInfo) IsCertified() bool {
	return info.Type == 0
}

func (info OfficialInfo) IsPersonal() bool {
	return info.Type == 1 || info.Type == 2
}

func (info OfficialInfo) IsOrganizational() bool {
	return info.Type >= 3 && info.Type <= 6
}

type MembershipType int

const (
	MemberTypeNone MembershipType = iota
	MemberTypeMonthly
	MemberTypeYearly
)

type VIPInfo struct {
	MemberShip MembershipType `json:"type"`
	Status     int
	ThemeType  int
}

func (info VIPInfo) IsAvailable() bool {
	return info.Status == 1
}

type VideoSubtitle struct {
	CanSubmit bool `json:"allow_submit"`
	List      []Subtitle
}

type Subtitle struct {
	Id          int
	Lan         string
	LanDoc      string
	IsLock      int
	AuthorMid   int
	SubtitleURL string `json:"subtitle_url"`
	Author      SubtitleAuthor
}

func (subtitle Subtitle) CanEdit() bool {
	return subtitle.IsLock == 0
}

type SubtitleAuthor struct {
	VideoOwner
	Sex           string
	Sign          string
	Rank          int
	Birthday      int
	IsFakeAccount int
	IsDeleted     int
}
