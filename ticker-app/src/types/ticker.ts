// Echte API-Antwort von der Wettstar API
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

// Payload beim Erstellen eines neuen Eintrags
export interface CreateEntryPayload {
  ticker_name: string
  title: string
  message: string
  type: number
  creator: string
  highlight: boolean
}

// Payload beim Bearbeiten eines Eintrags
export interface UpdateEntryPayload {
  ticker_id: number
  title: string
  message: string
  type: number
  creator: string
  highlight: boolean
  active: boolean
}

// Sortierung
export type SortField = 'created' | 'title' | 'ticker_id'
export type SortDirection = 'asc' | 'desc'

export interface SortConfig {
  field: SortField
  direction: SortDirection
}

// Status der App
export type AppStatus = 'idle' | 'loading' | 'success' | 'error'