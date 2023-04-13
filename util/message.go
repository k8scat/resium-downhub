package util

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var DingtalkClient *Dingtalk

const dingtalkSendAPI = "https://oapi.dingtalk.com/robot/send"

type Dingtalk struct {
	AccessToken string
	Secret      string
}

func NewDingtalk(accessToken, secret string) (client *Dingtalk, err error) {
	if accessToken == "" || secret == "" {
		err = errors.New("accessToken or secret cannot be empty")
		return
	}
	client = &Dingtalk{
		AccessToken: accessToken,
		Secret:      secret,
	}
	return
}

func (d *Dingtalk) Send(content string, atAll bool, atMobiles []string) error {
	at := map[string]interface{}{
		"isAtAll": atAll,
	}
	if atMobiles != nil && len(atMobiles) > 0 {
		at["atMobiles"] = atMobiles
	}
	data := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
		"at": at,
	}
	b, _ := json.Marshal(data)

	currentTimestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	signStr := sign(currentTimestamp, d.Secret)
	params := &url.Values{
		"access_token": {d.AccessToken},
		"timestamp":    {currentTimestamp},
		"sign":         {signStr},
	}

	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodPost, dingtalkSendAPI, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Send dingtalk message failed: %s, code: %d", string(b), resp.StatusCode)
	}
	return nil
}

func sign(t string, secret string) string {
	s := fmt.Sprintf("%s\n%s", t, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(s))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}

func Ding(content string) {
	if DingtalkClient == nil {
		log.Println("dingtalk client not initialized")
	} else {
		DingtalkClient.Send(content, false, nil)
	}
}
