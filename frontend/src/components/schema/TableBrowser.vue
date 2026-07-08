<template>
  <div class="h-full overflow-auto">
    <!-- Table data view -->
    <div ref="tableContainer" class="h-full w-full"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useConnectionsStore } from '../../stores/connections'
import type { TableInfo, BrowseOptions } from '../../types'

const props = defineProps<{
  table: TableInfo
  schema: string
}>()

const connectionsStore = useConnectionsStore()
const tableContainer = ref<HTMLElement | null>(null)
let tabulator: any = null

onMounted(() => {
  loadData()
})

onUnmounted(() => {
  if (tabulator) {
    tabulator.destroy()
  }
})

async function loadData() {
  if (!connectionsStore.currentConnectionId || !tableContainer.value) return

  const bindings = connectionsStore.getWailsBindings()
  const opts: BrowseOptions = {
    page: 1,
    page_size: 100,
  }

  const result = await bindings.BrowseTable(
    connectionsStore.currentConnectionId,
    props.schema,
    props.table.name,
    opts as any
  )

  if (result.error) return

  const { TabulatorFull } = await import('tabulator-tables')

  const columns = result.columns.map((col: any, idx: number) => ({
    title: col.name,
    field: `col_${idx}`,
    formatter: (cell: any) => {
      const val = cell.getValue()
      if (val === null || val === undefined) {
        return '<span style="color: #4B5563; font-style: italic;">NULL</span>'
      }
      return val
    },
  }))

  const data = result.rows.map((row: any[]) => {
    const rowData: Record<string, any> = {}
    row.forEach((val: any, idx: number) => {
      rowData[`col_${idx}`] = val
    })
    return rowData
  })

  tabulator = new TabulatorFull(tableContainer.value, {
    data,
    columns,
    layout: 'fitDataFill',
    height: '100%',
    progressiveRender: 'scroll',
    progressiveRenderSize: 100,
    progressiveRenderMargin: 350,
    movableColumns: true,
    resizableColumns: true,
  })
}
</script>
