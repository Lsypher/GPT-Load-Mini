<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">设置</h1>
      <p class="page-subtitle">配置 API 连接和身份验证</p>
    </div>

    <div class="settings-grid">
      <!-- API Configuration -->
      <div class="card settings-card">
        <div class="card-header">
          <div class="card-title-group">
            <svg class="card-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="3"/>
              <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
            </svg>
            <span class="card-title">API 配置</span>
          </div>
        </div>
        <div class="card-body">
          <div class="form-group">
            <label class="label">API 基础 URL</label>
            <input
              v-model="apiBaseUrl"
              type="text"
              class="input"
              placeholder="http://localhost:8080"
            />
          </div>
          <div class="form-group">
            <label class="label">认证密钥</label>
            <input
              v-model="authKey"
              type="text"
              class="input"
              placeholder="您的认证密钥"
            />
          </div>
        </div>
        <div class="card-footer">
          <button class="btn btn-warning" @click="handleReload">重新加载配置</button>
          <button class="btn btn-secondary" @click="handleTest" :disabled="testing">
            {{ testing ? '测试中...' : '测试连接' }}
          </button>
          <button class="btn btn-primary" @click="handleSave">保存</button>
        </div>
      </div>

      <!-- Appearance -->
      <div class="card settings-card">
        <div class="card-header">
          <div class="card-title-group">
            <svg class="card-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="5"/>
              <line x1="12" y1="1" x2="12" y2="3"/>
              <line x1="12" y1="21" x2="12" y2="23"/>
              <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
              <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
              <line x1="1" y1="12" x2="3" y2="12"/>
              <line x1="21" y1="12" x2="23" y2="12"/>
              <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
              <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
            </svg>
            <span class="card-title">外观</span>
          </div>
        </div>
        <div class="card-body">
          <div class="theme-selector">
            <button
              class="theme-option"
              :class="{ active: theme === 'dark' }"
              @click="setTheme('dark')"
            >
              <div class="theme-preview theme-preview-dark">
                <div class="preview-sidebar"></div>
                <div class="preview-content">
                  <div class="preview-header"></div>
                  <div class="preview-line"></div>
                  <div class="preview-line short"></div>
                </div>
              </div>
              <span class="theme-name">深色</span>
            </button>
            <button
              class="theme-option"
              :class="{ active: theme === 'light' }"
              @click="setTheme('light')"
            >
              <div class="theme-preview theme-preview-light">
                <div class="preview-sidebar"></div>
                <div class="preview-content">
                  <div class="preview-header"></div>
                  <div class="preview-line"></div>
                  <div class="preview-line short"></div>
                </div>
              </div>
              <span class="theme-name">浅色</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useSettingsStore } from '@/stores/settings'
import { healthApi } from '@/api/health'
import { settingsApi } from '@/api/settings'

const settingsStore = useSettingsStore()
const apiBaseUrl = ref(settingsStore.apiBaseUrl)
const authKey = ref(settingsStore.authKey)
const testing = ref(false)
const theme = ref<'dark' | 'light'>('dark')

function handleSave() {
  settingsStore.updateSettings(apiBaseUrl.value, authKey.value)
  ElMessage.success('设置已保存')
}

async function handleTest() {
  testing.value = true
  try {
    await healthApi.check(apiBaseUrl.value)
    ElMessage.success('连接成功')
  } catch {
    ElMessage.error('连接失败')
  } finally {
    testing.value = false
  }
}

async function handleReload() {
  try {
    await settingsApi.reloadConfig()
    ElMessage.success('配置已重新加载')
  } catch {}
}

function setTheme(t: 'dark' | 'light') {
  theme.value = t
  localStorage.setItem('theme', t)
  window.__setTheme?.(t)
}

onMounted(() => {
  const saved = localStorage.getItem('theme') as 'dark' | 'light' | null
  if (saved) {
    theme.value = saved
  }
})
</script>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.page-header {
  margin-bottom: var(--space-2);
}

.page-title {
  font-size: 1.25rem;
  font-weight: 600;
}

.page-subtitle {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-top: var(--space-1);
}

.settings-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(340px, 1fr));
  gap: var(--space-5);
}

.settings-card {
  display: flex;
  flex-direction: column;
}

.card-title-group {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.card-icon {
  width: 20px;
  height: 20px;
  color: var(--accent);
}

.card-title {
  font-size: 0.9375rem;
  font-weight: 600;
}

.card-body {
  padding: var(--space-5);
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  flex: 1;
}

.card-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  padding: var(--space-4) var(--space-5);
  border-top: 1px solid var(--border-subtle);
}

.btn-warning {
  background: var(--accent);
  color: var(--bg-primary);
  border: 1px solid var(--accent);
}

.btn-warning:hover {
  background: var(--accent-hover);
  border-color: var(--accent-hover);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.input {
  width: 100%;
}

/* Theme Selector */
.theme-selector {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-3);
}

.theme-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-4);
  background: var(--bg-primary);
  border: 2px solid var(--border);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.theme-option:hover {
  border-color: var(--text-tertiary);
}

.theme-option.active {
  border-color: var(--accent);
  background: var(--accent-muted);
}

.theme-preview {
  width: 100%;
  height: 80px;
  border-radius: var(--radius-md);
  overflow: hidden;
  display: flex;
  border: 1px solid var(--border);
}

.theme-preview-dark {
  background: #0d1117;
}

.theme-preview-light {
  background: #f6f8fa;
}

.preview-sidebar {
  width: 30%;
  height: 100%;
}

.theme-preview-dark .preview-sidebar {
  background: #161b22;
}

.theme-preview-light .preview-sidebar {
  background: #ffffff;
  border-right: 1px solid #d0d7de;
}

.preview-content {
  flex: 1;
  padding: var(--space-2);
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.theme-preview-dark .preview-header {
  background: #21262d;
}

.theme-preview-light .preview-header {
  background: #eaeef2;
}

.preview-header {
  height: 12px;
  border-radius: 2px;
  margin-bottom: var(--space-1);
}

.preview-line {
  height: 6px;
  border-radius: 2px;
}

.theme-preview-dark .preview-line {
  background: #30363d;
}

.theme-preview-light .preview-line {
  background: #d0d7de;
}

.preview-line.short {
  width: 60%;
}

.theme-name {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
}
</style>
