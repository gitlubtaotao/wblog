package main

import (
	"flag"
	
	"github.com/cihub/seelog"
)

func main() {
	//获取用户输入的log 文件名
	logConfigPath := flag.String("L", "../conf/seelog.xml", "log config file path")
	flag.Parse()
	logger, err := seelog.LoggerFromConfigAsFile(*logConfigPath)
	if err != nil {
		seelog.Error(err)
	}
	seelog.ReplaceLogger(logger)
	defer seelog.Flush()
	seelog.Critical("测试文件")
	
}
