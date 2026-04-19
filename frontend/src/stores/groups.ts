import { defineStore } from 'pinia'
import { ref } from 'vue'
import { groupApi } from '@/api/groups'
import type { Group } from '@/types'

export const useGroupsStore = defineStore('groups', () => {
  const groups = ref<Group[]>([])
  const loading = ref(false)

  async function fetchGroups() {
    loading.value = true
    try {
      const res = await groupApi.list()
      groups.value = res.data.data
    } finally {
      loading.value = false
    }
  }

  async function createGroup(data: Partial<Group>) {
    const res = await groupApi.create(data)
    groups.value.push(res.data.data)
    return res.data.data
  }

  async function updateGroup(id: number, data: Partial<Group>) {
    const res = await groupApi.update(id, data)
    const idx = groups.value.findIndex((g) => g.id === id)
    if (idx !== -1) groups.value[idx] = res.data.data
    return res.data.data
  }

  async function deleteGroup(id: number) {
    await groupApi.delete(id)
    groups.value = groups.value.filter((g) => g.id !== id)
  }

  return { groups, loading, fetchGroups, createGroup, updateGroup, deleteGroup }
})
