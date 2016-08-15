package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/gogap/logrus"
	"github.com/smallnest/goreq"
	"net/http"
)

func GetJSON(url string, response interface{}) (err error) {
	log.Debug("GET: ", url)
	httpResp, err := http.Get(url)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}
	return json.NewDecoder(httpResp.Body).Decode(response)
}

func PostJSON(url string, data []byte, response interface{}) (err error) {
	log.Debug("POST: ", url)
	bufffer := bytes.NewBuffer(data)
	httpResp, err := http.Post(url, "application/json;charset=utf-8", bufffer)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}
	return json.NewDecoder(httpResp.Body).Decode(response)
}

//get请求
func GetREQ(url string) (data []byte, err error) {
	log.Debug("GET: ", url)
	req, body, errarr := goreq.New().Get(url).End()
	if errarr != nil {
		err = errarr[0]
	} else if req.StatusCode == 200 {
		log.Debug("BEACK: ", body)
		if body == "" {
			err = errors.New("微信服务器异常")
		} else {
			data = []byte(body)
		}
	} else {
		err = errors.New(fmt.Sprintf("微信服务器异常,状态码:%d", req.StatusCode))
		log.Error(err)
	}
	return
}

//post请求
func PostREQ(url string, b []byte) (data []byte, err error) {
	log.Debug("POST: ", url, " --DATA: ", string(b))
	request := goreq.New().Post(url).SetHeader("Content-Type", "application/json")
	if b != nil {
		request.SendRawBytes(b)
	}
	req, body, errarr := request.End()
	if errarr != nil {
		err = errarr[0]
	} else if req.StatusCode == 200 {
		if len(body) < 500 {
			log.Debug("BACK: ", body)
		}
		if body == "" {
			err = errors.New("微信服务器异常")
		} else {
			data = []byte(body)
		}
	} else {
		err = errors.New(fmt.Sprintf("微信服务器异常,状态码:%d", req.StatusCode))
		log.Error(err)
	}

	return
}
