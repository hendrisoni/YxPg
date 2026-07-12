<template>
  <div class="h-full w-full flex flex-col items-center justify-center bg-navy-primary overflow-y-auto p-6">
    <!-- Logo -->
    <div class="flex items-center gap-3 mb-8">
      <svg class="w-12 h-12 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <ellipse cx="12" cy="5" rx="9" ry="3" />
        <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3" />
        <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5" />
      </svg>
      <div>
        <h1 class="text-3xl font-bold text-text-primary">YxPg</h1>
        <p class="text-sm text-text-secondary">PostgreSQL Desktop GUI</p>
      </div>
    </div>

    <!-- Connection cards -->
    <div class="w-full max-w-lg px-4">
      <div class="flex items-center justify-between mb-3">
        <h2 class="text-sm font-semibold text-text-primary">
          Connections
        </h2>
        <div class="flex items-center gap-1.5">
          <button
            @click="handleSync"
            :disabled="syncing"
            class="flex items-center gap-1 px-2.5 py-1 text-[11px] bg-navy-secondary border border-navy-border text-text-primary rounded-md hover:bg-navy-hover disabled:opacity-50 transition-colors"
          >
            <svg
              class="w-3 h-3"
              :class="{ 'animate-spin': syncing }"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
              <path d="M3 3v5h5" />
              <path d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16" />
              <path d="M21 21v-5h-5" />
            </svg>
            {{ syncing ? 'Syncing...' : 'Sync' }}
          </button>
          <button
            @click="handleSyncPgAdmin"
            :disabled="syncingPgAdmin"
            class="flex items-center gap-1 px-2.5 py-1 text-[11px] bg-navy-secondary border border-navy-border text-text-primary rounded-md hover:bg-navy-hover disabled:opacity-50 transition-colors"
          >
            <svg
              class="w-3 h-3"
              :class="{ 'animate-spin': syncingPgAdmin }"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67" />
            </svg>
            {{ syncingPgAdmin ? 'Syncing pgAdmin...' : 'Sync pgAdmin' }}
          </button>
          <button
            @click="handleNewConnection"
            class="flex items-center gap-1 px-2.5 py-1 text-xs bg-teal-accent text-navy-primary rounded-md font-medium hover:bg-teal-hover transition-colors"
          >
            <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 5v14M5 12h14" />
            </svg>
            New Connection
          </button>
        </div>
      </div>

      <!-- Connection list -->
      <div v-if="connectionsStore.connections.length > 0" class="space-y-1.5">
        <div
          v-for="conn in connectionsStore.connections"
          :key="conn.id"
          @click="handleConnect(conn)"
          class="group flex items-center gap-2.5 py-1.5 px-2.5 bg-navy-secondary border border-navy-border rounded-md cursor-pointer hover:border-teal-accent/50 hover:bg-navy-hover transition-all"
        >
          <div class="w-2 h-2 rounded-full flex-shrink-0" :style="{ backgroundColor: conn.color || '#00C9A7' }"></div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-1.5">
              <span class="text-xs font-semibold text-text-primary">{{ conn.name }}</span>
              <svg v-if="defaultConnectionId === conn.id" class="w-3.5 h-3.5 text-accent-amber fill-current flex-shrink-0" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2" title="Default Connection">
                <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
              </svg>
            </div>
            <div class="text-[10px] text-text-muted mt-0.5">{{ conn.host }}:{{ conn.port }}/{{ conn.database }}</div>
          </div>
          <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
            <button
              @click.stop="toggleDefault(conn.id)"
              class="p-1 rounded hover:bg-navy-tertiary text-text-muted hover:text-accent-amber transition-colors"
              :title="defaultConnectionId === conn.id ? 'Remove Default' : 'Set as Default'"
            >
              <svg class="w-3.5 h-3.5" :class="defaultConnectionId === conn.id ? 'fill-current text-accent-amber' : ''" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
              </svg>
            </button>
            <button
              @click.stop="handleTest(conn)"
              class="p-1 rounded hover:bg-navy-tertiary text-text-muted hover:text-accent-blue transition-colors"
              title="Test Connection"
            >
              <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
                <path d="M22 4 12 14.01l-3-3" />
              </svg>
            </button>
            <button
              @click.stop="handleEdit(conn)"
              class="p-1 rounded hover:bg-navy-tertiary text-text-muted hover:text-teal-accent transition-colors"
              title="Edit Connection"
            >
              <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M12 20h9" />
                <path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z" />
              </svg>
            </button>
            <button
              @click.stop="handleDelete(conn)"
              class="p-1 rounded hover:bg-navy-tertiary text-text-muted hover:text-accent-red transition-colors"
              title="Delete"
            >
              <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M3 6h18M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
              </svg>
            </button>
          </div>

          <!-- Test result indicator -->
          <div v-if="testResults[conn.id]" class="flex-shrink-0">
            <span v-if="testResults[conn.id]!.ok" class="text-[10px] text-accent-green">
              {{ testResults[conn.id]!.latency_ms }}ms
            </span>
            <span v-else class="text-[10px] text-accent-red">Failed</span>
          </div>
        </div>
      </div>

      <!-- Empty state -->
      <div v-else class="text-center py-12 border border-dashed border-navy-border rounded-lg">
        <svg class="w-10 h-10 mx-auto text-text-muted mb-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <ellipse cx="12" cy="5" rx="9" ry="3" />
          <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3" />
          <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5" />
        </svg>
        <p class="text-sm text-text-muted mb-2">No connections yet</p>
        <p class="text-xs text-text-muted">Create a connection to get started</p>
      </div>
    </div>

    <!-- Connection Dialog -->
    <ConnectionForm
      v-if="showConnectionForm"
      :connection="editingConnection"
      @close="closeConnectionForm"
      @save="handleSaveConnection"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useConnectionsStore } from '../stores/connections'
