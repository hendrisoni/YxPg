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
import { computed } from 'vue'
import { useUiStore } from '../../stores/ui'
import Modal from './Modal.vue'

defineEmits(['close'])

const uiStore = useUiStore()

const showFunctionsInSearch = computed({
  get: () => uiStore.settings.showFunctionsInSearch,
  set: (val: boolean) => uiStore.updateSetting('showFunctionsInSearch', val)
})

const pgBinPath = computed({
  get: () => uiStore.settings.pgBinPath,
  set: (val: string) => uiStore.updateSetting('pgBinPath', val)
})
</script>
