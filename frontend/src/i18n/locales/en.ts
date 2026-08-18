export default {
  common: {
    more: 'More',
    cancel: 'Cancel',
    homepage: 'Project homepage',
    minimize: 'Minimize',
    maximize: 'Maximize',
    close: 'Close',
  },
  search: {
    placeholder: 'Search name, path or ID…',
  },
  tabs: {
    all: 'All',
    missing: 'Missing',
    valid: 'Valid',
    special: 'Special',
  },
  status: {
    valid: 'Valid',
    missing: 'Path missing',
    special: 'Windows system path',
    unknown: 'Unknown',
  },
  list: {
    empty: 'No matching entries',
    noPath: '(no path)',
  },
  actions: {
    undo: 'Undo last cleanup',
    deleteSelected: 'Delete selected',
  },
  dialog: {
    deleteTitle: 'Delete selected entries',
    deleteMessage:
      'This will delete {n} notification area icon record(s). It only removes Windows-saved tray icon records and will not uninstall software or delete program files.',
    deleteConfirm: 'Delete',
    undoTitle: 'Undo last cleanup',
    undoMessage: 'This will restore the notification area icon records from the most recent backup.',
    undoConfirm: 'Restore',
  },
  message: {
    deleted: 'Deleted {n} record(s).\nThe Windows Settings page may need to be reopened to reflect changes.',
    restored: 'Restored {n} record(s).',
    nothingToUndo: 'Nothing to undo.',
  },
  settings: {
    language: 'Language',
    notesTitle: 'What happens if I delete an entry?',
    noteNotUninstalled: 'Not uninstalled software: it will recreate the entry on its next run',
    noteRunningApps: 'Running apps: their tray icon will be collapsed',
    about: 'About',
    description: 'Windows notification area icon records manager',
  },
}
