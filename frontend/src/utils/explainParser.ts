import type {
  ParsedExplainNode,
  PerformanceRecord,
  MetricComparison,
  QuickInsightItem,
  PerformanceBadgeType
} from '../types'

let nodeIdCounter = 0

/**
 * Normalizes PostgreSQL JSON Explain result into a structured ParsedExplainNode tree
 */
export function parsePostgresExplainJson(jsonStr: string): {
  planTree?: ParsedExplainNode
  executionTime: number
  planningTime: number
  totalCost: number
  startupCost: number
  planRows: number
  actualRows: number
  sharedHit: number
  sharedRead: number
  tempRead: number
  tempWritten: number
  scanTypes: string[]
  joinTypes: string[]
  topNodeType: string
} {
  try {
    let raw = JSON.parse(jsonStr)
    if (Array.isArray(raw) && raw.length > 0) {
      raw = raw[0]
    }

    const planningTime = typeof raw['Planning Time'] === 'number' ? raw['Planning Time'] : 0
    const executionTime = typeof raw['Execution Time'] === 'number' ? raw['Execution Time'] : 0
    const planObj = raw['Plan'] || raw

    const scanTypesSet = new Set<string>()
    const joinTypesSet = new Set<string>()

    let totalSharedHit = 0
    let totalSharedRead = 0
    let totalTempRead = 0
    let totalTempWritten = 0

    function parseNode(obj: any): ParsedExplainNode {
      nodeIdCounter++
      const nodeType = obj['Node Type'] || 'Result'
      const relName = obj['Relation Name']
      const alias = obj['Alias']
      const indexName = obj['Index Name']
      const joinType = obj['Join Type']
      const scanDir = obj['Scan Direction']

      if (nodeType.toLowerCase().includes('scan')) {
        scanTypesSet.add(nodeType)
      }
      if (nodeType.toLowerCase().includes('join') || joinType) {
        joinTypesSet.add(joinType ? `${joinType} Join` : nodeType)
      }

      const sHit = obj['Shared Hit Blocks'] || 0
      const sRead = obj['Shared Read Blocks'] || 0
      const tRead = obj['Temp Read Blocks'] || 0
      const tWrite = obj['Temp Written Blocks'] || 0

      totalSharedHit += sHit
      totalSharedRead += sRead
      totalTempRead += tRead
      totalTempWritten += tWrite

      const childPlans = Array.isArray(obj['Plans']) ? obj['Plans'].map((c: any) => parseNode(c)) : []

      return {
        id: `node_${nodeIdCounter}_${Date.now()}`,
        node_type: nodeType,
        relation_name: relName,
        alias: alias,
        index_name: indexName,
        scan_direction: scanDir,
        join_type: joinType,
        startup_cost: obj['Startup Cost'] || 0,
        total_cost: obj['Total Cost'] || 0,
        plan_rows: obj['Plan Rows'] || 0,
        plan_width: obj['Plan Width'] || 0,
        actual_time: obj['Actual Total Time'] !== undefined ? obj['Actual Total Time'] : obj['Actual Startup Time'],
        actual_rows: obj['Actual Rows'] || 0,
        actual_loops: obj['Actual Loops'] || 1,
        shared_hit_blocks: sHit,
        shared_read_blocks: sRead,
        shared_written_blocks: obj['Shared Written Blocks'] || 0,
        temp_read_blocks: tRead,
        temp_written_blocks: tWrite,
        filter: obj['Filter'],
        index_cond: obj['Index Cond'],
        recheck_cond: obj['Recheck Cond'],
        rows_removed_by_filter: obj['Rows Removed by Filter'],
        heap_fetches: obj['Heap Fetches'],
        workers_planned: obj['Workers Planned'],
        workers_launched: obj['Workers Launched'],
        plans: childPlans
      }
    }

    const planTree = planObj ? parseNode(planObj) : undefined

    return {
      planTree,
      executionTime,
      planningTime,
      totalCost: planTree ? planTree.total_cost : 0,
      startupCost: planTree ? planTree.startup_cost : 0,
      planRows: planTree ? planTree.plan_rows : 0,
      actualRows: planTree?.actual_rows ?? 0,
      sharedHit: totalSharedHit,
      sharedRead: totalSharedRead,
      tempRead: totalTempRead,
      tempWritten: totalTempWritten,
      scanTypes: Array.from(scanTypesSet),
      joinTypes: Array.from(joinTypesSet),
      topNodeType: planTree ? planTree.node_type : 'Unknown'
    }
  } catch (err) {
    console.error('Failed to parse EXPLAIN JSON:', err)
    return {
      executionTime: 0,
      planningTime: 0,
      totalCost: 0,
      startupCost: 0,
      planRows: 0,
      actualRows: 0,
      sharedHit: 0,
      sharedRead: 0,
      tempRead: 0,
      tempWritten: 0,
      scanTypes: [],
      joinTypes: [],
      topNodeType: 'Unknown'
    }
  }
}

