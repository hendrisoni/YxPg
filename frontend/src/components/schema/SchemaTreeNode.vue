<template>
  <div
    @dragover="handleWrapperDragOver($event)"
    @drop="handleWrapperDrop($event)"
  >
    <div
      class="flex items-center gap-1 pl-2 pr-12 py-1 cursor-pointer hover:bg-navy-hover transition-colors group relative border-l-2"
      :class="[
        node.id === selectedNodeId ? 'bg-navy-hover border-teal-accent text-text-primary' : 'border-transparent text-text-secondary',
        dragOver && node.type === 'category' ? 'bg-teal-accent/10 border-teal-accent' : ''
      ]"
      :style="{ paddingLeft: (level * 16 + 8) + 'px' }"
      draggable="true"
      @dragstart="handleDragStart($event, node.id)"
      @dragover="handleDragOver($event)"
      @dragenter="handleDragEnter($event)"
      @dragleave="handleDragLeave($event)"
      @drop="handleDrop($event, node.id)"
      @click="handleSelect"
      @dblclick="handleDoubleClick"
      @contextmenu.prevent="showContextMenu"
      :title="node.data ? `${node.data.connectionName || 'Database'} / ${node.data.database}.${node.data.schema}` : node.label"
    >
      <!-- Expand arrow -->
      <svg
        v-if="hasChildren"
        class="w-3 h-3 transition-transform duration-150 flex-shrink-0"
        :class="{ 'rotate-90': expanded }"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        @click.stop="toggleExpand"
      >
        <path d="M9 18l6-6-6-6" />
      </svg>
      <div v-else class="w-3 h-3 flex-shrink-0"></div>

      <!-- Icon -->
      <component :is="iconComponent" class="w-4 h-4 flex-shrink-0" :class="iconColor" />

      <!-- Label -->
      <span class="text-xs truncate flex-1" :class="labelColor">
        {{ node.label }}
      </span>


      
     

      <!-- Delete hover button -->
      <button
        @click.stop="$emit('delete-node', node.id)"
        class="opacity-0 group-hover:opacity-100 absolute right-7 top-1/2 -translate-y-1/2 p-0.5 rounded bg-navy-hover hover:bg-red-950/40 text-red-700 hover:text-red-500 transition-all z-10"
        title="Remove"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M18 6 6 18M6 6l12 12" />
        </svg>
      </button>
       <!-- Connection Initial Badge -->
      <span
        v-if="connectionInitial"
        class="text-[9px] font-bold px-1 py-0.5 rounded absolute right-1.5 top-1/2 -translate-y-1/2 transition-all z-10"
        :style="{ backgroundColor: connectionInitialsColor, color: ['#00C9A7', '#F59E0B', '#10B981', '#06B6D4', '#14B8A6'].includes(connectionInitialsColor) ? '#050505' : '#ffffff' }"
        :title="'Connection: ' + node.data.connectionName"
      >
        {{ connectionInitial }}
      </span>
    </div>

    <!-- Children -->
    <div v-if="expanded && node.children">
      <SchemaTreeNode
        v-for="child in node.children"
        :key="child.id"
        :node="child"
        :level="level + 1"
        :selected-node-id="selectedNodeId"
        @select-node="$emit('select-node', $event)"
        @open-table="(schema, name, connId) => $emit('open-table', schema, name, connId)"
        @open-query="(schema, name, connId) => $emit('open-query', schema, name, connId)"
        @copy-name="$emit('copy-name', $event)"
        @copy-select="$emit('copy-select', $event)"
        @view-ddl="(schema, name, connId) => $emit('view-ddl', schema, name, connId)"
        @drop-table="$emit('drop-table', $event)"
        @delete-node="$emit('delete-node', $event)"
        @rename-node="(id, label) => $emit('rename-node', id, label)"
      />
    </div>

    <!-- Context Menu -->
    <ContextMenu
      :show="showMenu"
      :x="menuX"
      :y="menuY"
      :items="contextMenuItems"
      @close="showMenu = false"
      @select="handleMenuSelect"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, h } from 'vue'
import type { TreeNode } from '../../types'
import { useWorkspaceStore } from '../../stores/workspace'
import { useConnectionsStore } from '../../stores/connections'
import { useSchemaStore } from '../../stores/schema'
import { useTabsStore } from '../../stores/tabs'
import ContextMenu from '../shared/ContextMenu.vue'
import type { ContextMenuItem } from '../shared/ContextMenu.vue'

const props = defineProps<{
  node: TreeNode
  level: number
  selectedNodeId?: string | null
}>()

const emit = defineEmits([
  'open-table',
  'open-query',
  'copy-name',
  'copy-select',
  'view-ddl',
  'drop-table',
  'delete-node',
  'rename-node',
  'select-node',
])

