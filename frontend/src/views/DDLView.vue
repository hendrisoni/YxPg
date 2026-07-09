<template>
  <div class="h-full flex flex-col overflow-hidden bg-navy-primary">
    <!-- Header tabs -->
    <div class="flex items-center gap-2 border-b border-navy-border bg-navy-secondary px-4 py-2 flex-shrink-0 select-none">
      <button
        @click="activeTab = 'columns'"
        class="px-3.5 py-1.5 text-xs rounded-md transition-colors flex items-center gap-1.5"
        :class="activeTab === 'columns' ? 'bg-teal-accent/15 text-teal-accent font-semibold border border-teal-accent/30' : 'text-text-secondary hover:text-text-primary hover:bg-navy-hover border border-transparent'"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="3" width="18" height="18" rx="2" /><path d="M3 9h18M3 15h18M9 3v18" />
        </svg>
        <span>Columns ({{ columnsList.length }})</span>
      </button>
      <button
        @click="activeTab = 'ddl'"
        class="px-3.5 py-1.5 text-xs rounded-md transition-colors flex items-center gap-1.5"
        :class="activeTab === 'ddl' ? 'bg-teal-accent/15 text-teal-accent font-semibold border border-teal-accent/30' : 'text-text-secondary hover:text-text-primary hover:bg-navy-hover border border-transparent'"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
          <polyline points="14 2 14 8 20 8" />
        </svg>
        <span>DDL</span>
      </button>
      
      <div class="flex-1"></div>
      
      <!-- Table descriptor -->
      <div class="text-xs text-text-muted font-mono bg-navy-tertiary/60 border border-navy-border/50 px-2.5 py-1 rounded">
        <span>{{ tab.schema || 'public' }}</span>
        <span class="text-teal-accent">.</span>
        <span class="text-text-primary font-semibold">{{ tab.table }}</span>
      </div>
    </div>

    <!-- Active Tab Panel -->
    <div class="flex-1 min-h-0 overflow-hidden relative">
      <!-- Columns View -->
      <div v-show="activeTab === 'columns'" class="h-full w-full flex flex-col overflow-hidden" @contextmenu="handleContainerContextMenu">
        <!-- Table container -->
        <div ref="tableContainer" class="flex-1 w-full overflow-hidden"></div>
        
        <!-- Status/Keyboard hints footer -->
        <div class="flex items-center justify-between px-4 py-2 bg-navy-secondary border-t border-navy-border text-[10px] text-text-muted select-none">
          <div class="flex items-center gap-4">
            <span>Right-click anywhere to open actions menu.</span>
            <span>•</span>
            <span class="flex items-center gap-1">
              <kbd class="px-1 bg-navy-tertiary border border-navy-border rounded text-[9px]">Ins</kbd> New Column
            </span>
            <span class="flex items-center gap-1">
              <kbd class="px-1 bg-navy-tertiary border border-navy-border rounded text-[9px]">Enter</kbd> Edit
            </span>
            <span class="flex items-center gap-1">
              <kbd class="px-1 bg-navy-tertiary border border-navy-border rounded text-[9px]">Del</kbd> Drop
            </span>
          </div>
          <div>
            <span>PostgreSQL Table Schema Properties</span>
          </div>
        </div>
      </div>

      <!-- DDL View -->
      <div v-show="activeTab === 'ddl'" class="h-full w-full flex flex-col overflow-hidden p-4">
        <div v-if="loadingDdl" class="flex-1 flex items-center justify-center">
          <div class="inline-block w-6 h-6 border-2 border-teal-accent border-t-transparent rounded-full animate-spin"></div>
        </div>
        <div v-show="!loadingDdl" class="flex-1 flex flex-col min-h-0 overflow-hidden rounded-lg">
          <!-- CodeMirror editor container -->
          <div ref="ddlEditorContainer" class="flex-1 w-full overflow-hidden min-h-0 bg-navy-tertiary border border-navy-border rounded-t-lg"></div>
          
          <div class="px-4 py-3 bg-navy-secondary border-t-0 border border-navy-border rounded-b-lg flex justify-end gap-2 flex-shrink-0">
            <button
              @click="copyDDL"
              class="px-4 py-1.5 text-xs border border-navy-border rounded-md text-text-secondary hover:bg-navy-hover hover:text-text-primary transition-colors flex items-center gap-1.5 cursor-pointer"
            >
              <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
              </svg>
              Copy DDL
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Floating Context Menu -->
    <Teleport to="body">
      <div
        v-if="showContextMenu"
        class="fixed z-50 min-w-[180px] bg-navy-secondary border border-navy-border rounded-lg shadow-xl py-1 select-none text-xs"
        :style="{ top: contextMenuY + 'px', left: contextMenuX + 'px' }"
      >
        <button
          @click="openNewColumnModal"
          class="w-full px-3 py-1.5 text-left text-text-primary hover:bg-navy-hover flex items-center justify-between cursor-pointer"
        >
          <span class="flex items-center gap-2">
            <svg class="w-3.5 h-3.5 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 5v14M5 12h14"/></svg>
            New Column...
          </span>
          <span class="text-text-muted text-[10px]">Ins</span>
        </button>
        <button
          @click="openEditColumnModal(selectedRowData)"
          :disabled="!selectedRowData"
          class="w-full px-3 py-1.5 text-left text-text-primary hover:bg-navy-hover disabled:opacity-40 disabled:hover:bg-transparent flex items-center justify-between cursor-pointer"
        >
          <span class="flex items-center gap-2">
            <svg class="w-3.5 h-3.5 text-accent-blue" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/></svg>
            Edit Column '{{ selectedRowData?.column_name || 'column' }}'...
          </span>
          <span class="text-text-muted text-[10px]">Enter</span>
        </button>
        <button
          @click="openRenameColumnModal(selectedRowData)"
          :disabled="!selectedRowData"
          class="w-full px-3 py-1.5 text-left text-text-primary hover:bg-navy-hover disabled:opacity-40 disabled:hover:bg-transparent flex items-center justify-between cursor-pointer"
        >
          <span class="flex items-center gap-2">
            <svg class="w-3.5 h-3.5 text-accent-amber" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20h9M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z"/></svg>
            Rename Column '{{ selectedRowData?.column_name || 'column' }}'...
          </span>
        </button>
        <div class="h-px bg-navy-border my-1"></div>
        <button
          @click="confirmDropColumn(selectedRowData)"
          :disabled="!selectedRowData"
          class="w-full px-3 py-1.5 text-left text-accent-red hover:bg-red-500/10 disabled:opacity-40 disabled:hover:bg-transparent flex items-center justify-between cursor-pointer"
        >
          <span class="flex items-center gap-2">
            <svg class="w-3.5 h-3.5 text-accent-red" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 6h18M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2M10 11v6M14 11v6"/></svg>
            Drop Column '{{ selectedRowData?.column_name || 'column' }}'
          </span>
          <span class="text-text-muted text-[10px]">Del</span>
        </button>
      </div>
    </Teleport>

    <!-- Column Add/Edit Modal -->
    <Modal
      :show="isColumnModalOpen"
      :title="columnModalMode === 'add' ? 'New Column' : `Edit Column '${originalColData?.column_name}'`"
      @close="isColumnModalOpen = false"
      size="md"
    >
      <form @submit.prevent="handleColumnSubmit" class="space-y-4 text-xs py-1">
        <!-- Error alert -->
        <div v-if="submitError" class="p-3 bg-red-500/10 border border-red-500/20 text-red-400 rounded-md">
          {{ submitError }}
        </div>

        <!-- Column Name -->
        <div class="flex flex-col gap-1">
          <label class="text-[11px] font-semibold text-text-secondary">Column Name</label>
          <input
            type="text"
            v-model="colName"
            placeholder="e.g. status_code"
            class="w-full text-xs bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-text-primary focus:border-teal-accent focus:outline-none focus:ring-0"
            required
            autoFocus
          />
        </div>

        <!-- Data Type Selection -->
        <div class="flex gap-3">
          <div class="flex-1 flex flex-col gap-1">
            <label class="text-[11px] font-semibold text-text-secondary">Data Type</label>
            <select
              v-model="dataType"
              class="w-full text-xs bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-text-primary focus:border-teal-accent focus:outline-none cursor-pointer"
            >
              <option v-for="t in commonDataTypes" :key="t.value" :value="t.value">{{ t.label }}</option>
              <option value="custom">Other...</option>
            </select>
          </div>
          <div v-if="dataType === 'custom'" class="flex-1 flex flex-col gap-1">
            <label class="text-[11px] font-semibold text-text-secondary">Custom Type</label>
            <input
              type="text"
              v-model="customDataType"
              placeholder="e.g. json, point"
              class="w-full text-xs bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-text-primary focus:border-teal-accent focus:outline-none focus:ring-0"
              required
            />
          </div>
        </div>

        <!-- Length & Precision details -->
        <div v-if="dataType === 'varchar' || dataType === 'numeric'" class="flex gap-3">
          <div class="flex-1 flex flex-col gap-1">
            <label class="text-[11px] font-semibold text-text-secondary">
              {{ dataType === 'varchar' ? 'Length' : 'Precision' }}
            </label>
            <input
              type="number"
              v-model.number="length"
              placeholder="e.g. 50"
              class="w-full text-xs bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-text-primary focus:border-teal-accent focus:outline-none focus:ring-0"
            />
          </div>
          <div v-if="dataType === 'numeric'" class="flex-1 flex flex-col gap-1">
            <label class="text-[11px] font-semibold text-text-secondary">Scale</label>
            <input
              type="number"
              v-model.number="scale"
              placeholder="e.g. 2"
              class="w-full text-xs bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-text-primary focus:border-teal-accent focus:outline-none focus:ring-0"
            />
          </div>
        </div>

        <!-- Not Null & Primary Key switches -->
        <div class="flex gap-6 py-2 px-3 bg-navy-tertiary/40 border border-navy-border rounded-md select-none">
          <label class="flex items-center gap-2 cursor-pointer text-text-secondary hover:text-text-primary transition-colors">
            <input
              type="checkbox"
              v-model="isNotNull"
              class="rounded bg-navy-tertiary border-navy-border text-teal-accent focus:ring-0 cursor-pointer"
            />
            <span>Not Null</span>
          </label>
          <label class="flex items-center gap-2 cursor-pointer text-text-secondary hover:text-text-primary transition-colors">
            <input
              type="checkbox"
              v-model="primaryKey"
              class="rounded bg-navy-tertiary border-navy-border text-teal-accent focus:ring-0 cursor-pointer"
            />
            <span>Primary Key</span>
          </label>
        </div>

        <!-- Default Value -->
        <div class="flex flex-col gap-1">
          <label class="text-[11px] font-semibold text-text-secondary">Default Value</label>
          <input
            type="text"
            v-model="defaultValue"
            placeholder="e.g. NULL, 'active', 0"
            class="w-full text-xs bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-text-primary focus:border-teal-accent focus:outline-none focus:ring-0"
          />
        </div>

        <!-- Column Comment -->
        <div class="flex flex-col gap-1">
          <label class="text-[11px] font-semibold text-text-secondary">Description / Comment</label>
          <input
            type="text"
            v-model="description"
            placeholder="Describe what this column is for"
            class="w-full text-xs bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-text-primary focus:border-teal-accent focus:outline-none focus:ring-0"
          />
        </div>
      </form>

      <template #footer>
        <button
          type="button"
          @click="isColumnModalOpen = false"
          class="px-4 py-2 text-xs text-text-secondary hover:text-text-primary transition-colors cursor-pointer"
        >
          Cancel
        </button>
        <button
          type="button"
          @click="handleColumnSubmit"
          :disabled="submitting"
          class="px-4 py-2 text-xs bg-teal-accent text-navy-primary rounded-md font-medium hover:bg-teal-hover transition-colors flex items-center gap-1.5 cursor-pointer disabled:opacity-50"
        >
          <div v-if="submitting" class="w-3.5 h-3.5 border-2 border-navy-primary border-t-transparent rounded-full animate-spin"></div>
          Save Column
        </button>
      </template>
    </Modal>

    <!-- Rename Column Dialog -->
    <Modal
      :show="isRenameModalOpen"
      :title="`Rename Column '${originalColData?.column_name}'`"
      @close="isRenameModalOpen = false"
      size="sm"
    >
      <div class="space-y-3 py-1 text-xs">
        <div class="flex flex-col gap-1.5">
          <label for="rename-col-input" class="text-[11px] font-semibold text-text-secondary">New Name</label>
          <input
            id="rename-col-input"
            type="text"
            v-model="renameColNewName"
            placeholder="e.g. code"
            class="w-full text-xs bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-text-primary focus:border-teal-accent focus:outline-none focus:ring-0"
            required
            @keydown.enter="handleRenameSubmit"
          />
        </div>
      </div>
      <template #footer>
        <button
          type="button"
          @click="isRenameModalOpen = false"
          class="px-4 py-1.5 text-xs text-text-secondary hover:text-text-primary transition-colors cursor-pointer"
        >
          Cancel
        </button>
        <button
          type="button"
          @click="handleRenameSubmit"
          class="px-4 py-1.5 text-xs bg-teal-accent text-navy-primary rounded-md font-medium hover:bg-teal-hover transition-colors cursor-pointer"
        >
          Rename
        </button>
      </template>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick, computed } from 'vue'