/**
 * Compare two execution trees and mark node diffs
 */
export function compareExplainTrees(
  prevTree?: ParsedExplainNode,
  currTree?: ParsedExplainNode
): { prevTreeWithDiff?: ParsedExplainNode; currTreeWithDiff?: ParsedExplainNode } {
  if (!prevTree || !currTree) {
    return { prevTreeWithDiff: prevTree, currTreeWithDiff: currTree }
  }

  function cloneAndDiff(
    currNode: ParsedExplainNode,
    prevNode?: ParsedExplainNode
  ): ParsedExplainNode {
    const isNodeChanged = !prevNode || currNode.node_type !== prevNode.node_type || currNode.index_name !== prevNode.index_name
    const diffStatus = isNodeChanged ? 'changed' : 'same'
    const prevType = prevNode ? (prevNode.index_name ? `${prevNode.node_type} (${prevNode.index_name})` : prevNode.node_type) : undefined

    const currPlans = currNode.plans || []
    const prevPlans = prevNode?.plans || []

    const diffedChildren = currPlans.map((cChild, idx) => cloneAndDiff(cChild, prevPlans[idx]))

    return {
      ...currNode,
      diff_status: diffStatus,
      prev_node_type: prevType,
      plans: diffedChildren
    }
  }

  const currTreeWithDiff = cloneAndDiff(currTree, prevTree)
  const prevTreeWithDiff = cloneAndDiff(prevTree, currTree)

  return { prevTreeWithDiff, currTreeWithDiff }
}

/**
 * Compare two performance records to compute metric comparisons
 */
