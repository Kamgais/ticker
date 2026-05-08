<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterView } from 'vue-router'
import LoginScreen from '@/components/LoginScreen.vue'
import DarkModeToggle from '@/components/DarkModeToggle.vue'

const isLoggedIn = ref(false)
const currentUser = ref('')

onMounted(() => {
  // Gespeicherte Session laden
  const savedUser = localStorage.getItem('auth_user')
  if (savedUser) {
    currentUser.value = savedUser
    isLoggedIn.value = true
  }

  // Dark Mode wiederherstellen
  if (localStorage.getItem('darkMode') === 'true') {
    document.documentElement.classList.add('dark')
  }
})

function handleLogin(username: string) {
  currentUser.value = username
  isLoggedIn.value = true
}

function handleLogout() {
  localStorage.removeItem('auth_user')
  isLoggedIn.value = false
  currentUser.value = ''
}
</script>

<template>
  <!-- Login Screen -->
  <LoginScreen v-if="!isLoggedIn" @login="handleLogin" />

  <!-- App -->
  <div v-else>
    <!-- Logout Bar -->
    <div class="bg-gray-100 dark:bg-gray-950 px-6 py-1.5 flex items-center justify-end gap-3 text-sm">
      <span class="text-gray-500 dark:text-gray-400">
        👤 {{ currentUser }}
      </span>
      <button
        @click="handleLogout"
        class="text-red-400 hover:text-red-600 dark:hover:text-red-300"
      >
        Abmelden
      </button>
    </div>

    <RouterView />
  </div>
</template>