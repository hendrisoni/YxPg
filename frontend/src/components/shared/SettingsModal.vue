<template>
  <Modal :show="true" title="Settings" @close="$emit('close')" size="sm">
    <div class="space-y-4 py-2 select-none">
      <div class="flex flex-col gap-3">
        <label class="flex items-start gap-3 cursor-pointer p-2 rounded hover:bg-navy-hover transition-colors group">
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
    </div>
  </Modal>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useUiStore } from '../../stores/ui'
import Modal from './Modal.vue'

defineEmits(['close'])

const uiStore = useUiStore()

const showFunctionsInSearch = computed({
  get: () => uiStore.settings.showFunctionsInSearch,
  set: (val: boolean) => uiStore.updateSetting('showFunctionsInSearch', val)
})
</script>
