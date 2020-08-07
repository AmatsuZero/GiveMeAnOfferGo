package bilibili_api

import (
	"fmt"
	"math"
	"strings"
)

func (size DanmukuFontSize) sizeRatio() float64 {
	return float64(size) / float64(DanmukuFontNormal)
}

// 字体是否被缩放过
func (size DanmukuFontSize) isScaled() bool {
	return size.sizeRatio() == 1
}

// 以 D 开头都是游客评论
func (danmuku Danmuku) isGuest() bool {
	return strings.HasPrefix(danmuku.UID, "D")
}

// 碰撞检测, 返回行号和时间偏移
func (danmuku Danmuku) detect(leaves []float64) (lineIndex int, offset float64) {
	beyonds := make([]float64, 0)
	for i, leave := range leaves {
		beyond := danmuku.Start - leave
		if beyond >= 0 { // 某一行有足够空间，直接返回行号和 0 偏移
			return i, 0
		}
		beyonds = append(beyonds, beyond)
	}
	// 所有行都没有空间了，那么找出哪一行能在最短时间内让出空间
	soon := float64(math.MinInt64)
	for i, beyond := range beyonds {
		if beyond > soon {
			soon = beyond
			lineIndex = i
		}
	}
	offset = -soon
	return
}

// 字符长度，1 个汉字当 2 个英文
func (danmuku Danmuku) displayLength() (l float64) {
	for _, r := range danmuku.Content {
		if r < 127 {
			l += 1
		} else {
			l += 2
		}
	}
	return
}

// 按用户自定义的字体大小来缩放的字体大小
func (danmuku Danmuku) scaledFontSize(c ASSConfig) int {
	return int(math.Ceil(float64(c.FontSize) * danmuku.FontSize.sizeRatio()))
}

/// 最长的行字符数
func (danmuku Danmuku) maxLength() (max int) {
	displayLen := func(str string) (l int) {
		for _, r := range str {
			if r < 127 {
				l += 1
			} else {
				l += 2
			}
		}
		return
	}
	for _, str := range strings.Split(danmuku.Content, "\n") {
		l := displayLen(str)
		if max < l {
			max = l
		}
	}
	return
}

/// 最长的行字符数
func (danmuku Danmuku) width(c ASSConfig) int {
	count := float64(danmuku.maxLength()) / 2
	return int(math.Ceil(float64(c.FontSize) * count))
}

// 整条字幕高度
func (danmuku Danmuku) height(c ASSConfig) int {
	count := len(strings.Split(danmuku.Content, "\n"))
	return count * c.FontSize
}

// 字幕出现和消失的水平坐标位置
func (danmuku Danmuku) hPos(c ASSConfig) (x1, x2 int) {
	switch danmuku.DanmukuType {
	case DanmukuTypeNormal: // 滚动字幕的水平位置参考点是整条字幕文本的中点
		return c.ScreenWidth + danmuku.width(c)/2, -danmuku.width(c) / 2
	default: // 默认在屏幕中间
		return c.ScreenWidth / 2, c.ScreenWidth / 2
	}
}

// 字幕每个字的移动的速度
func (danmuku Danmuku) speed(c ASSConfig) int {
	if danmuku.DanmukuType != DanmukuTypeNormal {
		return 0
	}
	// 基准时间，就是每个字的移动时间
	// 12 秒加上用户自定义的微调
	base := 12 + c.TuneDuration
	if base <= 0 {
		base = 0
	}
	return int(math.Ceil(float64(c.ScreenWidth) / base))
}

// 字幕出现和消失的垂直坐标位置
func (danmuku Danmuku) vPos(c ASSConfig, lineIndex int) (y1, y2 int) {
	switch danmuku.DanmukuType {
	case DanmukuTypeTop:
		return lineIndex * c.FontSize, lineIndex * c.FontSize
	case DanmukuTypeBottom:
		// 要让字幕不超出底部，减去高度
		y := c.ScreenHeight - (c.FontSize * lineIndex) - danmuku.height(c)
		// 再减去自定义的底部边距
		y -= c.BottomMargin
		return y, y
	case DanmukuTypeNormal:
		size := c.FontSize
		// 垂直位置，按基准字体大小算每一行的高度
		y := (lineIndex + 1) * size
		// 个别弹幕可能字体比基准要大，所以最上的一行还要避免挤出顶部屏幕
		// 坐标不能小于字体大小
		if y < danmuku.scaledFontSize(c) {
			y = danmuku.scaledFontSize(c)
		}
		return y, y
	default:
		return c.ScreenHeight / 2, c.ScreenHeight / 2
	}
}

