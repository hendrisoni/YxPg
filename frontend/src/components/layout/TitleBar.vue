<template>
  <div class="h-8 bg-[#060913] border-b border-navy-border flex items-center justify-between px-3 select-none flex-shrink-0 titlebar z-[100]">
    <!-- Left: App Icon & Name -->
    <div class="flex items-center gap-2 text-text-secondary w-1/3">
      <img src="/yxpg.png" class="w-3.5 h-3.5 object-contain" alt="YxPg Logo" />
      <span class="text-xs font-semibold tracking-wide text-text-muted">YxPg</span>
    </div>

    <!-- Center: App Title/Connection Status -->
    <div class="text-xs text-text-muted font-mono truncate max-w-[40%] text-center cursor-default flex-1 px-4">
      {{ titleText }}
    </div>

    <!-- Right: Window Controls -->
    <div class="flex items-center h-full w-1/3 justify-end">
      <!-- Minimize -->
      <button
        @click="minimizeWindow"
        class="h-full px-3 flex items-center justify-center text-text-muted hover:text-text-primary hover:bg-[#1e293b] transition-colors cursor-pointer no-drag"
        title="Minimize"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="5" y1="12" x2="19" y2="12" />
        </svg>
      </button>

      <!-- Maximize / Restore -->
      <button
        @click="maximizeWindow"
        class="h-full px-3 flex items-center justify-center text-text-muted hover:text-text-primary hover:bg-[#1e293b] transition-colors cursor-pointer no-drag"
        title="Maximize"
      >
        <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="3" width="18" height="18" rx="2" />
        </svg>
      </button>

      <!-- Close -->
      <button
        @click="closeWindow"
        class="h-full px-3 flex items-center justify-center text-text-muted hover:text-text-primary hover:bg-accent-red transition-colors cursor-pointer no-drag"
        title="Close"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18" />
          <line x1="6" y1="6" x2="18" y2="18" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { WindowMinimise, WindowToggleMaximise, Quit } from '../../../wailsjs/runtime/runtime'
import { useConnectionsStore } from '../../stores/connections'

const connectionsStore = useConnectionsStore()

const titleText = computed(() => {
  const current = connectionsStore.currentConnection
  if (current) {
    return `YxPg - ${current.name} (${current.host}:${current.port})`
  }
  return 'YxPg - Portable Postgres Studio'
})

function minimizeWindow() {
  WindowMinimise()
}

function maximizeWindow() {
  WindowToggleMaximise()
}

function closeWindow() {
  Quit()
}
</script>

<style scoped>
.titlebar {
  --wails-draggable:drag;
}
.no-drag {
  --wails-draggable:no-drag;
}
</style>