export function buildMetricComparisons(current: PerformanceRecord, previous?: PerformanceRecord): MetricComparison[] {
  if (!previous) {
    return [
      { name: 'Execution Time', unit: 'ms', previous: 0, current: current.execution_time, diff: 0, percentage: 0, status: 'same' },
      { name: 'Planning Time', unit: 'ms', previous: 0, current: current.planning_time, diff: 0, percentage: 0, status: 'same' },
      { name: 'Total Cost', unit: '', previous: 0, current: current.total_cost, diff: 0, percentage: 0, status: 'same' },
      { name: 'Actual Rows', unit: 'rows', previous: 0, current: current.actual_rows, diff: 0, percentage: 0, status: 'same' },
      { name: 'Estimated Rows', unit: 'rows', previous: 0, current: current.plan_rows, diff: 0, percentage: 0, status: 'same' },
      { name: 'Shared Hit', unit: 'blocks', previous: 0, current: current.shared_hit, diff: 0, percentage: 0, status: 'same' },
      { name: 'Shared Read', unit: 'blocks', previous: 0, current: current.shared_read, diff: 0, percentage: 0, status: 'same' },
      { name: 'Temp Read', unit: 'blocks', previous: 0, current: current.temp_read, diff: 0, percentage: 0, status: 'same' },
      { name: 'Temp Written', unit: 'blocks', previous: 0, current: current.temp_written, diff: 0, percentage: 0, status: 'same' }
    ]
  }

  function calcMetric(name: string, unit: string, curVal: number, prevVal: number, lowerIsBetter = true): MetricComparison {
    const diff = curVal - prevVal
    let percentage = 0
    if (prevVal !== 0) {
      percentage = Number(((diff / prevVal) * 100).toFixed(1))
    }
    let status: 'better' | 'worse' | 'same' = 'same'
    if (Math.abs(diff) > 0.001) {
      if (lowerIsBetter) {
        status = diff < 0 ? 'better' : 'worse'
      } else {
        status = diff > 0 ? 'better' : 'worse'
      }
    }

    return {
      name,
      unit,
      previous: prevVal,
      current: curVal,
      diff: Number(diff.toFixed(2)),
      percentage: Math.abs(percentage),
      status
    }
  }

  return [
    calcMetric('Execution Time', 'ms', current.execution_time, previous.execution_time, true),
    calcMetric('Planning Time', 'ms', current.planning_time, previous.planning_time, true),
    calcMetric('Total Cost', '', current.total_cost, previous.total_cost, true),
    calcMetric('Actual Rows', 'rows', current.actual_rows, previous.actual_rows, true),
    calcMetric('Estimated Rows', 'rows', current.plan_rows, previous.plan_rows, true),
    calcMetric('Shared Hit', 'blocks', current.shared_hit, previous.shared_hit, false),
    calcMetric('Shared Read', 'blocks', current.shared_read, previous.shared_read, true),
    calcMetric('Temp Read', 'blocks', current.temp_read, previous.temp_read, true),
    calcMetric('Temp Written', 'blocks', current.temp_written, previous.temp_written, true)
  ]
}

/**
 * Generate automatic plan summary sentences
 */
export function generatePlanSummary(current: PerformanceRecord, previous?: PerformanceRecord): string[] {
  if (!previous) {
    return [`Eksekusi pertama untuk query ini (${current.execution_time.toFixed(2)}ms, Cost: ${current.total_cost.toFixed(1)}).`]
  }

  const summaries: string[] = []

  // Execution Time
  const execDiff = current.execution_time - previous.execution_time
  if (Math.abs(execDiff) > 0.05) {
    const pct = previous.execution_time > 0 ? Math.abs((execDiff / previous.execution_time) * 100).toFixed(0) : '0'
    if (execDiff < 0) {
      summaries.push(`Execution Time turun ${pct}%.`)
    } else {
      summaries.push(`Execution Time meningkat ${pct}%.`)
    }
  } else {
    summaries.push(`Execution Time relatif stabil.`)
  }

  // Scan type change
  const prevScans = previous.scan_types.join(', ') || previous.top_node_type
  const currScans = current.scan_types.join(', ') || current.top_node_type
  if (prevScans && currScans && prevScans !== currScans) {
    summaries.push(`Planner berubah dari ${prevScans} menjadi ${currScans}.`)
  }

  // Cost change
  const costDiff = current.total_cost - previous.total_cost
  if (Math.abs(costDiff) > 0.5) {
    const pct = previous.total_cost > 0 ? Math.abs((costDiff / previous.total_cost) * 100).toFixed(0) : '0'
    if (costDiff < 0) {
      summaries.push(`Estimated Cost turun ${pct}%.`)
    } else {
      summaries.push(`Estimated Cost naik ${pct}%.`)
    }
  }

  // Shared Read change
  const readDiff = current.shared_read - previous.shared_read
  if (readDiff !== 0) {
    const pct = previous.shared_read > 0 ? Math.abs((readDiff / previous.shared_read) * 100).toFixed(0) : '0'
    if (readDiff < 0) {
      summaries.push(`Shared Read berkurang ${pct}%.`)
    } else {
      summaries.push(`Shared Read bertambah ${pct}%.`)
    }
  }

  // Planning Time
  const planDiff = current.planning_time - previous.planning_time
  if (Math.abs(planDiff) > 0.05) {
    const pct = previous.planning_time > 0 ? Math.abs((planDiff / previous.planning_time) * 100).toFixed(0) : '0'
    if (planDiff < 0) {
      summaries.push(`Planning Time turun ${pct}%.`)
    } else {
      summaries.push(`Planning Time meningkat ${pct}%.`)
    }
  }

  // Rows
  if (current.actual_rows === previous.actual_rows) {
    summaries.push(`Rows tetap sama (${current.actual_rows} rows).`)
  } else {
    summaries.push(`Jumlah Rows berubah dari ${previous.actual_rows} menjadi ${current.actual_rows}.`)
  }

  return summaries
}

