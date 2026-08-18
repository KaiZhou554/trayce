package trayicons

import "testing"

func TestDetermineStatus(t *testing.T) {
	neverExists := func(string) bool { return false }
	alwaysExists := func(string) bool { return true }

	cases := []struct {
		name   string
		path   string
		exists func(string) bool
		want   Status
	}{
		{"valid 文件存在", `D:\Program Files\Steam\steam.exe`, alwaysExists, StatusValid},
		{"missing 文件不存在", `C:\Users\no_such_user\AppData\Local\Figma\app-125.5.6\Figma.exe`, neverExists, StatusMissing},
		{"special GUID 前缀", `{F38BF404-1D43-42F2-9305-67DE0B28FC23}\explorer.exe`, neverExists, StatusSpecial},
		{"special GUID+子路径", `{6D809377-6AF0-444B-8957-A3773F02200E}\PowerToys\PowerToys.exe`, neverExists, StatusSpecial},
		{"special 小写 GUID", `{1ac14e77-02e7-4e5d-b744-2eb1ae5198b7}\Taskmgr.exe`, neverExists, StatusSpecial},
		{"special 即使文件存在也优先", `{F38BF404-1D43-42F2-9305-67DE0B28FC23}\explorer.exe`, alwaysExists, StatusSpecial},
		{"unknown 空路径", ``, neverExists, StatusUnknown},
		{"unknown 纯空白", `   `, neverExists, StatusUnknown},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := DetermineStatus(c.path, c.exists)
			if got != c.want {
				t.Errorf("DetermineStatus(%q) = %q, want %q", c.path, got, c.want)
			}
		})
	}
}

func TestDetermineStatusEnvVarFallback(t *testing.T) {
	// %VAR% 形式路径：先检查原样（失败），再展开环境变量检查（成功）→ valid
	called := []string{}
	exists := func(p string) bool {
		called = append(called, p)
		return len(called) == 2 // 第二次调用（展开后）返回 true
	}
	got := DetermineStatus(`%WINDIR%\notepad.exe`, exists)
	if got != StatusValid {
		t.Errorf("env-var path resolving to existing file should be valid, got %q", got)
	}
	if len(called) != 2 {
		t.Errorf("expected 2 existence checks (raw + expanded), got %d", len(called))
	}

	// 普通路径无 %VAR%：只检查一次
	called = nil
	exists2 := func(p string) bool { called = append(called, p); return false }
	got = DetermineStatus(`C:\plain\path.exe`, exists2)
	if got != StatusMissing {
		t.Errorf("plain missing path should be missing, got %q", got)
	}
	if len(called) != 1 {
		t.Errorf("plain path should be checked once, got %d calls", len(called))
	}
}

func TestIsSpecialPath(t *testing.T) {
	good := `{F38BF404-1D43-42F2-9305-67DE0B28FC23}\explorer.exe`
	if !IsSpecialPath(good) {
		t.Errorf("want special: %q", good)
	}
	bad := `C:\Program Files\Steam\steam.exe`
	if IsSpecialPath(bad) {
		t.Errorf("not special: %q", bad)
	}
	if IsSpecialPath(`F38BF404-1D43-42F2-9305-67DE0B28FC23\explorer.exe`) {
		t.Errorf("no braces should not be special")
	}
	if IsSpecialPath(`{not-a-guid}\explorer.exe`) {
		t.Errorf("invalid guid should not be special")
	}
}
