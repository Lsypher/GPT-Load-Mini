import { defineStore } from 'pinia'
import { ref } from 'vue'
import { keyApi } from '@/api/keys'
import type { APIKey } from '@/types'

export const useKeysStore = defineStore('keys', () => {
  const keys = ref<APIKey[]>([])
  const loading = ref(false)

  async function fetchKeys(groupId?: number) {
    loading.value = true
    try {
      const res = await keyApi.list(groupId)
      keys.value = res.data.data
    } finally {
      loading.value = false
    }
  }

  async function addKey(groupId: number, keyValue: string) {
    const res = await keyApi.add({ group_id: groupId, key_value: keyValue })
    keys.value.push(res.data.data)
    return res.data.data
  }

  async function deleteKey(id: number) {
    await keyApi.delete(id)
    keys.value = keys.value.filter((k) => k.id !== id)
  }

  async function restoreKey(id: number) {
    await keyApi.restore(id)
    const idx = keys.value.findIndex((k) => k.id === id)
    if (idx !== -1) keys.value[idx].status = 'active'
  }

  async function updateKey(id: number, data: { group_id?: number; key_value?: string }) {
    const res = await keyApi.update(id, data)
    const idx = keys.value.findIndex((k) => k.id === id)
    if (idx !== -1) keys.value[idx] = res.data.data
    return res.data.data
  }

  return { keys, loading, fetchKeys, addKey, updateKey, deleteKey, restoreKey }
})
