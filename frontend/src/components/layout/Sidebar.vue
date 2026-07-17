<template>
  <aside class="h-full bg-navy-secondary border-r border-navy-border flex flex-row"
    :class="{ 'transition-[width] duration-150': !uiStore.isResizingSidebar }" :style="{ width: sidebarWidth + 'px' }">
    <!-- Left Toolbar (VS Code-like Activity Bar) -->
    <div
      class="w-12 bg-[#090d16] border-r border-navy-border flex flex-col items-center py-4 justify-between flex-shrink-0">
      <!-- Top Icons -->
      <div class="flex flex-col items-center gap-5 w-full">
        <!-- Database/Home Button -->
        <button @click="tabsStore.activeTabId = null"
          class="p-2 rounded-lg text-teal-accent hover:bg-navy-hover transition-colors cursor-pointer"
          title="Go to Home">
          <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <ellipse cx="12" cy="5" rx="9" ry="3" />
            <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3" />
            <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5" />
          </svg>
        </button>

        <!-- Add Category Button -->
        <button @click="$emit('addCategory')"
          class="p-2 rounded-lg text-text-secondary hover:text-teal-accent hover:bg-navy-hover transition-colors cursor-pointer"
          title="Add Category">
          <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 5v14M5 12h14" />
          </svg>
        </button>

        <!-- Load / Save Workspace Button -->
        <button @click="$emit('openWorkspace')"
          class="p-2 rounded-lg text-text-secondary hover:text-teal-accent hover:bg-navy-hover transition-colors cursor-pointer"
          title="Load / Save Workspace">
          <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
            stroke-linecap="round" stroke-linejoin="round">
            <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v13a2 2 0 0 1-2 2z" />
            <polyline points="17 21 17 13 7 13 7 21" />
            <polyline points="7 3 7 8 15 8" />
          </svg>
        </button>

        <!-- Referential Integrity Button -->
        <button @click="$emit('openReferential')"
          class="p-2 rounded-lg text-text-secondary hover:text-teal-accent hover:bg-navy-hover transition-colors cursor-pointer"
          title="Referential Integrity Relation">
          <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="3" width="6" height="6" rx="1" />
            <rect x="15" y="15" width="6" height="6" rx="1" />
            <path d="M9 6h3a2 2 0 0 1 2 2v4a2 2 0 0 0 2 2h1" />
          </svg>
        </button>

        <!-- Functions & Triggers Button -->
        <button @click="$emit('openFunctionsTriggers')"
          class="p-2 rounded-lg text-text-secondary hover:text-teal-accent hover:bg-navy-hover transition-colors cursor-pointer"
          title="Functions & Triggers Manager">
          <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2" />
          </svg>
        </button>


      </div>

      <!-- Bottom Icons -->
      <div class="flex flex-col items-center w-full gap-3">
        <!-- Backup Button -->
        <button @click="$emit('openBackup')"
          class="p-2 rounded-lg text-text-secondary hover:text-teal-accent hover:bg-navy-hover transition-colors cursor-pointer"
          title="Backup Database (pg_dump)">
          <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
            stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
            <polyline points="17 8 12 3 7 8" />
            <line x1="12" y1="3" x2="12" y2="15" />
          </svg>
        </button>

        <!-- Settings Button -->
        <button @click="$emit('openSettings')"
          class="p-2 rounded-lg text-text-secondary hover:text-teal-accent hover:bg-navy-hover transition-colors cursor-pointer"
          title="Settings">
          <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
            stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="3" />
            <path
              d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Main Sidebar Content (Workspace explorer) -->
    <div class="flex-1 flex flex-col overflow-hidden bg-navy-secondary">
      <!-- Search -->
      <div class="px-3 py-2 border-b border-navy-border flex items-center gap-1.5">
        <div class="relative flex-1">
          <svg class="absolute left-2 top-1/2 -translate-y-1/2 w-4 h-4 text-text-muted" viewBox="0 0 24 24" fill="none"
            stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8" />
            <path d="M21 21l-4.35-4.35" />
          </svg>
          <input v-model="searchQuery" type="text" placeholder="Search workspace..."
            class="w-full pl-8 pr-3 py-1.5 text-xs bg-navy-tertiary border border-navy-border rounded-md text-text-primary placeholder-text-muted focus:border-teal-accent focus:outline-none" />
        </div>

        <!-- Expand / Collapse All Button (Icon Only) -->
        <button @click="toggleAllNodes"
          class="p-1.5 rounded text-text-muted hover:text-text-primary hover:bg-navy-hover transition-colors flex-shrink-0"
          :title="isAllExpanded ? 'Collapse All' : 'Expand All'">
          <svg v-if="isAllExpanded" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor"
            stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="m17 18-5-5-5 5M17 11l-5-5-5 5" />
          </svg>
          <svg v-else class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
            stroke-linecap="round" stroke-linejoin="round">
            <path d="m7 6 5 5 5-5M7 13l5 5 5-5" />
          </svg>
        </button>
      </div>

      <!-- Tree Container (Root Drop Zone) -->
      <div class="flex-1 overflow-y-auto overflow-x-hidden py-1" @dragover.prevent @drop="handleRootDrop">
        <!-- Loading indicator -->
        <div v-if="workspaceStore.loading" class="px-3 py-4 text-center">
          <div class="inline-block w-5 h-5 border-2 border-teal-accent border-t-transparent rounded-full animate-spin">
          </div>
          <p class="text-xs text-text-muted mt-2">Loading workspace...</p>
        </div>

        <!-- Workspace Tree -->
        <div v-else-if="filteredTree.length > 0">
          <SchemaTreeNode v-for="node in filteredTree" :key="node.id" :node="node" :level="0"
            :selected-node-id="selectedNodeId" @select-node="selectedNodeId = $event" @open-table="handleOpenTable"
            @open-query="handleOpenQuery" @copy-name="handleCopyName" @copy-select="handleCopySelect"
            @view-ddl="handleViewDDL" @drop-table="handleDropTable" @delete-node="handleDeleteNode"
            @rename-node="handleRenameNode" />
        </div>

        <!-- Empty state -->
        <div v-else class="px-3 py-8 text-center">
          <svg class="w-8 h-8 mx-auto text-text-muted" viewBox="0 0 24 24" fill="none" stroke="currentColor"
            stroke-width="1.5">
            <path d="M20 7H4a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2Z" />
            <path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16" />
          </svg>
          <p class="text-xs text-text-muted mt-2">Workspace is empty</p>
          <p class="text-[10px] text-text-muted mt-1">Press Ctrl+K to search and add tables</p>
          <button @click="$emit('addCategory')"
            class="mt-3 px-3 py-1.5 text-xs bg-teal-accent text-navy-primary rounded-md font-medium hover:bg-teal-hover transition-colors">
            Add Category
          </button>
        </div>
      </div>
    </div>

    <!-- Category Add/Rename Modal -->
    <Modal
      :show="isCategoryModalOpen"
      :title="categoryModalTitle"
      @close="isCategoryModalOpen = false"
      size="sm"
    >
      <div class="space-y-4 py-2 select-none">
        <div class="flex flex-col gap-1.5">
          <label for="category-name" class="text-xs font-semibold text-text-primary">
            Category Name
          </label>
          <input
            ref="categoryInputRef"
            v-model="categoryModalInput"
            id="category-name"
            type="text"
            placeholder="e.g. Production Databases"
            class="w-full text-xs bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-text-primary focus:border-teal-accent focus:outline-none placeholder-text-muted"
            required
            @keydown.enter="handleCategorySubmit"
          />
          <p v-if="categoryModalError" class="text-[10px] text-red-500">
            {{ categoryModalError }}
          </p>
        </div>
      </div>
      <template #footer>
        <button
          type="button"
          @click="isCategoryModalOpen = false"
          class="px-4 py-1.5 text-xs text-text-secondary hover:text-text-primary transition-colors cursor-pointer"
        >
          Cancel
        </button>
        <button
          type="button"
          @click="handleCategorySubmit"
          class="px-4 py-1.5 text-xs bg-teal-accent text-navy-primary rounded-md font-medium hover:bg-teal-hover transition-colors cursor-pointer"
        >
          {{ categoryModalMode === 'add' ? 'Add' : 'Rename' }}
        </button>
      </template>
    </Modal>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { useConnectionsStore } from '../../stores/connections'
