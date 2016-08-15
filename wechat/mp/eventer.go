package mp
/**
	各种事件的封装类
 */
//事件类型
type EventType string
// 微信支持的事件类型
const (
	MASSSENDJOBFINISH   EventType = "MASSSENDJOBFINISH"	//群发事件推送群发结果
)

