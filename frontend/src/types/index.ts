export interface Group {
  id: number
  name: string
  display_name: string
  channel_type: string
  upstream_url: string
  test_model: string
  sort: number
  proxy_api_key?: string
  created_at: string
  updated_at: string
}

export interface APIKey {
  id: number
  group_id: number
  key_value: string
  key_hash: string
  status: 'active' | 'invalid'
  failure_count: number
  last_used_at: string | null
  created_at: string
  updated_at: string
}

export interface RequestLog {
  id: string
  timestamp: string
  group_id: number
  group_name: string
  key_id: number
  model: string
  is_success: boolean
  source_ip: string
  status_code: number
  request_path: string
  duration_ms: number
  error_message: string
  is_stream: boolean
  request_type: string
}

export interface PaginatedLogs {
  data: RequestLog[]
  page: number
  page_size: number
  total: number
}

export interface Stats {
  total_keys: number
  active_keys: number
  total_requests: number
  error_rate: number
}

export interface ApiResponse<T> {
  code: number
  data: T
  message?: string
}
