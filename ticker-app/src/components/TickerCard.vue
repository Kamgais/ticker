<script setup lang="ts">
import { ref } from 'vue'
import { useTickerStore } from '@/stores/ticker'
import type { TickerEntry } from '@/types/ticker'
import TickerModal from '@/components/TickerModal.vue'

const props = defineProps<{ entry: TickerEntry }>()

const store = useTickerStore()
const isEditModalOpen = ref(false)
const isDeleting = ref(false)
const showDeleteConfirm = ref(false)

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleString('de-DE', {
    day: '2-digit', month: '2-digit', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  })
}

async function handleDelete() {
  isDeleting.value = true
  try {
    await store.deleteEntry(props.entry.ticker_id)
  } catch {
    // Fehler im Store
  } finally {
    isDeleting.value = false
    showDeleteConfirm.value = false
  }
}
</script>

<template>
  <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-4 border border-gray-100 dark:border-gray-700">

    <div class="flex items-start justify-between gap-2 mb-2">
      <h2 class="font-bold text-gray-800 dark:text-white text-lg">
        {{ entry.title }}
      </h2>
      <span
        v-if="entry.highlight"
        class="text-xs bg-yellow-100 dark:bg-yellow-900 text-yellow-700 dark:text-yellow-300 px-2 py-0.5 rounded-full shrink-0"
      >
        ⭐ Highlight
      </span>
    </div>

    <p class="text-gray-600 dark:text-gray-300 mb-3">
      {{ entry.message }}
    </p>

    <div class="flex items-center justify-between text-xs text-gray-400 dark:text-gray-500">
      <span>👤 {{ entry.creator }}</span>
      <span>🕐 {{ formatDate(entry.created) }}</span>
    </div>

    <div class="flex gap-2 mt-3 pt-3 border-t border-gray-100 dark:border-gray-700">
      <button
        @click="isEditModalOpen = true"
        class="flex-1 px-3 py-1.5 text-sm rounded border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700"
      >
        ✏️ Bearbeiten
      </button>

      <button
        v-if="!showDeleteConfirm"
        @click="showDeleteConfirm = true"
        class="flex-1 px-3 py-1.5 text-sm rounded border border-red-300 dark:border-red-700 text-red-500 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20"
      >
        🗑️ Löschen
      </button>

      <div v-else class="flex-1 flex gap-1">
        <button
          @click="showDeleteConfirm = false"
          class="flex-1 px-2 py-1.5 text-sm rounded border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300"
        >
          Abbrechen
        </button>
        <button
          @click="handleDelete"
          :disabled="isDeleting"
          class="flex-1 px-2 py-1.5 text-sm rounded bg-red-500 hover:bg-red-600 text-white disabled:opacity-50"
        >
          {{ isDeleting ? '...' : 'Sicher?' }}
        </button>
      </div>

    </div>
  </div>

  <TickerModal
    :is-open="isEditModalOpen"
    :entry-to-edit="entry"
    @close="isEditModalOpen = false"
  />
</template>