import { useConnectionsStore } from '../stores/connections'
import { useSchemaStore } from '../stores/schema'
import { useUiStore } from '../stores/ui'
import type { Tab } from '../types'
import Modal from '../components/shared/Modal.vue'

const props = defineProps<{
  tab: Tab
}>()

const connectionsStore = useConnectionsStore()
const schemaStore = useSchemaStore()
const uiStore = useUiStore()

const activeTab = ref<'columns' | 'ddl'>('columns')
const columnsList = ref<any[]>([])
const tableContainer = ref<HTMLElement | null>(null)
let tabulator: any = null

const ddlScript = ref('')
const loadingDdl = ref(false)
const ddlEditorContainer = ref<HTMLElement | null>(null)
let ddlEditor: any = null

// Context Menu State
const showContextMenu = ref(false)
const contextMenuX = ref(0)
const contextMenuY = ref(0)
const selectedRowData = ref<any>(null)

// Modal Column Form State
const isColumnModalOpen = ref(false)
const columnModalMode = ref<'add' | 'edit'>('add')
const originalColData = ref<any>(null)
const submitting = ref(false)

// Form Fields
const colName = ref('')
const dataType = ref('varchar')
const customDataType = ref('')
const length = ref<number | null>(null)
const scale = ref<number | null>(null)
const isNotNull = ref(false)
const primaryKey = ref(false)
const defaultValue = ref('')
const description = ref('')
const submitError = ref('')

