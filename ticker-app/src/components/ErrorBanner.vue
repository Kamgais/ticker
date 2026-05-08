<script setup lang="ts">
import { useTickerStore } from '@/stores/ticker'

const store = useTickerStore()
</script>

<template>
  <div
    v-if="store.errorMessage"
    class="bg-red-100 dark:bg-red-900 border border-red-300 dark:border-red-700 text-red-700 dark:text-red-200 px-4 py-3 rounded flex justify-between items-center"
  >
    <span>⚠️ {{ store.errorMessage }}</span>
    <button
      @click="store.errorMessage = null"
      class="text-red-500 hover:text-red-700 font-bold text-lg leading-none"
    >
      ✕
    </button>
  </div>
</template>