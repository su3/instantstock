package ctrl

import (
	"encoding/csv"
	"fmt"
	"instantstock/conf"
	"instantstock/util"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/axgle/mahonia"
)

type Fetch struct{}

func (c *Fetch) FetchCodes(codeList []string) {
	params := c.codesToParams(codeList)

	t := time.Now().Unix()
	timestamp := fmt.Sprintf("%d", t)
	URI := "/financehq/list=" + params
	accessKey := conf.SAE_ACCESS_KEY
	secretKey := conf.SAE_SECRET_KEY
	rawString := "GET" + "\n" + URI + "\n" + "x-sae-accesskey:" + accessKey + "\n" + "x-sae-timestamp:" + timestamp

	u := util.Util{}
	auth := "SAEV1_HMAC_SHA256 " + u.ComputeHmac256(rawString, secretKey)

	url := "http://g.sae.sina.com.cn" + URI
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-sae-accesskey", accessKey)
	req.Header.Set("x-sae-timestamp", timestamp)
	req.Header.Set("Authorization", auth)
	res, _ := client.Do(req)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll", err)
	}
	// fmt.Println(string(body))
	html := string(body)
	dec := mahonia.NewDecoder("gbk")
	html = dec.ConvertString(html)
	reader := csv.NewReader(strings.NewReader(html))
	csvData, _ := reader.ReadAll()
	for index := range csvData {
		ss := csvData[index]
		fmt.Println(ss)
	}
}

func (c *Fetch) codesToParams(codeList []string) string {
	fixedCodeList := []string{}
	for _, item := range codeList {
		var code string
		if strings.HasPrefix(item, "6") {
			code = "sh" + item
		} else {
			code = "sz" + item
		}
		fixedCodeList = append(fixedCodeList, code)
	}
	return strings.Join(fixedCodeList, ",")
}
