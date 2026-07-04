// WebSocket hook for React

import { useEffect, useState, useCallback, useRef } from 'react'
import { getWSClient, type WSConnectionStatus } from '@food-platform/api-client'
import type { WSMessage, WSMessageType } from '@food-platform/types'

export interface UseWebSocketOptions {
  onMessage?: <T = unknown>(message: WSMessage<T>) => void
  onStatusChange?: (status: WSConnectionStatus) => void
  autoConnect?: boolean
}

export function useWebSocket(
  channel: string | null,
  options: UseWebSocketOptions = {},
) {
  const { onMessage, onStatusChange, autoConnect = true } = options
  const [isConnected, setIsConnected] = useState(false)
  const [status, setStatus] = useState<WSConnectionStatus>('disconnected')
  const onMessageRef = useRef(onMessage)
  const onStatusRef = useRef(onStatusChange)

  useEffect(() => {
    onMessageRef.current = onMessage
    onStatusRef.current = onStatusChange
  })

  useEffect(() => {
    const ws = getWSClient()
    if (!ws || !channel) return

    if (autoConnect) {
      ws.connect()
    }

    // Subscribe to channel
    ws.subscribe(channel)

    // Subscribe to all messages
    const unsubscribeMessages = ws.on('*', (message) => {
      onMessageRef.current?.(message)
    })

    // Subscribe to status changes
    const unsubscribeStatus = ws.onStatus((newStatus) => {
      setStatus(newStatus)
      setIsConnected(newStatus === 'connected')
      onStatusRef.current?.(newStatus)
    })

    return () => {
      ws.unsubscribe(channel)
      unsubscribeMessages()
      unsubscribeStatus()
    }
  }, [channel, autoConnect])

  const send = useCallback(<T>(event: WSMessageType, payload: T) => {
    const ws = getWSClient()
    if (ws) {
      ws.send({
        event,
        payload,
        timestamp: new Date().toISOString(),
      })
    }
  }, [])

  return { isConnected, status, send }
}
