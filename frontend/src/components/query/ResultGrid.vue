<template>
  <div class="grid-container h-full flex flex-col overflow-hidden bg-navy-primary">
    <!-- Empty state -->
    <div v-if="!result" class="flex-1 flex items-center justify-center text-text-muted text-xs">
      Run a query to see results
    </div>

    <!-- Error state -->
    <div v-else-if="result.error" class="flex-1 flex flex-col items-center justify-center p-6 overflow-y-auto">
      <div class="max-w-xl w-full text-center space-y-3 bg-navy-secondary/40 border border-navy-border/60 rounded-xl p-6 shadow-lg">
        <svg class="w-10 h-10 mx-auto text-accent-red" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10" /><path d="m15 9-6 6M9 9l6 6" />
        </svg>
        <div>
          <h3 class="text-sm text-accent-red font-semibold">Query Error</h3>
          <p class="text-xs text-text-secondary mt-1.5 font-mono select-text text-left bg-navy-tertiary/60 border border-navy-border p-3 rounded-lg leading-relaxed whitespace-pre-wrap break-all">{{ result.error }}</p>
        </div>

        <!-- Read More / View SQL Toggle -->
        <div v-if="executedSql" class="text-left pt-2 border-t border-navy-border/40">
          <button
            type="button"
            @click="showFullSql = !showFullSql"
            class="text-[11px] text-teal-accent hover:underline flex items-center gap-1 cursor-pointer font-medium"
          >
            <span>{{ showFullSql ? 'Hide executed SQL' : 'Show executed SQL' }}</span>
            <svg
              class="w-3 h-3 transition-transform duration-150"
              :class="showFullSql ? 'rotate-180' : ''"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="m6 9 6 6 6-6" />
            </svg>
          </button>
          
          <div v-if="showFullSql" class="mt-2">
            <pre class="text-[10px] font-mono text-text-muted bg-navy-tertiary/40 border border-navy-border/50 rounded-lg p-3 overflow-x-auto max-h-48 select-text leading-normal whitespace-pre-wrap break-all">{{ executedSql }}</pre>
          </div>
        </div>
      </div>
    </div>

    <!-- Results table -->
    <div v-else-if="result.row_count > 0" class="flex-1 flex flex-col overflow-hidden">
      <!-- Tabulator Container -->
      <div class="flex-1 overflow-hidden relative" :class="showFilters ? 'show-grid-filters' : 'hide-grid-filters'">
        <div ref="tableContainer" class="h-full w-full"></div>
      </div>

      <!-- Footer -->
      <div class="grid-footer flex items-center justify-between border-t border-navy-border px-4 py-2 bg-navy-secondary text-xs text-text-secondary">
        <div class="flex items-center gap-3">
          <!-- Refresh -->
          <button @click="handleRefresh" class="p-1 rounded hover:bg-navy-hover text-text-muted hover:text-text-primary transition-colors" title="Refresh">
            <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67" />
            </svg>
          </button>

          <!-- Freeze Columns Dropdown -->
          <div class="relative flex-shrink-0" ref="freezeRef">
            <button @click.stop="toggleFreezeMenu" class="p-1 rounded hover:bg-navy-hover transition-colors flex items-center justify-center" :class="showFreezeMenu ? 'text-teal-accent bg-teal-accent/10' : 'text-text-muted hover:text-text-primary'" title="Freeze Columns">
              <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="12" y1="17" x2="12" y2="22" />
                <path d="M5 17h14v-1.76a2 2 0 0 0-.44-1.24l-2.78-3.48A2 2 0 0 1 15 9.28V5a1 1 0 0 0-1-1h-4a1 1 0 0 0-1 1v4.28a2 2 0 0 1-.78 1.24l-2.78 3.48a2 2 0 0 0-.44 1.24V17z" />
              </svg>
            </button>
            <div v-if="showFreezeMenu" class="absolute bottom-8 left-0 z-50 w-48 p-2.5 rounded-lg shadow-xl bg-navy-secondary border border-navy-border flex flex-col gap-1.5 animate-fade-in select-none">
              <div class="text-[10px] font-bold text-text-muted uppercase tracking-wider px-1 py-0.5">
                <span>Freeze Columns</span>
              </div>
              <div class="h-px bg-navy-border my-0.5"></div>
              <div v-if="freezeableColumns.length > 0" class="flex flex-col gap-1 max-h-48 overflow-y-auto pr-1 custom-scrollbar">
                <button
                  v-for="col in freezeableColumns"
                  :key="col.field"
                  @click.stop="toggleFreezeColumn(col.field)"
                  class="flex items-center justify-between px-2 py-1.5 rounded hover:bg-navy-hover cursor-pointer text-xs transition-colors w-full text-left"
                  :class="col.frozen ? 'text-teal-accent font-medium' : 'text-text-primary'"
                >
                  <span class="truncate flex-1 text-left">{{ col.title }}</span>
                  <svg v-if="col.frozen" class="w-3 h-3 text-teal-accent flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <line x1="12" y1="17" x2="12" y2="22" />
                    <path d="M5 17h14v-1.76a2 2 0 0 0-.44-1.24l-2.78-3.48A2 2 0 0 1 15 9.28V5a1 1 0 0 0-1-1h-4a1 1 0 0 0-1 1v4.28a2 2 0 0 1-.78 1.24l-2.78 3.48a2 2 0 0 0-.44 1.24V17z" />
                  </svg>
                  <svg v-else class="w-3 h-3 text-text-muted opacity-40 hover:opacity-100 flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <line x1="12" y1="17" x2="12" y2="22" />
                    <path d="M5 17h14v-1.76a2 2 0 0 0-.44-1.24l-2.78-3.48A2 2 0 0 1 15 9.28V5a1 1 0 0 0-1-1h-4a1 1 0 0 0-1 1v4.28a2 2 0 0 1-.78 1.24l-2.78 3.48a2 2 0 0 0-.44 1.24V17z" />
                  </svg>
                </button>
              </div>
              <div v-else class="text-[10px] text-text-muted text-center py-2">
                No columns to freeze
              </div>
            </div>
          </div>

          <!-- Export to XLSX -->
          <button @click="exportToXlsx" class="p-1 rounded hover:bg-navy-hover text-text-muted hover:text-text-primary transition-colors" title="Export to XLSX">
            <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z" />
              <polyline points="14 2 14 8 20 8" />
              <path d="M8 13h2v5H8z" />
              <path d="M12 15h2v3h-2z" />
              <path d="M16 12h2v6h-2z" />
            </svg>
          </button>
        </div>

        <div>
          <span v-if="tab?.type === 'table'">Showing {{ displayedCount }} of {{ totalRows }} rows</span>
          <span v-else>{{ result.row_count }} rows</span>
        </div>
      </div>
    </div>

    <!-- No rows -->
    <div v-else-if="result.row_count === 0" class="flex-1 flex items-center justify-center">
      <div class="text-center">
        <svg class="w-8 h-8 mx-auto text-text-muted mb-2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="3" width="18" height="18" rx="2" /><path d="M3 9h18M3 15h18M9 3v18" />
        </svg>
        <p class="text-sm text-text-muted">No rows returned</p>
        <p v-if="result.duration_ms" class="text-xs text-text-muted mt-1">
          {{ result.duration_ms }}ms
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onUnmounted, onMounted, nextTick } from 'vue'
import * as XLSX from 'xlsx'

