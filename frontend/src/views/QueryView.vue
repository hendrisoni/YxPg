<template>
  <div class="h-full flex flex-col overflow-hidden">
    <!-- Toolbar -->
    <div class="flex items-center gap-2 px-3 py-1.5 border-b border-navy-border bg-navy-secondary">
      <button @click="runQuery" :disabled="!isTargetConnected || isExecuting"
        class="flex items-center gap-1.5 px-3 py-1 text-xs bg-teal-accent text-navy-primary rounded font-medium hover:bg-teal-hover transition-colors disabled:opacity-50">
        <svg v-if="!isExecuting" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="currentColor">
          <polygon points="5 3 19 12 5 21 5 3" />
        </svg>
        <div v-else class="w-3.5 h-3.5 border-2 border-navy-primary border-t-transparent rounded-full animate-spin">
        </div>
        {{ isExecuting ? 'Running...' : 'Run' }}
      </button>

      <div class="w-px h-4 bg-navy-border"></div>

      <!-- Connection Selector -->
      <div class="flex items-center gap-1.5">
        <span class="text-[10px] text-text-secondary">Con:</span>
        <select v-model="selectedConnectionId"
          class="bg-navy-tertiary border border-navy-border text-xs text-text-primary rounded px-2 py-0.5 focus:border-teal-accent focus:outline-none cursor-pointer">
          <option value="">Default ({{ currentConnectionName }})</option>
          <option v-for="conn in connectionsStore.connections" :key="conn.id" :value="conn.id">
            {{ conn.name }} {{ connectionsStore.activeConnections.includes(conn.id) ? '●' : '○' }}
          </option>
        </select>

        <button @click="showCustomRunModal = true" :disabled="!isTargetConnected || isExecuting"
          class="flex items-center gap-1.5 px-2.5 py-1.5 text-[11px] bg-navy-secondary border border-navy-border text-text-primary rounded hover:bg-navy-hover transition-colors disabled:opacity-50"
          title="Run query with :date_from and :date_to replacement">
          <svg class="w-3.5 h-3.5 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor"
            stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
            <line x1="16" y1="2" x2="16" y2="6" />
            <line x1="8" y1="2" x2="8" y2="6" />
            <line x1="3" y1="10" x2="21" y2="10" />
          </svg>
          Date
        </button>

        <button @click="openQueryLogTab"
          class="flex items-center gap-1.5 px-2.5 py-1.5 text-[11px] bg-navy-secondary border border-navy-border text-text-primary rounded hover:bg-navy-hover transition-colors"
          title="Lihat 7 query terakhir">
          <svg class="w-3.5 h-3.5 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor"
            stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 8v4l3 3" />
            <circle cx="12" cy="12" r="10" />
          </svg>
        </button>
      </div>

      <div class="flex-1"></div>

      <!-- Query stats -->
      <div v-if="results.length === 1" class="flex items-center gap-2 text-xs text-text-muted">
        <span>{{ results[0].duration_ms || 0 }}ms</span>
        <span v-if="results[0].row_count > 0">• {{ results[0].row_count }} rows</span>
        <span v-if="results[0].error" class="text-accent-red">• Error</span>
      </div>
      <div v-else-if="results.length > 1" class="flex items-center gap-2 text-xs text-text-muted">
        <span>{{ results.length }} results</span>
      </div>
    </div>

    <!-- Editor + Results split -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Editor -->
      <div class="flex-1 min-h-0" :style="{ flexBasis: editorFlexBasis }" @dragover.prevent @drop="handleEditorDrop"
        @contextmenu.prevent="showContextMenu">
        <div ref="editorContainer" class="h-full w-full overflow-hidden"></div>
      </div>

      <!-- Context Menu -->
      <Teleport to="body">
        <div v-if="contextMenu.visible" class="ctx-menu" :style="{ top: contextMenu.y + 'px', left: contextMenu.x + 'px' }"
          @click.stop>
          <button class="ctx-menu-item" @click="ctxCopyAll">
            <svg class="ctx-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
              <path d="M14 13h4M14 17h4" />
            </svg>
            <span>Copy All</span>
            <span class="ctx-shortcut">Ctrl+Shift+C</span>
          </button>
          <div class="ctx-separator"></div>
          <button class="ctx-menu-item" @click="ctxCut">
            <svg class="ctx-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="6" cy="6" r="3" />
              <circle cx="6" cy="18" r="3" />
              <line x1="20" y1="4" x2="8.12" y2="15.88" />
              <line x1="14.47" y1="14.48" x2="20" y2="20" />
              <line x1="8.12" y1="8.12" x2="12" y2="12" />
            </svg>
            <span>Cut</span>
            <span class="ctx-shortcut">Ctrl+X</span>
          </button>
          <button class="ctx-menu-item" @click="ctxCopy">
            <svg class="ctx-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
            </svg>
            <span>Copy</span>
            <span class="ctx-shortcut">Ctrl+C</span>
          </button>
          <button class="ctx-menu-item" @click="ctxPaste">
            <svg class="ctx-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2" />
              <rect x="8" y="2" width="8" height="4" rx="1" ry="1" />
            </svg>
            <span>Paste</span>
            <span class="ctx-shortcut">Ctrl+V</span>
          </button>
          <div class="ctx-separator"></div>
          <button class="ctx-menu-item" @click="ctxFormatSql">
            <svg class="ctx-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 10H3M21 6H3M21 14H3M21 18H3" />
            </svg>
            <span>Format SQL</span>
            <span class="ctx-shortcut">Ctrl+Shift+F</span>
          </button>
          <button class="ctx-menu-item" @click="ctxReplace">
            <svg class="ctx-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8" />
              <path d="m21 21-4.35-4.35" />
              <path d="M8 11h6M11 8v6" />
            </svg>
            <span>Find & Replace</span>
            <span class="ctx-shortcut">Ctrl+H</span>
          </button>
          <div class="ctx-separator"></div>
          <button class="ctx-menu-item ctx-danger" @click="ctxClearEditor">
            <svg class="ctx-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 6h18M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
            </svg>
            <span>Clear Editor</span>
            <span class="ctx-shortcut">Ctrl+Shift+X</span>
          </button>
        </div>
      </Teleport>

      <!-- Resize handle -->
      <div class="h-1 bg-navy-border hover:bg-teal-accent cursor-row-resize transition-colors flex-shrink-0"
        @mousedown="startResize"></div>

      <!-- Results Area (Stacked for multiple) -->
      <div class="flex flex-col border-t border-navy-border overflow-hidden" :style="{ height: resultHeight + 'px' }">
        <div v-if="results.length > 0" class="flex-1 flex flex-col min-h-0">
          <!-- Stack each result -->
          <div 
            v-for="(res, idx) in results" 
            :key="idx"
            class="flex flex-col flex-1 min-h-[100px]"
            :class="{ 'border-b-4 border-navy-border/50': idx < results.length - 1 }"
          >
            <!-- Minimal header for multi-results to show which is which -->
            <div v-if="results.length > 1" class="px-2 py-1 bg-navy-tertiary border-b border-navy-border flex items-center justify-between flex-shrink-0 z-10">
              <span class="text-[10px] font-mono text-teal-accent/70 truncate max-w-[80%]">{{ res.raw_sql || `Result ${idx + 1}` }}</span>
              <div class="flex items-center gap-2 text-[10px] text-text-muted">
                <span v-if="res.error" class="text-accent-red">Error</span>
                <span v-else>{{ res.row_count }} rows</span>
                <span>{{ res.duration_ms }}ms</span>
              </div>
            </div>
            
            <div class="flex-1 min-h-0 overflow-hidden relative">
              <ResultGrid 
                :result="res" 
                :executed-sql="res.raw_sql || lastExecutedSql" 
                @refresh="runQuery" 
              />
            </div>
          </div>
        </div>
        
        <div v-else class="h-full flex flex-col items-center justify-center text-text-muted opacity-50">
          <svg class="w-8 h-8 mb-2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <rect x="3" y="3" width="18" height="18" rx="2" />
            <path d="M3 9h18M3 15h18M9 3v18" />
          </svg>
          <p class="text-xs">No result</p>
        </div>
      </div>
    </div>

    <!-- Custom Run parameters modal -->
    <QueryWinDate v-if="showCustomRunModal" :show="showCustomRunModal" @close="showCustomRunModal = false"
      @execute="handleCustomRun" />

  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import QueryWinDate from '../components/query/QueryWinDate.vue'
