import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import type { Tab } from '../types'

let tabCounter = 0

// Initialize tabCounter to prevent collisions
const savedTabsStr = localStorage.getItem('tabs:open_tabs')
if (savedTabsStr) {
  try {
    const savedTabs = JSON.parse(savedTabsStr) as Tab[]
    for (const tab of savedTabs) {
      const parts = tab.id.split('_')
      if (parts.length >= 2) {
        const count = parseInt(parts[1])
        if (!isNaN(count) && count > tabCounter) {
          tabCounter = count
        }
      }
    }
  } catch (e) {
    // Ignore parse errors on startup
  }
}

export const useTabsStore = defineStore('tabs', () => {
  const tabs = ref<Tab[]>(savedTabsStr ? JSON.parse(savedTabsStr) : [])
  const activeTabId = ref<string | null>(localStorage.getItem('tabs:active_tab_id'))

  // Persist state automatically using watchers
  watch(
    tabs,
    (newTabs) => {
      localStorage.setItem('tabs:open_tabs', JSON.stringify(newTabs))
    },
    { deep: true }
  )

  watch(
    activeTabId,
    (newId) => {
      if (newId) {
        localStorage.setItem('tabs:active_tab_id', newId)
      } else {
        localStorage.removeItem('tabs:active_tab_id')
      }
    }
  )

  const activeTab = computed(() => {
    if (!activeTabId.value) return null
    return tabs.value.find(t => t.id === activeTabId.value) || null
  })

  function createTab(type: Tab['type'], options: Partial<Tab> = {}): Tab {
    tabCounter++
    const id = `${type}_${tabCounter}_${Date.now()}`
    const tab: Tab = {
      id,
      title: options.title || getDefaultTitle(type, options),
      type,
      connectionId: options.connectionId,
      schema: options.schema,
      table: options.table,
      sql: options.sql || getDefaultSql(type, options),
      modified: false,
      data: options.data,
    }
    tabs.value.push(tab)
    activeTabId.value = id
    return tab
  }

  function closeTab(id: string) {
    const index = tabs.value.findIndex(t => t.id === id)
    if (index === -1) return

    tabs.value.splice(index, 1)

    if (activeTabId.value === id) {
      if (tabs.value.length > 0) {
        const newIndex = Math.min(index, tabs.value.length - 1)
        activeTabId.value = tabs.value[newIndex].id
      } else {
        activeTabId.value = null
      }
    }
  }

  function setActiveTab(id: string) {
    activeTabId.value = id
  }

  function updateTab(id: string, updates: Partial<Tab>) {
    const tab = tabs.value.find(t => t.id === id)
    if (tab) {
      Object.assign(tab, updates)
    }
  }

  function markModified(id: string, modified: boolean = true) {
    const tab = tabs.value.find(t => t.id === id)
    if (tab) {
      tab.modified = modified
    }
  }

  function getDefaultTitle(type: Tab['type'], options: Partial<Tab>): string {
    switch (type) {
      case 'query': return options.title || `Query ${tabCounter}`
      case 'table': return options.table || 'Table'
      case 'builder': return options.title || `Builder ${tabCounter}`
      case 'ddl': return options.title || options.table || 'DDL'
      case 'home': return 'Home'
      case 'log': return 'Query Log'
      default: return 'Untitled'
    }
  }

  function getDefaultSql(type: Tab['type'], options: Partial<Tab>): string {
    if (type === 'query' && options.table) {
      return `SELECT * FROM ${options.schema || 'public'}.${options.table}`
    }
    return ''
  }

  // Keyboard shortcut support
  function closeActiveTab() {
    if (activeTabId.value) {
      closeTab(activeTabId.value)
    }
  }

  function nextTab() {
    if (tabs.value.length <= 1) return
    const currentIdx = tabs.value.findIndex(t => t.id === activeTabId.value)
    const nextIdx = (currentIdx + 1) % tabs.value.length
    activeTabId.value = tabs.value[nextIdx].id
  }

  function prevTab() {
    if (tabs.value.length <= 1) return
    const currentIdx = tabs.value.findIndex(t => t.id === activeTabId.value)
    const prevIdx = (currentIdx - 1 + tabs.value.length) % tabs.value.length
    activeTabId.value = tabs.value[prevIdx].id
  }

  function closeAllTabs() {
    tabs.value = []
    activeTabId.value = null
  }

  function closeOtherTabs(id: string) {
    const keepTab = tabs.value.find(t => t.id === id)
    if (keepTab) {
      tabs.value = [keepTab]
      activeTabId.value = id
    } else {
      tabs.value = []
      activeTabId.value = null
    }
  }

  return {
    tabs,
    activeTabId,
    activeTab,
    createTab,
    closeTab,
    closeAllTabs,
    closeOtherTabs,
    setActiveTab,
    updateTab,
    markModified,
    closeActiveTab,
    nextTab,
    prevTab,
  }
})
