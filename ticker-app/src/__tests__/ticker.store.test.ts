import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useTickerStore } from '@/stores/ticker'
import type { TickerEntry } from '@/types/ticker'

// API mocken — wir testen den Store, nicht die API
vi.mock('@/services/tickerApi', () => ({
  tickerApi: {
    getEntries: vi.fn(),
    createEntry: vi.fn(),
    updateEntry: vi.fn(),
    deleteEntry: vi.fn(),
  }
}))

import { tickerApi } from '@/services/tickerApi'

// Testdaten
const mockEntries: TickerEntry[] = [
  {
    ticker_id: 1,
    ticker_name: 'wpftest',
    TeamtipId: 0,
    created: '2023-01-01T10:00:00Z',
    modified: '2023-01-01T10:00:00Z',
    title: 'Alpha Eintrag',
    message: 'Erste Nachricht',
    type: 1,
    creator: 'Cyril',
    highlight: true,
    automatic: false,
    active: true,
  },
  {
    ticker_id: 2,
    ticker_name: 'wpftest',
    TeamtipId: 0,
    created: '2023-02-01T10:00:00Z',
    modified: '2023-02-01T10:00:00Z',
    title: 'Beta Eintrag',
    message: 'Zweite Nachricht von Dominik',
    type: 1,
    creator: 'Dominik',
    highlight: false,
    automatic: false,
    active: true,
  },
  {
    ticker_id: 3,
    ticker_name: 'wpftest',
    TeamtipId: 0,
    created: '2023-03-01T10:00:00Z',
    modified: '2023-03-01T10:00:00Z',
    title: 'Gamma Eintrag',
    message: 'Dritte Nachricht',
    type: 1,
    creator: 'Cyril',
    highlight: false,
    automatic: false,
    active: true,
  },
]

