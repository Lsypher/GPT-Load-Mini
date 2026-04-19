import client from './client'
import type { Group, ApiResponse } from '@/types'

export const groupApi = {
  list: () => client.get<ApiResponse<Group[]>>('/api/groups'),

  get: (id: number) => client.get<ApiResponse<Group>>(`/api/groups/${id}`),

  create: (data: Partial<Group>) => client.post<ApiResponse<Group>>('/api/groups', data),

  update: (id: number, data: Partial<Group>) =>
    client.put<ApiResponse<Group>>(`/api/groups/${id}`, data),

  delete: (id: number) => client.delete(`/api/groups/${id}`),
}
