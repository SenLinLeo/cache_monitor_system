package main

import (
	"codis-monitor/conf"
	"codis-monitor/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func init() {
	// 设置日志格式为json格式
	log.SetFormatter(&log.JSONFormatter{})

	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	log.SetOutput(os.Stdout)

	// 设置日志级别为warn以上
	log.SetLevel(log.WarnLevel)
}

func main() {
	/** 初始化日志服务 **/
	fileObj, err := os.OpenFile("logs/logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Error("打开日志文件失败,error: ", err.Error())
		os.Exit(-1)
	}
	log.SetOutput(fileObj)

	var name = conf.Namespace
	namespace := &name
	checkSingleKeys := &name
	checkKeys := &name

	switch conf.Flag {
	case 1:
		log.Error("新建CodisExporter......")
		exporter, err := exporter.NewCodisExporter(
			exporter.RedisHost{Addrs: conf.Addrs, Passwords: conf.Passwords, Aliases: conf.Aliases},
			*namespace,
			*checkSingleKeys,
			*checkKeys,
		)
		if err != nil {
			log.Error("新建CodisExporter失败, error:", err)
			os.Exit(-1)
		}
		prometheus.MustRegister(exporter)
	case 2:
		log.Error("新建RedisExporter......")
		exporter, err := exporter.NewRedisExporter(
			exporter.RedisHost{Addrs: conf.Addrs, Passwords: conf.Passwords, Aliases: conf.Aliases},
			*namespace,
			*checkSingleKeys,
			*checkKeys,
		)
		if err != nil {
			log.Error("新建RedisExporter失败, error:", err)
			os.Exit(-1)
		}
		prometheus.MustRegister(exporter)
	default:
		log.Error("新建Exporter失败,参数有误Flag:", conf.Flag)
	}

	http.Handle(conf.MetricPath, promhttp.Handler())

	systemChan := make(chan int)

	/** 启动web服务，监听1010端口 **/
	go func() {
		log.Println("ListenAndServe at:172.30.60.194:1010")
		err := http.ListenAndServe("172.30.60.194:1010", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}

		systemChan <- 1
	}()

	for range systemChan {
		log.Error("缓存监控系统退出", )
	}

}
