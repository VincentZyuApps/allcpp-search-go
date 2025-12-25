package main

import (
	"cpp_search_go/internal/api"
	"cpp_search_go/internal/config"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 注册路由
	api.RegisterRoutes(r)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("🚀 CPP Search API 启动于 http://%s", addr)
	log.Printf("📖 使用方法: GET /search?msg=关键词")

	if err := r.Run(addr); err != nil {
		log.Fatalf("❌ 服务启动失败: %v", err)
	}
}
