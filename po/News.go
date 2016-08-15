package po

type TMpNews struct {
	Id               int64
	MediaId          string
	Title            string
	Author           string
	Digest           string
	ThumbUrl         string
	ThumbMediaId     string
	Url              string
	ContentSourceUrl string
}

func NewTMpNews(mediaId, title, author, digest, thumbUrl, thumbMediaId, url, contentSourceUrl string) *TMpNews {
	return &TMpNews{
		MediaId:          mediaId,
		Title:            title,
		Author:           author,
		Digest:           digest,
		ThumbUrl:         thumbUrl,
		ThumbMediaId:     thumbMediaId,
		Url:              url,
		ContentSourceUrl: contentSourceUrl,
	}
}
