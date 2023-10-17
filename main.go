package main

import (
	"compress/flate"
	"encoding/xml"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Chat struct {
	XMLName    xml.Name  `xml:"i"`
	ChatServer string    `xml:"chatserver"`
	ChatID     int       `xml:"chatid"`
	Mission    int       `xml:"mission"`
	MaxLimit   int       `xml:"maxlimit"`
	State      int       `xml:"state"`
	RealName   int       `xml:"real_name"`
	Source     string    `xml:"source"`
	Messages   []Message `xml:"d"`
}
type Message struct {
	Text string `xml:",chardata"`
	P    string `xml:"p,attr"`
}

type Result struct {
	AppearTime string
	SendTime   string
	FontSize   string
	Color      string
	Message    string
}

func main() {
	//发起请求
	url := "https://api.bilibili.com/x/v1/dm/list.so?oid=437586584"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("get failed,err:", err)
	}
	//处理编码
	reader := flate.NewReader(resp.Body)
	//解析xml
	xmlData, err := io.ReadAll(reader)
	var chat Chat
	err = xml.Unmarshal(xmlData, &chat)
	if err != nil {
		fmt.Println("unmarshal failed,err:", err)
		return
	}
	//// 打印解析得到的数据
	//fmt.Printf("Chat Server: %s\n", chat.ChatServer)
	//fmt.Printf("Chat ID: %d\n", chat.ChatID)
	//fmt.Printf("Mission: %d\n", chat.Mission)
	//fmt.Printf("Max Limit: %d\n", chat.MaxLimit)
	//fmt.Printf("State: %d\n", chat.State)
	//fmt.Printf("Real Name: %d\n", chat.RealName)
	//fmt.Printf("Source: %s\n", chat.Source)

	// 打印每条消息
	//for _, message := range chat.Messages {
	//	fmt.Printf("Message: %s, P: %s\n", message.Text, message.P)
	//}
	//
	var result []Result
	for _, message := range chat.Messages {
		res, err := handlerP(message.P)
		if err != nil {
			fmt.Println("handlerP failed,err:", err)
			return
		}
		res.Message = message.Text
		result = append(result, res)
	}

	//生成excel
	buildExcel(result)

}

// 生成excel
func buildExcel(data []Result) {
	// 创建一个工作表
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	f.SetColWidth("Sheet1", "A", "D", 20)
	f.SetColWidth("Sheet1", "E", "E", 200)
	// 设置单元格的值
	f.SetCellValue("Sheet1", "A1", "弹幕出现时间")
	f.SetCellValue("Sheet1", "B1", "发送时间")
	f.SetCellValue("Sheet1", "C1", "字体大小")
	f.SetCellValue("Sheet1", "D1", "颜色")
	f.SetCellValue("Sheet1", "E1", "弹幕内容")
	for i, v := range data {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), v.AppearTime)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), v.SendTime)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), v.FontSize)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), v.Color)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), v.Message)
	}
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err := f.SaveAs(time.Now().Format("2006-01-02") + ".xlsx"); err != nil {
		fmt.Println(err)
	}
}

// 处理颜色
func changeColor(s int) string {
	return fmt.Sprintf("#%06X", s)

}

// 处理数据
func handlerP(p string) (res Result, err error) {
	// 弹幕出现时间,模式,字体大小,颜色,发送时间戳,弹幕池,用户Hash,数据库ID,page
	//    <d p="91.35700,1,25,16777215,1679280461,0,c61a3d4c,1276656055551609600,10">愿人类荣光永存</d>
	parts := strings.Split(p, ",")
	appearTimeSecond, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return Result{}, err
	}
	appearTime := secondsToHourMinuteSecond(appearTimeSecond)
	fontSize := parts[2]

	colorInt, err := strconv.Atoi(parts[3])
	color := changeColor(colorInt)

	sendTimeUnix, err := strconv.Atoi(parts[4])
	sendTime := timeFormat(sendTimeUnix)

	return Result{AppearTime: appearTime, SendTime: sendTime, FontSize: fontSize, Color: color}, nil
}

// 处理发送时间(年月日 分时秒)
func timeFormat(timeUnix int) string {
	// 例子中的时间戳
	timestamp := int64(timeUnix)
	// 将时间戳转换为time.Time类型
	timeValue := time.Unix(timestamp, 0)
	// 将时间格式化为时分秒
	timeFormatted := timeValue.Format("06/01/02 15:04:05")
	return timeFormatted
}

// 处理弹幕出现在视频中的时间(视频第几分钟出现）
func secondsToHourMinuteSecond(seconds float64) string {
	hours := int(seconds / 3600)
	remainingSeconds := seconds - float64(hours*3600)
	minutes := int(remainingSeconds / 60)
	seconds = remainingSeconds - float64(minutes*60)

	return fmt.Sprintf("%02d:%02d:%05.2f", hours, minutes, seconds)
}
