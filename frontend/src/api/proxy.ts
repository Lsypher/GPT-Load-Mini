import client from './client'
import type { ApiResponse } from '@/types'

export interface ProxyTestRequest {
  group_name: string
  path: string
  method?: string
  body?: string
}

export interface ProxyTestResponse {
  status: number
  headers: Record<string, string>
  body: string
}

export const proxyApi = {
  test: (data: ProxyTestRequest) =>
    client.post<ApiResponse<ProxyTestResponse>>('/api/proxy/test', data),
}
