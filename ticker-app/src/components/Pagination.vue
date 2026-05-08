<script setup lang="ts">
import { useTickerStore } from '@/stores/ticker'

const store = useTickerStore()
</script>

<template>
  <div
    v-if="store.totalPages > 1"
    class="flex items-center justify-between pt-2"
  >

    <!-- Info -->
    <span class="text-sm text-gray-400 dark:text-gray-500">
      Seite {{ store.currentPage }} von {{ store.totalPages }}
      · {{ store.filteredEntries.length }} Einträge
    </span>

    <!-- Buttons -->
    <div class="flex gap-1">

      <!-- Erste Seite -->
      <button
        @click="store.setPage(1)"
        :disabled="store.currentPage === 1"
        class="px-2 py-1 text-sm rounded border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300 disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-700"
      >
        «
      </button>

      <!-- Zurück -->
      <button
        @click="store.setPage(store.currentPage - 1)"
        :disabled="store.currentPage === 1"
        class="px-3 py-1 text-sm rounded border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300 disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-700"
      >
        ‹
      </button>

      <!-- Seitenzahlen -->
      <button
        v-for="page in store.totalPages"
        :key="page"
        @click="store.setPage(page)"
        :class="[
          'px-3 py-1 text-sm rounded border ',
          page === store.currentPage
            ? 'bg-blue-500 border-blue-500 text-white'
            : 'border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700'
        ]"
      >
        {{ page }}
      </button>

      <!-- Weiter -->
      <button
        @click="store.setPage(store.currentPage + 1)"
        :disabled="store.currentPage === store.totalPages"
        class="px-3 py-1 text-sm rounded border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300 disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-700"
      >
        ›
      </button>

      <!-- Letzte Seite -->
      <button
        @click="store.setPage(store.totalPages)"
        :disabled="store.currentPage === store.totalPages"
        class="px-2 py-1 text-sm rounded border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300 disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-700"
      >
        »
      </button>

    </div>
  </div>
</template>