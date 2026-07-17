<template>
  <Teleport to="body">
    <Transition name="notification">
      <div
        v-if="show"
        class="fixed bottom-4 right-4 z-50 flex items-start gap-3 px-4 py-3 rounded-lg border shadow-xl max-w-sm animate-slide-in"
        :class="typeClass"
      >
        <!-- Icon -->
        <div class="flex-shrink-0 mt-0.5">
          <svg v-if="type === 'success'" class="w-4 h-4 text-accent-green" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 6 9 17l-5-5" />
          </svg>
          <svg v-else-if="type === 'error'" class="w-4 h-4 text-accent-red" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10" /><path d="m15 9-6 6M9 9l6 6" />
          </svg>
          <svg v-else-if="type === 'warning'" class="w-4 h-4 text-accent-amber" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z" />
            <path d="M12 9v4M12 17h.01" />
          </svg>
          <svg v-else class="w-4 h-4 text-accent-blue" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10" /><path d="M12 16v-4M12 8h.01" />
          </svg>
        </div>

        <!-- Content -->
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium text-text-primary">{{ title }}</p>
          <p v-if="message" class="text-xs text-text-secondary mt-0.5">{{ message }}</p>
          <button
            v-if="action"
            @click="handleAction"
            class="mt-1.5 px-2 py-0.5 rounded text-[10px] font-semibold bg-teal-accent/20 hover:bg-teal-accent/35 text-teal-accent border border-teal-accent/30 hover:border-teal-accent/50 transition-colors cursor-pointer"
          >
            {{ action.label }}
          </button>
        </div>

        <!-- Close -->
        <button
          @click="$emit('close')"
          class="flex-shrink-0 p-0.5 rounded hover:bg-white/10 text-text-muted hover:text-text-primary transition-colors"
        >
          <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M18 6 6 18M6 6l12 12" />
          </svg>
        </button>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  show: boolean
  type: 'success' | 'error' | 'warning' | 'info'
  title: string
  message?: string
  action?: {
    label: string
    callback: () => void
  }
}>()

const emit = defineEmits(['close'])

function handleAction() {
  if (props.action?.callback) {
    props.action.callback()
  }
  emit('close')
}

const typeClass = computed(() => {
  const classes = {
    success: 'bg-navy-secondary border-accent-green/30',
    error: 'bg-navy-secondary border-accent-red/30',
    warning: 'bg-navy-secondary border-accent-amber/30',
    info: 'bg-navy-secondary border-accent-blue/30',
  }
  return classes[props.type]
})
</script>

<style scoped>
.notification-enter-active,
.notification-leave-active {
  transition: all 0.3s ease;
}

.notification-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.notification-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
</style>
