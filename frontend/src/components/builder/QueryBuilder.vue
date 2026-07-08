<template>
  <div class="h-full flex flex-col overflow-hidden bg-navy-primary select-none">
    <!-- Top tabs and toolbar -->
    <div class="flex items-center justify-between bg-[#0b0f19] border-b border-navy-border h-9 px-3 flex-shrink-0">
      <!-- Tabs -->
      <div class="flex items-center h-full">
        <button
          v-for="tabName in ['Builder', 'Results']"
          :key="tabName"
          @click="currentTab = tabName"
          class="px-4 h-full text-xs font-semibold border-b-2 transition-colors duration-150"
          :class="currentTab === tabName ? 'border-teal-accent text-text-primary' : 'border-transparent text-text-secondary hover:text-text-primary'"
        >
          {{ tabName }}
        </button>
      </div>

      <!-- Action buttons -->
      <div class="flex items-center gap-2">
        <button
          @click="runQuery"
          :disabled="!isTargetConnected || isRunning || !generatedSQL"
          class="flex items-center gap-1.5 px-3 py-1 text-xs bg-teal-accent text-navy-primary rounded font-bold hover:bg-teal-hover transition-colors disabled:opacity-50"
        >
          <svg v-if="!isRunning" class="w-3 h-3 fill-current" viewBox="0 0 24 24">
            <polygon points="5 3 19 12 5 21 5 3" />
          </svg>
          <div v-else class="w-3 h-3 border-2 border-navy-primary border-t-transparent rounded-full animate-spin"></div>
          Run Query
        </button>

        <button
          @click="clearCanvas"
          class="p-1 rounded text-text-muted hover:text-text-primary hover:bg-navy-hover transition-colors"
          title="Clear Canvas"
        >
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 6h18M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
          </svg>
        </button>

        <div class="w-px h-4 bg-navy-border"></div>

        <!-- Connection Selector -->
        <div class="flex items-center gap-1.5">
          <span class="text-[10px] text-text-secondary">Connection:</span>
          <select
            v-model="selectedConnectionId"
            class="bg-navy-tertiary border border-navy-border text-xs text-text-primary rounded px-2 py-0.5 focus:border-teal-accent focus:outline-none cursor-pointer"
          >
            <option value="">Default ({{ currentConnectionName }})</option>
            <option
              v-for="conn in connectionsStore.connections"
              :key="conn.id"
              :value="conn.id"
            >
              {{ conn.name }} {{ connectionsStore.activeConnections.includes(conn.id) ? '●' : '○' }}
            </option>
          </select>
        </div>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- BUILDER CANVAS TAB -->
      <div v-show="currentTab === 'Builder'" class="flex-1 flex flex-col overflow-hidden">
        <!-- Canvas -->
        <div
          ref="canvasRef"
          class="flex-1 relative bg-navy-primary overflow-hidden cursor-default"
          @dragover.prevent
          @drop="handleCanvasDrop"
        >
          <!-- Grid Background (Visual aesthetic) -->
          <div class="absolute inset-0 bg-[radial-gradient(#1e293b_1px,transparent_1px)] [background-size:16px_16px] opacity-30"></div>

          <!-- Empty Canvas Message -->
          <div v-if="tables.length === 0" class="absolute inset-0 flex items-center justify-center pointer-events-none">
            <div class="text-center text-text-muted max-w-sm">
              <svg class="w-12 h-12 mx-auto mb-3 opacity-40 text-accent-amber" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="12" cy="12" r="3" /><path d="M12 3v6m0 6v6m-9-9h6m6 0h6" />
              </svg>
              <p class="text-sm font-semibold text-text-primary mb-1">Visual Query Canvas</p>
              <p class="text-xs">Drag and drop tables from the sidebar here to start building.</p>
              <p class="text-[11px] mt-2 text-text-muted/70">Drag column anchors together to form table JOINs.</p>
            </div>
          </div>

          <!-- SVG Joins and Connections Rendering Layer -->
          <svg class="absolute inset-0 pointer-events-none w-full h-full z-10">
            <!-- Saved Connections -->
            <g v-for="join in joins" :key="join.id">
              <template v-if="getJoinAnchors(join)">
                <path
                  :d="getBezierPath(
                    getJoinAnchors(join)!.from.x,
                    getJoinAnchors(join)!.from.y,
                    getJoinAnchors(join)!.to.x,
                    getJoinAnchors(join)!.to.y
                  )"
                  fill="none"
                  stroke="#00c9a7"
                  stroke-width="2.5"
                  class="drop-shadow-[0_0_4px_rgba(0,201,167,0.4)]"
                />
                <!-- Midpoint Delete Button -->
                <foreignObject
                  :x="getMidpoint(getJoinAnchors(join)!.from.x, getJoinAnchors(join)!.to.x) - 10"
                  :y="getMidpoint(getJoinAnchors(join)!.from.y, getJoinAnchors(join)!.to.y) - 10"
                  width="20"
                  height="20"
                  class="pointer-events-auto"
                >
                  <button
                    @click="removeJoin(join.id)"
                    class="w-5 h-5 flex items-center justify-center bg-accent-red text-white rounded-full border border-navy-primary shadow hover:scale-110 transition-transform cursor-pointer"
                    title="Remove JOIN"
                  >
                    <svg class="w-2.5 h-2.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                      <path d="M18 6 6 18M6 6l12 12" />
                    </svg>
                  </button>
                </foreignObject>
              </template>
            </g>

            <!-- Temporary Connection Line when dragging -->
            <path
              v-if="dragStartAnchor && dragStartAnchorPos"
              :d="getBezierPath(
                dragStartAnchorPos.x,
                dragStartAnchorPos.y,
                tempLineEnd.x,
                tempLineEnd.y
              )"
              fill="none"
              stroke="#fbbf24"
              stroke-width="2"
              stroke-dasharray="4"
              class="drop-shadow-[0_0_3px_rgba(251,191,36,0.5)]"
            />
          </svg>

          <!-- Table Schema Nodes -->
          <div
            v-for="table in tables"
            :key="table.id"
            class="absolute select-none bg-[#090d16] border border-[#202e42] rounded-lg shadow-xl z-20 flex flex-col max-h-[300px]"
            :class="{ 'border-teal-accent': activeTableId === table.id }"
            :style="{ left: table.x + 'px', top: table.y + 'px', zIndex: table.zIndex, width: (table.width || 260) + 'px' }"
            @mousedown="focusTable(table.id)"
          >
            <!-- Card Header -->
            <div
              class="flex items-center gap-1.5 px-2.5 py-1.5 bg-[#05080f] border-b border-[#182335] cursor-move flex-shrink-0"
              @mousedown="startDragTable($event, table)"
            >
              <svg class="w-3.5 h-3.5 text-accent-green flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="18" height="18" rx="2" /><path d="M3 9h18M3 15h18M9 3v18" />
              </svg>
              <!-- Select All Checkbox -->
              <input
                type="checkbox"
                :checked="table.columns.every(c => c.selected)"
                :indeterminate="table.columns.some(c => c.selected) && !table.columns.every(c => c.selected)"
                @change="toggleSelectAll(table)"
                @mousedown.stop
                class="w-3 h-3 rounded border-[#334155] bg-navy-tertiary text-teal-accent focus:ring-0 cursor-pointer flex-shrink-0"
                title="Select / Deselect All"
              />
              <span class="text-xs font-semibold text-text-primary truncate flex-1" :title="`${table.schema}.${table.name}`">
                {{ table.schema }}.{{ table.name }}
              </span>
              <span class="text-[9px] text-[#6b7280] bg-[#111827] px-1 py-0.5 rounded font-mono font-bold scale-90">
                {{ table.alias }}
              </span>

              <!-- Collapse Button -->
              <button
                @click.stop="table.collapsed = !table.collapsed; recalculatePositions()"
                class="p-0.5 rounded hover:bg-navy-hover text-text-muted hover:text-text-primary cursor-pointer"
              >
                <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line v-if="!table.collapsed" x1="5" y1="12" x2="19" y2="12" />
                  <path v-else d="M12 5v14M5 12h14" />
                </svg>
              </button>

              <!-- Close Button -->
              <button
                @click.stop="removeTable(table.id)"
                class="p-0.5 rounded hover:bg-navy-hover text-text-muted hover:text-accent-red cursor-pointer"
              >
                <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M18 6 6 18M6 6l12 12" />
                </svg>
              </button>
            </div>

            <!-- Column List -->
            <div v-show="!table.collapsed" class="flex-1 overflow-y-auto custom-scrollbar py-1">
              <div
                v-for="col in table.columns"
                :key="col.name"
                class="relative flex items-center gap-2 pl-6 pr-6 py-1 text-xs hover:bg-[#111927] group"
                :class="{ 'bg-[#0c1a2e]': isConnectionHover(table.id, col.name) }"
              >
                <!-- Left Anchor (incoming connections) -->
                <div
                  :id="`anchor-l-${table.id}-${col.name}`"
                  class="w-2.5 h-2.5 rounded-full border border-[#334155] bg-[#0b0f19] absolute left-1.5 top-1/2 -translate-y-1/2 cursor-crosshair hover:border-teal-accent hover:bg-teal-accent z-30 flex items-center justify-center"
                  @mousedown.stop.prevent="startConnection($event, table.id, col.name, 'left')"
                >
                  <div class="w-1 h-1 rounded-full bg-[#334155] group-hover:bg-teal-accent"></div>
                </div>

                <!-- Select Checkbox -->
                <input
                  type="checkbox"
                  v-model="col.selected"
                  class="w-3.5 h-3.5 rounded border-[#202e42] bg-navy-tertiary text-teal-accent focus:ring-0 cursor-pointer"
                />

                <!-- Primary Key Indicator -->
                <span v-if="col.isPrimaryKey" class="text-accent-amber text-[10px]" title="Primary Key">🔑</span>

                <!-- Column Name -->
                <span class="text-text-secondary truncate flex-1 font-mono text-[11px] group-hover:text-text-primary">
                  {{ col.name }}
                </span>

                <!-- Data Type Badge -->
                <span
                  class="text-[10px] font-bold font-mono w-5 text-right select-none pr-1"
                  :style="getDataTypeColor(col.type)"
                  :title="col.type"
                >
                  {{ getDataTypeLetter(col.type) }}
                </span>

                <!-- Right Anchor (outgoing connections) -->
                <div
                  :id="`anchor-r-${table.id}-${col.name}`"
                  class="w-2.5 h-2.5 rounded-full border border-teal-accent bg-[#0b0f19] absolute right-2 top-1/2 -translate-y-1/2 cursor-crosshair hover:bg-teal-accent z-30 flex items-center justify-center"
                  @mousedown.stop.prevent="startConnection($event, table.id, col.name, 'right')"
                >
                  <div class="w-1 h-1 rounded-full bg-teal-accent"></div>
                </div>
              </div>
            </div>

            <!-- Resize Handle -->
            <div
              v-show="!table.collapsed"
              class="absolute bottom-0 right-0 w-3.5 h-3.5 cursor-se-resize z-30 flex items-center justify-center text-text-muted hover:text-teal-accent bg-[#05080f] rounded-tl"
              @mousedown.stop.prevent="startResizeTable($event, table)"
              title="Drag to resize table width"
            >
              <svg class="w-2.5 h-2.5 opacity-60 hover:opacity-100" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <line x1="17" y1="21" x2="21" y2="17" />
                <line x1="12" y1="21" x2="21" y2="12" />
              </svg>
            </div>
          </div>
        </div>

        <!-- SQL Editor display panel -->
        <div class="h-40 border-t border-navy-border bg-[#070a12] flex flex-col overflow-hidden flex-shrink-0">
          <div class="flex items-center justify-between px-3 py-1.5 border-b border-navy-border bg-navy-secondary text-[10px] text-text-muted uppercase tracking-wider font-bold">
            <span>Generated SQL</span>
            <button @click="copySQL" class="hover:text-text-primary text-xs flex items-center gap-1 font-semibold normal-case cursor-pointer">
              <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="9" y="9" width="13" height="13" rx="2" /><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
              </svg>
              Copy
            </button>
          </div>
          <div class="flex-1 p-3 overflow-auto font-mono text-xs text-[#a5b4fc] bg-[#030712] whitespace-pre-wrap select-text">
            {{ generatedSQL || '-- Drag tables and check columns to generate query' }}
          </div>
        </div>
      </div>

      <!-- RESULTS TAB -->
      <div v-show="currentTab === 'Results'" class="flex-1 flex flex-col overflow-hidden">
        <ResultGrid :result="queryResult" :tab="tab" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useConnectionsStore } from '../../stores/connections'
