<script lang="ts">
	import type { Command } from '$lib/components/custom/command/command';
	import { createWebsocketStore, type WebsocketStore } from '$lib/stores/Websocket';
	import { onMount } from 'svelte';

	const formatDateAndTime = (date: string): string => {
		const d = new Date(date);
		// format DD.MM.YYYY HH:MM:SS
		return `${d.getDate().toString().padStart(2, '0')}.${(d.getMonth() + 1).toString().padStart(2, '0')}.${d.getFullYear()} ${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}:${d.getSeconds().toString().padStart(2, '0')}`;
	};

	let isResultVisible = $state(true);
	let isCommandVisible = $state(true);

	type InputProps = {
		command: Command;
	};

	let { command = $bindable() }: InputProps = $props();
	let typeUppercase = $derived(command.type.toUpperCase());
	let statusUppercase = $derived(command.status.toUpperCase());
	let lastUpdateString = $derived.by(() => {
		if (command.last_result_update) {
			return formatDateAndTime(command.last_result_update);
		}
		return 'N/A';
	});

	let websocketStore: WebsocketStore;

	const handleStatusChange = (event: MessageEvent) => {
		const dataFromWs = JSON.parse(event.data);
		if (dataFromWs.message_type === 'command-status-updated') {
			if (dataFromWs.id !== command.id) return;
			command.status = dataFromWs.status;
		}
	};

	const handleResult = (event: MessageEvent) => {
		const dataFromWs = JSON.parse(event.data);
		if (dataFromWs.message_type === 'command-result') {
			if (dataFromWs.command_id !== command.id) return;
			command.result = dataFromWs.result;
		}
	};

	const handleResultChunk = (event: MessageEvent) => {
		const dataFromWs = JSON.parse(event.data);
		if (dataFromWs.message_type === 'result-chunk') {
			if (dataFromWs.command_id !== command.id) return;
			command.last_result_update = dataFromWs.created_at;
			command.result_progress = dataFromWs.progress;
		}
	};

	const cancelCommand = () => {
		websocketStore.ws?.send(
			JSON.stringify({
				type: 'cancel-command',
				commandId: command.id
			})
		);
	};

	onMount(() => {
		websocketStore = createWebsocketStore();

		websocketStore.ws?.addEventListener('message', handleStatusChange);
		websocketStore.ws?.addEventListener('message', handleResult);
		websocketStore.ws?.addEventListener('message', handleResultChunk);

		return () => {
			websocketStore.ws?.removeEventListener('message', handleStatusChange);
			websocketStore.ws?.removeEventListener('message', handleResult);
			websocketStore.ws?.removeEventListener('message', handleResultChunk);
		};
	});
</script>

<div class="border-neutral bg-base-100 mx-3 rounded-2xl border-2 p-4 shadow-xl">
	<div>
		<div class="mb-3 flex justify-between">
			<h2 class="card-title flex-col items-center justify-center md:flex-row">
				<span class="badge badge-primary">{typeUppercase}</span>
				<span class="badge badge-secondary">{statusUppercase}</span>
			</h2>
			<div class="flex flex-col justify-around space-y-1.5 md:flex-row md:space-x-1 md:space-y-0">
				<button class="btn btn-sm" on:click={() => (isCommandVisible = !isCommandVisible)}
					>Toggle command</button
				>
				<button
					class="btn btn-sm"
					class:btn-disabled={!command.result}
					on:click={() => (isResultVisible = !isResultVisible)}>Toggle result</button
				>
				<button
					class="btn btn-sm"
					class:btn-disabled={command.status !== 'created'}
					on:click={cancelCommand}>Cancel</button
				>
			</div>
		</div>
		<p>Operator: {command.operator}</p>
		<table class="w-full">
			<tbody>
				<tr>
					<td class="text-left">Created:</td>
					<td class="text-right">{command.created_at}</td>
				</tr>
				<tr>
					<td class="text-left">Last result update:</td>
					<td class="text-right">{lastUpdateString}</td>
				</tr>
				<tr>
					<td class="text-left">Result retrieval progress:</td>
					<td class="text-right">{command.result_progress ?? 0}%</td>
				</tr>
			</tbody>
		</table>
		<div class="w-full overflow-x-auto">
			{#if isCommandVisible}
				<div class="divider m-0"></div>
				<pre><code>{command.command}</code></pre>
			{/if}
		</div>
		<div class="w-full overflow-x-auto">
			{#if isResultVisible && command.result}
				<div class="divider m-0"></div>
				<pre><code>{command.result}</code></pre>
			{/if}
		</div>
	</div>
</div>
