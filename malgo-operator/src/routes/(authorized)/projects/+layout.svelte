<script lang="ts">
	import { createWebsocketStore, type WebsocketStore } from '$lib/stores/Websocket';
	import { onDestroy, onMount } from 'svelte';
	import { browser } from '$app/environment';

	let websocketStore: WebsocketStore;

	onMount(() => {
		websocketStore = createWebsocketStore();
	});

	onDestroy(() => {
		if (browser) {
			websocketStore.disconnect();
		}
	});
</script>

<slot />
