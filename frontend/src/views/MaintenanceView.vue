<template>
  <div class="h-full w-full flex flex-col bg-navy-primary overflow-hidden">
    <!-- Header -->
    <div class="px-6 py-4 border-b border-navy-border flex items-center justify-between flex-shrink-0">
      <div>
        <h2 class="text-sm font-semibold text-text-primary flex items-center gap-2">
          <svg class="w-4 h-4 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z" />
          </svg>
          Database Maintenance
        </h2>
        <p class="text-[10px] text-text-secondary mt-0.5">Optimize, reclaim storage, and rebuild indices for your PostgreSQL database</p>
      </div>

      <!-- Status Badge -->
      <div class="flex items-center gap-2">
        <div v-if="running" class="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-accent-blue/15 border border-accent-blue/30 text-xs text-accent-blue font-medium">
          <div class="w-1.5 h-1.5 rounded-full bg-accent-blue animate-pulse"></div>
          Running Maintenance...
        </div>
        <div v-else-if="status === 'success'" class="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-accent-green/15 border border-accent-green/30 text-xs text-accent-green font-medium">
          <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <polyline points="20 6 9 17 4 12" />
          </svg>
          Success
        </div>
        <div v-else-if="status === 'error'" class="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-accent-red/15 border border-accent-red/30 text-xs text-accent-red font-medium">
          <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" />
          </svg>
          Failed
        </div>
        <div v-else class="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-navy-secondary border border-navy-border text-xs text-text-secondary">
          <div class="w-1.5 h-1.5 rounded-full bg-text-muted"></div>
          Idle
        </div>
      </div>
    </div>

    <!-- Main Content (Split Pane) -->
    <div class="flex-1 flex overflow-hidden p-6 gap-6">
      <!-- Left Pane: Form Settings (Scrollable) -->
      <div class="w-[420px] flex flex-col gap-4 overflow-y-auto pr-2 flex-shrink-0">
        <!-- 1. Target Connection & Database -->
        <div class="p-4 bg-navy-secondary border border-navy-border rounded-lg space-y-3.5">
          <div class="border-b border-navy-border pb-1.5 flex items-center justify-between">
            <h3 class="text-xs font-semibold text-text-primary flex items-center gap-1.5">
              <svg class="w-3.5 h-3.5 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <ellipse cx="12" cy="5" rx="9" ry="3" />
                <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3" />
                <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5" />
              </svg>
              Database Target
            </h3>

            <!-- Run Maintenance Button -->
            <button @click="startMaintenance" type="button" 
              class="px-3 py-1 text-[10px] bg-transparent border border-teal-accent text-teal-accent font-semibold hover:bg-teal-accent hover:text-navy-primary rounded-md transition-all duration-300 flex items-center gap-1.5 shadow-[0_0_12px_rgba(0,201,167,0.35)] hover:shadow-[0_0_20px_rgba(0,201,167,0.65)] cursor-pointer disabled:opacity-50 disabled:pointer-events-none disabled:shadow-none" 
              :disabled="running || !form.connectionId || !form.database">
              <svg v-if="running" class="w-3.5 h-3.5 border-2 border-current border-t-transparent rounded-full animate-spin" viewBox="0 0 24 24"></svg>
              <svg v-else class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polygon points="5 3 19 12 5 21 5 3" />
              </svg>
              {{ running ? 'Executing...' : 'Run Operation' }}
            </button>
          </div>

          <!-- Connection Select -->
          <div>
            <label class="block text-[10px] text-text-secondary mb-1">Target Connection</label>
            <select v-model="form.connectionId" class="w-full bg-navy-tertiary border border-navy-border rounded-md px-2 py-1.5 text-xs text-text-primary focus:border-teal-accent focus:outline-none" :disabled="running">
              <option value="" disabled>Select a connection...</option>
              <option v-for="conn in connectionsStore.connections" :key="conn.id" :value="conn.id">
                {{ conn.name }} ({{ conn.host }}:{{ conn.port }})
              </option>
            </select>
          </div>

          <!-- Database Select -->
          <div>
            <label class="block text-[10px] text-text-secondary mb-1">Target Database</label>
            <div class="relative">
              <select v-model="form.database" class="w-full bg-navy-tertiary border border-navy-border rounded-md px-2 py-1.5 text-xs text-text-primary focus:border-teal-accent focus:outline-none" :disabled="running || loadingDbs || !form.connectionId">
                <option value="" disabled>{{ loadingDbs ? 'Loading databases...' : 'Select a database...' }}</option>
                <option v-for="db in dbs" :key="db" :value="db">{{ db }}</option>
              </select>
              <div v-if="loadingDbs" class="absolute right-2.5 top-1/2 -translate-y-1/2">
                <div class="w-3.5 h-3.5 border-2 border-teal-accent border-t-transparent rounded-full animate-spin"></div>
              </div>
            </div>
          </div>
        </div>

        <!-- 2. Maintenance Operation -->
        <div class="p-4 bg-navy-secondary border border-navy-border rounded-lg space-y-4">
          <h3 class="text-xs font-semibold text-text-primary border-b border-navy-border pb-1.5 flex items-center gap-1.5">
            <svg class="w-3.5 h-3.5 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="9" />
              <path d="M12 8v8" />
              <path d="M8 12h8" />
            </svg>
            Maintenance Settings
          </h3>

          <!-- Operation Type Selection -->
          <div>
            <label class="block text-[10px] text-text-secondary mb-1.5">Operation Type</label>
            <div class="grid grid-cols-1 gap-2">
              <label class="flex items-start gap-2.5 p-2 bg-navy-tertiary/40 border border-navy-border rounded-md hover:border-teal-accent/30 cursor-pointer transition-colors" :class="{ 'border-teal-accent bg-teal-accent/5': form.task === 'vacuum_analyze' }">
                <input type="radio" v-model="form.task" value="vacuum_analyze" class="mt-0.5 accent-teal-accent" :disabled="running" />
                <div>
                  <span class="text-xs font-medium text-text-primary block">Vacuum Analyze</span>
                  <span class="text-[9px] text-text-secondary block mt-0.5">Reclaims storage from dead tuples, updates database planner query execution statistics</span>
                </div>
              </label>

              <label class="flex items-start gap-2.5 p-2 bg-navy-tertiary/40 border border-navy-border rounded-md hover:border-teal-accent/30 cursor-pointer transition-colors" :class="{ 'border-teal-accent bg-teal-accent/5': form.task === 'vacuum_full' }">
                <input type="radio" v-model="form.task" value="vacuum_full" class="mt-0.5 accent-teal-accent" :disabled="running" />
                <div>
                  <span class="text-xs font-medium text-text-primary block">Vacuum Full</span>
                  <span class="text-[9px] text-text-secondary block mt-0.5">Reclaims maximum storage back to OS by compacting tables. <strong class="text-accent-amber">Requires exclusive table locks.</strong></span>
                </div>
              </label>

              <label class="flex items-start gap-2.5 p-2 bg-navy-tertiary/40 border border-navy-border rounded-md hover:border-teal-accent/30 cursor-pointer transition-colors" :class="{ 'border-teal-accent bg-teal-accent/5': form.task === 'reindex' }">
                <input type="radio" v-model="form.task" value="reindex" class="mt-0.5 accent-teal-accent" :disabled="running" />
                <div>
                  <span class="text-xs font-medium text-text-primary block">Reindex Database</span>
                  <span class="text-[9px] text-text-secondary block mt-0.5">Rebuilds degraded indexes. Essential for databases with heavy transactional updates</span>
                </div>
              </label>
            </div>
          </div>

          <!-- Scope Selection -->
          <div>
            <label class="block text-[10px] text-text-secondary mb-1">Target Scope</label>
            <select v-model="form.scope" class="w-full bg-navy-tertiary border border-navy-border rounded-md px-2 py-1.5 text-xs text-text-primary focus:border-teal-accent focus:outline-none" :disabled="running">
              <option value="database">Entire Database</option>
              <option v-if="form.task === 'reindex'" value="schema">Specific Schema</option>
              <option value="table">Specific Table</option>
            </select>
          </div>

          <!-- Target Schema Select (Visible when scope is schema or table) -->
          <div v-if="form.scope === 'schema' || form.scope === 'table'">
            <label class="block text-[10px] text-text-secondary mb-1">Schema</label>
            <select v-model="form.schemaName" class="w-full bg-navy-tertiary border border-navy-border rounded-md px-2 py-1.5 text-xs text-text-primary focus:border-teal-accent focus:outline-none" :disabled="running || !form.connectionId">
              <option value="" disabled>Select a schema...</option>
              <option v-for="sch in schemas" :key="sch.name" :value="sch.name">{{ sch.name }}</option>
            </select>
          </div>

          <!-- Target Table Select (Visible when scope is table) -->
          <div v-if="form.scope === 'table'">
            <label class="block text-[10px] text-text-secondary mb-1">Table</label>
            <select v-model="form.tableName" class="w-full bg-navy-tertiary border border-navy-border rounded-md px-2 py-1.5 text-xs text-text-primary focus:border-teal-accent focus:outline-none" :disabled="running || !form.connectionId || !form.schemaName">
              <option value="" disabled>Select a table...</option>
              <option v-for="tbl in tables" :key="tbl.name" :value="tbl.name">{{ tbl.name }}</option>
            </select>
          </div>

          <!-- Options -->
          <div class="pt-2 border-t border-navy-border">
            <label class="flex items-center gap-2 p-1.5 hover:bg-navy-hover rounded cursor-pointer transition-colors" :class="{ 'opacity-50 pointer-events-none': running }">
              <input v-model="form.verbose" type="checkbox" class="w-3.5 h-3.5 accent-teal-accent rounded bg-navy-tertiary border-navy-border" />
              <span class="text-[10px] text-text-primary">Verbose Console Output</span>
            </label>
          </div>
        </div>

        <!-- 3. SQL Command Preview -->
        <div class="p-4 bg-[#090d16] border border-navy-border rounded-lg space-y-2">
          <span class="text-[9px] uppercase font-mono text-text-muted tracking-wider block select-none">SQL Statement Preview</span>
          <pre class="text-[10px] font-mono text-teal-accent/95 bg-[#050911] p-3 rounded border border-navy-border/40 whitespace-pre-wrap overflow-x-auto select-all leading-normal">{{ sqlPreview || '-- Selection incomplete --' }}</pre>
        </div>
      </div>

      <!-- Right Pane: Terminal Log Output Console -->
      <div class="flex-1 flex flex-col bg-navy-secondary border border-navy-border rounded-lg overflow-hidden">
        <!-- Terminal Header -->
        <div class="px-4 py-2 border-b border-navy-border bg-[#090d16] flex items-center justify-between flex-shrink-0 select-none">
          <div class="flex items-center gap-2 text-xs text-text-secondary font-medium">
            <div class="w-2.5 h-2.5 rounded-full bg-accent-red"></div>
            <div class="w-2.5 h-2.5 rounded-full bg-accent-amber"></div>
            <div class="w-2.5 h-2.5 rounded-full bg-accent-green"></div>
            <span class="ml-2 font-mono text-[10px] tracking-wide">MAINTENANCE CONSOLE LOG</span>
          </div>

          <div class="flex items-center gap-2">
            <button @click="copyLogs" class="p-1 hover:bg-navy-hover rounded text-text-muted hover:text-text-primary transition-colors cursor-pointer" title="Copy to Clipboard">
              <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2" /><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
              </svg>
            </button>
            <button @click="clearLogs" class="p-1 hover:bg-navy-hover rounded text-text-muted hover:text-accent-red transition-colors cursor-pointer" title="Clear Console">
              <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6" /><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Terminal Logs Window -->
        <div ref="terminalBody" class="flex-1 overflow-y-auto p-4 bg-[#050911] font-mono text-[11px] leading-relaxed text-text-primary">
          <div v-if="logs.length === 0" class="h-full flex items-center justify-center text-text-muted text-xs select-none">
            Maintenance task status and warning notices will stream here in real-time...
          </div>
          <div v-else class="space-y-1">
            <div v-for="(log, idx) in logs" :key="idx" class="whitespace-pre-wrap select-text selection:bg-teal-accent/30 selection:text-white" :class="getLogClass(log)">
              {{ log }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useConnectionsStore } from '../stores/connections'
