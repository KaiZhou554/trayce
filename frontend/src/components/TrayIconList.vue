<template>
  <n-scrollbar style="flex: 1; min-height: 0" content-style="padding: 12px 16px 8px;">
    <div class="space-y-2">
      <div
        v-for="entry in entries"
        :key="entry.id"
        class="flex items-center gap-3 rounded-lg px-3 py-2.5 bg-white/70 dark:bg-neutral-800/60 ring-1 ring-black/5 dark:ring-white/10"
      >
        <n-checkbox :checked="selected.has(entry.id)" @update:checked="toggle(entry.id)" />
        <img
          :src="iconSrc(entry)"
          class="w-8 h-8 rounded-sm object-contain shrink-0"
          alt=""
          draggable="false"
        />
        <div class="min-w-0 flex-1">
          <div class="flex items-center gap-2">
            <span class="text-sm font-medium text-neutral-800 dark:text-neutral-100 truncate">
              {{ displayName(entry) }}
            </span>
            <n-tag size="small" :type="tagType(entry.status)" :bordered="false" class="shrink-0">
              {{ t(`status.${entry.status}`) }}
            </n-tag>
          </div>
          <div class="text-xs text-neutral-500 dark:text-neutral-400 truncate" :title="entry.executablePath">
            {{ entry.executablePath || t('list.noPath') }}
            <span v-if="entry.publisher" class="text-neutral-400 dark:text-neutral-500"> · {{ entry.publisher }}</span>
          </div>
        </div>
      </div>

      <div v-if="entries.length === 0" class="text-center text-sm text-neutral-400 py-16">
        {{ t('list.empty') }}
      </div>
    </div>
  </n-scrollbar>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NCheckbox, NScrollbar, NTag } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import type { IconStatus, TrayIconEntry } from '../types'

const { t } = useI18n()

defineProps<{ entries: TrayIconEntry[] }>()
const emit = defineEmits<{ (e: 'selection', ids: Set<string>): void }>()

const selected = ref<Set<string>>(new Set())

const toggle = (id: string) => {
  const next = new Set(selected.value)
  next.has(id) ? next.delete(id) : next.add(id)
  selected.value = next
  emit('selection', next)
}

// 程序名 = 文件名去扩展名（可靠推断，不做花哨猜测）；无路径则用 Publisher/ID 兜底
const displayName = (e: TrayIconEntry): string => {
  const base = e.executablePath.split(/[\\/]/).pop() ?? ''
  const name = base.replace(/\.exe$/i, '')
  if (name) return name
  if (e.publisher) return e.publisher
  return e.id
}

const iconSrc = (e: TrayIconEntry) =>
  e.iconBase64 ? `data:image/png;base64,${e.iconBase64}` : PlaceholderIcon

const tagType = (s: IconStatus) =>
  ({ valid: 'success', missing: 'warning', special: 'info', unknown: 'default' })[s] as
    | 'success'
    | 'warning'
    | 'info'
    | 'default'

const PlaceholderIcon =
  'data:image/svg+xml;utf8,' +
  encodeURIComponent(
    '<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32"><rect width="32" height="32" rx="6" fill="#e5e7eb"/><path d="M8 8h12l4 4v12H8z" fill="none" stroke="#9ca3af" stroke-width="2"/></svg>',
  )

defineExpose({ clearSelection: () => { selected.value = new Set(); emit('selection', selected.value) } })
</script>
