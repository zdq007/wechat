package pay

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/zdq007/wechat/wechat"
	"hash"
	"reflect"
	"sort"
	"strings"
	"sync"
)

//对象池
var BodyBufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 16<<10)) // 16KB
	},
}

//签名参数并生成xml
func SignObjectToXML(paramObj interface{}, bodybuff *bytes.Buffer) error {
	t := reflect.TypeOf(paramObj)
	v := reflect.ValueOf(paramObj)
	if _, err := bodybuff.WriteString("<xml>"); err != nil {
		return err
	}

	fieldNum := t.NumField()
	params := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		field := t.Field(i)
		xmltag := field.Tag.Get("xml")
		isneed := field.Tag.Get("need")
		if xmltag == "" {
			xmltag = field.Name
		}
		if xmltag == "XMLName" || xmltag == "xml" || xmltag == "sign" {
			continue
		}
		vi := v.Field(i).Interface()
		val := fmt.Sprintf("%v", vi)
		add := false
		if valstr, ok := vi.(string); ok {
			if isneed == "1" || valstr != "" {
				add = true
			}
		} else {
			if isneed == "1" || val != "0" {
				add = true
			}
		}

		if add {
			bodybuff.WriteByte('<')
			bodybuff.WriteString(xmltag)
			bodybuff.WriteByte('>')
			bodybuff.WriteString(val)
			bodybuff.WriteString("</")
			bodybuff.WriteString(xmltag)
			bodybuff.WriteByte('>')
			params = append(params, xmltag+"="+val)
		}
	}
	sort.Strings(params)
	params = append(params, "key="+wechat.WX_PAY_KEY)
	stringSign := strings.Join(params, "&")
	fmt.Println("stringSign:", stringSign)
	sum := md5.Sum([]byte(stringSign))
	//fmt.Println("stringSign:", stringSign)
	sign := strings.ToUpper(hex.EncodeToString(sum[:]))
	//fmt.Println("sign:", sign)

	bodybuff.WriteString("<sign>")
	bodybuff.WriteString(sign)
	bodybuff.WriteString("</sign>")
	if _, err := bodybuff.WriteString("</xml>"); err != nil {
		return err
	}
	xml := bodybuff.String()
	fmt.Println("xml:", xml)
	return nil
}

//签名不生成xml，主要用在jssdk支付签名
func SignObject(paramObj interface{}) string {
	t := reflect.TypeOf(paramObj)
	v := reflect.ValueOf(paramObj)

	fieldNum := t.NumField()
	params := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		field := t.Field(i)
		xmltag := field.Tag.Get("xml")
		isneed := field.Tag.Get("need")
		if xmltag == "" {
			xmltag = field.Name
		}
		if xmltag == "XMLName" || xmltag == "xml" || xmltag == "sign" {
			continue
		}
		vi := v.Field(i).Interface()
		if valstr, ok := vi.(string); ok {
			if isneed == "1" || valstr != "" {
				params = append(params, xmltag+"="+valstr)
			}
		} else {
			val := fmt.Sprintf("%v", vi)
			if isneed == "1" || val != "0" {
				params = append(params, xmltag+"="+val)
			}
		}
	}
	sort.Strings(params)
	params = append(params, "key="+wechat.WX_PAY_KEY)
	stringSign := strings.Join(params, "&")
	sum := md5.Sum([]byte(stringSign))
	fmt.Println("stringSign:", stringSign)
	sign := strings.ToUpper(hex.EncodeToString(sum[:]))
	fmt.Println("sign:", sign)
	return sign

}

// 微信支付签名.
//  parameters: 待签名的参数集合
//  apiKey:     API密钥
//  fn:         func() hash.Hash, 如果 fn == nil 则默认用 md5.New
func SignMap(parameters map[string]string, apiKey string, fn func() hash.Hash) string {
	ks := make([]string, 0, len(parameters))
	for k := range parameters {
		if k == "sign" {
			continue
		}
		ks = append(ks, k)
	}
	sort.Strings(ks)

	if fn == nil {
		fn = md5.New
	}
	h := fn()
	buf := make([]byte, 256)
	for _, k := range ks {
		v := parameters[k]
		if v == "" {
			continue
		}

		buf = buf[:0]
		buf = append(buf, k...)
		buf = append(buf, '=')
		buf = append(buf, v...)
		buf = append(buf, '&')
		h.Write(buf)
	}
	buf = buf[:0]
	buf = append(buf, "key="...)
	buf = append(buf, apiKey...)
	h.Write(buf)

	signature := make([]byte, h.Size()*2)
	hex.Encode(signature, h.Sum(nil))
	return string(bytes.ToUpper(signature))
}