import { useWorkspaceStore } from '../../stores/workspace'
import { useTabsStore } from '../../stores/tabs'
import { useUiStore } from '../../stores/ui'
import type { TreeNode } from '../../types'
import SchemaTreeNode from '../schema/SchemaTreeNode.vue'
import Modal from '../shared/Modal.vue'

const emit = defineEmits(['toggle', 'newConnection', 'searchTable', 'openSettings', 'openWorkspace', 'openBackup', 'addCategory', 'openReferential', 'openFunctionsTriggers'])

const connectionsStore = useConnectionsStore()
const workspaceStore = useWorkspaceStore()
const tabsStore = useTabsStore()
const uiStore = useUiStore()

const searchQuery = ref('')
const isAllExpanded = ref(false)
const selectedNodeId = ref<string | null>(null)
const sidebarWidth = computed(() => uiStore.sidebarWidth)

const isCategoryModalOpen = ref(false)
const categoryModalMode = ref<'add' | 'rename'>('add')
const categoryModalInput = ref('')
const categoryModalNodeId = ref('')
const categoryModalError = ref('')
const categoryInputRef = ref<HTMLInputElement | null>(null)

const categoryModalTitle = computed(() => {
  return categoryModalMode.value === 'add' ? 'Add Category' : 'Rename Category'
})

watch(isCategoryModalOpen, (isOpen) => {
  if (isOpen) {
    categoryModalError.value = ''
    nextTick(() => {
      categoryInputRef.value?.focus()
      categoryInputRef.value?.select()
    })
  }
})

onMounted(() => {
  workspaceStore.loadWorkspace()
})

const filteredTree = computed(() => {
  if (!searchQuery.value) return workspaceStore.workspaceTree
  const query = searchQuery.value.toLowerCase()
  return filterTreeNodes(workspaceStore.workspaceTree, query)
})