import { useTabsStore } from '../../stores/tabs'
import { useUiStore } from '../../stores/ui'
import { useWorkspaceStore } from '../../stores/workspace'
import ResultGrid from '../query/ResultGrid.vue'
import type { Tab, QueryResult } from '../../types'

const props = defineProps<{
  tab: Tab
}>()

const connectionsStore = useConnectionsStore()
const tabsStore = useTabsStore()
const uiStore = useUiStore()
const workspaceStore = useWorkspaceStore()

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

const currentTab = ref('Builder')
const isRunning = ref(false)
const queryResult = ref<QueryResult | null>(null)

// Tables and Joins layout state
interface CanvasTable {
  id: string
  name: string
  schema: string
  alias: string
  x: number
  y: number
  width?: number
  zIndex: number
  collapsed: boolean
  columns: Array<{
    name: string
    type: string
    isPrimaryKey: boolean
    selected: boolean
  }>
}

interface CanvasJoin {
  id: string
  fromTableId: string
  fromColumn: string
  toTableId: string
  toColumn: string
  type: 'INNER' | 'LEFT' | 'RIGHT'
}

const tables = ref<CanvasTable[]>([])
const joins = ref<CanvasJoin[]>([])

let aliasCounter = 0
let maxZIndex = 10

function getDataTypeLetter(type: string): string {
  if (!type) return '?'
  return type.trim().charAt(0).toUpperCase()
}

