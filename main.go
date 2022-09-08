package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
)

// 访问网站，并提取网站标题
// 参数: httpurl
// 返回值: url: [网站标题, Server版本信息]
func GET(Httpurl string) {
	resp, err := http.Get(Httpurl)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	fmt.Println(resp.Header["Server"])

	c := make([]byte, 2048)
	var result string
	for {
		n, err := resp.Body.Read(c)
		if err != nil && io.EOF == err {
			break
		}
		result += string(c[:n])
	}

	re, err := regexp.Compile("<title>(.*?)</title>")
	match := re.MatchString(result)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(match)
	title := ""
	if match {
		title = re.FindString(result)
	}
	fmt.Println(title)

}

// 测试域名是否存活
// 参数: 域名字符串
// 返回值: 一个map对象
func Test(domain string) {
	record, _ := net.LookupIP(domain)

	if len(record) == 0 {
		fmt.Println("无效域名")
		return
	}
	fmt.Println(record[0])

	cur_httpurl := "http://" + domain
	cur_httpsurl := "https://" + domain

	GET(cur_httpsurl)
	GET(cur_httpurl)

}

// 文件保存，初步确定 csv格式
// 参数1: 保存的数据对象
// 参数2: 保存的文件名称
// 返回值: 是否成功保存
func SaveFile() {

}

func main() {

	domain := "dnslog.cn"
	Test(domain)
}
