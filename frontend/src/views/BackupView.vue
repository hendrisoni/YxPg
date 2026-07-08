<template>
  <div class="h-full w-full flex flex-col bg-navy-primary overflow-hidden">
    <!-- Header -->
    <div class="px-6 py-4 border-b border-navy-border flex items-center justify-between flex-shrink-0">
      <div>
        <h2 class="text-sm font-semibold text-text-primary">Database Backup</h2>
        <p class="text-[10px] text-text-secondary mt-0.5">Export database schema and data using pg_dump</p>
      </div>
      
      <!-- Status Badge -->
      <div class="flex items-center gap-2">
        <div v-if="running" class="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-accent-blue/15 border border-accent-blue/30 text-xs text-accent-blue font-medium">
          <div class="w-1.5 h-1.5 rounded-full bg-accent-blue animate-pulse"></div>
          Running Backup...
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
        <!-- 1. Connection Setting -->
        <div class="p-4 bg-navy-secondary border border-navy-border rounded-lg space-y-3.5">
          <div class="border-b border-navy-border pb-1.5 flex items-center justify-between">
            <h3 class="text-xs font-semibold text-text-primary flex items-center gap-1.5">
              <svg class="w-3.5 h-3.5 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <ellipse cx="12" cy="5" rx="9" ry="3" />
                <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3" />
                <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5" />
              </svg>
              Database Connection
            </h3>
            
            <!-- Neon Run Backup Button -->
            <button @click="startBackup" type="button" class="px-3 py-1 text-[10px] bg-transparent border border-teal-accent text-teal-accent font-semibold hover:bg-teal-accent hover:text-navy-primary rounded-md transition-all duration-300 flex items-center gap-1.5 shadow-[0_0_12px_rgba(0,201,167,0.35)] hover:shadow-[0_0_20px_rgba(0,201,167,0.65)] cursor-pointer disabled:opacity-50 disabled:pointer-events-none disabled:shadow-none" :disabled="running || !form.connectionId || !form.outputPath">
              <svg v-if="running" class="w-3.5 h-3.5 border-2 border-current border-t-transparent rounded-full animate-spin" viewBox="0 0 24 24"></svg>
              <svg v-else class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                <polyline points="17 8 12 3 7 8" />
                <line x1="12" y1="3" x2="12" y2="15" />
              </svg>
              {{ running ? 'Backing Up...' : 'Run Backup' }}
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
        </div>

        <!-- 2. Export Destination & Format -->
        <div class="p-4 bg-navy-secondary border border-navy-border rounded-lg space-y-3.5">
          <h3 class="text-xs font-semibold text-text-primary border-b border-navy-border pb-1.5 flex items-center gap-1.5">
            <svg class="w-3.5 h-3.5 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v13a2 2 0 0 1-2 2z" />
              <polyline points="17 21 17 13 7 13 7 21" />
              <polyline points="7 3 7 8 15 8" />
            </svg>
            Format & Destination
          </h3>

          <!-- Format Select -->
          <div>
            <label class="block text-[10px] text-text-secondary mb-1">Export Format</label>
            <select v-model="form.format" class="w-full bg-navy-tertiary border border-navy-border rounded-md px-2 py-1.5 text-xs text-text-primary focus:border-teal-accent focus:outline-none" :disabled="running">
              <option value="p">Plain SQL Script (.sql)</option>
              <option value="c">Custom Compressed Archive (.backup)</option>
              <option value="d">Directory Format Dump (folder)</option>
              <option value="t">Tar Archive (.tar)</option>
            </select>
          </div>

          <!-- Output Path -->
          <div>
            <label class="block text-[10px] text-text-secondary mb-1">Output Path</label>
            <div class="flex gap-2">
              <input v-model="form.outputPath" type="text" placeholder="Select destination file or folder..." class="flex-1 bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-xs text-text-primary focus:border-teal-accent focus:outline-none" :disabled="running" required />
              <button @click="handleBrowse" type="button" class="px-3 py-1.5 text-xs bg-navy-tertiary border border-navy-border rounded-md text-text-primary hover:bg-navy-hover hover:border-teal-accent/50 transition-colors cursor-pointer flex-shrink-0" :disabled="running">
                Browse...
              </button>
            </div>
            <span class="text-[9px] text-text-muted mt-1 block">
              {{ form.format === 'd' ? 'Choose an empty folder for directory format' : 'Choose a file destination' }}
            </span>
          </div>
        </div>

        <!-- 3. Advanced Settings -->
        <div class="p-4 bg-navy-secondary border border-navy-border rounded-lg space-y-3">
          <h3 class="text-xs font-semibold text-text-primary border-b border-navy-border pb-1.5 flex items-center gap-1.5">
            <svg class="w-3.5 h-3.5 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="3" />
              <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z" />
            </svg>
            Advanced Dump Options
          </h3>

          <div class="grid grid-cols-2 gap-2 text-xs">
            <label class="flex items-center gap-2 p-1.5 hover:bg-navy-hover rounded cursor-pointer transition-colors" :class="{ 'opacity-50 pointer-events-none': running }">
              <input v-model="form.schemaOnly" type="checkbox" class="w-3.5 h-3.5 accent-teal-accent rounded bg-navy-tertiary border-navy-border" />
              <span class="text-[10px] text-text-primary">Schema Only (-s)</span>
            </label>

            <label class="flex items-center gap-2 p-1.5 hover:bg-navy-hover rounded cursor-pointer transition-colors" :class="{ 'opacity-50 pointer-events-none': running }">
              <input v-model="form.dataOnly" type="checkbox" class="w-3.5 h-3.5 accent-teal-accent rounded bg-navy-tertiary border-navy-border" />
              <span class="text-[10px] text-text-primary">Data Only (-a)</span>
            </label>

            <label class="flex items-center gap-2 p-1.5 hover:bg-navy-hover rounded cursor-pointer transition-colors" :class="{ 'opacity-50 pointer-events-none': running }">
              <input v-model="form.clean" type="checkbox" class="w-3.5 h-3.5 accent-teal-accent rounded bg-navy-tertiary border-navy-border" />
              <span class="text-[10px] text-text-primary">Clean Objects (-c)</span>
            </label>

            <label class="flex items-center gap-2 p-1.5 hover:bg-navy-hover rounded cursor-pointer transition-colors" :class="{ 'opacity-50 pointer-events-none': running }">
              <input v-model="form.create" type="checkbox" class="w-3.5 h-3.5 accent-teal-accent rounded bg-navy-tertiary border-navy-border" />
              <span class="text-[10px] text-text-primary">Create DB (-C)</span>
            </label>

            <label class="flex items-center gap-2 p-1.5 hover:bg-navy-hover rounded cursor-pointer transition-colors" :class="{ 'opacity-50 pointer-events-none': running }">
              <input v-model="form.inserts" type="checkbox" class="w-3.5 h-3.5 accent-teal-accent rounded bg-navy-tertiary border-navy-border" />
              <span class="text-[10px] text-text-primary">Inserts instead of COPY</span>
            </label>

            <label class="flex items-center gap-2 p-1.5 hover:bg-navy-hover rounded cursor-pointer transition-colors" :class="{ 'opacity-50 pointer-events-none': running }">
              <input v-model="form.columnInserts" type="checkbox" class="w-3.5 h-3.5 accent-teal-accent rounded bg-navy-tertiary border-navy-border" />
              <span class="text-[10px] text-text-primary">Column Inserts</span>
            </label>

            <label class="flex items-center gap-2 p-1.5 hover:bg-navy-hover rounded cursor-pointer transition-colors" :class="{ 'opacity-50 pointer-events-none': running || !form.dataOnly }">
              <input v-model="form.disableTriggers" type="checkbox" class="w-3.5 h-3.5 accent-teal-accent rounded bg-navy-tertiary border-navy-border" :disabled="!form.dataOnly" />
              <span class="text-[10px] text-text-primary">Disable Triggers</span>
            </label>

            <label class="flex items-center gap-2 p-1.5 hover:bg-navy-hover rounded cursor-pointer transition-colors" :class="{ 'opacity-50 pointer-events-none': running }">
              <input v-model="form.verbose" type="checkbox" class="w-3.5 h-3.5 accent-teal-accent rounded bg-navy-tertiary border-navy-border" />
              <span class="text-[10px] text-text-primary">Verbose Output (-v)</span>
            </label>
          </div>
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
            <span class="ml-2 font-mono text-[10px] tracking-wide">CONSOLE LOG</span>
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
            Backup logs will stream here in real-time...
          </div>
          <div v-else class="space-y-1">
            <div v-for="(log, idx) in logs" :key="idx" class="whitespace-pre-wrap select-text selection:bg-teal-accent/30 selection:text-white" :class="getLogClass(log)">
              <template v-if="log.startsWith('__ACTION_OPEN_FOLDER__')">
                <button @click="openBackupFolder" class="inline-flex items-center gap-1.5 px-3 py-1.5 mt-2 text-[10px] font-semibold bg-teal-accent/10 border border-teal-accent/30 hover:border-teal-accent/80 hover:bg-teal-accent/20 text-teal-accent rounded transition-all cursor-pointer focus:outline-none shadow-[0_0_8px_rgba(0,201,167,0.1)]">
                  <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                    <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
                  </svg>
                  Open Backup Folder
                </button>
              </template>
              <template v-else>
                {{ log }}
              </template>
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
const errorMsg = ref('')
const logs = ref<string[]>([])
const terminalBody = ref<HTMLDivElement | null>(null)
const dbOptions = ref<string[]>([])
const loadingDbs = ref(false)