// Rename Modal State
const isRenameModalOpen = ref(false)
const renameColNewName = ref('')

const commonDataTypes = [
  { value: 'varchar', label: 'varchar (Character Varying)' },
  { value: 'text', label: 'text (Unlimited Text)' },
  { value: 'integer', label: 'integer (32-bit Integer)' },
  { value: 'bigint', label: 'bigint (64-bit Integer)' },
  { value: 'boolean', label: 'boolean (True/False)' },
  { value: 'numeric', label: 'numeric (Exact Decimal)' },
  { value: 'timestamp', label: 'timestamp (Date & Time)' },
  { value: 'date', label: 'date (Date Only)' },
  { value: 'uuid', label: 'uuid (UUID Identifier)' },
  { value: 'jsonb', label: 'jsonb (Binary JSON)' }
]

async function loadColumns() {
  const connId = props.tab.connectionId || connectionsStore.currentConnectionId
  if (!connId || !props.tab.table) return

  try {
    const cols = await schemaStore.loadColumns(connId, props.tab.schema || 'public', props.tab.table)
    columnsList.value = cols
    
    await nextTick()
    if (tableContainer.value) {
      if (tabulator) {
        tabulator.setData(cols)
      } else {
        const { TabulatorFull } = await import('tabulator-tables')
        tabulator = new TabulatorFull(tableContainer.value, {
          data: cols,
          columns: [
            {
              title: "",
              field: "key_icon",
              width: 45,
              hozAlign: "center",
              headerSort: false,
              formatter: (cell: any) => {
                const data = cell.getRow().getData()
                let iconsHtml = '<div class="flex items-center gap-1 justify-center h-full">'
                if (data.is_primary_key) {
                  iconsHtml += `<svg class="w-3.5 h-3.5 text-yellow-500 flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/></svg>`
                }
                iconsHtml += `<svg class="w-3.5 h-3.5 text-text-secondary flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2" /><path d="M3 9h18M3 15h18M9 3v18" /></svg>`
                iconsHtml += '</div>'
                return iconsHtml
              }
            },
            { title: "Column Name", field: "column_name", sorter: "string", widthGrow: 2 },
            { title: "Column Type", field: "data_type", sorter: "string", widthGrow: 2 },
            {
              title: "Key",
              field: "is_primary_key",
              sorter: "boolean",
              width: 120,
              formatter: (cell: any) => cell.getValue() ? "Primary Key" : ""
            },
            {
              title: "Not Null",
              field: "is_nullable",
              sorter: "boolean",
              width: 100,
              hozAlign: "center",
              formatter: (cell: any) => {
                const isNotNull = !cell.getValue()
                return `<input type="checkbox" ${isNotNull ? 'checked' : ''} disabled class="rounded bg-navy-tertiary border-navy-border text-teal-accent focus:ring-0 cursor-default w-3.5 h-3.5" />`
              }
            },
            {
              title: "Default",
              field: "column_default",
              sorter: "string",
              widthGrow: 2,
              formatter: (cell: any) => {
                const val = cell.getValue()
                return val === null || val === undefined ? "Null" : val
              }
            },
            { title: "Description", field: "comment", sorter: "string", widthGrow: 3 }
          ],
          layout: 'fitColumns',
          height: '100%',
          selectableRows: 1,
          selectableRowsRangeMode: 'click',
          rowContext: (e: any, row: any) => {
            e.preventDefault()
            e.stopPropagation()
            contextMenuX.value = e.clientX
            contextMenuY.value = e.clientY
            selectedRowData.value = row.getData()
            showContextMenu.value = true
            tabulator.deselectRow()
            row.select()
          },
          rowClick: (e: any, row: any) => {
            selectedRowData.value = row.getData()
          },
          placeholder: 'No columns found'
        })
      }
    }
  } catch (err: any) {
    console.error('Failed to load columns:', err)
  }
}