if (typeof window !== 'undefined') {
  ;(window as any).XLSX = XLSX
}
import { useConnectionsStore } from '../../stores/connections'
import type { QueryResult, Tab, BrowseOptions } from '../../types'

const props = defineProps<{
  result: QueryResult | null
  tab?: Tab
  executedSql?: string
}>()

const emit = defineEmits(['refresh'])

const connectionsStore = useConnectionsStore()
const displayedCount = ref(0)
const totalRows = ref(0)
const showFullSql = ref(false)

const tableContainer = ref<HTMLElement | null>(null)
let table: any = null

const showFilters = ref(false)
const showFreezeMenu = ref(false)
const freezeRef = ref<HTMLElement | null>(null)
const freezeableColumns = ref<{ field: string; title: string; frozen: boolean }[]>([])

watch(
  () => props.result,
  async (newResult) => {
    showFullSql.value = false
    if (newResult && newResult.row_count > 0 && !newResult.error) {
      await nextTick()
      renderTable(newResult)
    } else {
      destroyTable()
    }
  },
  { deep: true, immediate: true }
)

const closeFreezeMenu = (e: MouseEvent) => {
  if (freezeRef.value && !freezeRef.value.contains(e.target as Node)) {
    showFreezeMenu.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', closeFreezeMenu)
})