// 字幕坐标点的移动距离
func (danmuku Danmuku) distance(c ASSConfig) int {
	x1, x2 := danmuku.hPos(c)
	return x1 - x2
}

// 整条字幕的显示时间
func (danmuku Danmuku) duration(c ASSConfig) (d float64) {
	switch danmuku.DanmukuType {
	case DanmukuTypeNormal:
		if c.SameSpeed { // 每个弹幕的滚动速度都一样，辨认度好，适合观看剧集类视频。
			d = float64(danmuku.distance(c)) / float64(danmuku.speed(c))
		} else { // 每个弹幕的滚动速度都不一样，动态调整，辨认度低，适合观看 MTV 类视频
			base := 6 + c.TuneDuration
			if base <= 0 {
				base = 0
			}
			count := danmuku.maxLength() / 2
			switch {
			case count < 6:
				d = base + float64(count)
			case count < 12:
				d = base + float64(count)/2
			case count < 24:
				d = base + float64(count)/3
			default:
				d = base + 10
			}
		}
	default:
		base := 3 + c.TuneDuration
		if base <= 0 {
			base = 0
		}
		count := danmuku.maxLength() / 2
		switch {
		case count < 6:
			d = base + 1
		case count < 12:
			d = base + 2
		default:
			d = base + 3
		}
	}
	return
}

func (danmuku Danmuku) end(c ASSConfig) float64 {
	switch danmuku.DanmukuType {
	case DanmukuTypeNormal:
		//对于滚动样式弹幕来说，就是最后一个字符离开最右边缘的时间。
		//坐标是字幕中点，在屏幕外和内各有半个字幕宽度。
		//也就是跑过一个字幕宽度的路程
		d := float64(danmuku.width(c)) / float64(danmuku.speed(c))
		return danmuku.Start + d
	default:
		return danmuku.Start + danmuku.duration(c)
	}
}

type ASSConfig struct {
	FontSize                  int
	ScreenWidth, ScreenHeight int
	BottomMargin              int
	TuneDuration              float64
	SameSpeed                 bool // 是否同步方式计算每条弹幕的显示时长
	LineCount                 int
	DropOffset, CustomOffset  float64
	TemplateHeader            string
	FontName                  string
}

// ass文件的头部
const templateHeader = `
[Script Info]
ScriptType: v4.00+
Collisions: Normal
PlayResX: {width}
PlayResY: {height}

[V4+ Styles]
Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding
Style: Default,{fontname},54,&H00FFFFFF,&H00FFFFFF,&H00000000,&H00000000,0,0,0,0,100,100,0.00,0.00,1,2.00,0.00,2,30,30,120,0
Style: Alternate,{fontname},36,&H00FFFFFF,&H00FFFFFF,&H00000000,&H00000000,0,0,0,0,100,100,0.00,0.00,1,2.00,0.00,2,30,30,84,0
Style: Danmaku,{fontname},{fontsize},&H00FFFFFF,&H00FFFFFF,&H00000000,&H00000000,0,0,0,0,100,100,0.00,0.00,1,1.00,0.00,2,30,30,30,0

[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
`

func (c ASSConfig) String() string {
	header := c.TemplateHeader
	if len(header) == 0 {
		header = templateHeader
	}
	header = strings.ReplaceAll(header, "{width}", fmt.Sprintf("%d", c.ScreenWidth))
	header = strings.ReplaceAll(header, "{height}", fmt.Sprintf("%d", c.ScreenHeight))
	header = strings.ReplaceAll(header, "{fontname}", c.FontName)
	header = strings.ReplaceAll(header, "{fontsize}", fmt.Sprintf("%d", c.FontSize))
	return header
}

type assSubtitle struct {
	*Danmuku
	Offset    float64
	LineIndex int
	Config    ASSConfig
}

func (subtitle assSubtitle) Start() float64 {
	return subtitle.Danmuku.Start + subtitle.Offset
}

func (subtitle assSubtitle) End() float64 {
	return subtitle.end(subtitle.Config) + subtitle.duration(subtitle.Config)
}

