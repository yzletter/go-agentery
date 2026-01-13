package llm

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/cloudwego/eino-ext/components/model/ark"
)

// 单例模式全局变量
var (
	arkModel *ark.ChatModel
	arkOnce  sync.Once
)

// CreateArkModel 创建火山方舟 Model
func CreateArkModel() *ark.ChatModel {
	arkOnce.Do(func() {
		ctx := context.Background()

		// 火山参数需要传指针
		//timeout := 30 * time.Second
		//maxTokens := 3000
		//temperature := float32(1.0)

		var err error
		arkModel, err = ark.NewChatModel(ctx, &ark.ChatModelConfig{
			APIKey: os.Getenv("ARK_KEY"),
			Model:  "doubao-seed-1-8-251228",
			//Timeout:     &timeout,
			//MaxTokens:   &maxTokens,   // 返回最大的 Token 数
			//Temperature: &temperature, // 温度，温度越高越发散 [0, 2]
		})
		if err != nil {
			log.Fatal(err)
		}
	})

	return arkModel
}
