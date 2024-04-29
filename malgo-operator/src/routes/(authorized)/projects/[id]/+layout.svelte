<script lang="ts">
	import Sidebar from '$lib/components/sidebar/Sidebar.svelte';
	import SidebarEntry from '$lib/components/sidebar/SidebarEntry.svelte';
	import { navigating, page } from '$app/stores';
	import { browser } from '$app/environment';
	import { get } from 'svelte/store';
	import { onDestroy, onMount } from 'svelte';
	import { afterNavigate } from '$app/navigation';
	import Session from '$lib/components/custom/session/Session.svelte';
	import SessionLoading from '$lib/components/custom/session/SessionLoading.svelte';
	import SessionNotFound from '$lib/components/custom/session/SessionNotFound.svelte';
	import { createWebsocketStore, type WebsocketStore } from '$lib/stores/Websocket';

	type Session = {
		id: string;
		name: string;
		createdAt: string;
		heartbeatAt: string;
		data: object;
	};

	let { data } = $props();

	let loaded = $state(false);
	let sessions: Session[] = $state([]);
	let uniqueSessionIds = $derived(new Set<string>(sessions.map((session) => session.id)));

	let project = data.project!;

	let currentUrl = $state(get(page).url.pathname);
	let currentSuffix = $derived(currentUrl.split('/')[4] ?? 'commands');
	afterNavigate(() => {
		currentUrl = get(page).url.pathname;
	});

	const formatDate = (date: string): string => {
		const d = new Date(date);
		// format DD.MM.YYYY
		return `${d.getDate().toString().padStart(2, '0')}.${(d.getMonth() + 1).toString().padStart(2, '0')}.${d.getFullYear()}`;
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

	const handleNewSession = async (event: MessageEvent) => {
		const dataFromWs = JSON.parse(event.data);
		if (dataFromWs.message_type === 'session-new') {
			if (uniqueSessionIds.has(dataFromWs.session_id)) return;
			uniqueSessionIds.add(dataFromWs.session_id);
			Array.prototype.unshift.apply(sessions, [
				{
					id: dataFromWs.session_id,
					name: dataFromWs.name,
					createdAt: dataFromWs.created_at,
					heartbeatAt: dataFromWs.heartbeat,
					data: {}
				}
			]);
		}
	};

	let unsubscribe: () => void;

	let websocketStore: WebsocketStore;

	onMount(async () => {
		drawerResize();
		window.addEventListener('resize', drawerResize);
		unsubscribe = navigating.subscribe(drawerResize);

		sessions = await data.sessions;
		loaded = true;

		websocketStore = createWebsocketStore();
		await websocketStore.subscribeToProject(get(page).params.id);
		websocketStore.ws?.addEventListener('close', async () => {
			await websocketStore.subscribeToProject(get(page).params.id);
		});

		websocketStore.ws?.addEventListener('message', handleNewSession);
	});

	onDestroy(async () => {
		if (browser) {
			window.removeEventListener('resize', drawerResize);
			unsubscribe();

			websocketStore.ws?.removeEventListener('message', handleNewSession);

			websocketStore.ws?.removeEventListener('close', async () => {
				await websocketStore.subscribeToProject(get(page).params.id);
			});
			await websocketStore.unsubscribeFromProject();
		}
	});
</script>

<div class="flex w-full grow flex-col">
	<div
		id="project-header"
		class="flex w-full justify-around border-b-2 border-base-300 bg-base-200 py-2"
	>
		<span class="mt-1">{formatDate(project.startDate)}</span>
		<a href={`/projects/${project.id}`} class="btn btn-ghost btn-sm">{project.name}</a>
		<span class="mt-1">{formatDate(project.endDate)}</span>
	</div>
	<div id="sidebar-wrap" class="max-w-full">
		<Sidebar>
			<svelte:fragment slot="content"><slot /></svelte:fragment>
			<svelte:fragment slot="menu">
				{#await data.sessions}
					<li>
						<SessionLoading />
					</li>
				{:catch error}
					<li>
						<p class="text-center text-base-content">Error: {error.message}</p>
					</li>
				{/await}
				{#if loaded}
					{#if sessions.length === 0}
						<li>
							<SessionNotFound />
						</li>
					{:else}
						{#each sessions as session}
							{@const isSelected = currentUrl.includes(`/projects/${project.id}/${session.id}`)}
							<SidebarEntry
								href={`/projects/${project.id}/${session.id}/${currentSuffix}`}
								active={isSelected}
							>
								<Session
									id={session.id}
									bind:name={session.name}
									bind:createdAt={session.createdAt}
									bind:heartbeatAt={session.heartbeatAt}
								/>
							</SidebarEntry>
						{/each}
					{/if}
				{/if}
			</svelte:fragment>
		</Sidebar>
	</div>
</div>
