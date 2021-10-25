package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func InputUserPwd(user, passwd string, picByte *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("https://119.97.153.194:85/"),
		chromedp.WaitVisible(`#user`, chromedp.ByID),
		chromedp.Sleep(1 * time.Second),
		chromedp.SendKeys(`//*[@id="user"]`, user),
		chromedp.SendKeys(`//*[@id="password"]`, passwd),
		chromedp.Sleep(1 * time.Second),
		chromedp.Screenshot(`//*[@id="verify_code"]`, picByte),
	}
}

type IP struct {
	ip   string
	up   string
	down string
	all  string
}

func ParseSlice(s []*string) []IP {
	lip := make([]IP, 0, 1)
	for _, s := range s {
		var ip IP
		unnecessaryChar, _ := regexp.Compile("[ \\- ,]")
		fs := unnecessaryChar.ReplaceAllString(*s, "")
		ls := strings.Split(fs, "\n")
		ip.ip = ls[1]
		ip.up = ls[4]
		ip.down = ls[5]
		ip.all = ls[6]
		lip = append(lip, ip)
	}
	return lip
}

func GetTrString(n int, s *string) chromedp.Tasks {
	xpath := fmt.Sprintf("//*[@id=\"grid\"]/div/table/tbody/tr[%d]", n)
	return chromedp.Tasks{
		chromedp.TextContent(xpath, s),
	}
}

func Refund(user, pwd, id, softid string) string {
	client := &http.Client{}
	var req *http.Request
	var resp *http.Response
	var err error
	var body []byte
	urlString := "http://upload.chaojiying.net/Upload/ReportError.php"
	parameters := url.Values{}
	parameters.Add("user", user)
	parameters.Add("pass", pwd)
	parameters.Add("softid", softid)
	parameters.Add("id", id)
	req, err = http.NewRequest("POST", urlString, strings.NewReader(parameters.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0)")
	req.Header.Set("Connection", "Keep-Alive")
	if err != nil {
		log.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	r := gjson.Parse(string(body))
	if r.Get("err_str").String() != "OK" {
		fmt.Printf("refund Error : %s\n", r.Get("err_str").String())
		return r.Get("err_str").String()
	}
	return ""
}
func InputVerifyCode(verifyCode string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.SendKeys(`//*[@id="verify"]`, verifyCode),
		chromedp.WaitVisible(`//*[@id="button"]`),
		chromedp.Click(`//*[@id="button"]`),
		chromedp.Sleep(3 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetAllCookies().Do(ctx)
			if err != nil {
				return err
			}
			Cookie = cookies[0].Value
			return err
		}),
		chromedp.Sleep(1 * time.Second),
	}
}
func GetChromedp(ctx context.Context) (context.Context, context.CancelFunc) {
	options := chromedp.DefaultExecAllocatorOptions[:]
	options = append(options, chromedp.Flag("headless", false), chromedp.Flag("ignore-certificate-errors", "1"))
	c, cancel := chromedp.NewExecAllocator(ctx, options...)
	taskCtx, _ := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	return taskCtx, cancel
}
func getEncodedBase64(file []byte) string {
	encoded := base64.StdEncoding.EncodeToString(file)
	return encoded
}

