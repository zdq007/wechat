package mp
/**
	主动调用微信接口在此封装
 */
import (
	"github.com/zdq007/wechat/common"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	JSON "github.com/bitly/go-simplejson"
	"time"
	"github.com/zdq007/wechat/wechat"
)

/**
获取永久素材列表
*/
func BatchgetMaterial(offsize, size int) ([]byte, error) {
	mper.syncCheckToken()
	params := JSON.New()
	params.Set("type", "news")
	params.Set("offset", offsize)
	params.Set("count", size)
	data, _ := params.MarshalJSON()
	return common.PostREQ(fmt.Sprintf(wechat.BATCHGET_MATERIAL, mper.token.AccessToken), data)
}

//签名URL 授予JSSDK调用权限
func SignJSSDK(url string) *SignRespone {
	mper.syncCheckToken()
	noncestr := common.GetRandomString(12)
	timestamp := time.Now().Unix()
	singstr := fmt.Sprintf(wechat.SignJSSDKTemplate, mper.ticket.TicketStr, noncestr, timestamp, url)
	hsum := sha1.Sum([]byte(singstr))
	fmt.Println("jssdk Signature:", hex.EncodeToString(hsum[:]))
	return &SignRespone{Appid: wechat.WX_APP_ID, Noncestr: noncestr, Timestamp: timestamp, Signature: hex.EncodeToString(hsum[:])}
}

