<script lang="ts">
	import Command from '$lib/components/custom/command/Command.svelte';
	import type { Command as CommandType } from '$lib/components/custom/command/command';
	import { onDestroy, onMount } from 'svelte';
	import { createWebsocketStore, type WebsocketStore } from '$lib/stores/Websocket.js';
	import { page } from '$app/stores';
	import { get } from 'svelte/store';
	import { browser } from '$app/environment';

	let newestCommand: CommandType = $state({
		id: '',
		type: 'cmd',
		status: 'aaa',
		operator: 'bajo jajo',
		command: 'string',
		result: 'nic',
		created_at: 'today'
	});

	let commandToBeSent = $state('');

	let websocketStore: WebsocketStore;

	const sendNewCommand = () => {
		websocketStore.sendCommand(get(page).params.sessionid, commandToBeSent);
		commandToBeSent = '';
	};

	const handleNewCommand = (event: MessageEvent) => {
		try {
			newestCommand = JSON.parse(event.data);
		} catch (error) {
			console.error('Error parsing websocket message', error);
		}
	};

	onMount(async () => {
		websocketStore = createWebsocketStore();

		websocketStore.ws?.addEventListener('message', handleNewCommand);
	});

	onDestroy(() => {
		if (browser) {
			websocketStore.ws?.removeEventListener('message', handleNewCommand);
		}
	});
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
		<table class="table table-fixed">
			<thead>
				<tr>
					<th class="w-auto">Key</th>
					<th class="w-auto">Value</th>
					<th class="w-24">
						<button class="btn btn-sm">Add</button>
					</th>
				</tr>
			</thead>
			<tbody>
				<tr>
					<td>key</td>
					<td>value</td>
					<td class="flex flex-col items-center justify-center space-y-2">
						<button class="btn btn-sm">Edit</button>
						<button class="btn btn-error btn-sm">Del</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td class="flex flex-col items-center justify-center space-y-2">
						<button class="btn btn-sm">Edit</button>
						<button class="btn btn-error btn-sm">Del</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td class="flex flex-col items-center justify-center space-y-2">
						<button class="btn btn-sm">Edit</button>
						<button class="btn btn-error btn-sm">Del</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td class="flex flex-col items-center justify-center space-y-2">
						<button class="btn btn-sm">Edit</button>
						<button class="btn btn-error btn-sm">Del</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td class="flex flex-col items-center justify-center space-y-2">
						<button class="btn btn-sm">Edit</button>
						<button class="btn btn-error btn-sm">Del</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td class="flex flex-col items-center justify-center space-y-2">
						<button class="btn btn-sm">Edit</button>
						<button class="btn btn-error btn-sm">Del</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td class="flex flex-col items-center justify-center space-y-2">
						<button class="btn btn-sm">Edit</button>
						<button class="btn btn-error btn-sm">Del</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td class="flex flex-col items-center justify-center space-y-2">
						<button class="btn btn-sm">Edit</button>
						<button class="btn btn-error btn-sm">Del</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td class="flex flex-col items-center justify-center space-y-2">
						<button class="btn btn-sm">Edit</button>
						<button class="btn btn-error btn-sm">Del</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td class="flex flex-col items-center justify-center space-y-2">
						<button class="btn btn-sm">Edit</button>
						<button class="btn btn-error btn-sm">Del</button>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</div>
<div class="h-1/2 w-full overflow-y-auto py-3">
	<Command command={newestCommand} />
</div>
