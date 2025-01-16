package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bestk/native-printer/config"
	"github.com/bestk/native-printer/printer"
	"github.com/gorilla/websocket"
	"gopkg.in/natefinch/lumberjack.v2"
)

var upgrader = websocket.Upgrader{}

// PrintMessage 定义打印消息结构
type PrintMessage struct {
	Action  string `json:"action"`
	Printer string `json:"printer"`
	FileUrl string `json:"fileUrl"`
}

type PrintResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 初始化日志配置
func initLogger(config *config.Config) {
	// 创建多重写入器，同时写入文件和控制台
	multiWriter := io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   filepath.Join(config.Log.Folder, "app.log"), // 日志文件路径
		MaxSize:    10,                                          // 每个日志文件最大尺寸，单位 MB
		MaxBackups: 30,                                          // 保留的旧日志文件最大数量
		MaxAge:     7,                                           // 保留的旧日志文件最大天数
		Compress:   true,                                        // 是否压缩旧日志文件
	})

	// 设置日志输出到多重写入器
	log.SetOutput(multiWriter)

	// 设置日志格式
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// 处理打印请求
func handlePrintRequest(message *PrintMessage) error {
	// 验证必要参数
	if message.Action != "printPDF" {
		return fmt.Errorf("不支持的操作: %s", message.Action)
	}

	// 验证 FileUrl
	if message.FileUrl == "" {
		return fmt.Errorf("文件 URL 不能为空")
	}

	// 验证 URL 格式
	if _, err := url.Parse(message.FileUrl); err != nil {
		return fmt.Errorf("无效的文件 URL: %v", err)
	}

	// 创建临时文件夹用于存储下载的PDF
	tempDir := "./temp"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("创建临时目录失败: %v", err)
	}

	// 下载PDF文件
	tempFile := filepath.Join(tempDir, "temp.pdf")
	if err := downloadPDF(message.FileUrl, tempFile); err != nil {
		return fmt.Errorf("下载PDF失败: %v", err)
	}
	defer os.Remove(tempFile) // 处理完成后删除临时文件

	// 创建打印机实例
	printer, err := printer.NewPrinter()
	if err != nil {
		return fmt.Errorf("创建打印机实例失败: %v", err)
	}
	defer printer.Close()

	// 如果指定了打印机，则使用指定的打印机
	if message.Printer != "" {
		if err := printer.Open(message.Printer); err != nil {
			return fmt.Errorf("打开打印机失败: %v", err)
		}
	}

	// 打印文件
	if err := printer.Print(tempFile); err != nil {
		return fmt.Errorf("打印失败: %v", err)
	}

	return nil
}

// 下载PDF文件
func downloadPDF(url string, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP请求返回错误状态码: %d", resp.StatusCode)
	}

	// 检查 Content-Type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(strings.ToLower(contentType), "pdf") {
		return fmt.Errorf("文件类型不是 PDF: %s", contentType)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	return nil
}

// 修改 handleWebSocket 函数
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket 升级失败: %v", err)
		return
	}
	defer conn.Close()

	// 创建打印机实例获取打印机列表
	printerInstance, err := printer.NewPrinter()
	if err != nil {
		log.Printf("创建打印机实例失败: %v", err)
		response := PrintResponse{Code: 500, Message: fmt.Sprintf("创建打印机实例失败: %v", err)}
		sendResponse(conn, response)
		return
	}
	defer printerInstance.Close()

	// 获取打印机列表
	printers, err := printerInstance.ListPrinters()
	if err != nil {
		log.Printf("获取打印机列表失败: %v", err)
		response := PrintResponse{Code: 500, Message: fmt.Sprintf("获取打印机列表失败: %v", err)}
		sendResponse(conn, response)
		return
	}

	// 发送打印机列表给客户端
	response := PrintResponse{
		Code:    200,
		Message: "success",
		Data:    printers,
	}
	if err := sendResponse(conn, response); err != nil {
		log.Printf("发送打印机列表失败: %v", err)
		return
	}

	// 继续处理其他消息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("读取消息错误: %v", err)
			break
		}

		log.Printf("收到消息: %s", message)

		// 解析消息
		var printMsg PrintMessage
		if err := json.Unmarshal(message, &printMsg); err != nil {
			log.Printf("解析消息失败: %v", err)
			response := PrintResponse{Code: 400, Message: fmt.Sprintf("消息格式错误: %v", err)}
			sendResponse(conn, response)
			continue
		}

		// 处理打印请求
		if err := handlePrintRequest(&printMsg); err != nil {
			log.Printf("处理打印请求失败: %v", err)
			response := PrintResponse{Code: 500, Message: fmt.Sprintf("打印失败: %v", err)}
			sendResponse(conn, response)
			continue
		}

		// 发送成功响应
		response := PrintResponse{Code: 200, Message: "打印成功"}
		sendResponse(conn, response)
	}
}

// 添加一个辅助函数来发送响应
func sendResponse(conn *websocket.Conn, response PrintResponse) error {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("响应序列化失败: %v", err)
	}
	return conn.WriteMessage(websocket.TextMessage, responseJSON)
}

func main() {
	// 加载配置文件
	config, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("加载配置失败:", err)
	}

	// 初始化日志
	initLogger(config)
	log.Println("服务启动，日志系统初始化完成")

	// 配置 WebSocket，允许所有来源的请求
	upgrader.CheckOrigin = func(r *http.Request) bool {
		// 如果配置启用了 CORS，则允许所有来源
		if config.WebSocket.EnableCORS {
			return true
		}

		// 获取请求的 Origin
		origin := r.Header.Get("Origin")
		// 可以在这里添加允许的域名列表
		allowedOrigins := []string{
			"http://localhost:" + strconv.Itoa(config.WebSocket.Port),
			"http://127.0.0.1:" + strconv.Itoa(config.WebSocket.Port),
		}

		// 检查是否是允许的域名
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				return true
			}
		}

		return false
	}

	// 设置路由
	http.HandleFunc("/ws", handleWebSocket)

	// 启动服务器
	port := ":" + strconv.Itoa(config.WebSocket.Port)
	log.Printf("WebSocket 服务器启动在 %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("HTTP 服务器启动失败:", err)
	}
}