import { useConnectionsStore } from '../stores/connections'
import { useSchemaStore } from '../stores/schema'
import { useTabsStore } from '../stores/tabs'
import { useUiStore } from '../stores/ui'
import { useWorkspaceStore } from '../stores/workspace'
import { formatSQL } from '../utils/sql-formatter'
import type { Tab, QueryResult } from '../types'
import ResultGrid from '../components/query/ResultGrid.vue'

const props = defineProps<{
  tab: Tab
}>()

const connectionsStore = useConnectionsStore()
const tabsStore = useTabsStore()
const uiStore = useUiStore()
const workspaceStore = useWorkspaceStore()
const schemaStore = useSchemaStore()

const editorContainer = ref<HTMLElement | null>(null)
const isExecuting = ref(false)
const results = ref<QueryResult[]>([])
const activeResultIndex = ref(0)
const lastExecutedSql = ref('')

const showCustomRunModal = ref(false)

// Context menu state
const contextMenu = reactive({ visible: false, x: 0, y: 0 })

function showContextMenu(e: MouseEvent) {
  contextMenu.x = e.clientX
  contextMenu.y = e.clientY
  contextMenu.visible = true
  document.addEventListener('click', closeContextMenu, { once: true })
  document.addEventListener('keydown', closeOnEscape)
}

