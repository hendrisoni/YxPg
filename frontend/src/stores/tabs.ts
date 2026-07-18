import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import type { Tab } from '../types'
import { useUiStore } from './ui'

let tabCounter = 0
let queryCounter = 0

// Initialize tabCounter and queryCounter to prevent collisions
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
      if (tab.type === 'query') {
        const titleNum = parseInt(tab.title.replace(/\D/g, ''))
        if (!isNaN(titleNum) && titleNum > queryCounter) {
          queryCounter = titleNum
        }
      }
    }
  } catch (e) {
    // Ignore parse errors on startup
  }
}

const savedQueryCounter = parseInt(localStorage.getItem('tabs:query_counter') || '0')
const normalizedSavedQueryCounter = savedQueryCounter > 15 ? (savedQueryCounter % 15 || 15) : savedQueryCounter
if (normalizedSavedQueryCounter > queryCounter) {
  queryCounter = normalizedSavedQueryCounter
}

export const useTabsStore = defineStore('tabs', () => {
  const tabs = ref<Tab[]>(savedTabsStr ? JSON.parse(savedTabsStr) : [])
  const activeTabId = ref<string | null>(localStorage.getItem('tabs:active_tab_id'))

  // Create computed property that represents clean, serializable tabs state
  const serializableTabs = computed(() => {
    return tabs.value.map(tab => {
      const { data, ...rest } = tab
      let cleanData = undefined
      if (data) {
        const { results, queryResult, ...otherData } = data
        cleanData = otherData
      }
      return { ...rest, data: cleanData }
    })
  })

  // Persist state automatically using watchers on the serialized computed value
  watch(
    serializableTabs,
    (newTabsToSave) => {
      localStorage.setItem('tabs:open_tabs', JSON.stringify(newTabsToSave))
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

  function createTab(type: Tab['type'], options: Partial<Tab> = {}): Tab | null {
    if (type === 'query') {
      const queryTabsCount = tabs.value.filter(t => t.type === 'query').length
      if (queryTabsCount >= 25) {
        const uiStore = useUiStore()
        uiStore.addNotification({
          type: 'error',
          title: 'Limit Reached',
          message: 'Cannot add tab. Maximum of 25 query tabs allowed.'
        })
        return null
      }
    }

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
      case 'query': {
        if (!options.title) {
          queryCounter++
          if (queryCounter > 15) {
            queryCounter = 1
          }
          localStorage.setItem('tabs:query_counter', queryCounter.toString())
          return queryCounter.toString()
        }
        return options.title
      }
      case 'table': return options.table || 'Table'
      case 'builder': return options.title || `Builder ${tabCounter}`
      case 'ddl': return options.title || options.table || 'DDL'
      case 'home': return 'Home'
      case 'log': return 'Query Log'
      case 'backup': return 'Backup'
      case 'maintenance': return 'Maintenance'
      case 'referential': return 'Referential Integrity'
      case 'functions-triggers': return 'Functions & Triggers'
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
