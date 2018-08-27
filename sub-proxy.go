package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	consts "github.com/Basic-Components/sub-proxy/consts"
	loadconfig "github.com/Basic-Components/sub-proxy/loadconfig"
	proxy "github.com/Basic-Components/sub-proxy/proxy"

	log "github.com/sirupsen/logrus"
)

func closeHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		fmt.Println("Device shutdown!")
		os.Exit(0)
	}()

}
func makeConfig() loadconfig.Config {
	var (
		help        bool
		debug       bool
		conflate    bool
		sndHWM      int
		rcvHWM      int
		groupName   string
		frontendURL string
		backendURL  string
		logFormat   string
		logOutput   string
		configPath  string
	)
	// 解析命令行参数
	flag.BoolVar(&help, "help", false, "帮助命令")
	flag.BoolVar(&debug, "debug", false, "是否使用debug模式启动")
	flag.BoolVar(&conflate, "conflate", false, "是否只保留最后到来的消息")
	flag.IntVar(&sndHWM, "send_hwm", -1, "发送端的缓存长度")
	flag.IntVar(&rcvHWM, "receive_hwm", -1, "接收端的缓存长度")
	flag.StringVar(&groupName, "group_name", "", "使用该中间人的堆名")
	flag.StringVar(&frontendURL, "frontend_url", "", "前端连接的地址")
	flag.StringVar(&backendURL, "backend_url", "", "后端绑定的地址")
	flag.StringVar(&logFormat, "log_format", "", "设定log的形式")
	flag.StringVar(&logOutput, "log_output", "", "设定log输出的流位置")
	flag.StringVar(&configPath, "config_path", "", "设定log输出的流位置")
	flag.Parse()

	//构造config
	if help {
		fmt.Println("[" + consts.TYPE + "]:" + consts.NAME)
		fmt.Println("version:" + consts.VERSION)
		fmt.Println(consts.DESCRIPTION)
		flag.PrintDefaults()
		os.Exit(1)
	}
	var config = loadconfig.LoadConfig(configPath)
	config.Debug = debug
	config.Conflate = conflate
	if sndHWM >= 0 {
		config.SNDHWM = sndHWM
	}
	if rcvHWM >= 0 {
		config.RCVHWM = rcvHWM
	}
	if groupName != "" {
		config.GroupName = groupName
	}
	if frontendURL != "" {
		config.FrontendURL = frontendURL
	}
	if backendURL != "" {
		config.BackendURL = backendURL
	}
	if logFormat != "" {
		config.LogFormat = logFormat
	}
	if logOutput != "" {
		config.LogOutput = logOutput
	}
	return config
}

func main() {
	var config = makeConfig()
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
	switch config.LogFormat {
	case "json":
		{
			log.SetFormatter(&log.JSONFormatter{})
		}
	case "text":
		{
			log.SetFormatter(&log.TextFormatter{})
		}
	default:
		{
			log.SetFormatter(&log.JSONFormatter{})
		}
	}
	if config.LogOutput == "" {
		log.SetOutput(os.Stdout)
	} else {
		logFile, _ := os.OpenFile(config.LogOutput, os.O_CREATE|os.O_WRONLY, 0666)
		defer logFile.Close()
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
	}
	log.WithFields(log.Fields{
		consts.TYPE: consts.NAME,
	}).Info("Proxy for Stack [" + config.GroupName + "] start !")
	log.WithFields(log.Fields{
		consts.TYPE: consts.NAME,
	}).Info("Publisher should connect to url " + config.FrontendURL)
	log.WithFields(log.Fields{
		consts.TYPE: consts.NAME,
	}).Info("Subscriber should connect to url " + config.BackendURL)
	closeHandler()
	proxy.Run(config)
}
