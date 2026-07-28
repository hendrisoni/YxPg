<template>
  <Modal :show="show" title="Database Objects Search (Ctrl+K)" @close="$emit('close')" size="lg" resizable storage-key="yxpg_table_search_modal_size">
    <div class="flex flex-col h-full min-h-[350px]">
      <!-- Search Input & View Switcher -->
      <div class="flex items-center gap-2 mb-3">
        <!-- Search Input -->
        <div class="relative flex-1">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-text-muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8" />
            <path d="M21 21l-4.35-4.35" />
          </svg>
          <input
            ref="searchInput"
            v-model="query"
            type="text"
            placeholder="Search tables, views, or functions... (Use Up/Down arrows, Enter to add)"
            class="w-full pl-9 pr-3 py-2 text-xs bg-navy-tertiary border border-navy-border rounded-md text-text-primary placeholder-text-muted focus:border-teal-accent focus:outline-none"
            @keydown.down.prevent="selectNext"
            @keydown.up.prevent="selectPrev"
            @keydown.enter.prevent="confirmSelection"
          />
        </div>

        <!-- View Mode Switcher (List / Tree) -->
        <div class="flex items-center bg-navy-tertiary border border-navy-border rounded-md p-0.5 text-xs flex-shrink-0">
          <button
            @click="viewMode = 'list'"
            class="flex items-center gap-1.5 px-2.5 py-1 rounded text-[11px] font-medium transition-all"
            :class="viewMode === 'list' 
              ? 'bg-teal-accent/20 text-teal-accent border border-teal-accent/30 font-semibold' 
              : 'text-text-muted hover:text-text-primary border border-transparent'"
            title="List View"
          >
            <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="8" y1="6" x2="21" y2="6"></line>
              <line x1="8" y1="12" x2="21" y2="12"></line>
              <line x1="8" y1="18" x2="21" y2="18"></line>
              <line x1="3" y1="6" x2="3.01" y2="6"></line>
              <line x1="3" y1="12" x2="3.01" y2="12"></line>
              <line x1="3" y1="18" x2="3.01" y2="18"></line>
            </svg>
            List
          </button>
          <button
            @click="viewMode = 'tree'"
            class="flex items-center gap-1.5 px-2.5 py-1 rounded text-[11px] font-medium transition-all"
            :class="viewMode === 'tree' 
              ? 'bg-teal-accent/20 text-teal-accent border border-teal-accent/30 font-semibold' 
              : 'text-text-muted hover:text-text-primary border border-transparent'"
            title="Tree View (Database -> Schema -> Table/View)"
          >
            <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="3" width="6" height="5" rx="1" />
              <rect x="14" y="9" width="7" height="5" rx="1" />
              <rect x="14" y="16" width="7" height="5" rx="1" />
              <path d="M6 8v10h8M9 11.5h5" />
            </svg>
            Tree
          </button>
        </div>
      </div>

      <!-- Content Container -->
      <div class="flex-1 overflow-y-auto min-h-0 space-y-1 pr-1">
        <div v-if="loading" class="flex flex-col items-center justify-center h-full text-text-muted text-xs">
          <div class="inline-block w-6 h-6 border-2 border-teal-accent border-t-transparent rounded-full animate-spin mb-2"></div>
          Loading database objects...
        </div>
        <div v-else-if="filteredItems.length === 0" class="flex items-center justify-center h-full text-text-muted text-xs">
          No tables, views, or functions found matching "{{ query }}"
        </div>

        <!-- LIST VIEW -->
        <template v-else-if="viewMode === 'list'">
          <div
            v-for="(item, index) in filteredItems"
            :key="item.connection_id + '_' + item.schema + '_' + item.name + '_' + item.type"
            class="flex items-center justify-between p-2 rounded cursor-pointer transition-colors text-xs border"
            :class="index === selectedIndex 
              ? 'bg-navy-hover border-teal-accent/30 text-text-primary' 
              : 'bg-navy-secondary border-transparent hover:bg-navy-hover/50 text-text-secondary'"
            @click="selectedIndex = index"
            @dblclick="handleAdd(item)"
          >
            <!-- Item Details -->
            <div class="flex items-center gap-2 min-w-0">
              <!-- Icon -->
              <component :is="iconComponent(item.type)" class="w-3.5 h-3.5 flex-shrink-0" :class="iconColor(item.type)" />
              <!-- Name -->
              <span class="font-medium truncate text-text-primary">{{ item.name }}</span>
              <span class="text-text-muted text-[10px] truncate">({{ item.schema }})</span>
            </div>

            <!-- Connection Source & Type Badges -->
            <div class="flex items-center gap-2">
              <span class="px-1.5 py-0.5 rounded text-[10px] bg-navy-tertiary text-text-muted">
                {{ item.connection_name }} ({{ item.database_name }})
              </span>
              <span class="px-1.5 py-0.5 rounded text-[10px] uppercase font-bold" :class="typeBadgeClass(item.type)">
                {{ item.type }}
              </span>
              <button
                @click.stop="handleAdd(item, false)"
                class="ml-1 px-2 py-0.5 rounded text-[10px] font-semibold bg-teal-accent/20 hover:bg-teal-accent/35 text-teal-accent border border-teal-accent/30 hover:border-teal-accent/50 transition-all flex items-center gap-0.5 active:scale-95"
              >
                <svg class="w-2.5 h-2.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                  <path d="M5 12h14M12 5v14" />
                </svg>
                Add
              </button>
            </div>
          </div>
        </template>

        <!-- TREE VIEW -->
        <template v-else-if="viewMode === 'tree'">
          <div class="space-y-2">
            <div
              v-for="dbNode in treeData"
              :key="dbNode.key"
              class="rounded-md border border-navy-border/60 bg-navy-secondary/40 overflow-hidden"
            >
              <!-- Database Node Header -->
              <div
                class="flex items-center justify-between px-2.5 py-1.5 bg-navy-tertiary/70 hover:bg-navy-hover/60 cursor-pointer select-none border-b border-navy-border/40 text-xs"
                @click="toggleNode(dbNode.key)"
              >
                <div class="flex items-center gap-2 min-w-0">
                  <svg
                    class="w-3 h-3 text-text-muted transition-transform duration-150 flex-shrink-0"
                    :class="{ 'rotate-90': isNodeExpanded(dbNode.key) }"
                    viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                  >
                    <path d="M9 18l6-6-6-6" />
                  </svg>
                  <svg class="w-4 h-4 text-teal-accent flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <ellipse cx="12" cy="5" rx="9" ry="3" />
                    <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3" />
                    <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5" />
                  </svg>
                  <span class="font-bold text-text-primary truncate">
                    {{ dbNode.connection_name }} ({{ dbNode.database_name }})
                  </span>
                  <span class="text-[10px] px-1.5 py-0.2 rounded bg-teal-accent/10 text-teal-accent border border-teal-accent/20">
                    {{ dbNode.totalCount }} {{ dbNode.totalCount === 1 ? 'object' : 'objects' }}
                  </span>
                </div>
                
                <button
                  @click.stop="handleAddAll(getAllDbItems(dbNode))"
                  class="px-2 py-0.5 rounded text-[10px] font-semibold bg-teal-accent/20 hover:bg-teal-accent/35 text-teal-accent border border-teal-accent/30 transition-all flex items-center gap-1 active:scale-95"
                  title="Add all objects in this database"
                >
                  <svg class="w-2.5 h-2.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                    <path d="M5 12h14M12 5v14" />
                  </svg>
                  Add All
                </button>
              </div>

              <!-- Schemas inside Database Node -->
              <div v-if="isNodeExpanded(dbNode.key)" class="p-1 space-y-1">
                <div
                  v-for="schemaNode in dbNode.schemas"
                  :key="dbNode.key + '_' + schemaNode.name"
                  class="rounded border border-navy-border/30 bg-navy-secondary/30 overflow-hidden"
                >
                  <!-- Schema Node Header -->
                  <div
                    class="flex items-center justify-between px-2.5 py-1 hover:bg-navy-hover/50 cursor-pointer select-none text-[11px]"
                    @click="toggleNode(dbNode.key + '_' + schemaNode.name)"
                  >
                    <div class="flex items-center gap-2 min-w-0">
                      <svg
                        class="w-3 h-3 text-text-muted transition-transform duration-150 flex-shrink-0"
                        :class="{ 'rotate-90': isNodeExpanded(dbNode.key + '_' + schemaNode.name) }"
                        viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                      >
                        <path d="M9 18l6-6-6-6" />
                      </svg>
                      <svg class="w-3.5 h-3.5 text-accent-amber flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
                      </svg>
                      <span class="font-semibold text-text-secondary truncate">
                        {{ schemaNode.name }}
                      </span>
                      <span class="text-[10px] text-text-muted">
                        ({{ schemaNode.items.length }})
                      </span>
                    </div>

                    <button
                      @click.stop="handleAddAll(schemaNode.items)"
                      class="px-1.5 py-0.5 rounded text-[9px] font-medium bg-teal-accent/15 hover:bg-teal-accent/30 text-teal-accent border border-teal-accent/25 transition-all flex items-center gap-0.5 active:scale-95"
                      title="Add all objects in schema"
                    >
                      <svg class="w-2.5 h-2.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                        <path d="M5 12h14M12 5v14" />
                      </svg>
                      Add Schema
                    </button>
                  </div>

                  <!-- Items inside Schema Node -->
                  <div v-if="isNodeExpanded(dbNode.key + '_' + schemaNode.name)" class="pl-4 pr-1 py-1 space-y-0.5 border-t border-navy-border/20">
                    <div
                      v-for="item in schemaNode.items"
                      :key="item.connection_id + '_' + item.schema + '_' + item.name + '_' + item.type"
                      class="flex items-center justify-between p-1.5 rounded cursor-pointer transition-colors text-xs border border-transparent hover:bg-navy-hover/60 hover:border-navy-border/40 group"
                      @dblclick="handleAdd(item)"
                    >
                      <!-- Item Details -->
                      <div class="flex items-center gap-2 min-w-0">
                        <component :is="iconComponent(item.type)" class="w-3.5 h-3.5 flex-shrink-0" :class="iconColor(item.type)" />
                        <span class="font-medium truncate text-text-primary">{{ item.name }}</span>
                      </div>

                      <!-- Type Badge & Add Button -->
                      <div class="flex items-center gap-2">
                        <span class="px-1.5 py-0.5 rounded text-[9px] uppercase font-bold" :class="typeBadgeClass(item.type)">
                          {{ item.type }}
                        </span>
                        <button
                          @click.stop="handleAdd(item, false)"
                          class="px-2 py-0.5 rounded text-[10px] font-semibold bg-teal-accent/20 hover:bg-teal-accent/35 text-teal-accent border border-teal-accent/30 hover:border-teal-accent/50 transition-all flex items-center gap-0.5 active:scale-95"
                        >
                          <svg class="w-2.5 h-2.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                            <path d="M5 12h14M12 5v14" />
                          </svg>
                          Add
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>
  </Modal>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import type { CatalogItem } from '../../stores/workspace'
