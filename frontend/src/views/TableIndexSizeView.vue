<template>
  <div class="h-full flex flex-col bg-navy-primary text-text-primary overflow-hidden">
    <!-- Header / Toolbar -->
    <div class="p-3 bg-navy-secondary border-b border-navy-border flex flex-wrap items-center justify-between gap-3 flex-shrink-0">
      <!-- Left side: Title & Schema Selector -->
      <div class="flex items-center gap-3">
        <div class="flex items-center gap-2">
          <div class="p-1.5 rounded bg-teal-accent/10 text-teal-accent">
            <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="20" x2="18" y2="10" />
              <line x1="12" y1="20" x2="12" y2="4" />
              <line x1="6" y1="20" x2="6" y2="14" />
            </svg>
          </div>
          <div>
            <h1 class="text-sm font-semibold text-text-primary leading-none">Table & Index Size</h1>
            <p class="text-[11px] text-text-muted mt-0.5">Analisis ukuran penyimpanan per tabel dan indeks (Urutan terbesar)</p>
          </div>
        </div>

        <!-- Connection / Database Select -->
        <div class="flex items-center gap-1.5 ml-2">
          <label class="text-xs text-text-secondary">Conn:</label>
          <select
            v-model="selectedConnId"
            @change="onConnChange"
            class="text-xs bg-navy-tertiary border border-navy-border rounded px-2 py-1 text-text-primary focus:border-teal-accent focus:outline-none cursor-pointer"
          >
            <option v-for="conn in connectionsStore.connections" :key="conn.id" :value="conn.id">
              {{ conn.name }} ({{ conn.database }})
            </option>
          </select>
        </div>
      </div>

      <!-- Right side: Color legend, Search input & Refresh button -->
      <div class="flex items-center gap-2">
        <!-- Cached Data Status Indicator -->
        <div v-if="lastUpdatedTime" class="hidden md:flex items-center gap-1.5 text-[10px] bg-navy-tertiary/90 border border-navy-border/80 rounded px-2 py-1" :title="isFromCache ? 'Data dimuat dari penyimpanan tersimpan (' + lastUpdatedTime + '). Klik Refresh untuk menghitung ulang dari database.' : 'Data baru saja dihitung ulang pada ' + lastUpdatedTime">
          <svg class="w-3.5 h-3.5 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/>
            <polyline points="17 21 17 13 7 13 7 21"/>
            <polyline points="7 3 7 8 15 8"/>
          </svg>
          <span class="text-text-secondary font-mono">
            <span class="text-text-muted">Tersimpan:</span> {{ lastUpdatedTime }}
          </span>
          <span v-if="isFromCache" class="text-[9px] px-1 py-0.2 rounded bg-teal-accent/15 text-teal-accent border border-teal-accent/30 font-semibold font-mono">
            CACHED
          </span>
        </div>

        <!-- Color Legend Badge -->
        <div class="hidden sm:flex items-center gap-1.5 text-[10px] bg-navy-tertiary/70 border border-navy-border/70 rounded px-2 py-1 font-mono">
          <span class="text-text-muted font-sans font-medium mr-0.5">Level:</span>
          <span class="text-emerald-400 font-medium" title="Dibawah 100 MB">&lt;100MB</span>
          <span class="text-text-muted">•</span>
          <span class="text-amber-400 font-semibold" title="100 MB - 500 MB">100-500MB</span>
          <span class="text-text-muted">•</span>
          <span class="text-orange-400 font-bold" title="500 MB - 1 GB">500MB-1GB</span>
          <span class="text-text-muted">•</span>
          <span class="text-red-400 font-bold" title="Diatas 1 GB">&gt;1GB</span>
        </div>

        <div class="relative">
          <svg class="absolute left-2.5 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-text-muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8" />
            <path d="M21 21l-4.35-4.35" />
          </svg>
          <input
            v-model="searchFilter"
            type="text"
            placeholder="Cari tabel / indeks..."
            class="pl-8 pr-3 py-1 text-xs bg-navy-tertiary border border-navy-border rounded text-text-primary placeholder-text-muted focus:border-teal-accent focus:outline-none w-48"
          />
        </div>

        <!-- Show / Hide Schema Panel Toggle Button -->
        <button
          @click="showSchemaSummary = !showSchemaSummary"
          class="px-2 py-1 text-xs font-medium rounded border transition-all cursor-pointer flex items-center gap-1.5"
          :class="showSchemaSummary 
            ? 'bg-teal-accent/15 text-teal-300 border-teal-accent/40 shadow-sm' 
            : 'bg-navy-tertiary text-text-secondary border-navy-border hover:text-white hover:bg-navy-hover'"
          :title="showSchemaSummary ? 'Sembunyikan Panel Ringkasan Skema' : 'Tampilkan Panel Ringkasan Skema'"
        >
          <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path v-if="!showSchemaSummary" d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
            <circle v-if="!showSchemaSummary" cx="12" cy="12" r="3"></circle>
            <path v-if="showSchemaSummary" d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path>
            <line v-if="showSchemaSummary" x1="1" y1="1" x2="23" y2="23"></line>
          </svg>
          <span>{{ showSchemaSummary ? 'Sembunyikan Skema' : 'Ringkasan Skema' }}</span>
        </button>

        <button
          @click="loadData(true)"
          :disabled="loading"
          class="flex items-center gap-1 px-2.5 py-1 text-xs bg-teal-accent text-navy-primary rounded font-medium hover:bg-teal-hover transition-colors disabled:opacity-50 cursor-pointer"
          title="Hitung Ulang & Refresh Data dari Database"
        >
          <svg class="w-3.5 h-3.5" :class="{ 'animate-spin': loading }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21.5 2v6h-6M2.5 22v-6h6" />
            <path d="M2 11.5a10 10 0 0 1 18.8-4.3L21.5 8M22 12.5a10 10 0 0 1-18.8 4.3L2.5 16" />
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="flex-1 overflow-auto p-3">
      <!-- Loading State -->
      <div v-if="loading" class="h-full flex items-center justify-center">
        <div class="text-center">
          <div class="inline-block w-8 h-8 border-3 border-teal-accent border-t-transparent rounded-full animate-spin"></div>
          <p class="text-xs text-text-muted mt-2">Memuat data ukuran tabel & indeks...</p>
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="errorMessage" class="h-full flex items-center justify-center p-6 text-center">
        <div class="bg-red-500/10 border border-red-500/30 rounded-lg p-6 max-w-md">
          <svg class="w-10 h-10 text-red-400 mx-auto mb-2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10" />
            <line x1="12" y1="8" x2="12" />
            <line x1="12" y1="16" x2="12.01" y2="16" />
          </svg>
          <h3 class="text-sm font-semibold text-red-400">Gagal Memuat Data</h3>
          <p class="text-xs text-text-muted mt-1">{{ errorMessage }}</p>
          <button @click="loadData" class="mt-3 px-3 py-1 text-xs bg-navy-tertiary hover:bg-navy-hover text-text-primary rounded border border-navy-border">
            Coba Lagi
          </button>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredItems.length === 0" class="h-full flex items-center justify-center text-center p-6">
        <div>
          <svg class="w-10 h-10 text-text-muted mx-auto mb-2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M20 7H4a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2Z" />
            <path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16" />
          </svg>
          <p class="text-xs text-text-muted">Tidak ada data ukuran yang ditemukan.</p>
        </div>
      </div>

      <!-- Main Data Section: Schema Size Cards + Table Grid -->
      <div v-else class="space-y-3">
        <!-- Total Size Per Schema Summary Cards (Toggleable) -->
        <div v-if="showSchemaSummary && schemaSummaries.length > 0" class="bg-navy-secondary border border-navy-border rounded-lg p-2.5 shadow-md transition-all">
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-2">
              <div class="p-1 rounded bg-teal-accent/10 text-teal-accent">
                <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path>
                  <polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline>
                  <line x1="12" y1="22.08" x2="12" y2="12"></line>
                </svg>
              </div>
              <h3 class="text-xs font-bold text-text-primary uppercase tracking-wide flex items-center gap-2">
                <span>Ringkasan Ukuran Penyimpanan Per Schema</span>
                <span class="text-[10px] px-1.5 py-0.2 rounded-full bg-teal-accent/20 text-teal-accent border border-teal-accent/40 font-mono">
                  {{ schemaSummaries.length }} Schema
                </span>
              </h3>
            </div>
            <div class="flex items-center gap-2 text-[10px] text-text-muted">
              <span v-if="selectedSchema" class="text-teal-accent font-semibold flex items-center gap-1">
                Filter: {{ selectedSchema }}
                <button @click="selectedSchema = ''" class="text-rose-400 hover:underline cursor-pointer ml-1">(Reset)</button>
              </span>
              <span>Klik skema untuk memfilter</span>
              <button 
                @click="showSchemaSummary = false"
                class="p-1 rounded text-text-muted hover:text-white hover:bg-navy-tertiary transition-colors cursor-pointer ml-1"
                title="Sembunyikan Panel Ringkasan Skema"
              >
                <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>
          </div>

          <!-- Cards Grid for Each Schema -->
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-2">
            <div
              v-for="s in schemaSummaries"
              :key="s.name"
              @click="toggleSchemaFilter(s.name)"
              class="p-2 rounded-md border transition-all cursor-pointer select-none flex flex-col justify-between group"
              :class="selectedSchema === s.name 
                ? 'bg-teal-accent/15 border-teal-accent text-white shadow-md ring-1 ring-teal-accent/50' 
                : 'bg-navy-tertiary/70 border-navy-border/70 text-text-secondary hover:bg-navy-hover hover:border-slate-600 hover:text-white'"
            >
              <div class="flex items-center justify-between gap-1 mb-1">
                <span class="font-mono text-xs font-bold truncate group-hover:text-teal-accent" :class="selectedSchema === s.name ? 'text-teal-accent' : 'text-text-primary'">
                  {{ s.name }}
                </span>
                <span class="text-[9px] px-1 py-0.2 rounded font-mono" :class="selectedSchema === s.name ? 'bg-teal-accent/30 text-teal-100' : 'bg-navy-primary text-text-muted'">
                  {{ s.tableCount }} tbl
                </span>
              </div>

              <!-- Total Size & Percentage of DB -->
              <div class="mt-0.5 flex items-baseline justify-between gap-1">
                <span class="text-xs font-mono font-black" :class="getSizeColorClass(s.totalBytes)">
                  {{ formatBytes(s.totalBytes) }}
                </span>
                <span class="text-[9px] font-mono text-text-muted">
                  {{ s.percentOfTotal.toFixed(1) }}%
                </span>
              </div>

              <!-- Progress bar showing relative share -->
              <div class="mt-1.5 w-full bg-navy-primary h-1 rounded-full overflow-hidden border border-navy-border/50">
                <div 
                  class="h-full bg-gradient-to-r from-teal-500 to-emerald-400 rounded-full transition-all duration-300"
                  :style="{ width: `${Math.max(3, s.percentOfTotal)}%` }"
                ></div>
              </div>
            </div>
          </div>
        </div>

        <!-- Data Table -->
        <div class="border border-navy-border rounded-lg overflow-hidden bg-navy-secondary shadow-md">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="bg-navy-tertiary/80 text-text-secondary border-b border-navy-border select-none">
              <th class="py-2.5 px-3 w-10 text-center">#</th>
              <th class="py-2.5 px-3 w-10"></th>
              <th @click="toggleSort('schema')" class="py-2.5 px-3 cursor-pointer hover:text-teal-accent transition-colors">
                <div class="flex items-center gap-1">Schema <span v-if="sortBy === 'schema'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span></div>
              </th>
              <th @click="toggleSort('table_name')" class="py-2.5 px-3 cursor-pointer hover:text-teal-accent transition-colors">
                <div class="flex items-center gap-1">Nama Tabel <span v-if="sortBy === 'table_name'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span></div>
              </th>
              <th @click="toggleSort('row_count')" class="py-2.5 px-3 text-right cursor-pointer hover:text-teal-accent transition-colors">
                <div class="flex items-center justify-end gap-1">Baris <span v-if="sortBy === 'row_count'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span></div>
              </th>
              <th @click="toggleSort('dead_tuples')" class="py-2.5 px-3 text-right cursor-pointer hover:text-teal-accent transition-colors">
                <div class="flex items-center justify-end gap-1">Dead Tuples <span v-if="sortBy === 'dead_tuples'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span></div>
              </th>
              <th @click="toggleSort('free_bytes')" class="py-2.5 px-3 text-right cursor-pointer hover:text-teal-accent transition-colors">
                <div class="flex items-center justify-end gap-1">Free Space <span v-if="sortBy === 'free_bytes'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span></div>
              </th>
              <th @click="toggleSort('table_bytes')" class="py-2.5 px-3 text-right cursor-pointer hover:text-teal-accent transition-colors">
                <div class="flex items-center justify-end gap-1">Ukuran Tabel <span v-if="sortBy === 'table_bytes'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span></div>
              </th>
              <th @click="toggleSort('index_bytes')" class="py-2.5 px-3 text-right cursor-pointer hover:text-teal-accent transition-colors">
                <div class="flex items-center justify-end gap-1">Ukuran Indeks <span v-if="sortBy === 'index_bytes'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span></div>
              </th>
              <th @click="toggleSort('total_bytes')" class="py-2.5 px-3 text-right cursor-pointer hover:text-teal-accent transition-colors">
                <div class="flex items-center justify-end gap-1">Total Ukuran <span v-if="sortBy === 'total_bytes'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span></div>
              </th>
              <th class="py-2.5 px-3 w-44">Visual Ukuran</th>
              <th class="py-2.5 px-3 text-right w-24">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-navy-border/60">
            <template v-for="(item, index) in filteredItems" :key="item.schema + '.' + item.table_name">
              <tr class="hover:bg-navy-hover/50 transition-colors">
                <td class="py-2 px-3 text-center text-text-muted font-mono">{{ index + 1 }}</td>
                <td class="py-2 px-3 text-center">
                  <button
                    v-if="item.indexes && item.indexes.length > 0"
                    @click="toggleExpandRow(item.schema + '.' + item.table_name)"
                    class="p-0.5 rounded hover:bg-navy-tertiary text-text-muted hover:text-teal-accent transition-colors cursor-pointer"
                    :title="expandedRows.has(item.schema + '.' + item.table_name) ? 'Tutup Detail Indeks' : 'Lihat Detail Indeks'"
                  >
                    <svg class="w-4 h-4 transition-transform duration-150" :class="{ 'rotate-90 text-teal-accent': expandedRows.has(item.schema + '.' + item.table_name) }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M9 18l6-6-6-6" />
                    </svg>
                  </button>
                </td>
                <td class="py-2 px-3 font-mono text-text-secondary">{{ item.schema }}</td>
                <td class="py-2 px-3 font-semibold text-text-primary">
                  <div class="flex items-center gap-1.5">
                    <span>{{ item.table_name }}</span>
                    <span v-if="item.indexes && item.indexes.length > 0" class="text-[10px] text-text-muted font-normal">({{ item.indexes.length }} idx)</span>
                  </div>
                </td>
                <td class="py-2 px-3 text-right font-mono text-text-secondary">{{ formatNumber(item.row_count) }}</td>
                <td class="py-2 px-3 text-right font-mono">
                  <span v-if="item.dead_tuples > 0" class="text-amber-400 font-semibold" :title="`${formatNumber(item.dead_tuples)} dead tuples`">
                    {{ formatNumber(item.dead_tuples) }}
                  </span>
                  <span v-else class="text-text-muted">0</span>
                </td>
                <td class="py-2 px-3 text-right font-mono font-semibold">
                  <span v-if="item.free_bytes > 0" class="text-orange-400" :title="`Estimasi ruang kosong (bloat): ${formatBytes(item.free_bytes)}`">
                    {{ formatBytes(item.free_bytes) }}
                  </span>
                  <span v-else class="text-text-muted">0 B</span>
                </td>
                <td class="py-2 px-3 text-right font-mono font-semibold" :class="getSizeColorClass(item.table_bytes)">
                  {{ formatBytes(item.table_bytes) }}
                </td>
                <td class="py-2 px-3 text-right font-mono font-semibold" :class="getSizeColorClass(item.index_bytes)">
                  {{ formatBytes(item.index_bytes) }}
                </td>
                <td class="py-2 px-3 text-right font-mono font-bold text-sm" :class="getSizeColorClass(item.total_bytes)">
                  {{ formatBytes(item.total_bytes) }}
                </td>
                <td class="py-2 px-3">
                  <div class="w-full bg-navy-tertiary h-2.5 rounded-full overflow-hidden flex" title="Visualisasi perbandingan ukuran tabel dan indeks">
                    <div
                      class="bg-teal-accent h-full transition-all duration-300"
                      :style="{ width: getPercentage(item.table_bytes, maxTotalBytes) + '%' }"
                      title="Ukuran Tabel"
                    ></div>
                    <div
                      class="bg-amber-400 h-full transition-all duration-300"
                      :style="{ width: getPercentage(item.index_bytes, maxTotalBytes) + '%' }"
                      title="Ukuran Indeks"
                    ></div>
                  </div>
                </td>
                <td class="py-2 px-3 text-right">
                  <div class="flex items-center justify-end gap-1">
                    <button
                      @click="openColumnSizeTab(item.schema, item.table_name)"
                      class="p-1 rounded hover:bg-navy-tertiary text-text-secondary hover:text-teal-accent transition-colors cursor-pointer"
                      title="Lihat Detail Ukuran Per Kolom (Tab Baru)"
                    >
                      <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <line x1="12" y1="20" x2="12" y2="4" />
                        <line x1="6" y1="20" x2="6" y2="14" />
                        <line x1="18" y1="20" x2="18" y2="10" />
                      </svg>
                    </button>
                    <button
                      @click="openTable(item.schema, item.table_name)"
                      class="p-1 rounded hover:bg-navy-tertiary text-text-secondary hover:text-teal-accent transition-colors cursor-pointer"
                      title="Buka Data Tabel"
                    >
                      <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <rect x="3" y="3" width="18" height="18" rx="2" />
                        <path d="M3 9h18M9 21V9" />
                      </svg>
                    </button>
                    <button
                      @click="openQuery(item.schema, item.table_name)"
                      class="p-1 rounded hover:bg-navy-tertiary text-text-secondary hover:text-teal-accent transition-colors cursor-pointer"
                      title="Buat Query SELECT"
                    >
                      <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="4 17 10 11 4 5" />
                        <line x1="12" y1="19" x2="20" y2="19" />
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>

              <!-- Expandable Row for Index Breakdown -->
              <tr v-if="expandedRows.has(item.schema + '.' + item.table_name)" class="bg-navy-tertiary/40">
                <td colspan="10" class="p-3 pl-12 border-b border-navy-border">
                  <div class="bg-navy-secondary border border-navy-border rounded p-3.5">
                    <div class="flex items-center justify-between mb-2">
                      <h4 class="text-xs font-semibold text-text-primary flex items-center gap-1.5">
                        <svg class="w-3.5 h-3.5 text-amber-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <polygon points="12 2 2 7 12 12 22 7 12 2" />
                          <polyline points="2 17 12 22 22 17" />
                          <polyline points="2 12 12 17 22 12" />
                        </svg>
                        Rincian Indeks Tabel: <span class="text-teal-accent font-mono">{{ item.schema }}.{{ item.table_name }}</span>
                      </h4>
                      <span class="text-[11px] text-text-muted">Total: {{ item.indexes ? item.indexes.length : 0 }} Indeks</span>
                    </div>

                    <div v-if="!item.indexes || item.indexes.length === 0" class="text-[11px] text-text-muted italic py-1">
                      Tabel ini tidak memiliki indeks tambahan.
                    </div>
                    <table v-else class="w-full text-left text-[11px] border-collapse font-mono">
                      <thead>
                        <tr class="text-text-muted border-b border-navy-border/80 font-sans">
                          <th class="py-1 px-2">Nama Indeks</th>
                          <th class="py-1 px-2">Tipe Indeks</th>
                          <th class="py-1 px-2 text-right">Ukuran</th>
                          <th class="py-1 px-2 text-right">% Indeks</th>
                        </tr>
                      </thead>
                      <tbody class="divide-y divide-navy-border/40">
                        <tr v-for="idx in item.indexes" :key="idx.index_name" class="hover:bg-navy-hover/30">
                          <td class="py-1.5 px-2 text-text-primary font-semibold">{{ idx.index_name }}</td>
                          <td class="py-1.5 px-2 text-text-secondary uppercase text-[10px] font-sans">{{ idx.index_type || 'btree' }}</td>
                          <td class="py-1.5 px-2 text-right text-amber-400 font-bold">{{ formatBytes(idx.index_bytes) }}</td>
                          <td class="py-1.5 px-2 text-right text-text-muted">
                            {{ item.index_bytes > 0 ? ((idx.index_bytes / item.index_bytes) * 100).toFixed(1) + '%' : '0%' }}
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>
  </div>

    <!-- Summary Bar Footer -->
    <div class="px-4 py-2.5 bg-navy-secondary border-t border-navy-border flex flex-wrap items-center justify-between text-xs text-text-secondary flex-shrink-0">
      <div class="flex items-center gap-4">
        <span>Total Tabel: <strong class="text-text-primary">{{ filteredItems.length }}</strong></span>
        <span>Total Baris: <strong class="text-text-primary">{{ formatNumber(summaryStats.totalRows) }}</strong></span>
        <span v-if="summaryStats.totalDeadTuples > 0" class="text-amber-400">
          Dead Tuples: <strong class="text-amber-400 font-semibold">{{ formatNumber(summaryStats.totalDeadTuples) }}</strong>
        </span>
        <span v-if="summaryStats.totalFreeBytes > 0" class="text-orange-400">
          Free Space: <strong class="text-orange-400 font-semibold">{{ formatBytes(summaryStats.totalFreeBytes) }}</strong>
        </span>
      </div>
      <div class="flex items-center gap-5">
        <span class="flex items-center gap-1">
          <span class="w-2.5 h-2.5 rounded-full bg-teal-accent inline-block"></span>
          Ukuran Tabel: <strong class="text-text-primary">{{ formatBytes(summaryStats.totalTableBytes) }}</strong>
        </span>
        <span class="flex items-center gap-1">
          <span class="w-2.5 h-2.5 rounded-full bg-amber-400 inline-block"></span>
          Ukuran Indeks: <strong class="text-text-primary">{{ formatBytes(summaryStats.totalIndexBytes) }}</strong>
        </span>
        <span class="flex items-center gap-1 text-teal-accent font-semibold">
          Total DB Size: <strong class="text-teal-accent font-bold text-sm">{{ formatBytes(summaryStats.totalBytes) }}</strong>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useConnectionsStore } from '../stores/connections'
