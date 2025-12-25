package models

// Event 漫展事件结构
type Event struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Tag            string `json:"tag"`
	Location       string `json:"location"`
	Address        string `json:"address"`
	URL            string `json:"url"`
	Type           string `json:"type"`
	WannaGoCount   int    `json:"wannaGoCount"`
	CircleCount    int    `json:"circleCount"`
	DoujinshiCount int    `json:"doujinshiCount"`
	Time           string `json:"time"`
	AppLogoPicURL  string `json:"appLogoPicUrl"`
	LogoPicURL     string `json:"logoPicUrl"`
	Ended          string `json:"ended"`
	IsOnline       string `json:"isOnline"`
}

// APIResponse 统一响应结构
type APIResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// CPPAPIResponse CPP API 原始响应结构
type CPPAPIResponse struct {
	Result CPPResult `json:"result"`
}

type CPPResult struct {
	Total int            `json:"total"`
	List  []CPPEventItem `json:"list"`
}

// CPPEventItem CPP API 返回的事件项
type CPPEventItem struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Tag            string `json:"tag"`
	ProvName       string `json:"provName"`
	CityName       string `json:"cityName"`
	AreaName       string `json:"areaName"`
	EnterAddress   string `json:"enterAddress"`
	Type           string `json:"type"`
	WannaGoCount   int    `json:"wannaGoCount"`
	CircleCount    int    `json:"circleCount"`
	DoujinshiCount int    `json:"doujinshiCount"`
	EnterTime      int64  `json:"enterTime"`
	EndTime        int64  `json:"endTime"`
	StartTime      string `json:"startTime"`
	AppLogoPicURL  string `json:"appLogoPicUrl"`
	LogoPicURL     string `json:"logoPicUrl"`
	Enabled        int    `json:"enabled"`
	Ended          bool   `json:"ended"`
	IsOnline       int    `json:"isOnline"`
	EvmType        int    `json:"evmtype"`
}
