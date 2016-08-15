package mp
/**
	各种消息的封装类
 */

// 接收普通消息
type MsgType string
// 微信支持的消息类型
const (
	MsgTypeText       MsgType = "text"       // 文本消息
	MsgTypeImage      MsgType = "image"      // 图片消息
	MsgTypeVoice      MsgType = "voice"      // 语音消息
	MsgTypeVideo      MsgType = "video"      // 视频消息
	MsgTypeShortVideo MsgType = "shortvideo" // 小视频消息
	MsgTypeLocation   MsgType = "location"   // 地理位置消息
	MsgTypeLink       MsgType = "link"       // 链接消息
	MsgTypeMusic      MsgType = "music"      // 音乐消息
	MsgTypeNews       MsgType = "news"       // 图文消息
	MsgTypeWXCard     MsgType = "wxcard"     // 卡券，客服消息时使用
	MsgTypeEvent      MsgType = "event"      // 事件推送
)
