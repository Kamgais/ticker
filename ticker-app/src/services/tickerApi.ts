import axios, { type AxiosInstance, AxiosError } from 'axios'
import type { TickerEntry, CreateEntryPayload, UpdateEntryPayload } from '@/types/ticker'

const BASE_URL = ''
const MAX_RETRIES = 3
const RETRY_DELAY_MS = 1000

function sleep(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms))
}

class TickerApiService {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: BASE_URL,
      timeout: 10000,
      headers: { 'Content-Type': 'application/json' },
    })
  }

  private extractErrorMessage(error: AxiosError): string {
    if (!error.response) {
      return 'Keine Verbindung zum Server. Bitte Internetverbindung prüfen.'
    }
    switch (error.response.status) {
      case 400: return 'Ungültige Anfrage. Bitte Eingaben prüfen.'
      case 403: return 'Zugriff verweigert.'
      case 404: return 'Ressource nicht gefunden.'
      case 429: return 'Zu viele Anfragen. Bitte kurz warten.'
      case 500: return 'Serverfehler. Bitte später erneut versuchen.'
      default:  return `Fehler ${error.response.status}`
    }
  }

  private async withRetry<T>(fn: () => Promise<T>): Promise<T> {
    for (let attempt = 1; attempt <= MAX_RETRIES; attempt++) {
      try {
        return await fn()
      } catch (error) {
        const isLastAttempt = attempt === MAX_RETRIES
        const isNetworkError = error instanceof Error &&
          error.message.includes('Verbindung')

        if (isLastAttempt || !isNetworkError) throw error

        console.warn(`Versuch ${attempt} fehlgeschlagen. Retry in ${RETRY_DELAY_MS}ms...`)
        await sleep(RETRY_DELAY_MS * attempt)
      }
    }
    throw new Error('Alle Versuche fehlgeschlagen.')
  }

  async getEntries(): Promise<TickerEntry[]> {
    return this.withRetry(async () => {
      try {
        const response = await this.client.get<TickerEntry[]>('/api/ticker')
        return response.data
      } catch (error) {
        throw new Error(this.extractErrorMessage(error as AxiosError))
      }
    })
  }

  async createEntry(payload: CreateEntryPayload): Promise<TickerEntry> {
    return this.withRetry(async () => {
      try {
        const response = await this.client.post<TickerEntry>('/api/ticker', payload)
        return response.data
      } catch (error) {
        throw new Error(this.extractErrorMessage(error as AxiosError))
      }
    })
  }

  async updateEntry(payload: UpdateEntryPayload): Promise<TickerEntry> {
    return this.withRetry(async () => {
      try {
        const response = await this.client.patch<TickerEntry>('/api/ticker', payload)
        return response.data
      } catch (error) {
        throw new Error(this.extractErrorMessage(error as AxiosError))
      }
    })
  }

  async deleteEntry(tickerId: number): Promise<void> {
    return this.withRetry(async () => {
      try {
        await this.client.delete(`/api/ticker/${tickerId}`)
      } catch (error) {
        throw new Error(this.extractErrorMessage(error as AxiosError))
      }
    })
  }
}

export const tickerApi = new TickerApiService()