function destroyDdlEditor() {
  if (ddlEditor) {
    ddlEditor.destroy()
    ddlEditor = null
  }
}

async function loadDDL() {
  const connId = props.tab.connectionId || connectionsStore.currentConnectionId
  if (!connId || !props.tab.table) return

  loadingDdl.value = true
  destroyDdlEditor()
  try {
    const ddl = await schemaStore.getTableDDL(connId, props.tab.schema || 'public', props.tab.table)
    ddlScript.value = ddl

    await nextTick()
    if (ddlEditorContainer.value) {
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
        doc: ddl,
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

      ddlEditor = new EditorView({
        state,
        parent: ddlEditorContainer.value
      })
    }
  } catch (err: any) {
    console.error('Failed to load DDL:', err)
  } finally {
    loadingDdl.value = false
  }
}

function copyDDL() {
  navigator.clipboard.writeText(ddlScript.value)
  uiStore.addNotification({
    type: 'info',
    title: 'Copied',
    message: 'Table DDL copied to clipboard'
  })
}

const closeContextMenu = () => {
  showContextMenu.value = false
}

const handleContainerContextMenu = (e: MouseEvent) => {
  e.preventDefault()
  contextMenuX.value = e.clientX
  contextMenuY.value = e.clientY
  const rowEl = (e.target as HTMLElement).closest('.tabulator-row')
  if (!rowEl) {
    selectedRowData.value = null
    showContextMenu.value = true
  }
}

