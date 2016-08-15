/**
	本模块负责用户授权，并获获取授权的用户信息。
**/
package oauth

import (
	util "github.com/zdq007/wechat/common"
	"github.com/zdq007/wechat/wechat"
	"errors"
	"fmt"
)

type AuthError struct {
	Errcode int    //出错状态码
	Errmsg  string //出错信息
}
type Userinfo struct {
	//	用户的唯一标识
	OpenId string `json:"openid"`
	//	用户昵称
	Nickname string
	//	用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	Sex byte
	//	用户个人资料填写的省份
	Province string
	//	普通用户个人资料填写的城市
	City string
	//	国家，如中国为CN
	Country string
	//	用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	Headimgurl string
	//	只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。详见：获取用户个人信息（UnionID机制）
	Unionid string

	AuthError
}
type Token struct {
	AccessToken  string `json:"access_token"`  //网页授权接口调用凭证,注意：此access_token与基础支持的access_token不同
	ExpiresIn    int64  `json:"expires_in"`    // access_token接口调用凭证超时时间，单位（秒）
	RefreshToken string `json:"refresh_token"` //	用户刷新access_token
	Openid       string //用户唯一标识，请注意，在未关注公众号时，用户访问公众号的网页，也会产生一个用户和公众号唯一的OpenID
	Scope        string //用户授权的作用域，使用逗号（,）分隔

	AuthError
}

type Client struct {
	storage  *TokenStorage
	code     string
	token    *Token
	userinfo *Userinfo
}

type TokenStorage interface {
	Set(key, value interface{}) error //set value
	Get(key interface{}) interface{}  //get value
}

func NewOAuthClient(code string, storage *TokenStorage) (self *Client) {
	self = new(Client)
	self.code = code
	self.storage = storage
	return
}
func NewOAuthClientByDebug(oppenid, nickename, headimgurl string) (self *Client) {
	self = new(Client)
	self.userinfo = &Userinfo{OpenId: oppenid, Nickname: nickename, Headimgurl: headimgurl}
	return
}

//获取token
func (self *Client) GetAccessToken() string {
	return self.token.AccessToken
}

//获取用户信息
func (self *Client) GetUserInfo() (userinfo *Userinfo, err error) {
	if self.userinfo != nil {
		return self.userinfo, nil
	}
	if self.token == nil {
		if self.token, err = self.exchangeToken(self.code); err != nil {
			return
		}
	}
	err = self.loadUserInfo()
	if err != nil {
		fmt.Println("获取用户信息失败:", err)
	}
	userinfo = self.userinfo

	return
}

//code交换token
func (self *Client) exchangeToken(code string) (*Token, error) {
	var token Token
	err := util.GetJSON(fmt.Sprintf(wechat.WX_OAUATH_TOKEN, wechat.WX_APP_ID, wechat.WX_APP_SECRET, code), &token)
	if err != nil {
		return nil, err
	}
	if token.Errmsg != "" {
		return nil, errors.New(token.Errmsg)
	}
	return &token, err
}

//拉取用户信息
func (self *Client) loadUserInfo() (err error) {
	var userinfo Userinfo
	err = util.GetJSON(fmt.Sprintf(wechat.WX_OAUTH_USERINFO, self.token.AccessToken, self.token.Openid), &userinfo)
	if err != nil {
		return
	}
	fmt.Println("userinfo:", userinfo)
	if userinfo.Errmsg != "" {
		return errors.New(userinfo.Errmsg)
	} else {
		self.userinfo = &userinfo
		return
	}
}

//刷新token
func (self *Client) refreshToken() {

}

//检测token是否过期
func (self *Client) isExpiresToken() {

}
