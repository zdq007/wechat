/**
	本模块负责调用微信提供的支付api
**/
package pay

import (
	"time"

	"github.com/zdq007/wechat/common"
)

//微信公众号支付－统一下单
func WXMPPay(tradename, tradeno, ip, notigyUrl, openid string, totalfee int64) (paySign *JSPaySign, err error) {
	//调用微信统一下单接口
	var req = &UnifiedorderReq{
		Body:           tradename,
		OutTradeNo:     tradeno,
		TotalFee:       totalfee,
		SpbillCreateIp: ip,
		NotigyUrl:      notigyUrl,
		TradeType:      "JSAPI",
		Openid:         openid,
	}
	var res *UnifiedorderResp

	res, err = UnifiedOrder(req)
	if err != nil {
		return nil, err
	}
	//微信支付签名
	paySign = &JSPaySign{
		AppId:     res.Appid,
		Timestamp: time.Now().Unix(),
		NonceStr:  common.GetRandomString(12),
		Package:   "prepay_id=" + res.PrepayId,
		SignType:  "MD5",
		PaySign:   "",
	}
	paySign.PaySign = SignObject(paySign)

	return
}

//微信扫码支付 统一下单 模式二
func WXQRCodePay(tradename, tradeno, ip, notigyUrl string, totalfee int64) (res *UnifiedorderResp, err error) {
	//调用微信统一下单接口
	var req = &UnifiedorderReq{
		Body:           tradename,
		OutTradeNo:     tradeno,
		TotalFee:       totalfee,
		SpbillCreateIp: ip,
		NotigyUrl:      notigyUrl,
		TradeType:      "NATIVE",
		ProductId:      tradeno,
	}

	res, err = UnifiedOrder(req)
	if err != nil {
		return nil, err
	}

	return
}

//微信支付退款
func WXPayBack(tradeno string) {

}
