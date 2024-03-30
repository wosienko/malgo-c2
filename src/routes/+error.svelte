<script lang="ts">
	import { page } from '$app/stores';
	import { get } from 'svelte/store';
</script>

<div class="flex flex-col items-center justify-center">
	<h1>{get(page).status}</h1>
	<p>{get(page).error?.message}</p>
	<a href="/">Go back</a>
</div>