function getDataTypeColor(type: string): string {
  if (!type) return 'color: #9ca3af;'
  const t = type.toLowerCase()
  if (t.includes('int') || t.includes('num') || t.includes('double') || t.includes('float') || t.includes('real') || t.includes('decimal')) {
    return 'color: #34d399;' // emerald green for numbers
  }
  if (t.includes('char') || t.includes('text') || t.includes('varchar')) {
    return 'color: #60a5fa;' // blue for strings
  }
  if (t.includes('date') || t.includes('time') || t.includes('timestamp')) {
    return 'color: #fb7185;' // rose for dates/times
  }
  if (t.includes('bool')) {
    return 'color: #a78bfa;' // purple for booleans
  }
  return 'color: #fbbf24;' // amber for other/json/etc.
}

function getAlphabetAlias(index: number): string {
  let alias = ''
  let temp = index - 1
  while (temp >= 0) {
    alias = String.fromCharCode(97 + (temp % 26)) + alias
    temp = Math.floor(temp / 26) - 1
  }
  return alias || 'a'
}
const activeTableId = ref<string | null>(null)

// Dragging tables connection state
const canvasRef = ref<HTMLElement | null>(null)
const dragStartAnchor = ref<{ tableId: string; columnName: string } | null>(null)
const tempLineEnd = ref({ x: 0, y: 0 })
const anchorPositions = ref<Record<string, { x: number; y: number }>>({})

