package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// GET 访问网站，并提取网站标题
// 参数: http
// 返回值: url: [网站状态码, Server版本信息, 网站标题]
func GET(Httpurl string) list.List {
	//resultData :=make(map[string]list.List)

	l := list.List{}
	resp, err := http.Get(Httpurl)
	if err != nil {
		fmt.Println(err.Error())
		return list.List{}
	}
	//l.PushBack(Httpurl)

	defer resp.Body.Close()
	//fmt.Println(resp.Status)
	//fmt.Println(resp.Header)
	//resultData["status_code"] = resp.Status
	//resultData["server"] = resp.Header["Server"][0]

	l.PushBack(resp.Status)
	l.PushBack(resp.Header["Server"][0])

	c := make([]byte, 2048)
	var result string
	for {
		n, err := resp.Body.Read(c)
		if err != nil && io.EOF == err {
			break
		}
		result += string(c[:n])
		// 提前结束，其实我们只需要2048字节即可。
		break
	}

	//resultData["content"] = result
	l.PushBack(result[:])
	re, err := regexp.Compile("<title>(.*?)</title>")
	match := re.MatchString(result)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Println(match)
	title := ""
	if match {
		title = re.FindString(result)
	}
	//fmt.Println(title)
	//resultData["title"] = title
	l.PushBack(title)
	//fmt.Println(l)
	//resultData["Httpurl"] = l
	fmt.Println(l)

	return l
}

// Test 测试域名是否存活
// 参数: 域名字符串
// 返回值: 一个map对象
func Test(domain string) map[string]string {

	resultData := make(map[string]string)

	record, _ := net.LookupIP(domain)

	if len(record) == 0 {
		//fmt.Println("无效域名")
		return nil
	}
	fmt.Println(domain, record[0])

	//cur_httpurl := "http://" + domain
	//cur_httpsurl := "https://" + domain

	//httpcontent := GET(cur_httpurl)
	//for i:=httpcontent.Front(); i!= nil;i=i.Next(){
	//	fmt.Println(i.Value)
	//}
	//
	//fmt.Println(httpcontent)
	//
	//httpscontent := GET(cur_httpsurl)
	//fmt.Println(httpscontent)
	//for i:=httpcontent.Front(); i!= nil;i=i.Next(){
	//	fmt.Println(i.Value)
	//}
	resultData[domain] = record[0].String()
	return resultData
}

// SaveFile 文件保存，初步确定 csv格式
// 参数1: 保存的数据对象
// 参数2: 保存的文件名称
// 返回值: 是否成功保存
func SaveFile(list2 *list.List, savename string) bool {
	filecontent := ""

	//fmt.Println(list2)
	for i := list2.Front(); i != nil; i = i.Next() {
		fmt.Println(i.Value)
		xmap := i.Value.(map[string]string)
		for i := range xmap {
			fmt.Println("domain:", i, "IP:", xmap[i])
			filecontent += i + "," + xmap[i] + "\r\n"
		}
	}

	file, err := os.OpenFile(savename, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		fmt.Printf("err=%v", err)
		return false

	}

	defer file.Close()

	//使用缓存方式写入
	writer := bufio.NewWriter(file)
	count, w_err := writer.WriteString(filecontent)

	//需要使用Flush()将写入到writer缓存的数据真正写入到cw.txt文件中
	writer.Flush()
	if w_err != nil {
		fmt.Println("写入出错")
	} else {
		fmt.Printf("写入成功,共写入字节：%v", count)
	}
	return true
}

// ReadUrlFile 文件读取，按行读取
// 参数: url文件
// 返回值: url列表
// 备注: 要求处理url信息, 域名格式错误/不一致时(fofa的结果)
func ReadUrlFile(urlfile string) *list.List {
	l := list.New()

	file, err := os.Open(urlfile)
	if err != nil {
		fmt.Println("打开文件失败,", err)
	}
	reader := bufio.NewReader(file)
	//循环读取文件的内容
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF { //io.EOF表示文件的末尾
			break
		}
		domain := strings.Replace(str, "\n", "", 1)
		domain = strings.Replace(domain, "\r", "", 1)
		l.PushBack(domain)
		//fmt.Println(str)
	}
	return l
}

func main() {
	httpinfo := list.New()

	urllist := ReadUrlFile("target.txt")
	for i := urllist.Front(); i != nil; i = i.Next() {
		fmt.Println(i.Value.(string))
		domain := i.Value.(string)
		res := Test(domain)
		if res != nil {
			httpinfo.PushBack(res)
		}
	}
	//
	//domain := "baidu.cn"
	//res := Test(domain)
	//

	//
	//
	SaveFile(httpinfo, "muralist.csv")
}
