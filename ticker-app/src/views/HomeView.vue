<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useTickerStore } from '@/stores/ticker'
import TickerList from '@/components/TickerList.vue'
import SearchBar from '@/components/SearchBar.vue'
import ErrorBanner from '@/components/ErrorBanner.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import TickerModal from '@/components/TickerModal.vue'
import DarkModeToggle from '@/components/DarkModeToggle.vue'

const store = useTickerStore()
const isModalOpen = ref(false)

const REFRESH_INTERVAL = 10
const countdown = ref(REFRESH_INTERVAL)

let refreshInterval: ReturnType<typeof setInterval>
let countdownInterval: ReturnType<typeof setInterval>

onMounted(() => {
  store.fetchEntries()
  refreshInterval = setInterval(() => {
    store.fetchEntries()
    countdown.value = REFRESH_INTERVAL
  }, REFRESH_INTERVAL * 1000)
  countdownInterval = setInterval(() => {
    countdown.value--
  }, 1000)
})

onUnmounted(() => {
  clearInterval(refreshInterval)
  clearInterval(countdownInterval)
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors">

    <header class="bg-white dark:bg-gray-800 shadow px-6 py-4 flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-800 dark:text-white">
        📰 Ticker
        <span class="text-sm font-normal text-gray-400 ml-2">wpftest</span>
      </h1>

      <div class="flex items-center gap-4">

        <div class="flex items-center gap-2 text-sm text-gray-400 dark:text-gray-500">
          <svg
            v-if="store.status === 'loading'"
            class="w-4 h-4 animate-spin text-blue-500"
            fill="none" viewBox="0 0 24 24"
          >
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
          </svg>
          <span v-else class="w-2 h-2 rounded-full bg-green-400 inline-block" />
          <span v-if="store.status === 'loading'">Aktualisiere...</span>
          <span v-else>Refresh in {{ countdown }}s</span>
        </div>
        <DarkModeToggle />
        <button
          @click="isModalOpen = true"
          class="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg font-medium text-sm"
        >
          + Neuer Eintrag
        </button>

      </div>
    </header>


    <main class="max-w-3xl mx-auto px-4 py-8 space-y-4">
      <ErrorBanner />
      <SearchBar />
      <LoadingSpinner v-if="store.status === 'loading' && store.entries.length === 0" />
      <TickerList v-else />
    </main>

    <TickerModal
      :is-open="isModalOpen"
      @close="isModalOpen = false"
    />

  </div>
</template>