function filterTreeNodes(nodes: TreeNode[], query: string): TreeNode[] {
  return nodes.reduce<TreeNode[]>((acc, node) => {
    const matchesSelf = node.label.toLowerCase().includes(query)
    const filteredChildren = node.children ? filterTreeNodes(node.children, query) : []
    if (matchesSelf || filteredChildren.length > 0) {
      acc.push({
        ...node,
        children: filteredChildren.length > 0 ? filteredChildren : node.children,
        expanded: true,
      })
    }
    return acc
  }, [])
}

async function handleAddCategory() {
  const name = prompt('Enter category name:')
  if (name && name.trim()) {
    await workspaceStore.addCategory(name.trim())
  }
}

function handleRootDrop(event: DragEvent) {
  const nodeId = event.dataTransfer?.getData('text/plain')
  if (nodeId) {
    workspaceStore.moveNode(nodeId, null)
  }
}

interface NodeLocation {
  list: any[]
  index: number
  node: any
  parent: any | null
}

function findNodeLocation(nodes: any[], id: string, parent: any | null = null): NodeLocation | null {
  for (let i = 0; i < nodes.length; i++) {
    const node = nodes[i]
    if (node.id === id) {
      return { list: nodes, index: i, node: node, parent: parent }
    }
    if (node.children && node.children.length > 0) {
      const found = findNodeLocation(node.children, id, node)
      if (found) return found
    }
  }
  return null
}

function handleDeleteNode(id: string) {
  const loc = findNodeLocation(workspaceStore.workspaceTree, id)
  if (!loc) return

  const { list, index, node, parent } = loc
  
  // Remove node
  list.splice(index, 1)
  workspaceStore.saveWorkspace()

  uiStore.addNotification({
    type: 'info',
    title: 'Removed from Workspace',
    message: `"${node.label}" has been removed.`,
    duration: 10000,
    action: {
      label: 'Undo',
      callback: async () => {
        // Restore node
        list.splice(index, 0, node)
        if (parent) {
          parent.expanded = true
        }
        await workspaceStore.saveWorkspace()
        uiStore.addNotification({
          type: 'success',
          title: 'Restored',
          message: `"${node.label}" has been restored.`,
          duration: 3000
        })
      }
    }
  })
}

function handleRenameNode(id: string, currentLabel: string) {
  categoryModalMode.value = 'rename'
  categoryModalInput.value = currentLabel
  categoryModalNodeId.value = id
  isCategoryModalOpen.value = true
}

async function handleCategorySubmit() {
  const name = categoryModalInput.value.trim()
  if (!name) {
    categoryModalError.value = 'Category name cannot be empty'
    return
  }

  if (categoryModalMode.value === 'add') {
    await workspaceStore.addCategory(name)
  } else {
    await workspaceStore.renameNode(categoryModalNodeId.value, name)
  }

  isCategoryModalOpen.value = false
}

function setAllNodesExpandedState(nodes: any[], state: boolean) {
  for (const node of nodes) {
    node.expanded = state
    if (node.children && node.children.length > 0) {
      setAllNodesExpandedState(node.children, state)
    }
  }
}

function toggleAllNodes() {
  isAllExpanded.value = !isAllExpanded.value
  setAllNodesExpandedState(workspaceStore.workspaceTree, isAllExpanded.value)
  workspaceStore.workspaceTree = [...workspaceStore.workspaceTree]
  workspaceStore.saveWorkspace()
}


function handleOpenTable(schema: string, table: string, connectionId?: string) {
  const connId = connectionId || connectionsStore.currentConnectionId
  if (!connId) return
  tabsStore.createTab('table', {
    title: table,
    connectionId: connId,
    schema,
    table,
  })
}

function handleOpenQuery(schema: string, table: string, connectionId?: string) {
  const connId = connectionId || connectionsStore.currentConnectionId
  if (!connId) return
  tabsStore.createTab('query', {
    title: table,
    connectionId: connId,
    schema,
    table,
    sql: `SELECT * FROM ${schema}.${table}`,
  })
}

function handleCopyName(name: string) {
  navigator.clipboard.writeText(name)
  uiStore.addNotification({ type: 'info', title: 'Copied', message: `Copied "${name}" to clipboard` })
}

function handleCopySelect(schema: string, table: string) {
  const sql = `SELECT * FROM ${schema}.${table}`
  navigator.clipboard.writeText(sql)
  uiStore.addNotification({ type: 'info', title: 'Copied', message: 'SELECT query copied to clipboard' })
}

function handleViewDDL(schema: string, table: string, connectionId?: string) {
  const connId = connectionId || connectionsStore.currentConnectionId
  if (!connId) return
  tabsStore.createTab('ddl', {
    title: table,
    connectionId: connId,
    schema,
    table,
  })
}

function handleDropTable(schema: string, table: string) {
  uiStore.addNotification({ type: 'warning', title: 'Drop Table', message: `Drop ${schema}.${table}? (not implemented yet)` })
}

async function handleDisconnect() {
  if (connectionsStore.currentConnectionId) {
    await connectionsStore.disconnect(connectionsStore.currentConnectionId)
  }
}
</script>
