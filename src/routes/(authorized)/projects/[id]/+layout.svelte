<script lang="ts">
	import Sidebar from '$lib/components/sidebar/Sidebar.svelte';
	import SidebarEntry from '$lib/components/sidebar/SidebarEntry.svelte';
	import { page } from '$app/stores';
	import { get } from 'svelte/store';
	import { onMount } from 'svelte';
	import { afterNavigate } from '$app/navigation';

	let { data } = $props();

	let sessions = $state([]);

	let currentUrl = $state(get(page).url.pathname);
	afterNavigate(() => {
		currentUrl = get(page).url.pathname;
	});

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

	const formatDateAndTime = (date: string): string => {
		const d = new Date(date);
		// format DD.MM.YYYY HH:MM:SS
		return `${d.getDate().toString().padStart(2, '0')}.${(d.getMonth() + 1).toString().padStart(2, '0')}.${d.getFullYear()} ${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}:${d.getSeconds().toString().padStart(2, '0')}`;
	};

	onMount(async () => {
		sessions = await data.sessions;
	});
</script>

<div class="flex h-full w-full flex-col">
	<div class="flex w-full justify-around border-b-2 border-base-300 bg-base-200 py-2">
		<span class="mt-1">{formatDate(data.project.startDate)}</span>
		<a href={`/projects/${data.project.id}`} class="btn btn-ghost btn-sm">{data.project.name}</a>
		<span class="mt-1">{formatDate(data.project.endDate)}</span>
	</div>
	<div class="max-w-full grow overflow-y-auto">
		<Sidebar>
			<svelte:fragment slot="content"><slot /></svelte:fragment>
			<svelte:fragment slot="menu">
				{#if sessions.length === 0}
					<li>
						<div class="flex w-full flex-col">
							<div class="skeleton h-8 w-full"></div>
							<div class="skeleton h-4 w-full"></div>
							<div class="skeleton h-4 w-full"></div>
						</div>
					</li>
				{:else}
					{#each sessions as session}
						{@const isSelected = currentUrl.includes(`/projects/${data.project.id}/${session.id}`)}
						<SidebarEntry
							href={`/projects/${data.project.id}/${session.id}/commands`}
							active={isSelected}
						>
							<div class="flex flex-col">
								<span class="text-xl text-info">{session.name}</span>
								<table class="text-left">
									<tbody>
										<tr>
											<td>Created at:</td>
											<td>{formatDateAndTime(session.createdAt)}</td>
										</tr>
										<tr>
											<td>Last Heartbeat:</td>
											<td
												class={olderThanInSeconds(session.heartbeatAt, 3600)
													? 'text-error'
													: 'text-success'}>{formatDateAndTime(session.heartbeatAt)}</td
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
	</div>
</div>