import { useSchemaStore } from '../stores/schema'
import { useUiStore } from '../stores/ui'
import { useTabsStore } from '../stores/tabs'
import type { Connection, ConnectionTestResult } from '../types'
import ConnectionForm from '../components/connection/ConnectionForm.vue'

const router = useRouter()
const connectionsStore = useConnectionsStore()
const schemaStore = useSchemaStore()
const uiStore = useUiStore()
const tabsStore = useTabsStore()

const showConnectionForm = ref(false)
const editingConnection = ref<Connection | null>(null)
const testResults = ref<Record<string, ConnectionTestResult>>({})
const syncing = ref(false)
const syncingPgAdmin = ref(false)
const defaultConnectionId = ref<string | null>(localStorage.getItem('connections:default_connection_id'))

function toggleDefault(id: string) {
  if (defaultConnectionId.value === id) {
    defaultConnectionId.value = null
    localStorage.removeItem('connections:default_connection_id')
    uiStore.addNotification({
      type: 'info',
      title: 'Default Cleared',
      message: 'Default connection has been cleared.'
    })
  } else {
    defaultConnectionId.value = id
    localStorage.setItem('connections:default_connection_id', id)
    uiStore.addNotification({
      type: 'success',
      title: 'Default Set',
      message: 'Connection successfully set as default.'
    })
  }
}

function handleNewConnection() {
  editingConnection.value = null
  showConnectionForm.value = true
}

function handleEdit(conn: Connection) {
  editingConnection.value = conn
  showConnectionForm.value = true
}

function closeConnectionForm() {
  showConnectionForm.value = false
  editingConnection.value = null
}

async function handleSync() {
  syncing.value = true
  try {
    const count = await connectionsStore.syncServerConnections()
    uiStore.addNotification({
      type: 'success',
      title: 'Sync Successful',
      message: `Successfully synchronized ${count} connection(s) from server database`,
    })
  } catch (e: any) {
    uiStore.addNotification({
      type: 'error',
      title: 'Sync Failed',
      message: e.message || String(e),
    })
  } finally {
    syncing.value = false
  }
}

async function handleSyncPgAdmin() {
  syncingPgAdmin.value = true
  try {
    const bindings = (window as any).go.main.App
    const count = await bindings.SyncPgAdminConnections()
    uiStore.addNotification({
      type: 'success',
      title: 'Sync pgAdmin',
      message: `Successfully imported ${count} connections from pgAdmin.`,
    })
    await connectionsStore.loadConnections()
  } catch (err: any) {
    uiStore.addNotification({
      type: 'error',
      title: 'Sync failed',
      message: err.message || String(err),
    })
  } finally {
    syncingPgAdmin.value = false
  }
}

onMounted(async () => {
  await connectionsStore.loadConnections()
  
  // Auto-connect to default connection on startup
  if (defaultConnectionId.value && !connectionsStore.hasAutoConnected) {
    const conn = connectionsStore.connections.find(c => c.id === defaultConnectionId.value)
    if (conn) {
      connectionsStore.hasAutoConnected = true
      await handleConnect(conn)
    }
  }
})

async function handleConnect(conn: Connection) {
  try {
    await connectionsStore.connect(conn.id)
    await schemaStore.loadSchemas(conn.id)
    router.push('/workspace')
  } catch (e: any) {
    try {
      await connectionsStore.disconnect(conn.id)
    } catch (disErr) {
      // Ignore disconnect cleanup errors
    }
    uiStore.addNotification({
      type: 'error',
      title: 'Connection Failed',
      message: e.message || String(e),
    })
  }
}

async function handleTest(conn: Connection) {
  const result = await connectionsStore.testConnection(conn)
  testResults.value[conn.id] = result
  uiStore.addNotification({
    type: result.ok ? 'success' : 'error',
    title: result.ok ? 'Connection Successful' : 'Connection Failed',
    message: result.ok ? `Latency: ${result.latency_ms}ms` : result.message,
  })
}

async function handleDelete(conn: Connection) {
  await connectionsStore.deleteConnection(conn.id)
  if (defaultConnectionId.value === conn.id) {
    defaultConnectionId.value = null
    localStorage.removeItem('connections:default_connection_id')
  }
  uiStore.addNotification({
    type: 'success',
    title: 'Deleted',
    message: `Connection "${conn.name}" deleted`,
  })
}

async function handleSaveConnection(conn: Connection) {
  try {
    if (editingConnection.value) {
      await connectionsStore.updateConnection(conn)
      uiStore.addNotification({
        type: 'success',
        title: 'Updated',
        message: `Connection "${conn.name}" updated`,
      })
    } else {
      await connectionsStore.addConnection(conn)
      uiStore.addNotification({
        type: 'success',
        title: 'Saved',
        message: `Connection "${conn.name}" saved`,
      })
    }
    closeConnectionForm()
  } catch (e: any) {
    uiStore.addNotification({
      type: 'error',
      title: 'Save Failed',
      message: e.message || String(e),
    })
  }
}
</script>
