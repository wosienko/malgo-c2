<script lang="ts">
	import { page } from '$app/stores';
	import { get } from 'svelte/store';
	// TODO: Update once page is officially updated
</script>

<svelte:head>
	<title>MALGO - OOPS</title>
</svelte:head>

<div class="flex flex-col items-center justify-center space-y-3">
	<h1 class="text-4xl font-extrabold leading-none tracking-tight md:text-5xl lg:text-6xl">
		{get(page).status}
	</h1>
	<p class="text-lg">{get(page).error?.message}</p>
	<a href="/" class="btn">Go back</a>
</div>
