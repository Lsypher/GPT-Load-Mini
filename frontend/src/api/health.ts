import axios from 'axios'
import type { ApiResponse } from '@/types'

export const healthApi = {
  check: (baseUrl: string) =>
    axios.get<ApiResponse<{ status: string }>>(`${baseUrl}/api/health`),
}
