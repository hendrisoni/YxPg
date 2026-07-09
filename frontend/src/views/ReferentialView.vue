<template>
  <div class="h-full flex flex-col overflow-hidden bg-navy-primary select-none">
    <!-- Toolbar / Dropdown selection -->
    <div class="flex items-center gap-3 px-4 py-2 border-b border-navy-border bg-navy-secondary flex-shrink-0">
      <span class="text-xs font-semibold text-text-primary">Referential Integrity</span>
      
      <!-- Connection Selector -->
      <div class="flex items-center gap-1.5 ml-2">
        <span class="text-[10px] text-text-secondary font-semibold uppercase">Connection:</span>
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
            {{ conn.name }}
          </option>
        </select>
      </div>

      <!-- Constraint Selection (If multiple constraints exist for the dropped table) -->
      <div v-if="constraintsList.length > 1" class="flex items-center gap-1.5">
        <span class="text-[10px] text-text-secondary font-semibold uppercase">Select Relation:</span>
        <select
          :value="selectedConstraintIndex"
          @change="onConstraintChange"
          class="bg-navy-tertiary border border-navy-border text-xs text-text-primary rounded px-2 py-0.5 focus:border-teal-accent focus:outline-none cursor-pointer max-w-[250px] truncate"
        >
          <option v-for="(c, idx) in constraintsList" :key="idx" :value="idx">
            {{ c.constraint_name }} ({{ c.child_table }} -> {{ c.parent_table }})
          </option>
        </select>
      </div>
      
      <div class="flex-1"></div>
      
      <div v-if="selectedConstraint" class="text-[10px] text-text-muted font-mono bg-navy-tertiary/60 border border-navy-border/50 px-2 py-0.5 rounded">
        Active relation: {{ selectedConstraint.constraint_name }}
      </div>
    </div>

    <!-- Main Workspace Content -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Constraint Meta Header Details (Mockup Top Panel) -->
      <div v-if="selectedConstraint" class="bg-navy-secondary border-b border-navy-border p-4 grid grid-cols-2 md:grid-cols-3 gap-x-6 gap-y-3 flex-shrink-0 text-xs">
        <!-- Col 1: Nama Constraint & Table (Child) -->
        <div class="space-y-3">
          <div class="flex flex-col gap-1">
            <span class="text-[10px] text-text-muted font-semibold uppercase">Nama Constraint</span>
            <input
              type="text"
              readonly
              :value="selectedConstraint.constraint_name"
              class="bg-navy-tertiary border border-navy-border/60 rounded px-2 py-1 text-text-primary focus:outline-none font-mono text-[11px]"
            />
          </div>
          <div class="flex flex-col gap-1">
            <span class="text-[10px] text-text-muted font-semibold uppercase">Table (Child)</span>
            <input
              type="text"
              readonly
              :value="`${selectedConstraint.child_schema}.${selectedConstraint.child_table}`"
              class="bg-navy-tertiary border border-navy-border/60 rounded px-2 py-1 text-text-primary focus:outline-none font-mono text-[11px]"
            />
          </div>
          <div class="flex flex-col gap-1">
            <span class="text-[10px] text-text-muted font-semibold uppercase">Kolom (Child)</span>
            <input
              type="text"
              readonly
              :value="selectedConstraint.child_column"
              class="bg-navy-tertiary border border-navy-border/60 rounded px-2 py-1 text-text-primary focus:outline-none font-mono text-[11px]"
            />
          </div>
        </div>

        <!-- Col 2: Table (Parent) & Kolom (Parent) -->
        <div class="space-y-3">
          <div class="flex flex-col gap-1">
            <span class="text-[10px] text-text-muted font-semibold uppercase">Table (Parent)</span>
            <input
              type="text"
              readonly
              :value="`${selectedConstraint.parent_schema}.${selectedConstraint.parent_table}`"
              class="bg-navy-tertiary border border-navy-border/60 rounded px-2 py-1 text-text-primary focus:outline-none font-mono text-[11px]"
            />
          </div>
          <div class="flex flex-col gap-1">
            <span class="text-[10px] text-text-muted font-semibold uppercase">Kolom (Parent)</span>
            <input
              type="text"
              readonly
              :value="selectedConstraint.parent_column"
              class="bg-navy-tertiary border border-navy-border/60 rounded px-2 py-1 text-text-primary focus:outline-none font-mono text-[11px]"
            />
          </div>
          
          <!-- Status -->
          <div class="flex flex-col gap-1">
            <span class="text-[10px] text-text-muted font-semibold uppercase">Status</span>
            <div class="flex items-center gap-1.5 py-1 text-accent-green font-medium select-none">
              <svg class="w-4 h-4 text-accent-green fill-current" viewBox="0 0 24 24">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" />
              </svg>
              <span>OK</span>
            </div>
          </div>
        </div>

        <!-- Col 3: Referential Actions ON UPDATE / ON DELETE -->
        <div class="space-y-3 md:col-span-1 col-span-2">
          <span class="text-[10px] text-text-muted font-semibold uppercase block mb-1">Referential Action</span>
          <div class="flex flex-wrap gap-2 py-1">
            <div class="flex flex-col gap-1">
              <span class="text-[9px] text-text-muted uppercase">On Update</span>
              <span class="px-2 py-1 bg-blue-500/10 border border-blue-500/30 text-blue-400 font-bold rounded text-[10px] uppercase font-mono tracking-wider">
                {{ selectedConstraint.on_update || 'NO ACTION' }}
              </span>
            </div>
            <div class="flex flex-col gap-1">
              <span class="text-[9px] text-text-muted uppercase">On Delete</span>
              <span class="px-2 py-1 bg-amber-500/10 border border-amber-500/30 text-amber-400 font-bold rounded text-[10px] uppercase font-mono tracking-wider">
                {{ selectedConstraint.on_delete || 'NO ACTION' }}
              </span>
            </div>
            <div class="flex flex-col gap-1">
              <span class="text-[9px] text-text-muted uppercase">Match</span>
              <span class="px-2 py-1 bg-gray-500/10 border border-gray-500/30 text-gray-400 font-bold rounded text-[10px] uppercase font-mono tracking-wider">
                {{ selectedConstraint.match_type || 'SIMPLE' }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Relational Canvas Panel (Mockup Bottom Panel) -->
      <div
        ref="canvasRef"
        class="flex-1 relative overflow-hidden bg-[#070a12] p-6 flex flex-col justify-between"
        @dragover.prevent
        @drop="handleCanvasDrop"
      >
        <!-- Canvas grid background patterns -->
        <div class="absolute inset-0 bg-[radial-gradient(#1e293b_1px,transparent_1px)] [background-size:16px_16px] opacity-35 pointer-events-none"></div>

        <!-- Title -->
        <div class="absolute top-4 left-6 z-20 text-xs font-bold text-text-primary select-none flex items-center gap-2">
          <span>Relasi</span>
          <span v-if="selectedConstraint" class="text-[10px] font-normal text-text-muted font-mono">
            ({{ selectedConstraint.child_table }} &lt;-&gt; {{ selectedConstraint.parent_table }})
          </span>
        </div>

        <!-- Drop Instructions / Empty State -->
        <div v-if="!selectedConstraint" class="flex-1 flex flex-col items-center justify-center text-text-muted opacity-50 z-10 pointer-events-none select-none">
          <svg class="w-16 h-16 mb-4 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <rect x="3" y="3" width="6" height="6" rx="1" />
            <rect x="15" y="15" width="6" height="6" rx="1" />
            <path d="M9 6h3a2 2 0 0 1 2 2v4a2 2 0 0 0 2 2h1" />
          </svg>
          <h4 class="text-sm font-semibold text-text-primary mb-1">Referential Integrity visualizer</h4>
          <p class="text-xs text-center max-w-xs leading-relaxed">
            Drag a table from the left database tree sidebar and drop it onto this canvas to inspect and draw its relational Foreign Keys schema.
          </p>
        </div>

        <div v-else class="flex-1 flex items-center justify-between px-12 md:px-24 z-20 relative">
          <!-- Child Table Card -->
          <div
            data-table-type="child"
            class="w-64 bg-navy-secondary border border-navy-border rounded-lg shadow-2xl flex flex-col overflow-hidden"
          >
            <!-- Card Header (Blue for child) -->
            <div class="px-3 py-2 bg-blue-600/10 border-b border-navy-border flex items-center justify-between text-blue-400">
              <div class="flex items-center gap-2 truncate">
                <svg class="w-3.5 h-3.5 flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="3" width="18" height="18" rx="2" /><path d="M3 9h18M3 15h18M9 3v18" />
                </svg>
                <span class="text-xs font-semibold font-mono truncate">
                  {{ selectedConstraint.child_schema }}.{{ selectedConstraint.child_table }} (Child)
                </span>
              </div>
            </div>

            <!-- Card Body / Columns -->
            <div class="max-h-64 overflow-y-auto custom-scrollbar p-1 text-[11px]">
              <div
                v-for="col in childColumns"
                :key="col.column_name"
                :data-column-name="col.column_name"
                class="flex items-center justify-between px-2 py-1.5 rounded hover:bg-navy-hover transition-colors font-mono select-none"
                :class="col.column_name === selectedConstraint.child_column ? 'bg-blue-600/10 text-blue-400 font-bold border border-blue-500/20' : 'text-text-secondary'"
              >
                <div class="flex items-center gap-1.5 truncate">
                  <span v-if="col.is_primary_key" class="text-[9px] font-bold text-yellow-500 px-0.5 bg-yellow-500/10 rounded flex-shrink-0">PK</span>
                  <span class="truncate">{{ col.column_name }}</span>
                </div>
                <span class="text-[9.5px] text-text-muted flex-shrink-0">
                  {{ col.column_name === selectedConstraint.child_column ? '(FK)' : col.data_type }}
                </span>
              </div>
            </div>
          </div>

          <!-- Parent Table Card -->
          <div
            data-table-type="parent"
            class="w-64 bg-navy-secondary border border-navy-border rounded-lg shadow-2xl flex flex-col overflow-hidden"
          >
            <!-- Card Header (Green for parent) -->
            <div class="px-3 py-2 bg-green-600/10 border-b border-navy-border flex items-center justify-between text-green-400">
              <div class="flex items-center gap-2 truncate">
                <svg class="w-3.5 h-3.5 flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="3" width="18" height="18" rx="2" /><path d="M3 9h18M3 15h18M9 3v18" />
                </svg>
                <span class="text-xs font-semibold font-mono truncate">
                  {{ selectedConstraint.parent_schema }}.{{ selectedConstraint.parent_table }} (Parent)
                </span>
              </div>
            </div>

            <!-- Card Body / Columns -->
            <div class="max-h-64 overflow-y-auto custom-scrollbar p-1 text-[11px]">
              <div
                v-for="col in parentColumns"
                :key="col.column_name"
                :data-column-name="col.column_name"
                class="flex items-center justify-between px-2 py-1.5 rounded hover:bg-navy-hover transition-colors font-mono select-none"
                :class="col.column_name === selectedConstraint.parent_column ? 'bg-green-600/10 text-green-400 font-bold border border-green-500/20' : 'text-text-secondary'"
              >
                <div class="flex items-center gap-1.5 truncate">
                  <span v-if="col.is_primary_key" class="text-[9px] font-bold text-yellow-500 px-0.5 bg-yellow-500/10 rounded flex-shrink-0">PK</span>
                  <span class="truncate">{{ col.column_name }}</span>
                </div>
                <span class="text-[9.5px] text-text-muted flex-shrink-0">{{ col.data_type }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- SVG Relationship Overlay Layer -->
        <svg class="absolute inset-0 pointer-events-none w-full h-full z-10">
          <defs>
            <marker id="relation-arrow" viewBox="0 0 10 10" refX="10" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
              <path d="M 0 0 L 10 5 L 0 10 z" fill="#3b82f6" />
            </marker>
          </defs>
          <g v-if="positionsReady && selectedConstraint">
            <path
              :d="getBezierPath(lineStart.x, lineStart.y, lineEnd.x, lineEnd.y)"
              fill="none"
              stroke="#3b82f6"
              stroke-width="2.5"
              marker-end="url(#relation-arrow)"
              class="drop-shadow-[0_0_6px_rgba(59,130,246,0.35)]"
            />
            <!-- Labels (Child side: infinity label, Parent side: 1 label) -->
            <text :x="lineStart.x + 12" :y="lineStart.y - 6" fill="#3b82f6" class="text-[15px] font-bold">∞</text>
            <text :x="lineEnd.x - 16" :y="lineEnd.y - 6" fill="#3b82f6" class="text-[13px] font-bold">1</text>
          </g>
        </svg>

        <!-- Loading spinner overlay -->
        <div v-if="loading" class="absolute inset-0 bg-navy-primary/40 backdrop-blur-sm z-30 flex items-center justify-center">
          <div class="inline-block w-6 h-6 border-2 border-teal-accent border-t-transparent rounded-full animate-spin"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useConnectionsStore } from '../stores/connections'
import { useSchemaStore } from '../stores/schema'
import { useUiStore } from '../stores/ui'
import type { Tab } from '../types'

const props = defineProps<{
  tab: Tab
}>()

const connectionsStore = useConnectionsStore()
const schemaStore = useSchemaStore()
const uiStore = useUiStore()

const selectedConnectionId = ref('')
const loading = ref(false)

const constraintsList = ref<any[]>([])
const selectedConstraintIndex = ref(0)
const selectedConstraint = ref<any>(null)

const childColumns = ref<any[]>([])
const parentColumns = ref<any[]>([])

// Canvas refs and SVG positioning
const canvasRef = ref<HTMLElement | null>(null)
const lineStart = ref({ x: 0, y: 0 })
const lineEnd = ref({ x: 0, y: 0 })
const positionsReady = ref(false)

const currentConnectionName = computed(() => {
  return connectionsStore.currentConnection?.name || 'None'
})

// Query schema tables constraints from information_schema
async function fetchConstraintsForTable(connId: string, schema: string, table: string) {
  loading.value = true
  constraintsList.value = []
  selectedConstraint.value = null
  childColumns.value = []
  parentColumns.value = []
  positionsReady.value = false

  const sql = `
    SELECT
        c.conname AS constraint_name,
        ns1.nspname AS child_schema,
        t1.relname AS child_table,
        a1.attname AS child_column,
        ns2.nspname AS parent_schema,
        t2.relname AS parent_table,
        a2.attname AS parent_column,
        CASE c.confupdtype
            WHEN 'a' THEN 'NO ACTION'
            WHEN 'r' THEN 'RESTRICT'
            WHEN 'c' THEN 'CASCADE'
            WHEN 'n' THEN 'SET NULL'
            WHEN 'd' THEN 'SET DEFAULT'
            ELSE 'NO ACTION'
        END AS on_update,
        CASE c.confdeltype
            WHEN 'a' THEN 'NO ACTION'
            WHEN 'r' THEN 'RESTRICT'
            WHEN 'c' THEN 'CASCADE'
            WHEN 'n' THEN 'SET NULL'
            WHEN 'd' THEN 'SET DEFAULT'
            ELSE 'NO ACTION'
        END AS on_delete
    FROM
        pg_constraint c
        JOIN pg_class t1 ON c.conrelid = t1.oid
        JOIN pg_namespace ns1 ON t1.relnamespace = ns1.oid
        JOIN pg_class t2 ON c.confrelid = t2.oid
        JOIN pg_namespace ns2 ON t2.relnamespace = ns2.oid
        JOIN pg_attribute a1 ON a1.attrelid = c.conrelid AND a1.attnum = c.conkey[1]
        JOIN pg_attribute a2 ON a2.attrelid = c.confrelid AND a2.attnum = c.confkey[1]
    WHERE
        c.contype = 'f'
        AND (
            (ns1.nspname = '${schema}' AND t1.relname = '${table}')
            OR
            (ns2.nspname = '${schema}' AND t2.relname = '${table}')
        );
  `

  try {
    const bindings = connectionsStore.getWailsBindings()
    const result = await bindings.ExecuteQuery(connId, sql, 30)
    if (result.error) {
      throw new Error(result.error)
    }

    const list = parseQueryResult(result)
    if (list.length === 0) {
      uiStore.addNotification({
        type: 'info',
        title: 'No Constraints',
        message: `No foreign keys / referential integrity found for table "${schema}.${table}"`
      })
      return
    }

    constraintsList.value = list
    selectedConstraintIndex.value = 0
    await visualizeConstraint(list[0])
  } catch (err: any) {
    uiStore.addNotification({
      type: 'error',
      title: 'Query Failed',
      message: err.message || String(err)
    })
  } finally {
    loading.value = false
  }
}

async function visualizeConstraint(constraint: any) {
  selectedConstraint.value = constraint
  const connId = selectedConnectionId.value || connectionsStore.currentConnectionId
  if (!connId) return

  loading.value = true
  positionsReady.value = false
  
  try {
    // Load columns in parallel
    const [childCols, parentCols] = await Promise.all([
      schemaStore.loadColumns(connId, constraint.child_schema, constraint.child_table),
      schemaStore.loadColumns(connId, constraint.parent_schema, constraint.parent_table)
    ])
    childColumns.value = childCols
    parentColumns.value = parentCols
    
    await updateLinePositions()
  } catch (err: any) {
    console.error('Failed to load columns for visualizer:', err)
  } finally {
    loading.value = false
  }
}

function onConstraintChange(e: Event) {
  const index = parseInt((e.target as HTMLSelectElement).value)
  selectedConstraintIndex.value = index
  if (constraintsList.value[index]) {
    visualizeConstraint(constraintsList.value[index])
  }
}

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

// Coordinate calculation for SVG link lines
const getAnchorCoords = (tableType: 'child' | 'parent', columnName: string) => {
  const canvas = canvasRef.value
  if (!canvas) return { x: 0, y: 0 }
  
  const el = canvas.querySelector(`[data-table-type="${tableType}"] [data-column-name="${columnName}"]`)
  if (!el) return { x: 0, y: 0 }
  
  const canvasRect = canvas.getBoundingClientRect()
  const elRect = el.getBoundingClientRect()
  
  if (tableType === 'child') {
    return {
      x: elRect.right - canvasRect.left,
      y: elRect.top - canvasRect.top + (elRect.height / 2)
    }
  } else {
    return {
      x: elRect.left - canvasRect.left,
      y: elRect.top - canvasRect.top + (elRect.height / 2)
    }
  }
}

async function updateLinePositions() {
  positionsReady.value = false
  await nextTick()
  setTimeout(() => {
    if (!selectedConstraint.value) return
    const childCol = selectedConstraint.value.child_column
    const parentCol = selectedConstraint.value.parent_column
    
    lineStart.value = getAnchorCoords('child', childCol)
    lineEnd.value = getAnchorCoords('parent', parentCol)
    positionsReady.value = true
  }, 250)
}

function getBezierPath(x1: number, y1: number, x2: number, y2: number) {
  const dx = Math.abs(x2 - x1) * 0.45
  return `M ${x1} ${y1} C ${x1 + dx} ${y1}, ${x2 - dx} ${y2}, ${x2} ${y2}`
}

// Drag & Drop
function handleCanvasDrop(e: DragEvent) {
  e.preventDefault()
  const dataStr = e.dataTransfer?.getData('application/json')
  if (!dataStr) return
  
  try {
    const tableData = JSON.parse(dataStr)
    const connId = tableData.connectionId || selectedConnectionId.value || connectionsStore.currentConnectionId
    
    if (tableData.connectionId) {
      selectedConnectionId.value = tableData.connectionId
    }

    if (connId && tableData.schema && tableData.name) {
      fetchConstraintsForTable(connId, tableData.schema, tableData.name)
    }
  } catch (err: any) {
    console.error('Failed to parse drag drop table json:', err)
  }
}

onMounted(() => {
  window.addEventListener('resize', updateLinePositions)
  
  // If connection store already has a selected connection, preset
  if (connectionsStore.currentConnectionId) {
    selectedConnectionId.value = connectionsStore.currentConnectionId
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', updateLinePositions)
})

watch(() => selectedConnectionId.value, (newConnId) => {
  // Clear layout when switching connection
  if (newConnId) {
    constraintsList.value = []
    selectedConstraint.value = null
    childColumns.value = []
    parentColumns.value = []
    positionsReady.value = false
  }
})
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
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
