package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	product = `https://oapi.dingtalk.com/robot/send?access_token=84618677bead498eeeedc13a5147c3e74976e393e102c03294a2708d9ebff425`
	test    = `https://oapi.dingtalk.com/robot/send?access_token=54c9756ab0bc12626ae7332cef3694f8eb6f9c8111fea9c513da8338bed2ee52`
	LogFile = "brower_auto.log"
)

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

type IP struct {
	ip   string
	up   string
	down string
	all  string
}

func init() {

	developConfig := zap.NewDevelopmentConfig()
	developConfig.OutputPaths = append(developConfig.OutputPaths, LogFile)
	developConfig.ErrorOutputPaths = append(developConfig.OutputPaths, LogFile)
	Logger, _ = developConfig.Build()
	Sugar = Logger.Sugar()
}

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
func SendDingMsg(msg string, webHook string) {
	//请求地址模板
	//webHook := ``
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
	if resp.Status != "200 OK" {
		Sugar.Warn("Send msg error:", resp.Status, resp.Body)
	}
	//关闭请求
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if err != nil {
		Sugar.Warn(err)
	}
}

func ParseSlice(s []*string) []IP {
	lip := make([]IP, 0, 1)
	for _, s := range s {
		var ip IP
		unnecessaryChar, _ := regexp.Compile("[ \\- ,]")
		fs := unnecessaryChar.ReplaceAllString(*s, "")
		ls := strings.Split(fs, "\n")
		if ConvertUnit(ls[6]) > 1 {
			ip.ip = ls[1]
			ip.up = strconv.FormatFloat(ConvertUnit(ls[4]), 'f', 2, 64) + "GB"
			ip.down = strconv.FormatFloat(ConvertUnit(ls[5]), 'f', 2, 64) + "GB"
			ip.all = strconv.FormatFloat(ConvertUnit(ls[6]), 'f', 2, 64) + "GB"
			lip = append(lip, ip)
		}
	}
	return lip
}