onUnmounted(() => {
  destroyTable()
  document.removeEventListener('click', closeFreezeMenu)
})

function toggleFilters() {
  showFilters.value = !showFilters.value
  if (!showFilters.value) {
    table?.clearHeaderFilter()
  }
  setTimeout(() => {
    table?.redraw(true)
  }, 30)
}

function populateFreezeableColumns() {
  if (!table) return
  freezeableColumns.value = table.getColumns()
    .filter((col: any) => {
      const def = col.getDefinition()
      return def.field && def.field !== '__rownum'
    })
    .map((col: any) => {
      const def = col.getDefinition()
      let title = def.field
      if (props.result && props.result.columns) {
        const colIdx = parseInt(def.field.replace('col_', ''))
        if (!isNaN(colIdx) && props.result.columns[colIdx]) {
          title = props.result.columns[colIdx].name
        }
      }
      return {
        field: def.field,
        title: title,
        frozen: !!def.frozen
      }
    })
}

function toggleFreezeMenu() {
  showFreezeMenu.value = !showFreezeMenu.value
  if (showFreezeMenu.value) {
    populateFreezeableColumns()
  }
}

async function toggleFreezeColumn(field: string) {
  if (!table) return
  
  const columns = table.getColumns().filter((col: any) => {
    const def = col.getDefinition()
    return def.field && def.field !== '__rownum'
  })
  
  const targetIdx = columns.findIndex((col: any) => col.getField() === field)
  if (targetIdx === -1) return
  
  const targetCol = columns[targetIdx]
  const isCurrentlyFrozen = !!targetCol.getDefinition().frozen
  
  if (isCurrentlyFrozen) {
    // Unfreeze this column and all columns to the right of it
    for (let i = targetIdx; i < columns.length; i++) {
      await table.updateColumnDefinition(columns[i].getField(), { frozen: false })
    }
  } else {
    // Freeze this column and all columns to the left of it
    for (let i = 0; i <= targetIdx; i++) {
      await table.updateColumnDefinition(columns[i].getField(), { frozen: true })
    }
  }
  
  populateFreezeableColumns()
}

function handleRefresh() {
  emit('refresh')
}

function exportToXlsx() {
  if (!table) return
  const fileName = props.tab?.table ? `${props.tab.table}_export.xlsx` : 'query_export.xlsx'
  table.download('xlsx', fileName, { sheetName: 'Data' })
}

const oidToTypeName: Record<string, string> = {
  "16": "bool",
  "17": "bytea",
  "18": "char",
  "19": "name",
  "20": "int8",
  "21": "int2",
  "23": "int4",
  "24": "regproc",
  "25": "text",
  "26": "oid",
  "27": "tid",
  "28": "xid",
  "29": "cid",
  "114": "json",
  "142": "xml",
  "194": "pg_node_tree",
  "700": "float4",
  "701": "float8",
  "702": "abstime",
  "703": "reltime",
  "704": "tinterval",
  "705": "unknown",
  "718": "circle",
  "790": "money",
  "829": "macaddr",
  "650": "cidr",
  "869": "inet",
  "1000": "bool[]",
  "1001": "bytea[]",
  "1005": "int2[]",
  "1007": "int4[]",
  "1009": "text[]",
  "1014": "char[]",
  "1015": "varchar[]",
  "1016": "int8[]",
  "1021": "float4[]",
  "1022": "float8[]",
  "1028": "oid[]",
  "1042": "bpchar",
  "1043": "varchar",
  "1082": "date",
  "1083": "time",
  "1114": "timestamp",
  "1184": "timestamptz",
  "1186": "interval",
  "1231": "numeric[]",
  "1266": "timetz",
  "1560": "bit",
  "1562": "varbit",
  "1700": "numeric",
  "2278": "void",
  "2950": "uuid",
  "2951": "uuid[]",
  "3802": "jsonb",
  "3807": "jsonb[]",
}

