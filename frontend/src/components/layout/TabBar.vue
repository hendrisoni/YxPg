<template>
  <nav class="flex items-start bg-navy-secondary border-b border-navy-border min-h-9 h-auto select-none">
    <!-- Tabs wrap wrapper -->
    <div class="flex-1 flex flex-wrap items-center relative">
      <!-- Active Tab Neon Slider Line -->
      <div
        class="absolute h-[3px] bg-teal-accent shadow-[0_0_12px_#00C9A7,0_0_4px_#00C9A7] z-20 transition-all duration-300 ease-in-out"
        :style="sliderStyle"
      >
        <!-- Soft gradient glow cast below the indicator -->
        <div class="absolute top-[3px] left-0 right-0 h-4 bg-gradient-to-b from-teal-accent/15 to-transparent pointer-events-none blur-[2px]"></div>
      </div>

      <!-- Tabs -->
      <div
        v-for="tab in tabsStore.tabs"
        :key="tab.id"
        :id="`tab-item-${tab.id}`"
        @click="tabsStore.setActiveTab(tab.id)"
        @contextmenu.prevent="showTabContextMenu($event, tab)"
        class="group flex items-center gap-1.5 px-3 h-[35px] border-r border-b border-navy-border cursor-pointer transition-all relative min-w-0"
        :class="[
          tab.id === tabsStore.activeTabId
            ? 'bg-navy-primary text-text-primary font-medium'
            : 'bg-navy-secondary text-text-secondary hover:bg-navy-hover hover:text-text-primary'
        ]"
      >
        <!-- Tab icon -->
        <svg v-if="tab.type === 'query'" class="w-3.5 h-3.5 flex-shrink-0 text-accent-blue" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="16 18 22 12 16 6" /><polyline points="8 6 2 12 8 18" />
        </svg>
        <svg v-else-if="tab.type === 'table'" class="w-3.5 h-3.5 flex-shrink-0 text-accent-green" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="3" width="18" height="18" rx="2" /><path d="M3 9h18M3 15h18M9 3v18" />
        </svg>
        <svg v-else-if="tab.type === 'builder'" class="w-3.5 h-3.5 flex-shrink-0 text-orange-500 drop-shadow-[0_0_3px_rgba(249,115,22,0.85)]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <!-- Top Left Table -->
          <rect x="2" y="3" width="8" height="7" rx="1" />
          <path d="M2 6h8" />
          
          <!-- Top Right Table -->
          <rect x="14" y="3" width="8" height="7" rx="1" />
          <path d="M14 6h8" />
          
          <!-- Bottom Table -->
          <rect x="8" y="14" width="8" height="7" rx="1" />
          <path d="M8 17h8" />
          
          <!-- Connection Lines -->
          <path d="M10 6h4" />
          <path d="M6 10v2h12v-2" />
          <path d="M12 12v2" />
        </svg>
        <svg v-else-if="tab.type === 'ddl'" class="w-3.5 h-3.5 flex-shrink-0 text-accent-red" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
          <polyline points="14 2 14 8 20 8" />
        </svg>
        <svg v-else-if="tab.type === 'backup'" class="w-3.5 h-3.5 flex-shrink-0 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
          <polyline points="17 8 12 3 7 8" />
          <line x1="12" y1="3" x2="12" y2="15" />
        </svg>
        <svg v-else-if="tab.type === 'functions-triggers'" class="w-3.5 h-3.5 flex-shrink-0 text-purple-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2" />
        </svg>

        <!-- Modified indicator -->
        <span v-if="tab.modified" class="w-1.5 h-1.5 rounded-full bg-accent-amber flex-shrink-0"></span>

        <!-- Connection Server Label -->
        <span 
          v-if="getConnectionName(tab.connectionId)" 
          class="text-[9px] px-1 py-0.5 rounded leading-none font-bold select-none flex-shrink-0 border border-transparent"
          :style="{
            backgroundColor: getConnectionColor(tab.connectionId) ? getConnectionColor(tab.connectionId) + '15' : 'rgba(59, 130, 246, 0.15)',
            borderColor: getConnectionColor(tab.connectionId) ? getConnectionColor(tab.connectionId) + '30' : 'rgba(59, 130, 246, 0.3)',
            color: getConnectionColor(tab.connectionId) || '#3B82F6'
          }"
        >
          {{ getConnectionName(tab.connectionId) }}
        </span>

        <!-- Tab title -->
        <span class="text-xs whitespace-nowrap overflow-hidden text-ellipsis max-w-[120px]">
          {{ tab.title }}
        </span>

        <!-- Close button -->
        <button
          @click.stop="tabsStore.closeTab(tab.id)"
          class="ml-1 p-0.5 rounded opacity-0 group-hover:opacity-100 hover:bg-navy-hover text-text-muted hover:text-text-primary transition-all"
        >
          <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M18 6 6 18M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- New Tab Button -->
      <button
        @click="$emit('newTab')"
        class="flex items-center justify-center w-9 h-[35px] text-text-muted hover:text-text-primary hover:bg-navy-hover transition-colors border-r border-b border-navy-border flex-shrink-0"
        title="New Tab (Ctrl+T)"
      >
        <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 5v14M5 12h14" />
        </svg>
      </button>
    </div>

    <!-- Builder static action button on the right -->
    <button
      @click="$emit('openBuilder')"
      class="flex items-center gap-1.5 px-3 h-[35px] text-xs text-text-secondary hover:text-orange-400 hover:bg-orange-500/5 hover:shadow-[inset_0_-2px_0_0_#f97316] transition-all border-l border-b border-navy-border flex-shrink-0"
      title="Open Builder"
    >
      <svg class="w-3.5 h-3.5 text-orange-500 drop-shadow-[0_0_3px_rgba(249,115,22,0.85)]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <!-- Top Left Table -->
        <rect x="2" y="3" width="8" height="7" rx="1" />
        <path d="M2 6h8" />
        
        <!-- Top Right Table -->
        <rect x="14" y="3" width="8" height="7" rx="1" />
        <path d="M14 6h8" />
        
        <!-- Bottom Table -->
        <rect x="8" y="14" width="8" height="7" rx="1" />
        <path d="M8 17h8" />
        
        <!-- Connection Lines -->
        <path d="M10 6h4" />
        <path d="M6 10v2h12v-2" />
        <path d="M12 12v2" />
      </svg>
      <span class="hover:drop-shadow-[0_0_4px_rgba(249,115,22,0.5)] transition-all">Builder</span>
    </button>
  </nav>

  <!-- Tab Context Menu -->
  <ContextMenu
    :show="showMenu"
    :x="menuX"
    :y="menuY"
    :items="contextMenuItems"
    @close="showMenu = false"
    @select="handleMenuSelect"
  />
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useTabsStore } from '../../stores/tabs'
import { useUiStore } from '../../stores/ui'
import { useConnectionsStore } from '../../stores/connections'
import type { Tab } from '../../types'
import ContextMenu from '../shared/ContextMenu.vue'
import type { ContextMenuItem } from '../shared/ContextMenu.vue'