function closeContextMenu() {
  contextMenu.visible = false
  document.removeEventListener('keydown', closeOnEscape)
}

function closeOnEscape(e: KeyboardEvent) {
  if (e.key === 'Escape') closeContextMenu()
}

function ctxCut() {
  closeContextMenu()
  if (!editor) return
  const { from, to } = editor.state.selection.main
  const sel = editor.state.sliceDoc(from, to)
  if (sel) {
    navigator.clipboard.writeText(sel).then(() => {
      uiStore.addNotification({
        type: 'success',
        title: 'Cut',
        message: 'Selected text cut to clipboard',
        duration: 2000
      })
    })
    editor.dispatch({ changes: { from, to, insert: '' } })
  }
}

function ctxCopy() {
  closeContextMenu()
  if (!editor) return
  const sel = editor.state.sliceDoc(
    editor.state.selection.main.from,
    editor.state.selection.main.to
  )
  if (sel) {
    navigator.clipboard.writeText(sel).then(() => {
      uiStore.addNotification({
        type: 'success',
        title: 'Copied',
        message: 'Selected text copied to clipboard',
        duration: 2000
      })
    })
  }
}

function ctxCopyAll() {
  closeContextMenu()
  if (!editor) return
  const allText = editor.state.doc.toString()
  if (allText) {
    navigator.clipboard.writeText(allText).then(() => {
      uiStore.addNotification({
        type: 'success',
        title: 'Copied',
        message: 'All SQL text copied to clipboard',
        duration: 2000
      })
    })
  }
}

async function ctxPaste() {
  closeContextMenu()
  if (!editor) return
  const text = await navigator.clipboard.readText()
  if (text) {
    const { from, to } = editor.state.selection.main
    editor.dispatch({
      changes: { from, to, insert: text },
      selection: { anchor: from + text.length }
    })
  }
}

function ctxFormatSql() {
  closeContextMenu()
  formatSql()
}

async function ctxReplace() {
  closeContextMenu()
  if (!editor) return
  const { openSearchPanel } = await import('@codemirror/search')
  openSearchPanel(editor)
}

function ctxClearEditor() {
  closeContextMenu()
  clearEditor()
}

function saveToQueryLog(sql: string) {
  if (!sql || !sql.trim()) return
  let logs: string[] = []
  const stored = localStorage.getItem('yxpg:query_logs')
  if (stored) {
    try {
      logs = JSON.parse(stored)
    } catch {
      logs = []
    }
  }
  logs = logs.filter(q => q.trim() !== sql.trim())
  logs.unshift(sql)
  if (logs.length > 7) {
    logs = logs.slice(0, 7)
  }
  localStorage.setItem('yxpg:query_logs', JSON.stringify(logs))
}

function openQueryLogTab() {
  const existing = tabsStore.tabs.find(t => t.type === 'log')
  if (existing) {
    tabsStore.setActiveTab(existing.id)
  } else {
    tabsStore.createTab('log')
  }
}

