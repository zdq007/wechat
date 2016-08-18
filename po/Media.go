package po

import (
	"time"
)
/**
	素材组信息
 */
type TMpMedia struct {
	MediaId    string `gorm:"COLUMN:media_id;primary_key"`
	Type       string
	Synctime   int64
	Createtime int64
	Updatetime int64
	NewsList   []*TMpNews
}

func NewTMpMedia(mediaId, typestr string, createtime, updatetime int64) *TMpMedia {
	return &TMpMedia{
		MediaId:    mediaId,
		Type:       typestr,
		Synctime:   time.Now().Unix(),
		Createtime: createtime,
		Updatetime: updatetime,
		NewsList:   make([]*TMpNews, 0, 3),
	}
}
