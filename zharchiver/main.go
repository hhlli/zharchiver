package main

import (
	"fmt"
	"log"
	"net/http"

	"zharchiver/handlers"
	"zharchiver/models"
	"zharchiver/services"
)

func main() {
	db, err := models.InitDB("db/zharchiver.db")
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer db.Close()

	// 启动后台服务
	services.StartTelegramAutoBackup(db)
	services.StartTelegramBotListener(db)

	// 注册路由
	mux := http.NewServeMux()
	handler := handlers.RegisterRoutes(mux, db)

	port := 8080
	addr := fmt.Sprintf(":%d", port)
	log.Printf("服务已启动: http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}