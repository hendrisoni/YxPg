<template>
  <Modal :show="show" :title="isReadOnly ? 'View Row Details' : 'Edit Row Data'" @close="$emit('close')" size="lg">
    <form @submit.prevent="handleSubmit" class="space-y-4 max-h-[70vh] overflow-y-auto pr-1">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div v-for="(col, idx) in columns" :key="col.name" class="flex flex-col gap-1.5">
          <label :for="`field_${idx}`" class="text-xs font-semibold text-text-primary flex items-center justify-between">
            <span>{{ col.name }}</span>
            <span class="text-[9px] text-text-muted font-normal uppercase">{{ getTypeName(col.data_type) }}</span>
          </label>

          <!-- Input control depending on data type -->
          <!-- Date Field -->
          <input
            v-if="isDateField(col.data_type)"
            :id="`field_${idx}`"
            v-model="formValues[`col_${idx}`]"
            type="date"
            :disabled="isReadOnly"
            class="w-full text-xs bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-text-primary focus:border-teal-accent focus:outline-none placeholder-text-muted disabled:opacity-60 disabled:cursor-not-allowed"
          />

          <!-- Textarea for large text or json -->
          <textarea
            v-else-if="isLargeTextField(col.data_type)"
            :id="`field_${idx}`"
            v-model="formValues[`col_${idx}`]"
            rows="3"
            :disabled="isReadOnly"
            class="w-full text-xs bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-text-primary focus:border-teal-accent focus:outline-none placeholder-text-muted disabled:opacity-60 disabled:cursor-not-allowed font-mono"
          ></textarea>

          <!-- Boolean Select -->
          <select
            v-else-if="isBooleanField(col.data_type)"
            :id="`field_${idx}`"
            v-model="formValues[`col_${idx}`]"
            :disabled="isReadOnly"
            class="w-full text-xs bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-text-primary focus:border-teal-accent focus:outline-none placeholder-text-muted disabled:opacity-60 disabled:cursor-not-allowed"
          >
            <option :value="null">NULL</option>
            <option :value="true">true</option>
            <option :value="false">false</option>
          </select>

          <!-- Numeric Field -->
          <input
            v-else-if="isNumericField(col.data_type)"
            :id="`field_${idx}`"
            v-model.number="formValues[`col_${idx}`]"
            type="number"
            step="any"
            :disabled="isReadOnly"
            class="w-full text-xs bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-text-primary focus:border-teal-accent focus:outline-none placeholder-text-muted disabled:opacity-60 disabled:cursor-not-allowed text-right"
          />

          <!-- Default Text Field -->
          <input
            v-else
            :id="`field_${idx}`"
            v-model="formValues[`col_${idx}`]"
            type="text"
            :disabled="isReadOnly"
            class="w-full text-xs bg-navy-tertiary border border-navy-border rounded-md px-2.5 py-1.5 text-text-primary focus:border-teal-accent focus:outline-none placeholder-text-muted disabled:opacity-60 disabled:cursor-not-allowed"
          />
        </div>
      </div>

      <!-- Action Buttons -->
      <div v-if="!isReadOnly" class="flex justify-end gap-2 mt-4 pt-3 border-t border-navy-border">
        <button
          type="button"
          @click="$emit('close')"
          class="px-4 py-1.5 text-xs text-text-secondary hover:text-text-primary transition-colors cursor-pointer"
        >
          Cancel
        </button>
        <button
          type="submit"
          class="px-4 py-1.5 text-xs bg-teal-accent text-navy-primary rounded-md font-medium hover:bg-teal-hover transition-colors cursor-pointer"
        >
          Save Changes
        </button>
      </div>
    </form>
  </Modal>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import Modal from '../shared/Modal.vue'
import type { ColumnMeta, Tab } from '../../types'

const props = defineProps<{
  show: boolean
  rowData: Record<string, any>
  columns: ColumnMeta[]
  tab?: Tab
}>()

const emit = defineEmits(['close', 'save'])

const isReadOnly = computed(() => {
  return !(props.tab && props.tab.type === 'table')
})

const formValues = ref<Record<string, any>>({})

onMounted(() => {
  // Initialize form values from rowData
  props.columns.forEach((col, idx) => {
    const val = props.rowData[`col_${idx}`]
    
    // If it's a date type, ensure it's formatted as yyyy-MM-dd for the date input
    if (isDateField(col.data_type) && typeof val === 'string') {
      formValues.value[`col_${idx}`] = formatDateForInput(val)
    } else {
      formValues.value[`col_${idx}`] = val
    }
  })
})

function handleSubmit() {
  emit('save', { ...formValues.value })
}

// Helpers for data types identification
function getTypeName(oid: string): string {
  const o = parseInt(oid)
  if (isNaN(o)) return oid
  
  const types: Record<number, string> = {
    16: 'boolean',
    17: 'bytea',
    20: 'bigint',
    21: 'integer',
    23: 'integer',
    25: 'text',
    114: 'json',
    700: 'real',
    701: 'double precision',
    1042: 'char',
    1043: 'varchar',
    1082: 'date',
    1083: 'time',
    1114: 'timestamp',
    1184: 'timestamptz',
    1700: 'numeric',
    2950: 'uuid',
    3802: 'jsonb'
  }
  return types[o] || `type_${oid}`
}

function isDateField(oid: string): boolean {
  const name = getTypeName(oid).toLowerCase()
  return name.includes('date') || name.includes('time') || name.includes('timestamp')
}

function isLargeTextField(oid: string): boolean {
  const name = getTypeName(oid).toLowerCase()
  return name === 'text' || name.includes('json') || name.includes('bytea')
}

function isBooleanField(oid: string): boolean {
  const name = getTypeName(oid).toLowerCase()
  return name === 'boolean' || name === 'bool'
}

function isNumericField(oid: string): boolean {
  const name = getTypeName(oid).toLowerCase()
  return (
    name.includes('int') ||
    name.includes('num') ||
    name.includes('decimal') ||
    name.includes('real') ||
    name.includes('double') ||
    name.includes('float') ||
    name.includes('precision') ||
    name.includes('money')
  )
}

function formatDateForInput(val: string): string {
  const match = val.match(/^(\d{4}-\d{2}-\d{2})/)
  return match ? match[1] : val
}
</script>
