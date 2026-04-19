import client from './client'

export const settingsApi = {
  reloadConfig: () => client.post('/api/admin/reload-config'),
}