function getTypeName(type: string): string {
  if (!type) return "unknown"
  return oidToTypeName[type] || type
}

function parsePgNumeric(str: string): number | null {
  if (!str) return null
  const match = str.trim().match(/^\{(-?\d+)\s+(-?\d+)\s+(true|false)\s+(\w+)\s+(true|false)\}$/)
  if (match) {
    const intPart = match[1]
    const expPart = parseInt(match[2], 10)
    const isNaNVal = match[3] === 'true'
    const status = match[4]
    const isValid = match[5] === 'true'

    if (isNaNVal || !isValid) return null
    if (status === 'infinity') return Infinity
    if (status === '-infinity') return -Infinity

    return Number(intPart) * Math.pow(10, expPart)
  }
  return null
}

async function renderTable(result: QueryResult) {
  destroyTable()

  if (!tableContainer.value) return

  // Dynamically import Tabulator
  const { TabulatorFull } = await import('tabulator-tables')

  // Build columns from result metadata
  const columns = [
    {
      formatter: 'rownum',
      field: '__rownum',
      title: `<div class="filter-trigger-btn flex items-center justify-center cursor-pointer w-full h-full text-text-secondary hover:text-teal-accent transition-colors" title="Toggle Filters">
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/>
        </svg>
      </div>`,
      titleFormatter: 'html',
      download: false,
      width: 48,
      hozAlign: 'center',
      headerSort: false,
      resizable: false,
      frozen: true,
      headerHozAlign: 'center',
      headerClick: (e: any, column: any) => {
        toggleFilters()
      }
    },
    ...result.columns.map((col, idx) => {
      const isNumeric = (type: string): boolean => {
        if (!type) return false
        const name = getTypeName(type).toLowerCase()
        return (
          name.includes('int') ||
          name.includes('num') ||
          name.includes('decimal') ||
          name.includes('real') ||
          name.includes('double') ||
          name.includes('float') ||
          name.includes('precision') ||
          name.includes('money')
        )
      }
      const numericCol = isNumeric(col.data_type)
      const typeLabel = getTypeName(col.data_type)

      const alignClass = numericCol ? 'items-end text-right' : 'items-start text-left'
      const titleHtml = `<div class="flex flex-col ${alignClass} leading-tight"><span class="font-semibold text-text-primary text-xs">${col.name}</span><span class="text-[10px] text-text-muted mt-0.5 font-normal">${typeLabel}</span></div>`

      return {
        title: titleHtml,
        titleFormatter: "html",
        titleDownload: col.name,
        field: `col_${idx}`,
        headerTooltip: true,
        headerFilter: 'input',
        hozAlign: numericCol ? 'right' : 'left',
        formatter: (cell: any) => {
          const value = cell.getValue()
          if (value === null || value === undefined) {
            cell.getElement().classList.add('cell-null')
            return '<span style="color: #4B5563; font-style: italic;">NULL</span>'
          }

          // Numeric formatting
          if (numericCol) {
            let num = typeof value === 'number' ? value : NaN
            if (typeof value === 'string') {
              const parsed = parsePgNumeric(value)
              if (parsed !== null) {
                num = parsed
              } else {
                num = Number(value)
              }
            }
            if (!isNaN(num)) {
              return new Intl.NumberFormat('en-US', { maximumFractionDigits: 10 }).format(num)
            }
          }

          // Boolean formatting
          if (typeof value === 'boolean') {
            const cls = value ? 'cell-boolean cell-true' : 'cell-boolean cell-false'
            return `<span class="${cls}">${value}</span>`
          }

          // JSON formatting
          if (typeof value === 'string') {
            try {
              const parsed = JSON.parse(value)
              if (typeof parsed === 'object') {
                return `<span style="color: #3B82F6; cursor: pointer;" onclick="this.nextElementSibling.style.display=this.nextElementSibling.style.display==='none'?'block':'none'">{...}</span><pre style="display:none; margin:0; padding:4px; background:#111827; border-radius:4px; font-size:11px; max-width:400px; overflow:auto; color:#94A3B8;">${JSON.stringify(parsed, null, 2)}</pre>`
              }
            } catch {
              // Not JSON
            }

            // Date formatting
            if (/^\d{4}-\d{2}-\d{2}/.test(value)) {
              return `<span style="color: #06B6D4;">${value}</span>`
            }
          }

          return value
        },
        sorter: numericCol ? 'number' : 'string',
        headerSort: true,
      }
    })
  ]

  // Build table config
  const tableConfig: any = {
    columns: columns,
    layout: 'fitDataFill',
    height: '100%',
    selectableRows: 1,
    selectableRowsRangeMode: 'click',
    renderHorizontal: 'virtual',
    renderVertical: 'virtual',
    movableColumns: true,
    resizableColumns: true,
    headerSortTristate: true,
    placeholder: 'No data',
  }

  if (props.tab && props.tab.type === 'table') {
    tableConfig.ajaxURL = "wails://BrowseTable"
    tableConfig.progressiveLoad = "scroll"
    tableConfig.progressiveLoadScrollMargin = 350
    tableConfig.paginationSize = 100
    tableConfig.filterMode = "remote"
    tableConfig.sortMode = "remote"
    
    tableConfig.ajaxRequestFunc = async (url: string, config: any, params: any) => {
      const page = params.page || 1
      const size = params.size || 100
      const connId = props.tab?.connectionId || connectionsStore.currentConnectionId
      if (!connId || !props.tab?.table) {
        return { columns: [], rows: [], row_count: 0, total_count: 0 }
      }
      const bindings = connectionsStore.getWailsBindings()
      const opts: BrowseOptions = {
        page: page,
        page_size: size,
      }

      // Add remote sorting
      if (params.sorters && params.sorters.length > 0) {
        const sorter = params.sorters[0]
        const colIdx = parseInt(sorter.field.replace('col_', ''))
        if (!isNaN(colIdx) && result.columns[colIdx]) {
          opts.sort_by = result.columns[colIdx].name
          opts.sort_order = sorter.dir === 'asc' ? 'ASC' : 'DESC'
        }
      }

      // Add remote filtering
      if (params.filter && params.filter.length > 0) {
        opts.filters = params.filter.map((f: any) => {
          const colIdx = parseInt(f.field.replace('col_', ''))
          const colName = result.columns[colIdx]?.name || ''
          return {
            column: colName,
            operator: f.type === 'like' ? 'ILIKE' : '=',
            value: f.type === 'like' ? `%${f.value}%` : f.value
          }
        }).filter((f: any) => f.column)
      }

      return await bindings.BrowseTable(
        connId,
        props.tab.schema || 'public',
        props.tab.table,
        opts as any
      )
    }

    tableConfig.ajaxResponse = (url: string, params: any, response: any) => {
      if (!response || response.error) {
        throw new Error(response?.error || 'Failed to fetch data')
      }
      totalRows.value = response.total_count

      const transformedData = response.rows.map((row: any[]) => {
        const rowData: Record<string, any> = {}
        row.forEach((val: any, colIdx: number) => {
          rowData[`col_${colIdx}`] = val
        })
        return rowData
      })

      return {
        last_page: Math.max(1, Math.ceil(response.total_count / 100)),
        data: transformedData
      }
    }
  } else {
    // Local data mode for queries
    const localData = result.rows.map((row) => {
      const rowData: Record<string, any> = {}
      row.forEach((val, colIdx) => {
        rowData[`col_${colIdx}`] = val
      })
      return rowData
    })
    tableConfig.data = localData
    tableConfig.progressiveRender = 'scroll'
    tableConfig.progressiveRenderSize = 100
    tableConfig.progressiveRenderMargin = 350
  }

  // Create table
  table = new TabulatorFull(tableContainer.value, tableConfig)

  table.on('dataLoaded', () => {
    displayedCount.value = table.getData().length
  })
}

function destroyTable() {
  if (table) {
    table.destroy()
    table = null
  }
}
</script>
