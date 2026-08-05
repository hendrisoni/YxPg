<template>
  <div class="font-mono text-xs select-none">
    <div
      @click="$emit('select-node', node)"
      class="group flex items-start gap-2 p-2 my-1 rounded border transition-all cursor-pointer"
      :class="[
        isSelected ? 'border-teal-accent bg-teal-accent/10 shadow-[0_0_12px_rgba(0,201,167,0.15)]' : 'border-navy-border/80 bg-navy-tertiary/70 hover:bg-navy-hover hover:border-navy-border',
        diffClass
      ]"
    >
      <!-- Expand/Collapse toggle if has children -->
      <button
        v-if="node.plans && node.plans.length > 0"
        @click.stop="expanded = !expanded"
        class="mt-0.5 w-4 h-4 flex items-center justify-center rounded text-text-muted hover:text-white hover:bg-white/10"
      >
        <svg
          class="w-3 h-3 transition-transform duration-150"
          :class="{ 'rotate-90': expanded }"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <polyline points="9 18 15 12 9 6" />
        </svg>
      </button>
      <div v-else class="w-4 flex-shrink-0"></div>

      <!-- Main Node Information -->
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2 flex-wrap">
          <!-- Node Type Badge -->
          <span
            class="px-1.5 py-0.5 rounded text-[11px] font-bold tracking-tight shadow-sm"
            :class="nodeTypeBadgeClass"
          >
            {{ node.node_type }}
          </span>

          <!-- Relation / Index Name -->
          <span v-if="node.relation_name" class="text-text-primary font-semibold truncate">
            on <span class="text-teal-accent">{{ node.relation_name }}</span>
            <span v-if="node.alias" class="text-text-muted text-[10px]"> (as {{ node.alias }})</span>
          </span>

          <span v-if="node.index_name" class="text-amber-400 text-[11px] font-mono truncate">
            using {{ node.index_name }}
          </span>

          <!-- Diff status badge if changed -->
          <span v-if="node.diff_status === 'changed'" class="text-[9px] px-1 py-0.2 rounded bg-amber-500/20 text-amber-300 border border-amber-500/40">
            ★ Node Changed
          </span>
        </div>

        <!-- Metrics preview line -->
        <div class="flex items-center gap-3 mt-1 text-[10px] text-text-secondary flex-wrap">
          <span>Cost: <strong class="text-text-primary">{{ node.startup_cost.toFixed(1) }}..{{ node.total_cost.toFixed(1) }}</strong></span>
          <span v-if="node.actual_time !== undefined">Time: <strong class="text-teal-accent">{{ node.actual_time.toFixed(2) }}ms</strong></span>
          <span>Rows: <strong class="text-text-primary">{{ node.actual_rows !== undefined ? node.actual_rows : node.plan_rows }}</strong></span>
          <span v-if="node.actual_loops && node.actual_loops > 1">Loops: <strong class="text-purple-400">{{ node.actual_loops }}</strong></span>
        </div>

        <!-- Extra conditions preview -->
        <div v-if="node.filter" class="text-[10px] text-rose-300/80 truncate mt-0.5">
          Filter: <span class="font-mono opacity-90">{{ node.filter }}</span>
        </div>
        <div v-if="node.index_cond" class="text-[10px] text-emerald-300/80 truncate mt-0.5">
          Index Cond: <span class="font-mono opacity-90">{{ node.index_cond }}</span>
        </div>
      </div>
    </div>

    <!-- Recursive Sub-plans -->
    <div v-if="expanded && node.plans && node.plans.length > 0" class="pl-4 border-l border-navy-border/60 ml-2 space-y-1">
      <PlanNodeTree
        v-for="child in node.plans"
        :key="child.id"
        :node="child"
        :selected-node-id="selectedNodeId"
        @select-node="$emit('select-node', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { ParsedExplainNode } from '../../types'

const props = defineProps<{
  node: ParsedExplainNode
  selectedNodeId?: string
}>()

defineEmits<{
  (e: 'select-node', node: ParsedExplainNode): void
}>()

const expanded = ref(true)

const isSelected = computed(() => props.selectedNodeId === props.node.id)

const nodeTypeBadgeClass = computed(() => {
  const type = props.node.node_type.toLowerCase()
  if (type.includes('index scan') || type.includes('index only')) {
    return 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/40'
  }
  if (type.includes('seq scan')) {
    return 'bg-amber-500/20 text-amber-300 border border-amber-500/40'
  }
  if (type.includes('hash join') || type.includes('merge join')) {
    return 'bg-blue-500/20 text-blue-300 border border-blue-500/40'
  }
  if (type.includes('nested loop')) {
    return 'bg-purple-500/20 text-purple-300 border border-purple-500/40'
  }
  if (type.includes('bitmap')) {
    return 'bg-teal-500/20 text-teal-300 border border-teal-500/40'
  }
  if (type.includes('sort') || type.includes('aggregate')) {
    return 'bg-indigo-500/20 text-indigo-300 border border-indigo-500/40'
  }
  return 'bg-slate-700/50 text-slate-300 border border-slate-600/50'
})

const diffClass = computed(() => {
  if (props.node.diff_status === 'changed') {
    return 'border-amber-500/60 bg-amber-500/5'
  }
  return ''
})
</script>
