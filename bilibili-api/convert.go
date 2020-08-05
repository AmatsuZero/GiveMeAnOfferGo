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

// ass文件的头部
const header = `
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
func (danmuku Danmuku) scaledFontSize(c AssConfig) int {
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
func (danmuku Danmuku) width(c AssConfig) int {
	count := float64(danmuku.maxLength()) / 2
	return int(math.Ceil(float64(c.FontSize) * count))
}

// 整条字幕高度
func (danmuku Danmuku) height(c AssConfig) int {
	count := len(strings.Split(danmuku.Content, "\n"))
	return count * c.FontSize
}

// 字幕出现和消失的水平坐标位置
func (danmuku Danmuku) hPos(c AssConfig) (x1, x2 int) {
	switch danmuku.DanmukuType {
	case DanmukuTypeNormal: // 滚动字幕的水平位置参考点是整条字幕文本的中点
		return c.ScreenWidth + danmuku.width(c)/2, -danmuku.width(c) / 2
	default: // 默认在屏幕中间
		return c.ScreenWidth / 2, c.ScreenWidth / 2
	}
}

// 字幕出现和消失的垂直坐标位置
func (danmuku Danmuku) vPos(c AssConfig, lineIndex int) (y1, y2 int) {
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

// 整条字幕的显示时间
func (danmuku Danmuku) duration(c AssConfig) (d float64) {
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
	return
}

func (danmuku Danmuku) end(c AssConfig) float64 {
	return danmuku.Start + danmuku.duration(c)
}

// 秒数转 HH:MM:SS 格式
func (danmuku Danmuku) hhmmssTiming() string {
	if danmuku.Start < 0 {
		return "0:00:00.00"
	}
	i, d := math.Modf(danmuku.Start / 1)
	m, s := math.Modf(i / 60)
	h, m := math.Modf(m / 60)
	return fmt.Sprintf("%d:%02d:%02d.%02d", int(h), int(m), int(s), int(d*100))
}

type AssConfig struct {
	FontSize                  int
	ScreenWidth, ScreenHeight int
	BottomMargin              int
	TuneDuration              float64
}

func (danmuku Danmuku) convertToAss(config AssConfig) {

}