const form = reactive({
  connectionId: '',
  database: '',
  outputPath: '',
  format: 'c', // p = plain sql, c = custom compressed, d = directory, t = tar (c is default)
  schemaOnly: false,
  dataOnly: false,
  clean: false,
  create: false,
  inserts: false,
  columnInserts: false,
  disableTriggers: false,
  verbose: true,
})

// Initialize form
onMounted(() => {
  if (connectionsStore.connections.length === 0) {
    connectionsStore.loadConnections()
  }
  
  // Default to current connection if active
  if (connectionsStore.currentConnectionId) {
    form.connectionId = connectionsStore.currentConnectionId
  }

  // Hook into Wails events
  EventsOn('backup:log', (message: string) => {
    logs.value.push(message)
    scrollToBottom()
  })

  EventsOn('backup:status', (data: { status: 'success' | 'error'; message: string }) => {
    status.value = data.status
    running.value = false
    if (data.status === 'error') {
      errorMsg.value = data.message
      uiStore.addNotification({
        type: 'error',
        title: 'Backup Failed',
        message: data.message,
      })
    } else {
      uiStore.addNotification({
        type: 'success',
        title: 'Backup Complete',
        message: 'Database backup finished successfully.',
      })
      logs.value.push('__ACTION_OPEN_FOLDER__')
    }
    scrollToBottom()
  })
})

