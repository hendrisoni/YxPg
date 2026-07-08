<template>
  <Modal :show="true" title="Workspace Manager" @close="$emit('close')" size="lg">
    <div class="space-y-4 py-2 select-none">
      <!-- Top Action Bar (Save Current & Import) -->
      <div class="flex flex-col gap-3 p-3 bg-navy-tertiary border border-navy-border rounded-md">
        <h4 class="text-xs font-semibold text-text-primary">Save Current Workspace</h4>
        <div class="flex items-center gap-2">
          <input
            v-model="newWorkspaceName"
            type="text"
            placeholder="Workspace name (e.g. Production DBs, Testing)..."
            class="flex-1 px-3 py-1.5 text-xs bg-navy-secondary border border-navy-border rounded-md text-text-primary placeholder-text-muted focus:border-teal-accent focus:outline-none"
            @keyup.enter="saveCurrentWorkspace"
          />
          <button
            @click="saveCurrentWorkspace"
            class="px-3 py-1.5 text-xs bg-teal-accent text-navy-primary font-semibold rounded-md hover:bg-teal-hover transition-colors flex items-center gap-1 flex-shrink-0"
          >
            <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v13a2 2 0 0 1-2 2z" />
              <polyline points="17 21 17 13 7 13 7 21" />
              <polyline points="7 3 7 8 15 8" />
            </svg>
            Save Current
          </button>
        </div>

        <div class="border-t border-navy-border/50 pt-2 flex items-center justify-between">
          <span class="text-[10px] text-text-muted">Or import a workspace configuration file (.json)</span>
          <button
            @click="triggerImport"
            class="px-2.5 py-1 text-[11px] bg-navy-secondary border border-navy-border hover:border-text-muted text-text-primary rounded-md transition-colors flex items-center gap-1"
          >
            <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
              <polyline points="17 8 12 3 7 8" />
              <line x1="12" y1="3" x2="12" y2="15" />
            </svg>
            Import File
          </button>
          <input
            ref="fileInput"
            type="file"
            accept=".json"
            class="hidden"
            @change="handleFileImport"
          />
        </div>
      </div>

      <!-- Search Saved Workspaces -->
      <div v-if="savedWorkspaces.length > 0" class="relative">
        <svg class="absolute left-2.5 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-text-muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8" />
          <path d="M21 21l-4.35-4.35" />
        </svg>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Filter saved workspaces..."
          class="w-full pl-8 pr-3 py-1.5 text-xs bg-navy-tertiary border border-navy-border rounded-md text-text-primary placeholder-text-muted focus:border-teal-accent focus:outline-none"
        />
      </div>

      <!-- Saved Workspaces List -->
      <div class="space-y-2">
        <h4 class="text-xs font-semibold text-text-primary px-1">Saved Workspaces</h4>
        <div v-if="filteredWorkspaces.length === 0" class="text-center py-6 border border-dashed border-navy-border rounded-md bg-navy-tertiary/20">
          <svg class="w-8 h-8 mx-auto text-text-muted opacity-60" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
          </svg>
          <p class="text-xs text-text-muted mt-2">
            {{ searchQuery ? 'No matching workspaces found' : 'No saved workspaces yet' }}
          </p>
        </div>

        <div v-else class="max-h-[300px] overflow-y-auto space-y-2 pr-1">
          <div
            v-for="ws in filteredWorkspaces"
            :key="ws.id"
            class="flex items-center justify-between p-3 bg-navy-tertiary/60 border border-navy-border rounded-md hover:bg-navy-hover transition-colors group"
          >
            <div class="flex-1 min-w-0 pr-4">
              <div class="flex items-center gap-2">
                <span class="text-xs font-semibold text-text-primary truncate group-hover:text-teal-accent transition-colors">
                  {{ ws.name }}
                </span>
                <span class="text-[9px] px-1.5 py-0.5 rounded-full bg-navy-secondary text-text-muted font-medium">
                  {{ countNodes(ws.treeData) }} nodes
                </span>
              </div>
              <div class="text-[10px] text-text-muted mt-1 flex items-center gap-1.5">
                <span>Updated: {{ formatDate(ws.updatedAt) }}</span>
              </div>
            </div>

            <!-- Action buttons -->
            <div class="flex items-center gap-1.5">
              <!-- Load Workspace -->
              <button
                @click="loadWorkspace(ws)"
                title="Load Workspace"
                class="px-2.5 py-1 text-[11px] bg-teal-accent/15 text-teal-accent hover:bg-teal-accent hover:text-navy-primary font-medium rounded transition-colors"
              >
                Load
              </button>

              <!-- Overwrite -->
              <button
                @click="overwriteWorkspace(ws)"
                title="Overwrite with Current Workspace"
                class="p-1 rounded text-text-muted hover:text-text-primary hover:bg-navy-secondary transition-colors"
              >
                <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67" />
                </svg>
              </button>

              <!-- Export JSON -->
              <button
                @click="exportWorkspace(ws)"
                title="Export Workspace as File"
                class="p-1 rounded text-text-muted hover:text-text-primary hover:bg-navy-secondary transition-colors"
              >
                <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                  <polyline points="7 10 12 15 17 10" />
                  <line x1="12" y1="15" x2="12" y2="3" />
                </svg>
              </button>

              <!-- Delete -->
              <button
                @click="deleteWorkspace(ws)"
                title="Delete Saved Workspace"
                class="p-1 rounded text-text-muted hover:text-red-400 hover:bg-navy-secondary transition-colors"
              >
                <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6" />
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Modal>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useWorkspaceStore } from '../../stores/workspace'
import { useUiStore } from '../../stores/ui'
import Modal from './Modal.vue'

