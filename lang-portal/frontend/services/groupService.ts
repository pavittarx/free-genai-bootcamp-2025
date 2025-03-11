import { ref } from 'vue'

// Group interface matching backend schema
export interface Group {
  id: number
  name: string
  description: string
  created_at: string
}

// Pagination interface
interface PaginationParams {
  page?: number
  pageSize?: number
  search?: string
  sortBy?: string
  sortDirection?: 'asc' | 'desc'
}

// Base API URL
const BASE_URL = 'http://localhost:3000/api'

export function useGroupService() {
  const groups = ref<Group[]>([])
  const group = ref<Group | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchGroups(params: PaginationParams = {}) {
    loading.value = true
    error.value = null

    try {
      const queryParams = new URLSearchParams()
      if (params.page) queryParams.append('page', params.page.toString())
      if (params.pageSize) queryParams.append('pageSize', params.pageSize.toString())
      if (params.search) queryParams.append('search', params.search)
      if (params.sortBy) queryParams.append('sortBy', params.sortBy)
      if (params.sortDirection) queryParams.append('sortDirection', params.sortDirection)

      const response = await fetch(`${BASE_URL}/groups?${queryParams}`, {
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
      console.log('API Response:', data)
      
      // Directly use the groups array from the response
      groups.value = data.groups
      
      return {
        groups: groups.value,
        totalCount: data.total || 0
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
      const response = await fetch(`${BASE_URL}/groups/${id}`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        }
      })

      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to fetch group: ${response.status} - ${errorText}`)
      }

      const data = await response.json()
      
      // Correctly parse the nested group object
      group.value = {
        id: data.group.id,
        name: data.group.name,
        description: data.group.description,
        created_at: data.group.created_at
      }

      return group.value
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : 'An unknown error occurred while fetching group'
      
      error.value = errorMessage
      group.value = null
      
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    groups,
    group,
    loading,
    error,
    fetchGroups,
    getGroupById
  }
}
