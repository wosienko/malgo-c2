<script lang="ts">
	import EditableKeyValueEntry from '$lib/components/custom/session/key-value/EditableKeyValueEntry.svelte';
	import ValidatedInputWithLabel from '$lib/components/inputs/ValidatedInputWithHorizontalLabel.svelte';
	import ModalRunCancel from '$lib/components/modals/ModalRunCancel.svelte';
	import { fieldSchema } from '$lib/validationSchemas';
	import { createWebsocketStore, type WebsocketStore } from '$lib/stores/Websocket';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { get } from 'svelte/store';

	let websocketStore: WebsocketStore;

	type InputProps = {
		values: Map<string, string>;
	};

	let { values }: InputProps = $props();

	let keyValues = $state(values);
	let keysInUse: string[] = $derived(Array.from(keyValues.keys()));

	let showNewKeyModal = $state(() => {});

	let newKey = $state({
		key: '',
		value: ''
	});
	let newKeyValidation = $derived({
		key: fieldSchema
			.refine((value) => !keyValues.has(value), { message: 'Key already exists' })
			.safeParse(newKey.key),
		value: fieldSchema.safeParse(newKey.value)
	});
	let newKeyValid = $derived(newKeyValidation.key.success && newKeyValidation.value.success);

	const createNewKey = async (): Promise<boolean> => {
		websocketStore.ws?.send(
			JSON.stringify({
				type: 'session-new-key',
				sessionId: get(page).params.sessionid,
				key: newKey.key,
				value: newKey.value
			})
		);
		newKey = { key: '', value: '' };
		return true;
	};

	const listenForKeyValueChanges = (event: MessageEvent) => {
		const data = JSON.parse(event.data);
		if (data.message_type === 'session-key-value') {
			keyValues.set(data.key, data.value);
			keyValues = new Map(keyValues); // trigger reactivity
		}
	};

	const listenForKeyValueDeletes = (event: MessageEvent) => {
		const data = JSON.parse(event.data);
		if (data.message_type === 'session-key-value-delete') {
			keyValues.delete(data.key);
			keyValues = new Map(keyValues); // trigger reactivity
		}
	};

	onMount(() => {
		websocketStore = createWebsocketStore();

		websocketStore.ws?.addEventListener('message', listenForKeyValueChanges);
		websocketStore.ws?.addEventListener('message', listenForKeyValueDeletes);

		return () => {
			websocketStore.ws?.removeEventListener('message', listenForKeyValueChanges);
			websocketStore.ws?.removeEventListener('message', listenForKeyValueDeletes);
		};
	});
</script>

<ModalRunCancel
	id="new-session-key"
	title="New session key"
	message="Session key: "
	messageEmphasis=""
	btnClass="btn-success"
	btnText="Add key"
	btnDisabledCondition={!newKeyValid}
	onclickCallback={createNewKey}
	bind:showModal={showNewKeyModal}
	onHideModal={() => {
		newKey = { key: '', value: '' };
	}}
>
	<div class="space-y-8">
		<ValidatedInputWithLabel
			label="Key"
			type="text"
			id="new-key"
			name="key"
			bind:value={newKey.key}
			validation={newKeyValidation.key}
			classes="w-full max-w-xs"
		/>
		<ValidatedInputWithLabel
			label="Value"
			type="text"
			id="new-value"
			name="value"
			bind:value={newKey.value}
			validation={newKeyValidation.value}
			classes="w-full max-w-xs"
		/>
	</div>
</ModalRunCancel>

<table class="table table-fixed">
	<thead>
		<tr>
			<th class="w-auto">Key</th>
			<th class="w-auto">Value</th>
			<th class="w-24">
				<button class="btn btn-sm w-full" on:click={showNewKeyModal}>Add</button>
			</th>
		</tr>
	</thead>
	<tbody>
		{#each keyValues as val}
			{@const key = val[0]}
			{@const value = val[1]}
			<EditableKeyValueEntry {key} {value} {keysInUse} />
		{/each}
	</tbody>
</table>
