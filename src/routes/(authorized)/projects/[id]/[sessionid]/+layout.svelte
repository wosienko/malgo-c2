<script lang="ts">
	import { page } from '$app/stores';
	import { get } from 'svelte/store';
	import { afterNavigate, onNavigate } from '$app/navigation';
	import { onDestroy, onMount } from 'svelte';
	import { createWebsocketStore, type WebsocketStore } from '$lib/stores/Websocket';
	import { browser } from '$app/environment';

	const ALL_TABS = ['/commands', '/modules', '/history', '/settings'] as const;

	let { id, sessionid } = $state(get(page).params);

	let currentUrl = $state(get(page).url.pathname);

	afterNavigate(() => {
		currentUrl = get(page).url.pathname;
		id = get(page).params.id;
		sessionid = get(page).params.sessionid;
	});

	const capitalize = (s: string) => {
		return s.charAt(0).toUpperCase() + s.slice(1);
	};

	let websocketStore: WebsocketStore;

	onMount(async () => {
		websocketStore = createWebsocketStore();
		await websocketStore.subscribeToSession(get(page).params.sessionid);
	});

	onNavigate(async () => {
		if (browser && get(page).params.id !== '') {
			await websocketStore.unsubscribeFromSession();
			await websocketStore.subscribeToSession(get(page).params.sessionid);
		}
	});

	onDestroy(async () => {
		if (browser) {
			await websocketStore.unsubscribeFromSession();
		}
	});
</script>

<div class="flex h-full w-full flex-col pt-3">
	<div role="tablist" class="tabs tabs-bordered tabs-md w-full">
		{#each ALL_TABS as tab}
			<a
				role="tab"
				class="tab pb-9 hover:opacity-50"
				class:tab-active={currentUrl.endsWith(tab)}
				href={`/projects/${id}/${sessionid}${tab}`}
				>{capitalize(tab.replace('/', '').replace('-', ' '))}</a
			>
		{/each}
	</div>
	<div class="max-w-full grow overflow-x-auto">
		<slot />
	</div>
</div>
