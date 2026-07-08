<template>
  <div class="flex items-center gap-2 px-3 py-2 border-b border-navy-border bg-navy-secondary">
    <!-- Tabs for different details -->
    <div class="flex items-center gap-1">
      <button
        v-for="tab in detailTabs"
        :key="tab.key"
        @click="activeDetailTab = tab.key"
        class="px-2 py-1 text-xs rounded transition-colors"
        :class="activeDetailTab === tab.key ? 'bg-navy-tertiary text-text-primary' : 'text-text-muted hover:text-text-secondary hover:bg-navy-hover'"
      >
        {{ tab.label }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const activeDetailTab = ref('columns')

const detailTabs = [
  { key: 'columns', label: 'Columns' },
  { key: 'indexes', label: 'Indexes' },
  { key: 'foreignKeys', label: 'Foreign Keys' },
  { key: 'triggers', label: 'Triggers' },
  { key: 'ddl', label: 'DDL' },
]
</script>