const emit = defineEmits(['newTab', 'openBuilder'])

const tabsStore = useTabsStore()
const uiStore = useUiStore()
const connectionsStore = useConnectionsStore()

function getConnectionName(connId?: string): string {
  if (!connId) return ''
  const conn = connectionsStore.connections.find(c => c.id === connId)
  return conn ? conn.name : ''
}

function getConnectionColor(connId?: string): string {
  if (!connId) return ''
  const conn = connectionsStore.connections.find(c => c.id === connId)
  return conn?.color || ''
}

const showMenu = ref(false)
const menuX = ref(0)
const menuY = ref(0)
const targetTab = ref<Tab | null>(null)

const contextMenuItems = computed<ContextMenuItem[]>(() => {
  if (!targetTab.value) return []
  return [
    { label: 'Close Tab', action: 'close', shortcut: 'Ctrl+W' },
    { label: 'Close Others', action: 'close_others' },
    { label: 'Close All', action: 'close_all', danger: true }
  ]
})

function showTabContextMenu(event: MouseEvent, tab: Tab) {
  targetTab.value = tab
  menuX.value = event.clientX
  menuY.value = event.clientY
  showMenu.value = true
}

function handleMenuSelect(action: string) {
  if (!targetTab.value) return
  if (action === 'close') {
    tabsStore.closeTab(targetTab.value.id)
  } else if (action === 'close_others') {
    tabsStore.closeOtherTabs(targetTab.value.id)
  } else if (action === 'close_all') {
    tabsStore.closeAllTabs()
  }
  showMenu.value = false
}

// Neon slider position states and logic
const sliderStyle = ref({
  left: '0px',
  top: '0px',
  width: '0px',
  opacity: '0'
})

function updateSlider() {
  nextTick(() => {
    const activeTabId = tabsStore.activeTabId
    if (!activeTabId) {
      sliderStyle.value = { left: '0px', top: '0px', width: '0px', opacity: '0' }
      return
    }

    const activeEl = document.getElementById(`tab-item-${activeTabId}`)
    if (activeEl) {
      const offsetLeft = activeEl.offsetLeft
      const offsetTop = activeEl.offsetTop
      const width = activeEl.offsetWidth
      sliderStyle.value = {
        left: `${offsetLeft}px`,
        top: `${offsetTop}px`,
        width: `${width}px`,
        opacity: '1'
      }
    } else {
      sliderStyle.value = { left: '0px', top: '0px', width: '0px', opacity: '0' }
    }
  })
}

// Watchers
watch(() => tabsStore.activeTabId, () => {
  updateSlider()
}, { immediate: true })

watch(() => tabsStore.tabs, () => {
  updateSlider()
}, { deep: true })

onMounted(() => {
  updateSlider()
  window.addEventListener('resize', updateSlider)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateSlider)
})
</script>
