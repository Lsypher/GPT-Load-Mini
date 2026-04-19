<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">API 密钥</h1>
      <p class="page-subtitle">管理您的 API 密钥以进行身份验证</p>
    </div>

    <!-- Toolbar -->
    <div class="toolbar">
      <select v-model="filterGroupId" class="select" @change="loadKeys">
        <option :value="undefined">全部分组</option>
        <option v-for="g in groupsStore.groups" :key="g.id" :value="g.id">
          {{ g.display_name }}
        </option>
      </select>
      <div class="toolbar-spacer"></div>
      <button class="btn btn-secondary" @click="handleExport">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
          <polyline points="7 10 12 15 17 10"/>
          <line x1="12" y1="15" x2="12" y2="3"/>
        </svg>
        导出
      </button>
      <button class="btn btn-secondary" @click="showImportDialog = true">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
          <polyline points="17 8 12 3 7 8"/>
          <line x1="12" y1="3" x2="12" y2="15"/>
        </svg>
        导入
      </button>
      <button class="btn btn-primary" @click="showAddDialog = true">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="5" x2="12" y2="19"/>
          <line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        添加密钥
      </button>
    </div>

    <!-- Keys Table -->
    <div class="card">
      <div class="table-container">
        <table class="table">
          <thead>
            <tr>
              <th>ID</th>
              <th>分组</th>
              <th>密钥哈希</th>
              <th>状态</th>
              <th>失败次数</th>
              <th>最后使用</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="keysStore.loading">
              <td colspan="7" class="loading-cell">
                <span class="loading-skeleton" style="width: 100%; height: 20px;"></span>
              </td>
            </tr>
            <tr v-else-if="keysStore.keys.length === 0">
              <td colspan="7" class="empty-cell">
                <div class="empty-state">
                  <div class="empty-state-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                      <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
                    </svg>
                  </div>
                  <div class="empty-state-title">暂无密钥</div>
                  <p class="text-secondary">创建您的第一个 API 密钥</p>
                </div>
              </td>
            </tr>
            <tr v-else v-for="key in keysStore.keys" :key="key.id">
              <td class="cell-mono">{{ key.id }}</td>
              <td>{{ getGroupName(key.group_id) }}</td>
              <td class="cell-mono">{{ maskHash(key.key_hash) }}</td>
              <td>
                <span class="badge" :class="key.status === 'active' ? 'badge-success' : 'badge-danger'">
                  {{ key.status }}
                </span>
              </td>
              <td class="cell-mono">{{ key.failure_count }}</td>
              <td class="cell-mono">{{ formatDate(key.last_used_at) }}</td>
              <td>
                <div class="action-buttons">
                  <button
                    class="btn btn-ghost btn-sm"
                    @click="openEditDialog(key)"
                  >
                    编辑
                  </button>
                  <button
                    v-if="key.status !== 'active'"
                    class="btn btn-ghost btn-sm"
                    @click="handleRestore(key.id)"
                  >
                    恢复
                  </button>
                  <button
                    class="btn btn-ghost btn-sm text-danger"
                    @click="handleDelete(key.id)"
                  >
                    删除
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Add Dialog -->
    <div v-if="showAddDialog" class="dialog-overlay" @click.self="showAddDialog = false">
      <div class="dialog">
        <div class="dialog-header">
          <h3 class="dialog-title">添加 API 密钥</h3>
          <button class="btn btn-ghost btn-icon" @click="showAddDialog = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="dialog-body">
          <div class="form-group">
            <label class="label">分组</label>
            <select v-model="addForm.group_id" class="select" style="appearance: auto;">
              <option :value="0" disabled>选择分组</option>
              <option v-for="g in groupsStore.groups" :key="g.id" :value="g.id">
                {{ g.display_name }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label class="label">密钥值</label>
            <input
              v-model="addForm.key_value"
              type="text"
              class="input"
              placeholder="sk-..."
            />
          </div>
        </div>
        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="showAddDialog = false">取消</button>
          <button class="btn btn-primary" @click="handleAdd" :disabled="adding">
            {{ adding ? '添加中...' : '添加' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Edit Dialog -->
    <div v-if="showEditDialog" class="dialog-overlay" @click.self="showEditDialog = false">
      <div class="dialog">
        <div class="dialog-header">
          <h3 class="dialog-title">编辑 API 密钥</h3>
          <button class="btn btn-ghost btn-icon" @click="showEditDialog = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="dialog-body">
          <div class="form-group">
            <label class="label">Group</label>
            <select v-model="editForm.group_id" class="select" style="appearance: auto;">
              <option v-for="g in groupsStore.groups" :key="g.id" :value="g.id">
                {{ g.display_name }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label class="label">新密钥值（留空则保持不变）</label>
            <input
              v-model="editForm.key_value"
              type="text"
              class="input"
              placeholder="sk-..."
            />
          </div>
        </div>
        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="showEditDialog = false">取消</button>
          <button class="btn btn-primary" @click="handleEdit" :disabled="editing">
            {{ editing ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Import Dialog -->
    <div v-if="showImportDialog" class="dialog-overlay" @click.self="showImportDialog = false">
      <div class="dialog">
        <div class="dialog-header">
          <h3 class="dialog-title">导入 API 密钥</h3>
          <button class="btn btn-ghost btn-icon" @click="showImportDialog = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="dialog-body">
          <div class="form-group">
            <label class="label">Group</label>
            <select v-model="importForm.group_id" class="select" style="appearance: auto;">
              <option v-for="g in groupsStore.groups" :key="g.id" :value="g.id">
                {{ g.display_name }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label class="label">API 密钥（每行一个）</label>
            <textarea
              v-model="importForm.keys_text"
              class="input textarea"
              placeholder="sk-...&#10;sk-...&#10;sk-..."
              rows="6"
            />
          </div>
        </div>
        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="showImportDialog = false">取消</button>
          <button class="btn btn-primary" @click="handleImport" :disabled="importing">
            {{ importing ? '导入中...' : '导入' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import AppLayout from '@/components/layout/AppLayout.vue'
import { useKeysStore } from '@/stores/keys'
import { useGroupsStore } from '@/stores/groups'
import { maskHash, formatDate } from '@/utils/formatters'
import { keyApi } from '@/api/keys'

const keysStore = useKeysStore()
const groupsStore = useGroupsStore()

const filterGroupId = ref<number | undefined>()
const showAddDialog = ref(false)
const adding = ref(false)
const addForm = reactive({ group_id: 0, key_value: '' })

const showEditDialog = ref(false)
const editingId = ref<number | null>(null)
const editing = ref(false)
const editForm = reactive({ group_id: 0, key_value: '' })

const showImportDialog = ref(false)
const importing = ref(false)
const importForm = reactive({ group_id: 0, keys_text: '' })

function getGroupName(id: number) {
  return groupsStore.groups.find((g) => g.id === id)?.display_name || '-'
}

async function loadKeys() {
  await keysStore.fetchKeys(filterGroupId.value)
}

async function handleAdd() {
  if (!addForm.group_id || !addForm.key_value) {
    ElMessage.warning('请填写所有字段')
    return
  }
  adding.value = true
  try {
    await keysStore.addKey(addForm.group_id, addForm.key_value)
    ElMessage.success('密钥已添加')
    showAddDialog.value = false
    addForm.group_id = 0
    addForm.key_value = ''
  } catch {} finally {
    adding.value = false
  }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('确定要删除此密钥吗？', '警告', { type: 'warning' })
    await keysStore.deleteKey(id)
  } catch {}
}

async function handleRestore(id: number) {
  try {
    await keysStore.restoreKey(id)
    ElMessage.success('密钥已恢复')
  } catch {}
}

function openEditDialog(key: any) {
  editingId.value = key.id
  editForm.group_id = key.group_id
  editForm.key_value = ''
  showEditDialog.value = true
}

async function handleEdit() {
  if (!editingId.value) return
  editing.value = true
  try {
    const data: { group_id?: number; key_value?: string } = {}
    if (editForm.group_id) data.group_id = editForm.group_id
    if (editForm.key_value) data.key_value = editForm.key_value
    await keysStore.updateKey(editingId.value, data)
    ElMessage.success('密钥已更新')
    showEditDialog.value = false
  } catch {} finally {
    editing.value = false
  }
}

async function handleExport() {
  try {
    const res = await keyApi.export(filterGroupId.value)
    const blob = new Blob([JSON.stringify(res.data.data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `api-keys-${new Date().toISOString().slice(0, 10)}.json`
    a.click()
    URL.revokeObjectURL(url)
    ElMessage.success('密钥已导出')
  } catch {}
}

async function handleImport() {
  if (!importForm.group_id || !importForm.keys_text) {
    ElMessage.warning('请填写所有字段')
    return
  }
  const lines = importForm.keys_text.split('\n').filter(k => k.trim())
  if (lines.length === 0) {
    ElMessage.warning('没有可导入的密钥')
    return
  }
  importing.value = true
  try {
    const keys = lines.map(k => ({ group_id: importForm.group_id, key_value: k.trim() }))
    const res = await keyApi.import(keys)
    ElMessage.success(`已导入 ${res.data.data?.imported || 0} 个密钥，失败 ${res.data.data?.failed || 0} 个`)
    showImportDialog.value = false
    importForm.group_id = 0
    importForm.keys_text = ''
    await keysStore.fetchKeys()
  } catch {} finally {
    importing.value = false
  }
}

onMounted(() => {
  groupsStore.fetchGroups()
  keysStore.fetchKeys()
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

.toolbar .select {
  width: auto;
  min-width: 160px;
}

.btn .icon {
  width: 16px;
  height: 16px;
}

.btn-icon svg {
  width: 18px;
  height: 18px;
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

.empty-state {
  padding: var(--space-10);
}

.cell-mono {
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

/* Dialog */
.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: var(--space-4);
}

.dialog {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  width: 100%;
  max-width: 420px;
  box-shadow: var(--shadow);
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4) var(--space-5);
  border-bottom: 1px solid var(--border-subtle);
}

.dialog-title {
  font-size: 1rem;
  font-weight: 600;
}

.dialog-body {
  padding: var(--space-5);
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  padding: var(--space-4) var(--space-5);
  border-top: 1px solid var(--border-subtle);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.input, .select {
  width: 100%;
}

.textarea {
  resize: vertical;
  min-height: 100px;
  font-family: var(--font-mono);
  font-size: 0.8125rem;
}
</style>
