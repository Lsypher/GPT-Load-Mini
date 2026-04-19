import axios, { type AxiosInstance, type AxiosError, type InternalAxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import { useSettingsStore } from '@/stores/settings'

interface ApiResponse<T> {
  code: number
  data: T
  message?: string
}

const client: AxiosInstance = axios.create({
  baseURL: 'http://localhost:8080',
  timeout: 10000,
})

client.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const settings = useSettingsStore()
  if (settings.apiBaseUrl) {
    config.baseURL = settings.apiBaseUrl
  }
  if (settings.authKey && !config.url?.includes('/health')) {
    config.headers.set('X-Auth-Key', settings.authKey)
  }
  return config
})

client.interceptors.response.use(
  (response) => {
    const data = response.data as ApiResponse<unknown>
    if (data.code !== 0) {
      ElMessage.error(data.message || 'Request failed')
      return Promise.reject(new Error(data.message || 'Request failed'))
    }
    return response
  },
  (error: AxiosError<ApiResponse<unknown>>) => {
    const message = error.response?.data?.message || error.message || 'Network error'
    ElMessage.error(message)
    return Promise.reject(error)
  }
)

export default client
