<template>
  <div class="h-screen w-screen overflow-hidden flex">
    <!-- Sidebar Navigation -->
    <aside class="w-64 bg-gradient-to-b from-blue-50 to-blue-100 px-6 py-8 flex flex-col">
      <div class="text-left mb-6">
        <h2 class="text-xl font-bold text-gray-800">Language Portal</h2>
        <p class="text-xs text-gray-600 mt-1">Learning Tracker</p>
      </div>

      <nav class="flex-grow overflow-auto">
        <ul class="space-y-2">
          <li v-for="(link, index) in navLinks" :key="index">
            <NuxtLink 
              :to="link.to"
              class="flex items-center p-3 rounded-xl transition-all duration-300 group text-gray-700 hover:bg-blue-200 hover:text-blue-600"
            >
              <span class="text-lg mr-3">{{ link.icon }}</span>
              <span class="text-sm flex-grow">{{ link.label }}</span>
            </NuxtLink>
          </li>
        </ul>
      </nav>

      <div class="mt-4 pt-4 border-t border-gray-200 relative">
        <div 
          @click="toggleUserPopup"
          class="flex items-center cursor-pointer hover:bg-blue-100 p-2 rounded-xl transition-colors"
        >
          <div class="w-10 h-10 bg-blue-400 text-white rounded-full flex items-center justify-center mr-3">
            {{ getUserInitial }}
          </div>
          <div>
            <p class="text-sm font-semibold">{{ userName || 'Set Name' }}</p>
            <p class="text-xs text-gray-500">{{ userName ? 'Edit Profile' : 'Click to set name' }}</p>
          </div>
        </div>

        <!-- User Popup -->
        <div 
          v-if="showUserPopup"
          class="absolute bottom-full left-0 right-0 mb-2 bg-white rounded-xl shadow-lg border border-gray-200 p-4 z-50"
        >
          <div class="mb-4">
            <label class="block text-xs text-gray-600 mb-1">Your Name</label>
            <input 
              v-model="editedUserName"
              type="text" 
              placeholder="Enter your name"
              class="w-full px-3 py-2 text-sm border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
          </div>
          <div class="space-y-2">
            <button 
              @click="saveUserName"
              class="w-full bg-blue-500 text-white py-2 rounded-lg text-sm hover:bg-blue-600 transition-colors"
            >
              Save Name
            </button>
            <button 
              @click="showResetConfirmation = true; toggleUserPopup()"
              class="w-full bg-red-500 text-white py-2 rounded-lg text-sm hover:bg-red-600 transition-colors"
            >
              Reset / Clear
            </button>
            <button 
              @click="toggleUserPopup"
              class="w-full bg-gray-200 text-gray-700 py-2 rounded-lg text-sm hover:bg-gray-300 transition-colors"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </aside>

    <!-- Main Content -->
    <slot />

    <!-- Reset Confirmation Modal -->
    <div 
      v-if="showResetConfirmation" 
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
      <div class="bg-white rounded-2xl p-6 max-w-md w-full mx-4 shadow-2xl">
        <div class="text-center">
          <div class="mx-auto mb-4 w-12 h-12 bg-red-100 rounded-full flex items-center justify-center">
            <svg 
              class="w-6 h-6 text-red-600" 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24" 
              xmlns="http://www.w3.org/2000/svg"
            >
              <path 
                stroke-linecap="round" 
                stroke-linejoin="round" 
                stroke-width="2" 
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              ></path>
            </svg>
          </div>
          <h2 class="text-lg font-bold text-gray-800 mb-3">Reset All Activity?</h2>
          <p class="text-xs text-gray-600 mb-4">
            This action will permanently clear all your learning progress, 
            recent sessions, and saved activities. This cannot be undone.
          </p>
          <div class="flex justify-center space-x-3">
            <button 
              @click="showResetConfirmation = false"
              class="px-4 py-2 bg-gray-200 text-gray-800 rounded-lg text-xs hover:bg-gray-300 transition-colors"
            >
              Cancel
            </button>
            <button 
              @click="resetAllActivity"
              class="px-4 py-2 bg-red-500 text-white rounded-lg text-xs hover:bg-red-600 transition-colors"
            >
              Reset Anyway
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

// Navigation Links
const navLinks = [
  { to: '/words', icon: '📖', label: 'Words' },
  { to: '/groups', icon: '🏷️', label: 'Groups' },
  { to: '/study-activities', icon: '📚', label: 'Study Activities' },
  { to: '/sessions', icon: '⏱️', label: 'Sessions' }
]

const showResetConfirmation = ref(false)
const showUserPopup = ref(false)
const editedUserName = ref('')
const userName = ref('')

// Get user initial
const getUserInitial = computed(() => {
  if (userName.value) {
    return userName.value.charAt(0).toUpperCase()
  }
  return 'L'
})

// Toggle user popup
function toggleUserPopup() {
  showUserPopup.value = !showUserPopup.value
  if (showUserPopup.value) {
    editedUserName.value = userName.value || ''
  }
}

// Save user name
function saveUserName() {
  if (editedUserName.value.trim()) {
    userName.value = editedUserName.value.trim()
    toggleUserPopup()
  }
}

// Reset function to clear all data
function resetAllActivity() {
  // Implement reset logic here
  showResetConfirmation.value = false
}
</script>
