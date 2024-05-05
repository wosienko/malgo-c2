<script lang="ts">
	import { page } from '$app/stores';
	import { get } from 'svelte/store';
	import { browser } from '$app/environment';
	// TODO: Update once "page" store is officially updated
</script>

<svelte:head>
	<title>MALGO - OOPS</title>
</svelte:head>

<div class="flex flex-col items-center justify-center space-y-3">
	<h1 class="text-4xl font-extrabold leading-none tracking-tight md:text-5xl lg:text-6xl">
		{get(page).status}
	</h1>
	<p class="text-lg">{get(page).error?.message}</p>
	<button
		class="btn"
		on:click={() => {
			if (browser) window.history.back();
		}}>Go back</button
	>
</div>