/**
 * Generate max 5 quick insights with icons
 */
export function generateQuickInsights(current: PerformanceRecord, previous?: PerformanceRecord): QuickInsightItem[] {
  const insights: QuickInsightItem[] = []

  // 1. Scan check
  const hasIndexScan = current.scan_types.some(s => s.toLowerCase().includes('index'))
  const hasSeqScan = current.scan_types.some(s => s.toLowerCase().includes('seq'))
  if (hasIndexScan && !hasSeqScan) {
    insights.push({ type: 'positive', text: '✔ Menggunakan Index Scan secara efisien.' })
  } else if (hasSeqScan) {
    insights.push({ type: 'warning', text: '⚠ Terdapat Sequential Scan pada tabel.' })
  } else {
    insights.push({ type: 'positive', text: '✔ Tipe pencarian optimal.' })
  }

  // 2. Temp file check
  if (current.temp_read === 0 && current.temp_written === 0) {
    insights.push({ type: 'positive', text: '✔ Tidak ada Temp File (RAM mencukupi).' })
  } else {
    insights.push({ type: 'warning', text: `⚠ Menulis temp file disk (${current.temp_written} blocks).` })
  }

  // 3. Row estimation accuracy check
  if (current.plan_rows > 0 && current.actual_rows > 0) {
    const ratio = Math.max(current.plan_rows, current.actual_rows) / Math.min(current.plan_rows, current.actual_rows)
    if (ratio <= 3) {
      insights.push({ type: 'positive', text: '✔ Estimasi Rows akurat.' })
    } else {
      insights.push({ type: 'warning', text: '⚠ Estimasi Rows kurang akurat (pertimbangkan ANALYZE).' })
    }
  }

  // 4. Comparison checks if previous exists
  if (previous) {
    if (current.planning_time > previous.planning_time * 1.25 && (current.planning_time - previous.planning_time) > 0.2) {
      insights.push({ type: 'warning', text: '⚠ Planning Time meningkat.' })
    } else if (current.execution_time < previous.execution_time * 0.8) {
      insights.push({ type: 'positive', text: '✔ Performa waktu eksekusi membaik.' })
    }

    const hasNestedLoop = current.join_types.some(j => j.toLowerCase().includes('nested'))
    if (hasNestedLoop && current.actual_rows > 1000) {
      insights.push({ type: 'warning', text: '⚠ Nested Loop dapat dioptimalkan untuk dataset besar.' })
    }
  }

  return insights.slice(0, 5)
}

/**
 * Calculate Performance Badge (Excellent, Good, Average, Needs Optimization, Critical)
 */
export function calculatePerformanceBadge(current: PerformanceRecord, previous?: PerformanceRecord): PerformanceBadgeType {
  const t = current.execution_time

  if (current.temp_written > 100 || t > 2000) {
    return 'Critical'
  }

  if (previous) {
    const speedup = (previous.execution_time - current.execution_time) / previous.execution_time
    if (speedup >= 0.5 || t < 10) return 'Excellent'
    if (speedup >= 0.2 || t < 100) return 'Good'
    if (speedup <= -0.5 || t > 500) return 'Needs Optimization'
    if (t > 100) return 'Average'
  } else {
    if (t < 10) return 'Excellent'
    if (t < 100) return 'Good'
    if (t < 500) return 'Average'
    if (t <= 2000) return 'Needs Optimization'
  }

  return 'Average'
}
