<template>
  <div class="h-full w-full flex flex-col overflow-hidden bg-navy-primary p-4">
    <div class="flex-1 flex gap-4 min-h-0 overflow-hidden">
      
      <!-- Left Pane: Objects that depend on this table -->
      <div class="flex-1 flex flex-col min-h-0 bg-navy-secondary border border-navy-border rounded-lg overflow-hidden">
        <div class="px-4 py-2.5 bg-navy-tertiary border-b border-navy-border select-none flex-shrink-0">
          <h3 class="text-xs font-semibold text-text-primary truncate">
            Objects that depend on &lt;{{ props.tab.schema || 'public' }}.{{ props.tab.table }}&gt;
          </h3>
        </div>
        
        <div class="flex-1 overflow-y-auto p-4 relative">
          <div v-if="loading" class="absolute inset-0 flex items-center justify-center bg-navy-secondary/50">
            <div class="w-6 h-6 border-2 border-teal-accent border-t-transparent rounded-full animate-spin"></div>
          </div>
          
          <div v-else-if="dependentsTree.length === 0" class="h-full flex items-center justify-center text-xs text-text-muted select-none">
            No objects depend on this table.
          </div>
          
          <div v-else class="space-y-1">
            <div v-for="node in dependentsTree" :key="node.id" class="text-xs">
              <!-- Level 1: Schema/Tablespace -->
              <div 
                @click="toggleNode(node.id)"
                class="flex items-center gap-1.5 py-1.5 px-2 hover:bg-navy-hover rounded cursor-pointer select-none text-[11px] font-semibold text-text-primary"
              >
                <component :is="isExpanded(node.id) ? ChevronDown : ChevronRight" class="w-3.5 h-3.5 text-text-muted" />
                <Database v-if="node.id.startsWith('dependents-schema-')" class="w-3.5 h-3.5 text-teal-accent" />
                <Folder v-else class="w-3.5 h-3.5 text-accent-amber" />
                <span>{{ node.label }}</span>
                <span v-if="node.count !== undefined" class="text-[10px] text-text-muted font-normal ml-1">({{ node.count }})</span>
              </div>

              <!-- Level 2: Object Type Folder (Views, Sequences, etc.) -->
              <div v-show="isExpanded(node.id)" class="pl-4 border-l border-navy-border/60 ml-3.5 mt-0.5 space-y-0.5">
                <div v-for="subNode in node.children" :key="subNode.id">
                  <div 
                    @click="toggleNode(subNode.id)"
                    class="flex items-center gap-1.5 py-1 px-2 hover:bg-navy-hover rounded cursor-pointer select-none text-[11px] text-text-secondary"
                  >
                    <component :is="isExpanded(subNode.id) ? ChevronDown : ChevronRight" class="w-3 h-3 text-text-muted" />
                    <component :is="isExpanded(subNode.id) ? FolderOpen : Folder" class="w-3.5 h-3.5 text-accent-blue" />
                    <span>{{ subNode.label }}</span>
                    <span v-if="subNode.count !== undefined" class="text-[10px] text-text-muted font-normal ml-1">({{ subNode.count }})</span>
                  </div>

                  <!-- Level 3: Actual Object List -->
                  <div v-show="isExpanded(subNode.id)" class="pl-4 border-l border-navy-border/60 ml-3 mt-0.5 space-y-0.5">
                    <div 
                      v-for="item in subNode.children" 
                      :key="item.id"
                      class="flex items-center gap-2 py-1 px-2 hover:bg-navy-hover rounded select-text text-[11px] text-text-muted font-mono"
                    >
                      <component :is="getItemIcon(subNode.label)" class="w-3.5 h-3.5 flex-shrink-0 text-text-muted" />
                      <span class="hover:text-text-primary transition-colors cursor-text">{{ item.label }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Right Pane: Objects this table depends on -->
      <div class="flex-1 flex flex-col min-h-0 bg-navy-secondary border border-navy-border rounded-lg overflow-hidden">
        <div class="px-4 py-2.5 bg-navy-tertiary border-b border-navy-border select-none flex-shrink-0">
          <h3 class="text-xs font-semibold text-text-primary truncate">
            Objects that &lt;{{ props.tab.schema || 'public' }}.{{ props.tab.table }}&gt; depends on
          </h3>
        </div>
        
        <div class="flex-1 overflow-y-auto p-4 relative">
          <div v-if="loading" class="absolute inset-0 flex items-center justify-center bg-navy-secondary/50">
            <div class="w-6 h-6 border-2 border-teal-accent border-t-transparent rounded-full animate-spin"></div>
          </div>
          
          <div v-else-if="dependenciesTree.length === 0" class="h-full flex items-center justify-center text-xs text-text-muted select-none">
            This table does not depend on any objects.
          </div>
          
          <div class="space-y-1" v-else>
            <div v-for="node in dependenciesTree" :key="node.id" class="text-xs">
              <!-- Level 1: Schema/Tablespace -->
              <div 
                @click="toggleNode(node.id)"
                class="flex items-center gap-1.5 py-1.5 px-2 hover:bg-navy-hover rounded cursor-pointer select-none text-[11px] font-semibold text-text-primary"
              >
                <component :is="isExpanded(node.id) ? ChevronDown : ChevronRight" class="w-3.5 h-3.5 text-text-muted" />
                <Database v-if="node.id.startsWith('dependencies-schema-')" class="w-3.5 h-3.5 text-teal-accent" />
                <Folder v-else class="w-3.5 h-3.5 text-accent-amber" />
                <span>{{ node.label }}</span>
                <span v-if="node.count !== undefined" class="text-[10px] text-text-muted font-normal ml-1">({{ node.count }})</span>
              </div>

              <!-- Level 2: Object Type Folder (Views, Sequences, etc.) -->
              <div v-show="isExpanded(node.id)" class="pl-4 border-l border-navy-border/60 ml-3.5 mt-0.5 space-y-0.5">
                <div v-for="subNode in node.children" :key="subNode.id">
                  <div 
                    @click="toggleNode(subNode.id)"
                    class="flex items-center gap-1.5 py-1 px-2 hover:bg-navy-hover rounded cursor-pointer select-none text-[11px] text-text-secondary"
                  >
                    <component :is="isExpanded(subNode.id) ? ChevronDown : ChevronRight" class="w-3 h-3 text-text-muted" />
                    <component :is="isExpanded(subNode.id) ? FolderOpen : Folder" class="w-3.5 h-3.5 text-accent-blue" />
                    <span>{{ subNode.label }}</span>
                    <span v-if="subNode.count !== undefined" class="text-[10px] text-text-muted font-normal ml-1">({{ subNode.count }})</span>
                  </div>

                  <!-- Level 3: Actual Object List -->
                  <div v-show="isExpanded(subNode.id)" class="pl-4 border-l border-navy-border/60 ml-3 mt-0.5 space-y-0.5">
                    <div 
                      v-for="item in subNode.children" 
                      :key="item.id"
                      class="flex items-center gap-2 py-1 px-2 hover:bg-navy-hover rounded select-text text-[11px] text-text-muted font-mono"
                    >
                      <component :is="getItemIcon(subNode.label)" class="w-3.5 h-3.5 flex-shrink-0 text-text-muted" />
                      <span class="hover:text-text-primary transition-colors cursor-text">{{ item.label }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useConnectionsStore } from '../stores/connections'
import { useUiStore } from '../stores/ui'
import type { Tab } from '../types'
import * as App from '../../wailsjs/go/main/App'
import { 
  ChevronRight, 
  ChevronDown, 
  Folder, 
  FolderOpen, 
  Database, 
  Eye, 
  Binary, 
  Cpu, 
  Table, 
  Key, 
  FileText 
} from 'lucide-vue-next'

const props = defineProps<{
  tab: Tab
  activeTab: string
}>()

const connectionsStore = useConnectionsStore()
const uiStore = useUiStore()

const loading = ref(false)
const dependentsTree = ref<TreeNode[]>([])
const dependenciesTree = ref<TreeNode[]>([])
const expandedStates = ref<Record<string, boolean>>({})

interface TreeNode {
  id: string
  label: string
  count?: number
  isFolder: boolean
  children?: TreeNode[]
}

function toggleNode(nodeId: string) {
  if (expandedStates.value[nodeId] === undefined) {
    expandedStates.value[nodeId] = false
  } else {
    expandedStates.value[nodeId] = !expandedStates.value[nodeId]
  }
}

function isExpanded(nodeId: string): boolean {
  // Nodes are expanded by default (unless explicitly toggled off to false)
  return expandedStates.value[nodeId] !== false
}

function getItemIcon(label: string) {
  const name = label.toLowerCase()
  if (name.includes('views')) return Eye
  if (name.includes('sequences')) return Binary
  if (name.includes('functions')) return Cpu
  if (name.includes('tablespaces')) return Folder
  if (name.includes('tables')) return Table
  if (name.includes('indexes')) return Key
  return FileText
}

async function loadData() {
  const connId = props.tab.connectionId || connectionsStore.currentConnectionId
  if (!connId || !props.tab.table) return

  loading.value = true
  
  const schema = props.tab.schema || 'public'
  const table = props.tab.table
  const tableOid = `'"${schema}"."${table}"'::regclass`

  console.log('[DDLView_Depend] Loading dependencies for table:', `${schema}.${table}`, 'Oid cast:', tableOid)

  // 1. Query dependents (objects depending on our table)
  // Exclude Index relations ('i') to match pgAdmin's clean visual list
  const dependentsSql = `
    SELECT DISTINCT
        n.nspname AS schema_name,
        c.relname AS object_name,
        CASE 
            WHEN c.relkind = 'v' THEN 'Views'
            WHEN c.relkind = 'm' THEN 'Materialized Views'
            WHEN c.relkind = 'S' THEN 'Sequences'
            WHEN c.relkind = 'r' THEN 'Tables'
            ELSE 'Other'
        END AS object_type
    FROM pg_depend d
    JOIN pg_class c ON d.objid = c.oid
    LEFT JOIN pg_namespace n ON c.relnamespace = n.oid
    WHERE d.refobjid = ${tableOid} 
      AND d.objid <> ${tableOid}
      AND c.relkind IN ('v', 'm', 'S', 'r')

    UNION

    SELECT DISTINCT
        n.nspname AS schema_name,
        c.relname AS object_name,
        'Views' AS object_type
    FROM pg_depend d
    JOIN pg_rewrite r ON r.oid = d.objid
    JOIN pg_class c ON c.oid = r.ev_class
    JOIN pg_namespace n ON c.relnamespace = n.oid
    WHERE d.refobjid = ${tableOid} 
      AND d.classid = 'pg_rewrite'::regclass 
      AND c.oid <> ${tableOid}
      AND c.relkind = 'v'

    ORDER BY object_type, schema_name, object_name;
  `

  // 2. Query dependencies (what our table depends on)
  const dependenciesSql = `
    SELECT DISTINCT
        n.nspname AS schema_name,
        c.relname AS object_name,
        CASE 
            WHEN c.relkind = 'S' THEN 'Sequences'
            WHEN c.relkind = 'r' THEN 'Tables'
            ELSE 'Other'
        END AS object_type
    FROM pg_depend d
    JOIN pg_class c ON d.refobjid = c.oid
    LEFT JOIN pg_namespace n ON c.relnamespace = n.oid
    WHERE d.objid = ${tableOid} 
      AND d.refobjid <> ${tableOid}
      AND c.relkind IN ('S', 'r')

    UNION

    -- Sequences used by serial/identity default values (via pg_attrdef)
    SELECT DISTINCT
        n.nspname AS schema_name,
        c.relname AS object_name,
        'Sequences' AS object_type
    FROM pg_depend d
    JOIN pg_attrdef ad ON d.objid = ad.oid
    JOIN pg_class c ON d.refobjid = c.oid
    LEFT JOIN pg_namespace n ON c.relnamespace = n.oid
    WHERE ad.adrelid = ${tableOid}
      AND c.relkind = 'S'

    UNION

    SELECT DISTINCT
        n.nspname AS schema_name,
        p.proname || '()' AS object_name,
        'Functions' AS object_type
    FROM pg_depend d
    JOIN pg_proc p ON d.refobjid = p.oid
    JOIN pg_namespace n ON p.pronamespace = n.oid
    WHERE d.objid = ${tableOid} 
      AND d.refclassid = 'pg_proc'::regclass

    UNION

    SELECT 
        '' AS schema_name,
        t.spcname AS object_name,
        'Tablespaces' AS object_type
    FROM pg_class c
    JOIN pg_tablespace t ON c.reltablespace = t.oid
    WHERE c.oid = ${tableOid}

    UNION

    SELECT 
        '' AS schema_name,
        t.spcname AS object_name,
        'Tablespaces' AS object_type
    FROM pg_class c, pg_database d
    JOIN pg_tablespace t ON d.dattablespace = t.oid
    WHERE c.oid = ${tableOid} 
      AND c.reltablespace = 0 
      AND d.datname = current_database()

    ORDER BY object_type, schema_name, object_name;
  `

  try {
    const bindings = connectionsStore.getWailsBindings()
    
    console.log('[DDLView_Depend] Executing queries sequentially...')
    const depRes = await bindings.ExecuteQuery(connId, dependentsSql, 30)
    const reqRes = await bindings.ExecuteQuery(connId, dependenciesSql, 30)

    if (depRes.error) {
      console.error('[DDLView_Depend] Dependents query error:', depRes.error)
      throw new Error(depRes.error)
    }
    if (reqRes.error) {
      console.error('[DDLView_Depend] Dependencies query error:', reqRes.error)
      throw new Error(reqRes.error)
    }

    const parsedDependents = parseQueryResult(depRes)
    const parsedDependencies = parseQueryResult(reqRes)

    console.log('[DDLView_Depend] Dependents raw data:', parsedDependents)
    console.log('[DDLView_Depend] Dependencies raw data:', parsedDependencies)

    dependentsTree.value = buildTree(parsedDependents, 'dependents')
    dependenciesTree.value = buildTree(parsedDependencies, 'dependencies')
  } catch (err: any) {
    console.error('[DDLView_Depend] Load error details:', err)
    uiStore.addNotification({
      type: 'error',
      title: 'Load Dependencies Failed',
      message: err.message || String(err)
    })
  } finally {
    loading.value = false
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

function buildTree(items: any[], side: 'dependents' | 'dependencies'): TreeNode[] {
  const tree: TreeNode[] = []
  
  const schemas: { [schema: string]: { [type: string]: string[] } } = {}
  const rootTypes: { [type: string]: string[] } = {}

  items.forEach(item => {
    const { schema_name, object_name, object_type } = item
    if (!schema_name) {
      if (!rootTypes[object_type]) {
        rootTypes[object_type] = []
      }
      rootTypes[object_type].push(object_name)
    } else {
      if (!schemas[schema_name]) {
        schemas[schema_name] = {}
      }
      if (!schemas[schema_name][object_type]) {
        schemas[schema_name][object_type] = []
      }
      schemas[schema_name][object_type].push(object_name)
    }
  })

  // 1. Root Level items (like Tablespaces)
  Object.keys(rootTypes).sort().forEach(type => {
    const names = rootTypes[type].sort()
    tree.push({
      id: `${side}-root-type-${type}`,
      label: type,
      count: names.length,
      isFolder: true,
      children: names.map(name => ({
        id: `${side}-root-type-${type}-${name}`,
        label: name,
        isFolder: false
      }))
    })
  })

  // 2. Schema folders (like public) containing grouped item types
  Object.keys(schemas).sort().forEach(schema => {
    const types = schemas[schema]
    const schemaChildren: TreeNode[] = []

    Object.keys(types).sort().forEach(type => {
      const names = types[type].sort()
      schemaChildren.push({
        id: `${side}-schema-${schema}-${type}`,
        label: type,
        count: names.length,
        isFolder: true,
        children: names.map(name => ({
          id: `${side}-schema-${schema}-${type}-${name}`,
          label: name,
          isFolder: false
        }))
      })
    })

    tree.push({
      id: `${side}-schema-${schema}`,
      label: schema,
      isFolder: true,
      children: schemaChildren
    })
  })

  return tree
}

onMounted(() => {
  if (props.activeTab === 'dependencies') {
    loadData()
  }
})

watch(() => props.activeTab, (newVal) => {
  if (newVal === 'dependencies') {
    loadData()
  }
})

watch(() => props.tab, () => {
  if (props.activeTab === 'dependencies') {
    loadData()
  }
}, { deep: true })
</script>