const dragStartAnchorPos = computed(() => {
  if (!dragStartAnchor.value) return null
  const { tableId, columnName } = dragStartAnchor.value
  // Prefer right anchor for the drag start line
  return anchorPositions.value[`${tableId}:${columnName}:r`] 
    || anchorPositions.value[`${tableId}:${columnName}:l`] 
    || null
})

// Query generation computed sql
const generatedSQL = computed(() => {
  if (tables.value.length === 0) return ''

  const selectParts: string[] = []
  
  // Assemble columns
  for (const table of tables.value) {
    const selectedCols = table.columns.filter(c => c.selected)
    if (selectedCols.length > 0) {
      for (const col of selectedCols) {
        selectParts.push(`${table.alias}.${col.name}`)
      }
    }
  }

  // If no columns checked, default to selecting all columns from all tables
  if (selectParts.length === 0) {
    for (const table of tables.value) {
      selectParts.push(`${table.alias}.*`)
    }
  }

  // Resolve Join graph starting from the first table added
  const joinedTableIds = new Set<string>()
  const firstTable = tables.value[0]
  let fromClause = `${firstTable.schema}.${firstTable.name} ${firstTable.alias}`
  joinedTableIds.add(firstTable.id)

  let queryBuilderList = [...tables.value]
  let checkLoop = true
  let iterations = 0

  while (checkLoop && iterations < tables.value.length) {
    checkLoop = false
    iterations++
    
    for (let i = 0; i < queryBuilderList.length; i++) {
      const table = queryBuilderList[i]
      if (joinedTableIds.has(table.id)) continue

      // Look for a join relationship linking this table to any already joined table
      const matchingJoin = joins.value.find(j => {
        return (j.fromTableId === table.id && joinedTableIds.has(j.toTableId)) ||
               (j.toTableId === table.id && joinedTableIds.has(j.fromTableId))
      })

      if (matchingJoin) {
        const otherTableId = matchingJoin.fromTableId === table.id ? matchingJoin.toTableId : matchingJoin.fromTableId
        const otherTable = tables.value.find(t => t.id === otherTableId)
        
        if (otherTable) {
          const colSelf = matchingJoin.fromTableId === table.id ? matchingJoin.fromColumn : matchingJoin.toColumn
          const colOther = matchingJoin.fromTableId === table.id ? matchingJoin.toColumn : matchingJoin.fromColumn
          
          fromClause += `\nLEFT JOIN ${table.schema}.${table.name} ${table.alias} ON ${table.alias}.${colSelf} = ${otherTable.alias}.${colOther}`
          joinedTableIds.add(table.id)
          checkLoop = true
        }
      }
    }
  }

  // Comma joins fallback for remaining standalone tables (if any)
  for (const table of tables.value) {
    if (!joinedTableIds.has(table.id)) {
      fromClause += `, ${table.schema}.${table.name} ${table.alias}`
      joinedTableIds.add(table.id)
    }
  }

  return `SELECT\n  ${selectParts.join(',\n  ')}\nFROM ${fromClause}`
})

