package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"csz.net/tgstate/conf"
	"csz.net/tgstate/control"
	"csz.net/tgstate/utils"
)

var webPort string
var OptApi = true

func main() {
	//判断是否设置参数
	if conf.BotToken == "" || conf.ChannelName == "" {
		fmt.Println("请先设置Bot Token和对象")
		return
	}

	// 验证 Bot 配置
	if err := utils.ValidateBotConfig(); err != nil {
		log.Printf("Bot 配置验证失败: %v", err)
		fmt.Println("\n请检查以下配置:")
		fmt.Printf("1. Bot Token: %s\n", maskToken(conf.BotToken))
		fmt.Printf("2. 频道名称: %s\n", conf.ChannelName)
		fmt.Println("\n常见问题解决方案:")
		fmt.Println("- 确保 Bot Token 正确且有效")
		fmt.Println("- 确保频道名称格式正确 (如: @channelname 或 -1001234567890)")
		fmt.Println("- 确保 Bot 已被添加到频道并具有管理员权限")
		fmt.Println("- 确保 Bot 有发送消息的权限")
		return
	}

	go utils.BotDo()
	web()
}

// maskToken 遮蔽 Token 的敏感部分
func maskToken(token string) string {
	if len(token) < 10 {
		return "***"
	}
	return token[:10] + "***" + token[len(token)-4:]
}

func web() {
	http.HandleFunc(conf.FileRoute, control.D)
	http.HandleFunc("/s/", control.S) // 短链路由
	if OptApi {
		if conf.Pass != "" && conf.Pass != "none" {
			http.HandleFunc("/pwd", control.Pwd)
		}
		http.HandleFunc("/api", control.Middleware(control.UploadAPI))
		http.HandleFunc("/files", control.Middleware(control.FilesAPI))
		http.HandleFunc("/shortlinks", control.Middleware(control.ShortLinksAPI))
		http.HandleFunc("/", control.Middleware(control.Index))
	}

	if listener, err := net.Listen("tcp", ":"+webPort); err != nil {
		log.Fatalf("端口 %s 已被占用\n", webPort)
	} else {
		defer listener.Close()
		log.Printf("Http server start at %s\n", webPort)
		if err := http.Serve(listener, nil); err != nil {
			log.Fatal(err)
		}
	}
}

func init() {
	_ = godotenv.Load()

	flag.StringVar(&webPort, "port", "8088", "Web Port")
	flag.StringVar(&conf.BotToken, "token", os.Getenv("token"), "Bot Token")
	flag.StringVar(&conf.ChannelName, "target", os.Getenv("target"), "Channel Name or ID")
	flag.StringVar(&conf.Pass, "pass", os.Getenv("pass"), "Visit Password")
	flag.StringVar(&conf.ApiPass, "apiPass", os.Getenv("apiPass"), "API Visit Password")
	flag.StringVar(&conf.Mode, "mode", os.Getenv("mode"), "Run mode")
	flag.StringVar(&conf.BaseUrl, "url", os.Getenv("url"), "Base Url")
	flag.StringVar(&conf.AllowedExts, "exts", os.Getenv("exts"), "Allowed Exts")
	flag.StringVar(&conf.ProxyUrl, "proxyUrl", os.Getenv("proxyUrl"), "proxy url")
	flag.Parse()
	if conf.Mode == "m" {
		OptApi = false
	}
	if conf.Mode != "p" && conf.Mode != "m" {
		conf.Mode = "p"
	}
	_, err := control.InitDB()
	if err != nil {
		log.Fatal(err)
	}

}
