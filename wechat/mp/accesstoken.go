/**
	微信参数的初始化
	微信jssdk调用需要对发起调用的页面url签名，签名要用到微信分发的token。
	所以第一步获取token 和 ticket ，因为token有失效时间，这个模块还负责失效检查，保证token时刻有效
**/
package mp

import (
	"encoding/base64"
	"fmt"
	"github.com/go-ini/ini"
	log "github.com/gogap/logrus"
	util "github.com/zdq007/wechat/common"
	"github.com/zdq007/wechat/wechat"
	"os"
	"regexp"
	"sync"
	"time"
)

var (
	TimeStep = 60 * 20

	mper   *MP
	locker = new(sync.RWMutex)
)

type MP struct {
	token  *Token
	ticket *Ticket
					   //过期时间
	expiresTime   int64
	storageConfig *ini.File
	timerGear     byte // 0 1秒一次 1 2分钟一次
}

type SignRespone struct {
	Appid       string
	Noncestr    string
	Timestamp   int64
	Signature   string
	AccessToken string
	JsapiTicket string
	Url         string
}
type Token struct {
	AccessToken string `json:"access_token"`
	Expiress    int64  `json:"expires_in"`
}

type Ticket struct {
	TicketStr string `json:"ticket"`
	Expiress  int64  `json:"expires_in"`
}

/**
微信初始化方法
appId  			公众号id
appSecret 		公众号秘钥
token  			开发者token
encodingAESKey  密文模式下的加密串,明文模式下传空即可
*/
func Initialize(appId, appSecret, token, encodingAESKey string) {
	if matched, err := regexp.MatchString("^wx[0-9a-f]{16}$", appId); err != nil || !matched {
		log.Error("appId format error: %s", err)
	}
	if matched, err := regexp.MatchString("^[0-9a-f]{32}$", appSecret); err != nil || !matched {
		log.Error("appSecret format error: %s", err)
	}
	if matched, err := regexp.MatchString("^[0-9a-zA-Z]{3,32}$", token); err != nil || !matched {
		log.Error("token format error: %s", err)
	}
	if matched, err := regexp.MatchString("^[0-9a-zA-Z]{43}$", encodingAESKey); err != nil || !matched {
		log.Error("encodingAESKey format error: %s", err)
	}
	wechat.WX_APP_ID = appId
	wechat.WX_APP_SECRET = appSecret
	wechat.WX_DEVELOP_TOKEN = token

	if encodingAESKey != "" {
		var err error
		wechat.WX_DEVELOP_AESKEY, err = base64.StdEncoding.DecodeString(encodingAESKey + "=")
		if err != nil {
			log.Error("appSecret config error: %s", err)
			os.Exit(0)
		}
	}
	//获取access token
	fetchAccessToken()
}

/**
微信初始化方法
appId  			公众号id
appSecret 		公众号秘钥
*/
func InitializeSimple(appId, appSecret string) {
	if matched, err := regexp.MatchString("^wx[0-9a-f]{16}$", appId); err != nil || !matched {
		log.Error("appId format error: %s", err)
	}
	if matched, err := regexp.MatchString("^[0-9a-f]{32}$", appSecret); err != nil || !matched {
		log.Error("appSecret format error: %s", err)
	}

	wechat.WX_APP_ID = appId
	wechat.WX_APP_SECRET = appSecret

	//获取access token
	fetchAccessToken()
}

//开启mp后台worker，定时检测更新access token
func fetchAccessToken() {
	mper = new(MP)
	wxpro, err := ini.LooseLoad("wx.proc")
	if err != nil {
		log.Error("read sign config err:", err)
	} else {
		weSec := wxpro.Section("wechat")
		token := weSec.Key("token").MustString("")
		ticket := weSec.Key("ticket").MustString("")
		mper.expiresTime = weSec.Key("extime").MustInt64(0)
		mper.token = &Token{AccessToken: token, Expiress: mper.expiresTime}
		mper.ticket = &Ticket{TicketStr: ticket, Expiress: mper.expiresTime}
	}
	mper.storageConfig = wxpro
	go mper.tokenTimer()
}

//定期检查公众号token是否过期
func (self *MP) tokenTimer() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			log.Debug("AccessToken 定时检测开始")
			if self.isExpiress() {
				log.Debug("AccessToken 过期,重新获取")
				//更新token
				self.updateToken()
			} else {
				log.Debug("AccessToken 未过期")
			}
		//判定是否修改accesstoken检测时间间隔
			ticker = self.getTicker(ticker)
			log.Debug("AccessToken 定时检测结束")
		}
	}
}

//定时任务
func (self *MP) getTicker(ticker *time.Ticker) *time.Ticker {
	var timerGear byte //0档
	if !self.isExpiress() {
		//没有过期修改为1档
		timerGear = 1
	}
	if self.timerGear == timerGear {
		//如果已经是1档返回当前ticker
		return ticker
	} else {
		//产生对应档位档新ticker
		self.timerGear = timerGear
		ticker.Stop()
		if timerGear == 1 {
			log.Debug("AccessToken 未过期,定时器调整为6分钟一次")
			return time.NewTicker(60 * 6 * time.Second)
		} else {
			log.Debug("AccessToken 过期,定时器调整为6秒一次")
			return time.NewTicker(6 * time.Second)
		}
	}
}

//检测是否过期
func (self *MP) isExpiress() bool {
	return self.expiresTime < (time.Now().Unix() + int64(TimeStep))
}

//更新token
func (self *MP) updateToken() {
	self.token = self.getAccessToken()
	if self.token != nil {
		self.expiresTime = self.token.Expiress + time.Now().Unix()
		self.ticket = self.getJsApiTicket(self.token.AccessToken)
	}
	if self.token != nil && self.ticket != nil {
		//保存token到文件
		weSec := self.storageConfig.Section("wechat")
		weSec.NewKey("token", self.token.AccessToken)
		weSec.NewKey("ticket", self.ticket.TicketStr)
		weSec.NewKey("extime", fmt.Sprintf("%d", self.expiresTime))
		self.storageConfig.SaveTo("wx.proc")
		log.Debug("更新accesstoken成功")
	}
}

//获取JSAPI调用票据
func (self *MP) getJsApiTicket(accessToken string) *Ticket {
	url := fmt.Sprintf(wechat.WX_URL_MP_TICKET, accessToken)
	ticket := &Ticket{}
	err := util.GetJSON(url, ticket)
	if err == nil {
		return ticket
	}
	log.Error("获取JSAPI调用票据:", err)
	return nil
}

//获取公众号token
func (self *MP) getAccessToken() *Token {
	url := fmt.Sprintf(wechat.WX_URL_MP_ACCESSTOKEN, wechat.WX_APP_ID, wechat.WX_APP_SECRET)
	token := &Token{}
	err := util.GetJSON(url, token)
	if err == nil {
		return token
	}
	log.Error("获取公众号token:", err)
	return nil
}

//检测token是否过期。过期则同步更新token
func (self *MP) syncCheckToken() {
	locker.Lock()
	defer locker.Unlock()
	if self.isExpiress() {
		self.updateToken()
	}
}
