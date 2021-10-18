package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func getEncodedBase64(file []byte) string {
	encoded := base64.StdEncoding.EncodeToString(file)
	return encoded
}
func GetVerifyCode(user, pass, softid, codetype, len_min string, file []byte) string {
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
		return r.Get("err_no").String()
	}
	return r.Get("pic_str").String()
}

//func main() {
//	//http.PostForm("http://upload.chaojiying.net/Upload/Processing.php", url.Values{"user": "2822132073"})
//	pic, _ := ioutil.ReadFile("verify.png")
//	str := GetVerifyCode("2822132073", "fsl2000.", "3a90e8c04865c7d3ba2526ff47e9d11b", "1004", "4", pic)
//	fmt.Println(str)
//}

func main() {
	picByte := new([]byte)
	user := "admin"
	passwd := "5Hfw1!2@h!&!"
	options := chromedp.DefaultExecAllocatorOptions[:]
	options = append(options, chromedp.Flag("headless", false), chromedp.Flag("ignore-certificate-errors", "1"))
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()
	taskCtx, cancel := chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	defer cancel()
	err := chromedp.Run(taskCtx,
		chromedp.Navigate("https://119.97.153.194:85/"),
		chromedp.WaitVisible(`#user`, chromedp.ByID),
		chromedp.Sleep(1*time.Second),
		chromedp.SendKeys(`//*[@id="user"]`, user),
		chromedp.SendKeys(`//*[@id="password"]`, passwd),
		chromedp.Sleep(1*time.Second),
		chromedp.Screenshot(`//*[@id="verify_code"]`, picByte),
	)
	//pic.Write(*picByte)
	if err != nil {
		fmt.Println(err)
	}
	verifyCode := GetVerifyCode("2822132073", "fsl2000.", "3a90e8c04865c7d3ba2526ff47e9d11b", "1004", "4", *picByte)
	fmt.Println(verifyCode)
	err = chromedp.Run(taskCtx,
		chromedp.SendKeys(`//*[@id="verify"]`, verifyCode),
		chromedp.WaitVisible(`//*[@id="button"]`),
		chromedp.Click(`//*[@id="button"]`),
		chromedp.Sleep(3*time.Second),
		chromedp.WaitNotPresent(`//*[@id="error_msg"]`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetAllCookies().Do(ctx)
			if err != nil {
				return err
			}
			for i, cookie := range cookies {
				log.Printf("chrome cookie %d: %+v\n", i, cookie)
			}
			return nil
		}),
		chromedp.Sleep(300000000*time.Second),
	)
	if err != nil {
		fmt.Println(err)
	}
}