const handleKeyDown = (e: KeyboardEvent) => {
  if (activeTab.value !== 'columns') return
  
  if (document.activeElement?.tagName === 'INPUT' || document.activeElement?.tagName === 'TEXTAREA' || isColumnModalOpen.value || isRenameModalOpen.value) {
    return
  }

  if (e.key === 'Insert') {
    e.preventDefault()
    openNewColumnModal()
  } else if (e.key === 'Enter' && selectedRowData.value) {
    e.preventDefault()
    openEditColumnModal(selectedRowData.value)
  } else if (e.key === 'Delete' && selectedRowData.value) {
    e.preventDefault()
    confirmDropColumn(selectedRowData.value)
  }
}

// Modal Column Form Helpers
function openNewColumnModal() {
  originalColData.value = null
  colName.value = ''
  dataType.value = 'varchar'
  customDataType.value = ''
  length.value = null
  scale.value = null
  isNotNull.value = false
  primaryKey.value = false
  defaultValue.value = ''
  description.value = ''
  submitError.value = ''
  columnModalMode.value = 'add'
  isColumnModalOpen.value = true
}

function openEditColumnModal(col: any) {
  if (!col) return
  originalColData.value = col
  
  colName.value = col.column_name
  isNotNull.value = !col.is_nullable
  primaryKey.value = col.is_primary_key
  defaultValue.value = col.column_default || ''
  description.value = col.comment || ''
  
  const rawType = col.data_type.toLowerCase()
  
  const lengthMatch = rawType.match(/\((\d+)(?:,\s*(\d+))?\)/)
  if (lengthMatch) {
    length.value = parseInt(lengthMatch[1])
    scale.value = lengthMatch[2] ? parseInt(lengthMatch[2]) : null
  } else {
    length.value = null
    scale.value = null
  }
  
  let cleanType = rawType.replace(/\(\d+(?:,\s*\d+)?\)/, '').trim()
  if (cleanType === 'character varying') {
    cleanType = 'varchar'
  }
  
  const knownTypes = commonDataTypes.map(t => t.value)
  if (knownTypes.includes(cleanType)) {
    dataType.value = cleanType
    customDataType.value = ''
  } else {
    dataType.value = 'custom'
    customDataType.value = cleanType
  }
  
  submitError.value = ''
  columnModalMode.value = 'edit'
  isColumnModalOpen.value = true
}

