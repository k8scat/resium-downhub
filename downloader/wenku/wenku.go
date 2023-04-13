package wenku

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/k8scat/downhub/config"
	"github.com/k8scat/downhub/util"
	"github.com/tidwall/gjson"
)

const (
	wenkuGetDocInfoAPI = "https://wenku.baidu.com/api/doc/getdocinfo"
	wenkuDownloadAPI   = "https://wenku.baidu.com/user/submit/download"
	wenkuHost          = "wenku.baidu.com"
)

func wenkuHeaders(resourceUrl string) map[string]string {
	return map[string]string{
		"Host":       wenkuHost,
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36 Edg/85.0.564.60",
		"Referer":    resourceUrl,
	}
}

// 获取真实下载链接
func GetLocation(resourceUrl string) (string, error) {
	docInfo, err := getDocInfo(resourceUrl)
	if err != nil {
		return "", err
	}
	values := &url.Values{
		"doc_id":           {docInfo["docID"]},
		"storage":          {"0"},
		"downloadToken":    {docInfo["downloadToken"]},
		"req_vip_free_doc": {docInfo["isVipFreeDoc"]},
		"useTicket":        {"1"},
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 停止重定向，直接返回下载连接
		},
	}

	req, err := http.NewRequest(http.MethodPost, wenkuDownloadAPI, strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", wenkuHost)
	req.Header.Add("Referer", resourceUrl)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36 Edg/85.0.564.60")
	req.AddCookie(&http.Cookie{
		Name:  "BDUSS",
		Value: config.Config.BDUSS,
	})

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	location := resp.Header.Get("Location")
	if location == "" {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("parse location failed: %+v, raw: %s", resp, string(body))
	}
	return location, nil
}

// 获取文档信息
func getDocInfo(resourceUrl string) (info map[string]string, err error) {
	var result [][]string
	result, err = util.QuickRegexp(resourceUrl, `wenku\.baidu\.com\/view\/([a-z0-9]+)`)
	if err != nil {
		return
	}
	docID := result[0][1]
	params := &url.Values{
		"callback": {"json"},
		"doc_id":   {docID},
	}

	client := http.DefaultClient
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, wenkuGetDocInfoAPI, nil)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", wenkuHost)
	req.Header.Add("Referer", resourceUrl)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36 Edg/85.0.564.60")
	req.AddCookie(&http.Cookie{
		Name:  "BDUSS",
		Value: config.Config.BDUSS,
	})
	req.URL.RawQuery = params.Encode()

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result, err = util.QuickRegexp(string(b), `\/\*\*\/json\((.*)\)`)
	if err != nil {
		return
	}
	jsonStr := result[0][1]
	isPaymentDoc := gjson.Get(jsonStr, "isPaymentDoc").Int()
	if isPaymentDoc == 1 {
		err = fmt.Errorf("cannot download payment doc")
		return
	}
	downloadToken := gjson.Get(jsonStr, "downloadToken").String()
	isProfessionalDocDoc := gjson.Get(jsonStr, "docInfo.professionalDoc").String()
	var isVipFreeDoc string
	if isProfessionalDocDoc == "" {
		isVipFreeDoc = "0"
	} else {
		isVipFreeDoc = "1"
	}
	info = map[string]string{
		"docID":         docID,
		"downloadToken": downloadToken,
		"isVipFreeDoc":  isVipFreeDoc,
	}
	return
}
