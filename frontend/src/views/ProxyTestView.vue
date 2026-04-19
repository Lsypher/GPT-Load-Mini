<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">代理测试</h1>
      <p class="page-subtitle">测试代理请求并完整检查请求/响应</p>
    </div>

    <div class="test-container">
      <!-- Request Panel -->
      <div class="card request-card">
        <h3 class="card-title">请求</h3>

        <div class="form-group">
          <label class="label">分组</label>
          <el-select v-model="form.groupName" placeholder="选择分组" class="select-full">
            <el-option
              v-for="group in groups"
              :key="group.name"
              :label="group.display_name || group.name"
              :value="group.name"
            />
          </el-select>
        </div>

        <div class="form-row">
          <div class="form-group" style="flex: 2">
            <label class="label">路径</label>
            <input
              v-model="form.path"
              type="text"
              class="input"
              placeholder="v1/chat/completions"
            />
          </div>
          <div class="form-group" style="flex: 1">
            <label class="label">方法</label>
            <el-select v-model="form.method" class="select-full">
              <el-option label="POST" value="POST" />
              <el-option label="GET" value="GET" />
              <el-option label="PUT" value="PUT" />
              <el-option label="DELETE" value="DELETE" />
            </el-select>
          </div>
        </div>

        <div class="form-group">
          <label class="label">请求体 (JSON)</label>
          <textarea
            v-model="form.body"
            class="textarea"
            rows="8"
            placeholder='{"model":"gpt-3.5-turbo","messages":[{"role":"user","content":"Hello"}]}'
          ></textarea>
        </div>

        <div class="form-actions">
          <button class="btn btn-secondary" @click="loadGroups" :disabled="loadingGroups">
            刷新分组
          </button>
          <button class="btn btn-primary" @click="handleSend" :disabled="sending || !form.groupName || !form.path">
            {{ sending ? '发送中...' : '发送请求' }}
          </button>
        </div>
      </div>

      <!-- Response Panel -->
      <div class="card response-card" v-if="response">
        <div class="response-header">
          <h3 class="card-title">响应</h3>
          <span class="status-badge" :class="response.status >= 200 && response.status < 300 ? 'status-success' : 'status-error'">
            {{ response.status }} {{ response.status === 200 ? '成功' : response.status === 400 ? '错误请求' : response.status === 401 ? '未授权' : response.status === 404 ? '未找到' : response.status === 502 ? '网关错误' : response.status === 503 ? '服务不可用' : '' }}
          </span>
        </div>

        <!-- Response Headers -->
        <div class="collapsible" :class="{ collapsed: !showHeaders }">
          <button class="collapsible-header" @click="showHeaders = !showHeaders">
            <span>响应头 ({{ Object.keys(response.headers).length }})</span>
            <svg class="chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="6 9 12 15 18 9"/>
            </svg>
          </button>
          <div class="collapsible-content" v-show="showHeaders">
            <div class="headers-list">
              <div v-for="(value, key) in response.headers" :key="key" class="header-item">
                <span class="header-key">{{ key }}:</span>
                <span class="header-value">{{ value }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Response Body -->
        <div class="form-group">
          <label class="label">响应体</label>
          <pre class="response-body" :class="{ 'response-error': response.status >= 400 }">{{ formatBody(response.body) }}</pre>
        </div>
      </div>

      <!-- Error Panel -->
      <div class="card error-card" v-if="error">
        <h3 class="card-title error-title">错误</h3>
        <pre class="error-body">{{ error }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElSelect, ElOption } from 'element-plus'
import { useGroupsStore } from '@/stores/groups'
import { proxyApi, type ProxyTestResponse } from '@/api/proxy'

const groupsStore = useGroupsStore()
const groups = ref<typeof groupsStore.groups>([])
const loadingGroups = ref(false)
const sending = ref(false)
const response = ref<ProxyTestResponse | null>(null)
const error = ref('')
const showHeaders = ref(false)

const form = reactive({
  groupName: '',
  path: 'v1/chat/completions',
  method: 'POST',
  body: '{\n  "model": "gpt-3.5-turbo",\n  "messages": [\n    {\n      "role": "user",\n      "content": "Hello"\n    }\n  ]\n}',
})

async function loadGroups() {
  loadingGroups.value = true
  try {
    await groupsStore.fetchGroups()
    groups.value = groupsStore.groups
  } catch (e) {
    ElMessage.error('加载分组失败')
  } finally {
    loadingGroups.value = false
  }
}

async function handleSend() {
  if (!form.groupName || !form.path) {
    ElMessage.warning('分组和路径为必填项')
    return
  }

  sending.value = true
  error.value = ''
  response.value = null

  try {
    const res = await proxyApi.test({
      group_name: form.groupName,
      path: form.path,
      method: form.method,
      body: form.body || undefined,
    })
    response.value = res.data.data
  } catch (e: any) {
    error.value = e.message || '请求失败'
  } finally {
    sending.value = false
  }
}

function formatBody(body: string): string {
  if (!body) return ''
  try {
    return JSON.stringify(JSON.parse(body), null, 2)
  } catch {
    return body
  }
}

onMounted(loadGroups)
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

.test-container {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.card {
  padding: var(--space-5);
}

.card-title {
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: var(--space-4);
  color: var(--text-primary);
}

.request-card {
  max-width: 900px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  margin-bottom: var(--space-4);
}

.form-row {
  display: flex;
  gap: var(--space-4);
}

.select-full {
  width: 100%;
}

.input, .textarea {
  width: 100%;
}

.textarea {
  resize: vertical;
  min-height: 120px;
  font-family: var(--font-mono);
  font-size: 0.875rem;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  padding-top: var(--space-4);
  border-top: 1px solid var(--border-subtle);
}

.response-card {
  max-width: 900px;
}

.response-header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-4);
}

