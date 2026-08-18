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
