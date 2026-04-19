<template>
  <div class="page">
    <div class="page-header">
      <div class="breadcrumb">
        <router-link to="/groups" class="breadcrumb-link">分组</router-link>
        <svg class="breadcrumb-sep" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="9 18 15 12 9 6"/>
        </svg>
        <span class="breadcrumb-current">{{ isEdit ? '编辑分组' : '创建分组' }}</span>
      </div>
      <h1 class="page-title">{{ isEdit ? '编辑分组' : '创建分组' }}</h1>
      <p class="page-subtitle">{{ isEdit ? '更新分组配置' : '添加新的 API 路由分组' }}</p>
    </div>

    <div class="card form-card">
      <form @submit.prevent="handleSubmit" class="form">
        <div class="form-row">
          <div class="form-group">
            <label class="label">名称</label>
            <input
              v-model="form.name"
              type="text"
              class="input"
              placeholder="例如：openai"
            />
          </div>
          <div class="form-group">
            <label class="label">显示名称</label>
            <input
              v-model="form.display_name"
              type="text"
              class="input"
              placeholder="例如：OpenAI"
            />
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label class="label">渠道类型</label>
            <input
              v-model="form.channel_type"
              type="text"
              class="input"
              placeholder="例如：openai"
            />
          </div>
          <div class="form-group">
            <label class="label">排序顺序</label>
            <input
              v-model.number="form.sort"
              type="number"
              class="input"
              min="0"
            />
          </div>
        </div>

        <div class="form-group">
          <label class="label">上游 URL</label>
          <input
            v-model="form.upstream_url"
            type="text"
            class="input"
            placeholder="例如：https://api.openai.com"
          />
        </div>

        <div class="form-group">
          <label class="label">测试模型</label>
          <input
            v-model="form.test_model"
            type="text"
            class="input"
            placeholder="例如：gpt-3.5-turbo"
          />
        </div>

        <div class="form-group">
          <label class="label">代理密钥</label>
          <div class="input-with-button">
            <input
              v-model="form.proxy_api_key"
              type="text"
              class="input"
              placeholder="下游用户调用时需要携带此 Key，留空则无需认证"
            />
            <button type="button" class="btn btn-secondary" @click="generateApiKey">生成</button>
          </div>
        </div>

        <div class="form-actions form-actions-main">
          <button type="button" class="btn btn-secondary" @click="$router.back()">
            取消
          </button>
          <button type="submit" class="btn btn-primary" :disabled="saving">
            {{ saving ? '保存中...' : (isEdit ? '更新' : '创建') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useGroupsStore } from '@/stores/groups'
import { groupApi } from '@/api/groups'

const router = useRouter()
const route = useRoute()
const groupsStore = useGroupsStore()

const isEdit = computed(() => !!route.params.id)
const groupId = computed(() => Number(route.params.id))
const saving = ref(false)

const form = reactive({
  name: '',
  display_name: '',
  channel_type: '',
  upstream_url: '',
  test_model: '',
  sort: 0,
  proxy_api_key: '',
})

function generateApiKey() {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  for (let i = 0; i < 32; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  form.proxy_api_key = result
}

async function handleSubmit() {
  if (!form.name || !form.display_name) {
    ElMessage.warning('名称和显示名称为必填项')
    return
  }

  saving.value = true
  try {
    if (isEdit.value) {
      await groupsStore.updateGroup(groupId.value, form)
      ElMessage.success('分组已更新')
    } else {
      await groupsStore.createGroup(form)
      ElMessage.success('分组已创建')
    }
    router.push('/groups')
  } catch {
    // error handled by api client
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  if (isEdit.value) {
    const res = await groupApi.get(groupId.value)
    const g = res.data.data
    Object.assign(form, {
      name: g.name,
      display_name: g.display_name,
      channel_type: g.channel_type,
      upstream_url: g.upstream_url,
      test_model: g.test_model,
      sort: g.sort,
      proxy_api_key: g.proxy_api_key || '',
    })
  }
})
</script>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
  max-width: 700px;
}

.page-header {
  margin-bottom: var(--space-2);
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: 0.875rem;
  margin-bottom: var(--space-3);
}

.breadcrumb-link {
  color: var(--text-secondary);
  text-decoration: none;
}

.breadcrumb-link:hover {
  color: var(--accent);
}

.breadcrumb-sep {
  width: 14px;
  height: 14px;
  color: var(--text-tertiary);
}

.breadcrumb-current {
  color: var(--text-primary);
  font-weight: 500;
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

.form-card {
  padding: var(--space-6);
}

.form {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-4);
}

@media (max-width: 600px) {
  .form-row {
    grid-template-columns: 1fr;
  }
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  padding-top: var(--space-4);
  border-top: 1px solid var(--border-subtle);
}

.form-actions-main {
  margin-top: var(--space-4);
}

.input {
  width: 100%;
}

.input-with-button {
  display: flex;
  gap: var(--space-2);
}

.input-with-button .input {
  flex: 1;
}
</style>