describe('Ticker Store', () => {

  beforeEach(() => {
    setActivePinia(createPinia())
  })

  // ── State Tests ──────────────────────────────────────

  it('hat korrekten initialen State', () => {
    const store = useTickerStore()

    expect(store.entries).toEqual([])
    expect(store.status).toBe('idle')
    expect(store.errorMessage).toBeNull()
    expect(store.searchQuery).toBe('')
    expect(store.currentPage).toBe(1)
    expect(store.pageSize).toBe(3)
  })

  // ── Suche Tests ──────────────────────────────────────

  it('filtert Eintraege nach Titel', () => {
    const store = useTickerStore()
    store.entries = mockEntries

    store.setSearch('Alpha')

    expect(store.filteredEntries).toHaveLength(1)
    expect(store.filteredEntries[0].title).toBe('Alpha Eintrag')
  })

  it('filtert Eintraege nach Ersteller', () => {
    const store = useTickerStore()
    store.entries = mockEntries

    store.setSearch('Dominik')

    expect(store.filteredEntries).toHaveLength(1)
    expect(store.filteredEntries[0].creator).toBe('Dominik')
  })

  it('filtert Eintraege nach Nachricht', () => {
    const store = useTickerStore()
    store.entries = mockEntries

    store.setSearch('Zweite')

    expect(store.filteredEntries).toHaveLength(1)
    expect(store.filteredEntries[0].message).toContain('Zweite')
  })

  it('gibt alle Eintraege zurueck wenn Suche leer ist', () => {
    const store = useTickerStore()
    store.entries = mockEntries

    store.setSearch('')

    expect(store.filteredEntries).toHaveLength(3)
  })

  it('setzt currentPage auf 1 bei neuer Suche', () => {
    const store = useTickerStore()
    store.entries = mockEntries
    store.currentPage = 3

    store.setSearch('Alpha')

    expect(store.currentPage).toBe(1)
  })

  // ── Sortierung Tests ─────────────────────────────────

  it('sortiert Eintraege nach Titel aufsteigend', () => {
    const store = useTickerStore()
    store.entries = mockEntries

    store.setSort({ field: 'title', direction: 'asc' })

    expect(store.sortedEntries[0].title).toBe('Alpha Eintrag')
    expect(store.sortedEntries[1].title).toBe('Beta Eintrag')
    expect(store.sortedEntries[2].title).toBe('Gamma Eintrag')
  })

  it('sortiert Eintraege nach Titel absteigend', () => {
    const store = useTickerStore()
    store.entries = mockEntries

    store.setSort({ field: 'title', direction: 'desc' })

    expect(store.sortedEntries[0].title).toBe('Gamma Eintrag')
    expect(store.sortedEntries[2].title).toBe('Alpha Eintrag')
  })

  it('sortiert Eintraege nach ID aufsteigend', () => {
    const store = useTickerStore()
    store.entries = mockEntries

    store.setSort({ field: 'ticker_id', direction: 'asc' })

    expect(store.sortedEntries[0].ticker_id).toBe(1)
    expect(store.sortedEntries[2].ticker_id).toBe(3)
  })

  // ── Pagination Tests ─────────────────────────────────

  it('berechnet totalPages korrekt', () => {
    const store = useTickerStore()
    store.entries = mockEntries
    store.pageSize = 2

    expect(store.totalPages).toBe(2)
  })

  it('gibt korrekte Eintraege fuer Seite 1 zurueck', () => {
    const store = useTickerStore()
    store.entries = mockEntries
    store.pageSize = 2

    store.setPage(1)

    expect(store.paginatedEntries).toHaveLength(2)
  })

  it('gibt korrekte Eintraege fuer letzte Seite zurueck', () => {
    const store = useTickerStore()
    store.entries = [...mockEntries]
    store.pageSize = 2

    store.setSort({ field: 'ticker_id', direction: 'asc' })
    store.setPage(2)

    expect(store.paginatedEntries).toHaveLength(1)
    expect(store.paginatedEntries[0].ticker_id).toBe(3)
  })

  // ── API Actions Tests ────────────────────────────────

  it('laedt Eintraege erfolgreich', async () => {
    const store = useTickerStore()
    vi.mocked(tickerApi.getEntries).mockResolvedValue(mockEntries)

    await store.fetchEntries()

    expect(store.status).toBe('success')
    expect(store.entries).toHaveLength(3)
    expect(store.errorMessage).toBeNull()
  })

  it('setzt status auf error bei API-Fehler', async () => {
    const store = useTickerStore()
    vi.mocked(tickerApi.getEntries).mockRejectedValue(new Error('Netzwerkfehler'))

    await store.fetchEntries()

    expect(store.status).toBe('error')
    expect(store.errorMessage).toBe('Netzwerkfehler')
  })

  it('fuegt neuen Eintrag am Anfang der Liste hinzu', async () => {
    const store = useTickerStore()
    store.entries = [...mockEntries]

    const newEntry: TickerEntry = {
      ticker_id: 99,
      ticker_name: 'wpftest',
      TeamtipId: 0,
      created: '2023-04-01T10:00:00Z',
      modified: '2023-04-01T10:00:00Z',
      title: 'Neuer Eintrag',
      message: 'Neue Nachricht',
      type: 1,
      creator: 'Cyril',
      highlight: false,
      automatic: false,
      active: true,
    }

    vi.mocked(tickerApi.createEntry).mockResolvedValue(newEntry)

    await store.createEntry({
      ticker_name: 'wpftest',
      title: 'Neuer Eintrag',
      message: 'Neue Nachricht',
      creator: 'Cyril',
      type: 1,
      highlight: false,
    })

    expect(store.entries[0].ticker_id).toBe(99)
    expect(store.entries).toHaveLength(4)
  })

  it('entfernt Eintrag nach deleteEntry', async () => {
    const store = useTickerStore()
    store.entries = [...mockEntries]

    vi.mocked(tickerApi.deleteEntry).mockResolvedValue(undefined)

    await store.deleteEntry(1)

    expect(store.entries).toHaveLength(2)
    expect(store.entries.find(e => e.ticker_id === 1)).toBeUndefined()
  })

})