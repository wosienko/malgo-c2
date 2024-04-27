import { env } from '$env/dynamic/public';

let ws: WebSocket | null = null;
let reconnectInterval = 2000; // Interval in milliseconds

let consciousExit = false;

// Function to initiate a WebSocket connection
const connect = () => {
	if (ws !== null && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) {
		return;
	}

	ws = new WebSocket(env.PUBLIC_WEBSOCKET_URL);

	ws.onopen = () => {
		console.log('WebSocket connected');
		// Reset reconnect interval on successful connection
		reconnectInterval = 2000;
	};

	ws.onclose = () => {
		console.log('WebSocket disconnected');
		if (consciousExit) {
			consciousExit = false;
			return;
		}
		// Attempt to reconnect after a delay
		setTimeout(() => {
			console.log('Attempting to reconnect...');
			connect();
		}, reconnectInterval);

		// Increase reconnect interval for next time
		reconnectInterval = Math.min(reconnectInterval * 2, 10000);
	};

	ws.onerror = (error) => {
		console.error('WebSocket error:', error);
		ws!.close();
	};
};

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
			console.log('Subscribed to project', projectId);
		},
		unsubscribeFromProject: async () => {
			while (ws === null || ws.readyState !== WebSocket.OPEN) {
				await new Promise((resolve) => setTimeout(resolve, 100));
			}
			ws.send(JSON.stringify({ type: 'unsubscribeProject' }));
			console.log('Unsubscribed from project');
		},
		subscribeToSession: async (sessionId: string) => {
			while (ws === null || ws.readyState !== WebSocket.OPEN) {
				await new Promise((resolve) => setTimeout(resolve, 100));
			}
			ws.send(JSON.stringify({ type: 'subscribeSession', sessionId }));
			console.log('Subscribed to session', sessionId);
		},
		unsubscribeFromSession: async () => {
			while (ws === null || ws.readyState !== WebSocket.OPEN) {
				await new Promise((resolve) => setTimeout(resolve, 100));
			}
			ws.send(JSON.stringify({ type: 'unsubscribeSession' }));
			console.log('Unsubscribed from session');
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

export type WebsocketStore = ReturnType<typeof createWebsocketStore>;