const expanded = ref(props.node.expanded !== undefined ? props.node.expanded : false)
watch(() => props.node.expanded, (newVal) => {
  if (newVal !== undefined) {
    expanded.value = newVal
  }
})
const dragOver = ref(false)
const workspaceStore = useWorkspaceStore()
const connectionsStore = useConnectionsStore()
const schemaStore = useSchemaStore()
const tabsStore = useTabsStore()

const showMenu = ref(false)
const menuX = ref(0)
const menuY = ref(0)

const hasChildren = computed(() => {
  return props.node.children && props.node.children.length > 0
})

function toggleExpand() {
  if (hasChildren.value) {
    expanded.value = !expanded.value
    props.node.expanded = expanded.value
    workspaceStore.saveWorkspace()
  }
}

const iconComponent = computed(() => {
  const icons: Record<string, any> = {
    database: DatabaseIcon,
    table: TableIcon,
    view: ViewIcon,
    function: FunctionIcon,
    sequence: SequenceIcon,
    column: ColumnIcon,
    index: IndexIcon,
    type: TypeIcon,
    category: FolderIcon,
  }
  return icons[props.node.icon || props.node.type] || TableIcon
})

const iconColor = computed(() => {
  const colors: Record<string, string> = {
    database: 'text-accent-blue',
    table: 'text-accent-green',
    view: 'text-accent-amber',
    function: 'text-purple-400',
    sequence: 'text-cyan-400',
    column: 'text-text-secondary',
    index: 'text-orange-400',
    type: 'text-pink-400',
    category: 'text-amber-400',
  }
  return colors[props.node.icon || props.node.type] || 'text-text-secondary'
})

const connectionInitial = computed(() => {
  const name = props.node.data?.connectionName
  if (!name) return ''
  return name.trim().charAt(0).toUpperCase()
})

const connectionInitialsColor = computed(() => {
  const connId = props.node.data?.connectionId
  const connName = props.node.data?.connectionName
  
  let conn = null
  if (connId) {
    conn = connectionsStore.connections.find(c => c.id === connId)
  }
  if (!conn && connName) {
    conn = connectionsStore.connections.find(c => c.name === connName)
  }
  
  // If the user has explicitly customized the connection color (not the default Teal '#00C9A7'), use it
  if (conn?.color && conn.color !== '#00C9A7') return conn.color
  
  // Otherwise, use a deterministic color based on the initial letter
  const initial = connectionInitial.value
  if (!initial) return '#3B82F6'
  
  const colors = [
    '#3B82F6', // Blue
    '#10B981', // Green
    '#F59E0B', // Amber
    '#EF4444', // Red
    '#8B5CF6', // Purple
    '#EC4899', // Pink
    '#06B6D4', // Cyan
    '#14B8A6', // Teal
  ]
  const charCode = initial.charCodeAt(0)
  const index = charCode % colors.length
  return colors[index]
})

const labelColor = computed(() => {
  if (props.node.type === 'category') return 'text-text-primary font-semibold'
  return 'text-text-secondary'
})


// Selection & double-click
function handleSelect() {
  emit('select-node', props.node.id)
}

function handleDoubleClick() {
  if (props.node.type === 'category') {
    toggleExpand()
  } else if (props.node.data) {
    const { schema, name, connectionId } = props.node.data
    if (props.node.type === 'table' || props.node.type === 'view') {
      emit('open-table', schema, name, connectionId)
    } else if (props.node.type === 'function') {
      emit('open-query', schema, name, connectionId)
    }
  }
}

// Drag & Drop
function handleDragStart(event: DragEvent, id: string) {
  event.dataTransfer?.setData('text/plain', id)
  if (props.node.type === 'table' || props.node.type === 'view') {
    const tableData = {
      id: props.node.id,
      label: props.node.label,
      type: props.node.type,
      schema: props.node.data?.schema || 'public',
      name: props.node.data?.name || props.node.label,
      connectionId: props.node.data?.connectionId || connectionsStore.currentConnectionId
    }
    event.dataTransfer?.setData('application/json', JSON.stringify(tableData))
  }
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'copyMove'
  }
}

function handleDragOver(event: DragEvent) {
  if (props.node.type === 'category') {
    event.preventDefault()
    event.stopPropagation()
    if (event.dataTransfer) {
      event.dataTransfer.dropEffect = 'move'
    }
  }
}

function handleDragEnter(event: DragEvent) {
  if (props.node.type === 'category') {
    event.preventDefault()
    event.stopPropagation()
    dragOver.value = true
  }
}

function handleDragLeave(event: DragEvent) {
  if (props.node.type === 'category') {
    dragOver.value = false
  }
}

// Wrapper-level handlers: allow dropping anywhere on a category (including its children area)
function handleWrapperDragOver(event: DragEvent) {
  if (props.node.type === 'category') {
    event.preventDefault()
    event.stopPropagation()
  }
}

