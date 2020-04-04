package main

import (
	"flag"
	"github.com/cihub/seelog"
	"github.com/gitlubtaotao/wblog/system"
)

func main() {
	//配置环境变量
	configEnv := flag.String("env", "development", "set env development or production")
	flag.Parse()
	system.SetSeelogPath("../conf/seelog.xml")
	
	if err := system.LoadConfiguration(*configFilePath); err != nil {
		seelog.Critical("err parsing config log file", err)
		return
	}
	defer seelog.Flush()
	
}