import { useSchemaStore } from '../stores/schema'
import { useTabsStore } from '../stores/tabs'
import { useUiStore } from '../stores/ui'
import type { Tab, TableIndexSizeInfo } from '../types'
import * as App from '../../wailsjs/go/main/App'

const props = defineProps<{
  tab: Tab
}>()

const connectionsStore = useConnectionsStore()
const schemaStore = useSchemaStore()
const tabsStore = useTabsStore()
const uiStore = useUiStore()

const selectedConnId = ref<string>(props.tab.connectionId || connectionsStore.currentConnectionId || '')
const selectedSchema = ref<string>('')
const searchFilter = ref<string>('')

const items = ref<TableIndexSizeInfo[]>([])
const loading = ref<boolean>(false)
const errorMessage = ref<string>('')
const lastUpdatedTime = ref<string>('')
const isFromCache = ref<boolean>(false)
const showSchemaSummary = ref<boolean>(true)

// Sorting state for table view: default total_bytes descending (largest size first)
const sortBy = ref<keyof TableIndexSizeInfo>('total_bytes')
const sortDir = ref<'asc' | 'desc'>('desc')

const expandedRows = ref<Set<string>>(new Set())

interface SchemaSummary {
  name: string
  tableCount: number
  totalRows: number
  tableBytes: number
  indexBytes: number
  totalBytes: number
  percentOfTotal: number
}