onUnmounted(() => {
  EventsOff('backup:log')
  EventsOff('backup:status')
})

// Watch format and clear disableTriggers if dataOnly is unchecked
watch(() => form.dataOnly, (newVal) => {
  if (!newVal) {
    form.disableTriggers = false
  }
})

// Watch connection choice and load database options
watch(() => form.connectionId, (newVal) => {
  form.database = ''
  dbOptions.value = []
  
  if (!newVal) return

  const conn = connectionsStore.connections.find(c => c.id === newVal)
  if (conn) {
    form.database = conn.database
  }
}, { immediate: true })

// Watch database and format choices to update default output path
watch(() => [form.database, form.format], ([newDb, newFormat], [oldDb, oldFormat]) => {
  const getAutoPath = (db: string, fmt: string) => {
    if (!db) return ''
    const now = new Date()
    const yy = String(now.getFullYear()).slice(-2)
    const mm = String(now.getMonth() + 1).padStart(2, '0')
    const dd = String(now.getDate()).padStart(2, '0')
    const yymmdd = `${yy}${mm}${dd}`

    const extMap: Record<string, string> = {
      p: 'sql',
      c: 'backup',
      t: 'tar'
    }
    if (fmt === 'd') {
      return `D:\\${db}_${yymmdd}`
    }
    const ext = extMap[fmt] || 'backup'
    return `D:\\${db}_${yymmdd}.${ext}`
  }

  const expectedOldPath = oldDb ? getAutoPath(oldDb as string, oldFormat as string) : ''
  
  if (!form.outputPath || form.outputPath === expectedOldPath) {
    form.outputPath = getAutoPath(newDb as string, newFormat as string)
  }
})

