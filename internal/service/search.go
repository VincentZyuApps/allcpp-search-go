package service

import (
	"cpp_search_go/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	baseURL      = "https://www.allcpp.cn/allcpp/event/eventMainListV2.do"
	cdnPrefix    = "https://imagecdn3.allcpp.cn/upload"
	pageSize     = 10
	requestDelay = 300 * time.Millisecond
	// 用于获取所有漫展的大页面大小
	allPageSize = 100
)

// SearchEvents 搜索漫展事件
func SearchEvents(keyword string) ([]models.Event, int, error) {
	// 获取第一页数据
	firstPageData, err := fetchPage(keyword, 1)
	if err != nil {
		return nil, 0, err
	}

	total := firstPageData.Result.Total
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	// 解析第一页事件
	allEvents := parseEvents(firstPageData.Result.List)

	// 获取剩余页面
	for page := 2; page <= totalPages; page++ {
		time.Sleep(requestDelay)

		pageData, err := fetchPage(keyword, page)
		if err != nil {
			continue // 跳过失败的页面
		}

		events := parseEvents(pageData.Result.List)
		allEvents = append(allEvents, events...)
	}

	// 按时间排序（从近到远）
	sort.Slice(allEvents, func(i, j int) bool {
		timeI := parseTimeForSort(allEvents[i].Time)
		timeJ := parseTimeForSort(allEvents[j].Time)
		return timeI.Before(timeJ)
	})

	return allEvents, total, nil
}

// SearchAllEvents 获取所有漫展（不需要关键词）
func SearchAllEvents() ([]models.Event, int, error) {
	// 用空关键词和大页面获取所有漫展
	firstPageData, err := fetchPageWithSize("", 1, allPageSize)
	if err != nil {
		return nil, 0, err
	}

	total := firstPageData.Result.Total
	totalPages := int(math.Ceil(float64(total) / float64(allPageSize)))

	// 解析第一页事件
	allEvents := parseEvents(firstPageData.Result.List)

	// 获取剩余页面
	for page := 2; page <= totalPages; page++ {
		time.Sleep(requestDelay)

		pageData, err := fetchPageWithSize("", page, allPageSize)
		if err != nil {
			continue // 跳过失败的页面
		}

		events := parseEvents(pageData.Result.List)
		allEvents = append(allEvents, events...)
	}

	// 按时间排序（从近到远），排除已结束的放后面
	sort.Slice(allEvents, func(i, j int) bool {
		// 已结束的排后面
		if allEvents[i].Ended == "已结束" && allEvents[j].Ended != "已结束" {
			return false
		}
		if allEvents[i].Ended != "已结束" && allEvents[j].Ended == "已结束" {
			return true
		}
		// 同类按时间排序
		timeI := parseTimeForSort(allEvents[i].Time)
		timeJ := parseTimeForSort(allEvents[j].Time)
		return timeI.Before(timeJ)
	})

	return allEvents, total, nil
}

// FetchRawData 获取原始数据（调试用）
func FetchRawData(keyword string) (*models.CPPAPIResponse, error) {
	return fetchPage(keyword, 1)
}

// fetchPage 获取指定页面的数据
func fetchPage(keyword string, page int) (*models.CPPAPIResponse, error) {
	return fetchPageWithSize(keyword, page, pageSize)
}

// fetchPageWithSize 获取指定页面的数据（可自定义页面大小）
func fetchPageWithSize(keyword string, page int, size int) (*models.CPPAPIResponse, error) {
	params := url.Values{}
	params.Set("time", "8")
	params.Set("sort", "1")
	params.Set("keyword", keyword)
	params.Set("pageNo", fmt.Sprintf("%d", page))
	params.Set("pageSize", fmt.Sprintf("%d", size))

	requestURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Host", "www.allcpp.cn")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("errorWrap", "json")
	req.Header.Set("Origin", "https://cp.allcpp.cn")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36 Edg/143.0.0.0")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://cp.allcpp.cn/")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result models.CPPAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v", err)
	}

	return &result, nil
}

// parseEvents 解析事件列表
func parseEvents(items []models.CPPEventItem) []models.Event {
	events := make([]models.Event, 0, len(items))

	for _, item := range items {
		event := parseEvent(item)
		events = append(events, event)
	}

	return events
}