const schemaSummaries = computed<SchemaSummary[]>(() => {
  const map = new Map<string, SchemaSummary>()
  let overallTotalBytes = 0

  for (const item of items.value) {
    const sName = item.schema || 'public'
    overallTotalBytes += item.total_bytes || 0
    let existing = map.get(sName)
    if (!existing) {
      existing = {
        name: sName,
        tableCount: 0,
        totalRows: 0,
        tableBytes: 0,
        indexBytes: 0,
        totalBytes: 0,
        percentOfTotal: 0
      }
      map.set(sName, existing)
    }
    existing.tableCount += 1
    existing.totalRows += item.row_count || 0
    existing.tableBytes += item.table_bytes || 0
    existing.indexBytes += item.index_bytes || 0
    existing.totalBytes += item.total_bytes || 0
  }

  const list = Array.from(map.values())
  for (const s of list) {
    s.percentOfTotal = overallTotalBytes > 0 ? (s.totalBytes / overallTotalBytes) * 100 : 0
  }

  return list.sort((a, b) => b.totalBytes - a.totalBytes)
})

function toggleSchemaFilter(schemaName: string) {
  if (selectedSchema.value === schemaName) {
    selectedSchema.value = ''
  } else {
    selectedSchema.value = schemaName
  }
}

function getCacheKey(connId: string, schema: string): string {
  return `yxpg:table_index_sizes:${connId}:${schema || 'all'}`
}

