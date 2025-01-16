//go:build linux || darwin
// +build linux darwin

package printer

import (
	"os/exec"
	"strings"
)

type UnixPrinter struct {
	printerName string
}

func newSystemPrinter() (Printer, error) {
	return &UnixPrinter{}, nil
}

func (p *UnixPrinter) ListPrinters() ([]string, error) {
	cmd := exec.Command("lpstat", "-p")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// 解析输出获取打印机名称
	printers := []string{}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "printer") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				printers = append(printers, fields[1])
			}
		}
	}

	return printers, nil
}

func (p *UnixPrinter) Open(name string) error {
	p.printerName = name
	return nil
}

func (p *UnixPrinter) Print(filePath string) error {
	cmd := exec.Command("lp", "-d", p.printerName, filePath)
	return cmd.Run()
}

func (p *UnixPrinter) Close() error {
	return nil
}
