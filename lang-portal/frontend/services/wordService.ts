import { ref } from 'vue'

// Word interface matching backend schema
export interface Word {
  id: number
  hindi: string
  scrambled?: string
  hinglish: string
  english: string
  created_at: string
}

// Pagination interface
interface PaginationParams {
  page?: number
  pageSize?: number
  language?: string
}

// Default pagination options
const DEFAULT_PAGE_SIZE = 10

// Base API URL
const BASE_URL = 'http://localhost:3000/api'

export function useWordService() {
  const words = ref<Word[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchWords(params: PaginationParams = {}) {
    loading.value = true
    error.value = null

    try {
      const queryParams = new URLSearchParams()
      if (params.page) queryParams.append('page', params.page.toString())
      if (params.pageSize) queryParams.append('pageSize', params.pageSize.toString())
      else queryParams.append('pageSize', DEFAULT_PAGE_SIZE.toString())
      if (params.language) queryParams.append('language', params.language)

      const response = await fetch(`${BASE_URL}/words?${queryParams}`, {
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
      words.value = data.words || []
      
      return {
        words: words.value,
        totalCount: data.totalCount || 0
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

  async function searchWords(query: string, language?: string) {
    loading.value = true
    error.value = null

    try {
      const queryParams = new URLSearchParams()
      queryParams.append('query', query)
      if (language) queryParams.append('language', language)

      const response = await fetch(`${BASE_URL}/words/search?${queryParams}`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        }
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to search words: ${response.status} - ${errorText}`)
      }

      const data = await response.json()
      return {
        words: data.words || [],
        totalCount: data.totalCount || 0
      }
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : 'An unknown error occurred while searching words'
      
      error.value = errorMessage
      
      throw err
    } finally {
      loading.value = false
    }
  }

  async function getWordById(id: number) {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${BASE_URL}/words/${id}`, {
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

  async function getWordsByGroupId(groupId: number) {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${BASE_URL}/words/groups/${groupId}`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        }
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to fetch words for group ${groupId}: ${response.status} - ${errorText}`)
      }

      return await response.json()
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : `An unknown error occurred while fetching words for group ${groupId}`
      
      error.value = errorMessage
      throw err
    } finally {
      loading.value = false
    }
  }

  async function getRandomWord() {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${BASE_URL}/words/random`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        }
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to fetch random word: ${response.status} - ${errorText}`)
      }

      return await response.json()
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : 'An unknown error occurred while fetching a random word'
      
      error.value = errorMessage
      throw err
    } finally {
      loading.value = false
    }
  }

  async function createWord(wordData: { hindi: string; english: string; hinglish: string }) {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${BASE_URL}/words`, {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(wordData)
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to create word: ${response.status} - ${errorText}`)
      }

      return await response.json()
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : 'An unknown error occurred while creating a word'
      
      error.value = errorMessage
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateWord(id: number, wordData: { hindi: string; english: string; hinglish: string }) {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${BASE_URL}/words/${id}`, {
        method: 'PUT',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(wordData)
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to update word: ${response.status} - ${errorText}`)
      }

      return await response.json()
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : 'An unknown error occurred while updating a word'
      
      error.value = errorMessage
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteWord(id: number) {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${BASE_URL}/words/${id}`, {
        method: 'DELETE',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        }
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to delete word: ${response.status} - ${errorText}`)
      }

      return true
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : 'An unknown error occurred while deleting a word'
      
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
    searchWords,
    getWordById,
    getWordsByGroupId,
    getRandomWord,
    createWord,
    updateWord,
    deleteWord
  }
}
