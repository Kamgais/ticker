<script setup lang="ts">
import { ref, onMounted } from 'vue'

const isDark = ref(false)

onMounted(() => {
  // Gespeicherte Präferenz laden
  isDark.value = localStorage.getItem('darkMode') === 'true'
  applyDarkMode()
})

function toggle() {
  isDark.value = !isDark.value
  localStorage.setItem('darkMode', String(isDark.value))
  applyDarkMode()
}

function applyDarkMode() {
  if (isDark.value) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}
</script>

<template>
  <button
    @click="toggle"
    class="px-3 py-2 rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 text-sm"
    :title="isDark ? 'Light Mode' : 'Dark Mode'"
  >
    {{ isDark ? '☀️ Light' : '🌙 Dark' }}
  </button>
</template>