<script lang="ts">
	import { fade } from 'svelte/transition';
	import { createEventDispatcher, onMount } from 'svelte';
	import type { ApiError } from '$lib';

	export let issue: ApiError;

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
	out:fade
>
	<div class="col-span-3 text-wrap">
		<p>
			{issue.message}
		</p>
	</div>
	<div>
		<button on:click={() => dispatch('close')} class="btn-clear btn float-right">Hide</button>
	</div>
</div>
