# Websocket

Real-time communication between the operator and the client is done using WebSockets.

## Websocket Store

Websocket store is a global store that is used to communicate with the client. It is used to send commands, key-value pairs, and other data to the client.

It handles automatic reconnections.

```ts
export const createWebsocketStore = () => {
    // Immediately try to connect upon store creation
    connect();

    return {
        // Provide access to the WebSocket object with a getter
        get ws() {
            return ws;
        },
        subscribeToProject: async (projectId: string) => {
            while (ws === null || ws.readyState !== WebSocket.OPEN) {
                await new Promise((resolve) => setTimeout(resolve, 100));
            }
            ws.send(JSON.stringify({ type: 'subscribeProject', projectId }));
        },
        unsubscribeFromProject: async () => {
            while (ws === null || ws.readyState !== WebSocket.OPEN) {
                await new Promise((resolve) => setTimeout(resolve, 100));
            }
            ws.send(JSON.stringify({ type: 'unsubscribeProject' }));
        },
        subscribeToSession: async (sessionId: string) => {
            while (ws === null || ws.readyState !== WebSocket.OPEN) {
                await new Promise((resolve) => setTimeout(resolve, 100));
            }
            ws.send(JSON.stringify({ type: 'subscribeSession', sessionId }));
        },
        unsubscribeFromSession: async () => {
            while (ws === null || ws.readyState !== WebSocket.OPEN) {
                await new Promise((resolve) => setTimeout(resolve, 100));
            }
            ws.send(JSON.stringify({ type: 'unsubscribeSession' }));
        },
        disconnect: () => {
            if (ws === null) {
                console.error('WebSocket not connected');
                return;
            }
            consciousExit = true;
            ws.close();
        },
        sendCommand: (sessionId: string, command: string) => {
            if (ws === null || ws.readyState !== WebSocket.OPEN) {
                console.error('WebSocket not connected');
                return;
            }
            ws.send(JSON.stringify({ type: 'command', sessionId, command }));
        }
    };
};
```

## Listening to Events

Each object willing to listen to Websocket event should create its own listener.

!!! warning
    It is crucial to unsubscribe from the event when the object is destroyed.

```ts hl_lines="4-6 10-12"
onMount(async () => {
    websocketStore = createWebsocketStore();
    await websocketStore.subscribeToSession(get(page).params.sessionid);
    websocketStore.ws?.addEventListener('close', async () => {
        await websocketStore.subscribeToSession(get(page).params.sessionid);
    });
});
onDestroy(async () => {
    if (browser) {
        websocketStore.ws?.removeEventListener('close', async () => {
            await websocketStore.subscribeToSession(get(page).params.sessionid);
        });
        await websocketStore.unsubscribeFromSession();
    }
});
```

## Sending Messages

```ts
websocketStore.ws?.send(
    JSON.stringify({
        type: 'session-delete-key',
        sessionId: get(page).params.sessionid,
        key: key
    })
);
```
