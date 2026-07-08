<template>
  <div class="h-full flex flex-col overflow-hidden">
    <!-- Toolbar -->
    <div class="flex items-center gap-2 px-3 py-1.5 border-b border-navy-border bg-navy-secondary">
      <button
        @click="addColumn"
        class="flex items-center gap-1 px-2 py-1 text-xs bg-teal-accent text-navy-primary rounded font-medium hover:bg-teal-hover transition-colors"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 5v14M5 12h14" />
        </svg>
        Add Column
      </button>

      <button
        @click="applyChanges"
        :disabled="!hasChanges"
        class="flex items-center gap-1 px-2 py-1 text-xs border border-navy-border rounded text-text-secondary hover:bg-navy-hover transition-colors disabled:opacity-50"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M20 6 9 17l-5-5" />
        </svg>
        Apply DDL
      </button>

      <div class="flex-1"></div>

      <span class="text-xs text-text-muted">{{ columns.length }} columns</span>
    </div>

    <!-- Schema & Table name -->
    <div class="flex items-center gap-3 px-3 py-2 border-b border-navy-border bg-navy-secondary">
      <div class="flex items-center gap-2">
        <label class="text-xs text-text-muted">Schema:</label>
        <select v-model="schemaName" class="px-2 py-1 text-xs">
          <option value="public">public</option>
        </select>
      </div>
      <div class="flex items-center gap-2">
        <label class="text-xs text-text-muted">Table:</label>
        <input
          v-model="tableName"
          type="text"
          placeholder="new_table"
          class="px-2 py-1 text-xs w-48"
        />
      </div>
    </div>

    <!-- Column editor -->
    <div class="flex-1 overflow-auto">
      <table class="w-full text-xs">
        <thead class="bg-navy-secondary sticky top-0">
          <tr class="border-b border-navy-border">
            <th class="px-3 py-2 text-left text-text-muted font-medium">#</th>
            <th class="px-3 py-2 text-left text-text-muted font-medium">Column Name</th>
            <th class="px-3 py-2 text-left text-text-muted font-medium">Data Type</th>
            <th class="px-3 py-2 text-left text-text-muted font-medium">Length</th>
            <th class="px-3 py-2 text-left text-text-muted font-medium">Nullable</th>
            <th class="px-3 py-2 text-left text-text-muted font-medium">Default</th>
            <th class="px-3 py-2 text-left text-text-muted font-medium">PK</th>
            <th class="px-3 py-2 text-left text-text-muted font-medium">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(col, idx) in columns"
            :key="idx"
            class="border-b border-navy-border hover:bg-navy-hover transition-colors"
            :class="{
              'bg-accent-green/5': col.isNew,
              'bg-accent-amber/5': col.isModified,
              'bg-accent-red/5': col.isDeleted,
            }"
          >
            <td class="px-3 py-1.5 text-text-muted">{{ idx + 1 }}</td>
            <td class="px-3 py-1.5">
              <input
                v-model="col.name"
                type="text"
                class="w-full px-2 py-1 text-xs bg-navy-tertiary border border-navy-border rounded"
                placeholder="column_name"
              />
            </td>
            <td class="px-3 py-1.5">
              <select
                v-model="col.data_type"
                class="w-full px-2 py-1 text-xs bg-navy-tertiary border border-navy-border rounded"
              >
                <optgroup v-for="cat in pgTypeCategories" :key="cat" :label="cat">
                  <option v-for="t in getTypesForCategory(cat)" :key="t.value" :value="t.value">
                    {{ t.label }}
                  </option>
                </optgroup>
              </select>
            </td>
            <td class="px-3 py-1.5">
              <input
                v-model.number="col.length"
                type="number"
                class="w-20 px-2 py-1 text-xs bg-navy-tertiary border border-navy-border rounded"
                placeholder=""
              />
            </td>
            <td class="px-3 py-1.5 text-center">
              <input type="checkbox" v-model="col.is_nullable" class="rounded" />
            </td>
            <td class="px-3 py-1.5">
              <input
                v-model="col.default_value"
                type="text"
                class="w-full px-2 py-1 text-xs bg-navy-tertiary border border-navy-border rounded"
                placeholder=""
              />
            </td>
            <td class="px-3 py-1.5 text-center">
              <input type="checkbox" v-model="col.is_primary_key" class="rounded" />
            </td>
            <td class="px-3 py-1.5">
              <button
                @click="removeColumn(idx)"
                class="p-1 rounded hover:bg-navy-tertiary text-text-muted hover:text-accent-red transition-colors"
              >
                <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M3 6h18M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
                </svg>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- DDL Preview -->
    <div class="border-t border-navy-border bg-navy-secondary" style="height: 150px">
      <div class="flex items-center justify-between px-3 py-1 border-b border-navy-border">
        <span class="text-[10px] text-text-muted uppercase tracking-wider">DDL Preview</span>
      </div>
      <pre class="px-3 py-2 text-xs font-mono text-text-secondary overflow-auto h-full">{{ generatedDDL }}</pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useConnectionsStore } from '../../stores/connections'
