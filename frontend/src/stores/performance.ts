import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import type { PerformanceRecord } from '../types'
import { parsePostgresExplainJson } from '../utils/explainParser'

const STORAGE_KEY = 'yxpg:performance_records'

function hashQuery(sql: string): string {
  const normalized = sql.replace(/\s+/g, ' ').trim().toLowerCase()
  let hash = 5381
  for (let i = 0; i < normalized.length; i++) {
    hash = (hash * 33) ^ normalized.charCodeAt(i)
  }
  return (hash >>> 0).toString(16)
}

function deriveQueryName(sql: string): string {
  const clean = sql.replace(/--.*$/gm, '').replace(/\/\*[\s\S]*?\*\//g, '').trim()
  const firstLine = clean.split('\n')[0].trim()
  if (firstLine.length > 40) {
    return firstLine.substring(0, 40) + '...'
  }
  return firstLine || 'EXPLAIN ANALYZE'
}

export const usePerformanceStore = defineStore('performance', () => {
  // Load saved records from localStorage
  const savedStr = localStorage.getItem(STORAGE_KEY)
  const records = ref<PerformanceRecord[]>(savedStr ? JSON.parse(savedStr) : [])

  const activeRecordId = ref<string | null>(null)
  const compareRecordId = ref<string | null>(null)

  // Filters
  const searchFilter = ref('')
  const databaseFilter = ref('')
  const schemaFilter = ref('')
  const dateFilter = ref('')

  // Watch & persist to localStorage
  watch(
    records,
    (newRecs) => {
      // Keep up to 200 most recent records
      const capped = newRecs.slice(0, 200)
      localStorage.setItem(STORAGE_KEY, JSON.stringify(capped))
    },
    { deep: true }
  )

  // Filtered Records List
  const filteredRecords = computed(() => {
    return records.value.filter((rec) => {
      if (databaseFilter.value && rec.database !== databaseFilter.value) return false
      if (schemaFilter.value && rec.schema !== schemaFilter.value) return false
      if (dateFilter.value) {
        const recDate = rec.executed_at.substring(0, 10)
        if (recDate !== dateFilter.value) return false
      }
      if (searchFilter.value) {
        const q = searchFilter.value.toLowerCase()
        const matchSql = rec.sql.toLowerCase().includes(q)
        const matchName = rec.query_name.toLowerCase().includes(q)
        const matchDb = rec.database.toLowerCase().includes(q)
        if (!matchSql && !matchName && !matchDb) return false
      }
      return true
    })
  })

  // Active Record
  const activeRecord = computed(() => {
    if (!records.value.length) return undefined
    if (activeRecordId.value) {
      return records.value.find((r) => r.id === activeRecordId.value) || records.value[0]
    }
    return records.value[0]
  })

  // Same Query History (Timeline for current query hash)
  const activeQueryHistory = computed(() => {
    if (!activeRecord.value) return []
    return records.value
      .filter((r) => r.query_hash === activeRecord.value!.query_hash)
      .sort((a, b) => new Date(b.executed_at).getTime() - new Date(a.executed_at).getTime())
  })

  // Baseline Comparison Record
  const compareRecord = computed(() => {
    if (!activeRecord.value) return undefined
    if (compareRecordId.value) {
      const found = records.value.find((r) => r.id === compareRecordId.value)
      if (found) return found
    }
    // Default to the previous execution of the SAME query hash
    const history = activeQueryHistory.value
    const idx = history.findIndex((r) => r.id === activeRecord.value!.id)
    if (idx >= 0 && idx + 1 < history.length) {
      return history[idx + 1]
    }
    return undefined
  })

  // Add a new EXPLAIN result
  function recordExplainResult(
    connectionId: string,
    database: string,
    schema: string,
    sql: string,
    rawExplainJson: string
  ): PerformanceRecord {
    const qHash = hashQuery(sql)
    const qName = deriveQueryName(sql)

    const parsed = parsePostgresExplainJson(rawExplainJson)

    const newRecord: PerformanceRecord = {
      id: `perf_${Date.now()}_${Math.random().toString(36).substring(2, 7)}`,
      connection_id: connectionId,
      database: database || 'default',
      schema: schema || 'public',
      query_hash: qHash,
      query_name: qName,
      sql: sql,
      sql_length: sql.length,
      executed_at: new Date().toISOString(),
      execution_time: parsed.executionTime,
      planning_time: parsed.planningTime,
      total_time: parsed.executionTime + parsed.planningTime,
      total_cost: parsed.totalCost,
      startup_cost: parsed.startupCost,
      plan_rows: parsed.planRows,
      actual_rows: parsed.actualRows,
      shared_hit: parsed.sharedHit,
      shared_read: parsed.sharedRead,
      temp_read: parsed.tempRead,
      temp_written: parsed.tempWritten,
      scan_types: parsed.scanTypes,
      join_types: parsed.joinTypes,
      top_node_type: parsed.topNodeType,
      plan_tree: parsed.planTree,
      raw_explain_json: rawExplainJson
    }

    // Unshift to top of records list
    records.value.unshift(newRecord)
    activeRecordId.value = newRecord.id

    // Auto set baseline to previous run of same hash if available
    const sameHashHistory = records.value.filter((r) => r.query_hash === qHash && r.id !== newRecord.id)
    if (sameHashHistory.length > 0) {
      compareRecordId.value = sameHashHistory[0].id
    } else {
      compareRecordId.value = null
    }

    return newRecord
  }

  function setActiveRecord(id: string) {
    activeRecordId.value = id
    // Auto set baseline comparison
    const history = records.value
      .filter((r) => r.query_hash === activeRecord.value?.query_hash)
      .sort((a, b) => new Date(b.executed_at).getTime() - new Date(a.executed_at).getTime())
    const idx = history.findIndex((r) => r.id === id)
    if (idx >= 0 && idx + 1 < history.length) {
      compareRecordId.value = history[idx + 1].id
    } else {
      compareRecordId.value = null
    }
  }

  function setCompareRecord(id: string | null) {
    compareRecordId.value = id
  }

  function updateRemark(id: string, remark: string) {
    const rec = records.value.find(r => r.id === id)
    if (rec) {
      rec.remark = remark
    }
  }

  function deleteRecord(id: string) {
    records.value = records.value.filter(r => r.id !== id)
    if (activeRecordId.value === id) {
      activeRecordId.value = records.value[0]?.id || null
    }
    if (compareRecordId.value === id) {
      compareRecordId.value = records.value[1]?.id || null
    }
  }

  function clearHistory() {
    records.value = []
    activeRecordId.value = null
    compareRecordId.value = null
  }

  return {
    records,
    activeRecordId,
    compareRecordId,
    searchFilter,
    databaseFilter,
    schemaFilter,
    dateFilter,
    filteredRecords,
    activeRecord,
    compareRecord,
    activeQueryHistory,
    recordExplainResult,
    setActiveRecord,
    setCompareRecord,
    updateRemark,
    deleteRecord,
    clearHistory
  }
})
