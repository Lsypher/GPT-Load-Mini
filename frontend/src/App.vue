<template>
  <el-config-provider :locale="zhCn">
    <div :data-theme="theme">
      <router-view />
    </div>
  </el-config-provider>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'

const theme = ref<'dark' | 'light'>('dark')

onMounted(() => {
  const saved = localStorage.getItem('theme') as 'dark' | 'light' | null
  if (saved) {
    theme.value = saved
  } else if (window.matchMedia('(prefers-color-scheme: light)').matches) {
    theme.value = 'light'
  }
})

watch(theme, (val) => {
  localStorage.setItem('theme', val)
})

// Expose for header component
window.__setTheme = (t: 'dark' | 'light') => {
  theme.value = t
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}

#app {
  min-height: 100vh;
}
</style>