function handleCustomRun(dates: { dateFrom: string, dateTo: string }) {
  let sql = editor?.state.doc.toString() || ''
  if (!sql.trim()) return

  const formattedFrom = `'${dates.dateFrom}'::date`
  const formattedTo = `'${dates.dateTo}'::date`

  sql = sql.replace(/:date_from\b/gi, formattedFrom)
  sql = sql.replace(/:date_to\b/gi, formattedTo)

  runQuery(sql)
  showCustomRunModal.value = false
}
const resultHeight = ref(250)
const editorFlexBasis = ref('calc(100% - 251px)')

const isTargetConnected = computed(() => {
  const connId = props.tab.connectionId || connectionsStore.currentConnectionId
  if (!connId) return false
  return connectionsStore.activeConnections.includes(connId)
})

const selectedConnectionId = computed({
  get: () => props.tab.connectionId || '',
  set: (val: string) => {
    tabsStore.updateTab(props.tab.id, { connectionId: val || undefined })
  }
})

const currentConnectionName = computed(() => {
  return connectionsStore.currentConnection?.name || 'None'
})

let editor: any = null
let editorResizeObserver: ResizeObserver | null = null
let resizeStartY = 0
let resizeStartHeight = 0

const handleRunActiveQuery = () => {
  if (tabsStore.activeTab?.id === props.tab.id) {
    runQuery()
  }
}

const handleFormatActiveQuery = () => {
  if (tabsStore.activeTab?.id === props.tab.id) {
    formatSql()
  }
}

const handleGlobalKeydown = (e: KeyboardEvent) => {
  if (tabsStore.activeTab?.id !== props.tab.id) return

  const activeEl = document.activeElement
  if (activeEl && (activeEl.tagName === 'INPUT' || activeEl.tagName === 'TEXTAREA' || activeEl.getAttribute('contenteditable') === 'true')) {
    if (!editorContainer.value?.contains(activeEl)) {
      return
    }
  }

  const isCtrl = e.ctrlKey || e.metaKey
  const isShift = e.shiftKey

  if (isCtrl && isShift && e.key.toLowerCase() === 'c') {
    e.preventDefault()
    e.stopPropagation()
    ctxCopyAll()
  } else if (isCtrl && isShift && e.key.toLowerCase() === 'x') {
    e.preventDefault()
    e.stopPropagation()
    ctxClearEditor()
  }
}

onMounted(async () => {
  await nextTick()
  if (document.fonts) {
    await document.fonts.ready
  }
  initEditor()
  window.addEventListener('run-active-query', handleRunActiveQuery)
  window.addEventListener('format-active-query', handleFormatActiveQuery)
  window.addEventListener('keydown', handleGlobalKeydown, true)
})

onUnmounted(() => {
  if (editorResizeObserver) {
    editorResizeObserver.disconnect()
    editorResizeObserver = null
  }
  if (editor) {
    editor.destroy()
  }
  window.removeEventListener('run-active-query', handleRunActiveQuery)
  window.removeEventListener('format-active-query', handleFormatActiveQuery)
  window.removeEventListener('keydown', handleGlobalKeydown, true)
})

