import { ref } from 'vue'

// Group interface matching backend schema
export interface Group {
  id: number
  group: string
  created_at: string
}

// Word interface to match backend schema
export interface GroupWord {
  id: number
  hindi: string
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

export function useGroupService() {
  const groups = ref<Group[]>([])
  const groupWords = ref<GroupWord[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchGroups(params: PaginationParams = {}) {
    loading.value = true
    error.value = null
    groups.value = []

    try {
      const queryParams = new URLSearchParams()
      if (params.page) queryParams.append('page', params.page.toString())
      if (params.limit) queryParams.append('limit', params.limit.toString())
      if (params.search) queryParams.append('search', params.search)

      const response = await fetch(`${BASE_URL}/api/groups?${queryParams}`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        }
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to fetch groups: ${response.status} - ${errorText}`)
      }

      const data = await response.json()

      if (!data || !Array.isArray(data.groups)) {
        throw new Error('Invalid data structure received')
      }

      groups.value = data.groups
      
      return {
        groups: groups.value,
        total: data.total || 0
      }
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : 'An unknown error occurred while fetching groups'
      
      error.value = errorMessage
      groups.value = []
      
      throw err
    } finally {
      loading.value = false
    }
  }

  async function getGroupById(id: number) {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${BASE_URL}/api/groups/${id}`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        }
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to fetch group with id ${id}: ${response.status} - ${errorText}`)
      }

      return await response.json()
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : `An unknown error occurred while fetching group with id ${id}`
      
      error.value = errorMessage
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchGroupWords(groupId: number, params: PaginationParams = {}) {
    loading.value = true
    error.value = null
    groupWords.value = []

    try {
      const queryParams = new URLSearchParams()
      if (params.page) queryParams.append('page', params.page.toString())
      if (params.limit) queryParams.append('limit', params.limit.toString())
      if (params.search) queryParams.append('search', params.search)

      const response = await fetch(`${BASE_URL}/api/word-groups/${groupId}?${queryParams}`, {
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

      const data = await response.json()

      if (!data || !Array.isArray(data.words)) {
        throw new Error('Invalid data structure received')
      }

      groupWords.value = data.words
      
      return {
        words: groupWords.value,
        total: data.total || 0
      }
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : `An unknown error occurred while fetching group words`
      
      error.value = errorMessage
      groupWords.value = []
      
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    groups,
    groupWords,
    loading,
    error,
    fetchGroups,
    getGroupById,
    fetchGroupWords
  }
}
