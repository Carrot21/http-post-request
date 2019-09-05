package main

import (
	"flag"
	"http-post-request/config"
	"http-post-request/http"
	"log"
)

var (
	confFile = flag.String("confFile", "config/http.yml", "Configuration file.")

	apiIndex   = flag.Int("api_index", 1, "请输入接口的索引,以1为开始")
	userID     = flag.String("userid", "", "请输入UserID")
	chnnID     = flag.String("chnnid", "", "请输入ChnnID")
	indexID    = flag.String("index", "", "请输入索引ID")
	token      = flag.String("token", "", "请输入Token")
	goroutines = flag.Int("num", 1, "请输入并发数量")
)

func main() {
	flag.Parse()

	config.LoadConfig(*confFile)

	parseCommandArgs()
}

func parseCommandArgs() {
	if *apiIndex < 1 || *apiIndex > 7 {
		log.Println("输入API序号有误！请输入[1-7]闭区间的数字！")
		return
	}

	requestArgs := http.RequestArgs{
		UserID:  *userID,
		ChnnID:  *chnnID,
		IndexID: *indexID,
		Token:   *token,
	}

	http.HandleRequest(*apiIndex, requestArgs, *goroutines)
}
