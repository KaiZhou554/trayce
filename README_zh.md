[English](README.md) | 简体中文

# Trayce（通知区域图标管理器）

> 轻量、现代、无管理员权限的 Windows 11 实用工具 —— 查看并清理用户账户下残留的「其他系统托盘图标」记录。

基于 **Wails v2** 构建，Go 后端 + Vue 3 前端 + **Naive UI**，原生 Windows 毛玻璃窗口，无需联网。

## 特性

- **扫描查看** — 读取当前用户注册表 `HKCU\Control Panel\NotifyIconSettings` 的全部通知区域图标记录（只读，不修改注册表）
- **状态识别** — 按可执行文件存在性判断：`正常` / `路径失效` / `Windows 系统路径` / `未知`
  - 「路径失效」≠「软件已卸载」（升级、版本切换也会导致旧路径失效）
  - `{GUID}\...` 形式的 Windows 特殊路径会明确标注为「Windows 系统路径」
- **图标显示** — 直接解码 Windows 保存的 `IconSnapshot`（PNG）作为图标，数据损坏时显示占位图标
- **清理与撤销** — 可删除任意托盘记录（删除前确认并提示安全说明）；删除前自动备份到 `%LOCALAPPDATA%\unieditdept\trayce\backups\`，支持**撤销上次清理**
- **搜索与过滤** — 按名称 / 路径 / Publisher / ID 搜索；全部 / 失效 / 正常 / 特殊 一键过滤
- **国际化** — 简体中文 / English 界面，可在设置中切换（持久化于 localStorage）
- **权限友好** — 只操作当前用户注册表，无需管理员权限，不修改 HKLM、不删除程序文件、不联网

## 快速开始

### 环境要求

- **Go** 1.23+
- **Node.js** 18+
- **Wails CLI**（[安装指引](https://wails.io/docs/gettingstarted/installation)）

### 开发模式

```bash
# 安装前端依赖
cd frontend
npm install
cd ..

# 启动开发服务器（支持前端热更新）
wails dev
```

开发模式下可通过 http://localhost:34115 在浏览器中调试前端。

### 构建

```bash
wails build
```

产物位于 `build/bin/trayce.exe`。

### 测试

```bash
go test ./...
```

覆盖：状态判断（正常 / 失效 / 特殊 / 未知）、PNG 图标解码、排序与版本聚合、注册表扫描 / 删除 / 恢复（使用临时 HKCU 测试键，不触碰真实数据）、备份与撤销逻辑。

## 安全说明

- 注册表写操作**仅**针对 `HKCU\Control Panel\NotifyIconSettings` 及其子键
- 删除前强制确认，并说明「这只会删除 Windows 保存的通知区域图标记录，不会卸载软件，也不会删除程序文件」
- 每次删除前都会生成 JSON 备份（`%LOCALAPPDATA%\unieditdept\trayce\backups\`），可随时「撤销上次清理」
