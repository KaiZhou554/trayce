<template>
  <n-modal
    :show="show"
    preset="card"
    :title="t('common.settings')"
    :bordered="false"
    class="select-none"
    style="width: 420px; max-width: 90vw"
    @update:show="onShowChange"
  >
    <div class="space-y-5">
      <!-- 语言设置 -->
      <div>
        <div class="text-sm font-medium text-neutral-700 dark:text-neutral-200 mb-2">
          {{ t('settings.language') }}
        </div>
        <n-select v-model:value="lang" :options="langOptions" @update:value="onLangChange" />
      </div>

      <!-- 关于 -->
      <div>
        <div class="text-sm font-medium text-neutral-700 dark:text-neutral-200 mb-2">
          {{ t('settings.about') }}
        </div>
        <div class="text-sm text-neutral-600 dark:text-neutral-300">
          <div class="font-semibold text-neutral-800 dark:text-neutral-100">Trayce</div>
          <div class="text-xs text-neutral-400 mt-0.5">
            {{ t('settings.description') }} · v{{ version }}
          </div>
          <div class="flex items-center gap-1 mt-2 text-xs">
            <span class="text-neutral-400">{{ t('common.homepage') }}</span>
            <button class="text-sky-600 hover:underline" @click="openRepo">
              github.com/KaiZhou554/trayce
            </button>
          </div>
        </div>
      </div>
    </div>
  </n-modal>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NModal, NSelect } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'
import { currentLocale, setLocale, type Locale } from '../i18n'

const { t } = useI18n()

defineProps<{ show: boolean }>()
const emit = defineEmits<{ (e: 'update:show', v: boolean): void }>()

const version = '1.0.0'
const lang = ref<Locale>(currentLocale())
const langOptions = [
  { label: '简体中文', value: 'zh-CN' },
  { label: 'English', value: 'en' },
]

// 切换语言：更新全局 locale 并持久化到 localStorage
const onLangChange = (v: string) => {
  const locale = v as Locale
  lang.value = locale
  setLocale(locale)
}

const onShowChange = (v: boolean) => emit('update:show', v)

const openRepo = () => BrowserOpenURL('https://github.com/KaiZhou554/trayce')
</script>
