package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"strconv"
	"strings"

	"csz.net/tgstate/conf"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ValidateBotConfig 验证 Bot 配置
func ValidateBotConfig() error {
	if conf.BotToken == "" {
		return fmt.Errorf("bot token 未配置")
	}

	if conf.ChannelName == "" {
		return fmt.Errorf("频道名称未配置")
	}

	// 测试 Bot Token 是否有效
	bot, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		return fmt.Errorf("bot token 无效: %v", err)
	}

	// 获取 Bot 信息
	me, err := bot.GetMe()
	if err != nil {
		return fmt.Errorf("无法获取 bot 信息: %v", err)
	}

	log.Printf("Bot 配置验证成功: @%s (%s)", me.UserName, me.FirstName)
	return nil
}

func TgFileData(fileName string, fileData io.Reader) tgbotapi.FileReader {
	return tgbotapi.FileReader{
		Name:   fileName,
		Reader: fileData,
	}
}

func UpDocument(fileData tgbotapi.FileReader) string {
	bot, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Printf("创建 Bot API 实例失败: %v", err)
		return ""
	}

	// 验证配置
	if conf.ChannelName == "" {
		log.Println("错误: 频道名称未配置")
		return ""
	}

	log.Printf("正在上传文件 '%s' 到频道 '%s'", fileData.Name, conf.ChannelName)

	// Upload the file to Telegram
	params := tgbotapi.Params{
		"chat_id": conf.ChannelName, // Replace with the chat ID where you want to send the file
	}
	files := []tgbotapi.RequestFile{
		{
			Name: "document",
			Data: fileData,
		},
	}
	response, err := bot.UploadFiles("sendDocument", params, files)
	if err != nil {
		log.Printf("上传文件到 Telegram 失败: %v", err)
		log.Printf("请检查: 1) Bot Token 是否正确 2) 频道名称 '%s' 是否正确 3) Bot 是否已添加到频道并有发送权限", conf.ChannelName)
		return ""
	}
	var msg tgbotapi.Message
	if err := json.Unmarshal([]byte(response.Result), &msg); err != nil {
		log.Printf("解析 Telegram 响应失败: %v", err)
		log.Printf("响应内容: %s", response.Result)
		return ""
	}

	var resp string
	switch {
	case msg.Document != nil:
		resp = msg.Document.FileID
		log.Printf("文档上传成功，FileID: %s", resp)
	case msg.Audio != nil:
		resp = msg.Audio.FileID
		log.Printf("音频上传成功，FileID: %s", resp)
	case msg.Video != nil:
		resp = msg.Video.FileID
		log.Printf("视频上传成功，FileID: %s", resp)
	case msg.Sticker != nil:
		resp = msg.Sticker.FileID
		log.Printf("贴纸上传成功，FileID: %s", resp)
	default:
		log.Println("警告: 无法识别上传的文件类型")
	}

	if resp == "" {
		log.Println("错误: 未能获取文件ID")
	}

	return resp
}

func GetDownloadUrl(fileID string) (string, bool) {
	bot, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Printf("创建 Bot API 实例失败: %v", err)
		return "", false
	}
	// 使用 getFile 方法获取文件信息
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		log.Println("获取文件失败【" + fileID + "】")
		log.Println(err)
		return "", false
	}
	log.Println("获取文件成功【" + fileID + "】")
	// 获取文件下载链接
	fileURL := file.Link(conf.BotToken)
	return fileURL, true
}
func BotDo() {
	bot, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Println(err)
		return
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updatesChan := bot.GetUpdatesChan(u)
	for update := range updatesChan {
		var msg *tgbotapi.Message
		if update.Message != nil {
			msg = update.Message
		}
		if update.ChannelPost != nil {
			msg = update.ChannelPost
		}
		if msg != nil && msg.Text == "get" && msg.ReplyToMessage != nil {
			var fileID string
			switch {
			case msg.ReplyToMessage.Document != nil && msg.ReplyToMessage.Document.FileID != "":
				fileID = msg.ReplyToMessage.Document.FileID
			case msg.ReplyToMessage.Video != nil && msg.ReplyToMessage.Video.FileID != "":
				fileID = msg.ReplyToMessage.Video.FileID
			case msg.ReplyToMessage.Sticker != nil && msg.ReplyToMessage.Sticker.FileID != "":
				fileID = msg.ReplyToMessage.Sticker.FileID
			}
			if fileID != "" {
				newMsg := tgbotapi.NewMessage(msg.Chat.ID, strings.TrimSuffix(conf.BaseUrl, "/")+"/d/"+fileID)
				newMsg.ReplyToMessageID = msg.MessageID
				if !strings.HasPrefix(conf.ChannelName, "@") {
					if man, err := strconv.Atoi(conf.ChannelName); err == nil && int(msg.Chat.ID) == man {
						bot.Send(newMsg)
					}
				} else {
					bot.Send(newMsg)
				}
			}
		}
	}
}

// GenerateShortCode 生成短链码
func GenerateShortCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}

	return string(result)
}
