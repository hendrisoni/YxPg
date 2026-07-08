<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="show"
        class="fixed inset-0 z-50 flex items-start justify-center overflow-y-auto p-4 md:p-6"
        @click.self="$emit('close')"
      >
        <!-- Backdrop -->
        <div class="fixed inset-0 bg-black/60 backdrop-blur-sm -z-10" @click="$emit('close')"></div>

        <!-- Modal -->
        <div
          class="relative bg-navy-secondary border border-navy-border rounded-lg shadow-xl max-w-lg w-full my-auto z-10 animate-fade-in"
          :class="sizeClass"
        >
          <!-- Header -->
          <div v-if="title" class="flex items-center justify-between px-4 py-3 border-b border-navy-border">
            <h3 class="text-sm font-semibold text-text-primary">{{ title }}</h3>
            <button
              @click="$emit('close')"
              class="p-1 rounded hover:bg-navy-hover text-text-muted hover:text-text-primary transition-colors"
            >
              <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6 6 18M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Content -->
          <div class="px-4 py-3">
            <slot></slot>
          </div>

          <!-- Footer -->
          <div v-if="$slots.footer" class="px-4 py-3 border-t border-navy-border flex justify-end gap-2">
            <slot name="footer"></slot>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  show: boolean
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
}>(), {
  size: 'md',
})

defineEmits(['close'])

const sizeClass = computed(() => {
  const sizes = {
    sm: 'max-w-sm',
    md: 'max-w-lg',
    lg: 'max-w-2xl',
    xl: 'max-w-4xl',
  }
  return sizes[props.size]
})
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
