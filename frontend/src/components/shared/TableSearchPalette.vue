<template>
  <Modal :show="show" title="Database Objects Search (Ctrl+K)" @close="$emit('close')" size="lg">
    <div class="flex flex-col h-[450px]">
      <!-- Search Input -->
      <div class="relative mb-3">
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

      <!-- Items List -->
      <div class="flex-1 overflow-y-auto min-h-0 space-y-1 pr-1">
        <div v-if="loading" class="flex flex-col items-center justify-center h-full text-text-muted text-xs">
          <div class="inline-block w-6 h-6 border-2 border-teal-accent border-t-transparent rounded-full animate-spin mb-2"></div>
          Loading database objects...
        </div>
        <div v-else-if="filteredItems.length === 0" class="flex items-center justify-center h-full text-text-muted text-xs">
          No tables, views, or functions found matching "{{ query }}"
        </div>
        <div
          v-else
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
          </div>
        </div>
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

async function handleAdd(item: CatalogItem) {
  try {
    await workspaceStore.addObject(item)
    uiStore.addNotification({
      type: 'success',
      title: 'Added to Tree',
      message: `"${item.name}" added to workspace root.`,
    })
    emit('close')
  } catch (e: any) {
    uiStore.addNotification({
      type: 'error',
      title: 'Failed to Add',
      message: e.message || String(e),
    })
  }
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
