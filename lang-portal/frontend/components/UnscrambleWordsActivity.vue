<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 p-4">
    <div class="bg-gradient-to-br from-blue-100 to-blue-300 w-full max-w-2xl rounded-xl shadow-2xl overflow-hidden">
      <div class="p-8 flex flex-col min-h-[600px]">
        <div class="text-center mb-6">
          <h2 class="text-3xl font-bold text-gray-800">शब्द उलझाओ (Unscramble Words)</h2>
          <p class="text-gray-600">अक्षरों को सही क्रम में व्यवस्थित करें</p>
        </div>

        <div v-if="isLoading" class="flex-grow flex items-center justify-center">
          <div class="text-center">
            <div class="animate-pulse text-6xl mb-4">🧩</div>
            <p class="text-gray-600">चुनौतियाँ लोड हो रही हैं... (Loading challenges)</p>
            <p class="text-sm text-gray-500">
              चुनौती: {{ currentChallengeIndex + 1 }} / 10
            </p>
          </div>
        </div>

        <div v-else-if="error" class="flex-grow flex items-center justify-center">
          <div class="text-center">
            <div class="text-6xl mb-4">😕</div>
            <h3 class="text-xl font-semibold text-gray-800 mb-2">लोड करने में त्रुटि</h3>
            <p class="text-gray-600 mb-4">{{ errorMessage }}</p>
          </div>
        </div>

        <div v-else-if="activityCompleted" class="flex-grow flex items-center justify-center">
          <div class="text-center">
            <div class="text-6xl mb-4">🏆</div>
            <h3 class="text-xl font-semibold text-gray-800 mb-2">
              गतिविधि पूरी हुई (Activity Completed)
            </h3>
            <p class="text-2xl font-bold text-green-600 mb-4">
              आपका कुल स्कोर (Total Score): {{ score }} / 50
            </p>
            <div class="flex justify-center space-x-4">
              <button 
                @click="handleActivityEnd"
                class="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition"
              >
                समाप्त करें (Close)
              </button>
            </div>
          </div>
        </div>

        <div v-else-if="currentChallenge" class="flex-grow flex flex-col justify-between">
          <div class="space-y-6">
            <div class="flex justify-between items-center mb-4">
              <div class="flex space-x-2">
                <span class="text-sm font-medium text-gray-600">चुनौती (Challenge):</span>
                <span class="text-sm font-semibold text-blue-600">
                  {{ currentChallengeIndex + 1 }} / 10
                </span>
              </div>
              <div class="flex items-center space-x-2">
                <span class="text-sm font-medium text-gray-600">स्कोर (Score):</span>
                <span class="text-sm font-semibold text-green-600">{{ score }}</span>
              </div>
            </div>

            <div class="bg-white rounded-xl shadow-md p-6">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div class="bg-blue-50 rounded-lg p-4 shadow-inner">
                  <h3 class="text-lg font-semibold text-gray-800 mb-2">
                    उलझे हुए शब्द (Scrambled Word)
                  </h3>
                  <div class="flex justify-center space-x-2 mb-4">
                    <div 
                      v-for="(letter, index) in currentChallenge.scrambledWord.split('')" 
                      :key="index"
                      @click="selectLetter(index)"
                      class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center text-2xl font-bold text-blue-800 cursor-pointer hover:bg-blue-200 transition"
                    >
                      {{ letter }}
                    </div>
                  </div>
                </div>

                <div class="bg-green-50 rounded-lg p-4 shadow-inner">
                  <h3 class="text-lg font-semibold text-gray-800 mb-2">
                    आपका उत्तर (Your Answer)
                  </h3>
                  <div class="flex justify-center space-x-2 mb-4">
                    <div 
                      v-for="(letter, index) in userInput" 
                      :key="index"
                      @click="removeLetter(index)"
                      class="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center text-2xl font-bold text-green-800 cursor-pointer hover:bg-green-200 transition"
                    >
                      {{ letter }}
                    </div>
                  </div>
                </div>
              </div>

              <div class="mt-4 text-center">
                            
                <p class="text-lg text-gray-600">
                  अंग्रेजी (English): {{ currentWord?.english || '' }}
                </p>
              </div>
            </div>

            <div v-if="feedbackMessage" class="text-center">
              <p 
                :class="{
                  'text-green-600': feedbackMessage.type === 'success',
                  'text-red-600': feedbackMessage.type === 'error'
                }"
              >
                {{ feedbackMessage.text }}
              </p>
            </div>
          </div>

          <div class="mt-6 flex justify-center space-x-4">
            <button 
              @click="checkAnswer"
              class="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition disabled:opacity-50"
              :disabled="userInput.length !== currentChallenge.scrambledWord.length"
            >
              उत्तर जमा करें (Submit Answer)
            </button>
            <button 
              @click="skipChallenge"
              class="bg-gray-200 text-gray-700 px-6 py-3 rounded-lg hover:bg-gray-300 transition"
            >
              छोड़ें (Skip)
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject, onMounted, watch } from 'vue'
import { useQuery, type UseQueryOptions } from '@tanstack/vue-query'
import { apiService, type Word, type Session } from '~/services/api'
import axios from 'axios'

// Props definition
const props = defineProps<{
  activityId: string
}>()

// State management
const currentChallengeIndex = ref(0)
const userInput = ref<string[]>([])
const score = ref(0)
const activityCompleted = ref(false)
const session = ref<Session | null>(null)
const feedbackMessage = ref<{ type: 'success' | 'error'; text: string } | null>(null)

