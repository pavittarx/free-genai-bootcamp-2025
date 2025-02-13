<template>
  <div class="api-debugger">
    <h2>API Debugger</h2>
    <div v-if="loading" class="loading">Loading...</div>
    <div v-if="error" class="error">
      Error: {{ error }}
    </div>
    <div v-if="words.length" class="words-list">
      <h3>Words ({{ words.length }})</h3>
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Hindi</th>
            <th>Hinglish</th>
            <th>English</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="word in words" :key="word.id">
            <td>{{ word.id }}</td>
            <td>{{ word.hindi }}</td>
            <td>{{ word.hinglish }}</td>
            <td>{{ word.english }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <button @click="testFetch">Test Fetch Words</button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useWordService, Word } from '~/services/wordService'

const { words, loading, error, fetchWords } = useWordService()

async function testFetch() {
  await fetchWords({
    page: 1,
    limit: 10
  })
}

onMounted(testFetch)
</script>

<style scoped>
.api-debugger {
  background-color: #f0f0f0;
  padding: 20px;
  margin: 20px;
  border-radius: 8px;
}

.loading {
  color: blue;
}

.error {
  color: red;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}
</style>
