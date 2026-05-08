<script setup lang="ts">
import { useTickerStore } from '@/stores/ticker'
import type { SortField, SortDirection } from '@/types/ticker'

const store = useTickerStore()

function handleSortField(event: Event) {
  const field = (event.target as HTMLSelectElement).value as SortField
  store.setSort({ field, direction: store.sortConfig.direction })
}

function toggleSortDirection() {
  const direction: SortDirection =
    store.sortConfig.direction === 'asc' ? 'desc' : 'asc'
  store.setSort({ field: store.sortConfig.field, direction })
}
</script>

<template>
  <div class="flex flex-col sm:flex-row gap-2">

    <!-- Suchfeld -->
    <div class="relative flex-1">
      <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400">🔍</span>
      <input
        :value="store.searchQuery"
        @input="store.setSearch(($event.target as HTMLInputElement).value)"
        type="text"
        placeholder="Suchen nach Titel, Nachricht, Ersteller..."
        class="w-full pl-9 pr-4 py-2 rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      <button
        v-if="store.searchQuery"
        @click="store.setSearch('')"
        class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
      >
        ✕
      </button>
    </div>

    <!-- Sortierung -->
    <div class="flex gap-2">

      <!-- Sortierfeld -->
      <select
        :value="store.sortConfig.field"
        @change="handleSortField"
        class="px-3 py-2 rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
      >
        <option value="created">📅 Datum</option>
        <option value="title">🔤 Titel</option>
        <option value="ticker_id">🔢 ID</option>
      </select>

      <!-- Richtung -->
      <button
        @click="toggleSortDirection"
        class="px-3 py-2 rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 text-sm"
        :title="store.sortConfig.direction === 'asc' ? 'Aufsteigend' : 'Absteigend'"
      >
        {{ store.sortConfig.direction === 'asc' ? '⬆️ Aufsteigend' : '⬇️ Absteigend' }}
      </button>

    </div>
  </div>
</template>