function loadFromCache(connId: string, schema: string): boolean {
  try {
    const key = getCacheKey(connId, schema)
    const raw = localStorage.getItem(key)
    if (!raw) return false
    const parsed = JSON.parse(raw)
    if (parsed && Array.isArray(parsed.data) && parsed.data.length > 0) {
      items.value = parsed.data
      if (parsed.timestamp) {
        const d = new Date(parsed.timestamp)
        const timeStr = d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
        const dateStr = d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short' })
        lastUpdatedTime.value = `${timeStr} (${dateStr})`
      } else {
        lastUpdatedTime.value = 'Tersimpan'
      }
      isFromCache.value = true
      return true
    }
  } catch (e) {
    console.warn('Failed to load table index sizes from cache:', e)
  }
  return false
}

function saveToCache(connId: string, schema: string, data: TableIndexSizeInfo[]) {
  try {
    const key = getCacheKey(connId, schema)
    const now = Date.now()
    const payload = {
      timestamp: now,
      data
    }
    localStorage.setItem(key, JSON.stringify(payload))
    const d = new Date(now)
    const timeStr = d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
    const dateStr = d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short' })
    lastUpdatedTime.value = `${timeStr} (${dateStr})`
    isFromCache.value = false
  } catch (e) {
    console.warn('Failed to save table index sizes to cache:', e)
  }
}

