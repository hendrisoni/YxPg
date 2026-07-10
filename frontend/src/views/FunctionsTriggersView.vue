<template>
  <div class="h-full w-full flex flex-col overflow-hidden bg-navy-primary">
    <!-- View wrapper -->
    <div class="flex-1 flex gap-4 min-h-0 overflow-hidden p-4">

      <!-- Left Panel: Connection select & Lists -->
      <div
        class="w-[360px] flex-shrink-0 flex flex-col min-h-0 bg-navy-secondary border border-navy-border rounded-lg p-4 space-y-4">

        <!-- Connection selector & Actions -->
        <div class="flex gap-2 items-center select-none flex-shrink-0">
          <select v-model="selectedConnectionId" @change="loadData"
            class="flex-1 text-xs bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-text-primary focus:border-teal-accent focus:outline-none cursor-pointer">
            <option v-for="conn in connectionsStore.connections" :key="conn.id" :value="conn.id">
              {{ conn.name }} ({{ conn.database }})
            </option>
          </select>
          <button @click="loadData" :disabled="loading || !selectedConnectionId"
            class="p-1.5 border border-navy-border rounded-md text-text-secondary hover:bg-navy-hover hover:text-text-primary transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
            title="Refresh list">
            <RefreshCw class="w-4 h-4" :class="loading ? 'animate-spin' : ''" />
          </button>
        </div>

        <!-- Search Bar & Toggle -->
        <div class="flex gap-2 items-center select-none flex-shrink-0">
          <div class="relative flex-1">
            <Search class="absolute left-2.5 top-1/2 -translate-y-1/2 w-4 h-4 text-text-muted" />
            <input v-model="searchQuery" type="text" placeholder="Search by name or code content..."
              class="w-full pl-9 pr-8 py-1.5 text-xs bg-navy-tertiary border border-navy-border rounded-md text-text-primary placeholder-text-muted focus:border-teal-accent focus:outline-none font-sans" />
            <button v-if="searchQuery" @click="searchQuery = ''"
              class="absolute right-2.5 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-primary text-[10px] cursor-pointer">
              Clear
            </button>
          </div>
          <button @click="showTriggersCategory = !showTriggersCategory"
            class="p-1.5 border border-navy-border rounded-md text-text-secondary transition-all cursor-pointer flex-shrink-0"
            :class="showTriggersCategory ? 'bg-purple-500/20 border-purple-500/40 text-purple-400 shadow-[0_0_12px_rgba(168,85,247,0.3)]' : 'hover:bg-navy-hover hover:text-text-primary'"
            title="Toggle Trigger Category Folders">
            <svg class="w-5 h-5 transition-all duration-300"
              :class="showTriggersCategory ? 'fill-purple-400/30 drop-shadow-[0_0_5px_rgba(168,85,247,0.8)]' : 'fill-none'"
              viewBox="0 0 24 24" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2" />
            </svg>
          </button>
        </div>

        <!-- Scrollable Lists -->
        <div class="flex-1 overflow-y-auto pr-1 space-y-3 min-h-0">
          <div v-if="loading" class="h-40 flex items-center justify-center">
            <div class="w-6 h-6 border-2 border-teal-accent border-t-transparent rounded-full animate-spin"></div>
          </div>

          <div v-else-if="!selectedConnectionId"
            class="h-40 flex items-center justify-center text-xs text-text-muted text-center select-none p-4">
            Select a connection to load PL/SQL Functions & Triggers.
          </div>

          <div v-else class="space-y-4">
            <!-- 1. CATEGORIZED VIEW (showTriggersCategory = true) -->
            <div v-if="showTriggersCategory" class="space-y-4">
              <!-- Triggers Category -->
              <div class="space-y-1">
                <div @click="showTriggers = !showTriggers"
                  class="flex items-center gap-1.5 py-1 px-2 hover:bg-navy-hover rounded cursor-pointer select-none text-[11px] font-semibold text-text-primary uppercase tracking-wider">
                  <component :is="showTriggers ? ChevronDown : ChevronRight" class="w-3.5 h-3.5 text-text-muted" />
                  <Zap class="w-3.5 h-3.5 text-purple-400" />
                  <span>Triggers</span>
                  <span class="text-[10px] text-text-muted font-normal ml-1">({{ filteredTriggers.length }})</span>
                </div>

                <div v-show="showTriggers" class="pl-3 border-l border-navy-border/60 ml-3.5 mt-0.5 space-y-0.5">
                  <div v-if="filteredTriggers.length === 0"
                    class="py-2 pl-3 text-[10px] text-text-muted italic select-none">
                    No triggers match search query.
                  </div>
                  <div v-for="tg in filteredTriggers" :key="tg.schema_name + '.' + tg.table_name + '.' + tg.trigger_name"
                    @click="selectItem(tg, 'trigger')"
                    class="group flex flex-col p-2.5 rounded-lg border hover:bg-navy-hover cursor-pointer transition-colors"
                    :class="selectedItem && selectedItem.trigger_name === tg.trigger_name ? 'border-accent-blue bg-accent-blue/5' : 'border-transparent bg-navy-tertiary/40'">
                    <div class="flex items-center justify-between min-w-0">
                      <span class="font-mono text-xs font-semibold text-text-primary truncate" :title="tg.trigger_name">
                        {{ tg.trigger_name }}
                      </span>
                      <span
                        class="px-1 py-0.5 rounded bg-navy-hover border border-navy-border text-[8px] text-text-muted font-mono uppercase">
                        {{ tg.schema_name }}
                      </span>
                    </div>
                    <div class="text-[9px] text-accent-amber font-mono truncate mt-1">
                      ON {{ tg.schema_name }}.{{ tg.table_name }}
                    </div>
                  </div>
                </div>
              </div>

              <!-- Functions Category -->
              <div class="space-y-1">
                <div @click="showFunctions = !showFunctions"
                  class="flex items-center gap-1.5 py-1 px-2 hover:bg-navy-hover rounded cursor-pointer select-none text-[11px] font-semibold text-text-primary uppercase tracking-wider">
                  <component :is="showFunctions ? ChevronDown : ChevronRight" class="w-3.5 h-3.5 text-text-muted" />
                  <Cpu class="w-3.5 h-3.5 text-accent-blue" />
                  <span>Functions</span>
                  <span class="text-[10px] text-text-muted font-normal ml-1">({{ filteredFunctions.length }})</span>
                </div>

                <div v-show="showFunctions" class="pl-3 border-l border-navy-border/60 ml-3.5 mt-0.5 space-y-0.5">
                  <div v-if="filteredFunctions.length === 0"
                    class="py-2 pl-3 text-[10px] text-text-muted italic select-none">
                    No functions match search query.
                  </div>
                  <div v-for="fn in filteredFunctions" :key="fn.schema_name + '.' + fn.function_name"
                    @click="selectItem(fn, 'function')"
                    class="group flex flex-col p-2.5 rounded-lg border hover:bg-navy-hover cursor-pointer transition-colors"
                    :class="selectedItem && selectedItem.function_name === fn.function_name ? 'border-accent-blue bg-accent-blue/5' : 'border-transparent bg-navy-tertiary/40'">
                    <div class="flex items-center justify-between min-w-0">
                      <span class="font-mono text-xs font-semibold text-text-primary truncate" :title="fn.function_name">
                        {{ fn.function_name }}()
                      </span>
                      <span
                        class="px-1 py-0.5 rounded bg-navy-hover border border-navy-border text-[8px] text-text-muted font-mono uppercase">
                        {{ fn.schema_name }}
                      </span>
                    </div>
                    <div class="flex items-center gap-2 mt-1">
                      <span class="text-[9px] text-teal-accent font-sans uppercase font-bold tracking-wider">
                        {{ fn.language_name }}
                      </span>
                      <span v-if="fn.function_type === 'trigger'"
                        class="px-1 py-0.2 rounded bg-purple-500/10 border border-purple-500/20 text-[8px] text-purple-400 font-bold uppercase tracking-wider">
                        Trigger Fn
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 2. FLAT UNIFIED VIEW (showTriggersCategory = false, default) -->
            <div v-else class="space-y-1">
              <div v-if="filteredFunctions.length === 0" class="py-2 pl-3 text-[10px] text-text-muted italic select-none">
                No items match search query.
              </div>
              <div v-for="fn in filteredFunctions" :key="fn.schema_name + '.' + fn.function_name"
                @click="selectItem(fn, 'function')"
                class="group flex flex-col p-2.5 rounded-lg border hover:bg-navy-hover cursor-pointer transition-colors"
                :class="selectedItem && selectedItem.function_name === fn.function_name ? 'border-accent-blue bg-accent-blue/5' : 'border-transparent bg-navy-tertiary/40'">
                <div class="flex items-center justify-between min-w-0">
                  <span class="font-mono text-xs font-semibold text-text-primary truncate" :title="fn.function_name">
                    {{ fn.function_name }}()
                  </span>
                  <span
                    class="px-1 py-0.5 rounded bg-navy-hover border border-navy-border text-[8px] text-text-muted font-mono uppercase">
                    {{ fn.schema_name }}
                  </span>
                </div>

                <!-- Trigger Function Metadata with Table Info -->
                <div v-if="fn.function_type === 'trigger'" class="mt-1 space-y-1">
                  <div v-if="getTriggerTableInfo(fn)" class="text-[9.5px] text-accent-amber font-mono truncate">
                    ON {{ getTriggerTableInfo(fn) }}
                  </div>
                  <div v-for="tg in getAttachedTriggers(fn)" :key="tg.trigger_name" class="flex items-center gap-1 mt-0.5">
                    <Zap class="w-3.5 h-3.5 text-purple-400 flex-shrink-0" />
                    <span class="text-[9px] text-text-muted font-mono truncate">
                      Trigger: {{ tg.trigger_name }}
                    </span>
                  </div>
                  <div v-if="getAttachedTriggers(fn).length === 0" class="text-[9px] text-text-muted italic">
                    Trigger function (no active triggers)
                  </div>
                </div>

                <!-- Standard Function Metadata -->
                <div v-else class="flex items-center gap-2 mt-1">
                  <span class="text-[9px] text-teal-accent font-sans uppercase font-bold tracking-wider">
                    {{ fn.language_name }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Panel: CodeMirror SQL Preview -->
      <div
        class="flex-1 flex flex-col min-h-0 bg-navy-secondary border border-navy-border rounded-lg overflow-hidden relative">
        <div v-if="selectedItem" class="flex-1 flex flex-col min-h-0 overflow-hidden">
          <div
            class="px-4 py-2.5 bg-navy-tertiary border-b border-navy-border select-none flex-shrink-0 flex items-center justify-between">
            <div class="flex items-center gap-2 min-w-0 pr-2">
              <Zap v-if="selectedType === 'trigger'" class="w-4 h-4 text-purple-400 flex-shrink-0" />
              <Cpu v-else class="w-4 h-4 text-accent-blue flex-shrink-0" />
              <h3 class="text-xs font-semibold text-text-primary font-mono truncate">
                {{ selectedType === 'trigger' ? 'TRIGGER' : 'FUNCTION' }}
                {{ selectedItem.schema_name }}.{{ selectedType === 'trigger' ? selectedItem.trigger_name :
                  selectedItem.function_name }}
              </h3>
            </div>

            <button @click="copyDdl"
              class="px-2.5 py-1 text-[11px] border border-navy-border rounded text-text-secondary hover:bg-navy-hover hover:text-text-primary transition-colors flex items-center gap-1.5 cursor-pointer">
              <Copy class="w-3.5 h-3.5" />
              <span>Copy SQL</span>
            </button>
          </div>

          <div class="flex-1 w-full overflow-hidden min-h-0 bg-[#121212] relative">
            <div ref="editorContainer" class="h-full w-full overflow-hidden"></div>
          </div>
        </div>

        <div v-else
          class="flex-1 flex flex-col items-center justify-center p-6 text-center select-none text-text-muted space-y-2">
          <svg class="w-12 h-12 text-navy-border" viewBox="0 0 24 24" fill="none" stroke="currentColor"
            stroke-width="1.5">
            <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2" />
          </svg>
          <div class="text-xs font-medium text-text-secondary">No Function or Trigger Selected</div>
          <div class="text-[11px]">Select any item from the explorer on the left to view its definition.</div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { useConnectionsStore } from '../stores/connections'
import { useUiStore } from '../stores/ui'
import type { Tab } from '../types'
import {
  ChevronRight,
  ChevronDown,
  Zap,
  Cpu,
  Search,
  RefreshCw,
  Copy
} from 'lucide-vue-next'

const props = defineProps<{
  tab: Tab
}>()

const connectionsStore = useConnectionsStore()
const uiStore = useUiStore()

const loading = ref(false)
const selectedConnectionId = ref<string>('')
const searchQuery = ref('')

const showTriggers = ref(true)
const showFunctions = ref(true)
const showTriggersCategory = ref(false)

const functionsList = ref<any[]>([])
const triggersList = ref<any[]>([])

const selectedItem = ref<any | null>(null)
const selectedType = ref<'trigger' | 'function' | null>(null)
const previewCode = ref('')

const editorContainer = ref<HTMLElement | null>(null)
let editor: any = null

function parseQueryResult(result: any): any[] {
  if (!result || !result.rows) return []
  const cols = result.columns
  return result.rows.map((row: any[]) => {
    const obj: any = {}
    cols.forEach((col: any, idx: number) => {
      const colName = col && typeof col === 'object' ? col.name : col
      obj[colName] = row[idx]
    })
    return obj
  })
}

const filteredTriggers = computed(() => {
  const query = searchQuery.value.toLowerCase().trim()
  if (!query) return triggersList.value

  return triggersList.value.filter(tg => {
    const nameMatch = tg.trigger_name.toLowerCase().includes(query)
    const tableMatch = tg.table_name.toLowerCase().includes(query)
    const contentMatch = tg.definition && tg.definition.toLowerCase().includes(query)
    return nameMatch || tableMatch || contentMatch
  })
})

const filteredFunctions = computed(() => {
  const query = searchQuery.value.toLowerCase().trim()
  if (!query) return functionsList.value

  return functionsList.value.filter(fn => {
    const nameMatch = fn.function_name.toLowerCase().includes(query)
    const contentMatch = fn.definition && fn.definition.toLowerCase().includes(query)
    return nameMatch || contentMatch
  })
})

function formatTriggerDef(def: string): string {
  if (!def) return ''
  let formatted = def.trim().replace(/\s+/g, ' ')
  formatted = formatted.replace(/\s+(BEFORE|AFTER|INSTEAD OF)\s+/i, '\n$1 ')
  formatted = formatted.replace(/\s+(ON)\s+/i, '\n$1 ')
  formatted = formatted.replace(/\s+(FOR\s+(?:EACH\s+)?(?:ROW|STATEMENT))\s+/i, '\n$1 ')
  formatted = formatted.replace(/\bEXECUTE\s+FUNCTION\b/gi, 'EXECUTE PROCEDURE')
  return formatted
}

async function selectItem(item: any, type: 'trigger' | 'function') {
  selectedItem.value = item
  selectedType.value = type

  let code = ''
  if (type === 'trigger') {
    const schema = item.schema_name
    const fnName = item.function_name
    const fnDef = item.function_definition || ''
    const tgDef = item.definition || ''

    code = `DROP FUNCTION IF EXISTS ${schema}.${fnName}() CASCADE;\n\n${fnDef}\n\n${formatTriggerDef(tgDef)};`
  } else if (type === 'function' && item.function_type === 'trigger') {
    const schema = item.schema_name
    const fnName = item.function_name
    const fnDef = item.definition || ''

    // Find matching triggers that use this function
    const matchingTriggers = triggersList.value.filter(tg =>
      tg.function_name === fnName &&
      tg.schema_name === schema
    )

    let triggerBlocks = ''
    if (matchingTriggers.length > 0) {
      triggerBlocks = '\n\n' + matchingTriggers.map(tg => `${formatTriggerDef(tg.definition)};`).join('\n\n')
    }

    code = `DROP FUNCTION IF EXISTS ${schema}.${fnName}() CASCADE;\n\n${fnDef}${triggerBlocks}`
  } else {
    code = item.definition || ''
  }

  previewCode.value = code

  await nextTick()
  initEditor(code)
}

async function destroyEditor() {
  if (editor) {
    editor.destroy()
    editor = null
  }
}

async function initEditor(code: string) {
  destroyEditor()
  if (!editorContainer.value) return

  try {
    const { EditorView, lineNumbers, highlightSpecialChars, drawSelection } = await import('@codemirror/view')
    const { syntaxHighlighting, HighlightStyle } = await import('@codemirror/language')
    const { tags } = await import('@lezer/highlight')
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

    const darkTheme = EditorView.theme({
      '&': {
        backgroundColor: '#121212',
        color: '#A9B7C6',
        height: '100%',
      },
      '.cm-scroller': {
        fontFamily: "'DM Mono', ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace",
        fontSize: '13px',
        lineHeight: '20px',
      },
      '.cm-gutters': {
        backgroundColor: '#0C0C0C',
        color: '#858585',
        borderRight: '1px solid #202020',
      }
    })

    const state = EditorState.create({
      doc: code,
      extensions: [
        lineNumbers(),
        highlightSpecialChars(),
        drawSelection(),
        syntaxHighlighting(darculaHighlightStyle),
        sqlLang(),
        darkTheme,
        EditorState.readOnly.of(true),
        EditorView.editable.of(false)
      ]
    })

    editor = new EditorView({
      state,
      parent: editorContainer.value
    })
  } catch (err) {
    console.error('Failed to initialize CodeMirror editor:', err)
  }
}

function copyDdl() {
  if (!previewCode.value) return
  navigator.clipboard.writeText(previewCode.value)
  uiStore.addNotification({
    type: 'info',
    title: 'Copied',
    message: 'SQL definition copied to clipboard'
  })
}

function getAttachedTriggers(fn: any) {
  if (!fn || fn.function_type !== 'trigger') return []
  return triggersList.value.filter(tg =>
    tg.function_name === fn.function_name &&
    tg.schema_name === fn.schema_name
  )
}

function getTriggerTableInfo(fn: any): string {
  const attached = getAttachedTriggers(fn)
  if (attached.length === 0) return ''
  const tables = [...new Set(attached.map(tg => `${tg.schema_name}.${tg.table_name}`))]
  return tables.join(', ')
}

async function loadData() {
  const connId = selectedConnectionId.value
  if (!connId) return

  loading.value = true
  selectedItem.value = null
  destroyEditor()

  const functionsSql = `
    SELECT 
        n.nspname AS schema_name,
        p.proname AS function_name,
        l.lanname AS language_name,
        pg_get_functiondef(p.oid) AS definition,
        CASE 
            WHEN t.typname = 'trigger' THEN 'trigger'
            ELSE 'normal'
        END AS function_type
    FROM pg_proc p
    JOIN pg_namespace n ON p.pronamespace = n.oid
    JOIN pg_language l ON p.prolang = l.oid
    LEFT JOIN pg_type t ON p.prorettype = t.oid
    WHERE n.nspname NOT IN ('pg_catalog', 'information_schema')
      AND l.lanname IN ('plpgsql', 'sql')
    ORDER BY schema_name, function_name;
  `

  const triggersSql = `
    SELECT 
        n.nspname AS schema_name,
        c.relname AS table_name,
        t.tgname AS trigger_name,
        p.proname AS function_name,
        pg_get_functiondef(p.oid) AS function_definition,
        pg_get_triggerdef(t.oid) AS definition
    FROM pg_trigger t
    JOIN pg_class c ON t.tgrelid = c.oid
    JOIN pg_namespace n ON c.relnamespace = n.oid
    JOIN pg_proc p ON t.tgfoid = p.oid
    WHERE NOT t.tgisinternal
      AND n.nspname NOT IN ('pg_catalog', 'information_schema')
    ORDER BY schema_name, table_name, trigger_name;
  `

  try {
    const bindings = connectionsStore.getWailsBindings()

    // Execute sequentially to prevent context cancels
    const funcsRes = await bindings.ExecuteQuery(connId, functionsSql, 30)
    const trigRes = await bindings.ExecuteQuery(connId, triggersSql, 30)

    if (funcsRes.error) throw new Error(funcsRes.error)
    if (trigRes.error) throw new Error(trigRes.error)

    functionsList.value = parseQueryResult(funcsRes)
    triggersList.value = parseQueryResult(trigRes)
  } catch (err: any) {
    uiStore.addNotification({
      type: 'error',
      title: 'Load Functions & Triggers Failed',
      message: err.message || String(err)
    })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  // Use tab connectionId or fallback to current connection
  selectedConnectionId.value = props.tab.connectionId || connectionsStore.currentConnectionId || ''
  if (selectedConnectionId.value) {
    loadData()
  }
})

watch(() => props.tab.connectionId, (newId) => {
  if (newId) {
    selectedConnectionId.value = newId
    loadData()
  }
})
</script>
