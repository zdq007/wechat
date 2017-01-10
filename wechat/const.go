package wechat

//微信平台配置
var (
	//公众号id
	WX_APP_ID string
	//公众号秘钥
	WX_APP_SECRET string
	//开发者TOKEN
	WX_DEVELOP_TOKEN string
	// 消息加解密密钥
	WX_DEVELOP_AESKEY []byte

	//[支付]
	//商户id
	WX_PAY_MCHID string
	//商户秘钥
	WX_PAY_KEY string
)

//oauth
const (
	//OAUATH2认证 snsapi_userinfo snsapi_base
	WX_OAUATH_CODE    = `https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect`
	WX_OAUATH_TOKEN   = `https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code`
	WX_OAUTH_USERINFO = `https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN`
)

//mp accesstoken ticket
const (
	//获取TICKET
	WX_URL_MP_TICKET = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=jsapi&access_token=%s"
	//获取Access_token
	WX_URL_MP_ACCESSTOKEN = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

//pay
const (
	//统一下单url
	WX_PAY_UNIFIEDORDER = `https://api.mch.weixin.qq.com/pay/unifiedorder`

	//支付接口返回状态
	ReturnCodeSuccess = "SUCCESS"
	ReturnCodeFail    = "FAIL"

	//支付接口业务状态
	ResultCodeSuccess = "SUCCESS"
	ResultCodeFail    = "FAIL"
)

//mp
const (
	//注意这里参数名必须全部小写，且必须有序
	SignJSSDKTemplate = `jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s`
	SignAdderTemplate = `accesstoken=%s&appid=%s&noncestr=%s&timestamp=%d&url=%s`
	//获取永久素材
	BATCHGET_MATERIAL = `https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token=%s`
)
