<script setup lang="ts">
import { useTickerStore } from '@/stores/ticker'
import TickerCard from '@/components/TickerCard.vue'
import Pagination from '@/components/Pagination.vue'

const store = useTickerStore()
</script>

<template>
  <div class="space-y-4">

    <!-- Keine Ergebnisse -->
    <div
      v-if="store.paginatedEntries.length === 0"
      class="text-center text-gray-400 dark:text-gray-500 py-12"
    >
      <p class="text-4xl mb-2">📭</p>
      <p>Keine Einträge gefunden.</p>
    </div>

    <!-- Liste -->
    <div v-else class="space-y-3">
      <TickerCard
        v-for="entry in store.paginatedEntries"
        :key="entry.ticker_id"
        :entry="entry"
      />
    </div>

    <!-- Pagination -->
    <Pagination />

  </div>
</template>