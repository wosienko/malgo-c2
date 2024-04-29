<script lang="ts">
	import Command from '$lib/components/custom/command/Command.svelte';
	import type { Command as CommandType } from '$lib/components/custom/command/command';
	import { onDestroy, onMount } from 'svelte';
	import { createWebsocketStore, type WebsocketStore } from '$lib/stores/Websocket.js';
	import { page } from '$app/stores';
	import { get } from 'svelte/store';
	import { browser } from '$app/environment';
	import CommandLoading from '$lib/components/custom/command/CommandLoading.svelte';
	import { afterNavigate } from '$app/navigation';
	import KeyValueTable from '$lib/components/custom/session/key-value/KeyValueTable.svelte';
	import KeyValueLoading from '$lib/components/custom/session/key-value/KeyValueLoading.svelte';

	let { data } = $props();

	// undefined - no commands yet
	// null - loading
	let newestCommand: CommandType | undefined | null = $state(null);

	let commandToBeSent = $state('');

	let websocketStore: WebsocketStore;

	const sendNewCommand = () => {
		websocketStore.sendCommand(get(page).params.sessionid, commandToBeSent);
		commandToBeSent = '';
	};

	const handleNewCommand = (event: MessageEvent) => {
		try {
			const data = JSON.parse(event.data);
			if (data.message_type !== 'new-command') return;
			newestCommand = data;
			newestCommand!.created_at = formatDateAndTime(newestCommand!.created_at);
		} catch (error) {
			console.error('Error parsing websocket message', error);
		}
	};

	onMount(async () => {
		newestCommand = await data.command;

		websocketStore = createWebsocketStore();

		websocketStore.ws?.addEventListener('message', handleNewCommand);
	});

	afterNavigate(async () => {
		newestCommand = null;
		newestCommand = await fetch(
			`/api/projects/${get(page).params.id}/sessions/${get(page).params.sessionid}/command`
		)
			.then((res) => res.json())
			.catch(() => undefined);
	});

	onDestroy(() => {
		if (browser) {
			websocketStore.ws?.removeEventListener('message', handleNewCommand);
		}
	});

	const formatDateAndTime = (date: string): string => {
		const d = new Date(date);
		// format DD.MM.YYYY HH:MM:SS
		return `${d.getDate().toString().padStart(2, '0')}.${(d.getMonth() + 1).toString().padStart(2, '0')}.${d.getFullYear()} ${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}:${d.getSeconds().toString().padStart(2, '0')}`;
	};

	const convertObjectToMap = (obj: Record<string, string>): Map<string, string> => {
		const map = new Map<string, string>();
		for (const [key, value] of Object.entries(obj)) {
			map.set(key, value);
		}
		return map;
	};
</script>

<div class="flex h-1/2 w-full justify-around border-b-2 border-neutral pb-6">
	<div class="m-3 flex h-full flex-1 flex-col">
		<textarea
			class="textarea textarea-warning w-full grow resize-none"
			bind:value={commandToBeSent}
			placeholder="Command"
		></textarea>
		<button
			class="btn btn-warning mt-3 w-full"
			class:btn-disabled={commandToBeSent === ''}
			on:click={sendNewCommand}>Run</button
		>
	</div>
	<div class="m-3 hidden h-full flex-1 overflow-auto md:block">
		{#await data.keyValue}
			<KeyValueLoading />
		{:then keyValue}
			<KeyValueTable values={convertObjectToMap(keyValue)} />
		{/await}
	</div>
</div>
<div class="h-1/2 w-full overflow-y-auto py-3">
	{#if newestCommand === null}
		<CommandLoading />
	{:else if newestCommand}
		<Command command={newestCommand} />
	{:else}
		<p class="text-center">No commands yet</p>
	{/if}
</div>
