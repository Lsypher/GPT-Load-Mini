import client from './client'
import type { PaginatedLogs, ApiResponse } from '@/types'

export const logApi = {
  list: (params?: { page?: number; page_size?: number }) =>
    client.get<ApiResponse<PaginatedLogs>>('/api/logs', { params }),
}
