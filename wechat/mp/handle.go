package mp
/**
	解析微信消息和事件
	微信TOKEN认证自动返回
	将解析出来的msg对象对调给用户对应的方法
	//部分参考
	//"github.com/arstd/weixin"
 */
import (
	"fmt"
	"io/ioutil"
	"sort"
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	log "github.com/gogap/logrus"
	"encoding/xml"
	"github.com/gogap/errors"
	"strings"
	"github.com/zdq007/wechat/wechat"
	"github.com/zdq007/wechat/wechat/tools"
)
//处理消息
type MSGHandle struct{
	MsgHandleDelegate  MsgHandleDelegater
}


func NewHandle() *MSGHandle{
	return new(MSGHandle)
}
//推送的消息或者事件
func (self *MSGHandle) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method, " : ", req.RequestURI)
	signature := req.FormValue("signature")
	timestamp := req.FormValue("timestamp")
	nonce := req.FormValue("nonce")

	if  !ValidateURL(wechat.WX_DEVELOP_TOKEN, timestamp, nonce, signature){
		http.Error(resp, "validate url error, request not from weixin?", http.StatusUnauthorized)
		return
	}

	fmt.Println("验证通过")
	switch req.Method  {
	case "GET":
		resp.Write([]byte(req.FormValue("echostr")))
	case "POST":
		self.processMessage(resp, req)
	default:
		http.Error(resp, "validate url error, request not from weixin?", http.StatusUnauthorized)
	}
}
//解析消息
func (self *MSGHandle) processMessage(resp http.ResponseWriter, req *http.Request)  {
	// 读取报文
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error(err)
		http.Error(resp, "read body error", http.StatusNotAcceptable)
		return
	}
	defer  req.Body.Close()
	q := req.URL.Query()
	timestamp := q.Get("timestamp")
	nonce := q.Get("nonce")
	encryptType := q.Get("encrypt_type")
	msgSignature := q.Get("msg_signature")
	msg,err:= self.parseBody(encryptType,timestamp,nonce,msgSignature,body);
	if err!=nil{
		log.Error(err)
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("消息体:",string(body))

	fmt.Println(msg)

	// 处理消息
	self.handleMessage(msg)

	//不需特殊回复的一律回复空(表示接收成功)
	resp.Write([]byte(""))

}
// 处理消息
func (self *MSGHandle) handleMessage(msg * Message){
	switch msg.MsgType {
	case MsgTypeEvent:
		self.HandleEvent(msg)
	//事件

	}

}
func (self *MSGHandle) HandleEvent(msg *Message) {
	switch msg.Event {
	case MASSSENDJOBFINISH: //群发消息结束
		if self.MsgHandleDelegate!=nil{
			self.MsgHandleDelegate.MsgSendStatus(msg)
		}
	}
}
//前面三个参数是用来解析加密消息的
func (self *MSGHandle) parseBody(encryptType, timestamp, nonce, msgSignature string, body []byte) (msg *Message, err error) {
	msg = &Message{}
	// 如果报文被加密了，先要验签解密
	if encryptType == "aes" {
		encMsg := &EncMessage{}
		// 解析加密的 xml
		err = xml.Unmarshal(body, encMsg)
		if err != nil {
			return nil, err
		}
		msg.ToUserName = encMsg.ToUserName
		msg.Encrypt = encMsg.Encrypt

		if ValidateURLEnc(wechat.WX_DEVELOP_TOKEN, timestamp, nonce, encMsg.Encrypt) != msgSignature {
			return nil, errors.New("check signature error")
		}

		body, err = util.DecryptMsg(encMsg.Encrypt, wechat.WX_DEVELOP_AESKEY, wechat.WX_APP_ID)
		if err != nil {
			return nil, err
		}
		log.Debug("receive: %s", body)
	}

	// 解析 xml
	err = xml.Unmarshal(body, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

//验证开发者token 密文
func ValidateURLEnc(token, timestamp, nonce, encryptedMsg string) (signature string) {
	strs := sort.StringSlice{token, timestamp, nonce, encryptedMsg}
	strs.Sort()

	buf := make([]byte, 0, len(token)+len(timestamp)+len(nonce)+len(encryptedMsg))

	buf = append(buf, strs[0]...)
	buf = append(buf, strs[1]...)
	buf = append(buf, strs[2]...)
	buf = append(buf, strs[3]...)

	hashsum := sha1.Sum(buf)
	return hex.EncodeToString(hashsum[:])
}

//验证开发者token 明文模式
func ValidateURL(token, timestamp, nonce, signature string) bool {
	if token == "" || timestamp == "" || nonce == "" || signature == "" {
		return false
	}
	strs := sort.StringSlice{token, timestamp, nonce}
	strs.Sort()

	hashsum := sha1.Sum([]byte(strings.Join(strs,"")))
	return hex.EncodeToString(hashsum[:]) == signature
}
