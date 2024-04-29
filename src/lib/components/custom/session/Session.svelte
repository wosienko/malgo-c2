<script lang="ts">
	import { onMount } from 'svelte';
	import { createWebsocketStore, type WebsocketStore } from '$lib/stores/Websocket';

	type InputProps = {
		id: string;
		name: string;
		createdAt: string;
		heartbeatAt: string;
	};

	let { id, name, createdAt, heartbeatAt }: InputProps = $props();

	const olderThanInSeconds = (date: string, seconds: number): boolean => {
		const d = new Date(date);
		const now = new Date();
		return (now.getTime() - d.getTime()) / 1000 > seconds;
	};

	const formatDateAndTime = (date: string): string => {
		const d = new Date(date);
		// format DD.MM.YYYY HH:MM:SS
		return `${d.getDate().toString().padStart(2, '0')}.${(d.getMonth() + 1).toString().padStart(2, '0')}.${d.getFullYear()} ${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}:${d.getSeconds().toString().padStart(2, '0')}`;
	};

	let websocketStore: WebsocketStore;

	const handleSessionRename = (event: MessageEvent) => {
		const dataFromWs = JSON.parse(event.data);
		if (dataFromWs.message_type === 'session-renamed') {
			if (dataFromWs.session_id !== id) return;
			name = dataFromWs.name;
		}
	};

	onMount(() => {
		websocketStore = createWebsocketStore();

		websocketStore.ws?.addEventListener('message', handleSessionRename);

		return () => {
			websocketStore.ws?.removeEventListener('message', handleSessionRename);
		};
	});
</script>

<div class="flex flex-col">
	<span class="w-64 text-xl text-info">{name}</span>
	<table class="text-left">
		<tbody>
			<tr>
				<td>Created at:</td>
				<td>{formatDateAndTime(createdAt)}</td>
			</tr>
			<tr>
				<td>Last Heartbeat:</td>
				<td class={olderThanInSeconds(heartbeatAt, 3600) ? 'text-error' : 'text-success'}
					>{formatDateAndTime(heartbeatAt)}</td
				>
			</tr>
		</tbody>
	</table>
</div>