//func GetHtmlFile() []byte {
//	urlString := "https://119.97.153.194:85/CFluxStatistic.php?sid=" + Cookie
//	fmt.Println(urlString)
//	tr := &http.Transport{
//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//	}
//	client := &http.Client{Transport: tr}
//	var req *http.Request
//	var resp *http.Response
//	var err error
//	var body []byte
//	parameters := url.Values{}
//	parameters.Add("stat_method", "ip_user")
//	parameters.Add("rank_base", "total_flow")
//	parameters.Add("custDate", "0")
//	parameters.Add("startDateTime", "2021-10-22 07:00")
//	parameters.Add("endDateTime", "2021-10-22 08:00")
//	parameters.Add("view_num", "10")
//	parameters.Add("schedule", "0")
//	parameters.Add("sourceType", "")
//	parameters.Add("sourceIp", "")
//	parameters.Add("sourceUser", "")
//	parameters.Add("sourceGroup", "")
//	parameters.Add("sub_group", "")
//	parameters.Add("appTypeText", "所有应用")
//	parameters.Add("appType", "0")
//	parameters.Add("appName", "0")
//	parameters.Add("graph_type", "0")
//	parameters.Add("formState", "unwrap")
//	parameters.Add("curPage", "1")
//	parameters.Add("type", "report")
//	parameters.Add("sid", Cookie)
//	parameters.Add("progressive", "true")
//	req, err = http.NewRequest("POST", urlString, strings.NewReader(parameters.Encode()))
//	if err != nil {
//		log.Fatal(err)
//	}
//	refer := "https://119.97.153.194:85/CFluxStatistic.php?custDate=0&startDateTime=2021-10-22%2007%3A00&endDateTime=2021-10-22%2008%3A00&schedule=0&sourceType=all&sourceIp=&sourceUser=&sourceGroup=&sub_group=1&appTypeText=%E6%89%80%E6%9C%89%E5%BA%94%E7%94%A8&appType=0&appName=0&stat_method=app_type&rank_base=total_flow&graph_type=0&view_num=10&formState=wrap&curPage=1&act=report&type=view&sid=" + Cookie
//	req.Header.Set("Accept", "text/html, */*;q=0.01")
//	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
//	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
//	req.Header.Set("Connection", "keep-alive")
//	req.Header.Set("Content-Length", "369")
//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
//	req.Header.Set("Cookie", "PHPSESSID="+Cookie)
//	req.Header.Set("Host", "119.97.153.194")
//	req.Header.Set("Origin", "https")
//	req.Header.Set("Referer", refer)
//	//req.Header.Set("sec-ch-ua",'')
//	req.Header.Set("sec-ch-ua-mobile", "?0")
//	//req.Header.Set("sec-ch-ua-platform"," "Windows"")
//	req.Header.Set("Sec-Fetch-Dest", "empty")
//	req.Header.Set("Sec-Fetch-Mode", "cors")
//	req.Header.Set("Sec-Fetch-Site", "same-origin")
//	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36")
//	req.Header.Set("X-Requested-With", "XMLHttpRequest")
//
//	if err != nil {
//		log.Fatal(err)
//	}
//	resp, err = client.Do(req)
//	fmt.Println(resp.Header)
//	fmt.Println(resp.Cookies())
//	fmt.Println(resp.Request)
//	fmt.Println(resp.Status)
//	fmt.Println(resp.Uncompressed)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//	body, _ = ioutil.ReadAll(resp.Body)
//	return body
//}