func (subtitle assSubtitle) Position() struct{ x1, x2, y1, y2 int } {
	x1, x2 := subtitle.hPos(subtitle.Config)
	y1, y2 := subtitle.vPos(subtitle.Config, subtitle.LineIndex)
	return struct{ x1, x2, y1, y2 int }{x1: x1, x2: x2, y1: y1, y2: y2}
}

// 秒数转 HH:MM:SS 格式
func (subtitle assSubtitle) s2hms(secs float64) string {
	if secs < 0 {
		return "0:00:00.00"
	}
	i, d := math.Modf(secs / 1)
	m, s := math.Modf(i / 60)
	h, m := math.Modf(m / 60)
	return fmt.Sprintf("%d:%02d:%02d.%02d", int(h), int(m), int(s), int(d*100))
}

func (subtitle assSubtitle) StartMarkup() string {
	return subtitle.s2hms(subtitle.Start())
}

func (subtitle assSubtitle) EndMarkup() string {
	return subtitle.s2hms(subtitle.End())
}

func (subtitle assSubtitle) ColorMarkup() string {
	// 白色不用特别标记
	if strings.ToUpper(subtitle.Color) == "FFFFFF" {
		return ""
	}
	return "\\\\c&H" + subtitle.Color
}

func (subtitle assSubtitle) FontSizeMarkup() string {
	if subtitle.FontSize.isScaled() {
		return fmt.Sprintf("\\\\fs%d", subtitle.FontSize)
	}
	return ""
}

func (subtitle assSubtitle) StyleMarkup() string {
	pos := subtitle.Position()
	switch subtitle.DanmukuType {
	case DanmukuTypeNormal:
		return fmt.Sprintf("\\\\move(%d, %d, %d, %d)", pos.x1, pos.y2, pos.x2, pos.y2)
	default:
		return fmt.Sprintf("\\\\a6\\\\pos(%d, %d)", pos.x1, pos.y1)
	}
}

func (subtitle assSubtitle) LayerMarkup() string {
	switch subtitle.DanmukuType {
	case DanmukuTypeNormal:
		return "-2"
	default:
		return "-3"
	}
}

func (subtitle assSubtitle) ContentMarkup() string {
	return fmt.Sprintf("{%v}%v", strings.Join([]string{
		subtitle.StyleMarkup(),
		subtitle.ColorMarkup(),
		subtitle.BorderMarkup(),
		subtitle.FontSizeMarkup(),
	}, ""), subtitle.Content)
}

func (subtitle assSubtitle) String() string {
	return fmt.Sprintf("Dialogue: %v,%v,%v,Danmaku,,0000,0000,0000,,%v",
		subtitle.LayerMarkup(),
		subtitle.StartMarkup(),
		subtitle.EndMarkup(),
		subtitle.ContentMarkup())
}

func (subtitle assSubtitle) BorderMarkup() string {
	return ""
}

func convertToAss(c ASSConfig, danmukuArr []*Danmuku) (dropoutCnt int, file string) {
	collisions := map[DanmukuType][]float64{
		DanmukuTypeNormal: make([]float64, c.LineCount),
		DanmukuTypeTop:    make([]float64, c.LineCount),
		DanmukuTypeBottom: make([]float64, c.LineCount),
	}
	subtitles := make([]assSubtitle, 0)
	for _, danmuku := range danmukuArr {
		collision, ok := collisions[danmuku.DanmukuType]
		if !ok {
			continue
		}
		lineIndex, waitingOffset := danmuku.detect(collision)
		if waitingOffset > c.DropOffset { // 超过容忍的偏移量，丢弃掉此条弹幕
			continue
		}
		collision[lineIndex] = danmuku.end(c) + waitingOffset
		// 再加上自定义偏移
		offset := waitingOffset + c.CustomOffset
		subtitles = append(subtitles, assSubtitle{
			Danmuku:   danmuku,
			Offset:    offset,
			LineIndex: lineIndex,
			Config:    c,
		})
	}
	header := fmt.Sprintf("%v", c)
	events := make([]string, 0, len(subtitles))
	for _, s := range subtitles {
		events = append(events, fmt.Sprintf("%v", s))
	}
	return len(danmukuArr) - len(subtitles), header + strings.Join(events, "\n")
}