function handleWrapperDrop(event: DragEvent) {
  if (props.node.type !== 'category') return
  event.preventDefault()
  event.stopPropagation()
  dragOver.value = false
  const nodeId = event.dataTransfer?.getData('text/plain')
  if (nodeId && nodeId !== props.node.id) {
    workspaceStore.moveNode(nodeId, props.node.id)
  }
}

function handleDrop(event: DragEvent, targetId: string) {
  if (props.node.type !== 'category') return
  event.preventDefault()
  event.stopPropagation()
  dragOver.value = false
  const nodeId = event.dataTransfer?.getData('text/plain')
  if (nodeId && nodeId !== targetId) {
    workspaceStore.moveNode(nodeId, targetId)
  }
}

const contextMenuItems = computed<ContextMenuItem[]>(() => {
  if (props.node.type === 'category') {
    return [
      { label: 'Rename', action: 'rename' },
    ]
  }
  if (props.node.type === 'table' || props.node.type === 'view') {
    return [
      { label: `Open ${props.node.type === 'table' ? 'Table' : 'View'}`, action: 'open' },
      { label: 'Generate SELECT (DML)', action: 'dml' },
      { label: 'Generate CREATE DDL', action: 'ddl' }
    ]
  }
  return []
})

function showContextMenu(event: MouseEvent) {
  if (props.node.type === 'table' || props.node.type === 'view' || props.node.type === 'category') {
    event.preventDefault()
    menuX.value = event.clientX
    menuY.value = event.clientY
    showMenu.value = true
  }
}

async function handleMenuSelect(action: string) {
  showMenu.value = false

  if (action === 'rename') {
    emit('rename-node', props.node.id, props.node.label)
    return
  }

  const connectionId = props.node.data?.connectionId || connectionsStore.currentConnectionId
  const schema = props.node.data?.schema
  const table = props.node.label
  
  if (!connectionId || !schema || !table) return

  if (action === 'open') {
    emit('open-table', schema, table, connectionId)
  } else if (action === 'dml') {
    try {
      const cols = await schemaStore.loadColumns(connectionId, schema, table)
      const colNames = cols.length > 0 ? cols.map(c => c.column_name).join(',\n  ') : '*'
      const sql = `SELECT\n  ${colNames}\nFROM ${schema}.${table};`
      
      tabsStore.createTab('query', {
        title: `SELECT ${table}`,
        connectionId,
        schema,
        table,
        sql,
      })
    } catch (e: any) {
      console.error('Failed to generate SELECT query:', e)
    }
  } else if (action === 'ddl') {
    try {
      const ddl = await schemaStore.getTableDDL(connectionId, schema, table)
      tabsStore.createTab('query', {
        title: `${table} DDL`,
        connectionId,
        schema,
        table,
        sql: ddl,
      })
    } catch (e: any) {
      console.error('Failed to generate DDL:', e)
    }
  }
}

// Simple icon components
function DatabaseIcon() {
  return h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2', class: 'w-4 h-4' }, [
    h('ellipse', { cx: '12', cy: '5', rx: '9', ry: '3' }),
    h('path', { d: 'M21 12c0 1.66-4 3-9 3s-9-1.34-9-3' }),
    h('path', { d: 'M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5' }),
  ])
}

function TableIcon() {
  return h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2', class: 'w-4 h-4' }, [
    h('rect', { x: '3', y: '3', width: '18', height: '18', rx: '2' }),
    h('path', { d: 'M3 9h18M3 15h18M9 3v18' }),
  ])
}

function ViewIcon() {
  return h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2', class: 'w-4 h-4' }, [
    h('path', { d: 'M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z' }),
    h('circle', { cx: '12', cy: '12', r: '3' }),
  ])
}

function FunctionIcon() {
  return h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2', class: 'w-4 h-4' }, [
    h('path', { d: 'M8 3H5a2 2 0 0 0-2 2v14c0 1.1.9 2 2 2h3' }),
    h('path', { d: 'M16 3h3a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-3' }),
    h('path', { d: 'M9 12h6M12 9v6' }),
  ])
}

function SequenceIcon() {
  return h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2', class: 'w-4 h-4' }, [
    h('path', { d: 'M4 12h16M4 12l4-4M4 12l4 4M20 12l-4-4M20 12l-4 4' }),
  ])
}

function ColumnIcon() {
  return h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2', class: 'w-4 h-4' }, [
    h('rect', { x: '8', y: '2', width: '8', height: '20', rx: '1' }),
  ])
}

function IndexIcon() {
  return h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2', class: 'w-4 h-4' }, [
    h('path', { d: 'M4 7h16M4 12h16M4 17h10' }),
  ])
}

function TypeIcon() {
  return h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2', class: 'w-4 h-4' }, [
    h('path', { d: 'M4 7V4h16v3M9 20h6M12 4v16' }),
  ])
}

function FolderIcon() {
  return h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2', class: 'w-4 h-4' }, [
    h('path', { d: 'M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z' }),
  ])
}
</script>
