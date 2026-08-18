package trayicons

import (
	"path/filepath"
	"sort"
	"strings"
)

var statusRank = map[Status]int{
	StatusMissing: 0,
	StatusValid:   1,
	StatusSpecial: 2,
	StatusUnknown: 3,
}

// BaseName 从可执行路径提取程序名（文件名去掉扩展名）。
// 不做花哨的软件名猜测 —— 文件名即名称。
func BaseName(exePath string) string {
	base := filepath.Base(exePath)
	if ext := filepath.Ext(base); ext != "" {
		base = strings.TrimSuffix(base, ext)
	}
	return base
}

// SortEntries 排序：状态（失效优先）-> 程序名 -> 完整路径。
// 相同程序名的多条记录（多版本）自然聚合在一起。
func SortEntries(entries []TrayIconEntry) []TrayIconEntry {
	out := append([]TrayIconEntry(nil), entries...)
	sort.SliceStable(out, func(i, j int) bool {
		ri, rj := statusRank[out[i].Status], statusRank[out[j].Status]
		if ri != rj {
			return ri < rj
		}
		ni, nj := strings.ToLower(BaseName(out[i].ExecutablePath)), strings.ToLower(BaseName(out[j].ExecutablePath))
		if ni != nj {
			return ni < nj
		}
		return strings.ToLower(out[i].ExecutablePath) < strings.ToLower(out[j].ExecutablePath)
	})
	return out
}
