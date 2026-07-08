<template>
  <Teleport to="body">
    <div
      v-if="show"
      class="fixed z-50 bg-navy-secondary border border-navy-border rounded-lg shadow-xl py-1 min-w-[160px] animate-fade-in"
      :style="{ left: x + 'px', top: y + 'px' }"
      @click.stop
    >
      <button
        v-for="item in items"
        :key="item.label"
        @click="handleClick(item)"
        class="w-full px-3 py-1.5 text-xs text-left flex items-center gap-2 transition-colors"
        :class="[
          item.danger ? 'text-accent-red hover:bg-accent-red/10' : 'text-text-secondary hover:bg-navy-hover hover:text-text-primary',
          item.disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'
        ]"
        :disabled="item.disabled"
      >
        <span v-if="item.icon" class="w-4 h-4 flex-shrink-0">{{ item.icon }}</span>
        <span>{{ item.label }}</span>
        <span v-if="item.shortcut" class="ml-auto text-text-muted">{{ item.shortcut }}</span>
      </button>
    </div>

    <!-- Backdrop to close -->
    <div
      v-if="show"
      class="fixed inset-0 z-40"
      @click="$emit('close')"
      @contextmenu.prevent="$emit('close')"
    ></div>
  </Teleport>
</template>

<script setup lang="ts">
export interface ContextMenuItem {
  label: string
  icon?: string
  shortcut?: string
  danger?: boolean
  disabled?: boolean
  action: string
}

defineProps<{
  show: boolean
  x: number
  y: number
  items: ContextMenuItem[]
}>()

const emit = defineEmits(['close', 'select'])

function handleClick(item: ContextMenuItem) {
  if (!item.disabled) {
    emit('select', item.action)
    emit('close')
  }
}
</script>