// Recursively find database objects in Pinia Tree Nodes
function findNodeInTree(nodes: any[], id: string): any | null {
  for (const node of nodes) {
    if (node.id === id) return node
    if (node.children && node.children.length > 0) {
      const found = findNodeInTree(node.children, id)
      if (found) return found
    }
  }
  return null
}

// Canvas Drag-and-drop table loader
async function handleCanvasDrop(e: DragEvent) {
  e.preventDefault()
  const nodeId = e.dataTransfer?.getData('text/plain')
  if (!nodeId) return

  const node = findNodeInTree(workspaceStore.workspaceTree, nodeId)
  if (!node || (node.type !== 'table' && node.type !== 'view')) return

  const schemaName = node.data?.schema || 'public'
  const tableName = node.data?.name || node.data?.table || node.label
  const connId = node.data?.connectionId || connectionsStore.currentConnectionId

  if (!connId) return

  if (node.data?.connectionId) {
    tabsStore.updateTab(props.tab.id, { connectionId: node.data.connectionId })
  }

  try {
    const bindings = connectionsStore.getWailsBindings()
    const columnsData = await bindings.GetColumns(connId, schemaName, tableName)

    const canvas = canvasRef.value
    let dropX = 100
    let dropY = 100
    if (canvas) {
      const rect = canvas.getBoundingClientRect()
      dropX = Math.max(10, e.clientX - rect.left - 120)
      dropY = Math.max(10, e.clientY - rect.top - 20)
    }

    aliasCounter++
    const alias = getAlphabetAlias(aliasCounter)

    const newTable: CanvasTable = {
      id: `table_${Date.now()}_${Math.random().toString(36).substr(2, 5)}`,
      name: tableName,
      schema: schemaName,
      alias: alias,
      x: dropX,
      y: dropY,
      width: 260,
      zIndex: maxZIndex + 1,
      collapsed: false,
      columns: columnsData.map((col: any) => ({
        name: col.column_name,
        type: col.data_type,
        isPrimaryKey: col.is_primary_key,
        selected: false
      }))
    }

    tables.value.push(newTable)
    maxZIndex++
    recalculatePositions()
  } catch (err: any) {
    uiStore.addNotification({
      type: 'error',
      title: 'Failed to load table schema',
      message: err.message,
    })
  }
}

