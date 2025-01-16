package printer

type Printer interface {
	// 列出所有打印机
	ListPrinters() ([]string, error)
	// 打开打印机
	Open(name string) error
	// 打印文件
	Print(filePath string) error
	// 关闭打印机
	Close() error
}

// 创建打印机实例
func NewPrinter() (Printer, error) {
	return newSystemPrinter()
}
