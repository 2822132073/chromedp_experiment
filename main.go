package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"io"
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
func Refund(user, pwd, id, softid string) {
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, err = ioutil.ReadAll(resp.Body)
	r := gjson.Parse(string(body))
	if r.Get("err_str").String() != "OK" {
		fmt.Printf("refund Error : %s\n", r.Get("err_str").String())
	}
}
func InputVerifyCode(verifyCode string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.SendKeys(`//*[@id="verify"]`, verifyCode),
		chromedp.WaitVisible(`//*[@id="button"]`),
		chromedp.Click(`//*[@id="button"]`),
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
func GetContent(begin, end string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Sleep(3 * time.Second),
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
		chromedp.SetValue(`//*[@id="cpStart"]`, begin),
		chromedp.SetValue(`//*[@id="cpEnd"]`, end),
		//chromedp.Click(`//*[@id="cpStart"]`),
		//chromedp.SendKeys(`//*[@id="cpStart"]`, "\b\b\b\b\b\b\b\b\b\b\b\b\b"+begin),
		//chromedp.Click(`//*[@id="cpEnd"]`),
		//chromedp.SendKeys(`//*[@id="cpEnd"]`, "\b\b\b\b\b\b\b\b\b\b\b\b\b"+end),
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
	for i := 1; i < 11; i++ {
		s := new(string)
		_ = chromedp.Run(taskCtx, GetTrString(i, s))
		SLice = append(SLice, s)
	}
	ip := ParseSlice(SLice)
	return ip
}
func Login(taskCtx context.Context) bool {
	s := new(string)
	picByte := new([]byte)
	user := "admin"
	passwd := "5Hfw1!2@h!&!"
	err := chromedp.Run(taskCtx, InputUserPwd(user, passwd, picByte))
	if err != nil {
		fmt.Println(err)
	}
	r := GetVerifyCode("2822132073", "fsl2000.", "3a90e8c04865c7d3ba2526ff47e9d11b", "1004", "4", *picByte)
	verifyCode := r.Get("pic_str").Str
	err = chromedp.Run(taskCtx, InputVerifyCode(verifyCode))
	if err != nil {
		log.Fatal(err)
	}
	err = chromedp.Run(taskCtx, chromedp.InnerHTML(`/`, s))
	if err != nil {
		log.Fatal(err)
	}
	ok, err := regexp.MatchString(".*<input class=\"input_text\" type=\"password\" autocomplete=\"off\" disablautocomplete=\"\" id=\"password\" style=\"background-color: rgb(226, 237, 252);\">.*", *s)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		return true
	} else {
		Refund("2822132073", "fsl2000.", r.Get("pic_id").String(), "3a90e8c04865c7d3ba2526ff47e9d11b")
		return false
	}
}
func attemptLogin() (context.Context, context.CancelFunc) {
	for true {
		taskCtx, cancel := GetChromedp(context.Background())
		ok := Login(taskCtx)
		if !ok {
			cancel()
			fmt.Println("Login Failed , log in again !")
			time.Sleep(3 * time.Second)
		} else {
			fmt.Println("Login Success !!")
			return taskCtx, cancel
		}
	}
	return nil, nil
}
func Run() {
	taskCtx, cancel := attemptLogin()
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if err != nil {
		fmt.Println(err)
	}
}
func main() {
	Run()
}
