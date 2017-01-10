package po

/**
图文素材信息
*/
type TMpNews struct {
	Id               int64 `gorm:"primary_key"`
	MediaId          string
	Title            string
	Author           string
	Digest           string
	ThumbUrl         string
	ThumbMediaId     string
	Url              string
	ContentSourceUrl string
	UpdateTime       int64
}

func NewTMpNews(mediaId, title, author, digest, thumbUrl, thumbMediaId, url, contentSourceUrl string, updateTime int64) *TMpNews {
	return &TMpNews{
		MediaId:          mediaId,
		Title:            title,
		Author:           author,
		Digest:           digest,
		ThumbUrl:         thumbUrl,
		ThumbMediaId:     thumbMediaId,
		Url:              url,
		ContentSourceUrl: contentSourceUrl,
		UpdateTime:       updateTime,
	}
}
