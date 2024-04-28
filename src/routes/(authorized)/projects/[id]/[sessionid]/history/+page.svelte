<script lang="ts">
	import { afterNavigate, invalidate, pushState } from '$app/navigation';
	import { get } from 'svelte/store';
	import { page } from '$app/stores';
	import Command from '$lib/components/custom/command/Command.svelte';
	import CommandLoading from '$lib/components/custom/command/CommandLoading.svelte';
	import { onMount } from 'svelte';
	import { createWebsocketStore, type WebsocketStore } from '$lib/stores/Websocket';

	let { data } = $props();

	let pageLimit = $state(5);
	let currentPage = $state(1);
	let offsetFromOriginalPage = 0;

	const urlPrefix = get(page).url.pathname;
	let { id, sessionid } = get(page).params;

	const loadPage = (pageNo: number) => {
		currentPage = pageNo;
		pushState(`${urlPrefix}?page=${currentPage}&pageSize=${pageLimit}`, {});
		data.commands = fetch(
			`/api/projects/${id}/sessions/${sessionid}/commands?page=${currentPage}&pageSize=${pageLimit}`
		).then((res) => res.json());

		offsetFromOriginalPage = 0;
	};

	afterNavigate(() => {
		invalidate((url) => {
			return url.pathname.startsWith(`/api/projects/${id}/sessions/${sessionid}/commands`);
		});
		id = get(page).params.id;
		sessionid = get(page).params.sessionid;

		offsetFromOriginalPage = 0;
	});

	let websocketStore: WebsocketStore;

	const handleNewCommand = (event: MessageEvent) => {
		try {
			const data = JSON.parse(event.data);
			if (data.message_type !== 'new-command') return;
			data.count = data.count + 1;
			if (offsetFromOriginalPage == 0) {
				currentPage += 1;
			}
			offsetFromOriginalPage += 1;
			offsetFromOriginalPage = offsetFromOriginalPage % pageLimit;
		} catch (error) {
			console.error('Error parsing websocket message', error);
		}
	};

	onMount(() => {
		websocketStore = createWebsocketStore();

		websocketStore.ws?.addEventListener('message', handleNewCommand);

		return () => {
			websocketStore.ws?.removeEventListener('message', handleNewCommand);
		};
	});
</script>

<div class="mt-3 flex justify-center space-x-3">
	<div class="join">
		<button
			class="btn join-item btn-sm"
			class:btn-disabled={currentPage === 1}
			on:click={() => loadPage(currentPage - 1)}>«</button
		>
		<button class="btn join-item btn-sm hover:cursor-auto hover:border-base-200 hover:bg-base-200"
			>Page {currentPage}</button
		>
		{#await data.count}
			<button class="btn btn-disabled join-item btn-sm">»</button>
		{:then count}
			<button
				class="btn join-item btn-sm"
				class:btn-disabled={currentPage * pageLimit >= count}
				on:click={() => loadPage(currentPage + 1)}>»</button
			>
		{/await}
	</div>
	<div class="mt-1 md:flex-1">
		<div class="flex items-center justify-center space-x-3">
			<span>Limit: {pageLimit}</span>
			<input
				type="range"
				min="1"
				max="20"
				bind:value={pageLimit}
				class="range range-sm max-w-32"
				on:mouseup={() => loadPage(currentPage)}
				on:touchend={() => loadPage(currentPage)}
			/>
		</div>
	</div>
</div>
<div class="mt-5 space-y-5">
	{#await data.commands}
		<CommandLoading />
	{:then commands}
		{#if commands === undefined || commands.length === 0}
			<p class="text-center">No commands found</p>
		{:else}
			{#each commands as command}
				<Command {command} />
			{/each}
		{/if}
	{/await}
</div>