function openRenameColumnModal(col: any) {
  if (!col) return
  originalColData.value = col
  renameColNewName.value = col.column_name
  isRenameModalOpen.value = true
}

async function handleColumnSubmit() {
  const connId = props.tab.connectionId || connectionsStore.currentConnectionId
  if (!connId || !props.tab.table) return

  submitError.value = ''
  
  const name = colName.value.trim()
  if (!name) {
    submitError.value = 'Column name is required'
    return
  }

  const finalType = dataType.value === 'custom' ? customDataType.value.trim() : dataType.value
  if (!finalType) {
    submitError.value = 'Data type is required'
    return
  }

  const schema = props.tab.schema || 'public'
  const table = props.tab.table

  let sql = ''
  
  if (columnModalMode.value === 'add') {
    const typeWithLength = finalType + (length.value ? `(${length.value}${scale.value ? `,${scale.value}` : ''})` : '')
    sql = `ALTER TABLE "${schema}"."${table}" ADD COLUMN "${name}" ${typeWithLength}`

    if (isNotNull.value) {
      sql += ' NOT NULL'
    }
    if (defaultValue.value.trim()) {
      sql += ` DEFAULT ${defaultValue.value.trim()}`
    }
    sql += ';'

    if (primaryKey.value) {
      sql += `\nALTER TABLE "${schema}"."${table}" ADD PRIMARY KEY ("${name}");`
    }
    if (description.value.trim()) {
      sql += `\nCOMMENT ON COLUMN "${schema}"."${table}"."${name}" IS '${description.value.trim().replace(/'/g, "''")}';`
    }
  } else {
    const orig = originalColData.value
    const origName = orig.column_name
    
    const sqlStatements: string[] = []
    
    const newTypeWithLength = finalType + (length.value ? `(${length.value}${scale.value ? `,${scale.value}` : ''})` : '')
    
    if (name !== origName) {
      sqlStatements.push(`ALTER TABLE "${schema}"."${table}" RENAME COLUMN "${origName}" TO "${name}";`)
    }

    const origTypeNormalized = orig.data_type.toLowerCase()
    if (finalType !== origTypeNormalized || length.value !== orig.character_maximum_length) {
      sqlStatements.push(`ALTER TABLE "${schema}"."${table}" ALTER COLUMN "${name}" TYPE ${newTypeWithLength} USING "${name}"::${finalType};`)
    }

    const origIsNotNull = !orig.is_nullable
    if (isNotNull.value !== origIsNotNull) {
      if (isNotNull.value) {
        sqlStatements.push(`ALTER TABLE "${schema}"."${table}" ALTER COLUMN "${name}" SET NOT NULL;`)
      } else {
        sqlStatements.push(`ALTER TABLE "${schema}"."${table}" ALTER COLUMN "${name}" DROP NOT NULL;`)
      }
    }

    const origDefault = orig.column_default || ''
    const newDefault = defaultValue.value.trim()
    if (newDefault !== origDefault) {
      if (newDefault) {
        sqlStatements.push(`ALTER TABLE "${schema}"."${table}" ALTER COLUMN "${name}" SET DEFAULT ${newDefault};`)
      } else {
        sqlStatements.push(`ALTER TABLE "${schema}"."${table}" ALTER COLUMN "${name}" DROP DEFAULT;`)
      }
    }

    const origComment = orig.comment || ''
    const newComment = description.value.trim()
    if (newComment !== origComment) {
      sqlStatements.push(`COMMENT ON COLUMN "${schema}"."${table}"."${name}" IS '${newComment.replace(/'/g, "''")}';`)
    }

    const origIsPrimaryKey = orig.is_primary_key
    if (primaryKey.value !== origIsPrimaryKey) {
      if (primaryKey.value) {
        sqlStatements.push(`ALTER TABLE "${schema}"."${table}" ADD PRIMARY KEY ("${name}");`)
      } else {
        sqlStatements.push(`ALTER TABLE "${schema}"."${table}" DROP CONSTRAINT IF EXISTS "${table}_pkey";`)
      }
    }

    if (sqlStatements.length === 0) {
      isColumnModalOpen.value = false
      return
    }
    sql = sqlStatements.join('\n')
  }

  submitting.value = true
  try {
    const bindings = connectionsStore.getWailsBindings()
    const result = await bindings.ExecuteQuery(connId, sql, 30)
    if (result.error) {
      throw new Error(result.error)
    }

    uiStore.addNotification({
      type: 'success',
      title: columnModalMode.value === 'add' ? 'Column Added' : 'Column Updated',
      message: `Successfully saved changes for column "${name}"`
    })
    
    isColumnModalOpen.value = false
    loadColumns()
  } catch (err: any) {
    submitError.value = err.message || 'Operation failed'
  } finally {
    submitting.value = false
  }
}

