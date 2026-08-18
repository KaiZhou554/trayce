export type IconStatus = 'valid' | 'missing' | 'special' | 'unknown'

export interface TrayIconEntry {
  id: string
  iconGuid: string
  publisher: string
  executablePath: string
  iconBase64: string
  status: IconStatus
  isSpecialPath: boolean
}

export interface DeleteResult {
  deleted: number
  backupPath: string
}

export const STATUS_LABEL: Record<IconStatus, string> = {
  valid: '正常',
  missing: '路径失效',
  special: 'Windows 系统路径',
  unknown: '未知',
}

export const STATUS_ORDER: IconStatus[] = ['missing', 'valid', 'special', 'unknown']
