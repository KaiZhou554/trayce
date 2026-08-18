package trayicons

import (
	"errors"
	"os"
	"testing"
)

type fakeSource struct {
	entries   []RawEntry
	deleted   []string
	restored  []RawEntry
	deleteErr func() error
}

func (f *fakeSource) List() ([]RawEntry, error) { return f.entries, nil }
func (f *fakeSource) Delete(id string) error {
	f.deleted = append(f.deleted, id)
	if f.deleteErr != nil {
		return f.deleteErr()
	}
	return nil
}
func (f *fakeSource) Restore(es []RawEntry) error { f.restored = es; return nil }

func TestBuildEntries(t *testing.T) {
	raw := []RawEntry{
		{ID: "1", ExecutablePath: `D:\Program Files\Steam\steam.exe`, IconSnapshot: []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}},
		{ID: "2", ExecutablePath: `{F38BF404-1D43-42F2-9305-67DE0B28FC23}\explorer.exe`},
		{ID: "3", ExecutablePath: `C:\nope\missing.exe`},
		{ID: "4"},
	}
	entries := BuildEntries(raw, func(p string) bool { return false })
	if entries[0].Status != StatusMissing {
		t.Errorf("first sorted entry should be missing, got %q", entries[0].Status)
	}
	byID := map[string]TrayIconEntry{}
	for _, e := range entries {
		byID[e.ID] = e
	}
	if byID["2"].Status != StatusSpecial || !byID["2"].IsSpecialPath {
		t.Errorf("special entry wrong: %+v", byID["2"])
	}
	if byID["3"].Status != StatusMissing {
		t.Errorf("missing entry wrong: %+v", byID["3"])
	}
	if byID["4"].Status != StatusUnknown {
		t.Errorf("unknown entry wrong: %+v", byID["4"])
	}
	if byID["1"].IconBase64 == "" {
		t.Error("valid png should produce base64")
	}
}

// DeleteEntriesGuarded 删除任意状态的记录（不限制状态）；不存在的 id 忽略
func TestDeleteEntriesGuardedDeletesAnyStatus(t *testing.T) {
	f := &fakeSource{entries: []RawEntry{
		{ID: "m1", ExecutablePath: `C:\nope\a.exe`},
		{ID: "v1", ExecutablePath: `D:\ok\b.exe`},
		{ID: "s1", ExecutablePath: `{F38BF404-1D43-42F2-9305-67DE0B28FC23}\explorer.exe`},
		{ID: "u1", ExecutablePath: ``},
	}}
	deleted, _, err := DeleteEntriesGuarded(f, []string{"m1", "v1", "s1", "u1", "nope"}, func([]RawEntry) (string, error) {
		return `X:\backup.json`, nil
	})
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if deleted != 4 {
		t.Errorf("deleted = %d, want 4", deleted)
	}
	if len(f.deleted) != 4 {
		t.Errorf("deleted %v, want 4 ids", f.deleted)
	}
}

func TestDeleteEntriesBackupFailureAborts(t *testing.T) {
	f := &fakeSource{entries: []RawEntry{
		{ID: "m1", ExecutablePath: `C:\nope\a.exe`},
	}}
	// 备份失败时必须中止删除
	_, _, err := DeleteEntriesGuarded(f, []string{"m1"}, func([]RawEntry) (string, error) {
		return "", errors.New("backup failed")
	})
	if err == nil {
		t.Error("want error when backup fails")
	}
	if len(f.deleted) != 0 {
		t.Errorf("deleted %v despite backup failure", f.deleted)
	}
}

func TestDeleteEntriesGuardedHappyPath(t *testing.T) {
	f := &fakeSource{entries: []RawEntry{
		{ID: "m1", ExecutablePath: `C:\nope\a.exe`},
		{ID: "m2", ExecutablePath: `C:\nope\b.exe`},
	}}
	deleted, path, err := DeleteEntriesGuarded(f, []string{"m1", "m2"}, func([]RawEntry) (string, error) {
		return `X:\backup.json`, nil
	})
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if deleted != 2 {
		t.Errorf("deleted = %d, want 2", deleted)
	}
	if path != `X:\backup.json` {
		t.Errorf("backup path = %q", path)
	}
	if len(f.deleted) != 2 || f.deleted[0] != "m1" || f.deleted[1] != "m2" {
		t.Errorf("deleted = %v", f.deleted)
	}
}

func TestDeleteEntriesGuardedEmptyRequest(t *testing.T) {
	f := &fakeSource{}
	deleted, _, err := DeleteEntriesGuarded(f, nil, func([]RawEntry) (string, error) {
		return "", errors.New("should not be called")
	})
	if err != nil || deleted != 0 {
		t.Errorf("empty request: deleted=%d err=%v", deleted, err)
	}
}

func TestDeleteEntriesGuardedPartialFailureCountsSuccesses(t *testing.T) {
	f := &fakeSource{
		entries: []RawEntry{
			{ID: "m1", ExecutablePath: `C:\nope\a.exe`},
			{ID: "m2", ExecutablePath: `C:\nope\b.exe`},
			{ID: "m3", ExecutablePath: `C:\nope\c.exe`},
		},
	}
	f.deleteErr = func() error {
		if len(f.deleted) >= 2 {
			return errors.New("second delete failed")
		}
		return nil
	}
	deleted, path, err := DeleteEntriesGuarded(f, []string{"m1", "m2", "m3"}, func([]RawEntry) (string, error) {
		return `X:\backup.json`, nil
	})
	if err == nil {
		t.Error("want error on partial failure")
	}
	if deleted != 1 { // 只有 m1 成功，m2 失败即中止
		t.Errorf("deleted = %d, want 1 (actual successes)", deleted)
	}
	if path != `X:\backup.json` {
		t.Errorf("backup path = %q", path)
	}
}

func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	existing := dir + `\exists.txt`
	if err := os.WriteFile(existing, []byte("x"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if !FileExists(existing) {
		t.Error("existing file should be found")
	}
	if FileExists(dir + `\missing.txt`) {
		t.Error("missing file should not be found")
	}
	if !FileExists(dir) {
		t.Error("existing directory should be found (os.Stat succeeds)")
	}
}
