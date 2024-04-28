<script lang="ts">
	import ValidatedInput from '$lib/components/inputs/ValidatedInput.svelte';
	import { fieldSchema } from '$lib/validationSchemas';
	import { createWebsocketStore, type WebsocketStore } from '$lib/stores/Websocket';
	import { onMount } from 'svelte';
	import { get } from 'svelte/store';
	import { page } from '$app/stores';

	type InputProps = {
		key: string;
		value: string;
		keysInUse: string[];
	};

	let { key, value, keysInUse }: InputProps = $props();

	let websocketStore: WebsocketStore;

	// state
	let editableKey = $state(key);
	let editableValue = $state(value);

	// validation
	const keyCheck = $derived(
		fieldSchema
			.refine((value) => !keysInUse.includes(value) || value === key, { message: 'Key in use' })
			.safeParse(editableKey)
	);
	const valueCheck = $derived(fieldSchema.safeParse(editableValue));

	let currentlyEdited = $state(false);
	const startEditing = () => {
		currentlyEdited = true;
	};
	const finishEditing = () => {
		currentlyEdited = false;
		if (key === editableKey && value === editableValue) {
			return;
		}
		if (key !== editableKey) {
			deleteItem();
			websocketStore.ws?.send(
				JSON.stringify({
					type: 'session-new-key',
					sessionId: get(page).params.sessionid,
					key: editableKey,
					value: editableValue
				})
			);
		} else {
			websocketStore.ws?.send(
				JSON.stringify({
					type: 'session-update-key',
					sessionId: get(page).params.sessionid,
					key: editableKey,
					value: editableValue
				})
			);
		}
	};
	const cancelEditing = () => {
		currentlyEdited = false;
		editableKey = key;
		editableValue = value;
	};

	const deleteItem = () => {
		websocketStore.ws?.send(
			JSON.stringify({
				type: 'session-delete-key',
				sessionId: get(page).params.sessionid,
				key: key
			})
		);
	};

	onMount(() => {
		websocketStore = createWebsocketStore();
	});
</script>

<tr>
	{#if !currentlyEdited}
		<td>{key}</td>
		<td>{value}</td>
		<td class="flex flex-col space-y-2">
			<button class="btn btn-sm" on:click={startEditing}>Edit</button>
			<button class="btn btn-error btn-sm" on:click={deleteItem}>Del</button>
		</td>
	{:else}
		<td>
			<div class="md:-mt-1">
				<ValidatedInput
					type="text"
					id="name"
					name="name"
					bind:value={editableKey}
					validation={keyCheck}
					classes="input-sm mb-5 w-full"
				/>
			</div>
		</td>
		<td>
			<div class="md:-mt-1">
				<ValidatedInput
					type="text"
					id="name"
					name="name"
					bind:value={editableValue}
					validation={valueCheck}
					classes="input-sm mb-5 w-full"
				/>
			</div>
		</td>
		<td class="flex flex-col space-y-2">
			<button
				class="btn btn-success btn-sm"
				class:btn-disabled={!keyCheck.success ||
					!valueCheck.success ||
					(editableKey === key && editableValue === value)}
				on:click={finishEditing}>Save</button
			>
			<button class="btn btn-sm" on:click={cancelEditing}>Cancel</button>
		</td>
	{/if}
</tr>