// Log lines styling
function getLogClass(line: string) {
  if (line.startsWith('__ACTION_OPEN_FOLDER__')) {
    return ''
  }
	const lower = line.toLowerCase()
	if (lower.includes('error') || lower.includes('failed') || lower.includes('fatal')) {
		return 'text-accent-red font-semibold'
	}
	if (lower.includes('success') || lower.includes('complete')) {
		return 'text-accent-green font-semibold'
	}
	if (lower.includes('warning')) {
		return 'text-accent-amber'
	}
	if (lower.includes('starting') || lower.includes('running command')) {
		return 'text-accent-blue opacity-85'
	}
	return 'text-text-primary/90'
}

// Reset Form fields
function resetForm() {
  form.database = ''
  form.outputPath = ''
  form.format = 'c'
  form.schemaOnly = false
  form.dataOnly = false
  form.clean = false
  form.create = false
  form.inserts = false
  form.columnInserts = false
  form.disableTriggers = false
  form.verbose = true
  
  if (connectionsStore.currentConnectionId) {
    form.connectionId = connectionsStore.currentConnectionId
  } else {
    form.connectionId = ''
  }
  status.value = 'idle'
  errorMsg.value = ''
}

// File Dialogue browse helper
async function handleBrowse() {
  try {
    const extMap: Record<string, string> = {
      p: 'sql',
      c: 'backup',
      t: 'tar'
    }

    if (form.format === 'd') {
      const folderPath = await App.BrowseBackupFolder()
      if (folderPath) {
        form.outputPath = folderPath
      }
    } else {
      const extension = extMap[form.format] || 'sql'
      const defaultFilename = `${form.database || 'backup'}_${new Date().toISOString().slice(0, 10)}.${extension}`
      
      const filePath = await App.BrowseBackupFile(defaultFilename)
      if (filePath) {
        form.outputPath = filePath
      }
    }
  } catch (e: any) {
    uiStore.addNotification({
      type: 'error',
      title: 'Browse Error',
      message: e.message || String(e)
    })
  }
}

// Start pg_dump backup execution
async function startBackup() {
  if (running.value) return
  
  logs.value = []
  status.value = 'idle'
  errorMsg.value = ''
  running.value = true

  try {
    // Send standard fields + global pgBinPath setting
    const opts = {
      connection_id: form.connectionId,
      database: form.database,
      output_path: form.outputPath,
      format: form.format,
      schema_only: form.schemaOnly,
      data_only: form.dataOnly,
      clean: form.clean,
      create: form.create,
      inserts: form.inserts,
      column_inserts: form.columnInserts,
      disable_triggers: form.disableTriggers,
      verbose: form.verbose,
      pg_bin_path: uiStore.settings.pgBinPath
    }

    await App.StartBackup(opts)
  } catch (e: any) {
    running.value = false
    status.value = 'error'
    const msg = e.message || String(e)
    logs.value.push(`Error executing: ${msg}`)
    uiStore.addNotification({
      type: 'error',
      title: 'Execution Error',
      message: msg
    })
  }
}

// Clear Terminal console
function clearLogs() {
  logs.value = []
}

// Copy Terminal console to Clipboard
function copyLogs() {
  if (logs.value.length === 0) return
  const cleanLogs = logs.value.filter(line => !line.startsWith('__ACTION_OPEN_FOLDER__'))
  navigator.clipboard.writeText(cleanLogs.join('\n'))
  uiStore.addNotification({
    type: 'info',
    title: 'Copied',
    message: 'Backup logs copied to clipboard.'
  })
}

// Open backup directory in system file explorer
async function openBackupFolder() {
  if (!form.outputPath) return
  try {
    await App.OpenFolder(form.outputPath)
  } catch (e: any) {
    uiStore.addNotification({
      type: 'error',
      title: 'Failed to Open Folder',
      message: e.message || String(e)
    })
  }
}

// Scroll log body container down
function scrollToBottom() {
  nextTick(() => {
    if (terminalBody.value) {
      terminalBody.value.scrollTop = terminalBody.value.scrollHeight
    }
  })
}
</script>

<style scoped>
/* Custom styled thin terminal scrollbar */
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