// Draggable table card actions
const draggingTable = ref<CanvasTable | null>(null)
let dragOffset = { x: 0, y: 0 }

function focusTable(tableId: string) {
  activeTableId.value = tableId
  const table = tables.value.find(t => t.id === tableId)
  if (table) {
    maxZIndex++
    table.zIndex = maxZIndex
  }
}

function startDragTable(e: MouseEvent, table: CanvasTable) {
  focusTable(table.id)
  draggingTable.value = table
  dragOffset = {
    x: e.clientX - table.x,
    y: e.clientY - table.y
  }
  window.addEventListener('mousemove', dragTable)
  window.addEventListener('mouseup', endDragTable)
}

function dragTable(e: MouseEvent) {
  if (!draggingTable.value) return
  const canvas = canvasRef.value
  if (!canvas) return
  const rect = canvas.getBoundingClientRect()

  let newX = e.clientX - dragOffset.x
  let newY = e.clientY - dragOffset.y

  draggingTable.value.x = Math.max(10, Math.min(rect.width - 250, newX))
  draggingTable.value.y = Math.max(10, Math.min(rect.height - 50, newY))

  recalculatePositions()
}

function endDragTable() {
  draggingTable.value = null
  window.removeEventListener('mousemove', dragTable)
  window.removeEventListener('mouseup', endDragTable)
}

const resizingTable = ref<CanvasTable | null>(null)
let resizeStartWidth = 0
let resizeStartMouseX = 0

function startResizeTable(e: MouseEvent, table: CanvasTable) {
  focusTable(table.id)
  resizingTable.value = table
  resizeStartWidth = table.width || 260
  resizeStartMouseX = e.clientX
  window.addEventListener('mousemove', resizeTable)
  window.addEventListener('mouseup', endResizeTable)
}

function resizeTable(e: MouseEvent) {
  if (!resizingTable.value) return
  const deltaX = e.clientX - resizeStartMouseX
  const newWidth = Math.max(200, Math.min(600, resizeStartWidth + deltaX))
  resizingTable.value.width = newWidth
  recalculatePositions()
}

function endResizeTable() {
  resizingTable.value = null
  window.removeEventListener('mousemove', resizeTable)
  window.removeEventListener('mouseup', endResizeTable)
}

function removeTable(id: string) {
  tables.value = tables.value.filter(t => t.id !== id)
  joins.value = joins.value.filter(j => j.fromTableId !== id && j.toTableId !== id)
  recalculatePositions()
}

function toggleSelectAll(table: CanvasTable) {
  const allSelected = table.columns.every(c => c.selected)
  table.columns.forEach(c => c.selected = !allSelected)
}

function clearCanvas() {
  tables.value = []
  joins.value = []
  aliasCounter = 0
}

// Connection dragging and anchoring
const hoverTarget = ref<{ tableId: string; columnName: string } | null>(null)

function isConnectionHover(tableId: string, colName: string) {
  if (!hoverTarget.value || !dragStartAnchor.value) return false
  return hoverTarget.value.tableId === tableId && hoverTarget.value.columnName === colName
}

