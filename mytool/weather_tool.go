package mytool

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type WeatherParams struct {
	CityCode string `json:"city_code" jsonschema:"required,description=城市的编码"`
	// description中如果包含,还需要转义（,在tag中是分隔的意思），此时建议用jsonschema_description
	Day int `json:"day" jsonschema:"required,enum=0,enum=1,enum=2,enum=3,enum=4" jsonschema_description:"获取现在和未来的气象数据。0表示当前实时的气象数据，1表示今天的气象数据，2表示明天的气象数据，3表示后天的气象数据，4表示大后天的气象数据，最大只能到4"`
}

type WeatherCast struct {
	Date string // 日期
	Week string // 周几
	// 白天
	DayWeather       string // 天气
	DayTemperature   string `json:"daytemp"`  // 温度
	DayWindDirection string `json:"daywind"`  // 风向
	DayWindPowerr    string `json:"daypower"` // 风力
	// 晚上
	NightWeather       string
	NightTemperature   string `json:"nighttemp"`
	NightWindDirection string `json:"nightwind"`
	NightWindPowerr    string `json:"nightpower"`
}

type LiveWeather struct {
	Weather       string
	Temperature   string
	WindDirection string
	WindPower     string
	Humidity      string
}

type WeatherInfo struct {
	Status    string        `json:"status"`
	Lives     []LiveWeather `json:"lives"` // 实况天气
	Forecasts []struct {    // 预报天气
		Casts []WeatherCast // 预报数据 list 结构，元素 cast,按顺序为当天、第二天、第三天、第四天的预报数据
	}
}

func WeatherForecast(GaodeCityCode string, live bool) *WeatherInfo {
	url := "https://restapi.amap.com/v3/weather/weatherInfo?" + "key=" + GaoDeKey + "&city=" + GaodeCityCode
	// 查询实时天气还是所有天气
	if !live {
		url += "&extensions=all"
	}

	resp, err := HttpGet(url)
	if err != nil {
		log.Println("获取天气失败")
		return nil
	}

	var weather WeatherInfo
	err = sonic.Unmarshal(resp, &weather)

	if weather.Status != "1" {
		log.Println("天气反序列化失败")
		return nil
	}

	return &weather
}

func CreateWeatherTool() tool.InvokableTool {
	weatherTool, err := utils.InferTool(
		"weather_tool",
		"获取气象数据，包含天气、温度、风力、风向等，既可以获得当前的实时气象数据，也可以获取未来3天的气象数据",
		func(ctx context.Context, params WeatherParams) (output string, err error) {
			if params.Day == 0 {
				// 当前实时天气
				weather := WeatherForecast(params.CityCode, true)
				if weather == nil {
					return "", errors.New("获取天气失败")
				}

				return fmt.Sprintf("现在的天气情况是: %s，温度 %s℃，湿度 %s%%，风力 %s级，风向 %s", weather.Lives[0].Weather, weather.Lives[0].Temperature, weather.Lives[0].Humidity, weather.Lives[0].WindPower, weather.Lives[0].WindDirection), nil
			} else if params.Day >= 1 && params.Day <= 4 {
				weathers := WeatherForecast(params.CityCode, false)
				if weathers == nil {
					return "", errors.New("获取天气失败")
				}
				//fmt.Println(weathers)
				dayWeather := weathers.Forecasts[0].Casts[params.Day-1]
				return fmt.Sprintf("%s 白天的天气情况是: %s，温度 %s℃，风力 %s级，风向 %s，晚上的天气情况是: %s，温度 %s℃，风力 %s级，风向 %s", dayWeather.Date,
					dayWeather.DayWeather, dayWeather.DayTemperature, dayWeather.DayWindPowerr, dayWeather.DayWindDirection,
					dayWeather.NightWeather, dayWeather.NightTemperature, dayWeather.NightWindPowerr, dayWeather.NightWindDirection), nil
			} else {
				return "", errors.New("day 参数不合法")
			}
		},
	)
	if err != nil {
		log.Println("创建天气工具失败")
		return nil
	}

	return weatherTool
}
