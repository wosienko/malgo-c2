<script lang="ts">
	import EditableEntry from '$lib/components/custom/session/settings/EditableEntry.svelte';
	import { createWebsocketStore, type WebsocketStore } from '$lib/stores/Websocket.js';
	import { onMount } from 'svelte';
	import { get } from 'svelte/store';
	import { page } from '$app/stores';

	let { data } = $props();

	let websocketStore: WebsocketStore;

	const handleSessionRename = (event: MessageEvent) => {
		const dataFromWs = JSON.parse(event.data);
		if (dataFromWs.message_type === 'session-renamed') {
			if (dataFromWs.session_id !== get(page).params.sessionid) return;
			data.session_name = dataFromWs.name;
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

<div class="flex h-full flex-col items-center space-y-3">
	<table class="table table-fixed overflow-x-auto text-center">
		<thead>
			<tr>
				<th>Setting</th>
				<th>Value</th>
				<th></th>
			</tr>
		</thead>
		<tbody>
			<EditableEntry
				key="Session name"
				ws_type="session-rename"
				value={data.session_name}
				value_key="name"
			/>
		</tbody>
	</table>
</div>
