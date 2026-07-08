<template>
  <div class="h-full flex flex-col overflow-hidden">
    <QueryBuilder :tab="tab" />
  </div>
</template>

<script setup lang="ts">
import type { Tab } from '../types'
import QueryBuilder from '../components/builder/QueryBuilder.vue'

defineProps<{
  tab: Tab
}>()
</script>
