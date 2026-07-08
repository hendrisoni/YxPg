<template>
  <Modal :show="true" :title="connection ? 'Edit Connection' : 'New Connection'" @close="$emit('close')">
    <form @submit.prevent="handleSave" class="space-y-4">
      <!-- Name -->
      <div>
        <label class="block text-xs text-text-secondary mb-1">Connection Name</label>
        <input
          v-model="form.name"
          type="text"
          placeholder="My Database"
          class="w-full"
          required
        />
      </div>

      <!-- Host & Port -->
      <div class="grid grid-cols-3 gap-3">
        <div class="col-span-2">
          <label class="block text-xs text-text-secondary mb-1">Host</label>
          <input
            v-model="form.host"
            type="text"
            placeholder="localhost"
            class="w-full"
            required
          />
        </div>
        <div>
          <label class="block text-xs text-text-secondary mb-1">Port</label>
          <input
            v-model.number="form.port"
            type="number"
            placeholder="5432"
            class="w-full"
            required
          />
        </div>
      </div>

      <!-- Database -->
      <div>
        <label class="block text-xs text-text-secondary mb-1">Database</label>
        <input
          v-model="form.database"
          type="text"
          placeholder="postgres"
          class="w-full"
          required
        />
      </div>

      <!-- Username & Password -->
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="block text-xs text-text-secondary mb-1">Username</label>
          <input
            v-model="form.username"
            type="text"
            placeholder="postgres"
            class="w-full"
            required
          />
        </div>
        <div>
          <label class="block text-xs text-text-secondary mb-1">Password</label>
          <input
            v-model="form.password"
            type="password"
            placeholder="••••••"
            class="w-full"
          />
        </div>
      </div>

      <!-- SSL Mode -->
      <div>
        <label class="block text-xs text-text-secondary mb-1">SSL Mode</label>
        <select v-model="form.ssl_mode" class="w-full">
          <option value="disable">Disable</option>
          <option value="require">Require</option>
          <option value="verify-full">Verify Full</option>
        </select>
      </div>

      <!-- Color -->
      <div>
        <label class="block text-xs text-text-secondary mb-1">Color Label</label>
        <div class="flex gap-2">
          <button
            v-for="color in colorOptions"
            :key="color.value"
            type="button"
            @click="form.color = color.value"
            class="w-6 h-6 rounded-full border-2 transition-all"
            :class="form.color === color.value ? 'border-white scale-110' : 'border-transparent'"
            :style="{ backgroundColor: color.value }"
          ></button>
        </div>
      </div>

      <!-- Import from URL -->
      <div>
        <label class="block text-xs text-text-secondary mb-1">Or import from URL</label>
        <input
          v-model="connectionUrl"
          type="text"
          placeholder="postgres://user:pass@host:5432/dbname"
          class="w-full text-xs"
          @blur="parseUrl"
        />
      </div>

      <!-- Test result -->
      <div v-if="testResult" class="flex items-center gap-2 text-xs">
        <div
          class="w-2 h-2 rounded-full"
          :class="testResult.ok ? 'bg-accent-green' : 'bg-accent-red'"
        ></div>
        <span :class="testResult.ok ? 'text-accent-green' : 'text-accent-red'">
          {{ testResult.ok ? `Connected (${testResult.latency_ms}ms)` : testResult.message }}
        </span>
      </div>

      <!-- Actions -->
      <div class="flex items-center justify-between pt-2">
        <button
          type="button"
          @click="handleTest"
          :disabled="testing"
          class="flex items-center gap-1.5 px-3 py-1.5 text-xs border border-navy-border rounded-md text-text-secondary hover:bg-navy-hover transition-colors disabled:opacity-50"
        >
          <div v-if="testing" class="w-3 h-3 border border-teal-accent border-t-transparent rounded-full animate-spin"></div>
          <svg v-else class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
            <path d="M22 4 12 14.01l-3-3" />
          </svg>
          Test Connection
        </button>

        <div class="flex gap-2">
          <button
            type="button"
            @click="$emit('close')"
            class="px-3 py-1.5 text-xs text-text-secondary hover:bg-navy-hover rounded-md transition-colors"
          >
            Cancel
          </button>
          <button
            type="submit"
            class="px-4 py-1.5 text-xs bg-teal-accent text-navy-primary rounded-md font-medium hover:bg-teal-hover transition-colors"
          >
            Save
          </button>
        </div>
      </div>
    </form>
  </Modal>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useConnectionsStore } from '../../stores/connections'
import type { Connection, ConnectionTestResult } from '../../types'
import Modal from '../shared/Modal.vue'
import * as App from '../../../wailsjs/go/main/App'

const props = defineProps<{
  connection?: Connection | null
}>()

const emit = defineEmits(['close', 'save'])

const connectionsStore = useConnectionsStore()

const colorOptions = [
  { value: '#00C9A7', label: 'Teal' },
  { value: '#3B82F6', label: 'Blue' },
  { value: '#F59E0B', label: 'Amber' },
  { value: '#EF4444', label: 'Red' },
  { value: '#10B981', label: 'Green' },
  { value: '#8B5CF6', label: 'Purple' },
  { value: '#EC4899', label: 'Pink' },
  { value: '#06B6D4', label: 'Cyan' },
]

const form = reactive({
  name: props.connection?.name || '',
  host: props.connection?.host || 'localhost',
  port: props.connection?.port || 5432,
  database: props.connection?.database || 'postgres',
  username: props.connection?.username || 'postgres',
  password: props.connection?.password || '',
  ssl_mode: props.connection?.ssl_mode || 'disable',
  color: props.connection?.color || colorOptions[Math.floor(Math.random() * colorOptions.length)].value,
})

onMounted(async () => {
  if (!props.connection) {
    try {
      const defaults = await App.GetDefaultConnectionConfig()
      if (defaults) {
        if (defaults.host) form.host = defaults.host
        if (defaults.port) form.port = defaults.port
        if (defaults.database) form.database = defaults.database
        if (defaults.username) form.username = defaults.username
        if (defaults.password !== undefined) form.password = defaults.password
      }
    } catch (err) {
      console.error('Failed to load default connection config:', err)
    }
  }
})

const connectionUrl = ref('')
const testResult = ref<ConnectionTestResult | null>(null)
const testing = ref(false)

function parseUrl() {
  if (!connectionUrl.value) return
  try {
    const url = new URL(connectionUrl.value)
    form.host = url.hostname || 'localhost'
    form.port = parseInt(url.port) || 5432
    form.database = url.pathname.replace('/', '') || 'postgres'
    form.username = url.username || 'postgres'
    form.password = url.password || ''
  } catch {
    // Invalid URL, ignore
  }
}

async function handleTest() {
  testing.value = true
  testResult.value = null
  try {
    testResult.value = await connectionsStore.testConnection({ ...props.connection, ...form } as Connection)
  } catch (e: any) {
    testResult.value = { ok: false, latency_ms: 0, message: e.message }
  } finally {
    testing.value = false
  }
}

function handleSave() {
  const conn: Connection = {
    id: props.connection?.id || '',
    ...form,
    created_at: props.connection?.created_at || new Date().toISOString(),
  }
  emit('save', conn)
}
</script>
