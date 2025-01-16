//go:build windows
// +build windows

package printer

import (
	"fmt"
	"os"

	"github.com/alexbrainman/printer"
)

type WindowsPrinter struct {
	printer *printer.Printer
}

func newSystemPrinter() (Printer, error) {
	return &WindowsPrinter{}, nil
}

func (p *WindowsPrinter) ListPrinters() ([]string, error) {
	return printer.ReadNames()
}

func (p *WindowsPrinter) Open(name string) error {
	printer, err := printer.Open(name)
	if err != nil {
		return err
	}
	p.printer = printer
	return nil
}

func (p *WindowsPrinter) Print(filePath string) error {
	if p.printer == nil {
		return fmt.Errorf("打印机未初始化")
	}

	if err := p.printer.StartDocument("Test", "RAW"); err != nil {
		return fmt.Errorf("开始文档失败: %v", err)
	}

	dat, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}
	_, err = p.printer.Write(dat)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}
	if err := p.printer.EndDocument(); err != nil {
		return fmt.Errorf("结束文档失败: %v", err)
	}

	return p.Close()
}

func (p *WindowsPrinter) Close() error {
	if p.printer != nil {
		return p.printer.Close()
	}
	return nil
}
