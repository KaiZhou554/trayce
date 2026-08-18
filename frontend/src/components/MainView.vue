<template>
  <div class="h-screen bg-white/20 dark:bg-neutral-900/20 ring-1 ring-black/5 dark:ring-white/10 overflow-hidden flex flex-col transition-colors">
    <TitleBar
      :is-maximized="isMaximized"
      :app-icon="AppIcon"
      @minimize="WindowMinimise"
      @maximize="toggleMaximize"
      @close="Quit"
      @open-settings="settingsOpen = true"
    />

    <!-- 工具栏：搜索 + 过滤 -->
    <div class="px-4 pt-3 pb-2 flex flex-col gap-2 shrink-0">
      <n-input v-model:value="query" size="large" clearable :placeholder="t('search.placeholder')">
        <template #prefix>
          <Search24Regular class="w-4 h-4 text-neutral-400" />
        </template>
      </n-input>
      <n-radio-group :value="filter" size="small" @update:value="onFilterChange">
        <n-radio-button v-for="tab in tabOptions" :key="tab.value" :value="tab.value">
          {{ tab.label }}
        </n-radio-button>
      </n-radio-group>
    </div>

    <TrayIconList ref="listRef" :entries="filtered" @selection="selectedIds = $event" />

    <!-- 底部操作栏 -->
    <div class="px-4 py-3 shrink-0 flex items-center justify-between gap-3 border-t border-black/5 dark:border-white/10">
      <n-button :disabled="!canUndo" @click="askUndo">
        <template #icon><ArrowUndo24Regular class="w-4 h-4" /></template>
        {{ t('actions.undo') }}
      </n-button>
      <n-button type="error" :disabled="selectedIds.size === 0" @click="askDelete">
        {{ t('actions.deleteSelected') }}（{{ selectedIds.size }}）
      </n-button>
    </div>

    <ConfirmDialog
      :open="dialog.open"
      :title="dialog.title"
      :message="dialog.message"
      :confirm-text="dialog.confirmText"
      @confirm="onDialogConfirm"
      @cancel="dialog.open = false"
    />

    <SettingsDialog :show="settingsOpen" @update:show="settingsOpen = $event" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useMessage, NButton, NInput, NRadioGroup, NRadioButton } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { WindowMinimise, WindowMaximise, WindowUnmaximise, WindowIsMaximised, Quit } from '../../wailsjs/runtime/runtime'
import { Scan, DeleteEntries, UndoLastCleanup } from '../../wailsjs/go/main/App'
import { Search24Regular, ArrowUndo24Regular } from '@vicons/fluent'
import type { IconStatus, TrayIconEntry } from '../types'
import AppIcon from '../assets/appicon.png'
import TitleBar from './TitleBar.vue'
import TrayIconList from './TrayIconList.vue'
import ConfirmDialog from './ConfirmDialog.vue'
import SettingsDialog from './SettingsDialog.vue'

const { t } = useI18n()
const message = useMessage()

const all = ref<TrayIconEntry[]>([])
const query = ref('')
const filter = ref<'all' | IconStatus>('all')
const isMaximized = ref(false)
const settingsOpen = ref(false)
const selectedIds = ref<Set<string>>(new Set())
const listRef = ref<InstanceType<typeof TrayIconList> | null>(null)
const canUndo = ref(true)

const countOf = (key: string) =>
  key === 'all' ? all.value.length : all.value.filter(e => e.status === key).length

const tabOptions = computed(() => [
  { label: `${t('tabs.all')} (${countOf('all')})`, value: 'all' },
  { label: `${t('tabs.missing')} (${countOf('missing')})`, value: 'missing' },
  { label: `${t('tabs.valid')} (${countOf('valid')})`, value: 'valid' },
  { label: `${t('tabs.special')} (${countOf('special')})`, value: 'special' },
])

const onFilterChange = (v: string) => {
  filter.value = v as 'all' | IconStatus
}

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase()
  return all.value.filter(e => {
    if (filter.value !== 'all' && e.status !== filter.value) return false
    if (!q) return true
    return (
      e.id.toLowerCase().includes(q) ||
      e.publisher.toLowerCase().includes(q) ||
      e.executablePath.toLowerCase().includes(q)
    )
  })
})

const refresh = async () => {
  // wailsjs 生成的类型里 status 是 string，这里窄化为 IconStatus
  const rows = await Scan()
  all.value = rows.map(r => ({ ...r, status: r.status as IconStatus }))
  selectedIds.value = new Set()
  listRef.value?.clearSelection()
}

const dialog = ref<{
  open: boolean
  title: string
  message: string
  confirmText: string
  action: 'delete' | 'undo'
}>({ open: false, title: '', message: '', confirmText: '', action: 'delete' })

const askDelete = () => {
  dialog.value = {
    open: true,
    title: t('dialog.deleteTitle'),
    message: t('dialog.deleteMessage', { n: selectedIds.value.size }),
    confirmText: t('dialog.deleteConfirm'),
    action: 'delete',
  }
}

const askUndo = () => {
  dialog.value = {
    open: true,
    title: t('dialog.undoTitle'),
    message: t('dialog.undoMessage'),
    confirmText: t('dialog.undoConfirm'),
    action: 'undo',
  }
}

const onDialogConfirm = async () => {
  dialog.value.open = false
  try {
    if (dialog.value.action === 'delete') {
      const res = await DeleteEntries([...selectedIds.value])
      message.success(t('message.deleted', { n: res.deleted }))
    } else {
      const res = await UndoLastCleanup()
      message.success(res.deleted > 0 ? t('message.restored', { n: res.deleted }) : t('message.nothingToUndo'))
    }
    await refresh()
  } catch (err) {
    message.error(`${err}`)
  }
}

const toggleMaximize = async () => {
  if (isMaximized.value) {
    await WindowUnmaximise()
    isMaximized.value = false
  } else {
    await WindowMaximise()
    isMaximized.value = true
  }
}

onMounted(async () => {
  isMaximized.value = await WindowIsMaximised()
  await refresh()
})
</script>
