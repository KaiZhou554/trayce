export default {
  common: {
    settings: '设置',
    cancel: '取消',
    homepage: '项目主页',
  },
  search: {
    placeholder: '搜索名称、路径或 ID…',
  },
  tabs: {
    all: '全部',
    missing: '失效',
    valid: '正常',
    special: '特殊',
  },
  status: {
    valid: '正常',
    missing: '路径失效',
    special: 'Windows 系统路径',
    unknown: '未知',
  },
  list: {
    empty: '没有匹配的记录',
    noPath: '(无路径)',
  },
  actions: {
    undo: '撤销上次清理',
    deleteSelected: '删除所选记录',
  },
  dialog: {
    deleteTitle: '确认删除所选记录',
    deleteMessage:
      '将删除 {n} 条通知区域图标记录。这只会删除 Windows 保存的通知区域图标记录，不会卸载软件，也不会删除程序文件。',
    deleteConfirm: '删除记录',
    undoTitle: '撤销上次清理',
    undoMessage: '将根据最近一次备份，恢复被清理的通知区域图标记录。',
    undoConfirm: '恢复',
  },
  message: {
    deleted: '已删除 {n} 条记录。\nWindows 设置页面可能需要重新打开才能看到变化。',
    restored: '已恢复 {n} 条记录。',
    nothingToUndo: '没有可撤销的记录。',
  },
  settings: {
    language: '语言 / Language',
    about: '关于',
    description: 'Windows 通知区域图标记录管理器',
  },
}
