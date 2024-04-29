<script lang="ts">
	import { fade } from 'svelte/transition';
	import { createEventDispatcher, onMount } from 'svelte';

	export let message: string;

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
	class="alert toast alert-success toast-center toast-top mt-5 grid max-w-xs grid-flow-col"
	in:fade
	out:fade
>
	<div class="col-span-3 text-wrap">
		<p>
			{message}
		</p>
	</div>
	<div>
		<button on:click={() => dispatch('close')} class="btn-clear btn float-right">Hide</button>
	</div>
</div>
