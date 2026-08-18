<template>
  <n-modal
    :show="open"
    preset="card"
    :title="title"
    :bordered="false"
    style="width: 420px; max-width: 90vw"
    @update:show="onShowChange"
  >
    <p class="text-sm text-neutral-600 dark:text-neutral-300 leading-relaxed">{{ message }}</p>
    <div v-if="danger" class="mt-3 text-xs text-amber-600 dark:text-amber-400 leading-relaxed">
      这只会删除 Windows 保存的通知区域图标记录，不会卸载软件，也不会删除程序文件。
    </div>
    <template #footer>
      <div class="flex justify-end gap-2">
        <n-button @click="$emit('cancel')">取消</n-button>
        <n-button type="error" @click="$emit('confirm')">{{ confirmText }}</n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { NButton, NModal } from 'naive-ui'

defineProps<{
  open: boolean
  title: string
  message: string
  confirmText: string
  danger?: boolean
}>()
const emit = defineEmits<{ (e: 'confirm'): void; (e: 'cancel'): void }>()

// 按 ESC / 点击遮罩关闭时视为取消
const onShowChange = (show: boolean) => {
  if (!show) emit('cancel')
}
</script>