// Fetch random word function
const fetchRandomWord = async (index: number): Promise<Word> => {
  try {
    const words = await apiService.getRandomWords(1)
    
    if (!Array.isArray(words) || words.length === 0) {
      throw new Error('No words returned from server')
    }
    
    const word = words[0]
    if (!word || typeof word !== 'object' || !word.hindi || !word.scrambled) {
      throw new Error('Invalid word structure')
    }
    
    return word
  } catch (error) {
    console.error('Error fetching word:', error)
    throw error
  }
}

// Current challenge management
const currentChallenge = ref<{
  word: string
  scrambledWord: string
} | null>(null)

// Query for fetching words
const { 
  data: queryData, 
  refetch, 
  isLoading, 
  error 
} = useQuery<Word, Error>({
  queryKey: ['randomWord', currentChallengeIndex.value],
  queryFn: () => fetchRandomWord(currentChallengeIndex.value),
  enabled: currentChallengeIndex.value < 10,
  retry: 3,
  staleTime: 0
})

// Watch query data and update current challenge
watch(queryData, (newWord) => {
  if (newWord) {
    currentChallenge.value = {
      word: newWord.hindi,
      scrambledWord: newWord.scrambled
    }
  }
}, { immediate: true })

// Error message computation
const errorMessage = computed(() => {
  return error.value ? error.value.message : 'अज्ञात त्रुटि (Unknown error)'
})

// Letter selection methods
const selectLetter = (index: number) => {
  if (!currentChallenge.value) return
  
  const scrambledWord = currentChallenge.value.scrambledWord
  const letter = scrambledWord.split('')[index]
  
  if (userInput.value.length < scrambledWord.length) {
    userInput.value.push(letter)
  }
}

const removeLetter = (index: number) => {
  userInput.value.splice(index, 1)
}

// Answer checking method
const checkAnswer = async () => {
  if (!currentChallenge.value || !session.value) {
    console.warn('Cannot check answer: missing data')
    return
  }

  const userInputValue = userInput.value.join('')
  const correctAnswer = currentChallenge.value.scrambledWord
  const isCorrect = userInputValue.toLowerCase() === correctAnswer.toLowerCase()
  const challengeScore = isCorrect ? 5 : 0
  
  try {
    // Submit session activity
    await apiService.submitActivity({
      session_id: session.value.id,
      activity_id: props.activityId,
      challenge: currentChallenge.value.word,
      answer: correctAnswer,
      input: userInputValue,
      score: challengeScore
    })

    // Update feedback message
    feedbackMessage.value = {
      type: isCorrect ? 'success' : 'error',
      text: isCorrect 
        ? 'बधाई हो! आपका उत्तर सही है। (Congratulations! Your answer is correct.)' 
        : 'क्षमा करें, यह उत्तर सही नहीं है। (Sorry, this answer is incorrect.)'
    }

    // Update score
    score.value += challengeScore

    // Handle challenge progression
    if (currentChallengeIndex.value === 9) {
      activityCompleted.value = true
      await endActivity()
    } else {
      currentChallengeIndex.value++
      await resetChallenge()
    }
  } catch (error) {
    console.error('Answer submission error:', error)
    feedbackMessage.value = {
      type: 'error',
      text: 'उत्तर जमा करने में त्रुटि (Error submitting answer)'
    }
  }
}

// Skip challenge method
const skipChallenge = async () => {
  if (!currentChallenge.value || !session.value) return

  try {
    // Submit skipped activity
    await apiService.submitActivity({
      session_id: session.value.id,
      activity_id: props.activityId,
      challenge: currentChallenge.value.word,
      answer: currentChallenge.value.scrambledWord,
      input: '',
      score: 0
    })

    // Handle challenge progression
    if (currentChallengeIndex.value === 9) {
      activityCompleted.value = true
      await endActivity()
    } else {
      currentChallengeIndex.value++
      await resetChallenge()
    }
  } catch (error) {
    console.error('Skip challenge error:', error)
  }
}

// Reset challenge method
const resetChallenge = async () => {
  // Clear previous input
  userInput.value = []
  feedbackMessage.value = null

  // Explicitly refetch the word
  try {
    await refetch()
  } catch (error) {
    console.error('Challenge reset error:', error)
  }
}

// End activity method
const endActivity = async () => {
  try {
    if (session.value) {
      // Close session with final score
      await apiService.closeSession(session.value.id, score.value)
      
      // Mark activity as completed
      activityCompleted.value = true
    }
  } catch (error) {
    console.error('Error ending activity:', error)
  }
}

// Popup control injection
const popupControl = inject<{ close?: () => void }>('popupControl', {})

// Handle activity end
const handleActivityEnd = () => {
  if (typeof popupControl.close === 'function') {
    popupControl.close()
  }
}

// Initialize session on component mount
onMounted(async () => {
  try {
    // Create session
    session.value = await apiService.createSession(props.activityId)
    
    // Fetch first word
    await refetch()
  } catch (error) {
    console.error('Initialization error:', error)
    
    // Set error feedback
    feedbackMessage.value = {
      type: 'error',
      text: 'सत्र शुरू करने में त्रुटि (Error starting session)'
    }
    
    // Close popup on critical error
    if (typeof popupControl.close === 'function') {
      popupControl.close()
    }
  }
})
</script>

<style scoped>
/* Existing styles remain the same */
</style>
