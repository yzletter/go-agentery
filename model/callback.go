package llm

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	cbutils "github.com/cloudwego/eino/utils/callbacks"
)

/*
callbacks.AppendGlobalHandlers()注入全局Handler，提供切面能力，以无侵入的方式注入logging, tracing, metrics等功能。
callbacks.AppendGlobalHandlers(new(LoggerCallbacks)) // 调试的时候可以打开，查看每一个Node的输入和输出
callbacks.AppendGlobalHandlers(GetChatModelInputCallback()) // 调试的时候可以打开，查看ChatModel的输入
*/
type LoggerCallbacks struct{}

func (l *LoggerCallbacks) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	log.Printf("[INPUT] name: %v, type: %v, component: %v, input: %v", info.Name, info.Type, info.Component, input)
	return ctx
}

func (l *LoggerCallbacks) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	// output里还包含 reasoning content 推理过程，可以看一下大模型错在哪儿了
	log.Printf("[OUTPUT] name: %v, type: %v, component: %v, output: %v", info.Name, info.Type, info.Component, output)
	return ctx
}

func (l *LoggerCallbacks) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	log.Printf("[ERROR] name: %v, type: %v, component: %v, error: %v", info.Name, info.Type, info.Component, err)
	return ctx
}

func (l *LoggerCallbacks) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo, input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	return ctx
}

func (l *LoggerCallbacks) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {
	return ctx
}

/**
如果一个 Handler，不想关注所有的 5 个触发时机，只想关注一部分，比如只关注 OnStart，建议使用 callbacks.NewHandlerBuilder().OnStartFn(...).Build()。
如果不想关注所有的组件类型，只想关注特定组件，比如 ChatModel，建议使用 github.com/cloudwego/eino/utils/callbacks.NewHandlerHelper().ChatModel(...).Handler()，可以只接收 ChatModel 的回调并拿到一个具体类型的 CallbackInput/CallbackOutput。
*/

func GetStartCallback() callbacks.Handler {
	return callbacks.NewHandlerBuilder().
		OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
			log.Printf("[INPUT] name: %v, type: %v, component: %v, input: %v", info.Name, info.Type, info.Component, input)
			return ctx
		}).Build() // Builder设计模式
}

func GetEndCallback() callbacks.Handler {
	return callbacks.NewHandlerBuilder().OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
		// output里还包含 reasoning content 推理过程，可以看一下大模型错在哪儿了
		log.Printf("[OUTPUT] name: %v, type: %v, component: %v, output: %v", info.Name, info.Type, info.Component, output)
		return ctx
	}).Build() // Builder设计模式
}

/*
给特定节点加Callback
*/
func GetChatModelInputCallback() callbacks.Handler {
	// 只给ChatModel注入CallbackHandler
	return cbutils.NewHandlerHelper().ChatModel(&cbutils.ModelCallbackHandler{
		// 这里只定义了OnStart，当然还可以指定OnEnd、OnError等
		OnStart: func(ctx context.Context, info *callbacks.RunInfo, input *model.CallbackInput) context.Context {
			// 打印ChatModel的输入
			fmt.Printf("\n[ChatModel Input] component: %s\n", info.Name)
			for i, msg := range input.Messages {
				fmt.Printf("  Message %d [%s]: %s\n", i+1, msg.Role, msg.Content)
				if len(msg.ToolCalls) > 0 {
					fmt.Printf("    Tool Calls: %d\n", len(msg.ToolCalls))
					for j, tc := range msg.ToolCalls {
						fmt.Printf("      %d. %s: %s\n", j+1, tc.Function.Name, tc.Function.Arguments)
					}
				}
			}
			return ctx
		},
	}).Handler()
}

func GetToolInputCallback() callbacks.Handler {
	// 只给Tool注入CallbackHandler
	return cbutils.NewHandlerHelper().Tool(&cbutils.ToolCallbackHandler{
		// 这里只定义了OnStart，当然还可以指定OnEnd、OnError等
		OnStart: func(ctx context.Context, info *callbacks.RunInfo, input *tool.CallbackInput) context.Context {
			// 打印ChatModel的输入
			fmt.Printf("\n[Tool Input] component: %s, args: %s\n", info.Name, input.ArgumentsInJSON)
			return ctx
		},
	}).Handler()
}
