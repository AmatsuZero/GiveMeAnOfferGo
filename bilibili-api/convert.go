package bilibili_api

import (
	"fmt"
	"math"
	"strings"
)

func (size DanmukuFontSize) sizeRatio() int {
	return int(size / DanmukuFontNormal)
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

// 生成时间点的ass格式表示：`0:00:00.00`
func (danmuku Danmuku) timeSpot() string {
	h, f := math.Modf(danmuku.Timing / 3600)
	m, f := math.Modf(f * 60)
	return fmt.Sprintf("%d:%02d:%05.2f", int(h), int(m), f*60)
}

type AssConfig struct {
	FontSize                  int
	ScreenWidth, ScreenHeight int
}

func (danmuku Danmuku) convertToAss(config AssConfig) {

}
