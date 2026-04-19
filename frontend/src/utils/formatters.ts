export function formatDate(dateStr: string | null): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleString('zh-CN')
}

export function formatDuration(ms: number): string {
  if (ms < 1000) return `${ms}ms`
  return `${(ms / 1000).toFixed(2)}s`
}

export function formatErrorRate(rate: number): string {
  return `${(rate * 100).toFixed(1)}%`
}

export function maskKey(key: string): string {
  if (!key || key.length < 8) return '***'
  return key.slice(0, 4) + '***' + key.slice(-4)
}

export function maskHash(hash: string): string {
  if (!hash || hash.length < 12) return '***'
  return hash.slice(0, 6) + '...' + hash.slice(-4)
}
