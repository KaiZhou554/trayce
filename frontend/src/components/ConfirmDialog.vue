<template>
  <n-modal
    :show="open"
    preset="card"
    :title="title"
    :bordered="false"
    class="select-none"
    style="width: 420px; max-width: 90vw"
    @update:show="onShowChange"
  >
    <p class="text-sm text-neutral-600 dark:text-neutral-300 leading-relaxed">{{ message }}</p>
    <template #footer>
      <div class="flex justify-end gap-2">
        <n-button @click="$emit('cancel')">{{ t('common.cancel') }}</n-button>
        <n-button type="error" @click="$emit('confirm')">{{ confirmText }}</n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { NButton, NModal } from 'naive-ui'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps<{
  open: boolean
  title: string
  message: string
  confirmText: string
}>()
const emit = defineEmits<{ (e: 'confirm'): void; (e: 'cancel'): void }>()

// 按 ESC / 点击遮罩关闭时视为取消
const onShowChange = (show: boolean) => {
  if (!show) emit('cancel')
}
</script>