import { useUiStore } from '../../stores/ui'
import { PG_TYPE_CATEGORIES, getTypesByCategory } from '../../utils/type-mapper'
import type { Tab, ColumnDefinition } from '../../types'

const props = defineProps<{
  tab: Tab
}>()

const connectionsStore = useConnectionsStore()
const uiStore = useUiStore()

const schemaName = ref(props.tab.schema || 'public')
const tableName = ref(props.tab.table || '')
const columns = ref<ExtendedColumn[]>([])
const hasChanges = ref(false)

interface ExtendedColumn extends ColumnDefinition {
  isNew?: boolean
  isModified?: boolean
  isDeleted?: boolean
}

const pgTypeCategories = PG_TYPE_CATEGORIES

function getTypesForCategory(category: string) {
  return getTypesByCategory(category)
}

onMounted(async () => {
  const connId = props.tab.connectionId || connectionsStore.currentConnectionId
  if (props.tab.table && connId) {
    await loadExistingTable()
  } else {
    // Default columns for new table
    columns.value = [
      createDefaultColumn('id', 'bigint', false, true),
    ]
  }
})

async function loadExistingTable() {
  const connId = props.tab.connectionId || connectionsStore.currentConnectionId
  if (!connId) return
  try {
    const bindings = connectionsStore.getWailsBindings()
    const cols = await bindings.GetColumns(
      connId,
      schemaName.value,
      tableName.value
    )
    columns.value = cols.map((c: any) => ({
      name: c.column_name,
      data_type: c.data_type,
      length: c.character_maximum_length,
      is_nullable: c.is_nullable,
      default_value: c.default_value || '',
      is_primary_key: c.is_primary_key,
      is_unique: false,
      is_auto_increment: false,
      isNew: false,
      isModified: false,
      isDeleted: false,
    }))
  } catch (e: any) {
    uiStore.addNotification({ type: 'error', title: 'Load Failed', message: e.message })
  }
}

function createDefaultColumn(name: string, dataType: string, nullable: boolean = true, pk: boolean = false): ExtendedColumn {
  return {
    name,
    data_type: dataType,
    is_nullable: nullable,
    default_value: '',
    is_primary_key: pk,
    is_unique: false,
    is_auto_increment: false,
    isNew: true,
  }
}

function addColumn() {
  columns.value.push(createDefaultColumn(`column_${columns.value.length + 1}`, 'text'))
  hasChanges.value = true
}

function removeColumn(idx: number) {
  columns.value.splice(idx, 1)
  hasChanges.value = true
}

const generatedDDL = computed(() => {
  if (!tableName.value || columns.value.length === 0) {
    return '-- Add a table name and columns to generate DDL'
  }

  const colDefs = columns.value.map(col => {
    let def = `  ${col.name} ${col.data_type}`
    if (col.length) def += `(${col.length})`
    if (!col.is_nullable) def += ' NOT NULL'
    if (col.default_value) def += ` DEFAULT ${col.default_value}`
    if (col.is_unique) def += ' UNIQUE'
    return def
  })

  const pkCols = columns.value.filter(c => c.is_primary_key)
  if (pkCols.length > 0) {
    colDefs.push(`  CONSTRAINT ${tableName.value}_pkey PRIMARY KEY (${pkCols.map(c => c.name).join(', ')})`)
  }

  return `CREATE TABLE ${schemaName.value}.${tableName.value} (\n${colDefs.join(',\n')}\n);`
})

async function applyChanges() {
  const connId = props.tab.connectionId || connectionsStore.currentConnectionId
  if (!connId) return

  try {
    const bindings = connectionsStore.getWailsBindings()
    await bindings.CreateTable(connId, {
      schema: schemaName.value,
      table_name: tableName.value,
      columns: columns.value.map(c => ({
        name: c.name,
        data_type: c.data_type,
        length: c.length,
        is_nullable: c.is_nullable,
        default_value: c.default_value,
        is_primary_key: c.is_primary_key,
        is_unique: c.is_unique,
        is_auto_increment: c.is_auto_increment,
      })),
    } as any)

    uiStore.addNotification({
      type: 'success',
      title: 'Table Created',
      message: `${schemaName.value}.${tableName.value} created successfully`,
    })
  } catch (e: any) {
    uiStore.addNotification({
      type: 'error',
      title: 'Create Failed',
      message: e.message,
    })
  }
}
</script>
