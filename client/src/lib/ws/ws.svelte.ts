import type { OutMessage, InMessage } from "./messages";

export const CONNECTION_STATUS = {
  CONNECTED: 'connected',
  DISCONNECTED: 'disconnected',
  ERROR: 'error',
} as const;

export class WS {
  socket: WebSocket | null = null;
  messages: InMessage[] = $state([]);
  connectionStatus: string = $state(CONNECTION_STATUS.DISCONNECTED);

  constructor(url: string = "ws://localhost:8080") {
    let ws = new WebSocket(url)
    this.socket = ws

    ws.onopen = () => {
        this.socket = ws;
        this.connectionStatus = CONNECTION_STATUS.CONNECTED;
        console.log('WebSocket connected');
    };
    
    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        console.log('WebSocket message received:', data);
        this.messages.push(data);
    };
    
    ws.onclose = () => {
        this.socket = null;
        this.connectionStatus = CONNECTION_STATUS.DISCONNECTED;
        console.log('WebSocket disconnected');
    };
    
    ws.onerror = (error) => {
        this.connectionStatus = CONNECTION_STATUS.ERROR;
        console.error('WebSocket error:', error);
    };
  }

  send(message: OutMessage) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(message));
    } else {
      console.error("WebSocket is not open. Cannot send message:", this.socket?.readyState, message);
    }
  }
}