<script lang="ts">
	import Sidebar from '$lib/components/sidebar/Sidebar.svelte';
	import SidebarEntry from '$lib/components/sidebar/SidebarEntry.svelte';
	import { navigating, page } from '$app/stores';
	import { browser } from '$app/environment';
	import { get } from 'svelte/store';
	import { onDestroy, onMount } from 'svelte';
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

	const drawerResize = () => {
		// set height to the height of the window
		let headerHeight = document.getElementById('header')!.getBoundingClientRect();
		let footerHeight = document.getElementById('footer')!.getBoundingClientRect();
		let projectHeaderHeight = document.getElementById('project-header')!.getBoundingClientRect();
		let target1 = document.querySelector('.drawer-content') as HTMLElement;
		let target2 = document.querySelector('.drawer-side') as HTMLElement;

		let width = window.innerWidth > 0 ? window.innerWidth : screen.width;

		target1.style.height = `calc(100dvh - ${headerHeight.height + footerHeight.height + projectHeaderHeight.height}px)`;

		if (width < 1024) {
			target2.style.height = `100dvh`;
		} else {
			target2.style.height = `calc(100dvh - ${headerHeight.height + footerHeight.height + projectHeaderHeight.height}px)`;
		}
	};

	let unsubscribe: () => void;

	onMount(async () => {
		sessions = await data.sessions;

		drawerResize();
		window.addEventListener('resize', drawerResize);
		unsubscribe = navigating.subscribe(drawerResize);
	});

	onDestroy(() => {
		if (browser) {
			window.removeEventListener('resize', drawerResize);
			unsubscribe();
		}
	});
</script>

<div class="flex w-full grow flex-col">
	<div
		id="project-header"
		class="flex w-full justify-around border-b-2 border-base-300 bg-base-200 py-2"
	>
		<span class="mt-1">{formatDate(data.project.startDate)}</span>
		<a href={`/projects/${data.project.id}`} class="btn btn-ghost btn-sm">{data.project.name}</a>
		<span class="mt-1">{formatDate(data.project.endDate)}</span>
	</div>
	<div id="sidebar-wrap" class="max-w-full">
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
