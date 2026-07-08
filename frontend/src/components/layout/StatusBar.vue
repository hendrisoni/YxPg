<template>
  <footer class="flex items-center justify-between px-3 h-6 bg-navy-secondary border-t border-navy-border text-xs select-none">
    <!-- Left side -->
    <div class="flex items-center gap-3">
      <!-- Connection info -->
      <div v-if="connectionInfo" class="flex items-center gap-1.5 text-text-secondary">
        <div class="w-1.5 h-1.5 rounded-full bg-accent-green"></div>
        <span>{{ connectionInfo.db }}@{{ connectionInfo.host }}</span>
      </div>
      <div v-else class="text-text-muted">No connection</div>

      <!-- Query stats -->
      <div v-if="lastQueryDuration > 0" class="flex items-center gap-1.5 text-text-muted">
        <span>{{ lastQueryDuration }}ms</span>
        <span v-if="lastRowCount > 0">· {{ lastRowCount }} rows</span>
      </div>
    </div>

    <!-- Right side -->
    <div class="flex items-center gap-3 text-text-muted">
      <span>{{ currentTime }}</span>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useConnectionsStore } from '../../stores/connections'

const connectionsStore = useConnectionsStore()

const currentTime = ref('')
const lastQueryDuration = ref(0)
const lastRowCount = ref(0)
let timeInterval: ReturnType<typeof setInterval> | null = null

const connectionInfo = computed(() => {
  if (!connectionsStore.currentConnection) return null
  return {
    host: connectionsStore.currentConnection.host,
    db: connectionsStore.currentConnection.database,
    user: connectionsStore.currentConnection.username,
  }
})

function updateTime() {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('en-US', { hour12: false })
}

onMounted(() => {
  updateTime()
  timeInterval = setInterval(updateTime, 1000)
})

onUnmounted(() => {
  if (timeInterval) clearInterval(timeInterval)
})

defineExpose({
  setQueryStats(duration: number, rowCount: number) {
    lastQueryDuration.value = duration
    lastRowCount.value = rowCount
  },
})
</script>
