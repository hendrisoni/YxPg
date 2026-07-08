<template>
  <nav class="flex items-center bg-navy-secondary border-b border-navy-border h-9 select-none">
    <!-- Tabs scrollable wrapper -->
    <div class="flex-1 flex items-center h-full overflow-x-auto no-scrollbar">
      <!-- Tabs -->
      <div class="flex items-center h-full">
        <div
          v-for="tab in tabsStore.tabs"
          :key="tab.id"
          @click="tabsStore.setActiveTab(tab.id)"
          @contextmenu.prevent="showTabContextMenu($event, tab)"
          class="group flex items-center gap-1.5 px-3 h-full border-r border-navy-border cursor-pointer transition-colors relative min-w-0"
          :class="[
            tab.id === tabsStore.activeTabId
              ? 'bg-navy-primary text-text-primary'
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
          <svg v-else-if="tab.type === 'builder'" class="w-3.5 h-3.5 flex-shrink-0 text-accent-amber" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3" /><path d="M12 3v6m0 6v6m-9-9h6m6 0h6" />
          </svg>
          <svg v-else-if="tab.type === 'ddl'" class="w-3.5 h-3.5 flex-shrink-0 text-accent-red" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
            <polyline points="14 2 14 8 20 8" />
          </svg>

          <!-- Modified indicator -->
          <span v-if="tab.modified" class="w-1.5 h-1.5 rounded-full bg-accent-amber flex-shrink-0"></span>

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
      </div>

      <!-- New Tab Button -->
      <button
        @click="$emit('newTab')"
        class="flex items-center justify-center w-9 h-full text-text-muted hover:text-text-primary hover:bg-navy-hover transition-colors border-r border-navy-border flex-shrink-0"
        title="New Tab (Ctrl+T)"
      >
        <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 5v14M5 12h14" />
        </svg>
      </button>
    </div>

    <!-- Query Builder static action button on the right -->
    <button
      @click="$emit('openBuilder')"
      class="flex items-center gap-1.5 px-3 h-full text-xs text-text-secondary hover:text-text-primary hover:bg-navy-hover transition-colors border-l border-navy-border flex-shrink-0"
      title="Open Query Builder"
    >
      <svg class="w-3.5 h-3.5 text-accent-amber" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="3" /><path d="M12 3v6m0 6v6m-9-9h6m6 0h6" />
      </svg>
      <span>Query Builder</span>
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
import { ref, computed } from 'vue'
import { useTabsStore } from '../../stores/tabs'
import { useUiStore } from '../../stores/ui'
import type { Tab } from '../../types'
import ContextMenu from '../shared/ContextMenu.vue'
import type { ContextMenuItem } from '../shared/ContextMenu.vue'

const emit = defineEmits(['newTab', 'openBuilder'])

const tabsStore = useTabsStore()
const uiStore = useUiStore()

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
</script>
