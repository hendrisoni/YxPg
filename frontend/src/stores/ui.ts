import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Notification } from '../types'
import * as App from '../../wailsjs/go/main/App'

export const useUiStore = defineStore('ui', () => {
  const sidebarOpen = computed(() => true)
  const sidebarWidth = ref(280)
  const historyPanelOpen = ref(false)
  const commandPaletteOpen = ref(false)
  const notifications = ref<Notification[]>([])
  const queryHistoryOpen = ref(false)

  // Layout state
  const resultPanelHeight = ref(300)
  const isResizingSidebar = ref(false)
  const isResizingResult = ref(false)

  function toggleSidebar() {
    // Disabled: sidebar is always visible
  }

  function toggleHistory() {
    historyPanelOpen.value = !historyPanelOpen.value
  }

  function toggleCommandPalette() {
    commandPaletteOpen.value = !commandPaletteOpen.value
  }

  function addNotification(notification: Omit<Notification, 'id'>) {
    const id = `notif_${Date.now()}_${Math.random().toString(36).slice(2, 7)}`
    const notif = { id, ...notification }
    notifications.value.push(notif)

    // Auto-remove after duration
    const duration = notification.duration || 5000
    if (duration > 0) {
      setTimeout(() => {
        removeNotification(id)
      }, duration)
    }
  }

  function removeNotification(id: string) {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index !== -1) {
      notifications.value.splice(index, 1)
    }
  }

  function clearNotifications() {
    notifications.value = []
  }

  // Sidebar resizing
  function setSidebarWidth(width: number) {
    sidebarWidth.value = Math.max(200, Math.min(500, width))
  }

  // Result panel resizing
  function setResultPanelHeight(height: number) {
    resultPanelHeight.value = Math.max(100, Math.min(600, height))
  }

  // Settings state
  const settings = ref({
    showFunctionsInSearch: localStorage.getItem('settings:showFunctionsInSearch') !== 'false',
    pgBinPath: '',
    dbPath: ''
  })

  async function loadSettings() {
    try {
      const path = await App.GetPgBinPath()
      settings.value.pgBinPath = path || ''
      const dbPath = await App.GetDbPath()
      settings.value.dbPath = dbPath || ''
    } catch (e) {
      console.error('Failed to load settings from config:', e)
    }
  }

  function updateSetting(key: 'showFunctionsInSearch' | 'pgBinPath' | 'dbPath', value: any) {
    if (key === 'showFunctionsInSearch') {
      settings.value.showFunctionsInSearch = value
      localStorage.setItem(`settings:${key}`, String(value))
    } else if (key === 'pgBinPath') {
      settings.value.pgBinPath = value
      App.SavePgBinPath(value).catch((e: any) => {
        console.error('Failed to save pgBinPath to config:', e)
      })
    } else if (key === 'dbPath') {
      settings.value.dbPath = value
      App.SetDbPath(value).catch((e: any) => {
        console.error('Failed to save dbPath:', e)
      })
    }
  }

  async function selectDbPath() {
    try {
      const selected = await App.SelectDbFile()
      if (selected) {
        settings.value.dbPath = selected
        await saveAndRefreshDb(selected)
      }
    } catch (e: any) {
      addNotification({
        type: 'error',
        title: 'Gagal Pilih Database',
        message: e?.message || String(e)
      })
    }
  }

  async function saveAndRefreshDb(targetPath?: string) {
    const p = targetPath || settings.value.dbPath
    if (!p) return
    try {
      await App.SetDbPath(p)
      const current = await App.GetDbPath()
      settings.value.dbPath = current
      addNotification({
        type: 'success',
        title: 'Database SQLite Berhasil Di-refresh',
        message: `Aktif: ${current}`
      })
    } catch (e: any) {
      addNotification({
        type: 'error',
        title: 'Gagal Refresh Database',
        message: e?.message || String(e)
      })
      throw e
    }
  }

  return {
    sidebarOpen,
    sidebarWidth,
    historyPanelOpen,
    commandPaletteOpen,
    notifications,
    queryHistoryOpen,
    resultPanelHeight,
    isResizingSidebar,
    isResizingResult,
    toggleSidebar,
    toggleHistory,
    toggleCommandPalette,
    addNotification,
    removeNotification,
    clearNotifications,
    setSidebarWidth,
    setResultPanelHeight,
    settings,
    updateSetting,
    loadSettings,
    selectDbPath,
    saveAndRefreshDb,
  }
}
)