const availableSchemas = computed(() => {
  const set = new Set<string>()
  for (const item of items.value) {
    if (item.schema) set.add(item.schema)
  }
  return Array.from(set).sort()
})

const maxTotalBytes = computed(() => {
  let max = 0
  for (const item of items.value) {
    if (item.total_bytes > max) max = item.total_bytes
  }
  return max || 1
})

const filteredItems = computed(() => {
  let res = items.value

  if (selectedSchema.value) {
    res = res.filter(i => i.schema === selectedSchema.value)
  }

  if (searchFilter.value.trim()) {
    const q = searchFilter.value.toLowerCase().trim()
    res = res.filter(i => {
      const matchTable = i.table_name.toLowerCase().includes(q) || i.schema.toLowerCase().includes(q)
      const matchIndex = i.indexes?.some(idx => idx.index_name.toLowerCase().includes(q))
      return matchTable || matchIndex
    })
  }

  // Sort results
  return res.slice().sort((a, b) => {
    const valA = a[sortBy.value]
    const valB = b[sortBy.value]

    if (valA === undefined) return 1
    if (valB === undefined) return -1

    if (typeof valA === 'string' && typeof valB === 'string') {
      const strA = valA.toLowerCase()
      const strB = valB.toLowerCase()
      if (strA < strB) return sortDir.value === 'asc' ? -1 : 1
      if (strA > strB) return sortDir.value === 'asc' ? 1 : -1
      return 0
    }

    const numA = typeof valA === 'number' ? valA : 0
    const numB = typeof valB === 'number' ? valB : 0
    if (numA < numB) return sortDir.value === 'asc' ? -1 : 1
    if (numA > numB) return sortDir.value === 'asc' ? 1 : -1
    return 0
  })
})