async function initEditor() {
  if (!editorContainer.value) return

  // Dynamically import CodeMirror
  const { EditorView, lineNumbers, highlightActiveLineGutter, highlightSpecialChars, drawSelection, dropCursor, rectangularSelection, crosshairCursor, highlightActiveLine, keymap } = await import('@codemirror/view')
  const { foldGutter, indentOnInput, syntaxHighlighting, HighlightStyle, bracketMatching, foldKeymap } = await import('@codemirror/language')
  const { tags } = await import('@lezer/highlight')
  const { history, defaultKeymap, historyKeymap } = await import('@codemirror/commands')
  const { searchKeymap } = await import('@codemirror/search')
  const { closeBrackets, autocompletion, closeBracketsKeymap, completionKeymap } = await import('@codemirror/autocomplete')
  const { lintKeymap } = await import('@codemirror/lint')
  const { sql: sqlLang } = await import('@codemirror/lang-sql')
  const { EditorState } = await import('@codemirror/state')

  const darculaHighlightStyle = HighlightStyle.define([
    { tag: tags.keyword, color: '#CC7832', fontWeight: 'bold' },
    { tag: tags.modifier, color: '#CC7832', fontWeight: 'bold' },
    { tag: tags.string, color: '#6A8759' },
    { tag: tags.number, color: '#6897BB' },
    { tag: tags.integer, color: '#6897BB' },
    { tag: tags.float, color: '#6897BB' },
    { tag: tags.comment, color: '#808080', fontStyle: 'italic' },
    { tag: tags.lineComment, color: '#808080', fontStyle: 'italic' },
    { tag: tags.blockComment, color: '#808080', fontStyle: 'italic' },
    { tag: tags.variableName, color: '#A9B7C6' },
    { tag: tags.propertyName, color: '#A9B7C6' },
    { tag: tags.name, color: '#A9B7C6' },
    { tag: tags.operator, color: '#A9B7C6' },
    { tag: tags.punctuation, color: '#A9B7C6' },
    { tag: tags.bool, color: '#CC7832', fontWeight: 'bold' },
  ])

  const editorFontFamily = "'DM Mono', ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace"
  const editorSetup = [
    lineNumbers(),
    highlightActiveLineGutter(),
    highlightSpecialChars(),
    history(),
    foldGutter(),
    dropCursor(),
    EditorState.allowMultipleSelections.of(true),
    indentOnInput(),
    syntaxHighlighting(darculaHighlightStyle),
    bracketMatching(),
    closeBrackets(),
    autocompletion(),
    rectangularSelection(),
    crosshairCursor(),
    highlightActiveLine(),
    keymap.of([
      ...closeBracketsKeymap,
      ...defaultKeymap,
      ...searchKeymap,
      ...historyKeymap,
      ...foldKeymap,
      ...completionKeymap,
      ...lintKeymap,
    ]),
  ]

  const darkTheme = EditorView.theme({
    '&': {
      backgroundColor: '#151515',
      color: '#A9B7C6',
      height: '100%',
    },
    '.cm-scroller': {
      fontFamily: editorFontFamily,
      fontSize: '14px',
      lineHeight: '22px',
    },
    '.cm-content': {
      caretColor: '#BBBBBB',
      fontFamily: editorFontFamily,
      fontSize: '14px',
      lineHeight: '22px',
      padding: '10px 0',
    },
    '.cm-line': {
      lineHeight: '22px',
      padding: '0 6px',
    },
    '.cm-cursor': {
      borderLeftColor: '#BBBBBB',
    },
    '&.cm-focused .cm-cursor': {
      borderLeftColor: '#BBBBBB',
    },
    '.cm-activeLine': {
      backgroundColor: '#222222',
    },
    '.cm-selectionBackground': {
      backgroundColor: '#1e3a8a !important',
    },
    '&.cm-focused .cm-selectionBackground': {
      backgroundColor: '#1e3a8a !important',
    },
    '.cm-selectionMatch, .cm-selectionMatch-main': {
      backgroundColor: 'rgba(30, 58, 138, 0.4) !important',
    },
    '.cm-gutters': {
      backgroundColor: '#151515',
      color: '#858585',
      borderRight: 'none',
    },
    '.cm-activeLineGutter': {
      backgroundColor: '#222222',
      color: '#C6C6C6',
    },
    '.cm-foldPlaceholder': {
      backgroundColor: '#2D2D2D',
      color: '#808080',
      border: 'none',
    },
    '.cm-panels': {
      backgroundColor: '#151515',
      borderBottom: '1px solid #222222',
    },
    '.cm-panels-bottom': {
      borderTop: '1px solid #222222',
      borderBottom: 'none',
    },
    '.cm-search': {
      padding: '6px 12px',
      display: 'flex',
      flexWrap: 'wrap',
      alignItems: 'center',
      gap: '8px',
    },
    '.cm-search .cm-textfield': {
      backgroundColor: '#222222',
      border: '1px solid #2A2A2A',
      color: '#CCCCCC',
      borderRadius: '4px',
      padding: '4px 8px',
      fontSize: '12px',
      outline: 'none',
      width: '160px',
    },
    '.cm-search .cm-textfield:focus': {
      borderColor: '#CC7832',
      boxShadow: '0 0 0 1px #CC7832',
    },
    '.cm-search button': {
      background: '#2A2A2A',
      backgroundImage: 'none',
      color: '#CCCCCC',
      border: 'none',
      borderRadius: '4px',
      width: '26px',
      height: '26px',
      display: 'inline-flex',
      alignItems: 'center',
      justifyContent: 'center',
      cursor: 'pointer',
      transition: 'all 0.15s ease',
      fontSize: '0',
      padding: '0',
    },
    '.cm-search button:hover': {
      background: '#CC7832',
      backgroundImage: 'none',
      color: '#151515',
    },
    '.cm-search button:active': {
      background: '#B36224',
      backgroundImage: 'none',
    },
    '.cm-search button[name=next]::before': {
      content: '""',
      display: 'inline-block',
      width: '14px',
      height: '14px',
      backgroundColor: 'currentColor',
      mask: "url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='m6 9 6 6 6-6'/%3E%3C/svg%3E\") no-repeat center / contain",
      WebkitMask: "url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='m6 9 6 6 6-6'/%3E%3C/svg%3E\") no-repeat center / contain",
    },
    '.cm-search button[name=prev]::before': {
      content: '""',
      display: 'inline-block',
      width: '14px',
      height: '14px',
      backgroundColor: 'currentColor',
      mask: "url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='m18 15-6-6-6 6'/%3E%3C/svg%3E\") no-repeat center / contain",
      WebkitMask: "url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='m18 15-6-6-6 6'/%3E%3C/svg%3E\") no-repeat center / contain",
    },
    '.cm-search button[name=select]::before': {
      content: '""',
      display: 'inline-block',
      width: '14px',
      height: '14px',
      backgroundColor: 'currentColor',
      mask: "url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Crect x='3' y='3' width='18' height='18' rx='2'/%3E%3Cpath d='M7 8h10M7 12h10M7 16h10'/%3E%3C/svg%3E\") no-repeat center / contain",
      WebkitMask: "url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Crect x='3' y='3' width='18' height='18' rx='2'/%3E%3Cpath d='M7 8h10M7 12h10M7 16h10'/%3E%3C/svg%3E\") no-repeat center / contain",
    },
    '.cm-search button[name=replace]::before': {
      content: '""',
      display: 'inline-block',
      width: '14px',
      height: '14px',
      backgroundColor: 'currentColor',
      mask: "url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8'/%3E%3Cpath d='M3 3v5h5M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16'/%3E%3Cpath d='M16 16h5v5'/%3E%3C/svg%3E\") no-repeat center / contain",
      WebkitMask: "url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8'/%3E%3Cpath d='M3 3v5h5M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16'/%3E%3Cpath d='M16 16h5v5'/%3E%3C/svg%3E\") no-repeat center / contain",
    },
    '.cm-search button[name=replaceAll]::before': {
      content: '""',
      display: 'inline-block',
      width: '14px',
      height: '14px',
      backgroundColor: 'currentColor',
      mask: "url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8'/%3E%3Cpath d='M3 3v5h5M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16'/%3E%3Cpath d='M16 16h5v5M7 12h10M12 7v10'/%3E%3C/svg%3E\") no-repeat center / contain",
      WebkitMask: "url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8'/%3E%3Cpath d='M3 3v5h5M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16'/%3E%3Cpath d='M16 16h5v5M7 12h10M12 7v10'/%3E%3C/svg%3E\") no-repeat center / contain",
    },
    '.cm-search label': {
      color: '#94A3B8',
      fontSize: '12px',
      display: 'flex',
      alignItems: 'center',
      gap: '4px',
      cursor: 'pointer',
      userSelect: 'none',
    },
    '.cm-search input[type=checkbox]': {
      accentColor: '#CC7832',
      cursor: 'pointer',
    },
    '.cm-search button[name=close]': {
      marginLeft: 'auto',
      background: 'transparent',
      backgroundImage: 'none',
      color: '#94A3B8',
      borderRadius: '4px',
      display: 'inline-flex',
      alignItems: 'center',
      justifyContent: 'center',
    },
    '.cm-search button[name=close]::before': {
      content: '""',
      display: 'inline-block',
      width: '14px',
      height: '14px',
      backgroundColor: 'currentColor',
      mask: "url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M18 6 6 18M6 6l12 12'/%3E%3C/svg%3E\") no-repeat center / contain",
      WebkitMask: "url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M18 6 6 18M6 6l12 12'/%3E%3C/svg%3E\") no-repeat center / contain",
    },
    '.cm-search button[name=close]:hover': {
      background: 'rgba(239, 68, 68, 0.2)',
      backgroundImage: 'none',
      color: '#EF4444',
    },
    '.cm-search br': {
      flexBasis: '100%',
      height: '0',
    },
  })

  const state = EditorState.create({
    doc: props.tab.sql || '',
    extensions: [
      editorSetup,
      sqlLang(),
      darkTheme,
      EditorView.updateListener.of((update) => {
        if (update.docChanged) {
          const doc = update.state.doc.toString()
          tabsStore.updateTab(props.tab.id, { sql: doc })
          tabsStore.markModified(props.tab.id, true)
        }
      }),
    ],
  })

  editor = new EditorView({
    state,
    parent: editorContainer.value,
  })

  editorResizeObserver = new ResizeObserver(() => {
    editor?.requestMeasure()
  })
  editorResizeObserver.observe(editorContainer.value)

  requestAnimationFrame(() => {
    editor?.requestMeasure()
  })
}

