package conf

var BotToken string
var ChannelName string
var Pass string
var ApiPass string
var Mode string
var BaseUrl string
var AllowedExts string
var ProxyUrl string

type UploadResponse struct {
	Code         int    `json:"code"`
	Message      string `json:"message"`
	ImgUrl       string `json:"url"`
	ProxyUrl     string `json:"proxyUrl"`
	ShortUrl     string `json:"shortUrl"`
	ShortFileUrl string `json:"shortFileUrl"`
	Name         string `json:"name"`
	ChunkId      string `json:"chunkId,omitempty"`
}

type ResponseResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const FileRoute = "/d/"
