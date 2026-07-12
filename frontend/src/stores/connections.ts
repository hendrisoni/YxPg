import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Connection, ConnectionTestResult, SchemaInfo, TableInfo, ColumnInfo } from '../types'
import * as App from '../../wailsjs/go/main/App'

// Wails runtime bindings declaration
declare global {
  interface Window {
    Go: any
  }
}

export const useConnectionsStore = defineStore('connections', () => {
  const connections = ref<Connection[]>([])
  const activeConnections = ref<string[]>([])
  const currentConnectionId = ref<string | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const hasAutoConnected = ref(false)

  const currentConnection = computed(() => {
    if (!currentConnectionId.value) return null
    return connections.value.find(c => c.id === currentConnectionId.value) || null
  })

  const isConnected = computed(() => {
    return currentConnectionId.value !== null
  })

  async function loadConnections() {
    loading.value = true
    try {
      connections.value = await App.ListConnections()
      activeConnections.value = await App.GetActiveConnections()
      if (!currentConnectionId.value && activeConnections.value.length > 0) {
        currentConnectionId.value = activeConnections.value[0]
      }
    } catch (e: any) {
      error.value = e.message || String(e)
    } finally {
      loading.value = false
    }
  }

  async function addConnection(conn: Connection) {
    loading.value = true
    try {
      await App.AddConnection(conn as any)
      await loadConnections()
    } catch (e: any) {
      error.value = e.message || String(e)
      throw e
    } finally {
      loading.value = false
    }
  }

  async function updateConnection(conn: Connection) {
    loading.value = true
    try {
      await App.UpdateConnection(conn as any)
      await loadConnections()
    } catch (e: any) {
      error.value = e.message || String(e)
      throw e
    } finally {
      loading.value = false
    }
  }

  async function testConnection(conn: Connection): Promise<ConnectionTestResult> {
    try {
      return await App.TestConnection(conn as any)
    } catch (e: any) {
      return { ok: false, latency_ms: 0, message: e.message || String(e) }
    }
  }

  async function deleteConnection(id: string) {
    loading.value = true
    try {
      await App.DeleteConnection(id)
      if (currentConnectionId.value === id) {
        currentConnectionId.value = null
      }
      await loadConnections()
    } catch (e: any) {
      error.value = e.message || String(e)
    } finally {
      loading.value = false
    }
  }

  async function connect(id: string) {
    loading.value = true
    error.value = null
    try {
      await App.Connect(id)
      currentConnectionId.value = id
      activeConnections.value = await App.GetActiveConnections()
    } catch (e: any) {
      error.value = e.message || String(e)
      throw e
    } finally {
      loading.value = false
    }
  }

  async function disconnect(id: string) {
    try {
      await App.Disconnect(id)
      if (currentConnectionId.value === id) {
        currentConnectionId.value = null
      }
      activeConnections.value = await App.GetActiveConnections()
    } catch (e: any) {
      error.value = e.message || String(e)
    }
  }

  async function syncServerConnections() {
    loading.value = true
    error.value = null
    try {
      const count = await (App as any).SyncServerConnections()
      await loadConnections()
      return count
    } catch (e: any) {
      error.value = e.message || String(e)
      throw e
    } finally {
      loading.value = false
    }
  }

  function getWailsBindings() {
    return App
  }

  return {
    connections,
    activeConnections,
    currentConnectionId,
    currentConnection,
    isConnected,
    loading,
    error,
    hasAutoConnected,
    loadConnections,
    addConnection,
    updateConnection,
    testConnection,
    deleteConnection,
    connect,
    disconnect,
    syncServerConnections,
    getWailsBindings,
  }
})