import { useUiStore } from '../stores/ui'
import * as App from '../../wailsjs/go/main/App'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'

const connectionsStore = useConnectionsStore()
const uiStore = useUiStore()

const running = ref(false)
const status = ref<'idle' | 'success' | 'error'>('idle')
const logs = ref<string[]>([])
const terminalBody = ref<HTMLDivElement | null>(null)
const dbs = ref<string[]>([])
const schemas = ref<any[]>([])
const tables = ref<any[]>([])
const loadingDbs = ref(false)

const form = reactive({
  connectionId: '',
  database: '',
  task: 'vacuum_analyze', // vacuum_analyze, vacuum_full, reindex
  scope: 'database', // database, schema, table
  schemaName: 'public',
  tableName: '',
  verbose: true,
})

// Dynamic SQL statement preview generator
const sqlPreview = computed(() => {
  if (!form.database) return ''
  switch (form.task) {
    case 'vacuum_analyze':
      if (form.scope === 'table' && form.schemaName && form.tableName) {
        return `VACUUM (ANALYZE${form.verbose ? ', VERBOSE' : ''}) "${form.schemaName}"."${form.tableName}";`
      }
      return `VACUUM (ANALYZE${form.verbose ? ', VERBOSE' : ''});`
    case 'vacuum_full':
      if (form.scope === 'table' && form.schemaName && form.tableName) {
        return `VACUUM (FULL, ANALYZE${form.verbose ? ', VERBOSE' : ''}) "${form.schemaName}"."${form.tableName}";`
      }
      return `VACUUM (FULL, ANALYZE${form.verbose ? ', VERBOSE' : ''});`
    case 'reindex':
      switch (form.scope) {
        case 'schema':
          return `REINDEX ${form.verbose ? '(VERBOSE) ' : ''}SCHEMA "${form.schemaName}";`
        case 'table':
          if (form.schemaName && form.tableName) {
            return `REINDEX ${form.verbose ? '(VERBOSE) ' : ''}TABLE "${form.schemaName}"."${form.tableName}";`
          }
          return ''
        default:
          return `REINDEX ${form.verbose ? '(VERBOSE) ' : ''}DATABASE "${form.database}";`
      }
    default:
      return ''
  }
})

