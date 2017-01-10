/**
	本模块负责发送http请求比获取响应数据，对数据做安全验证。
**/
package pay

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	util "github.com/zdq007/wechat/common"
	"github.com/zdq007/wechat/wechat"
	"net/http"
)

/**业务方法，对微信平台的调用接口**/

/**
* 统一下单接口 公众号支付
**/
func UnifiedOrder(req *UnifiedorderReq) (res *UnifiedorderResp, err error) {
	req.Appid = wechat.WX_APP_ID
	req.MchId = wechat.WX_PAY_MCHID
	req.NonceStr = util.GetRandomString(12)
	res = new(UnifiedorderResp)
	if err = PostXML(wechat.WX_PAY_UNIFIEDORDER, *req, res); err != nil {
		return
	}
	if res.ResultCode != wechat.ResultCodeSuccess {
		err = errors.New(fmt.Sprintf("return_msg:%v, order_num: %v, result_code: %v, result_msg: %v", res.ReturnMsg, req.OutTradeNo, res.ResultCode, res.ErrCodeDes))
		return
	}
	if err = VilidataResp(*res); err != nil {
		return
	}
	return res, nil
}

/**工具方法**/

//校验返回结构
func VilidataResp(res Resp) (err error) {
	// 判断协议状态
	ReturnCode := res.GetReturnCode()
	if ReturnCode != wechat.ReturnCodeSuccess {
		err = errors.New(fmt.Sprintf("return_code: %v, return_msg: %v", ReturnCode, res.GetReturnMsg()))
		return
	}
	appId := res.GetAppId()
	// 安全考虑, 做下验证
	if appId != wechat.WX_APP_ID {
		err = fmt.Errorf("appid mismatch, have: %v, want: %v", appId, wechat.WX_APP_ID)
		return
	}
	mchId := res.GetMchId()
	if mchId != wechat.WX_PAY_MCHID {
		err = fmt.Errorf("mch_id mismatch, have: %v, want: %v", mchId, wechat.WX_PAY_MCHID)
		return
	}

	signature := res.GetSignature()
	if signature == "" {
		err = errors.New("no sign parameter")
		return
	}

	signature1 := SignObject(res)
	if signature != signature1 {
		err = fmt.Errorf("check signature failed, input: %v, local: %v", signature, signature1)
		return
	}
	fmt.Println("校验结构成功!")
	return
}

// 微信支付通用请求方法.
//  注意: err == nil 表示协议状态都为 SUCCESS(return_code == SUCCESS).
func PostXML(url string, req interface{}, res Resp) (err error) {
	body := BodyBufferPool.Get().(*bytes.Buffer)
	body.Reset()
	defer BodyBufferPool.Put(body)
	err = SignObjectToXML(req, body)
	if err != nil {
		return
	}
	httpResp, err := http.Post(url, "text/xml;charset=utf-8", body)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}
	err = xml.NewDecoder(httpResp.Body).Decode(res)
	if err != nil {
		return
	}
	fmt.Println("res:", res)
	return
}
