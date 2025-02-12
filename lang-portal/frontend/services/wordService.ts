import { ref } from 'vue'

// Word interface matching backend schema
export interface Word {
  id: number
  hindi: string
  scrambled: string
  hinglish: string
  english: string
  created_at: string
}

// Pagination interface
interface PaginationParams {
  page?: number
  limit?: number
  search?: string
}

// Base API URL
const BASE_URL = 'http://localhost:3000'

export function useWordService() {
  const words = ref<Word[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchWords(params: PaginationParams = {}) {
    loading.value = true
    error.value = null
    words.value = []

    try {
      const queryParams = new URLSearchParams()
      if (params.page) queryParams.append('page', params.page.toString())
      if (params.limit) queryParams.append('limit', params.limit.toString())
      if (params.search) queryParams.append('search', params.search)

      const response = await fetch(`${BASE_URL}/api/words?${queryParams}`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        }
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to fetch words: ${response.status} - ${errorText}`)
      }

      const data = await response.json()

      if (!data || !Array.isArray(data.words)) {
        throw new Error('Invalid data structure received')
      }

      words.value = data.words
      
      return {
        words: words.value,
        total: data.total || 0
      }
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : 'An unknown error occurred while fetching words'
      
      error.value = errorMessage
      words.value = []
      
      throw err
    } finally {
      loading.value = false
    }
  }

  async function getWordById(id: number) {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${BASE_URL}/api/words/${id}`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        }
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to fetch word with id ${id}: ${response.status} - ${errorText}`)
      }

      return await response.json()
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : `An unknown error occurred while fetching word with id ${id}`
      
      error.value = errorMessage
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    words,
    loading,
    error,
    fetchWords,
    getWordById
  }
}
