import client from './client'
import type { APIKey, ApiResponse } from '@/types'

export const keyApi = {
  list: (groupId?: number) => {
    const params = groupId ? { group_id: groupId } : {}
    return client.get<ApiResponse<APIKey[]>>('/api/keys', { params })
  },

  add: (data: { group_id: number; key_value: string }) =>
    client.post<ApiResponse<APIKey>>('/api/keys', data),

  update: (id: number, data: { group_id?: number; key_value?: string }) =>
    client.put<ApiResponse<APIKey>>(`/api/keys/${id}`, data),

  delete: (id: number) => client.delete(`/api/keys/${id}`),

  restore: (id: number) => client.post(`/api/keys/${id}/restore`),

  export: (groupId?: number) => {
    const params = groupId ? { group_id: groupId } : {}
    return client.get<ApiResponse<{ group_id: number; key_hash: string; status: string; created_at: string }[]>>('/api/keys/export', { params })
  },

  import: (keys: { group_id: number; key_value: string }[]) =>
    client.post<ApiResponse<{ imported: number; failed: number }>>('/api/keys/import', keys),
}
