package main

import (
	"fmt"
	iconv "github.com/djimenez/iconv-go"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Hello struct{}

//func (h Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	//s := "飞得更高"
//	//name := "汪峰"
//	//url := "http://box.zhangmen.baidu.com/x?op=12&count=1&title=" + s + "$$" + name + "$$$$"
//	//url := "http://music.soso.com/music.cgi?ty=getsongurls&w=" + s + "&pl=" + name

//	url := "http://cgi.music.soso.com/fcgi-bin/m.q?w=andy&source=1&t=1"

//	response, _ := http.Get(url)
//	defer response.Body.Close()
//	body, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(string(body))
//	//fmt.Fprintln(w, string(body))
//	//str := GetMP3URL(body)
//	//fmt.Fprintln(w, "<!DOCTYPE HTML>")
//	//fmt.Fprintln(w, "<html>")
//	//fmt.Fprintln(w, "<body>")
//	//fmt.Fprintln(w, "<audio src=\""+str+"\" controls=\"controls\">")
//	//fmt.Fprintln(w, "Your browser does not support the audio tag.")
//	//fmt.Fprintln(w, "</audio>")
//	////fmt.Fprint(w, "<a href='"+str+"'>"+str+"</a>")
//	//fmt.Fprintln(w, "</body>")
//	//fmt.Fprintln(w, "</html")
//}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("starting....")
	r.ParseForm()
	queryArr := r.Form["w"]
	fmt.Println(queryArr)
	if len(queryArr) < 1 {
		return
	}
	querystr := queryArr[0]
	in := []byte(querystr)
	out := make([]byte, len(in))

	iconv.Convert(in, out, "utf-8", "gbk")

	fmt.Printf("%v\n", string(out))

	v := url.Values{}
	v.Add("p", "1")
	v.Add("source", "1")
	v.Add("t", "1")
	v.Add("w", string(out))
	//s := "andy"
	url := "http://cgi.music.soso.com/fcgi-bin/m.q?" + v.Encode()

	//fmt.Println(url)

	//response, _ := http.Get(url)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	cgiUrl := "http://soso.music.qq.com/fcgi-bin/fcg_song.fcg"
	cgiRes, _ := http.Get(cgiUrl)
	defer cgiRes.Body.Close()
	fmt.Printf(":::::%v\n", cgiRes.Header)
	urls, songName, singer, _ := GetSosoMP3URL(string(body))
	//w.Header().Add("Referer", "http://soso.music.qq.com/fcgi-bin/fcg_song.fcg")
	//w.Header().Add("Set-Cookie", "qqmusic_sosokey=4D96476733A6D833E90FEA9E590408D171B92452775E15FB")
	//w.Header().Add("Set-Cookie", "qqmusic_fromtag=10")
	//w.Header().Add("Set-Cookie", "domain=qq.com")
	fmt.Fprintln(w, "<!DOCTYPE HTML>")
	fmt.Fprintln(w, "<html>")

	fmt.Fprintln(w, "<body>")

	for i, value := range urls {
		fmt.Fprintln(w, "<audio src=\""+value+"\" controls=\"controls\">")
		fmt.Fprintln(w, "Your browser does not support the audio tag.")
		fmt.Fprintln(w, "</audio>")
		fmt.Fprintln(w, songName[i])
		fmt.Fprintln(w, "&nbsp&nbsp&nbsp")
		fmt.Fprintln(w, singer[i])
		fmt.Fprintln(w, "<br>")
	}
	fmt.Fprintln(w, "</body>")
	fmt.Fprintln(w, "</html>")

}

func Start() {
	fmt.Println("ready for connect..")
	//var h Hello
	http.HandleFunc("/", serveHTTP)
	http.ListenAndServe("localhost:4000", nil)
}