import { useWorkspaceStore } from '../../stores/workspace'
import { useUiStore } from '../../stores/ui'
import Modal from './Modal.vue'
import { h } from 'vue'

interface SchemaNodeData {
  name: string
  items: CatalogItem[]
}

interface DatabaseNodeData {
  key: string
  connection_name: string
  database_name: string
  schemas: SchemaNodeData[]
  totalCount: number
}

const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits(['close'])

const workspaceStore = useWorkspaceStore()
const uiStore = useUiStore()

const query = ref('')
const items = ref<CatalogItem[]>([])
const loading = ref(false)
const selectedIndex = ref(0)
const searchInput = ref<HTMLInputElement | null>(null)

// View mode state
const viewMode = ref<'list' | 'tree'>((localStorage.getItem('yxpg_palette_view_mode') as 'list' | 'tree') || 'list')
const collapsedNodes = ref<Record<string, boolean>>({})

watch(viewMode, (newVal) => {
  localStorage.setItem('yxpg_palette_view_mode', newVal)
})

// Load items on open
watch(() => props.show, async (newVal) => {
  if (newVal) {
    query.value = ''
    selectedIndex.value = 0
    loading.value = true
    items.value = await workspaceStore.fetchCatalog()
    loading.value = false
    
    // Auto focus
    nextTick(() => {
      if (searchInput.value) {
        searchInput.value.focus()
      }
    })
  }
})

