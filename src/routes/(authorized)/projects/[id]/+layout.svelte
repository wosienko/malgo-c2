<script lang="ts">
	import Sidebar from '$lib/components/sidebar/Sidebar.svelte';
	import SidebarEntry from '$lib/components/sidebar/SidebarEntry.svelte';
	import { page } from '$app/stores';
	import { get } from 'svelte/store';

	let { data } = $props();

	const formatDate = (date: string): string => {
		const d = new Date(date);
		// format DD.MM.YYYY
		return `${d.getDate().toString().padStart(2, '0')}.${(d.getMonth() + 1).toString().padStart(2, '0')}.${d.getFullYear()}`;
	};
</script>

<div class="flex w-full justify-around border-b-2 border-base-300 bg-base-200 py-2">
	<span class="mt-1">{formatDate(data.project.startDate)}</span>
	<a href={get(page).url.pathname} class="btn btn-ghost btn-sm">{data.project.name}</a>
	<span class="mt-1">{formatDate(data.project.endDate)}</span>
</div>
<Sidebar>
	<div slot="content"><slot /></div>
	<svelte:fragment slot="menu">
		<SidebarEntry href="/">Entry 1</SidebarEntry>
	</svelte:fragment>
</Sidebar>
