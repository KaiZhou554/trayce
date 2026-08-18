package trayicons

import "testing"

func TestSortEntries(t *testing.T) {
	entries := []TrayIconEntry{
		{ID: "1", ExecutablePath: `D:\Program Files\Steam\steam.exe`, Status: StatusValid},
		{ID: "2", ExecutablePath: `C:\Users\x\AppData\Local\Figma\app-126.1.2\Figma.exe`, Status: StatusMissing},
		{ID: "3", ExecutablePath: `{F38BF404-1D43-42F2-9305-67DE0B28FC23}\explorer.exe`, Status: StatusSpecial},
		{ID: "4", ExecutablePath: `C:\Users\x\AppData\Local\Figma\app-125.5.6\Figma.exe`, Status: StatusMissing},
		{ID: "5", ExecutablePath: `D:\tools\FastCopy\FastCopy.exe`, Status: StatusValid},
		{ID: "6", ExecutablePath: ``, Status: StatusUnknown},
	}
	got := SortEntries(entries)
	// 状态优先级：missing(0) < valid(1) < special(2) < unknown(3)
	// 同状态按程序名：Figma(4,2) < FastCopy(5) < steam(1)；Figma 版本内按完整路径字典序 app-125.5.6 < app-126.1.2
	wantOrder := []string{"4", "2", "5", "1", "3", "6"}
	if len(got) != len(wantOrder) {
		t.Fatalf("len = %d, want %d", len(got), len(wantOrder))
	}
	for i, id := range wantOrder {
		if got[i].ID != id {
			t.Errorf("position %d = %q, want %q (full order %v)", i, got[i].ID, id, idsOf(got))
		}
	}
	// Figma 两个版本必须相邻（版本聚合）
	figma := []string{}
	for _, e := range got {
		if BaseName(e.ExecutablePath) == "Figma" {
			figma = append(figma, e.ID)
		}
	}
	if len(figma) != 2 || figma[0] != "4" || figma[1] != "2" {
		t.Errorf("Figma versions not grouped: %v", figma)
	}
}

func TestBaseName(t *testing.T) {
	cases := []struct{ in, want string }{
		{`C:\Users\x\AppData\Local\Figma\app-126.1.2\Figma.exe`, "Figma"},
		{`D:\Program Files\Steam\steam.exe`, "steam"},
		{`{F38BF404-1D43-42F2-9305-67DE0B28FC23}\explorer.exe`, "explorer"},
		{``, ""},
	}
	for _, c := range cases {
		if got := BaseName(c.in); got != c.want {
			t.Errorf("BaseName(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func idsOf(es []TrayIconEntry) []string {
	out := make([]string, 0, len(es))
	for _, e := range es {
		out = append(out, e.ID)
	}
	return out
}
