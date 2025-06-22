package main

import (
	"log"
	"net/http"
	"strconv"

	"back/db"

	"github.com/spf13/viper"
)

func main() {
	initConfig()

	if err := db.Init(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	srv := NewServer()

	port := viper.GetInt("port")
	addr := ":" + strconv.Itoa(port)
	log.Printf("服务器运行在 http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, srv))
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()
}
