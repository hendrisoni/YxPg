<template>
  <div class="h-full flex flex-col bg-navy-primary text-text-primary overflow-hidden font-ui">
    <!-- Top Header Bar -->
    <header class="flex items-center justify-between px-4 py-2 bg-navy-secondary border-b border-navy-border flex-shrink-0">
      <div class="flex items-center gap-3">
        <div class="w-7 h-7 rounded-lg bg-teal-accent/15 border border-teal-accent/30 flex items-center justify-center text-teal-accent">
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
          </svg>
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h1 class="text-sm font-bold text-text-primary tracking-wide">Performance Analyzer</h1>
            <!-- Performance Badge -->
            <span
              v-if="activeRecord"
              class="px-2 py-0.5 text-[10px] font-bold rounded-full border shadow-sm uppercase tracking-wider"
              :class="badgeClass"
            >
              {{ performanceBadge }}
            </span>
          </div>
          <p class="text-[11px] text-text-muted">EXPLAIN ANALYZE performance tracking & execution plan comparison</p>
        </div>
      </div>

      <!-- Action buttons & status info -->
      <div class="flex items-center gap-2">
        <button
          @click="clearHistory"
          class="px-2.5 py-1 text-xs bg-navy-tertiary border border-navy-border hover:bg-rose-500/10 hover:text-rose-400 hover:border-rose-500/30 text-text-secondary rounded transition-colors cursor-pointer"
        >
          Clear History
        </button>
      </div>
    </header>

    <!-- Main Content Area split into History Sidebar + Dashboard -->
    <div class="flex-1 flex overflow-hidden">
      <!-- Left Sidebar: History List (2-Query Selection via Checkboxes) -->
      <div class="w-64 bg-navy-secondary border-r border-navy-border flex flex-col flex-shrink-0 overflow-hidden">
        <!-- Sidebar Header info -->
        <div class="p-3 border-b border-navy-border flex items-center justify-between">
          <div class="text-xs font-bold text-text-muted uppercase tracking-wider">
            History (Pilih 2 Query)
          </div>
          <span class="text-[10px] text-teal-accent font-mono">
            {{ (activeRecord ? 1 : 0) + (compareRecord ? 1 : 0) }}/2 selected
          </span>
        </div>

        <!-- History Records List -->
        <div class="flex-1 overflow-y-auto p-2 space-y-1.5">
          <div v-if="perfStore.records.length === 0" class="text-center py-8 text-xs text-text-muted">
            Belum ada riwayat EXPLAIN.<br />Jalankan EXPLAIN ANALYZE pada Query View.
          </div>

          <div
            v-for="rec in perfStore.records"
            :key="rec.id"
            class="p-2 rounded border transition-all flex items-start gap-2 text-xs group"
            :class="[
              rec.id === activeRecord?.id
                ? 'bg-teal-accent/10 border-teal-accent/60 text-text-primary'
                : rec.id === compareRecord?.id
                ? 'bg-amber-500/10 border-amber-500/50 text-text-primary'
                : 'bg-navy-tertiary/40 border-navy-border/60 hover:bg-navy-hover text-text-secondary'
            ]"
          >
            <!-- Checkbox for 2-Query Comparison -->
            <input
              type="checkbox"
              :checked="rec.id === activeRecord?.id || rec.id === compareRecord?.id"
              @change="toggleSelectForCompare(rec.id)"
              class="mt-1 w-3.5 h-3.5 accent-teal-accent cursor-pointer flex-shrink-0"
              title="Centang untuk memilih & membandingkan 2 query"
            />

            <!-- Card Content -->
            <div class="flex-1 min-w-0 cursor-pointer" @click="toggleSelectForCompare(rec.id)">
              <div class="flex items-center justify-between font-mono text-[10px] text-text-muted mb-1 gap-1">
                <!-- Inline Remark Textbox directly replacing database label -->
                <input
                  v-model="rec.remark"
                  type="text"
                  placeholder="Isi remark..."
                  @click.stop
                  @input="perfStore.updateRemark(rec.id, rec.remark || '')"
                  class="flex-1 min-w-0 bg-navy-primary border border-navy-border/80 focus:border-teal-accent focus:bg-navy-tertiary px-1.5 py-0.5 rounded text-[10px] text-teal-accent placeholder-text-muted/60 focus:outline-none transition-colors font-sans"
                />
                <span class="flex-shrink-0 text-text-muted font-mono text-[10px] pl-1">{{ formatTime(rec.executed_at) }}</span>
              </div>
              <div class="font-semibold text-text-primary truncate mb-1" :title="rec.sql">
                {{ rec.query_name }}
              </div>
              <div class="flex items-center justify-between text-[10px]">
                <span class="font-mono text-teal-accent font-medium">{{ rec.execution_time.toFixed(2) }}ms</span>
                <span class="text-text-muted font-mono">Cost: {{ rec.total_cost.toFixed(0) }}</span>
              </div>
            </div>

            <!-- Delete Button -->
            <button
              @click.stop="perfStore.deleteRecord(rec.id)"
              class="text-text-muted hover:text-rose-400 p-1 rounded hover:bg-rose-500/10 transition-colors flex-shrink-0 cursor-pointer"
              title="Hapus riwayat eksekusi ini"
            >
              <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6" />
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Center & Right: Dashboard Body -->
      <div v-if="activeRecord" class="flex-1 flex overflow-hidden">
        <!-- Main Scrollable Dashboard -->
        <div class="flex-1 overflow-y-auto p-4 space-y-5">
          <!-- Query Information -->
          <section class="bg-navy-secondary border border-navy-border rounded-lg p-3">
            <div class="flex items-center justify-between mb-2">
              <div class="flex items-center gap-2 text-xs font-semibold text-text-muted uppercase tracking-wider">
                <svg class="w-3.5 h-3.5 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" /><polyline points="14 2 14 8 20 8" />
                </svg>
                Query Information
              </div>
              <button
                @click="showFullSql = !showFullSql"
                class="text-[10px] text-teal-accent hover:underline cursor-pointer"
              >
                {{ showFullSql ? 'Hide SQL' : 'View Full SQL' }}
              </button>
            </div>

            <div class="grid grid-cols-2 md:grid-cols-6 gap-2 text-xs mb-2">
              <div class="bg-navy-primary p-2 rounded border border-navy-border">
                <div class="text-[10px] text-text-muted">Database</div>
                <div class="font-mono text-text-primary font-medium truncate">{{ activeRecord.database }}</div>
              </div>
              <div class="bg-navy-primary p-2 rounded border border-navy-border">
                <div class="text-[10px] text-text-muted">Schema</div>
                <div class="font-mono text-text-primary font-medium truncate">{{ activeRecord.schema }}</div>
              </div>
              <div class="bg-navy-primary p-2 rounded border border-navy-border flex flex-col justify-between">
                <div class="text-[10px] text-text-muted">Remark / Catatan</div>
                <input
                  v-model="activeRecord.remark"
                  type="text"
                  placeholder="Isi remark..."
                  @input="perfStore.updateRemark(activeRecord.id, activeRecord.remark || '')"
                  class="w-full bg-navy-secondary border border-navy-border/60 focus:border-teal-accent px-1.5 py-0.5 rounded text-xs text-teal-accent placeholder-text-muted/60 focus:outline-none transition-colors font-sans mt-0.5"
                />
              </div>
              <div class="bg-navy-primary p-2 rounded border border-navy-border">
                <div class="text-[10px] text-text-muted">Query Hash</div>
                <div class="font-mono text-teal-accent font-medium truncate" :title="activeRecord.query_hash">
                  {{ activeRecord.query_hash }}
                </div>
              </div>
              <div class="bg-navy-primary p-2 rounded border border-navy-border">
                <div class="text-[10px] text-text-muted">SQL Length</div>
                <div class="font-mono text-text-primary font-medium">{{ activeRecord.sql_length }} chars</div>
              </div>
              <div class="bg-navy-primary p-2 rounded border border-navy-border">
                <div class="text-[10px] text-text-muted">Timestamp</div>
                <div class="font-mono text-text-primary font-medium text-[11px] truncate">
                  {{ new Date(activeRecord.executed_at).toLocaleString() }}
                </div>
              </div>
            </div>

            <!-- SQL Text snippet or expanded view -->
            <div v-if="showFullSql" class="mt-2 p-2 bg-navy-primary border border-navy-border rounded font-mono text-xs text-teal-300 max-h-40 overflow-y-auto whitespace-pre-wrap">
              {{ activeRecord.sql }}
            </div>
          </section>

          <!-- Summary Cards -->
          <section>
            <h2 class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-2">Performa Latest Execution</h2>
            <div class="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-8 gap-2">
              <!-- Execution Time -->
              <div class="bg-navy-secondary border border-navy-border p-2.5 rounded-lg">
                <div class="text-[10px] text-text-muted">Execution Time</div>
                <div class="text-sm font-bold font-mono text-teal-accent my-0.5">
                  {{ activeRecord.execution_time.toFixed(2) }}<span class="text-[10px] font-normal text-text-muted">ms</span>
                </div>
                <span class="inline-flex items-center gap-1 text-[9px] px-1.5 py-0.2 rounded font-medium" :class="getDeltaStatusBadgeClass('Execution Time')">
                  {{ getDeltaText('Execution Time') }}
                </span>
              </div>

              <!-- Planning Time -->
              <div class="bg-navy-secondary border border-navy-border p-2.5 rounded-lg">
                <div class="text-[10px] text-text-muted">Planning Time</div>
                <div class="text-sm font-bold font-mono text-text-primary my-0.5">
                  {{ activeRecord.planning_time.toFixed(2) }}<span class="text-[10px] font-normal text-text-muted">ms</span>
                </div>
                <span class="inline-flex items-center gap-1 text-[9px] px-1.5 py-0.2 rounded font-medium" :class="getDeltaStatusBadgeClass('Planning Time')">
                  {{ getDeltaText('Planning Time') }}
                </span>
              </div>

              <!-- Total Time -->
              <div class="bg-navy-secondary border border-navy-border p-2.5 rounded-lg">
                <div class="text-[10px] text-text-muted">Total Time</div>
                <div class="text-sm font-bold font-mono text-text-primary my-0.5">
                  {{ activeRecord.total_time.toFixed(2) }}<span class="text-[10px] font-normal text-text-muted">ms</span>
                </div>
                <span class="inline-flex items-center gap-1 text-[9px] px-1.5 py-0.2 rounded font-medium" :class="getDeltaStatusBadgeClass('Total Time')">
                  {{ getDeltaText('Total Time') }}
                </span>
              </div>

              <!-- Rows -->
              <div class="bg-navy-secondary border border-navy-border p-2.5 rounded-lg">
                <div class="text-[10px] text-text-muted">Actual Rows</div>
                <div class="text-sm font-bold font-mono text-text-primary my-0.5">
                  {{ activeRecord.actual_rows }}
                </div>
                <span class="inline-flex items-center gap-1 text-[9px] px-1.5 py-0.2 rounded font-medium" :class="getDeltaStatusBadgeClass('Actual Rows')">
                  {{ getDeltaText('Actual Rows') }}
                </span>
              </div>

              <!-- Total Cost -->
              <div class="bg-navy-secondary border border-navy-border p-2.5 rounded-lg">
                <div class="text-[10px] text-text-muted">Total Cost</div>
                <div class="text-sm font-bold font-mono text-amber-400 my-0.5">
                  {{ activeRecord.total_cost.toFixed(1) }}
                </div>
                <span class="inline-flex items-center gap-1 text-[9px] px-1.5 py-0.2 rounded font-medium" :class="getDeltaStatusBadgeClass('Total Cost')">
                  {{ getDeltaText('Total Cost') }}
                </span>
              </div>

              <!-- Buffers -->
              <div class="bg-navy-secondary border border-navy-border p-2.5 rounded-lg">
                <div class="text-[10px] text-text-muted">Buffers (Hit/Read)</div>
                <div class="text-sm font-bold font-mono text-purple-300 my-0.5">
                  {{ activeRecord.shared_hit }}/{{ activeRecord.shared_read }}
                </div>
                <span class="inline-flex items-center gap-1 text-[9px] px-1.5 py-0.2 rounded font-medium bg-slate-800 text-slate-400">
                  Shared
                </span>
              </div>

              <!-- Scan Type -->
              <div class="bg-navy-secondary border border-navy-border p-2.5 rounded-lg">
                <div class="text-[10px] text-text-muted">Scan Type</div>
                <div class="text-xs font-bold text-teal-300 truncate my-0.5" :title="activeRecord.scan_types.join(', ') || 'N/A'">
                  {{ activeRecord.scan_types[0] || activeRecord.top_node_type }}
                </div>
                <span class="text-[9px] text-text-muted truncate block">Planner</span>
              </div>

              <!-- Join Type -->
              <div class="bg-navy-secondary border border-navy-border p-2.5 rounded-lg">
                <div class="text-[10px] text-text-muted">Join Type</div>
                <div class="text-xs font-bold text-indigo-300 truncate my-0.5" :title="activeRecord.join_types.join(', ') || 'None'">
                  {{ activeRecord.join_types[0] || 'None' }}
                </div>
                <span class="text-[9px] text-text-muted truncate block">Planner</span>
              </div>
            </div>
          </section>

          <!-- Plan Summary & Quick Insight -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Plan Summary -->
            <div class="bg-navy-secondary border border-navy-border rounded-lg p-3">
              <h3 class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-2 flex items-center gap-1.5">
                <svg class="w-3.5 h-3.5 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10" /><line x1="12" y1="16" x2="12" y2="12" /><line x1="12" y1="8" x2="12.01" y2="8" />
                </svg>
                Plan Summary & Ringkasan Perubahan
              </h3>
              <ul class="space-y-1.5 text-xs text-text-primary">
                <li v-for="(sum, idx) in planSummaries" :key="idx" class="flex items-start gap-2">
                  <span class="text-teal-accent font-bold">•</span>
                  <span>{{ sum }}</span>
                </li>
              </ul>
            </div>

            <!-- Quick Insight -->
            <div class="bg-navy-secondary border border-navy-border rounded-lg p-3">
              <h3 class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-2 flex items-center gap-1.5">
                <svg class="w-3.5 h-3.5 text-amber-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M8.5 14.5A2.5 2.5 0 0 0 11 12c0-1.38-.5-2-1-3-1.072-2.143-.224-4.054 2-6 .5 2.5 2 4.9 4 6.5 2 1.6 3 3.5 3 5.5a7 7 0 1 1-14 0c0-1.153.433-2.294 1-3a2.5 2.5 0 0 0 2.5 2.5z" />
                </svg>
                Quick Insights (Maks 5 Poin)
              </h3>
              <div class="space-y-1.5 text-xs">
                <div v-for="(ins, idx) in quickInsights" :key="idx" class="flex items-start gap-2">
                  <span :class="ins.type === 'positive' ? 'text-emerald-400' : ins.type === 'warning' ? 'text-amber-400' : 'text-slate-400'">
                    {{ ins.text }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Trend Chart -->
          <section v-if="trendPoints.length > 1" class="bg-navy-secondary border border-navy-border rounded-lg p-3">
            <h2 class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-2 flex items-center gap-1.5">
              <svg class="w-3.5 h-3.5 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="23 6 13.5 15.5 8.5 10.5 1 18" /><polyline points="17 6 23 6 23 12" />
              </svg>
              Trend Performa Executions (Execution Time vs Cost)
            </h2>
            <div class="h-32 w-full bg-navy-primary rounded p-2 border border-navy-border relative flex items-center justify-center">
              <svg class="w-full h-full" viewBox="0 0 500 100" preserveAspectRatio="none">
                <!-- Grid Lines -->
                <line x1="0" y1="25" x2="500" y2="25" stroke="#202020" stroke-dasharray="2" />
                <line x1="0" y1="50" x2="500" y2="50" stroke="#202020" stroke-dasharray="2" />
                <line x1="0" y1="75" x2="500" y2="75" stroke="#202020" stroke-dasharray="2" />

                <!-- Execution Time Line (Teal) -->
                <polyline
                  fill="none"
                  stroke="#00C9A7"
                  stroke-width="2.5"
                  :points="trendSvgPointsExec"
                />

                <!-- Cost Line (Amber) -->
                <polyline
                  fill="none"
                  stroke="#F59E0B"
                  stroke-width="1.5"
                  stroke-dasharray="3,3"
                  :points="trendSvgPointsCost"
                />

                <!-- Points -->
                <circle
                  v-for="(pt, i) in trendPoints"
                  :key="i"
                  :cx="pt.x"
                  :cy="pt.yExec"
                  r="3.5"
                  fill="#00C9A7"
                  class="hover:r-5 transition-all cursor-pointer"
                >
                  <title>{{ pt.label }}: {{ pt.execTime.toFixed(2) }}ms</title>
                </circle>
              </svg>
            </div>
            <div class="flex items-center justify-between text-[10px] text-text-muted mt-1 px-1">
              <div class="flex items-center gap-3">
                <span class="flex items-center gap-1"><span class="w-2.5 h-0.5 bg-teal-accent"></span> Execution Time (ms)</span>
                <span class="flex items-center gap-1"><span class="w-2.5 h-0.5 bg-amber-400 border-dashed"></span> Cost</span>
              </div>
              <span>{{ trendPoints.length }} executions recorded</span>
            </div>
          </section>

          <!-- Performance Comparison Table -->
          <section class="bg-navy-secondary border border-navy-border rounded-lg p-3">
            <div class="flex items-center justify-between mb-2">
              <h2 class="text-xs font-semibold text-text-muted uppercase tracking-wider flex items-center gap-1.5">
                <svg class="w-3.5 h-3.5 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="20" x2="18" y2="10" /><line x1="12" y1="20" x2="12" y2="4" /><line x1="6" y1="20" x2="6" y2="14" />
                </svg>
                Performance Comparison Table
              </h2>
              <span v-if="compareRecord" class="text-[10px] text-text-muted">
                Comparing vs <strong class="text-teal-accent">{{ formatTime(compareRecord.executed_at) }}</strong>
              </span>
            </div>

            <div class="overflow-x-auto">
              <table class="w-full text-left text-xs border-collapse">
                <thead>
                  <tr class="border-b border-navy-border text-text-muted text-[10px] uppercase font-mono bg-navy-primary">
                    <th class="py-1.5 px-2">Metric</th>
                    <th class="py-1.5 px-2">Previous</th>
                    <th class="py-1.5 px-2">Current</th>
                    <th class="py-1.5 px-2">Difference</th>
                    <th class="py-1.5 px-2">Percentage</th>
                    <th class="py-1.5 px-2">Status</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-navy-border/60">
                  <tr v-for="m in metricComparisons" :key="m.name" class="hover:bg-navy-hover/50">
                    <td class="py-1.5 px-2 font-medium text-text-primary">{{ m.name }}</td>
                    <td class="py-1.5 px-2 font-mono text-text-muted">
                      {{ compareRecord ? `${m.previous} ${m.unit}` : '-' }}
                    </td>
                    <td class="py-1.5 px-2 font-mono font-bold text-text-primary">
                      {{ m.current }} <span class="text-[10px] text-text-muted font-normal">{{ m.unit }}</span>
                    </td>
                    <td class="py-1.5 px-2 font-mono" :class="m.status === 'better' ? 'text-emerald-400' : m.status === 'worse' ? 'text-rose-400' : 'text-slate-400'">
                      {{ compareRecord ? (m.diff > 0 ? `+${m.diff}` : m.diff) : '-' }}
                    </td>
                    <td class="py-1.5 px-2 font-mono" :class="m.status === 'better' ? 'text-emerald-400' : m.status === 'worse' ? 'text-rose-400' : 'text-slate-400'">
                      {{ compareRecord ? `${m.percentage}%` : '-' }}
                    </td>
                    <td class="py-1.5 px-2">
                      <span
                        v-if="compareRecord"
                        class="px-1.5 py-0.2 rounded text-[9px] font-bold uppercase"
                        :class="m.status === 'better' ? 'bg-emerald-500/20 text-emerald-300' : m.status === 'worse' ? 'bg-rose-500/20 text-rose-300' : 'bg-slate-800 text-slate-400'"
                      >
                        {{ m.status === 'better' ? 'Membaik' : m.status === 'worse' ? 'Memburuk' : 'Tetap' }}
                      </span>
                      <span v-else class="text-text-muted text-[10px]">First Run</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </section>

          <!-- Execution Plan Comparison -->
          <section class="bg-navy-secondary border border-navy-border rounded-lg p-3">
            <h2 class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-2 flex items-center gap-1.5">
              <svg class="w-3.5 h-3.5 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="18" height="18" rx="2" /><line x1="12" y1="3" x2="12" y2="21" />
              </svg>
              Execution Plan Side-by-Side Comparison
            </h2>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
              <!-- Previous Plan -->
              <div class="bg-navy-primary p-2 rounded border border-navy-border flex flex-col h-80">
                <div class="text-[11px] font-bold text-text-muted border-b border-navy-border pb-1 mb-2 flex items-center justify-between">
                  <span>Previous Plan</span>
                  <span v-if="compareRecord" class="font-mono text-teal-accent text-[10px]">{{ compareRecord.execution_time.toFixed(2) }}ms</span>
                </div>
                <div class="flex-1 overflow-y-auto pr-1">
                  <PlanNodeTree
                    v-if="comparedTrees.prevTreeWithDiff"
                    :node="comparedTrees.prevTreeWithDiff"
                    :selected-node-id="selectedNode?.id"
                    @select-node="selectedNode = $event"
                  />
                  <div v-else class="h-full flex items-center justify-center text-xs text-text-muted italic">
                    Belum ada riwayat pembanding sebelumnya.
                  </div>
                </div>
              </div>

              <!-- Current Plan -->
              <div class="bg-navy-primary p-2 rounded border border-navy-border flex flex-col h-80">
                <div class="text-[11px] font-bold text-teal-accent border-b border-navy-border pb-1 mb-2 flex items-center justify-between">
                  <span>Current Plan (Latest)</span>
                  <span class="font-mono text-[10px]">{{ activeRecord.execution_time.toFixed(2) }}ms</span>
                </div>
                <div class="flex-1 overflow-y-auto pr-1">
                  <PlanNodeTree
                    v-if="comparedTrees.currTreeWithDiff"
                    :node="comparedTrees.currTreeWithDiff"
                    :selected-node-id="selectedNode?.id"
                    @select-node="selectedNode = $event"
                  />
                </div>
              </div>
            </div>
          </section>

          <!-- Timeline -->
          <section class="bg-navy-secondary border border-navy-border rounded-lg p-3">
            <h2 class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-2 flex items-center gap-1.5">
              <svg class="w-3.5 h-3.5 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10" /><polyline points="12 6 12 12 16 14" />
              </svg>
              Execution History Timeline (Klik untuk memilih Baseline Pembanding)
            </h2>

            <div class="overflow-x-auto">
              <table class="w-full text-left text-xs border-collapse">
                <thead>
                  <tr class="border-b border-navy-border text-text-muted text-[10px] uppercase font-mono bg-navy-primary">
                    <th class="py-1.5 px-2">Timestamp</th>
                    <th class="py-1.5 px-2">Exec Time</th>
                    <th class="py-1.5 px-2">Plan Time</th>
                    <th class="py-1.5 px-2">Rows</th>
                    <th class="py-1.5 px-2">Cost</th>
                    <th class="py-1.5 px-2">Role</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-navy-border/60 font-mono">
                  <tr
                    v-for="item in perfStore.activeQueryHistory"
                    :key="item.id"
                    @click="perfStore.setCompareRecord(item.id)"
                    class="cursor-pointer hover:bg-navy-hover transition-colors"
                    :class="[
                      item.id === activeRecord.id ? 'bg-teal-accent/10 font-bold' : item.id === compareRecord?.id ? 'bg-amber-500/10' : ''
                    ]"
                  >
                    <td class="py-1.5 px-2 text-text-primary">{{ new Date(item.executed_at).toLocaleString() }}</td>
                    <td class="py-1.5 px-2 text-teal-accent">{{ item.execution_time.toFixed(2) }}ms</td>
                    <td class="py-1.5 px-2 text-text-secondary">{{ item.planning_time.toFixed(2) }}ms</td>
                    <td class="py-1.5 px-2 text-text-primary">{{ item.actual_rows }}</td>
                    <td class="py-1.5 px-2 text-amber-400">{{ item.total_cost.toFixed(1) }}</td>
                    <td class="py-1.5 px-2">
                      <span v-if="item.id === activeRecord.id" class="px-1.5 py-0.2 rounded text-[9px] bg-teal-accent text-navy-primary font-bold">
                        ACTIVE CURRENT
                      </span>
                      <span v-else-if="item.id === compareRecord?.id" class="px-1.5 py-0.2 rounded text-[9px] bg-amber-400 text-navy-primary font-bold">
                        BASELINE COMPARE
                      </span>
                      <span v-else class="text-[9px] text-text-muted">Click to Compare</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </section>
        </div>

        <!-- Detail Execution Plan Panel -->
        <div v-if="selectedNode" class="w-72 bg-navy-secondary border-l border-navy-border p-3 flex flex-col flex-shrink-0 overflow-y-auto">
          <div class="flex items-center justify-between border-b border-navy-border pb-2 mb-3">
            <div class="font-bold text-xs text-text-primary">Node Details</div>
            <button @click="selectedNode = null" class="text-text-muted hover:text-white text-xs">✕</button>
          </div>

          <div class="space-y-2 text-xs">
            <div class="bg-navy-primary p-2 rounded border border-navy-border">
              <div class="text-[10px] text-text-muted">Node Type</div>
              <div class="font-bold text-teal-accent font-mono">{{ selectedNode.node_type }}</div>
            </div>

            <div v-if="selectedNode.relation_name" class="bg-navy-primary p-2 rounded border border-navy-border">
              <div class="text-[10px] text-text-muted">Relation Name</div>
              <div class="font-semibold text-text-primary font-mono">{{ selectedNode.relation_name }}</div>
            </div>

            <div v-if="selectedNode.index_name" class="bg-navy-primary p-2 rounded border border-navy-border">
              <div class="text-[10px] text-text-muted">Index Name</div>
              <div class="font-semibold text-amber-300 font-mono">{{ selectedNode.index_name }}</div>
            </div>

            <div class="grid grid-cols-2 gap-1.5">
              <div class="bg-navy-primary p-2 rounded border border-navy-border">
                <div class="text-[10px] text-text-muted">Startup Cost</div>
                <div class="font-mono text-text-primary">{{ selectedNode.startup_cost.toFixed(1) }}</div>
              </div>
              <div class="bg-navy-primary p-2 rounded border border-navy-border">
                <div class="text-[10px] text-text-muted">Total Cost</div>
                <div class="font-mono text-amber-400 font-bold">{{ selectedNode.total_cost.toFixed(1) }}</div>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-1.5">
              <div class="bg-navy-primary p-2 rounded border border-navy-border">
                <div class="text-[10px] text-text-muted">Actual Time</div>
                <div class="font-mono text-teal-accent font-bold">
                  {{ selectedNode.actual_time !== undefined ? `${selectedNode.actual_time.toFixed(2)}ms` : '-' }}
                </div>
              </div>
              <div class="bg-navy-primary p-2 rounded border border-navy-border">
                <div class="text-[10px] text-text-muted">Loops</div>
                <div class="font-mono text-purple-300">{{ selectedNode.actual_loops || 1 }}</div>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-1.5">
              <div class="bg-navy-primary p-2 rounded border border-navy-border">
                <div class="text-[10px] text-text-muted">Plan Rows</div>
                <div class="font-mono text-text-primary">{{ selectedNode.plan_rows }}</div>
              </div>
              <div class="bg-navy-primary p-2 rounded border border-navy-border">
                <div class="text-[10px] text-text-muted">Actual Rows</div>
                <div class="font-mono text-text-primary font-bold">{{ selectedNode.actual_rows !== undefined ? selectedNode.actual_rows : '-' }}</div>
              </div>
            </div>

            <!-- Buffers -->
            <div class="bg-navy-primary p-2 rounded border border-navy-border">
              <div class="text-[10px] text-text-muted mb-1">Shared & Temp Buffers</div>
              <div class="grid grid-cols-2 gap-1 text-[11px] font-mono text-text-secondary">
                <div>Hit: <span class="text-text-primary">{{ selectedNode.shared_hit_blocks || 0 }}</span></div>
                <div>Read: <span class="text-text-primary">{{ selectedNode.shared_read_blocks || 0 }}</span></div>
                <div>Temp R: <span class="text-text-primary">{{ selectedNode.temp_read_blocks || 0 }}</span></div>
                <div>Temp W: <span class="text-text-primary">{{ selectedNode.temp_written_blocks || 0 }}</span></div>
              </div>
            </div>

            <!-- Filter -->
            <div v-if="selectedNode.filter" class="bg-navy-primary p-2 rounded border border-navy-border">
              <div class="text-[10px] text-text-muted">Filter</div>
              <div class="font-mono text-rose-300 text-[11px] break-words">{{ selectedNode.filter }}</div>
            </div>

            <!-- Index Cond -->
            <div v-if="selectedNode.index_cond" class="bg-navy-primary p-2 rounded border border-navy-border">
              <div class="text-[10px] text-text-muted">Index Condition</div>
              <div class="font-mono text-emerald-300 text-[11px] break-words">{{ selectedNode.index_cond }}</div>
            </div>

            <!-- Recheck Cond -->
            <div v-if="selectedNode.recheck_cond" class="bg-navy-primary p-2 rounded border border-navy-border">
              <div class="text-[10px] text-text-muted">Recheck Condition</div>
              <div class="font-mono text-teal-300 text-[11px] break-words">{{ selectedNode.recheck_cond }}</div>
            </div>

            <!-- Workers -->
            <div v-if="selectedNode.workers_planned" class="bg-navy-primary p-2 rounded border border-navy-border">
              <div class="text-[10px] text-text-muted">Workers</div>
              <div class="font-mono text-text-primary">
                Planned: {{ selectedNode.workers_planned }} / Launched: {{ selectedNode.workers_launched || 0 }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State when no active record selected -->
      <div v-else class="flex-1 flex flex-col items-center justify-center text-text-muted p-8 text-center">
        <svg class="w-12 h-12 text-navy-border mb-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
        </svg>
        <div class="text-sm font-semibold text-text-primary mb-1">Belum Ada Performance Record</div>
        <p class="text-xs max-w-sm">
          Jalankan <strong>EXPLAIN ANALYZE</strong> pada Query View untuk merekam dan menganalisis performa query secara otomatis.
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { usePerformanceStore } from '../stores/performance'
import type { ParsedExplainNode } from '../types'
import {
  compareExplainTrees,
  buildMetricComparisons,
  generatePlanSummary,
  generateQuickInsights,
  calculatePerformanceBadge
} from '../utils/explainParser'

import PlanNodeTree from '../components/performance/PlanNodeTree.vue'

const perfStore = usePerformanceStore()

const showFullSql = ref(false)
const selectedNode = ref<ParsedExplainNode | null>(null)

const activeRecord = computed(() => perfStore.activeRecord)
const compareRecord = computed(() => perfStore.compareRecord)

const performanceBadge = computed(() => {
  if (!activeRecord.value) return 'Average'
  return calculatePerformanceBadge(activeRecord.value, compareRecord.value)
})

const badgeClass = computed(() => {
  switch (performanceBadge.value) {
    case 'Excellent':
      return 'bg-emerald-500/20 text-emerald-300 border-emerald-500/40'
    case 'Good':
      return 'bg-teal-500/20 text-teal-300 border-teal-500/40'
    case 'Average':
      return 'bg-slate-700/50 text-slate-300 border-slate-600/50'
    case 'Needs Optimization':
      return 'bg-amber-500/20 text-amber-300 border-amber-500/40'
    case 'Critical':
      return 'bg-rose-500/20 text-rose-300 border-rose-500/40'
    default:
      return 'bg-slate-700/50 text-slate-300 border-slate-600/50'
  }
})

const metricComparisons = computed(() => {
  if (!activeRecord.value) return []
  return buildMetricComparisons(activeRecord.value, compareRecord.value)
})

const planSummaries = computed(() => {
  if (!activeRecord.value) return []
  return generatePlanSummary(activeRecord.value, compareRecord.value)
})

const quickInsights = computed(() => {
  if (!activeRecord.value) return []
  return generateQuickInsights(activeRecord.value, compareRecord.value)
})

const comparedTrees = computed(() => {
  if (!activeRecord.value) return {}
  return compareExplainTrees(compareRecord.value?.plan_tree, activeRecord.value.plan_tree)
})

// Trend Chart SVG Points
const trendPoints = computed(() => {
  const history = perfStore.activeQueryHistory.slice().reverse()
  if (history.length === 0) return []

  const maxExec = Math.max(...history.map(h => h.execution_time), 1)
  const maxCost = Math.max(...history.map(h => h.total_cost), 1)

  return history.map((h, idx) => {
    const x = history.length > 1 ? (idx / (history.length - 1)) * 480 + 10 : 250
    const yExec = 90 - (h.execution_time / maxExec) * 75
    const yCost = 90 - (h.total_cost / maxCost) * 75
    return {
      x,
      yExec,
      yCost,
      execTime: h.execution_time,
      cost: h.total_cost,
      label: new Date(h.executed_at).toLocaleTimeString()
    }
  })
})

const trendSvgPointsExec = computed(() => {
  return trendPoints.value.map(p => `${p.x},${p.yExec}`).join(' ')
})

const trendSvgPointsCost = computed(() => {
  return trendPoints.value.map(p => `${p.x},${p.yCost}`).join(' ')
})

function formatTime(isoStr: string) {
  try {
    const d = new Date(isoStr)
    return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  } catch (e) {
    return isoStr
  }
}

function getDeltaStatusBadgeClass(metricName: string) {
  const m = metricComparisons.value.find(x => x.name === metricName)
  if (!m || !compareRecord.value) return 'bg-slate-800 text-slate-400'
  if (m.status === 'better') return 'bg-emerald-500/20 text-emerald-300'
  if (m.status === 'worse') return 'bg-rose-500/20 text-rose-300'
  return 'bg-slate-800 text-slate-400'
}

function getDeltaText(metricName: string) {
  const m = metricComparisons.value.find(x => x.name === metricName)
  if (!m || !compareRecord.value) return 'Unchanged'
  if (m.status === 'better') return `▲ Membaik (-${m.percentage}%)`
  if (m.status === 'worse') return `▼ Memburuk (+${m.percentage}%)`
  return 'Tetap'
}

function toggleSelectForCompare(id: string) {
  const isCurrentActive = perfStore.activeRecordId === id
  const isCurrentCompare = perfStore.compareRecordId === id

  if (isCurrentActive) {
    perfStore.activeRecordId = perfStore.compareRecordId
    perfStore.compareRecordId = null
  } else if (isCurrentCompare) {
    perfStore.compareRecordId = null
  } else {
    if (!perfStore.activeRecordId) {
      perfStore.activeRecordId = id
    } else if (!perfStore.compareRecordId) {
      perfStore.compareRecordId = id
    } else {
      perfStore.compareRecordId = perfStore.activeRecordId
      perfStore.activeRecordId = id
    }
  }
}

function clearHistory() {
  if (confirm('Clear all EXPLAIN performance history?')) {
    perfStore.clearHistory()
  }
}
</script>
