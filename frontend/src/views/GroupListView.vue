<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">分组管理</h1>
      <p class="page-subtitle">配置 API 路由分组和上游设置</p>
    </div>

    <!-- Toolbar -->
    <div class="toolbar">
      <div class="toolbar-spacer"></div>
      <router-link to="/groups/new" class="btn btn-primary">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="5" x2="12" y2="19"/>
          <line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        创建分组
      </router-link>
    </div>

    <!-- Groups Table -->
    <div class="card">
      <div class="table-container">
        <table class="table">
          <thead>
            <tr>
              <th>ID</th>
              <th>名称</th>
              <th>显示名称</th>
              <th>渠道</th>
              <th>上游 URL</th>
              <th>测试模型</th>
              <th>排序</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="groupsStore.loading">
              <td colspan="8" class="loading-cell">
                <span v-for="i in 3" :key="i" class="loading-skeleton" style="display: block; margin: 8px 0; height: 20px;"></span>
              </td>
            </tr>
            <tr v-else-if="groupsStore.groups.length === 0">
              <td colspan="8" class="empty-cell">
                <div class="empty-state">
                  <div class="empty-state-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                      <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
                    </svg>
                  </div>
                  <div class="empty-state-title">暂无分组</div>
                  <p class="text-secondary">创建您的第一个分组</p>
                </div>
              </td>
            </tr>
            <tr v-else v-for="group in groupsStore.groups" :key="group.id">
              <td class="cell-mono">{{ group.id }}</td>
              <td class="cell-mono">{{ group.name }}</td>
              <td>{{ group.display_name }}</td>
              <td class="cell-mono">{{ group.channel_type }}</td>
              <td class="cell-url" :title="group.upstream_url">{{ group.upstream_url }}</td>
              <td class="cell-mono">{{ group.test_model }}</td>
              <td class="cell-mono">{{ group.sort }}</td>
              <td>
                <div class="action-buttons">
                  <router-link :to="`/groups/${group.id}/edit`" class="btn btn-ghost btn-sm">
                    编辑
                  </router-link>
                  <button class="btn btn-ghost btn-sm text-danger" @click="handleDelete(group.id)">
                    删除
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { ElMessageBox } from 'element-plus'
import { useGroupsStore } from '@/stores/groups'

const groupsStore = useGroupsStore()

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('确定要删除此分组吗？', '警告', { type: 'warning' })
    await groupsStore.deleteGroup(id)
  } catch {}
}

onMounted(() => groupsStore.fetchGroups())
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

.card {
  overflow: hidden;
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

.cell-url {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: var(--font-mono);
  font-size: 0.8125rem;
}

.action-buttons {
  display: flex;
  gap: var(--space-1);
}

.text-danger {
  color: var(--danger);
}
</style>
