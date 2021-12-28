package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

// 全局日志配置初始化
// 1.日志同时输出到到文件和标准输出；
// 2.保留最近 7 天日志，一天一个日志文件；
func init() {

	/*
			日志轮转相关函数
				WithLinkName 		为最新的日志建立软连接
				WithRotationTime	设置日志分割的时间，隔多久分割一次
				WithMaxAge 和 WithRotationCount 二者只能设置一个

		  		WithMaxAge	 		设置文件清理前的最长保存时间
		  		WithRotationCount 	设置文件清理前最多保存的个数
	*/

	logfile := "./third_party/logger/logfile/"
	fsWriter, err := rotatelogs.New(
		logfile+"%Y-%m-%d.log",
		rotatelogs.WithMaxAge(time.Duration(168)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		panic(err)
	}

	multiWriter := io.MultiWriter(fsWriter, os.Stdout)
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(multiWriter)
	log.SetLevel(log.InfoLevel)

}
