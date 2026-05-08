<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  login: [username: string]
}>()

const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

// Simulierte Benutzer
const FAKE_USERS = [
  { username: 'admin', password: 'admin123' },
  { username: 'cyril', password: 'cyril123' },
]

async function handleLogin() {
  if (!username.value.trim() || !password.value.trim()) {
    error.value = 'Bitte alle Felder ausfüllen.'
    return
  }

  loading.value = true
  error.value = ''

  // Simulierte Verzögerung wie echte API
  await new Promise(resolve => setTimeout(resolve, 800))

  const user = FAKE_USERS.find(
    u => u.username === username.value && u.password === password.value
  )

  if (user) {
    localStorage.setItem('auth_user', username.value)
    emit('login', username.value)
  } else {
    error.value = 'Falscher Benutzername oder Passwort.'
  }

  loading.value = false
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center px-4">
    <div class="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-sm p-8">

      <!-- Logo -->
      <div class="text-center mb-8">
        <p class="text-5xl mb-3">📰</p>
        <h1 class="text-2xl font-bold text-gray-800 dark:text-white">Ticker App</h1>
        <p class="text-sm text-gray-400 mt-1">Bitte melde dich an</p>
      </div>

      <!-- Hinweis -->
      <div class="bg-blue-50 dark:bg-blue-900/30 border border-blue-200 dark:border-blue-700 rounded-lg p-3 mb-6 text-xs text-blue-600 dark:text-blue-300">
        <p class="font-medium mb-1">Demo-Zugangsdaten:</p>
        <p>👤 admin / 🔑 admin123</p>
        <p>👤 cyril / 🔑 cyril123</p>
      </div>

      <!-- Formular -->
      <form @submit.prevent="handleLogin" class="space-y-4">

        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Benutzername
          </label>
          <input
            v-model="username"
            type="text"
            placeholder="Benutzername eingeben..."
            class="w-full px-3 py-2 rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Passwort
          </label>
          <input
            v-model="password"
            type="password"
            placeholder="Passwort eingeben..."
            class="w-full px-3 py-2 rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <!-- Fehler -->
        <p v-if="error" class="text-sm text-red-500">
          ⚠️ {{ error }}
        </p>

        <button
          type="submit"
          :disabled="loading"
          class="w-full px-4 py-2 bg-blue-500 hover:bg-blue-600 disabled:opacity-50 text-white rounded-lg font-medium"
        >
          {{ loading ? 'Anmelden...' : 'Anmelden' }}
        </button>

      </form>
    </div>
  </div>
</template>