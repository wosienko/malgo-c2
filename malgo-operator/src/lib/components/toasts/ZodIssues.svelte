<script lang="ts">
	import type { ZodIssue } from 'zod';
	import { fade } from 'svelte/transition';
	import { createEventDispatcher, onMount } from 'svelte';

	export let issues: ZodIssue[];

	const dispatch = createEventDispatcher();

	onMount(() => {
		let timeout = setTimeout(() => {
			dispatch('close');
		}, 3000);

		return () => {
			clearTimeout(timeout);
		};
	});
</script>

<div
	class="alert toast alert-error toast-center toast-top mt-5 grid max-w-xs grid-flow-col"
	in:fade
	out:fade
>
	<div class="col-span-3 text-wrap">
		{#each issues as { message, path }}
			<p>
				{String(path[0]).charAt(0).toUpperCase() + String(path[0]).slice(1)} -
				{message}
			</p>
		{/each}
	</div>
	<div>
		<button on:click={() => dispatch('close')} class="btn-clear btn float-right">Hide</button>
	</div>
</div>