// Initialize connections list
onMounted(() => {
  if (connectionsStore.connections.length === 0) {
    connectionsStore.loadConnections()
  }

  // Pre-fill active connection
  if (connectionsStore.currentConnectionId) {
    form.connectionId = connectionsStore.currentConnectionId
  }

  // Bind Wails real-time backend log streaming
  EventsOn('maintenance:log', (message: string) => {
    logs.value.push(message)
    scrollToBottom()
  })

  EventsOn('maintenance:status', (data: { status: 'success' | 'error'; message: string }) => {
    status.value = data.status
    running.value = false
    if (data.status === 'error') {
      uiStore.addNotification({
        type: 'error',
        title: 'Maintenance Failed',
        message: data.message,
      })
    } else {
      uiStore.addNotification({
        type: 'success',
        title: 'Maintenance Complete',
        message: 'Database optimization completed successfully.',
      })
    }
    scrollToBottom()
  })
})

onUnmounted(() => {
  EventsOff('maintenance:log')
  EventsOff('maintenance:status')
})

// When connectionId changes, fetch databases and schema info
watch(() => form.connectionId, async (newVal) => {
  form.database = ''
  form.tableName = ''
  dbs.value = []
  schemas.value = []
  tables.value = []
  
  if (!newVal) return

  loadingDbs.value = true
  try {
    const list = await App.GetDatabases(newVal)
    dbs.value = list || []
    
    // Auto select connection default database if available
    const conn = connectionsStore.connections.find(c => c.id === newVal)
    if (conn && dbs.value.includes(conn.database)) {
      form.database = conn.database
    } else if (dbs.value.length > 0) {
      form.database = dbs.value[0]
    }

    // Load schemas
    const schList = await App.GetSchemas(newVal)
    schemas.value = schList || []
    if (schemas.value.length > 0) {
      const hasPublic = schemas.value.some(s => s.name === 'public')
      form.schemaName = hasPublic ? 'public' : schemas.value[0].name
    }
  } catch (e) {
    console.error('Failed to load connection data details:', e)
  } finally {
    loadingDbs.value = false
  }
}, { immediate: true })

