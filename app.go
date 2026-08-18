package main

import (
	"context"
	"fmt"

	"golang.org/x/sys/windows/registry"

	"trayce/internal/trayicons"
)

// App struct
type App struct {
	ctx    context.Context
	source trayicons.Source
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		source: trayicons.RegistrySource{
			Root: registry.CURRENT_USER,
			Path: `Control Panel\NotifyIconSettings`,
		},
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Scan 扫描当前用户的通知区域图标记录（只读，不修改注册表）。
func (a *App) Scan() ([]trayicons.TrayIconEntry, error) {
	raw, err := a.source.List()
	if err != nil {
		return nil, err
	}
	return trayicons.BuildEntries(raw, trayicons.FileExists), nil
}

// DeleteResult 返回删除操作的摘要
type DeleteResult struct {
	Deleted    int    `json:"deleted"`
	BackupPath string `json:"backupPath"`
}

// DeleteEntries 删除指定 id 的通知区域图标记录（删除前自动备份）。
// 不限制记录状态；部分删除失败时错误消息会带上已删数量与备份位置，便于用户撤销。
func (a *App) DeleteEntries(ids []string) (*DeleteResult, error) {
	deleted, backupPath, err := trayicons.DeleteEntriesGuarded(a.source, ids, trayicons.SaveBackup)
	if err != nil {
		if deleted > 0 {
			return nil, fmt.Errorf("已删除 %d 条后出错，其余已中止：%v；备份位于 %s", deleted, err, backupPath)
		}
		return nil, err
	}
	return &DeleteResult{Deleted: deleted, BackupPath: backupPath}, nil
}

// UndoLastCleanup 撤销上次清理：从最近一次备份恢复注册表记录。
func (a *App) UndoLastCleanup() (*DeleteResult, error) {
	path, entries, err := trayicons.LatestBackup()
	if err != nil {
		return nil, err
	}
	if entries == nil {
		return &DeleteResult{Deleted: 0, BackupPath: ""}, nil
	}
	if err := a.source.Restore(entries); err != nil {
		return nil, err
	}
	return &DeleteResult{Deleted: len(entries), BackupPath: path}, nil
}