interface SavedWorkspace {
  id: string
  name: string
  createdAt: number
  updatedAt: number
  treeData: any[]
}

defineEmits(['close'])

const workspaceStore = useWorkspaceStore()
const uiStore = useUiStore()

const savedWorkspaces = ref<SavedWorkspace[]>([])
const newWorkspaceName = ref('')
const searchQuery = ref('')
const fileInput = ref<HTMLInputElement | null>(null)

onMounted(() => {
  loadSavedWorkspacesList()
})

function loadSavedWorkspacesList() {
  const data = localStorage.getItem('yxpg_saved_workspaces')
  if (data) {
    try {
      savedWorkspaces.value = JSON.parse(data)
    } catch (e) {
      savedWorkspaces.value = []
    }
  } else {
    savedWorkspaces.value = []
  }
}

function saveSavedWorkspacesList() {
  localStorage.setItem('yxpg_saved_workspaces', JSON.stringify(savedWorkspaces.value))
}

const filteredWorkspaces = computed(() => {
  if (!searchQuery.value) return savedWorkspaces.value
  const query = searchQuery.value.toLowerCase()
  return savedWorkspaces.value.filter(ws => ws.name.toLowerCase().includes(query))
})

// Save current workspace store tree data under a name
function saveCurrentWorkspace() {
  const name = newWorkspaceName.value.trim()
  if (!name) {
    uiStore.addNotification({
      type: 'warning',
      title: 'Validation Error',
      message: 'Please enter a name for the workspace.',
    })
    return
  }

  const existingIndex = savedWorkspaces.value.findIndex(
    ws => ws.name.toLowerCase() === name.toLowerCase()
  )

  const treeCopy = JSON.parse(JSON.stringify(workspaceStore.workspaceTree))

  if (existingIndex !== -1) {
    // Overwrite existing by name
    if (confirm(`A workspace named "${name}" already exists. Overwrite it?`)) {
      savedWorkspaces.value[existingIndex].treeData = treeCopy
      savedWorkspaces.value[existingIndex].updatedAt = Date.now()
      saveSavedWorkspacesList()
      uiStore.addNotification({
        type: 'success',
        title: 'Workspace Updated',
        message: `Workspace "${name}" has been updated.`,
      })
      newWorkspaceName.value = ''
    }
  } else {
    // Create new
    const newWorkspace: SavedWorkspace = {
      id: `ws_${Date.now()}_${Math.random().toString(36).substr(2, 5)}`,
      name,
      createdAt: Date.now(),
      updatedAt: Date.now(),
      treeData: treeCopy,
    }
    savedWorkspaces.value.push(newWorkspace)
    saveSavedWorkspacesList()
    uiStore.addNotification({
      type: 'success',
      title: 'Workspace Saved',
      message: `Workspace "${name}" has been saved.`,
    })
    newWorkspaceName.value = ''
  }
}

