import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

const STORAGE_KEY = 'gpt-load-mini-settings'

interface Settings {
  apiBaseUrl: string
  authKey: string
}

function loadFromStorage(): Settings {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) return JSON.parse(stored)
  } catch {}
  return { apiBaseUrl: 'http://localhost:8080', authKey: '' }
}

export const useSettingsStore = defineStore('settings', () => {
  const settings = loadFromStorage()

  const apiBaseUrl = ref(settings.apiBaseUrl)
  const authKey = ref(settings.authKey)

  function saveSettings() {
    localStorage.setItem(STORAGE_KEY, JSON.stringify({
      apiBaseUrl: apiBaseUrl.value,
      authKey: authKey.value,
    }))
  }

  watch([apiBaseUrl, authKey], saveSettings)

  function updateSettings(url: string, key: string) {
    apiBaseUrl.value = url
    authKey.value = key
  }

  return { apiBaseUrl, authKey, updateSettings }
})
