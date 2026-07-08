import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { SchemaInfo, TableInfo, ColumnInfo, IndexInfo, FKInfo, ViewInfo, FunctionInfo, SequenceInfo, TypeInfo, TreeNode } from '../types'
import { useConnectionsStore } from './connections'

export const useSchemaStore = defineStore('schema', () => {
  const schemas = ref<SchemaInfo[]>([])
  const tables = ref<TableInfo[]>([])
  const columns = ref<ColumnInfo[]>([])
  const indexes = ref<IndexInfo[]>([])
  const foreignKeys = ref<FKInfo[]>([])
  const views = ref<ViewInfo[]>([])
  const functions = ref<FunctionInfo[]>([])
  const sequences = ref<SequenceInfo[]>([])
  const types = ref<TypeInfo[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const treeData = ref<TreeNode[]>([])

  async function loadSchemas(connId: string) {
    const bindings = useConnectionsStore().getWailsBindings()
    loading.value = true
    error.value = null
    try {
      schemas.value = await bindings.GetSchemas(connId)
      await loadTreeForConnection(connId)
    } catch (e: any) {
      error.value = e.message
      treeData.value = []
      throw e
    } finally {
      loading.value = false
    }
  }

  async function loadTables(connId: string, schema: string) {
    const bindings = useConnectionsStore().getWailsBindings()
    try {
      const schemaTables = await bindings.GetTables(connId, schema)
      // Update tree data
      const schemaNode = treeData.value.find(n => n.id === `schema_${schema}`)
      if (schemaNode) {
        schemaNode.children = schemaTables.map(t => ({
          id: `table_${schema}_${t.name}`,
          label: t.name,
          type: 'table' as const,
          icon: 'table',
          data: t,
        }))
      }
      return schemaTables
    } catch (e: any) {
      error.value = e.message
      return []
    }
  }

  async function loadColumns(connId: string, schema: string, table: string) {
    const bindings = useConnectionsStore().getWailsBindings()
    try {
      columns.value = await bindings.GetColumns(connId, schema, table)
      return columns.value
    } catch (e: any) {
      error.value = e.message
      return []
    }
  }

  async function loadIndexes(connId: string, schema: string, table: string) {
    const bindings = useConnectionsStore().getWailsBindings()
    try {
      indexes.value = await bindings.GetIndexes(connId, schema, table)
      return indexes.value
    } catch (e: any) {
      error.value = e.message
      return []
    }
  }

  async function loadForeignKeys(connId: string, schema: string, table: string) {
    const bindings = useConnectionsStore().getWailsBindings()
    try {
      foreignKeys.value = await bindings.GetForeignKeys(connId, schema, table)
      return foreignKeys.value
    } catch (e: any) {
      error.value = e.message
      return []
    }
  }

  async function loadViews(connId: string, schema: string) {
    const bindings = useConnectionsStore().getWailsBindings()
    try {
      views.value = await bindings.GetViews(connId, schema)
      return views.value
    } catch (e: any) {
      error.value = e.message
      return []
    }
  }

  async function loadFunctions(connId: string, schema: string) {
    const bindings = useConnectionsStore().getWailsBindings()
    try {
      functions.value = await bindings.GetFunctions(connId, schema)
      return functions.value
    } catch (e: any) {
      error.value = e.message
      return []
    }
  }

  async function loadSequences(connId: string, schema: string) {
    const bindings = useConnectionsStore().getWailsBindings()
    try {
      sequences.value = await bindings.GetSequences(connId, schema)
      return sequences.value
    } catch (e: any) {
      error.value = e.message
      return []
    }
  }

  async function loadTypes(connId: string, schema: string) {
    const bindings = useConnectionsStore().getWailsBindings()
    try {
      types.value = await bindings.GetTypes(connId, schema)
      return types.value
    } catch (e: any) {
      error.value = e.message
      return []
    }
  }

  async function getTableDDL(connId: string, schema: string, table: string): Promise<string> {
    const bindings = useConnectionsStore().getWailsBindings()
    try {
      return await bindings.GetTableDDL(connId, schema, table)
    } catch (e: any) {
      error.value = e.message
      return ''
    }
  }

  async function loadTreeForConnection(connId: string) {
    const bindings = useConnectionsStore().getWailsBindings()
    treeData.value = []

    for (const schema of schemas.value) {
      const schemaNode: TreeNode = {
        id: `schema_${schema.name}`,
        label: schema.name,
        type: 'schema',
        icon: 'database',
        children: [],
        expanded: false,
      }

      let tableNodes: TreeNode[] = []
      let viewNodes: TreeNode[] = []
      let funcNodes: TreeNode[] = []
      let seqNodes: TreeNode[] = []
      let typeNodes: TreeNode[] = []

      // Load tables
      try {
        const schemaTables = await bindings.GetTables(connId, schema.name)
        tableNodes = schemaTables.map(t => ({
          id: `table_${schema.name}_${t.name}`,
          label: t.name,
          type: 'table' as const,
          icon: t.type === 'view' ? 'view' : 'table',
          data: t,
        }))
      } catch (e) {
        console.error(`Failed to load tables for schema ${schema.name}:`, e)
      }

      // Load views
      try {
        const schemaViews = await bindings.GetViews(connId, schema.name)
        viewNodes = schemaViews.map(v => ({
          id: `view_${schema.name}_${v.name}`,
          label: v.name,
          type: 'view' as const,
          icon: 'view',
          data: v,
        }))
      } catch (e) {
        console.error(`Failed to load views for schema ${schema.name}:`, e)
      }

      // Load functions
      try {
        const schemaFuncs = await bindings.GetFunctions(connId, schema.name)
        funcNodes = schemaFuncs.map(f => ({
          id: `func_${schema.name}_${f.name}`,
          label: f.name,
          type: 'function' as const,
          icon: 'function',
          data: f,
        }))
      } catch (e) {
        console.error(`Failed to load functions for schema ${schema.name}:`, e)
      }

      // Load sequences
      try {
        const schemaSeqs = await bindings.GetSequences(connId, schema.name)
        seqNodes = schemaSeqs.map(s => ({
          id: `seq_${schema.name}_${s.name}`,
          label: s.name,
          type: 'sequence' as const,
          icon: 'sequence',
          data: s,
        }))
      } catch (e) {
        console.error(`Failed to load sequences for schema ${schema.name}:`, e)
      }

      // Load custom types
      try {
        const schemaTypes = await bindings.GetTypes(connId, schema.name)
        typeNodes = schemaTypes.map(t => ({
          id: `type_${schema.name}_${t.name}`,
          label: t.name,
          type: 'function' as const,
          icon: 'type',
          data: t,
        }))
      } catch (e) {
        console.error(`Failed to load types for schema ${schema.name}:`, e)
      }

      schemaNode.children = [...tableNodes, ...viewNodes, ...funcNodes, ...seqNodes, ...typeNodes]
      treeData.value.push(schemaNode)
    }
  }

  async function refreshSchema(connId: string) {
    const bindings = useConnectionsStore().getWailsBindings()
    loading.value = true
    try {
      await bindings.RefreshSchema(connId)
      await loadSchemas(connId)
    } catch (e: any) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  return {
    schemas,
    tables,
    columns,
    indexes,
    foreignKeys,
    views,
    functions,
    sequences,
    types,
    treeData,
    loading,
    error,
    loadSchemas,
    loadTables,
    loadColumns,
    loadIndexes,
    loadForeignKeys,
    loadViews,
    loadFunctions,
    loadSequences,
    loadTypes,
    getTableDDL,
    refreshSchema,
    loadTreeForConnection,
  }
})
