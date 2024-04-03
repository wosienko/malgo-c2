<script lang="ts">
	import type { SafeParseReturnType } from 'zod';

	type InputProps = {
		id: string;
		name: string;
		type: string;
		value: string;
		validation: SafeParseReturnType<string, string>;
		classes?: string;
	};

	let {
		id,
		name,
		type,
		// eslint-disable-next-line no-undef
		value = $bindable(),
		// eslint-disable-next-line no-undef
		validation,
		classes = ''
	}: InputProps = $props();
</script>

{#if !validation.success}
	<p class="mb-1.5 text-xs text-error">
		{validation.error.errors[0].message.replace('String', '')}
	</p>
{:else}
	<p class="mb-1.5 text-xs text-transparent">For formatting sake</p>
{/if}
<input
	{type}
	{id}
	{name}
	autocomplete="off"
	class={`input input-bordered ${classes}`}
	bind:value
	class:input-error={!validation.success}
/>