const summaryStats = computed(() => {
  let totalRows = 0
  let totalDeadTuples = 0
  let totalFreeBytes = 0
  let totalTableBytes = 0
  let totalIndexBytes = 0
  let totalBytes = 0

  for (const item of filteredItems.value) {
    totalRows += item.row_count || 0
    totalDeadTuples += item.dead_tuples || 0
    totalFreeBytes += item.free_bytes || 0
    totalTableBytes += item.table_bytes || 0
    totalIndexBytes += item.index_bytes || 0
    totalBytes += item.total_bytes || 0
  }

  return {
    totalRows,
    totalDeadTuples,
    totalFreeBytes,
    totalTableBytes,
    totalIndexBytes,
    totalBytes
  }
})

onMounted(() => {
  if (!selectedConnId.value && connectionsStore.connections.length > 0) {
    selectedConnId.value = connectionsStore.connections[0].id
  }
  loadData(false)
})

watch(() => props.tab.connectionId, (newConnId) => {
  if (newConnId && newConnId !== selectedConnId.value) {
    selectedConnId.value = newConnId
    loadData(false)
  }
})

watch(selectedConnId, (newConnId) => {
  if (newConnId && props.tab && props.tab.id && props.tab.connectionId !== newConnId) {
    tabsStore.updateTab(props.tab.id, { connectionId: newConnId })
  }
})

