package api

import (
	"cpp_search_go/internal/models"
	"cpp_search_go/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.Engine) {
	r.GET("/search", handleSearch)
	r.GET("/search_all", handleSearchAll)
	r.GET("/", handleIndex)
}

// handleIndex 首页
func handleIndex(c *gin.Context) {
	c.JSON(http.StatusOK, models.APIResponse{
		Code: 200,
		Msg:  "CPP Search API - Go Version",
		Data: map[string]string{
			"usage":     "GET /search?msg=关键词",
			"usage_all": "GET /search_all 获取所有漫展",
			"debug":     "GET /search?msg=关键词&debug=raw",
			"version":   "1.0.0",
		},
	})
}

// handleSearchAll 获取所有漫展
func handleSearchAll(c *gin.Context) {
	// 用空关键词搜索，CPP API 会返回所有漫展
	events, total, err := service.SearchAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code: 500,
			Msg:  "获取数据失败: " + err.Error(),
			Data: []interface{}{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "所有漫展",
		"total": total,
		"data":  events,
	})
}

// handleSearch 搜索漫展
func handleSearch(c *gin.Context) {
	msg := strings.TrimSpace(c.Query("msg"))
	debug := c.Query("debug")

	// 检查关键词
	if msg == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code: 400,
			Msg:  "请提供搜索关键词",
			Data: []interface{}{},
		})
		return
	}

	// 调试模式：返回原始数据
	if debug == "raw" || debug == "response" {
		rawData, err := service.FetchRawData(msg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Code: 500,
				Msg:  "获取数据失败: " + err.Error(),
				Data: []interface{}{},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":     200,
			"msg":      "原始数据（debug模式）",
			"total":    rawData.Result.Total,
			"raw_data": rawData.Result.List,
		})
		return
	}

	// 正常搜索
	events, _, err := service.SearchEvents(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code: 500,
			Msg:  "获取数据失败: " + err.Error(),
			Data: []interface{}{},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code: 200,
		Msg:  msg,
		Data: events,
	})
}
