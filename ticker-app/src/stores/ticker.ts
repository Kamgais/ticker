import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { tickerApi } from '@/services/tickerApi'
import type {
  TickerEntry,
  CreateEntryPayload,
  UpdateEntryPayload,
  SortConfig,
  AppStatus,
} from '@/types/ticker'

export const useTickerStore = defineStore('ticker', () => {

  const entries = ref<TickerEntry[]>([])
  const status = ref<AppStatus>('idle')
  const errorMessage = ref<string | null>(null)
  const searchQuery = ref('')
  const sortConfig = ref<SortConfig>({ field: 'created', direction: 'desc' })
  const currentPage = ref(1)
  const pageSize = ref(3)

  const filteredEntries = computed(() => {
    const q = searchQuery.value.toLowerCase().trim()
    if (!q) return entries.value
    return entries.value.filter(e =>
      e.title.toLowerCase().includes(q) ||
      e.message.toLowerCase().includes(q) ||
      e.creator.toLowerCase().includes(q)
    )
  })

  const sortedEntries = computed(() => {
    return [...filteredEntries.value].sort((a, b) => {
      const field = sortConfig.value.field
      const dir = sortConfig.value.direction === 'asc' ? 1 : -1

      if (field === 'created') {
        return (new Date(a.created).getTime() - new Date(b.created).getTime()) * dir
      }
      if (field === 'title') {
        return a.title.localeCompare(b.title) * dir
      }
      if (field === 'ticker_id') {
        return (a.ticker_id - b.ticker_id) * dir
      }
      return 0
    })
  })

  const totalPages = computed(() =>
    Math.ceil(sortedEntries.value.length / pageSize.value)
  )

  const paginatedEntries = computed(() => {
    const start = (currentPage.value - 1) * pageSize.value
    const end = start + pageSize.value
    return sortedEntries.value.slice(start, end)
  })

async function fetchEntries() {
  status.value = 'loading'
  errorMessage.value = null
  try {
    const all = await tickerApi.getEntries()
    // Nur aktive Einträge anzeigen
    entries.value = all.filter(e => e.active)
    status.value = 'success'
  } catch (e) {
    status.value = 'error'
    errorMessage.value = (e as Error).message
  }
}


  async function createEntry(payload: CreateEntryPayload) {
    try {
      const newEntry = await tickerApi.createEntry(payload)
      entries.value.unshift(newEntry)
    } catch (e) {
      errorMessage.value = (e as Error).message
      throw e
    }
  }

  async function updateEntry(payload: UpdateEntryPayload) {
    try {
      const updated = await tickerApi.updateEntry(payload)
      const index = entries.value.findIndex(e => e.ticker_id === payload.ticker_id)
      if (index !== -1) entries.value[index] = updated
    } catch (e) {
      errorMessage.value = (e as Error).message
      throw e
    }
  }

  async function deleteEntry(tickerId: number) {
    try {
      await tickerApi.deleteEntry(tickerId)
      entries.value = entries.value.filter(e => e.ticker_id !== tickerId)
    } catch (e) {
      errorMessage.value = (e as Error).message
      throw e
    }
  }

  function setSearch(query: string) {
    searchQuery.value = query
    currentPage.value = 1
  }

  function setSort(config: SortConfig) {
    sortConfig.value = config
    currentPage.value = 1
  }

  function setPage(page: number) {
    currentPage.value = page
  }

  return {
    entries,
    status,
    errorMessage,
    searchQuery,
    sortConfig,
    currentPage,
    pageSize,
    filteredEntries,
    sortedEntries,
    totalPages,
    paginatedEntries,
    fetchEntries,
    createEntry,
    updateEntry,
    deleteEntry,
    setSearch,
    setSort,
    setPage,
  }
})