import { createI18n } from 'vue-i18n'
import zhCN from './locales/zh-CN'
import en from './locales/en'

export type Locale = 'zh-CN' | 'en'

const STORAGE_KEY = 'trayce-lang'

function loadSavedLocale(): Locale {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved === 'zh-CN' || saved === 'en') return saved
  } catch {
    // localStorage 不可用时回退默认
  }
  return 'zh-CN'
}

export const i18n = createI18n({
  legacy: false,
  locale: loadSavedLocale(),
  fallbackLocale: 'zh-CN',
  messages: { 'zh-CN': zhCN, en },
})

// 切换语言并持久化到 localStorage
export function setLocale(locale: Locale) {
  i18n.global.locale.value = locale
  try {
    localStorage.setItem(STORAGE_KEY, locale)
  } catch {
    // ignore
  }
}

export function currentLocale(): Locale {
  return i18n.global.locale.value as Locale
}
