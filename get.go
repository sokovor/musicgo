package main

import (
	"encoding/xml"
	"fmt"
	iconv "github.com/djimenez/iconv-go"
	"regexp"
	"strconv"
	"strings"
)

func GetSosoMP3URL(data string) (urls []string, songName []string, singer []string, suc bool) {
	suc = true
	re, _ := regexp.Compile("@@(.*)@@.*@@(.*)@@.*@@.*@@.*@@.*@@FI(http.*);;")
	list := re.FindAllStringSubmatch(data, 3)
	if len(list) < 1 {
		suc = false
		return
	}
	fmt.Println("执行。。。。。。。。")
	//fmt.Printf("%v\n", list)
	urls = make([]string, len(list))
	songName = make([]string, len(list))
	singer = make([]string, len(list))
	converter, _ := iconv.NewConverter("GBK", "utf-8")
	for i, value := range list {
		songName[i], _ = converter.ConvertString(value[1])
		singer[i], _ = converter.ConvertString(value[2])
		//fmt.Printf("%v\n%v\n%v\n%v\n", value[0], value[1], value[2], value[3])
		subre, _ := regexp.Compile("(\\d)+")
		result := subre.FindAllString(value[3], 2)

		if result == nil || len(result) != 2 {
			urls[i] = value[3]
		} else {
			steam, _ := strconv.ParseInt(result[0], 10, 32)
			songid, _ := strconv.ParseInt(result[1], 10, 32)
			songid = songid - 12000000 + 30000000
			urls[i] = "http://stream" + strconv.FormatInt((steam+10), 10) + ".qqmusic.qq.com/" + strconv.FormatInt(songid, 10) + ".mp3"
		}
	}
	return
}

func GetMP3URL(data []byte) string {

	type Result struct {
		XMLName xml.Name `xml:"result"`
		Count   int      `xml:"count"`
		Encode  []string `xml:"url>encode"`
		Decode  []string `xml:"url>decode"`
		Typed   []string `xml:"url>type"`
	}
	var v Result

	out := make([]byte, len(data))

	iconv.Convert(data, out, "gb2312", "utf-8")
	fmt.Printf("%v\n", string(out))
	fString := strings.Replace(string(out), "gb2312", "utf-8", 1)
	fmt.Printf("%v\n", fString)
	err := xml.Unmarshal([]byte(fString), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return "error"
	}
	if v.Count < 0 {
		return "没有找到"
	}

	for i, value := range v.Typed {
		if value == "8" {
			return v.Encode[i][0:strings.LastIndex(v.Encode[i], "/")+1] + v.Decode[i]
		}
	}
	return "没有mp3格式"

}