// Filter items
const filteredItems = computed(() => {
  let list = items.value
  if (!uiStore.settings.showFunctionsInSearch) {
    list = list.filter(item => item.type !== 'function')
  }
  if (!query.value) return list
  const q = query.value.toLowerCase()
  return list.filter(item => 
    item.name.toLowerCase().includes(q) ||
    item.schema.toLowerCase().includes(q) ||
    item.connection_name.toLowerCase().includes(q)
  )
})

// Compute Tree Hierarchy
const treeData = computed(() => {
  const map = new Map<string, { connection_name: string; database_name: string; schemaMap: Map<string, CatalogItem[]> }>()

  for (const item of filteredItems.value) {
    const dbKey = `${item.connection_id}_${item.database_name}`
    if (!map.has(dbKey)) {
      map.set(dbKey, {
        connection_name: item.connection_name,
        database_name: item.database_name,
        schemaMap: new Map<string, CatalogItem[]>()
      })
    }
    const dbEntry = map.get(dbKey)!
    const schemaKey = item.schema || 'public'
    if (!dbEntry.schemaMap.has(schemaKey)) {
      dbEntry.schemaMap.set(schemaKey, [])
    }
    dbEntry.schemaMap.get(schemaKey)!.push(item)
  }

  const result: DatabaseNodeData[] = []
  for (const [dbKey, dbVal] of map.entries()) {
    const schemas: SchemaNodeData[] = []
    let totalCount = 0
    for (const [schemaName, itemsList] of dbVal.schemaMap.entries()) {
      schemas.push({
        name: schemaName,
        items: itemsList
      })
      totalCount += itemsList.length
    }
    result.push({
      key: dbKey,
      connection_name: dbVal.connection_name,
      database_name: dbVal.database_name,
      schemas,
      totalCount
    })
  }

  return result
})

