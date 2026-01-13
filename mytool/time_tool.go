package mytool

import (
	"context"
	"log"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type TimeRequest struct {
	//TimeZone string `json:"time_zone" jsonschema:"required,description=当地时区"`	// 必传加 required
	TimeZone string `json:"time_zone" jsonschema:"description=当地时区"`
}

// GetTime 获取当前时间的函数
func GetTime(ctx context.Context, TimeZone *TimeRequest) (string, error) {
	// 没传时区
	if len(TimeZone.TimeZone) == 0 {
		TimeZone.TimeZone = "Asia/Shanghai"
	}

	// 获取时区
	loc, err := time.LoadLocation(TimeZone.TimeZone)
	if err != nil {
		return "", err
	}

	// 获取当前时间
	now := time.Now().In(loc).Format("2006-01-02 15:04:05")

	return now, nil
}

// CreateTimeTool 把 GetTime 转为 Tool
func CreateTimeTool() tool.InvokableTool {
	tool, err := utils.InferTool(
		"time_tool", // 工具名
		"获取用户当前的时间", // 工具描述
		GetTime,
	)

	if err != nil {
		log.Fatal(err)
	}

	return tool
}
