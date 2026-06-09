import { useEffect, useRef, useCallback, useState } from 'react'

export function useWebSocket(url: string) {
  const ws = useRef<WebSocket | null>(null)
  const [lastMessage, setLastMessage] = useState<string | null>(null)
  const [connected, setConnected] = useState(false)

  useEffect(() => {
    const socket = new WebSocket(url)
    ws.current = socket

    socket.onopen = () => setConnected(true)
    socket.onclose = () => setConnected(false)
    socket.onmessage = (event) => setLastMessage(event.data)

    return () => {
      socket.close()
    }
  }, [url])

  const send = useCallback((data: unknown) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(data))
    }
  }, [])

  return { send, lastMessage, connected }
}
