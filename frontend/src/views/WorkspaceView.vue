<template>
  <div class="h-screen w-screen flex flex-col overflow-hidden bg-navy-primary">
    <!-- Sidebar -->
    <div class="flex flex-1 overflow-hidden">
      <Sidebar
        @new-connection="showNewConnection = true"
        @search-table="showTableSearch = !showTableSearch"
        @open-settings="showSettings = true"
        @open-workspace="showWorkspaceModal = true"
      />

      <!-- Resize Handle -->
      <div
        class="w-[3px] cursor-col-resize select-none z-50 hover:bg-teal-accent active:bg-teal-accent transition-colors duration-150"
        @mousedown="startResizeSidebar"
      ></div>

      <!-- Main content -->
      <div class="flex-1 flex flex-col overflow-hidden">
        <!-- Tab bar -->
        <TabBar @new-tab="createNewTab" @open-builder="openQueryBuilder" />

        <!-- Tab content -->
        <div class="flex-1 overflow-hidden">
          <HomeView v-if="!tabsStore.activeTab" />
          <QueryView
            v-else-if="tabsStore.activeTab.type === 'query'"
            :tab="tabsStore.activeTab"
            :key="'query-' + tabsStore.activeTab.id"
          />
          <TableView
            v-else-if="tabsStore.activeTab.type === 'table'"
            :tab="tabsStore.activeTab"
            :key="'table-' + tabsStore.activeTab.id"
          />
          <BuilderView
            v-else-if="tabsStore.activeTab.type === 'builder'"
            :tab="tabsStore.activeTab"
            :key="'builder-' + tabsStore.activeTab.id"
          />
          <DDLView
            v-else-if="tabsStore.activeTab.type === 'ddl'"
            :tab="tabsStore.activeTab"
            :key="'ddl-' + tabsStore.activeTab.id"
          />
          <QueryLogView
            v-else-if="tabsStore.activeTab.type === 'log'"
            :key="'log-' + tabsStore.activeTab.id"
          />
        </div>
      </div>
    </div>

    <!-- Status bar -->
    <StatusBar ref="statusBar" />

    <!-- Notifications -->
    <Notification
      v-for="notif in uiStore.notifications"
      :key="notif.id"
      :show="true"
      :type="notif.type"
      :title="notif.title"
      :message="notif.message"
      @close="uiStore.removeNotification(notif.id)"
    />

    <!-- New Connection Dialog -->
    <ConnectionForm
      v-if="showNewConnection"
      @close="showNewConnection = false"
      @save="handleSaveConnection"
    />

    <!-- Database Objects Search (Ctrl+K) -->
    <TableSearchPalette
      :show="showTableSearch"
      @close="showTableSearch = false"
    />

    <!-- Settings Dialog -->
    <SettingsModal
      v-if="showSettings"
      @close="showSettings = false"
    />

    <!-- Workspace Management Dialog -->
    <WorkspaceModal
      v-if="showWorkspaceModal"
      @close="showWorkspaceModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useConnectionsStore } from '../stores/connections'
import { useSchemaStore } from '../stores/schema'
import { useTabsStore } from '../stores/tabs'
import { useUiStore } from '../stores/ui'
import { setupKeyboardShortcuts } from '../utils/shortcuts'
import type { Connection } from '../types'

import Sidebar from '../components/layout/Sidebar.vue'
import TabBar from '../components/layout/TabBar.vue'
import StatusBar from '../components/layout/StatusBar.vue'
import Notification from '../components/shared/Notification.vue'
import ConnectionForm from '../components/connection/ConnectionForm.vue'
import TableSearchPalette from '../components/shared/TableSearchPalette.vue'
import SettingsModal from '../components/shared/SettingsModal.vue'
import WorkspaceModal from '../components/shared/WorkspaceModal.vue'
import HomeView from './HomeView.vue'
import QueryView from './QueryView.vue'
import TableView from './TableView.vue'
import BuilderView from './BuilderView.vue'
import DDLView from './DDLView.vue'
import QueryLogView from './QueryLogView.vue'

const connectionsStore = useConnectionsStore()
const schemaStore = useSchemaStore()
const tabsStore = useTabsStore()
const uiStore = useUiStore()

const showNewConnection = ref(false)
const showTableSearch = ref(false)
const showSettings = ref(false)
const showWorkspaceModal = ref(false)
const statusBar = ref<InstanceType<typeof StatusBar> | null>(null)

function startResizeSidebar(e: MouseEvent) {
  e.preventDefault()
  uiStore.isResizingSidebar = true
  
  const handleMouseMove = (event: MouseEvent) => {
    uiStore.setSidebarWidth(event.clientX)
  }
  
  const handleMouseUp = () => {
    uiStore.isResizingSidebar = false
    window.removeEventListener('mousemove', handleMouseMove)
    window.removeEventListener('mouseup', handleMouseUp)
  }
  
  window.addEventListener('mousemove', handleMouseMove)
  window.addEventListener('mouseup', handleMouseUp)
}

// Keyboard shortcuts
onMounted(() => {
  connectionsStore.loadConnections()
  const cleanup = setupKeyboardShortcuts({
    newQueryTab: () => createNewTab(),
    closeTab: () => tabsStore.closeActiveTab(),
    toggleSidebar: () => uiStore.toggleSidebar(),
    toggleHistory: () => uiStore.toggleHistory(),
    commandPalette: () => uiStore.toggleCommandPalette(),
    refreshSchema: () => {
      if (connectionsStore.currentConnectionId) {
        schemaStore.refreshSchema(connectionsStore.currentConnectionId)
      }
    },
    newConnection: () => {
      showNewConnection.value = true
    },
    toggleTableSearch: () => {
      showTableSearch.value = !showTableSearch.value
    },
    runQuery: () => {
      window.dispatchEvent(new CustomEvent('run-active-query'))
    },
    formatSql: () => {
      window.dispatchEvent(new CustomEvent('format-active-query'))
    },
    openBuilder: () => {
      tabsStore.createTab('builder', { title: 'Query Builder' })
    },
    openDesigner: () => {
      tabsStore.createTab('ddl', { title: 'Table Designer' })
    },
  })

  onUnmounted(cleanup)
})

function createNewTab() {
  tabsStore.createTab('query', {
    title: `Query ${tabsStore.tabs.length + 1}`,
    connectionId: connectionsStore.currentConnectionId || undefined,
  })
}

function openQueryBuilder() {
  tabsStore.createTab('builder', { title: 'Query Builder' })
}

async function handleSaveConnection(conn: Connection) {
  try {
    await connectionsStore.addConnection(conn)
    showNewConnection.value = false
    uiStore.addNotification({
      type: 'success',
      title: 'Saved',
      message: `Connection "${conn.name}" saved`,
    })
  } catch (e: any) {
    uiStore.addNotification({
      type: 'error',
      title: 'Save Failed',
      message: e.message || String(e),
    })
  }
}
</script>
