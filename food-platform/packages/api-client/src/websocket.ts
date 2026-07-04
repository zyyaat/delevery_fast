// WebSocket client with auto-reconnect

import type { WSMessage, WSMessageType } from '@food-platform/types'
import { uuid } from '@food-platform/utils'

export interface WSClientConfig {
  url: string
  token: string
  protocols?: string[]
  reconnectInterval?: number
  maxReconnectAttempts?: number
  heartbeatInterval?: number
}

export type WSMessageHandler = (message: WSMessage) => void
export type WSStatusHandler = (status: WSConnectionStatus) => void

export type WSConnectionStatus = 'connecting' | 'connected' | 'disconnected' | 'reconnecting' | 'failed'

export class WSClient {
  private ws: WebSocket | null = null
  private config: WSClientConfig
  private reconnectAttempts = 0
  private heartbeatTimer: ReturnType<typeof setInterval> | null = null
  private messageHandlers = new Map<WSMessageType | '*', Set<WSMessageHandler>>()
  private statusHandlers = new Set<WSStatusHandler>()
  private isManualClose = false

  constructor(config: WSClientConfig) {
    this.config = {
      reconnectInterval: 1000,
      maxReconnectAttempts: 10,
      heartbeatInterval: 30000,
      ...config,
    }
  }

  connect() {
    this.isManualClose = false
    this.setStatus('connecting')
    this.doConnect()
  }

  private doConnect() {
    try {
      this.ws = new WebSocket(this.config.url, this.config.protocols ?? ['food-platform.v1'])
      this.setupEventListeners()
    } catch (error) {
      this.handleDisconnect()
    }
  }

  private setupEventListeners() {
    if (!this.ws) return

    this.ws.onopen = () => {
      this.reconnectAttempts = 0
      this.setStatus('connected')

      // Send auth
      this.send({
        event: 'auth',
        payload: { token: this.config.token },
        timestamp: new Date().toISOString(),
        id: uuid(),
      })

      // Start heartbeat
      this.startHeartbeat()
    }

    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data) as WSMessage
        this.handleMessage(message)
      } catch (error) {
        console.error('Failed to parse WS message:', error)
      }
    }

    this.ws.onerror = () => {
      // Error will be followed by close
    }

    this.ws.onclose = () => {
      this.stopHeartbeat()
      if (!this.isManualClose) {
        this.handleDisconnect()
      }
    }
  }

  private handleMessage(message: WSMessage) {
    // Notify specific handlers
    const specificHandlers = this.messageHandlers.get(message.event)
    specificHandlers?.forEach((handler) => handler(message))

    // Notify wildcard handlers
    const wildcardHandlers = this.messageHandlers.get('*')
    wildcardHandlers?.forEach((handler) => handler(message))
  }

  private handleDisconnect() {
    if (this.reconnectAttempts >= (this.config.maxReconnectAttempts ?? 10)) {
      this.setStatus('failed')
      return
    }

    this.setStatus('reconnecting')
    const delay = Math.min(
      (this.config.reconnectInterval ?? 1000) * 2 ** this.reconnectAttempts,
      30000,
    )
    this.reconnectAttempts++

    setTimeout(() => {
      this.doConnect()
    }, delay)
  }

  private startHeartbeat() {
    const interval = this.config.heartbeatInterval ?? 30000
    this.heartbeatTimer = setInterval(() => {
      this.send({
        event: 'ping',
        payload: {},
        timestamp: new Date().toISOString(),
        id: uuid(),
      })
    }, interval)
  }

  private stopHeartbeat() {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer)
      this.heartbeatTimer = null
    }
  }

  send(message: Omit<WSMessage, 'id'> & { id?: string }) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(
        JSON.stringify({
          ...message,
          id: message.id ?? uuid(),
        }),
      )
    }
  }

  subscribe(channel: string) {
    this.send({
      event: 'subscribe',
      payload: { channel },
      timestamp: new Date().toISOString(),
      id: uuid(),
    })
  }

  unsubscribe(channel: string) {
    this.send({
      event: 'unsubscribe',
      payload: { channel },
      timestamp: new Date().toISOString(),
      id: uuid(),
    })
  }

  on<T extends WSMessage = WSMessage>(
    event: WSMessageType | '*',
    handler: (message: T) => void,
  ): () => void {
    if (!this.messageHandlers.has(event)) {
      this.messageHandlers.set(event, new Set())
    }
    const handlers = this.messageHandlers.get(event)!
    handlers.add(handler as WSMessageHandler)

    return () => {
      handlers.delete(handler as WSMessageHandler)
    }
  }

  onStatus(handler: WSStatusHandler): () => void {
    this.statusHandlers.add(handler)
    return () => {
      this.statusHandlers.delete(handler)
    }
  }

  private setStatus(status: WSConnectionStatus) {
    this.statusHandlers.forEach((handler) => handler(status))
  }

  disconnect() {
    this.isManualClose = true
    this.stopHeartbeat()
    this.ws?.close()
    this.ws = null
    this.setStatus('disconnected')
  }
}

// ============ React hook helper (singleton) ============

let _wsInstance: WSClient | null = null

export function initWSClient(config: WSClientConfig): WSClient {
  _wsInstance = new WSClient(config)
  return _wsInstance
}

export function getWSClient(): WSClient | null {
  return _wsInstance
}
