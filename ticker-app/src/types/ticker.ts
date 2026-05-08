export interface TickerEntry {
  ticker_id: number
  ticker_name: string
  TeamtipId: number
  created: string
  modified: string
  title: string
  message: string
  type: number
  creator: string
  highlight: boolean
  automatic: boolean
  active: boolean
}

export interface CreateEntryPayload {
  ticker_name: string
  title: string
  message: string
  type: number
  creator: string
  highlight: boolean
}

export interface UpdateEntryPayload {
  ticker_id: number
  title: string
  message: string
  type: number
  creator: string
  highlight: boolean
  active: boolean
}

export type SortField = 'created' | 'title' | 'ticker_id'
export type SortDirection = 'asc' | 'desc'

export interface SortConfig {
  field: SortField
  direction: SortDirection
}

export type AppStatus = 'idle' | 'loading' | 'success' | 'error'