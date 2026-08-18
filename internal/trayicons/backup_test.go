package trayicons

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadLatestBackup(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("LOCALAPPDATA", dir) // BackupDir 读该变量

	entries := []RawEntry{
		{ID: "10066634233605885324", ExecutablePath: `C:\Users\x\AppData\Local\Figma\app-125.5.6\Figma.exe`, Publisher: "Figma"},
		{ID: "1015239650541068951", IconSnapshot: []byte{0x89, 0x50, 0x4E, 0x47}},
	}
	path, err := SaveBackup(entries)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if !filepath.IsAbs(path) || filepath.Dir(path) != filepath.Join(dir, "unieditdept", "trayce", "backups") {
		t.Errorf("backup path = %q", path)
	}

	latest, loaded, err := LatestBackup()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if latest != path {
		t.Errorf("latest = %q, want %q", latest, path)
	}
	if len(loaded) != 2 || loaded[0].ID != entries[0].ID || string(loaded[1].IconSnapshot) != string(entries[1].IconSnapshot) {
		t.Errorf("loaded mismatch: %+v", loaded)
	}
}

func TestSaveBackupTwiceDoesNotOverwrite(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("LOCALAPPDATA", dir)

	p1, err := SaveBackup([]RawEntry{{ID: "first", ExecutablePath: `C:\a\b.exe`}})
	if err != nil {
		t.Fatalf("save1: %v", err)
	}
	p2, err := SaveBackup([]RawEntry{{ID: "second", ExecutablePath: `C:\c\d.exe`}})
	if err != nil {
		t.Fatalf("save2: %v", err)
	}
	if p1 == p2 {
		t.Fatalf("two backups in same second must not share a filename: %q", p1)
	}
	latest, entries, err := LatestBackup()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if latest != p2 || len(entries) != 1 || entries[0].ID != "second" {
		t.Errorf("latest should be second backup, got %q %v", latest, entries)
	}
}

func TestLatestBackupNone(t *testing.T) {
	t.Setenv("LOCALAPPDATA", t.TempDir())
	path, entries, err := LatestBackup()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if path != "" || entries != nil {
		t.Errorf("want empty result, got %q %v", path, entries)
	}
}

func TestBackupDirFallbackToAppData(t *testing.T) {
	t.Setenv("LOCALAPPDATA", "")
	t.Setenv("APPDATA", t.TempDir())
	dir, err := BackupDir()
	if err != nil {
		t.Fatalf("backupdir: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("dir not created: %v", err)
	}
}
