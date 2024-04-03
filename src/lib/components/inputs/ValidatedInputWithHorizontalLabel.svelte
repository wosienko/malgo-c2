<script lang="ts">
	import ValidatedInput from '$lib/components/inputs/ValidatedInput.svelte';
	import type { SafeParseReturnType } from 'zod';

	type Props = {
		id: string;
		label: string;
		name: string;
		type: string;
		value: string;
		validation: SafeParseReturnType<string, string>;
		classes?: string;
	};

	// eslint-disable-next-line no-undef
	let { id, label, name, type, value = $bindable(), validation, classes = '' }: Props = $props();
</script>

<div class="flex items-center justify-between">
	<label for={id}>{label}</label>
	<div class="-mt-4 flex flex-col items-center">
		<ValidatedInput {id} {name} {type} bind:value {validation} {classes} />
	</div>
</div>