async function handleRenameSubmit() {
  const name = renameColNewName.value.trim()
  if (!name) return

  const origName = originalColData.value.column_name
  if (name === origName) {
    isRenameModalOpen.value = false
    return
  }

  const connId = props.tab.connectionId || connectionsStore.currentConnectionId
  if (!connId || !props.tab.table) return

  const schema = props.tab.schema || 'public'
  const table = props.tab.table
  const sql = `ALTER TABLE "${schema}"."${table}" RENAME COLUMN "${origName}" TO "${name}";`

  try {
    const bindings = connectionsStore.getWailsBindings()
    const result = await bindings.ExecuteQuery(connId, sql, 30)
    if (result.error) {
      throw new Error(result.error)
    }

    uiStore.addNotification({
      type: 'success',
      title: 'Column Renamed',
      message: `Renamed "${origName}" to "${name}"`
    })
    
    isRenameModalOpen.value = false
    loadColumns()
  } catch (err: any) {
    uiStore.addNotification({
      type: 'error',
      title: 'Rename Failed',
      message: err.message || String(err)
    })
  }
}

async function confirmDropColumn(col: any) {
  if (!col) return
  const colName = col.column_name
  
  if (confirm(`Are you sure you want to drop column "${colName}"? This action cannot be undone.`)) {
    const connId = props.tab.connectionId || connectionsStore.currentConnectionId
    if (!connId || !props.tab.table) return

    const schema = props.tab.schema || 'public'
    const table = props.tab.table
    const sql = `ALTER TABLE "${schema}"."${table}" DROP COLUMN "${colName}" CASCADE;`

    try {
      const bindings = connectionsStore.getWailsBindings()
      const result = await bindings.ExecuteQuery(connId, sql, 30)
      if (result.error) {
        throw new Error(result.error)
      }

      uiStore.addNotification({
        type: 'success',
        title: 'Column Dropped',
        message: `Column "${colName}" dropped successfully.`
      })
      
      loadColumns()
    } catch (err: any) {
      uiStore.addNotification({
        type: 'error',
        title: 'Drop Failed',
        message: err.message || String(err)
      })
    }
  }
}

onMounted(() => {
  document.addEventListener('click', closeContextMenu)
  window.addEventListener('keydown', handleKeyDown)
  console.log('DDLView mounted for table:', props.tab.table)
  loadColumns()
})

onUnmounted(() => {
  document.removeEventListener('click', closeContextMenu)
  window.removeEventListener('keydown', handleKeyDown)
  if (tabulator) {
    tabulator.destroy()
  }
  destroyDdlEditor()
})

watch(activeTab, (newTab) => {
  console.log('DDLView activeTab changed to:', newTab)
  if (newTab === 'ddl') {
    loadDDL()
  } else {
    loadColumns()
  }
})

watch(() => props.tab, () => {
  console.log('DDLView props.tab changed:', props.tab.table)
  if (activeTab.value === 'ddl') {
    loadDDL()
  } else {
    loadColumns()
  }
}, { deep: true })
</script>

<style scoped>
/* Scoped overrides to style Tabulator inside DDL Properties view beautifully */
:deep(.tabulator) {
  border: none !important;
  font-size: 11px !important;
}

:deep(.tabulator-col-title) {
  font-weight: 600 !important;
  color: var(--text-primary) !important;
}
</style>
