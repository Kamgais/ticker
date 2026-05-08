<script setup lang="ts">
import { ref, watch } from 'vue'
import { useTickerStore } from '@/stores/ticker'
import type { TickerEntry } from '@/types/ticker'

const props = defineProps<{
  isOpen: boolean
  entryToEdit?: TickerEntry | null
}>()

const emit = defineEmits<{
  close: []
}>()

const store = useTickerStore()
const loading = ref(false)

const title = ref('')
const message = ref('')
const creator = ref('')
const highlight = ref(false)

watch(() => props.isOpen, (open) => {
  if (open) {
    if (props.entryToEdit) {
      title.value = props.entryToEdit.title
      message.value = props.entryToEdit.message
      creator.value = props.entryToEdit.creator
      highlight.value = props.entryToEdit.highlight
    } else {
      title.value = ''
      message.value = ''
      creator.value = ''
      highlight.value = false
    }
  }
})

async function handleSubmit() {
  if (!title.value.trim() || !message.value.trim() || !creator.value.trim()) return

  loading.value = true
  try {
    if (props.entryToEdit) {
      await store.updateEntry({
        ticker_id: props.entryToEdit.ticker_id,
        title: title.value.trim(),
        message: message.value.trim(),
        creator: creator.value.trim(),
        highlight: highlight.value,
        type: props.entryToEdit.type,
        active: props.entryToEdit.active,
      })
    } else {
      await store.createEntry({
        ticker_name: 'wpftest',
        title: title.value.trim(),
        message: message.value.trim(),
        creator: creator.value.trim(),
        highlight: highlight.value,
        type: 1,
      })
    }
    emit('close')
  } catch {
    // Fehler wird im Store als errorMessage gesetzt
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <Teleport to="body">
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
      @click.self="emit('close')"
    >
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-md mx-4 p-6">
        <div class="flex items-center justify-between mb-6">
          <h2 class="text-xl font-bold text-gray-800 dark:text-white">
            {{ entryToEdit ? 'Eintrag bearbeiten' : 'Neuer Eintrag' }}
          </h2>
          <button
            @click="emit('close')"
            class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 text-2xl leading-none"
          >
            ✕
          </button>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Titel *
            </label>
            <input
              v-model="title"
              type="text"
              placeholder="Titel eingeben..."
              class="w-full px-3 py-2 rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Nachricht *
            </label>
            <textarea
              v-model="message"
              rows="3"
              placeholder="Nachricht eingeben..."
              class="w-full px-3 py-2 rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Ersteller *
            </label>
            <input
              v-model="creator"
              type="text"
              placeholder="Dein Name..."
              class="w-full px-3 py-2 rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div class="flex items-center gap-3">
            <input
              v-model="highlight"
              type="checkbox"
              id="highlight"
              class="w-4 h-4 accent-blue-500"
            />
            <label for="highlight" class="text-sm text-gray-700 dark:text-gray-300">
              ⭐ Als Highlight markieren
            </label>
          </div>

          <div class="flex gap-3 pt-2">
            <button
              type="button"
              @click="emit('close')"
              class="flex-1 px-4 py-2 rounded border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700"
            >
              Abbrechen
            </button>
            <button
              type="submit"
              :disabled="loading || !title.trim() || !message.trim() || !creator.trim()"
              class="flex-1 px-4 py-2 rounded bg-blue-500 hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed text-white font-medium"
            >
              {{ loading ? 'Speichern...' : entryToEdit ? 'Aktualisieren' : 'Erstellen' }}
            </button>
          </div>

        </form>
      </div>
    </div>
  </Teleport>
</template>