package trayicons

import (
	"os"
	"regexp"
	"strings"
)

// guidPrefixRe 匹配 {GUID}\ 开头的 Windows 特殊路径
var guidPrefixRe = regexp.MustCompile(`^\{[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}\}\\`)

// IsSpecialPath 判断路径是否为 {GUID}\... 形式（如 {1AC14E77-...}\Taskmgr.exe）
func IsSpecialPath(path string) bool {
	return guidPrefixRe.MatchString(path)
}

// envVarRe 匹配 Windows 环境变量引用 %VAR%
var envVarRe = regexp.MustCompile(`%([^%]+)%`)

// expandWindowsEnv 展开 Windows 风格的 %VAR% 路径引用；
// 未定义或空值的变量保持原样（保守，不因展开失败误判）。
func expandWindowsEnv(path string) string {
	return envVarRe.ReplaceAllStringFunc(path, func(m string) string {
		if v := os.Getenv(m[1 : len(m)-1]); v != "" {
			return v
		}
		return m
	})
}

// DetermineStatus 根据可执行路径判断状态。fileExists 可注入以便测试。
func DetermineStatus(exePath string, fileExists func(string) bool) Status {
	if strings.TrimSpace(exePath) == "" {
		return StatusUnknown
	}
	if IsSpecialPath(exePath) {
		return StatusSpecial
	}
	if fileExists(exePath) {
		return StatusValid
	}
	// 兼容 %VAR% 形式路径：展开环境变量后再次检查（保守，减少误判）
	expanded := expandWindowsEnv(exePath)
	if expanded != exePath && fileExists(expanded) {
		return StatusValid
	}
	return StatusMissing
}
