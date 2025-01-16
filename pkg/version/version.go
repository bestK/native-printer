package version

import "fmt"

var (
	// 编译时通过 -ldflags 设置这些变量
	Version    = "dev"
	CommitHash = "none"
	BuildTime  = "unknown"
)

// GetVersionInfo 返回版本信息
func GetVersionInfo() string {
	return fmt.Sprintf("\nVersion: %s\nCommit: %s\nBuildTime: %s",
		Version,
		CommitHash,
		BuildTime,
	)
}

func GetAppName() string {
	return fmt.Sprintf("Native Printer %s", Version)
}
