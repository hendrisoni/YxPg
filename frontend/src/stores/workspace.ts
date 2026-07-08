import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { TreeNode } from '../types'
import { useConnectionsStore } from './connections'

export interface CatalogItem {
  connection_id: string
  connection_name: string
  database_name: string
  schema: string
  name: string
  type: 'table' | 'view' | 'function'
}

export const useWorkspaceStore = defineStore('workspace', () => {
  const workspaceTree = ref<TreeNode[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Load workspace tree from backend
  async function loadWorkspace() {
    const bindings = useConnectionsStore().getWailsBindings()
    loading.value = true
    try {
      const dataStr = await (bindings as any).LoadWorkspace()
      workspaceTree.value = JSON.parse(dataStr)
    } catch (e: any) {
      error.value = e.message || String(e)
      workspaceTree.value = []
    } finally {
      loading.value = false
    }
  }

  // Save workspace tree to backend
  async function saveWorkspace() {
    const bindings = useConnectionsStore().getWailsBindings()
    try {
      const dataStr = JSON.stringify(workspaceTree.value)
      await (bindings as any).SaveWorkspace(dataStr)
    } catch (e: any) {
      error.value = e.message || String(e)
    }
  }

  // Add a custom category to the root of the tree
  async function addCategory(name: string) {
    const categoryNode: TreeNode = {
      id: `category_${Date.now()}_${Math.random().toString(36).substr(2, 5)}`,
      label: name,
      type: 'category',
      icon: 'database', // category folder icon
      children: [],
      expanded: true,
    }
    workspaceTree.value.push(categoryNode)
    await saveWorkspace()
  }

  // Add an item (from Ctrl+K catalog search) to the root of the tree
  async function addObject(item: CatalogItem) {
    const nodeId = `item_${item.connection_id}_${item.schema}_${item.name}_${item.type}`
    
    // Check if it already exists in the tree (recursively)
    if (findNodeById(workspaceTree.value, nodeId)) {
      return // Already exists
    }

    const node: TreeNode = {
      id: nodeId,
      label: item.name,
      type: item.type,
      icon: item.type,
      children: [],
      data: {
        connectionId: item.connection_id,
        connectionName: item.connection_name,
        database: item.database_name,
        schema: item.schema,
        name: item.name,
        type: item.type,
      }
    }

    workspaceTree.value.push(node)
    await saveWorkspace()
  }

  // Helper to recursively find a node by ID
  function findNodeById(nodes: TreeNode[], id: string): TreeNode | null {
    for (const node of nodes) {
      if (node.id === id) return node
      if (node.children && node.children.length > 0) {
        const found = findNodeById(node.children, id)
        if (found) return found
      }
    }
    return null
  }

  // Helper to recursively remove a node by ID and return it
  function removeNodeById(nodes: TreeNode[], id: string): TreeNode | null {
    for (let i = 0; i < nodes.length; i++) {
      const node = nodes[i]
      if (node.id === id) {
        const [removed] = nodes.splice(i, 1)
        return removed
      }
      if (node.children && node.children.length > 0) {
        const removed = removeNodeById(node.children, id)
        if (removed) return removed
      }
    }
    return null
  }

  // Delete a node by ID
  async function deleteNode(id: string) {
    removeNodeById(workspaceTree.value, id)
    await saveWorkspace()
  }

  // Move a node to a target category (or root if targetCategoryId is null)
  async function moveNode(nodeId: string, targetCategoryId: string | null) {
    // 1. Remove node from its current place
    const node = removeNodeById(workspaceTree.value, nodeId)
    if (!node) return

    // 2. Add to new place
    if (targetCategoryId === null) {
      workspaceTree.value.push(node)
    } else {
      const targetCategory = findNodeById(workspaceTree.value, targetCategoryId)
      if (targetCategory && targetCategory.type === 'category') {
        if (!targetCategory.children) {
          targetCategory.children = []
        }
        targetCategory.children.push(node)
      } else {
        // Fallback: put back at root if category not found
        workspaceTree.value.push(node)
      }
    }
    await saveWorkspace()
  }

  // Fetch the active database catalog (for Ctrl+K palette)
  async function fetchCatalog(): Promise<CatalogItem[]> {
    const bindings = useConnectionsStore().getWailsBindings()
    try {
      const result = await (bindings as any).GetSearchCatalog()
      console.log('[fetchCatalog] result:', result, 'length:', result?.length)
      return result || []
    } catch (e: any) {
      console.error('[fetchCatalog] Failed to fetch search catalog:', e)
      return []
    }
  }

  return {
    workspaceTree,
    loading,
    error,
    loadWorkspace,
    saveWorkspace,
    addCategory,
    addObject,
    deleteNode,
    moveNode,
    fetchCatalog,
  }
})
