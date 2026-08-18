package trayicons

import (
	"os"
)

// FileExists 生产环境使用：os.Stat 存在性检查。
// 只有明确「不存在」（os.IsNotExist）才返回 false；
// 权限拒绝、路径过长等无法确认的情况一律视为存在（保守 —— 宁可不删，不要误删）。
func FileExists(p string) bool {
	_, err := os.Stat(p)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

// BuildEntries 把注册表原始记录映射为展示模型（含状态与图标），并排序。
func BuildEntries(raw []RawEntry, fileExists func(string) bool) []TrayIconEntry {
	entries := make([]TrayIconEntry, 0, len(raw))
	for _, r := range raw {
		entries = append(entries, TrayIconEntry{
			ID:             r.ID,
			IconGuid:       r.IconGuid,
			Publisher:      r.Publisher,
			ExecutablePath: r.ExecutablePath,
			IconBase64:     EncodeIconBase64(r.IconSnapshot),
			Status:         DetermineStatus(r.ExecutablePath, fileExists),
			IsSpecialPath:  IsSpecialPath(r.ExecutablePath),
		})
	}
	return SortEntries(entries)
}

// DeleteEntriesGuarded 执行「校验 -> 备份 -> 删除」。
// 服务端只校验「记录存在」：不存在的 id 忽略，不限制记录状态（用户可删除任意托盘记录）。
// 备份失败则中止删除。返回 (删除数量, 备份文件路径, error)。
func DeleteEntriesGuarded(src Source, requested []string, saveBackup func([]RawEntry) (string, error)) (int, string, error) {
	if len(requested) == 0 {
		return 0, "", nil
	}
	all, err := src.List()
	if err != nil {
		return 0, "", err
	}
	byID := map[string]RawEntry{}
	for _, e := range all {
		byID[e.ID] = e
	}
	targets := []RawEntry{}
	for _, id := range requested {
		if e, ok := byID[id]; ok {
			targets = append(targets, e)
		}
	}
	if len(targets) == 0 {
		return 0, "", nil
	}
	backupPath, err := saveBackup(targets)
	if err != nil {
		return 0, "", err
	}
	deleted := 0
	for _, e := range targets {
		if err := src.Delete(e.ID); err != nil {
			return deleted, backupPath, err // 返回实际成功数，便于前端提示与撤销
		}
		deleted++
	}
	return deleted, backupPath, nil
}
