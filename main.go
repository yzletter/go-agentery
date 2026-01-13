package main

import (
	"context"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

// ParseUrlParams 解析路由
func ParseUrlParams(rawQuery string) map[string]string {
	params := make(map[string]string, 10)
	args := strings.Split(rawQuery, "&")
	for _, arg := range args {
		arr := strings.Split(arg, "=")
		if len(arr) != 2 {
			continue
		}
		key, _ := url.QueryUnescape(arr[0])
		value, _ := url.QueryUnescape(arr[1])

		params[key] = value
	}
	return params
}

func Chat(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Fatal(ok)
	}

	w.Header().Add("Content-Type", "text/event-stream; charset=utf-8") // 标识响应为事件流。charset=utf-8是为了解决中文乱码
	w.Header().Add("Cache-Control", "no-cache")                        // 防止浏览器缓存响应，确保实时性
	w.Header().Add("Connection", "keep-alive")                         // 保持连接开放，支持持续流式传输

	// 解析路由
	params := ParseUrlParams(r.URL.RawQuery)
	message := params["msg"]

	// 调 Agent
	ctx := context.Background()
	messages := make([]adk.Message, 0, 10)
	messages = append(messages, &schema.Message{Role: schema.User, Content: message})
	runner := GetRunner()
	iter := runner.Run(ctx, messages)

	// http://127.0.0.1:5678/chat?msg=golang是什么&id=121&dd=111
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}

		if event.Err != nil {
			log.Printf("Read LLM Failed, %s", event.Err)
			break
		}

		//msg, err := event.Output.MessageOutput.GetMessage()
		//if err != nil {
		//	log.Printf("Read LLM Failed, %s", event.Err)
		//} else {
		//	// 判断是不是大模型输出的
		//	if msg != nil && msg.Role == schema.System {
		//		fmt.Println(msg)
		//		w.Write([]byte(msg.Content))
		//	}
		//}
		s := event.Output.MessageOutput.MessageStream
		for {
			msg, err := s.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Printf("Stream Recv Failed, %s", err)
			}

			if msg != nil {
				// 拿不到角色
				w.Write([]byte(msg.Content))
				flusher.Flush()
			}
		}
	}

}

func main() {
	mux := http.NewServeMux()

	// 图标
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("./favicon.png")
		if err != nil {
			return
		}
		defer file.Close()
		io.Copy(w, file)
	})

	// 页面
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./chat.html")
		if err != nil {
			log.Println("Template Create Failed")
			return
		}
		err = tmpl.Execute(w, map[string]string{"url": "http://127.0.0.1:5678/chat"})
		if err != nil {
			log.Println("Template Excute Failed")
			return
		}
	})

	// Chat
	mux.HandleFunc("GET /chat", Chat)

	if err := http.ListenAndServe("127.0.0.1:5678", mux); err != nil {
		panic(err)
	}
}
