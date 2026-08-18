package trayicons

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/sys/windows/registry"
)

// newTestSource 使用 HKCU 下随机临时键，测完即删，不触碰真实 NotifyIconSettings。
// 结构模拟真实场景：<parent>\<id> 是带 4 个值的子键。
func newTestSource(t *testing.T) (Source, func()) {
	t.Helper()
	parent := fmt.Sprintf(`Software\TrayceTest\%d`, time.Now().UnixNano())
	child := parent + `\10066634233605885324`
	_, _, err := registry.CreateKey(registry.CURRENT_USER, parent, registry.CREATE_SUB_KEY)
	if err != nil {
		t.Fatalf("create parent: %v", err)
	}
	k, _, err := registry.CreateKey(registry.CURRENT_USER, child, registry.SET_VALUE)
	if err != nil {
		t.Fatalf("create key: %v", err)
	}
	if err := k.SetStringValue("IconGuid", "{11111111-2222-3333-4444-555555555555}"); err != nil {
		t.Fatalf("set iconGuid: %v", err)
	}
	if err := k.SetStringValue("ExecutablePath", `D:\Program Files\Steam\steam.exe`); err != nil {
		t.Fatalf("set exePath: %v", err)
	}
	if err := k.SetStringValue("Publisher", "Steam"); err != nil {
		t.Fatalf("set publisher: %v", err)
	}
	if err := k.SetBinaryValue("IconSnapshot", []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x01, 0x02}); err != nil {
		t.Fatalf("set snapshot: %v", err)
	}
	k.Close()

	src := RegistrySource{Root: registry.CURRENT_USER, Path: parent}
	cleanup := func() {
		_ = registry.DeleteKey(registry.CURRENT_USER, child)
		_ = registry.DeleteKey(registry.CURRENT_USER, parent)
	}
	return src, cleanup
}

func TestRegistrySourceList(t *testing.T) {
	src, cleanup := newTestSource(t)
	defer cleanup()
	entries, err := src.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.ExecutablePath != `D:\Program Files\Steam\steam.exe` || e.Publisher != "Steam" {
		t.Errorf("unexpected entry: %+v", e)
	}
	if e.IconGuid != "{11111111-2222-3333-4444-555555555555}" {
		t.Errorf("iconGuid = %q", e.IconGuid)
	}
	if len(e.IconSnapshot) != 10 {
		t.Errorf("iconSnapshot len = %d, want 10", len(e.IconSnapshot))
	}
}

func TestRegistrySourceDelete(t *testing.T) {
	src, cleanup := newTestSource(t)
	defer cleanup()
	if err := src.Delete("someKeyThatDoesNotExist"); err == nil {
		t.Error("deleting non-existent key should error")
	}
	entries, _ := src.List()
	if err := src.Delete(entries[0].ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	after, _ := src.List()
	if len(after) != 0 {
		t.Errorf("after delete want 0 entries, got %d", len(after))
	}
}

func TestRegistrySourceRestore(t *testing.T) {
	src, cleanup := newTestSource(t)
	defer cleanup()
	orig, _ := src.List()
	if err := src.Delete(orig[0].ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := src.Restore(orig); err != nil {
		t.Fatalf("restore: %v", err)
	}
	after, _ := src.List()
	if len(after) != 1 {
		t.Fatalf("after restore want 1 entry, got %d", len(after))
	}
	if after[0].ID != orig[0].ID || after[0].ExecutablePath != orig[0].ExecutablePath {
		t.Errorf("restored mismatch: %+v", after[0])
	}
}

func TestRegistrySourceRestoreInvalidID(t *testing.T) {
	src, cleanup := newTestSource(t)
	defer cleanup()
	bad := []RawEntry{
		{ID: `evil\..\x`, ExecutablePath: `C:\a.exe`}, // 含反斜杠，试图嵌套
		{ID: `/abs/path`, ExecutablePath: `C:\a.exe`},  // 含斜杠
		{ID: `..`, ExecutablePath: `C:\a.exe`},         // 相对路径
		{ID: ``, ExecutablePath: `C:\a.exe`},           // 空
	}
	for _, e := range bad {
		if err := src.Restore([]RawEntry{e}); err == nil {
			t.Errorf("want error for invalid ID %q", e.ID)
		}
	}
}

func TestRegistrySourceRestoreSkipsExisting(t *testing.T) {
	src, cleanup := newTestSource(t)
	defer cleanup()
	orig, _ := src.List()
	// 键仍存在时恢复：跳过、不覆盖、不报错
	if err := src.Restore(orig); err != nil {
		t.Fatalf("restore existing: %v", err)
	}
	after, _ := src.List()
	if len(after) != 1 || after[0].ID != orig[0].ID {
		t.Fatalf("want original entry intact, got %+v", after)
	}
	// 值未被破坏
	if after[0].ExecutablePath != `D:\Program Files\Steam\steam.exe` {
		t.Errorf("existing value was overwritten: %q", after[0].ExecutablePath)
	}
}
