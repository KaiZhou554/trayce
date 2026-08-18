package trayicons

// Status 表示一条托盘图标记录的状态
type Status string

const (
	StatusValid   Status = "valid"   // 文件存在，正常
	StatusMissing Status = "missing" // 文件不存在，路径失效（≠ 已卸载）
	StatusSpecial Status = "special" // {GUID}\... Windows 系统路径，保守处理
	StatusUnknown Status = "unknown" // 无路径或无法判断
)

// TrayIconEntry 是返回给前端的展示模型
type TrayIconEntry struct {
	ID             string `json:"id"`
	IconGuid       string `json:"iconGuid"`
	Publisher      string `json:"publisher"`
	ExecutablePath string `json:"executablePath"`
	IconBase64     string `json:"iconBase64"` // PNG 图标的 base64（无则空）
	Status         Status `json:"status"`
	IsSpecialPath  bool   `json:"isSpecialPath"`
}

// RawEntry 是注册表原始数据（扫描中间态 / 备份 / 恢复用）
type RawEntry struct {
	ID             string `json:"id"`
	IconGuid       string `json:"iconGuid"`
	ExecutablePath string `json:"executablePath"`
	Publisher      string `json:"publisher"`
	IconSnapshot   []byte `json:"iconSnapshot"`
}

// Backup 是一次删除前的备份文件内容
type Backup struct {
	CreatedAt string     `json:"createdAt"` // RFC3339
	Entries   []RawEntry `json:"entries"`
}