function isNodeExpanded(nodeKey: string): boolean {
  // If user typed a query, keep nodes expanded so matching items are visible
  if (query.value.trim().length > 0) return true
  return !collapsedNodes.value[nodeKey]
}

function toggleNode(nodeKey: string) {
  collapsedNodes.value[nodeKey] = !collapsedNodes.value[nodeKey]
}

function getAllDbItems(dbNode: DatabaseNodeData): CatalogItem[] {
  const res: CatalogItem[] = []
  for (const schemaNode of dbNode.schemas) {
    res.push(...schemaNode.items)
  }
  return res
}

// Reset selected index when query changes
watch(query, () => {
  selectedIndex.value = 0
})

function selectNext() {
  if (filteredItems.value.length > 0) {
    selectedIndex.value = (selectedIndex.value + 1) % filteredItems.value.length
  }
}

function selectPrev() {
  if (filteredItems.value.length > 0) {
    selectedIndex.value = (selectedIndex.value - 1 + filteredItems.value.length) % filteredItems.value.length
  }
}

function confirmSelection() {
  if (filteredItems.value.length > 0 && selectedIndex.value < filteredItems.value.length) {
    handleAdd(filteredItems.value[selectedIndex.value])
  }
}

async function handleAdd(item: CatalogItem, shouldClose: boolean = true) {
  try {
    await workspaceStore.addObject(item)
    uiStore.addNotification({
      type: 'success',
      title: 'Added to Tree',
      message: `"${item.name}" added to workspace root.`,
    })
    if (shouldClose) {
      emit('close')
    }
  } catch (e: any) {
    uiStore.addNotification({
      type: 'error',
      title: 'Failed to Add',
      message: e.message || String(e),
    })
  }
}

async function handleAddAll(itemsList: CatalogItem[]) {
  let addedCount = 0
  for (const item of itemsList) {
    try {
      await workspaceStore.addObject(item)
      addedCount++
    } catch (e) {
      // ignore individual failures
    }
  }
  uiStore.addNotification({
    type: 'success',
    title: 'Added to Workspace',
    message: `Added ${addedCount} object(s) to workspace.`,
  })
}

// Icons and style helpers
function iconComponent(type: string) {
  const icons: Record<string, any> = {
    table: TableIcon,
    view: ViewIcon,
    function: FunctionIcon,
  }
  return icons[type] || TableIcon
}

function iconColor(type: string) {
  const colors: Record<string, string> = {
    table: 'text-accent-green',
    view: 'text-accent-amber',
    function: 'text-purple-400',
  }
  return colors[type] || 'text-text-secondary'
}

function typeBadgeClass(type: string) {
  const classes: Record<string, string> = {
    table: 'bg-accent-green/10 text-accent-green',
    view: 'bg-accent-amber/10 text-accent-amber',
    function: 'bg-purple-500/10 text-purple-400',
  }
  return classes[type] || 'bg-navy-tertiary text-text-secondary'
}

// Simple icons
function TableIcon() {
  return h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2', class: 'w-4 h-4' }, [
    h('rect', { x: '3', y: '3', width: '18', height: '18', rx: '2' }),
    h('path', { d: 'M3 9h18M3 15h18M9 3v18' }),
  ])
}

function ViewIcon() {
  return h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2', class: 'w-4 h-4' }, [
    h('path', { d: 'M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z' }),
    h('circle', { cx: '12', cy: '12', r: '3' }),
  ])
}

function FunctionIcon() {
  return h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2', class: 'w-4 h-4' }, [
    h('path', { d: 'M8 3H5a2 2 0 0 0-2 2v14c0 1.1.9 2 2 2h3' }),
    h('path', { d: 'M16 3h3a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-3' }),
    h('path', { d: 'M9 12h6M12 9v6' }),
  ])
}
</script>

