package trayicons

import "encoding/base64"

// pngMagic 是 PNG 文件头（IconSnapshot 的实际内容就是 PNG）
var pngMagic = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

// EncodeIconBase64 校验 PNG 头后转 base64；数据无效返回空串（前端显示占位图标）
func EncodeIconBase64(snapshot []byte) string {
	if len(snapshot) < len(pngMagic) {
		return ""
	}
	for i, b := range pngMagic {
		if snapshot[i] != b {
			return ""
		}
	}
	return base64.StdEncoding.EncodeToString(snapshot)
}
