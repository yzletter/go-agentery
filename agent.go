package main

import (
	"context"
	"sync"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	llm "github.com/yzletter/go-agentery/model"
	"github.com/yzletter/go-agentery/mytool"
)

var (
	runner *adk.Runner
	once   sync.Once
)

func GetRunner() *adk.Runner {
	once.Do(func() {
		ctx := context.Background()

		// 创建 Model
		model := llm.CreateArkModel()

		// 创建工具集
		tools := []tool.BaseTool{
			mytool.CreateTimeTool(),
			mytool.CreateLocationTool(),
			mytool.CreateWeatherTool(),
		}

		// 创建 Agent
		agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
			Name:        "my_agent",
			Description: "我的万能AI助手",
			Instruction: "", // SystemMessage，支持FString渲染
			Model:       model,
			ToolsConfig: adk.ToolsConfig{
				ToolsNodeConfig: compose.ToolsNodeConfig{
					Tools: tools,
				},
			},
		})
		if err != nil {
			panic(err)
		}

		// 创建 Runner
		runner = adk.NewRunner(ctx, adk.RunnerConfig{
			Agent:           agent, // 指定哪个 Agent
			EnableStreaming: true,  // 允许流式
			CheckPointStore: nil,   // 中断恢复
		})
	})

	return runner
}