// Helper to parse created/selected into tables from SQL query
function detectAndExtractCreatedTable(sql: string): { schema: string; name: string } | null {
  const cleanSql = sql.replace(/\/\*[\s\S]*?\*\/|--.*$/gm, '') // Remove SQL comments
  
  // 1. SELECT ... INTO [TEMP | TEMPORARY | UNLOGGED]? [TABLE]? <table_name> ...
  const selectIntoRegex = /\bINTO\s+(?:TEMPORARY\s+|TEMP\s+|UNLOGGED\s+)?(?:TABLE\s+)?([a-zA-Z0-9_"\.]+)/i
  const selectIntoMatch = cleanSql.match(selectIntoRegex)
  if (selectIntoMatch) {
    return parseTableName(selectIntoMatch[1])
  }
  
  // 2. CREATE [GLOBAL/LOCAL TEMP/UNLOGGED] TABLE [IF NOT EXISTS] <table_name> ...
  const createTableRegex = /\bCREATE\s+(?:GLOBAL\s+|LOCAL\s+)?(?:TEMPORARY\s+|TEMP\s+|UNLOGGED\s+)?TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?([a-zA-Z0-9_"\.]+)/i
  const createTableMatch = cleanSql.match(createTableRegex)
  if (createTableMatch) {
    return parseTableName(createTableMatch[1])
  }
  
  return null
}

function parseTableName(rawName: string): { schema: string; name: string } {
  const cleanName = rawName.replace(/"/g, '') // Remove quotes
  if (cleanName.includes('.')) {
    const parts = cleanName.split('.')
    return {
      schema: parts[0],
      name: parts[1]
    }
  }
  return {
    schema: 'public',
    name: cleanName
  }
}

async function runQuery(customSql?: string | MouseEvent) {
  const connId = props.tab.connectionId || connectionsStore.currentConnectionId
  if (!connId || isExecuting.value) return

  const originalSql = editor?.state.doc.toString() || ''
  const sql = (typeof customSql === 'string' && customSql) ? customSql : originalSql
  if (!sql.trim()) return

  // Save original SQL to query history log
  saveToQueryLog(originalSql)

  lastExecutedSql.value = sql

  isExecuting.value = true
  results.value = []
  activeResultIndex.value = 0

  try {
    const bindings = connectionsStore.getWailsBindings()
    
    // Check if there are multiple statements via a simple semicolon check
    // Ensure we don't count trailing semicolons
    const trimmedSql = sql.trim()
    const hasMultipleStatements = trimmedSql.includes(';') && trimmedSql.indexOf(';') < trimmedSql.length - 1

    if (hasMultipleStatements && (bindings as any).ExecuteMultipleQueries) {
      const resArray = await (bindings as any).ExecuteMultipleQueries(connId, sql, 30)
      if (resArray && resArray.length > 0) {
        results.value = resArray
      }
    } else {
      const result = await bindings.ExecuteQuery(connId, sql, 30)
      results.value = [result]
    }

    // Detect if a table was created and add it to the workspace tree
    const firstResult = results.value.length > 0 ? results.value[0] : null
    if (firstResult && !firstResult.error) {
      const createdTable = detectAndExtractCreatedTable(sql)
      if (createdTable) {
        const conn = connectionsStore.connections.find(c => c.id === connId)
        const connName = conn?.name || 'Database'
        const dbName = conn?.database || 'postgres'
        
        await workspaceStore.addObject({
          connection_id: connId,
          connection_name: connName,
          database_name: dbName,
          schema: createdTable.schema,
          name: createdTable.name,
          type: 'table',
        })
        
        // Reload schema store tables to sync autocomplete/metadata
        try {
          await schemaStore.loadTables(connId, createdTable.schema)
        } catch (err) {
          console.warn('[runQuery] Failed to reload tables list:', err)
        }
      }
    }
  } catch (e: any) {
    results.value = [{
      columns: [],
      rows: [],
      row_count: 0,
      total_count: 0,
      duration_ms: 0,
      error: e.message || String(e),
      query_type: 'error',
    }]
  } finally {
    isExecuting.value = false
  }
}

function formatSql() {
  if (!editor) return
  const doc = editor.state.doc.toString()
  const formatted = formatSQL(doc)
  editor.dispatch({
    changes: { from: 0, to: editor.state.doc.length, insert: formatted },
  })
}

function clearEditor() {
  if (!editor) return
  editor.dispatch({
    changes: { from: 0, to: editor.state.doc.length, insert: '' },
  })
}

async function handleEditorDrop(event: DragEvent) {
  event.preventDefault()
  const dataStr = event.dataTransfer?.getData('application/json')
  if (!dataStr) return

  try {
    const tableInfo = JSON.parse(dataStr)
    if (tableInfo && (tableInfo.type === 'table' || tableInfo.type === 'view')) {
      if (tableInfo.connectionId) {
        tabsStore.updateTab(props.tab.id, { connectionId: tableInfo.connectionId })
      }
      const connId = tableInfo.connectionId || props.tab.connectionId || connectionsStore.currentConnectionId
      let fieldsStr = '*'

      if (connId) {
        const schemaStore = useSchemaStore()
        const cols = await schemaStore.loadColumns(connId, tableInfo.schema, tableInfo.name)
        if (cols && cols.length > 0) {
          fieldsStr = cols.map(c => c.column_name).join(', ')
        }
      }

      const sql = `SELECT ${fieldsStr} FROM ${tableInfo.schema}.${tableInfo.name};\n`
      if (editor) {
        const range = editor.state.selection.main
        editor.dispatch({
          changes: {
            from: range.from,
            to: range.to,
            insert: sql
          },
          selection: { anchor: range.from + sql.length }
        })
        editor.focus()
      }
    }
  } catch (err) {
    console.error('Failed to parse dropped table data:', err)
  }
}

function startResize(e: MouseEvent) {
  resizeStartY = e.clientY
  resizeStartHeight = resultHeight.value

  const onMove = (e: MouseEvent) => {
    const delta = resizeStartY - e.clientY
    resultHeight.value = Math.max(100, Math.min(600, resizeStartHeight + delta))
    editorFlexBasis.value = `calc(100% - ${resultHeight.value + 1}px)`
    editor?.requestMeasure()
  }

  const onUp = () => {
    document.removeEventListener('mousemove', onMove)
    document.removeEventListener('mouseup', onUp)
  }

  document.addEventListener('mousemove', onMove)
  document.addEventListener('mouseup', onUp)
}
</script>

<style scoped>
.ctx-menu {
  position: fixed;
  z-index: 9999;
  min-width: 200px;
  background: #1a1a1a;
  border: 1px solid #2a2a2a;
  border-radius: 8px;
  padding: 4px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5), 0 0 0 1px rgba(255,255,255,0.04);
  backdrop-filter: blur(12px);
  animation: ctx-fade-in 0.12s ease-out;
}

@keyframes ctx-fade-in {
  from { opacity: 0; transform: scale(0.96) translateY(-4px); }
  to   { opacity: 1; transform: scale(1) translateY(0); }
}

.ctx-menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 6px 10px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: #c8d0da;
  font-size: 12px;
  cursor: pointer;
  transition: background 0.1s, color 0.1s;
  text-align: left;
}

.ctx-menu-item:hover {
  background: #00C9A7;
  color: #0a0a0a;
}

.ctx-menu-item:hover .ctx-shortcut {
  color: rgba(10, 10, 10, 0.6);
}

.ctx-menu-item.ctx-danger:hover {
  background: #EF4444;
  color: #fff;
}

.ctx-icon {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
}

.ctx-shortcut {
  margin-left: auto;
  font-size: 10px;
  color: #4b5563;
}

.ctx-separator {
  height: 1px;
  background: #2a2a2a;
  margin: 4px 8px;
}
</style>