// When schemaName changes, fetch tables list
watch(() => [form.connectionId, form.schemaName], async ([connId, schemaName]) => {
  form.tableName = ''
  tables.value = []
  if (!connId || !schemaName) return

  try {
    const tblList = await App.GetTables(connId, schemaName)
    tables.value = tblList || []
    if (tables.value.length > 0) {
      form.tableName = tables.value[0].name
    }
  } catch (e) {
    console.error('Failed to load table details list:', e)
  }
})

// Reset scope to database/table for non-reindex operations
watch(() => form.task, (newTask) => {
  if (newTask !== 'reindex' && form.scope === 'schema') {
    form.scope = 'database'
  }
})

// Execute Maintenance Task
async function startMaintenance() {
  if (running.value) return

  logs.value = []
  status.value = 'idle'
  running.value = true

  try {
    const opts = {
      connection_id: form.connectionId,
      database: form.database,
      task: form.task,
      scope: form.scope,
      schema_name: form.schemaName,
      table_name: form.tableName,
      verbose: form.verbose,
    }

    await App.StartMaintenance(opts)
  } catch (e: any) {
    running.value = false
    status.value = 'error'
    const msg = e.message || String(e)
    logs.value.push(`Fatal Error: ${msg}`)
    uiStore.addNotification({
      type: 'error',
      title: 'Execution Error',
      message: msg,
    })
  }
}

