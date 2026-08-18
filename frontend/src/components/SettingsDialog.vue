<template>
  <n-modal
    :show="show"
    preset="card"
    title="设置"
    :bordered="false"
    class="select-none"
    style="width: 420px; max-width: 90vw"
    @update:show="onShowChange"
  >
    <div class="space-y-5">
      <!-- 语言设置 -->
      <div>
        <div class="text-sm font-medium text-neutral-700 dark:text-neutral-200 mb-2">语言 / Language</div>
        <n-select v-model:value="lang" :options="langOptions" />
        <div class="text-xs text-neutral-400 mt-1">i18n 支持即将到来，当前仅保存选择。</div>
      </div>

      <!-- 关于 -->
      <div>
        <div class="text-sm font-medium text-neutral-700 dark:text-neutral-200 mb-2">关于</div>
        <div class="text-sm text-neutral-600 dark:text-neutral-300">
          <div class="font-semibold text-neutral-800 dark:text-neutral-100">Trayce</div>
          <div class="text-xs text-neutral-400 mt-0.5">Windows 通知区域图标记录管理器 · v1.0.0</div>
          <div class="flex items-center gap-1 mt-2 text-xs">
            <span class="text-neutral-400">项目主页</span>
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
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'

defineProps<{ show: boolean }>()
const emit = defineEmits<{ (e: 'update:show', v: boolean): void }>()

const lang = ref('zh-CN')
const langOptions = [
  { label: '简体中文', value: 'zh-CN' },
  { label: 'English', value: 'en' },
]

const onShowChange = (v: boolean) => emit('update:show', v)

const openRepo = () => BrowserOpenURL('https://github.com/KaiZhou554/trayce')
</script>
