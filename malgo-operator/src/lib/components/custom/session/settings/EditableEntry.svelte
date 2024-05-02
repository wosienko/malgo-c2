<script lang="ts">
	import ValidatedInput from '$lib/components/inputs/ValidatedInput.svelte';
	import { fieldSchema } from '$lib/validationSchemas';
	import { createWebsocketStore, type WebsocketStore } from '$lib/stores/Websocket';
	import { onMount } from 'svelte';
	import { get } from 'svelte/store';
	import { page } from '$app/stores';
	import { afterNavigate } from '$app/navigation';

	type InputProps = {
		key: string;
		ws_type: string;
		value: string;
		value_key: string;
	};

	function createSendType<T extends InputProps>(
		props: T,
		type: string,
		sessionId: string,
		value: string
	): { type: string; sessionId: string; [key: string]: string } {
		return {
			type,
			sessionId,
			[props.value_key]: value // Using value_key to set the dynamic property
		};
	}

	let { key, value, ws_type, value_key }: InputProps = $props();
	let disposable: InputProps = { key, value, ws_type, value_key };

	let websocketStore: WebsocketStore;

	// state
	let editableValue = $state(value);

	// validation
	const valueCheck = $derived(fieldSchema.safeParse(editableValue));

	let currentlyEdited = $state(false);
	const startEditing = () => {
		currentlyEdited = true;
	};
	const finishEditing = () => {
		currentlyEdited = false;

		websocketStore.ws?.send(
			JSON.stringify(createSendType(disposable, ws_type, get(page).params.sessionid, editableValue))
		);
	};
	const cancelEditing = () => {
		currentlyEdited = false;
		editableValue = value;
	};

	onMount(() => {
		websocketStore = createWebsocketStore();
	});

	afterNavigate(() => {
		editableValue = value;
	});
</script>

<tr>
	{#if !currentlyEdited}
		<td>{key}</td>
		<td>{value}</td>
		<td>
			<button class="btn btn-sm" on:click={startEditing}>Edit</button>
		</td>
	{:else}
		<td>
			{key}
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
		<td class="flex flex-col items-center space-y-3">
			<button
				class="btn btn-success btn-sm max-w-16"
				class:btn-disabled={!valueCheck.success || editableValue === value}
				on:click={finishEditing}>Save</button
			>
			<button class="btn btn-sm max-w-16" on:click={cancelEditing}>Cancel</button>
		</td>
	{/if}
</tr>