function startConnection(e: MouseEvent, tableId: string, columnName: string, side: 'left' | 'right') {
  e.stopPropagation()
  e.preventDefault()
  
  dragStartAnchor.value = { tableId, columnName }

  const canvas = canvasRef.value
  if (!canvas) return
  const rect = canvas.getBoundingClientRect()
  tempLineEnd.value = {
    x: e.clientX - rect.left,
    y: e.clientY - rect.top
  }

  window.addEventListener('mousemove', dragConnection)
  window.addEventListener('mouseup', endConnection)
}

function dragConnection(e: MouseEvent) {
  const canvas = canvasRef.value
  if (!canvas) return
  const rect = canvas.getBoundingClientRect()
  const mx = e.clientX - rect.left
  const my = e.clientY - rect.top
  tempLineEnd.value = { x: mx, y: my }

  // Find nearest anchor for hover highlight
  const nearest = findNearestAnchor(mx, my)
  hoverTarget.value = nearest
}

function findNearestAnchor(mx: number, my: number): { tableId: string; columnName: string } | null {
  if (!dragStartAnchor.value) return null
  const SNAP_RADIUS = 30
  let closest: { tableId: string; columnName: string; dist: number } | null = null

  for (const table of tables.value) {
    if (table.id === dragStartAnchor.value.tableId) continue
    for (const col of table.columns) {
      // Check both left and right anchors
      for (const side of ['l', 'r']) {
        const key = `${table.id}:${col.name}:${side}`
        const pos = anchorPositions.value[key]
        if (!pos) continue
        const dist = Math.sqrt((mx - pos.x) ** 2 + (my - pos.y) ** 2)
        if (dist < SNAP_RADIUS && (!closest || dist < closest.dist)) {
          closest = { tableId: table.id, columnName: col.name, dist }
        }
      }
    }
  }
  return closest ? { tableId: closest.tableId, columnName: closest.columnName } : null
}

function endConnection(e: MouseEvent) {
  if (dragStartAnchor.value) {
    const canvas = canvasRef.value
    if (canvas) {
      const rect = canvas.getBoundingClientRect()
      const mx = e.clientX - rect.left
      const my = e.clientY - rect.top
      const target = findNearestAnchor(mx, my)

      if (target && target.tableId !== dragStartAnchor.value.tableId) {
        // Check for duplicate join
        const exists = joins.value.some(j =>
          (j.fromTableId === dragStartAnchor.value!.tableId && j.fromColumn === dragStartAnchor.value!.columnName && j.toTableId === target.tableId && j.toColumn === target.columnName) ||
          (j.toTableId === dragStartAnchor.value!.tableId && j.toColumn === dragStartAnchor.value!.columnName && j.fromTableId === target.tableId && j.fromColumn === target.columnName)
        )

        if (!exists) {
          joins.value.push({
            id: `join_${Date.now()}`,
            fromTableId: dragStartAnchor.value.tableId,
            fromColumn: dragStartAnchor.value.columnName,
            toTableId: target.tableId,
            toColumn: target.columnName,
            type: 'INNER'
          })

          uiStore.addNotification({
            type: 'success',
            title: 'JOIN Created',
            message: `${dragStartAnchor.value.columnName} ↔ ${target.columnName}`,
          })
        }
      }
    }
  }

  dragStartAnchor.value = null
  hoverTarget.value = null
  window.removeEventListener('mousemove', dragConnection)
  window.removeEventListener('mouseup', endConnection)
}

function removeJoin(joinId: string) {
  joins.value = joins.value.filter(j => j.id !== joinId)
}