// Load a saved workspace into the active workspace store
async function loadWorkspace(ws: SavedWorkspace) {
  if (confirm(`Load workspace "${ws.name}"? This will replace your current sidebar setup.`)) {
    try {
      workspaceStore.workspaceTree = JSON.parse(JSON.stringify(ws.treeData))
      await workspaceStore.saveWorkspace()
      uiStore.addNotification({
        type: 'success',
        title: 'Workspace Loaded',
        message: `Workspace "${ws.name}" loaded successfully.`,
      })
    } catch (err: any) {
      uiStore.addNotification({
        type: 'error',
        title: 'Load Failed',
        message: err.message || String(err),
      })
    }
  }
}

// Overwrite a saved workspace with the current workspace store state
function overwriteWorkspace(ws: SavedWorkspace) {
  if (confirm(`Are you sure you want to overwrite "${ws.name}" with your current workspace layout?`)) {
    ws.treeData = JSON.parse(JSON.stringify(workspaceStore.workspaceTree))
    ws.updatedAt = Date.now()
    saveSavedWorkspacesList()
    uiStore.addNotification({
      type: 'success',
      title: 'Workspace Saved',
      message: `Workspace "${ws.name}" overwritten with current layout.`,
    })
  }
}

// Delete a saved workspace
function deleteWorkspace(ws: SavedWorkspace) {
  if (confirm(`Are you sure you want to delete workspace "${ws.name}"?`)) {
    savedWorkspaces.value = savedWorkspaces.value.filter(item => item.id !== ws.id)
    saveSavedWorkspacesList()
    uiStore.addNotification({
      type: 'info',
      title: 'Workspace Deleted',
      message: `Workspace "${ws.name}" deleted.`,
    })
  }
}

// Export a workspace to a JSON file
function exportWorkspace(ws: SavedWorkspace) {
  try {
    const dataStr = JSON.stringify(ws, null, 2)
    const blob = new Blob([dataStr], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `yxpg-workspace-${ws.name.replace(/[^a-z0-9]/gi, '_').toLowerCase()}.json`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  } catch (err: any) {
    uiStore.addNotification({
      type: 'error',
      title: 'Export Failed',
      message: err.message || String(err),
    })
  }
}

// Trigger import file input click
function triggerImport() {
  fileInput.value?.click()
}

// Handle imported JSON file
function handleFileImport(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      const content = e.target?.result as string
      const parsed = JSON.parse(content)

      // Validation
      if (!parsed.name || !Array.isArray(parsed.treeData)) {
        throw new Error('Invalid workspace file format. Must contain a name and a treeData array.')
      }

      // Check for duplicate name, rename if needed
      let importedName = parsed.name
      let counter = 1
      while (savedWorkspaces.value.some(ws => ws.name.toLowerCase() === importedName.toLowerCase())) {
        importedName = `${parsed.name} (${counter})`
        counter++
      }

      const importedWorkspace: SavedWorkspace = {
        id: `ws_${Date.now()}_${Math.random().toString(36).substr(2, 5)}`,
        name: importedName,
        createdAt: parsed.createdAt || Date.now(),
        updatedAt: parsed.updatedAt || Date.now(),
        treeData: parsed.treeData,
      }

      savedWorkspaces.value.push(importedWorkspace)
      saveSavedWorkspacesList()

      uiStore.addNotification({
        type: 'success',
        title: 'Workspace Imported',
        message: `Successfully imported "${importedWorkspace.name}".`,
      })
    } catch (err: any) {
      uiStore.addNotification({
        type: 'error',
        title: 'Import Failed',
        message: err.message || 'Could not parse workspace JSON.',
      })
    } finally {
      // Clear input
      target.value = ''
    }
  }
  reader.readAsText(file)
}

// Count nodes in tree recursively
function countNodes(nodes: any[]): number {
  if (!nodes || !Array.isArray(nodes)) return 0
  let count = 0
  for (const node of nodes) {
    count++
    if (node.children && node.children.length > 0) {
      count += countNodes(node.children)
    }
  }
  return count
}

// Format date
function formatDate(timestamp: number): string {
  return new Date(timestamp).toLocaleString(undefined, {
    dateStyle: 'short',
    timeStyle: 'short',
  })
}
</script>
