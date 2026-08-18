<template>
  <n-config-provider :theme="theme">
    <n-message-provider>
      <MainView />
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { darkTheme, useOsTheme, NConfigProvider, NMessageProvider } from 'naive-ui'
import MainView from './components/MainView.vue'

const osTheme = useOsTheme()
const theme = computed(() => (osTheme.value === 'dark' ? darkTheme : null))
</script>