function onConnChange() {
  if (props.tab && props.tab.id && selectedConnId.value) {
    tabsStore.updateTab(props.tab.id, { connectionId: selectedConnId.value })
  }
  loadData(true)
}

async function loadData(forceRefresh = false) {
  if (!selectedConnId.value) {
    errorMessage.value = 'Silakan pilih koneksi database terlebih dahulu.'
    return
  }

  errorMessage.value = ''

  // Try to load from local cache if not forcing refresh
  if (!forceRefresh) {
    const hasCache = loadFromCache(selectedConnId.value, '')
    if (hasCache) {
      return
    }
  }

  loading.value = true

  try {
    if (!connectionsStore.activeConnections.includes(selectedConnId.value)) {
      try {
        await connectionsStore.connect(selectedConnId.value)
      } catch (connErr) {
        console.warn('Store connect warning:', connErr)
      }
    }
    const res = await App.GetTableIndexSizes(selectedConnId.value, '')
    items.value = res || []
    saveToCache(selectedConnId.value, '', items.value)
  } catch (err: any) {
    console.error('Failed to get table index sizes:', err)
    errorMessage.value = err?.message || String(err)
    uiStore.addNotification({
      type: 'error',
      title: 'Gagal Memuat Ukuran',
      message: errorMessage.value
    })
  } finally {
    loading.value = false
  }
}

