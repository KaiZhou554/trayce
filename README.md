English | [简体中文](README_zh.md)

# Trayce (Notification Area Icon Manager)

> A lightweight, modern Windows 11 utility that runs without administrator privileges — view and clean up stale "Other system tray icons" records left in your user account.

Built with **Wails v2** (Go backend + Vue 3 frontend + **Naive UI**), native Windows acrylic window, fully offline.

## Features

- **Scan & view** — Read all notification area icon records from `HKCU\Control Panel\NotifyIconSettings` (read-only, never modifies the registry)
- **Status detection** — Determines status by checking the executable's existence: `Valid` / `Path missing` / `Windows system path` / `Unknown`
  - "Path missing" ≠ "Uninstalled" (upgrades and version switches also leave old paths behind)
  - `{GUID}\...` special Windows paths are clearly labeled as "Windows system path"
- **Icon display** — Decodes Windows' saved `IconSnapshot` (PNG) directly as the icon; shows a placeholder when the data is corrupted
- **Cleanup & undo** — Delete any tray icon record (with a confirmation dialog and safety note); every deletion is backed up to `%LOCALAPPDATA%\unieditdept\trayce\backups\` first, and you can **undo the last cleanup** anytime
- **Search & filter** — Search by name / path / Publisher / ID; one-click filters for All / Missing / Valid / Special
- **i18n** — Simplified Chinese and English UI, switchable in Settings (persisted in localStorage)
- **Privacy-friendly** — Only touches the current user's registry, no administrator rights, no HKLM writes, no file deletion, no network access

## Getting Started

### Prerequisites

- **Go** 1.23+
- **Node.js** 18+
- **Wails CLI** ([installation guide](https://wails.io/docs/gettingstarted/installation))

### Development

```bash
# Install frontend dependencies
cd frontend
npm install
cd ..

# Start the dev server (with frontend hot reload)
wails dev
```

In dev mode you can debug the frontend in a browser at http://localhost:34115.

### Build

```bash
wails build
```

The binary is produced at `build/bin/trayce.exe`.

### Tests

```bash
go test ./...
```

Covers: status detection (valid / missing / special / unknown), PNG icon decoding, sorting & version grouping, registry scan / delete / restore (using temporary HKCU test keys — never touches real data), and backup & undo logic.

## Safety Notes

- Registry write operations target **only** `HKCU\Control Panel\NotifyIconSettings` and its subkeys
- Deletion always asks for confirmation and explains: "This only removes Windows-saved notification area icon records — it will not uninstall software or delete program files."
- A JSON backup (`%LOCALAPPDATA%\unieditdept\trayce\backups\`) is created before every deletion, so you can always "undo the last cleanup"