func GetVerifyCode(user, pass, softid, codetype, len_min string, file []byte) gjson.Result {
	client := &http.Client{}
	var req *http.Request
	var resp *http.Response
	var err error
	var body []byte
	urlString := "http://upload.chaojiying.net/Upload/Processing.php"
	parameters := url.Values{}
	parameters.Add("user", user)
	parameters.Add("pass", pass)
	parameters.Add("softid", softid)
	//http://www.chaojiying.com/price.html
	parameters.Add("codetype", codetype)
	parameters.Add("len_min", len_min)
	parameters.Add("file_base64", getEncodedBase64(file))
	req, err = http.NewRequest("POST", urlString, strings.NewReader(parameters.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0)")
	req.Header.Set("Connection", "Keep-Alive")
	if err != nil {
		log.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var b = string(body)
	r := gjson.Parse(b)
	if r.Get("err_str").String() != "OK" && "0" != r.Get("err_no").String() {
		fmt.Printf("Inquire Faild,Failed code is %s", r.Get("err_no").String())
		return r
	}
	return r
}

func Decode(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func GetContent(begin, end string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Click(`//*[@id="Accordion1"]/div[1]/div[2]/div/dl/dd[3]/a`),
		chromedp.Sleep(1 * time.Second),
		chromedp.Click(`//*[@id="fush"]`),
		chromedp.Sleep(500 * time.Millisecond),
		chromedp.Click(`//*[@id="type4"]`),
		chromedp.Sleep(500 * time.Millisecond),
		chromedp.SendKeys(`//*[@id="timpSelector"]`, "自定义"),
		chromedp.Sleep(500 * time.Millisecond),
		chromedp.Click(`//*[@id="timpSelector"]`),
		chromedp.Sleep(500 * time.Millisecond),
		chromedp.Click(`//*[@id="cpStart"]`),
		chromedp.SendKeys(`//*[@id="cpStart"]`, "\b\b\b\b\b\b\b\b\b\b\b\b\b"+begin),
		chromedp.Click(`//*[@id="cpEnd"]`),
		chromedp.SendKeys(`//*[@id="cpEnd"]`, "\b\b\b\b\b\b\b\b\b\b\b\b\b"+end),
		chromedp.Click(`//*[@id="sure"]`),
		chromedp.Sleep(500 * time.Microsecond),
		chromedp.Click(`//*[@id="querybtn"]/span/strong`),
		chromedp.Sleep(500 * time.Millisecond),
		chromedp.Click(`//*[@id="newtab"]`),
		chromedp.Sleep(200 * time.Millisecond),
		chromedp.Click(`//*[@id="sure"]`),
	}
}
func GetTime() (string, string) {
	subOneHour, _ := time.ParseDuration("-1h")
	now := time.Now()
	ago := now.Add(subOneHour)
	end := now.Format("15:04")
	begin := ago.Format("15:04")
	return begin, end
}

func InqueryData(begin, end string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Sleep(1 * time.Second),
		chromedp.Click(`//*[@id="querybtn"]/span/strong`),
		chromedp.Sleep(500 * time.Millisecond),
		chromedp.Click(`//*[@id="cpStart"]`),
		chromedp.SendKeys(`//*[@id="cpStart"]`, "\b\b\b\b\b\b"+begin),
		chromedp.Click(`//*[@id="cpEnd"]`),
		chromedp.SendKeys(`//*[@id="cpEnd"]`, "\b\b\b\b\b\b"+end),
		chromedp.Click(`//*[@id="sure"]`),
	}
}

func GetIpSlice(taskCtx context.Context) []IP {
	SLice := make([]*string, 0, 10)
	//fmt.Println(Cookie)
	for i := 1; i < 11; i++ {
		s := new(string)
		_ = chromedp.Run(taskCtx, GetTrString(i, s))
		SLice = append(SLice, s)
	}
	ip := ParseSlice(SLice)
	return ip
}

func Login(taskCtx context.Context) bool {
	var chromeNodes []*cdp.Node
	picByte := new([]byte)
	user := "admin"
	passwd := "5Hfw1!2@h!&!"
	err := chromedp.Run(taskCtx, InputUserPwd(user, passwd, picByte))
	if err != nil {
		fmt.Println(err)
	}
	r := GetVerifyCode("2822132073", "fsl2000.", "3a90e8c04865c7d3ba2526ff47e9d11b", "1004", "4", *picByte)
	verifyCode := r.Get("pic_str").Str
	//fmt.Println(verifyCode)
	taskCtx, _ = context.WithTimeout(taskCtx, 5*time.Second)
	err = chromedp.Run(taskCtx, InputVerifyCode(verifyCode))
	if err != nil {
		log.Fatal(err)
	}
	err = chromedp.Run(taskCtx, chromedp.Nodes(`//*[@id="user_id"]`, &chromeNodes))
	if err != nil || chromeNodes == nil {
		picByte = new([]byte)
		e := Refund("2822132073", "fsl2000.", r.Get("pic_id").String(), "3a90e8c04865c7d3ba2526ff47e9d11b")
		if e != "" {
			fmt.Printf("Refund Error: %s\n", e)
		} else {
			fmt.Println("Success Refund !")
		}
		return false
	} else {
		log.Println("Success Login!")
		return true
	}
}

func attemptLogin() (context.Context, context.CancelFunc) {
	for true {
		taskCtx, cancel := GetChromedp(context.Background())
		ok := Login(taskCtx)
		if !ok {
			fmt.Println("登录失败!")
			cancel()
			time.Sleep(10 * time.Second)
		} else {
			return taskCtx, cancel
		}
	}
	return nil, nil
}

func Run() {
	taskCtx, cancel := attemptLogin()
	if taskCtx == nil {
		fmt.Println("登录失败!")
		return
	}
	defer cancel()
	err := chromedp.Run(taskCtx, GetContent(GetTime()))
	i := 0
	for true {
		fmt.Println(i)
		b, e := GetTime()
		err = chromedp.Run(taskCtx, InqueryData(b, e))
		if err != nil {
			log.Fatalln(err)
		}
		ip := GetIpSlice(taskCtx)
		msg := fmt.Sprintf("B%s --> %s 数据如下\n%-25s%-25s%-25s%-25s\n", b, e, "IP", "发送", "接收", "all")
		for _, i := range ip {
			msg = msg + fmt.Sprintf("%-25s%-20s%-20s%-20s\n", i.ip, i.up, i.down, i.all)
		}
		fmt.Printf("%s --> %s 数据如下\n", b, e)
		if strings.Contains(b, "00") || strings.Contains(e, "00") {
			fmt.Println(msg)
			SendDingMsg(msg)
		} else {
			fmt.Println(msg)
		}
		time.Sleep(60 * time.Second)
		i++
	}
	if err != nil {
		fmt.Println(err)
	}
}
func SendDingMsg(msg string) {
	//请求地址模板
	webHook := `https://oapi.dingtalk.com/robot/send?access_token=2e83953407705096cf645dbc21cbcab9e3047ddaecb1a2afb3a7cb8091a40435`
	content := `{"msgtype": "text",
		"text": {"content": "` + msg + `test"}
	}`
	//创建一个请求
	req, err := http.NewRequest("POST", webHook, strings.NewReader(content))
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	//设置请求头
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	//发送请求
	resp, err := client.Do(req)
	//关闭请求
	defer resp.Body.Close()

	if err != nil {
		fmt.Println(err)
	}
}

var Cookie string

func main() {
	Run()
}
