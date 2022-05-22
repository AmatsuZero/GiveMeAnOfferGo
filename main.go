package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"math/big"
	"strconv"
	"strings"
)

const sourceName = "files/source.xlsx"
const targetName = "files/book.xlsx"
const sourceSheetName = "公寓484套"
const targetSheetName = "公寓楼"

type rowInfo []string

type roomInfo struct {
	roomNo      string // 房间号
	area        string // 面积
	presentUser string // 现状态
}

type ownerInfo struct {
	rooms []roomInfo
	name  string // 业主
}

func main() {
	source, err := excelize.OpenFile(sourceName)
	if err != nil {
		fmt.Println(err)
		return
	}

	allInfos := extractAllInfo(source)     // 提取所有信息成字典
	if err := source.Close(); err != nil { // 关闭数据源表格
		fmt.Println(err)
	}

	dest, err := excelize.OpenFile(targetName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := dest.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	updateSheets(dest, allInfos)
	err = dest.Save() // 保存
	if err != nil {
		fmt.Println(err)
	}
}

func updateSheets(sheet *excelize.File, src map[string]*ownerInfo) {
	rows, err := sheet.GetRows(targetSheetName)
	if err != nil {
		fmt.Println(err)
		return
	}
	rowsCnt := len(rows)
	for i := 2; i < rowsCnt; i++ {
		row := rowInfo(rows[i]) // 这一行的信息
		name := row.getValueFromIndex(1)
		realRowNum := i + 1 // 行号跟数组下标不一样
		if len(name) == 0 {
			fmt.Printf("❌ 第 %v 行人名为空\n", realRowNum)
			continue
		}
		names := strings.Split(name, "/") // 可能名字有多个
		for _, n := range names {
			info, ok := src[n]
			if !ok {
				fmt.Printf("❌ 第 %v 行没有找到这个人：%v\n", realRowNum, n)
				// 设置样式
				fillStyle(sheet, realRowNum)
				continue
			}
			// 填写数量
			axis := "E" + strconv.Itoa(realRowNum)
			sheet.SetCellValue(targetSheetName, axis, len(info.rooms))
			// 填写面积
			axis = "F" + strconv.Itoa(realRowNum)
			sheet.SetCellValue(targetSheetName, axis, info.SumOfArea())
			// 填写备注
			axis = "I" + strconv.Itoa(realRowNum)
			sheet.SetCellValue(targetSheetName, axis, info.RoomNos())
		}
	}
}

func fillStyle(sheet *excelize.File, realRowNum int) {
	// 统一设置背景色为黄色 / 居中样式 / 黑色边框
	style, _ := sheet.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFFF00"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})
	start := "A" + strconv.Itoa(realRowNum)
	end := "I" + strconv.Itoa(realRowNum)
	sheet.SetCellStyle(targetSheetName, start, end, style)
}

func extractAllInfo(sheet *excelize.File) (source map[string]*ownerInfo) {
	rows, err := sheet.GetRows(sourceSheetName)
	if err != nil {
		fmt.Println(err)
		return
	}

	source = make(map[string]*ownerInfo)
	rowsCnt := len(rows)

	for i := 3; i < rowsCnt; i++ { // 跳过表头
		row := rowInfo(rows[i]) // 这一行的信息
		if row == nil {         // 跳过空行
			continue
		}
		name := row.getValueFromIndex(2)
		if len(name) == 0 {
			fmt.Printf("❌ 第 %v 行人名为空\n", i+1)
			continue
		}
		// 有的房子包含多个业主，这样各计一次
		names := strings.Split(name, "/")
		for _, n := range names {
			info, ok := source[n]
			if !ok {
				info = &ownerInfo{
					name:  n,
					rooms: make([]roomInfo, 0),
				}
				source[n] = info
			}
			if ri := newRoomInfoFromRow(row); ri.IsValid() {
				info.rooms = append(info.rooms, ri)
			}
		}
	}
	return
}

func newRoomInfoFromRow(r rowInfo) roomInfo {
	return roomInfo{
		roomNo:      r.getValueFromIndex(1),
		presentUser: r.getValueFromIndex(3),
		area:        r.getValueFromIndex(4),
	}
}

func (r roomInfo) IsValid() bool { // 至少有房间号和面积才算有效
	return len(r.roomNo) > 0 && len(r.area) > 0
}

func (r rowInfo) getValueFromIndex(idx int) string {
	if idx >= len(r) {
		return ""
	}
	return r[idx]
}

func (o ownerInfo) SumOfArea() string { // 这里用高精度求值
	const precision = 200
	if len(o.rooms) == 0 {
		return "0"
	} else if len(o.rooms) == 1 {
		return o.rooms[0].area
	}
	result, _ := new(big.Float).SetPrec(precision).SetString(o.rooms[0].area)
	for i := 1; i < len(o.rooms); i++ {
		r := o.rooms[i]
		num, _ := new(big.Float).SetPrec(precision).SetString(r.area)
		result = result.Add(result, num)
	}
	return result.String()
}

func (o ownerInfo) RoomNos() string {
	var arr []string
	for _, r := range o.rooms {
		arr = append(arr, r.roomNo)
	}
	return strings.Join(arr, "、")
}
