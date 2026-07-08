<template>
  <div class="flex items-center justify-center" v-if="loading">
    <div
      class="border-2 border-teal-accent border-t-transparent rounded-full animate-spin"
      :class="sizeClass"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  loading?: boolean
  size?: 'sm' | 'md' | 'lg'
}>(), {
  loading: true,
  size: 'md',
})

const sizeClass = computed(() => {
  const sizes = { sm: 'w-4 h-4', md: 'w-6 h-6', lg: 'w-8 h-8' }
  return sizes[props.size]
})
</script>