func ConvertUnit(s string) float64 {
	n, err := regexp.Compile("[0-9]*")
	if err != nil {
		Sugar.Warn(err)
	}
	if strings.Contains(s, "KB") {
		i, err := strconv.ParseFloat(n.FindString(s), 64)
		if err != nil {
			Sugar.Warn(err)
		}
		// KB转换成MB,再转换成GB
		i = i / 1024 / 1024
		//s = strconv.FormatFloat(i, 'f', 4, 64) + "GB"
		return i
	} else if strings.Contains(s, "MB") {
		i, err := strconv.ParseFloat(n.FindString(s), 64)
		if err != nil {
			Sugar.Warn(err)
		}
		i = i / 1024
		return i
	}
	i, err := strconv.ParseFloat(n.FindString(s), 64)
	if err != nil {
		Sugar.Warn(err)
	}
	return i
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
		Sugar.Warn(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0)")
	req.Header.Set("Connection", "Keep-Alive")
	if err != nil {
		Sugar.Warn(err)
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
		chromedp.Sleep(3 * time.Second),
	}
}
func GetChromedp(ctx context.Context) (context.Context, context.CancelFunc) {
	options := chromedp.DefaultExecAllocatorOptions[:]
	options = append(options, chromedp.Flag("headless", false), chromedp.Flag("ignore-certificate-errors", "1"))
	cctx, cancel := context.WithCancel(ctx)
	c, cancel := chromedp.NewExecAllocator(cctx, options...)
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
		Sugar.Warn(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0)")
	req.Header.Set("Connection", "Keep-Alive")
	if err != nil {
		Sugar.Warn(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		Sugar.Warn(err)
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
		chromedp.ActionFunc(func(ctx context.Context) error {
			Sugar.Info("click 流量统计 success")
			return nil
		}),
		chromedp.Sleep(5 * time.Second),
		chromedp.Click(`//*[@id="fush"]`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			Sugar.Info("click 详细信息 success")
			return nil
		}),
		chromedp.Click(`//*[@id="type4"]`),
		chromedp.SendKeys(`//*[@id="timpSelector"]`, "自定义"),
		chromedp.Click(`//*[@id="timpSelector"]`),
		chromedp.SetValue(`//*[@id="cpStart"]`, begin),
		chromedp.SetValue(`//*[@id="cpEnd"]`, end),
		chromedp.Click(`//*[@id="sure"]`),
		chromedp.Click(`//*[@id="querybtn"]/span/strong`),
		chromedp.Click(`//*[@id="newtab"]`),
		chromedp.Click(`//*[@id="sure"]`),
	}
}
func GetTime() (string, string) {
	subOneHour, _ := time.ParseDuration("-1h")
	now := time.Now()
	ago := now.Add(subOneHour)
	end := now.Format("15:04")
	begin := ago.Format("15:04")
	if now.Hour() >= ago.Hour() {
		return begin, end
	}
	return "", ""
}
func InqueryData(begin, end string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Sleep(1 * time.Second),
		chromedp.Click(`//*[@id="querybtn"]/span/strong`),
		chromedp.Sleep(500 * time.Millisecond),
		chromedp.Click(`//*[@id="cpStart"]`),
		chromedp.SetValue(`//*[@id="cpStart"]`, begin),
		chromedp.Click(`//*[@id="cpEnd"]`),
		chromedp.SetValue(`//*[@id="cpEnd"]`, end),
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
	r := GetVerifyCode("", "", "", "1004", "4", *picByte)
	verifyCode := r.Get("pic_str").Str
	Logger.Info("verify code", zap.String("verifyCode", verifyCode))
	err = chromedp.Run(taskCtx, InputVerifyCode(verifyCode))
	if err != nil {
		Sugar.Warn(err)
	}
	err = chromedp.Run(taskCtx, chromedp.InnerHTML(`/`, s))
	if err != nil {
		Sugar.Warn(err)
	}
	ok := strings.Contains(*s, "欢迎登录")
	if !ok {
		return true
	} else {
		Refund("", "", r.Get("pic_id").String(), "")
		return false
	}
}
func attemptLogin(ctx context.Context) (context.Context, context.CancelFunc) {
	for true {
		taskCtx, cancel := GetChromedp(ctx)
		ok := Login(taskCtx)
		if !ok {
			cancel()
			Logger.Info("Login Failed , log in again !")
			time.Sleep(3 * time.Second)
		} else {
			Logger.Info("Login Success !!")
			return taskCtx, cancel
		}
	}
	return nil, nil
}
func GetDateFromWebAndSendMsg(ctx context.Context, Live chan int, dead chan int) {
	taskCtx, cancel := attemptLogin(ctx)
	err := chromedp.Run(taskCtx, GetContent(GetTime()))
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		for true {
			select {
			case <-Live:
				Logger.Info("Receive signal from channel")
				b, e := GetTime()
				if b != "" {
					err = chromedp.Run(taskCtx, InqueryData(b, e))
					if err != nil {
						Sugar.Error(err)
						return
					}
					msg := GenerateMsg(taskCtx, b, e)
					if strings.Contains(strings.Split(e, ":")[1], "00") {
						fmt.Println(msg)
						SendDingMsg(msg, product)
						Logger.Info("Send msg to AlterGroup")
					} else {
						fmt.Println(msg)
						SendDingMsg(msg, test)
						Logger.Info("Send msg to DebugGroup")
					}
				}
			}
		}
	}()
	go func() {
		for true {
			select {
			case <-dead:
				Sugar.Warn("goroutine receive exit signal,aborting.....")
				cancel()
				return
			}
		}
	}()
}
func GenerateMsg(taskCtx context.Context, b, e string) string {
	ip := GetIpSlice(taskCtx)
	msg := fmt.Sprintf("B%s --> %s 数据如下\n%-25s%-18s%-19s%-20s\n", b, e, "IP", "发送", "接收", "all")
	for _, i := range ip {
		msg = msg + fmt.Sprintf("%-25s%-20s%-20s%-20s\n", i.ip, i.up, i.down, i.all)
	}
	return msg
}
func Run() {
	Live := make(chan int)
	Dead := make(chan int)
	ctx := context.Background()
	GetDateFromWebAndSendMsg(ctx, Live, Dead)
	go func() {
		for true {
			Logger.Info("Send signal to channel")
			Live <- 1
			time.Sleep(60 * time.Second)
		}
	}()
	i := 0
	for true {
		select {
		case <-Live:
			i++
			Logger.Warn("Detect getData dead", zap.Int("number", i))
			if i >= 2 {
				Logger.Warn("Restart GetData goroutine")
				GetDateFromWebAndSendMsg(ctx, Live, Dead)
				Sugar.Warn("Detect last goroutine happen error,send dead signal to the gorouitine")
				Dead <- 1
				i = 0
			}
		default:
			time.Sleep(5 * time.Second)
		}
	}
}

func main() {
	Run()
}
