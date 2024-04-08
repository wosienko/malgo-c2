<script lang="ts">
	import Sidebar from '$lib/components/sidebar/Sidebar.svelte';
	import SidebarEntry from '$lib/components/sidebar/SidebarEntry.svelte';
	import { page } from '$app/stores';
	import { get } from 'svelte/store';
	import { onMount } from 'svelte';

	let { data } = $props();

	let sessions = $state([]);

	const formatDate = (date: string): string => {
		const d = new Date(date);
		// format DD.MM.YYYY
		return `${d.getDate().toString().padStart(2, '0')}.${(d.getMonth() + 1).toString().padStart(2, '0')}.${d.getFullYear()}`;
	};

	const olderThanInSeconds = (date: string, seconds: number): boolean => {
		const d = new Date(date);
		const now = new Date();
		return (now.getTime() - d.getTime()) / 1000 > seconds;
	};

	onMount(async () => {
		sessions = await data.sessions;
	});
</script>

<div class="flex w-full justify-around border-b-2 border-base-300 bg-base-200 py-2">
	<span class="mt-1">{formatDate(data.project.startDate)}</span>
	<a href={get(page).url.pathname} class="btn btn-ghost btn-sm">{data.project.name}</a>
	<span class="mt-1">{formatDate(data.project.endDate)}</span>
</div>
<Sidebar>
	<div slot="content"><slot /></div>
	<svelte:fragment slot="menu">
		{#if sessions.length === 0}
			<li>Sessions loading</li>
		{:else}
			{#each sessions as session}
				{@const isSelected =
					get(page).url.pathname === `/projects/${data.project.id}/${session.id}`}
				<SidebarEntry href={`/projects/${data.project.id}/${session.id}`} active={isSelected}>
					<div class="flex flex-col" class:bg-neutral={isSelected}>
						<span class="text-xl text-info">{session.name}</span>
						<table class="text-left" class:text-neutral-content={isSelected}>
							<tbody>
								<tr>
									<td>Created at:</td>
									<td>{session.createdAt}</td>
								</tr>
								<tr>
									<td>Last Heartbeat:</td>
									<td
										class={olderThanInSeconds(session.heartbeatAt, 3600)
											? 'text-error'
											: 'text-success'}>{session.heartbeatAt}</td
									>
								</tr>
							</tbody>
						</table>
					</div>
				</SidebarEntry>
			{/each}
		{/if}
	</svelte:fragment>
</Sidebar>
