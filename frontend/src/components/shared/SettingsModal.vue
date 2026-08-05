<template>
  <Modal :show="true" title="Settings" @close="$emit('close')" size="md">
    <div class="space-y-4 py-2 select-none">
      <!-- Checkbox Settings -->
      <div class="flex flex-col gap-3">
        <label class="flex items-start gap-3 cursor-pointer p-2.5 rounded-lg bg-navy-secondary border border-navy-border hover:border-teal-accent/30 hover:bg-navy-hover transition-all group">
          <input
            v-model="showFunctionsInSearch"
            type="checkbox"
            class="mt-0.5 w-4 h-4 rounded text-teal-accent bg-navy-tertiary border-navy-border focus:ring-teal-accent focus:ring-offset-navy-secondary focus:ring-2 cursor-pointer accent-teal-accent"
          />
          <div class="flex flex-col">
            <span class="text-xs font-medium text-text-primary group-hover:text-teal-accent transition-colors">Show functions in search table</span>
            <span class="text-[10px] text-text-muted mt-0.5">Include database functions in Ctrl+K search results</span>
          </div>
        </label>
      </div>

      <!-- SQLite Database Location Setting -->
      <div class="flex flex-col gap-1.5 p-3 rounded-lg bg-navy-secondary border border-navy-border">
        <div class="flex items-center justify-between">
          <span class="text-xs font-semibold text-text-primary">Lokasi Database SQLite (.sys)</span>
          <span class="text-[10px] text-teal-accent font-medium">System Storage</span>
        </div>
        <div class="flex items-center gap-2 mt-1">
          <input
            v-model="dbPath"
            type="text"
            placeholder="D:\Yx\YxPg\yxpg.sys"
            class="flex-1 text-xs bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-text-primary focus:border-teal-accent focus:outline-none"
          />
          <button
            @click="handleBrowseDb"
            type="button"
            class="px-2.5 py-1.5 text-xs font-medium bg-navy-tertiary hover:bg-navy-hover border border-navy-border text-text-primary rounded-md transition-colors flex items-center gap-1.5 cursor-pointer shrink-0"
            title="Pilih lokasi file database SQLite (.sys)"
          >
            <svg class="w-3.5 h-3.5 text-teal-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
            </svg>
            Ambil Alamat DB
          </button>
          <button
            @click="handleRefreshDb"
            type="button"
            :disabled="isRefreshing"
            class="px-2.5 py-1.5 text-xs font-medium bg-teal-accent/10 hover:bg-teal-accent/20 border border-teal-accent/40 text-teal-accent rounded-md transition-colors flex items-center gap-1.5 cursor-pointer shrink-0 disabled:opacity-50"
            title="Refresh database SQLite & muat ulang data"
          >
            <svg class="w-3.5 h-3.5" :class="{ 'animate-spin': isRefreshing }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Refresh DB
          </button>
        </div>
        <span class="text-[10px] text-text-muted mt-1 leading-normal">
          Database SQLite menyimpan koneksi, history query, & workspace (default: <code>D:\Yx\YxPg\yxpg.sys</code>).
        </span>
      </div>

      <!-- PostgreSQL Bin Path Setting -->
      <div class="flex flex-col gap-1.5 p-3 rounded-lg bg-navy-secondary border border-navy-border">
        <div class="flex items-center justify-between">
          <span class="text-xs font-semibold text-text-primary">PostgreSQL Bin Path</span>
          <span class="text-[10px] text-teal-accent font-medium">For pg_dump / backups</span>
        </div>
        <input
          v-model="pgBinPath"
          type="text"
          placeholder="C:\Program Files\PostgreSQL\16\bin"
          class="w-full text-xs mt-1 bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-text-primary focus:border-teal-accent focus:outline-none"
        />
        <span class="text-[10px] text-text-muted mt-1 leading-normal">
          Specify the directory containing <code>pg_dump</code> (leave empty if it is already in your system's PATH).
        </span>
      </div>
    </div>
  </Modal>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useUiStore } from '../../stores/ui'
import { useConnectionsStore } from '../../stores/connections'
import { useWorkspaceStore } from '../../stores/workspace'
import Modal from './Modal.vue'

defineEmits(['close'])

const uiStore = useUiStore()
const connectionsStore = useConnectionsStore()
const workspaceStore = useWorkspaceStore()

const isRefreshing = ref(false)

onMounted(() => {
  uiStore.loadSettings()
})

const showFunctionsInSearch = computed({
  get: () => uiStore.settings.showFunctionsInSearch,
  set: (val: boolean) => uiStore.updateSetting('showFunctionsInSearch', val)
})

const pgBinPath = computed({
  get: () => uiStore.settings.pgBinPath,
  set: (val: string) => uiStore.updateSetting('pgBinPath', val)
})

const dbPath = computed({
  get: () => uiStore.settings.dbPath,
  set: (val: string) => uiStore.updateSetting('dbPath', val)
})

async function handleBrowseDb() {
  await uiStore.selectDbPath()
  await connectionsStore.loadConnections()
  await workspaceStore.loadWorkspace()
}

async function handleRefreshDb() {
  isRefreshing.value = true
  try {
    await uiStore.saveAndRefreshDb()
    await connectionsStore.loadConnections()
    await workspaceStore.loadWorkspace()
  } finally {
    isRefreshing.value = false
  }
}
</script>
