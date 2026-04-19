<template>
  <div class="dashboard">
    <!-- Stats Grid -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-card-header">
          <div class="stat-card-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
            </svg>
          </div>
        </div>
        <div class="stat-card-label">密钥总数</div>
        <div class="stat-card-value">{{ stats.total_keys }}</div>
      </div>

      <div class="stat-card">
        <div class="stat-card-header">
          <div class="stat-card-icon success">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
              <polyline points="22 4 12 14.01 9 11.01"/>
            </svg>
          </div>
        </div>
        <div class="stat-card-label">活跃密钥</div>
        <div class="stat-card-value success">{{ stats.active_keys }}</div>
      </div>

      <div class="stat-card">
        <div class="stat-card-header">
          <div class="stat-card-icon info">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/>
            </svg>
          </div>
        </div>
        <div class="stat-card-label">请求总数</div>
        <div class="stat-card-value info">{{ formatNumber(stats.total_requests) }}</div>
      </div>

      <div class="stat-card">
        <div class="stat-card-header">
          <div class="stat-card-icon" :class="errorRateClass">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
          </div>
        </div>
        <div class="stat-card-label">错误率</div>
        <div class="stat-card-value" :class="errorRateClass">{{ formatErrorRate(stats.error_rate) }}</div>
      </div>
    </div>

    <!-- Recent Logs -->
    <div class="card">
      <div class="card-header">
        <span class="card-title">最近日志</span>
        <router-link to="/logs" class="btn btn-ghost btn-sm">
          查看全部
          <svg class="icon-sm" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M5 12h14M12 5l7 7-7 7"/>
          </svg>
        </router-link>
      </div>
      <div class="table-container">
        <table class="table">
          <thead>
            <tr>
              <th>时间</th>
              <th>分组</th>
              <th>模型</th>
              <th>状态</th>
              <th>耗时</th>
              <th>来源 IP</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="recentLogs.length === 0">
              <td colspan="6" class="empty-cell">暂无日志</td>
            </tr>
            <tr v-for="log in recentLogs" :key="log.id">
              <td class="cell-mono">{{ formatTime(log.timestamp) }}</td>
              <td>{{ log.group_name }}</td>
              <td class="cell-mono">{{ log.model }}</td>
              <td>
                <span class="badge" :class="log.is_success ? 'badge-success' : 'badge-danger'">
                  {{ log.status_code || 'ERR' }}
                </span>
              </td>
              <td class="cell-mono">{{ formatDuration(log.duration_ms) }}</td>
              <td class="cell-mono">{{ log.source_ip }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { statsApi } from '@/api/stats'
import { logApi } from '@/api/logs'
import type { Stats, RequestLog } from '@/types'

const stats = ref<Stats>({ total_keys: 0, active_keys: 0, total_requests: 0, error_rate: 0 })
const recentLogs = ref<RequestLog[]>([])

const errorRateClass = computed(() => {
  if (stats.value.error_rate > 0.05) return 'danger'
  if (stats.value.error_rate > 0.02) return 'warning'
  return ''
})

function formatNumber(num: number): string {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toString()
}

function formatErrorRate(rate: number): string {
  return (rate * 100).toFixed(2) + '%'
}

function formatTime(timestamp: string): string {
  if (!timestamp) return '-'
  const d = new Date(timestamp)
  return d.toLocaleTimeString('en-US', { hour12: false })
}

function formatDuration(ms: number): string {
  if (!ms) return '-'
  if (ms < 1000) return ms + 'ms'
  return (ms / 1000).toFixed(2) + 's'
}

async function load() {
  try {
    const [statsRes, logsRes] = await Promise.all([
      statsApi.get(),
      logApi.list()
    ])
    stats.value = statsRes.data.data
    const logsData = logsRes.data.data?.data ?? logsRes.data.data ?? []
    recentLogs.value = Array.isArray(logsData) ? logsData.slice(0, 10) : []
  } catch (e) {
    console.error('Failed to load dashboard data:', e)
  }
}

let timer: ReturnType<typeof setInterval>

onMounted(() => {
  load()
  timer = setInterval(load, 30000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--space-4);
}

@media (max-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}

.stat-card-icon {
  background: var(--accent-muted);
  color: var(--accent);
}

.stat-card-icon.success {
  background: var(--success-muted);
  color: var(--success);
}

.stat-card-icon.info {
  background: var(--info-muted);
  color: var(--info);
}

.stat-card-icon.danger {
  background: var(--danger-muted);
  color: var(--danger);
}

.stat-card-icon.warning {
  background: var(--warning-muted);
  color: var(--warning);
}

.stat-card-icon svg {
  width: 20px;
  height: 20px;
}

.stat-card-value.danger {
  color: var(--danger);
}

.stat-card-value.warning {
  color: var(--warning);
}

.card-title {
  font-size: 0.9375rem;
  font-weight: 600;
}

.btn-ghost {
  color: var(--text-secondary);
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
}

.btn-ghost:hover {
  color: var(--accent);
}

.icon-sm {
  width: 14px;
  height: 14px;
}

.empty-cell {
  text-align: center;
  color: var(--text-tertiary);
  padding: var(--space-8) !important;
}

.table td {
  vertical-align: middle;
}
</style>
