/**
	支付请求和响应模型
**/
package pay

import (
	"encoding/xml"
)

//need 是否为必填项
//返回接口
type Resp interface {
	GetSignature() string
	GetReturnCode() string
	GetReturnMsg() string
	GetResultCode() string
	GetAppId() string
	GetMchId() string
}

//js发起支付参数
type JSPaySign struct {
	AppId     string `xml:"appId"`
	Timestamp int64  `xml:"timeStamp"`
	NonceStr  string `xml:"nonceStr"`
	Package   string `xml:"package"`
	SignType  string `xml:"signType"`
	PaySign   string
}

//统一下单请求参数
type UnifiedorderReq struct {
	XMLName        xml.Name `xml:"xml" need:"1"`
	Appid          string   `xml:"appid" need:"1"`            //公众账号ID
	MchId          string   `xml:"mch_id" need:"1"`           //商户号
	DeviceInfo     string   `xml:"device_info" need:"0"`      //设备号
	NonceStr       string   `xml:"nonce_str" need:"1"`        //随机字符串
	Sign           string   `xml:"sign" need:"1"`             //签名
	Body           string   `xml:"body" need:"1"`             //商品描述
	Detail         string   `xml:"detail" need:"0"`           //商品详情
	Attach         string   `xml:"attach" need:"0"`           //附加数据
	OutTradeNo     string   `xml:"out_trade_no" need:"1"`     //商户订单号
	FeeType        string   `xml:"fee_type" need:"0"`         //货币类型
	TotalFee       int64    `xml:"total_fee" need:"1"`        //总金额
	SpbillCreateIp string   `xml:"spbill_create_ip" need:"1"` //终端IP
	TimeStart      string   `xml:"time_start" need:"0"`       //交易起始时间
	TimeExpire     string   `xml:"time_expire" need:"0"`      //交易结束时间
	GoodsTag       string   `xml:"goods_tag" need:"0"`        //商品标记
	NotigyUrl      string   `xml:"notify_url" need:"1"`       //通知地址
	TradeType      string   `xml:"trade_type" need:"1"`       //交易类型
	ProductId      string   `xml:"product_id" need:"0"`       //商品ID
	LimitPoay      string   `xml:"limit_pay" need:"0"`        //指定支付方式
	Openid         string   `xml:"openid" need:"0"`           //用户在商户appid下的唯一标识
}

//统一下单返回结果
type UnifiedorderResp struct {
	ReturnCode string `xml:"return_code" need:"1"`  //返回状态码
	ReturnMsg  string `xml:"return_msg" need:"0"`   //返回信息
	Appid      string `xml:"appid" need:"1"`        //公众账号ID
	MchId      string `xml:"mch_id" need:"1"`       //商户号
	DeviceInfo string `xml:"device_info" need:"0"`  //设备号
	NonceStr   string `xml:"nonce_str" need:"1"`    //随机字符串
	Sign       string `xml:"sign" need:"1"`         //签名
	ResultCode string `xml:"result_code" need:"1"`  //业务结果
	ErrCode    string `xml:"err_code" need:"0"`     //错误代码
	ErrCodeDes string `xml:"err_code_des" need:"0"` //错误代码描述
	TradeType  string `xml:"trade_type" need:"1"`   //交易类型
	PrepayId   string `xml:"prepay_id" need:"1"`    //预支付交易会话标识
	CodeUrl    string `xml:"code_url" need:"0"`     //二维码链接
}

func (self UnifiedorderResp) GetReturnCode() string {
	return self.ReturnCode
}
func (self UnifiedorderResp) GetReturnMsg() string {
	return self.ReturnMsg
}
func (self UnifiedorderResp) GetResultCode() string {
	return self.ResultCode
}
func (self UnifiedorderResp) GetSignature() string {
	return self.Sign
}
func (self UnifiedorderResp) GetAppId() string {
	return self.Appid
}
func (self UnifiedorderResp) GetMchId() string {
	return self.MchId
}

//支付回调结果
type PayResp struct {
	ReturnCode    string `xml:"return_code" need:"1"`    //返回状态码
	ReturnMsg     string `xml:"return_msg" need:"0"`     //返回信息
	Appid         string `xml:"appid" need:"1"`          //公众账号ID
	MchId         string `xml:"mch_id" need:"1"`         //商户号
	DeviceInfo    string `xml:"device_info" need:"0"`    //设备号
	NonceStr      string `xml:"nonce_str" need:"1"`      //随机字符串
	Sign          string `xml:"sign" need:"1"`           //签名
	ResultCode    string `xml:"result_code" need:"1"`    //业务结果
	ErrCode       string `xml:"err_code" need:"0"`       //错误代码
	ErrCodeDes    string `xml:"err_code_des" need:"0"`   //错误代码描述
	Openid        string `xml:"openid" need:"1"`         //用户标识
	IsSubscribe   string `xml:"is_subscribe" need:"0"`   //是否关注公众账号
	TradeType     string `xml:"trade_type" need:"1"`     //交易类型
	BankType      string `xml:"bank_type" need:"1"`      //付款银行
	TotalFee      int64  `xml:"total_fee" need:"1"`      //总金额
	FeeType       string `xml:"fee_type" need:"0"`       //货币种类
	CashFee       int    `xml:"cash_fee" need:"1"`       //现金支付金额
	CashFeeType   string `xml:"cash_fee_type" need:"0"`  //现金支付货币类型
	CouponFee     int64  `xml:"coupon_fee" need:"0"`     //代金券或立减优惠金额
	CouponCount   int    `xml:"coupon_count" need:"0"`   //代金券或立减优惠使用数量
	CouponIdN     string `xml:"coupon_id_$n" need:"0"`   //代金券或立减优惠ID
	CouponFeeN    int64  `xml:"coupon_fee_$n" need:"0"`  //单个代金券或立减优惠支付金额
	TransactionId string `xml:"transaction_id" need:"1"` //微信支付订单号
	OutTradeNo    string `xml:"out_trade_no" need:"1"`   //商户订单号
	Attach        string `xml:"attach" need:"0"`         //商家数据包
	TimeEnd       string `xml:"time_end" need:"1"`       //支付完成时间
}

func (self PayResp) GetReturnCode() string {
	return self.ReturnCode
}
func (self PayResp) GetReturnMsg() string {
	return self.ReturnMsg
}
func (self PayResp) GetResultCode() string {
	return self.ResultCode
}
func (self PayResp) GetSignature() string {
	return self.Sign
}
func (self PayResp) GetAppId() string {
	return self.Appid
}
func (self PayResp) GetMchId() string {
	return self.MchId
}