function toggleSort(field: keyof TableIndexSizeInfo) {
  if (sortBy.value === field) {
    sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortBy.value = field
    sortDir.value = field === 'table_name' || field === 'schema' ? 'asc' : 'desc'
  }
}

function toggleExpandRow(key: string) {
  if (expandedRows.value.has(key)) {
    expandedRows.value.delete(key)
  } else {
    expandedRows.value.add(key)
  }
}

function formatBytes(bytes: number): string {
  if (bytes === 0 || !bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  const num = bytes / Math.pow(k, i)
  return `${num.toFixed(num >= 100 || i === 0 ? 0 : 2)} ${sizes[i]}`
}

function getSizeColorClass(bytes: number): string {
  if (!bytes || bytes <= 0) return 'text-text-muted font-normal'
  const mb = 1024 * 1024
  const gb = 1024 * 1024 * 1024

  if (bytes >= 5 * gb) return 'text-red-500 font-extrabold'
  if (bytes >= 1 * gb) return 'text-red-400 font-bold'
  if (bytes >= 500 * mb) return 'text-orange-400 font-bold'
  if (bytes >= 100 * mb) return 'text-amber-400 font-semibold'
  return 'text-emerald-400 font-medium'
}

function formatNumber(num: number): string {
  return new Intl.NumberFormat().format(num || 0)
}

function getPercentage(part: number, total: number): number {
  if (!total || total <= 0) return 0
  return Math.min(100, Math.max(0, (part / total) * 100))
}

function openColumnSizeTab(schema: string, table: string) {
  if (!selectedConnId.value) return
  tabsStore.createTab('table-column-size', {
    title: `${table} Column Size`,
    connectionId: selectedConnId.value,
    schema,
    table
  })
}

function openTable(schema: string, table: string) {
  if (!selectedConnId.value) return
  tabsStore.createTab('table', {
    title: table,
    connectionId: selectedConnId.value,
    schema,
    table
  })
}

function openQuery(schema: string, table: string) {
  if (!selectedConnId.value) return
  tabsStore.createTab('query', {
    title: table,
    connectionId: selectedConnId.value,
    schema,
    table,
    sql: `SELECT * FROM ${schema}.${table} LIMIT 100`
  })
}
</script>
