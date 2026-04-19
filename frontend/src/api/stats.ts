import client from './client'
import type { Stats, ApiResponse } from '@/types'

export const statsApi = {
  get: () => client.get<ApiResponse<Stats>>('/api/stats'),
}
