<script lang="ts">
	import type { Command } from '$lib/components/custom/command/command';
	import { createWebsocketStore, type WebsocketStore } from '$lib/stores/Websocket';
	import { onMount } from 'svelte';

	let isResultVisible = $state(true);
	let isCommandVisible = $state(true);

	type InputProps = {
		command: Command;
	};

	let { command = $bindable() }: InputProps = $props();
	let typeUppercase = $derived(command.type.toUpperCase());
	let statusUppercase = $derived(command.status.toUpperCase());

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

		return () => {
			websocketStore.ws?.removeEventListener('message', handleStatusChange);
			websocketStore.ws?.removeEventListener('message', handleResult);
		};
	});
</script>

<div class="mx-3 rounded-2xl border-2 border-neutral bg-base-100 p-4 shadow-xl">
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
					<td class="text-right">{command.last_result_update ?? 'No result'}</td>
				</tr>
			</tbody>
		</table>
		{#if isCommandVisible}
			<div class="divider m-0"></div>
			<p class="whitespace-pre-line">{command.command}</p>
		{/if}
		{#if isResultVisible && command.result}
			<div class="divider m-0"></div>
			<p class="whitespace-pre-line">
				{command.result}
			</p>
		{/if}
	</div>
</div>
