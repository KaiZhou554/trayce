package trayicons

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// Source 抽象注册表读写，便于 App 层注入 fake 做单元测试。
type Source interface {
	List() ([]RawEntry, error)
	Delete(id string) error
	Restore(entries []RawEntry) error
}

// RegistrySource 是真实的注册表实现。
// 只操作 Root\Path 及其子键 —— 生产环境 Root=HKCU, Path=Control Panel\NotifyIconSettings。
type RegistrySource struct {
	Root registry.Key
	Path string
}

// List 枚举所有子键并读取 4 个值。
// 单个子键读取失败时跳过该键（保守，不因一条坏记录中断整个扫描）。
func (s RegistrySource) List() ([]RawEntry, error) {
	k, err := registry.OpenKey(s.Root, s.Path, registry.READ|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return nil, err
	}
	defer k.Close()
	names, err := k.ReadSubKeyNames(0)
	if err != nil {
		return nil, err
	}
	entries := make([]RawEntry, 0, len(names))
	for _, name := range names {
		sk, err := registry.OpenKey(k, name, registry.READ)
		if err != nil {
			continue
		}
		iconGuid, _, _ := sk.GetStringValue("IconGuid")
		exePath, _, _ := sk.GetStringValue("ExecutablePath")
		publisher, _, _ := sk.GetStringValue("Publisher")
		snapshot, _, _ := sk.GetBinaryValue("IconSnapshot")
		sk.Close()
		entries = append(entries, RawEntry{
			ID: name, IconGuid: iconGuid, ExecutablePath: exePath,
			Publisher: publisher, IconSnapshot: snapshot,
		})
	}
	return entries, nil
}

// Delete 删除子键（等价 Remove-Item，但直接调用 Registry API）。
func (s RegistrySource) Delete(id string) error {
	return registry.DeleteKey(s.Root, s.Path+`\`+id)
}

// Restore 重新创建子键并写回全部值（用于撤销清理）。
// 安全校验：ID 必须是合法子键名（拒绝 \ /、空、.、..）；键已存在则跳过，不覆盖新数据。
func (s RegistrySource) Restore(entries []RawEntry) error {
	parent, err := registry.OpenKey(s.Root, s.Path, registry.CREATE_SUB_KEY)
	if err != nil {
		return err
	}
	defer parent.Close()
	for _, e := range entries {
		if !validKeyName(e.ID) {
			return fmt.Errorf("invalid subkey name %q", e.ID)
		}
		k, existed, err := registry.CreateKey(parent, e.ID, registry.SET_VALUE)
		if err != nil {
			return err
		}
		if existed {
			k.Close()
			continue // 键已存在（系统或应用重建），跳过不覆盖
		}
		if err := k.SetStringValue("IconGuid", e.IconGuid); err != nil {
			k.Close()
			return err
		}
		if err := k.SetStringValue("ExecutablePath", e.ExecutablePath); err != nil {
			k.Close()
			return err
		}
		if err := k.SetStringValue("Publisher", e.Publisher); err != nil {
			k.Close()
			return err
		}
		if len(e.IconSnapshot) > 0 {
			if err := k.SetBinaryValue("IconSnapshot", e.IconSnapshot); err != nil {
				k.Close()
				return err
			}
		}
		k.Close()
	}
	return nil
}

// validKeyName 校验注册表子键名合法：非空、非相对路径、不含路径分隔符或 NUL。
func validKeyName(id string) bool {
	if id == "" || id == "." || id == ".." {
		return false
	}
	if strings.ContainsAny(id, `\/`) {
		return false
	}
	return !strings.ContainsRune(id, 0)
}