// Coordinate calculations and position updates
function recalculatePositions() {
  nextTick(() => {
    const canvas = canvasRef.value
    if (!canvas) return
    const canvasRect = canvas.getBoundingClientRect()

    const newPositions: Record<string, { x: number; y: number }> = {}

    for (const table of tables.value) {
      if (table.collapsed) {
        for (const col of table.columns) {
          // Collapsed: all anchors point to the right-side of header
          newPositions[`${table.id}:${col.name}:r`] = {
            x: table.x + 240,
            y: table.y + 16
          }
          newPositions[`${table.id}:${col.name}:l`] = {
            x: table.x,
            y: table.y + 16
          }
        }
      } else {
        for (const col of table.columns) {
          for (const side of ['l', 'r']) {
            const anchorId = `anchor-${side}-${table.id}-${col.name}`
            const el = document.getElementById(anchorId)
            if (el) {
              const rect = el.getBoundingClientRect()
              newPositions[`${table.id}:${col.name}:${side}`] = {
                x: rect.left - canvasRect.left + rect.width / 2,
                y: rect.top - canvasRect.top + rect.height / 2
              }
            }
          }
        }
      }
    }
    anchorPositions.value = newPositions
  })
}

// Pick the best anchor pair (left/right) for a join line
function getJoinAnchors(join: CanvasJoin): { from: { x: number; y: number }; to: { x: number; y: number } } | null {
  const fromR = anchorPositions.value[`${join.fromTableId}:${join.fromColumn}:r`]
  const fromL = anchorPositions.value[`${join.fromTableId}:${join.fromColumn}:l`]
  const toR = anchorPositions.value[`${join.toTableId}:${join.toColumn}:r`]
  const toL = anchorPositions.value[`${join.toTableId}:${join.toColumn}:l`]

  if (!fromR || !fromL || !toR || !toL) return null

  // Pick the pair with shortest horizontal distance
  const pairs = [
    { from: fromR, to: toL },
    { from: fromR, to: toR },
    { from: fromL, to: toL },
    { from: fromL, to: toR },
  ]

  let best = pairs[0]
  let bestDist = Math.abs(best.from.x - best.to.x) + Math.abs(best.from.y - best.to.y)

  for (let i = 1; i < pairs.length; i++) {
    const d = Math.abs(pairs[i].from.x - pairs[i].to.x) + Math.abs(pairs[i].from.y - pairs[i].to.y)
    if (d < bestDist) {
      best = pairs[i]
      bestDist = d
    }
  }

  return best
}

// Line rendering path helpers
function getBezierPath(x1: number, y1: number, x2: number, y2: number) {
  const dx = Math.abs(x2 - x1) * 0.5
  return `M ${x1} ${y1} C ${x1 + dx} ${y1}, ${x2 - dx} ${y2}, ${x2} ${y2}`
}

function getMidpoint(val1: number, val2: number) {
  return val1 + (val2 - val1) / 2
}

function copySQL() {
  if (generatedSQL.value) {
    navigator.clipboard.writeText(generatedSQL.value)
    uiStore.addNotification({
      type: 'info',
      title: 'Copied',
      message: 'SQL query copied to clipboard',
    })
  }
}

// SQL Execution
async function runQuery() {
  const connId = props.tab.connectionId || connectionsStore.currentConnectionId
  if (!connId || !generatedSQL.value) return
  isRunning.value = true
  try {
    const bindings = connectionsStore.getWailsBindings()
    const result = await bindings.ExecuteQuery(
      connId,
      generatedSQL.value,
      30
    )
    queryResult.value = result
    if (result.error) {
      uiStore.addNotification({ type: 'error', title: 'Query Error', message: result.error })
    } else {
      uiStore.addNotification({ type: 'success', title: 'Success', message: `${result.row_count} rows returned` })
      currentTab.value = 'Results'
    }
  } catch (e: any) {
    uiStore.addNotification({ type: 'error', title: 'Failed', message: e.message })
  } finally {
    isRunning.value = false
  }
}

// Watchers for reactive layout positions
watch(() => tables.value, () => {
  recalculatePositions()
}, { deep: true })

onMounted(() => {
  recalculatePositions()
  window.addEventListener('resize', recalculatePositions)
})

onUnmounted(() => {
  window.removeEventListener('resize', recalculatePositions)
})
</script>

<style scoped>
/* Scroller visual aesthetic overrides */
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #202e42;
  border-radius: 9999px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #334155;
}
</style>