// Console utilities
function clearLogs() {
  logs.value = []
}

function copyLogs() {
  if (logs.value.length === 0) return
  navigator.clipboard.writeText(logs.value.join('\n'))
  uiStore.addNotification({
    type: 'info',
    title: 'Copied',
    message: 'Maintenance console logs copied to clipboard.',
  })
}

function scrollToBottom() {
  nextTick(() => {
    if (terminalBody.value) {
      terminalBody.value.scrollTop = terminalBody.value.scrollHeight
    }
  })
}

// Log line semantic styling matching terminal color states
function getLogClass(line: string) {
  const lower = line.toLowerCase()
  if (lower.includes('failed') || lower.includes('error') || lower.includes('fatal')) {
    return 'text-accent-red font-semibold'
  }
  if (lower.includes('success') || lower.includes('completed successfully')) {
    return 'text-accent-green font-semibold'
  }
  if (lower.includes('warning') || lower.includes('warning:')) {
    return 'text-accent-amber'
  }
  if (lower.startsWith('executing query:') || lower.includes('starting database maintenance')) {
    return 'text-accent-blue opacity-85'
  }
  // PostgreSQL notice levels
  if (lower.includes('[info]') || lower.startsWith('[info]')) {
    return 'text-text-primary opacity-80'
  }
  if (lower.includes('[notice]') || lower.startsWith('[notice]')) {
    return 'text-teal-accent/80'
  }
  return 'text-text-primary/95'
}
</script>

<style scoped>
/* Scrollbar */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}
::-webkit-scrollbar-track {
  background: #050911;
}
::-webkit-scrollbar-thumb {
  background: #1e293b;
  border-radius: 3px;
}
::-webkit-scrollbar-thumb:hover {
  background: #334155;
}
</style>
