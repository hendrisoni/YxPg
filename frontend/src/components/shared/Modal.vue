<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="show"
        class="fixed inset-0 z-50 flex items-center justify-center overflow-y-auto p-4 md:p-6"
        @click.self="$emit('close')"
      >
        <!-- Backdrop -->
        <div class="fixed inset-0 bg-black/60 backdrop-blur-sm -z-10" @click="$emit('close')"></div>

        <!-- Modal Box -->
        <div
          ref="modalRef"
          class="relative bg-navy-secondary border border-navy-border rounded-lg shadow-xl w-full my-auto z-10 flex flex-col max-h-[90vh] animate-fade-in"
          :class="[sizeClass, { 'select-none': isResizing }]"
          :style="customStyle"
        >
          <!-- Header -->
          <div v-if="title" class="flex items-center justify-between px-4 py-3 border-b border-navy-border flex-shrink-0">
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
          <div class="px-4 py-3 flex-1 flex flex-col min-h-0 overflow-hidden">
            <slot></slot>
          </div>

          <!-- Footer -->
          <div v-if="$slots.footer" class="px-4 py-3 border-t border-navy-border flex justify-end gap-2 flex-shrink-0">
            <slot name="footer"></slot>
          </div>

          <!-- Bottom-Right Resize Handle Grip -->
          <div
            v-if="resizable"
            class="absolute bottom-0 right-0 w-4 h-4 cursor-se-resize flex items-center justify-center text-text-muted hover:text-teal-accent transition-colors z-20"
            @mousedown="startResize"
            title="Drag to resize modal"
          >
            <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15L15 21M21 9L9 21" />
            </svg>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'

const props = withDefaults(defineProps<{
  show: boolean
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
  resizable?: boolean
  storageKey?: string
}>(), {
  size: 'md',
  resizable: false,
})

const emit = defineEmits(['close'])

const modalRef = ref<HTMLElement | null>(null)
const modalWidth = ref<number | null>(null)
const modalHeight = ref<number | null>(null)
const isResizing = ref(false)

const sizeClass = computed(() => {
  if (modalWidth.value) return ''
  const sizes = {
    sm: 'max-w-sm',
    md: 'max-w-lg',
    lg: 'max-w-2xl',
    xl: 'max-w-4xl',
  }
  return sizes[props.size]
})

const customStyle = computed(() => {
  const styles: Record<string, string> = {}
  if (modalWidth.value) {
    styles.width = `${modalWidth.value}px`
    styles.maxWidth = '95vw'
  }
  if (modalHeight.value) {
    styles.height = `${modalHeight.value}px`
    styles.maxHeight = '90vh'
  }
  return styles
})

// Load saved modal size from localStorage when shown
watch(() => props.show, (newVal) => {
  if (newVal && props.storageKey) {
    const saved = localStorage.getItem(props.storageKey)
    if (saved) {
      try {
        const { w, h } = JSON.parse(saved)
        if (w) modalWidth.value = w
        if (h) modalHeight.value = h
      } catch (e) {}
    }
  }
})

let startX = 0
let startY = 0
let startWidth = 0
let startHeight = 0

function startResize(e: MouseEvent) {
  if (!props.resizable || !modalRef.value) return
  e.preventDefault()
  isResizing.value = true
  startX = e.clientX
  startY = e.clientY
  const rect = modalRef.value.getBoundingClientRect()
  startWidth = rect.width
  startHeight = rect.height

  window.addEventListener('mousemove', handleMouseMove)
  window.addEventListener('mouseup', handleMouseUp)
}

function handleMouseMove(e: MouseEvent) {
  if (!isResizing.value) return
  const dx = e.clientX - startX
  const dy = e.clientY - startY
  const newW = Math.max(380, Math.min(window.innerWidth * 0.95, startWidth + dx))
  const newH = Math.max(250, Math.min(window.innerHeight * 0.9, startHeight + dy))
  modalWidth.value = newW
  modalHeight.value = newH
}

function handleMouseUp() {
  if (isResizing.value) {
    isResizing.value = false
    window.removeEventListener('mousemove', handleMouseMove)
    window.removeEventListener('mouseup', handleMouseUp)
    if (props.storageKey && modalWidth.value && modalHeight.value) {
      localStorage.setItem(props.storageKey, JSON.stringify({ w: modalWidth.value, h: modalHeight.value }))
    }
  }
}

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && props.show) {
    emit('close')
  }
}

watch(() => props.show, (newVal) => {
  if (newVal) {
    window.addEventListener('keydown', handleKeydown)
  } else {
    window.removeEventListener('keydown', handleKeydown)
  }
})

onMounted(() => {
  if (props.show) {
    window.addEventListener('keydown', handleKeydown)
  }
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
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
