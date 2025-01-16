# Native Printer

Native Printer 是一个跨平台的本地打印服务，支持 Windows 和 Unix-like 系统。通过 WebSocket 接口提供网络打印功能。

## 功能特性

-   支持 Windows 和 Unix-like (Linux/macOS) 系统
-   WebSocket 接口，易于集成
-   支持 PDF 文件打印
-   自动获取本地打印机列表
-   可配置的日志系统
-   跨域支持（可配置）

## 系统要求

-   Go 1.23 或更高版本
-   Windows 系统需要安装打印机驱动
-   Unix-like 系统需要 CUPS 支持

## 安装

```bash
git clone https://github.com/bestk/native-printer.git
cd native-printer
go build
```

## 配置

配置文件 `config.yaml` 示例：

```yaml
log:
    folder: './logs'
websocket:
    port: 8080
    enableCORS: true
```

## 使用方法

1. 启动服务：

```bash
./native-printer
```

2. WebSocket 客户端示例：

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

// 打印文件
ws.send(
    JSON.stringify({
        action: 'printPDF',
        printer: 'Microsoft Print to PDF',
        fileUrl: 'https://example.com/document.pdf',
    }),
);
```

响应格式：

```json
{
    "code": 200,
    "message": "success",
    "data": null
}
```

## 开发

### 目录结构

```

├── config/ # 配置相关
├── printer/ # 打印机实现
├── pkg/ # 公共包
├── logs/ # 日志文件
├── main.go # 主程序
└── config.yaml # 配置文件
```

### 调试

使用 VS Code 进行调试，已配置 Windows 和 Mac 环境的调试配置。

## 许可证

MIT

## 贡献指南

欢迎提交 Issue 和 Pull Request。

## 更新日志

[版本更新记录]
