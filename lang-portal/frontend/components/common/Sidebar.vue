<template>
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
  </aside>
</template>

<script setup lang="ts">
const navLinks = [
  { to: '/words', icon: '📖', label: 'Words' },
  { to: '/groups', icon: '🏷️', label: 'Groups' },
  { to: '/study-activities', icon: '📚', label: 'Study Activities' },
  { to: '/sessions', icon: '⏱️', label: 'Sessions' }
]
</script>
