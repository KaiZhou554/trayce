<template>
  <header style="--wails-draggable:drag" class="h-12 flex items-center pl-4 shrink-0 select-none">
    <img :src="appIcon" class="h-5 w-5" alt="" />
    <span class="ml-3 text-sm font-medium text-neutral-700 dark:text-neutral-200">Trayce</span>

    <div style="--wails-draggable:no-drag" class="ml-auto flex items-center gap-1">
      <button
        class="dark:text-white w-8 h-8 flex items-center justify-center rounded-sm hover:bg-neutral-400/30 dark:hover:bg-neutral-900/30 active:bg-neutral-400/40 dark:active:bg-neutral-900/40 text-lg leading-none transition-colors"
        :title="t('common.settings')"
        @click="$emit('openSettings')"
      >
        <Icon><MoreHorizontal20Regular /></Icon>
      </button>
      <button
        @click="$emit('minimize')"
        class="dark:text-white w-12 h-12 flex items-center justify-center hover:bg-neutral-400/50 dark:hover:bg-neutral-900/50 active:bg-neutral-400 dark:active:bg-neutral-900 text-sm transition-colors"
        title="最小化"
      >
        <Icon><Subtract48Regular /></Icon>
      </button>
      <button
        @click="$emit('maximize')"
        class="dark:text-white w-12 h-12 flex items-center justify-center hover:bg-neutral-400/50 dark:hover:bg-neutral-900/50 active:bg-neutral-400 dark:active:bg-neutral-900 text-sm transition-colors"
        title="最大化"
      >
        <Icon>
          <component :is="isMaximized ? SquareMultiple24Regular : Tab28Regular" />
        </Icon>
      </button>
      <button
        @click="$emit('close')"
        class="dark:text-white w-12 h-12 flex items-center justify-center hover:bg-red-500 active:hover:bg-red-500/50 hover:text-white text-sm transition-colors"
        title="关闭"
      >
        <Icon><Dismiss48Regular /></Icon>
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
import {
  Dismiss48Regular,
  Tab28Regular,
  Subtract48Regular,
  SquareMultiple24Regular,
  MoreHorizontal20Regular,
} from '@vicons/fluent'
import { Icon } from '@vicons/utils'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps<{
  isMaximized: boolean
  appIcon: string
}>()

defineEmits<{
  openSettings: []
  minimize: []
  maximize: []
  close: []
}>()
</script>