.response-header .card-title {
  margin-bottom: 0;
}

.status-badge {
  padding: var(--space-1) var(--space-3);
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  font-weight: 500;
  font-family: var(--font-mono);
}

.status-success {
  background: var(--success);
  color: white;
}

.status-error {
  background: var(--danger);
  color: white;
}

.collapsible {
  margin-bottom: var(--space-4);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.collapsible-header {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-3) var(--space-4);
  background: var(--bg-tertiary);
  border: none;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
}

.collapsible-header:hover {
  background: var(--bg-hover);
}

.chevron {
  width: 16px;
  height: 16px;
  transition: transform var(--transition-fast);
}

.collapsible.collapsed .chevron {
  transform: rotate(-90deg);
}

.collapsible-content {
  border-top: 1px solid var(--border);
}

.headers-list {
  padding: var(--space-3) var(--space-4);
  max-height: 200px;
  overflow-y: auto;
}

.header-item {
  display: flex;
  gap: var(--space-2);
  font-size: 0.8125rem;
  font-family: var(--font-mono);
  padding: var(--space-1) 0;
}

.header-key {
  color: var(--accent);
  font-weight: 500;
}

.header-value {
  color: var(--text-secondary);
  word-break: break-all;
}

.response-body {
  background: var(--bg-tertiary);
  padding: var(--space-4);
  border-radius: var(--radius-md);
  font-family: var(--font-mono);
  font-size: 0.8125rem;
  overflow-x: auto;
  max-height: 400px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
}

.response-error {
  color: var(--danger);
}

.error-card {
  max-width: 900px;
  border: 1px solid var(--danger);
}

.error-title {
  color: var(--danger);
}

.error-body {
  background: var(--bg-tertiary);
  padding: var(--space-4);
  border-radius: var(--radius-md);
  font-family: var(--font-mono);
  font-size: 0.875rem;
  color: var(--danger);
  white-space: pre-wrap;
  margin: 0;
}
</style>
