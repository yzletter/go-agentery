package mytool

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

var (
	GaoDeKey = os.Getenv("GAODE_KEY")
)

func HttpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http Get %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http response StatusCode %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error: %w", err)
	}
	return data, nil
}

// GetOutboundIP 获取用户当前 IP
func GetOutboundIP() string {
	data, err := HttpGet("https://httpbin.org/ip")
	if err != nil {
		log.Printf("获取对外 IP 失败: %s", err)
		return ""
	}

	mp := make(map[string]string, 1)
	sonic.Unmarshal(data, &mp)
	return mp["origin"]
}

type Location struct {
	Status    string
	Province  string
	City      string
	Rectangle string
	Longitude float64
	Latitude  float64
	AdCode    string
}

func GetMyLocation() (*Location, error) {
	ip := GetOutboundIP()
	url := "https://restapi.amap.com/v3/ip?key=" + GaoDeKey + "&ip=" + ip
	data, err := HttpGet(url)
	if err != nil {
		log.Printf("获取对外 IP 失败: %s", err)
		return nil, err
	}

	var location Location
	err = sonic.Unmarshal(data, &location)
	if err != nil {
		log.Printf("JSON 反序列化失败: %s", err)
		return nil, err
	} else if location.Status != "1" {
		log.Printf("GetMyLocation Failed: %s", err)
		return nil, errors.New("GetMyLocation Failed")
	}

	return &location, nil
}

func CreateLocationTool() tool.InvokableTool {
	locationTool, err := utils.InferTool(
		"location_tool",
		"获取当前地理位置,包括省、城市（含城市名称和城市编码）",
		func(ctx context.Context, input struct{}) (string, error) { // 通过匿名函数满足 Tool 所需要的标签
			location, _ := GetMyLocation()
			if location == nil {
				return "", errors.New("获取当前地理位置失败")
			}
			return fmt.Sprintf("当前城市：%s, 城市编码：%s, 属于省份：%s", location.City, location.AdCode, location.Province), nil
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	return locationTool
}