// parseEvent 解析单个事件
func parseEvent(item models.CPPEventItem) models.Event {
	// 检查是否取消
	isCancelled := item.Enabled == 5

	// 处理名称
	name := item.Name
	if isCancelled && !strings.Contains(name, "(已取消)") {
		name += "(已取消)"
	}

	// 添加状态标签
	if !isCancelled {
		statusTag := getEventStatusTag(item)
		if statusTag != "" && !strings.Contains(name, statusTag) {
			name += statusTag
		}
	}

	// 处理线上/线下
	isOnline := "线下"
	if item.IsOnline == 1 {
		isOnline = "线上"
	}

	return models.Event{
		ID:             item.ID,
		Name:           name,
		Tag:            item.Tag,
		Location:       parseLocation(item),
		Address:        item.EnterAddress,
		URL:            fmt.Sprintf("https://www.allcpp.cn/allcpp/event/event.do?event=%d", item.ID),
		Type:           parseType(item),
		WannaGoCount:   item.WannaGoCount,
		CircleCount:    item.CircleCount,
		DoujinshiCount: item.DoujinshiCount,
		Time:           parseTime(item),
		AppLogoPicURL:  parseImageURL(item.AppLogoPicURL, item.LogoPicURL),
		LogoPicURL:     parseImageURLSingle(item.LogoPicURL),
		Ended:          parseEnded(item),
		IsOnline:       isOnline,
	}
}

// getEventStatusTag 获取展会状态标签
func getEventStatusTag(item models.CPPEventItem) string {
	if item.EnterTime == 0 {
		return ""
	}

	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	enterTime := time.UnixMilli(item.EnterTime).In(loc)
	enterDate := time.Date(enterTime.Year(), enterTime.Month(), enterTime.Day(), 0, 0, 0, 0, loc)

	endTime := enterTime
	if item.EndTime > 0 {
		endTime = time.UnixMilli(item.EndTime).In(loc)
	}
	endDate := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 0, loc)

	// 如果展会已结束
	if now.After(endDate) {
		return ""
	}

	// 如果展会正在进行中
	endDateStart := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, loc)
	if !today.Before(enterDate) && !today.After(endDateStart) {
		return "(进行中)"
	}

	// 如果展会还未开始，计算剩余天数
	if now.Before(enterTime) {
		daysLeft := int(math.Ceil(enterDate.Sub(today).Hours() / 24))
		if daysLeft > 0 {
			return fmt.Sprintf("(还有%d天开始)", daysLeft)
		}
	}

	return ""
}

// parseLocation 解析地点
func parseLocation(item models.CPPEventItem) string {
	parts := []string{}
	if item.ProvName != "" {
		parts = append(parts, item.ProvName)
	}
	if item.CityName != "" {
		parts = append(parts, item.CityName)
	}
	if item.AreaName != "" {
		parts = append(parts, item.AreaName)
	}
	return strings.Join(parts, " ")
}

// parseType 解析类型
func parseType(item models.CPPEventItem) string {
	if item.Type != "" {
		return item.Type
	}

	typeMap := map[int]string{
		0: "综合展",
		1: "ONLY",
		2: "茶会",
		3: "漫展",
	}

	if t, ok := typeMap[item.EvmType]; ok {
		return t
	}

	// 从标签猜测类型
	tag := strings.ToUpper(item.Tag)
	if strings.Contains(tag, "ONLY") {
		return "ONLY"
	}
	if strings.Contains(tag, "茶会") || strings.Contains(tag, "茶话会") {
		return "茶会"
	}

	return "综合展"
}

// parseTime 解析时间
func parseTime(item models.CPPEventItem) string {
	if item.EnterTime > 0 {
		t := time.UnixMilli(item.EnterTime)
		return t.Format("2006-01-02 15:04:05")
	}
	if item.StartTime != "" {
		return item.StartTime
	}
	return ""
}

// parseImageURL 解析图片URL（appLogoPicUrl）
func parseImageURL(appLogoURL, logoURL string) string {
	url := appLogoURL

	if url != "" && !strings.HasPrefix(url, "http") {
		url = cdnPrefix + url
	}

	// 如果为空，尝试从logoURL生成
	if url == "" && logoURL != "" {
		re := regexp.MustCompile(`\?.*$`)
		url = re.ReplaceAllString(logoURL, "")
		if !strings.HasPrefix(url, "http") {
			url = cdnPrefix + url
		}
	}

	return url
}

// parseImageURLSingle 解析单个图片URL
func parseImageURLSingle(url string) string {
	if url != "" && !strings.HasPrefix(url, "http") {
		return cdnPrefix + url
	}
	return url
}

// parseEnded 判断是否结束
func parseEnded(item models.CPPEventItem) string {
	// 优先判断 enabled 状态
	switch item.Enabled {
	case 1:
		return "已结束"
	case 2:
		return "筹备中"
	case 5:
		return "已取消"
	}

	// 根据结束时间判断
	if item.EndTime > 0 {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		endTime := time.UnixMilli(item.EndTime).In(loc)
		endDate := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 0, loc)

		if time.Now().In(loc).After(endDate) {
			return "已结束"
		}
	}

	// 判断 ended 字段
	if item.Ended {
		return "已结束"
	}

	return "未结束"
}

// parseTimeForSort 解析时间用于排序
func parseTimeForSort(timeStr string) time.Time {
	if timeStr == "" {
		return time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	}

	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	}
	return t
}
