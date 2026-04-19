<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">请求日志</h1>
      <p class="page-subtitle">监控 API 请求历史和错误</p>
    </div>

    <!-- Toolbar -->
    <div class="toolbar">
      <button class="btn btn-secondary" @click="load" :disabled="loading">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="23 4 23 10 17 10"/>
          <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/>
        </svg>
        刷新
      </button>
      <button
        class="btn"
        :class="autoRefresh ? 'btn-primary' : 'btn-secondary'"
        @click="toggleAutoRefresh"
      >
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <polyline points="12 6 12 12 16 14"/>
        </svg>
        {{ autoRefresh ? '停止自动' : '自动刷新' }}
      </button>
      <div class="toolbar-spacer"></div>
      <span v-if="autoRefresh" class="auto-refresh-indicator">
        <span class="pulse"></span>
        每 10 秒自动刷新
      </span>
    </div>

    <!-- Logs Table -->
    <div class="card">
      <div class="table-container" :class="{ 'table-loading': loading }">
        <table class="table">
          <thead>
            <tr>
              <th>时间</th>
              <th>分组</th>
              <th>模型</th>
              <th>路径</th>
              <th>状态</th>
              <th>耗时</th>
              <th>来源 IP</th>
              <th>流式</th>
              <th>错误</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading && logs.length === 0">
              <td colspan="9" class="loading-cell">
                <span v-for="i in 5" :key="i" class="loading-skeleton" style="display: block; margin: 8px 0; height: 20px;"></span>
              </td>
            </tr>
            <tr v-else-if="logs.length === 0">
              <td colspan="9" class="empty-cell">
                <div class="empty-state">
                  <div class="empty-state-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                      <polyline points="14 2 14 8 20 8"/>
                      <line x1="16" y1="13" x2="8" y2="13"/>
                      <line x1="16" y1="17" x2="8" y2="17"/>
                      <polyline points="10 9 9 9 8 9"/>
                    </svg>
                  </div>
                  <div class="empty-state-title">暂无日志</div>
                  <p class="text-secondary">请求日志将显示在这里</p>
                </div>
              </td>
            </tr>
            <tr v-else v-for="log in logs" :key="log.id">
              <td class="cell-mono">{{ formatTime(log.timestamp) }}</td>
              <td>{{ log.group_name }}</td>
              <td class="cell-mono">{{ log.model }}</td>
              <td class="cell-mono cell-path" :title="log.request_path">{{ log.request_path }}</td>
              <td>
                <span class="badge" :class="log.is_success ? 'badge-success' : 'badge-danger'">
                  {{ log.status_code || 'ERR' }}
                </span>
              </td>
              <td class="cell-mono">{{ formatDuration(log.duration_ms) }}</td>
              <td class="cell-mono">{{ log.source_ip }}</td>
              <td>
                <span class="badge" :class="log.is_stream ? 'badge-info' : 'badge-muted'">
                  {{ log.is_stream ? '是' : '否' }}
                </span>
              </td>
              <td class="cell-error">
                <el-tooltip :content="log.error_message || ''" placement="top" :disabled="!log.error_message">
                  <span>{{ log.error_message || '-' }}</span>
                </el-tooltip>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Pagination -->
    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="sizes, prev, pager, next"
        @current-change="handlePageChange"
        @size-change="load"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { ElTooltip } from 'element-plus'
import AppLayout from '@/components/layout/AppLayout.vue'
import { logApi } from '@/api/logs'
import type { RequestLog } from '@/types'

const logs = ref<RequestLog[]>([])
const loading = ref(false)
const autoRefresh = ref(false)
const timer = ref<ReturnType<typeof setInterval> | null>(null)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

function formatTime(timestamp: string): string {
  if (!timestamp) return '-'
  const d = new Date(timestamp)
  return d.toLocaleString('en-US', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  })
}

function formatDuration(ms: number): string {
  if (!ms && ms !== 0) return '-'
  if (ms < 1000) return ms + 'ms'
  return (ms / 1000).toFixed(2) + 's'
}

async function load() {
  loading.value = true
  try {
    const res = await logApi.list({ page: currentPage.value, page_size: pageSize.value })
    logs.value = res.data.data.data || []
    total.value = res.data.data.total || 0
  } catch (e) {
    console.error('Failed to load logs:', e)
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  currentPage.value = page
  load()
}

function toggleAutoRefresh() {
  if (timer.value) {
    clearInterval(timer.value)
    timer.value = null
    autoRefresh.value = false
  } else {
    load()
    timer.value = setInterval(load, 10000)
    autoRefresh.value = true
  }
}

onMounted(load)

onUnmounted(() => {
  if (timer.value) clearInterval(timer.value)
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

.toolbar {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  flex-wrap: wrap;
}

.btn .icon {
  width: 16px;
  height: 16px;
}

.auto-refresh-indicator {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: 0.8125rem;
  color: var(--text-secondary);
}

.pulse {
  width: 8px;
  height: 8px;
  background: var(--success);
  border-radius: 50%;
  animation: pulse-animation 1.5s infinite;
}

@keyframes pulse-animation {
  0% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(1.2); }
  100% { opacity: 1; transform: scale(1); }
}

.card {
  overflow: hidden;
}

.table-container {
  overflow-x: auto;
  transition: opacity var(--transition-fast);
}

.table-loading {
  opacity: 0.6;
  pointer-events: none;
}

.loading-cell {
  padding: var(--space-4) !important;
}

.empty-cell {
  padding: 0 !important;
}

.cell-mono {
  font-family: var(--font-mono);
  font-size: 0.8125rem;
}

.cell-path {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cell-error {
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--danger);
  font-size: 0.8125rem;
}

.table td {
  vertical-align: middle;
}

.loading-skeleton {
  background: linear-gradient(
    90deg,
    var(--bg-tertiary) 25%,
    var(--bg-hover) 50%,
    var(--bg-tertiary) 75%
  );
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  border-radius: var(--radius-sm);
}

@keyframes skeleton-loading {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: var(--space-4) 0;
}
</style>
