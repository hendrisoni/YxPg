<template>
  <div class="h-full flex flex-col overflow-hidden">
    <!-- Placeholder -->
    <div class="flex-1 flex items-center justify-center">
      <div class="text-center">
        <svg class="w-16 h-16 mx-auto text-text-muted mb-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
          <polyline points="14 2 14 8 20 8" />
          <line x1="16" y1="13" x2="8" y2="13" />
          <line x1="16" y1="17" x2="8" y2="17" />
        </svg>
        <h3 class="text-lg font-semibold text-text-primary mb-2">DDL Editor</h3>
        <p class="text-sm text-text-muted mb-4">Design tables and manage schema changes</p>
        <p class="text-xs text-text-muted">Coming in Milestone 5</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Tab } from '../types'

const props = defineProps<{
  tab: Tab
}>()
</script>
