package trayicons

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

// backupSeq 保证同一时间戳内文件名唯一（Windows 时钟存在 100ns 节拍，仅靠纳秒时间戳可能撞车）
var backupSeq atomic.Uint64

// BackupDir 返回备份目录（%LOCALAPPDATA%\unieditdept\trayce\backups），不存在则创建。
func BackupDir() (string, error) {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" {
		base = os.Getenv("APPDATA")
	}
	dir := filepath.Join(base, "unieditdept", "trayce", "backups")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

// SaveBackup 把即将删除的记录写为 JSON 备份，返回文件路径。
func SaveBackup(entries []RawEntry) (string, error) {
	dir, err := BackupDir()
	if err != nil {
		return "", err
	}
	seq := backupSeq.Add(1)
	path := filepath.Join(dir, fmt.Sprintf("trayce-backup-%s-%06d.json", time.Now().Format("20060102-150405.000000000"), seq))
	data, err := json.MarshalIndent(Backup{
		CreatedAt: time.Now().Format(time.RFC3339),
		Entries:   entries,
	}, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// LatestBackup 返回最近一次备份的文件路径与内容；无备份时返回 ("", nil, nil)。
// 按文件修改时间取最新，兼容旧版秒级时间戳文件名的备份。
func LatestBackup() (string, []RawEntry, error) {
	dir, err := BackupDir()
	if err != nil {
		return "", nil, err
	}
	matches, err := filepath.Glob(filepath.Join(dir, "trayce-backup-*.json"))
	if err != nil {
		return "", nil, err
	}
	if len(matches) == 0 {
		return "", nil, nil
	}
	latest := matches[0]
	var latestMod time.Time
	var latestName string
	for _, m := range matches {
		fi, err := os.Stat(m)
		if err != nil {
			continue
		}
		// mtime 优先；同 mtime（同一秒内多次备份）时按文件名后缀（纳秒时间戳）决胜
		if fi.ModTime().After(latestMod) || (fi.ModTime().Equal(latestMod) && m > latestName) {
			latestMod = fi.ModTime()
			latestName = m
			latest = m
		}
	}
	data, err := os.ReadFile(latest)
	if err != nil {
		return "", nil, err
	}
	var b Backup
	if err := json.Unmarshal(data, &b); err != nil {
		return "", nil, err
	}
	return latest, b.Entries, nil
}
