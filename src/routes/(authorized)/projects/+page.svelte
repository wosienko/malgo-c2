<script lang="ts">
	import { pushState } from '$app/navigation';
	import { page } from '$app/stores';
	import { get } from 'svelte/store';
	import { onMount } from 'svelte';

	let { data } = $props();
	let count = $state(0);
	let allView = $state(get(page).url.searchParams.get('all') === 'true');
	let allPage = $state(Number(get(page).url.searchParams.get('page')) || 1);
	let allLimit = $state(Number(get(page).url.searchParams.get('limit')) || 6);

	const currentStatus = (startDate: string, endDate: string) => {
		const start = new Date(startDate).setHours(0, 0, 0, 0);
		const end = new Date(endDate).setHours(23, 59, 59, 999);
		const now = new Date();
		if (start > now) return 'UPCOMING';
		if (start <= now && end >= now) return 'ONGOING';
		return 'FINISHED';
	};

	const loadPage = (page: number) => {
		allPage = page;
		pushState(`/projects?all=${allView}&page=${allPage}&pageSize=${allLimit}`, {});
		data.projects = fetch(`/api/projects?all=${allView}&page=${allPage}&pageSize=${allLimit}`).then(
			(res) => res.json()
		);
	};

	const toggleAllProjects = () => {
		allView = !allView;
		pushState(`/projects?all=${allView}&page=${allPage}&pageSize=${allLimit}`, {});
		if (allView) {
			data.projects = fetch(
				`/api/projects?all=${allView}&page=${allPage}&pageSize=${allLimit}`
			).then((res) => res.json());
		} else {
			data.projects = fetch(`/api/projects?page=1&pageSize=999`).then((res) => res.json());
		}
	};

	onMount(async () => {
		count = await data.count;
	});
</script>

<div class="flex h-full w-full flex-col items-center">
	<button class="btn btn-primary" on:click={toggleAllProjects}>Toggle view of all projects</button>
	{#if allView}
		<div class="mt-3 flex justify-center space-x-3">
			<div class="join">
				<button
					class="btn join-item btn-sm"
					class:btn-disabled={allPage === 1}
					on:click={() => loadPage(allPage - 1)}>«</button
				>
				<button
					class="btn join-item btn-sm hover:cursor-auto hover:border-base-200 hover:bg-base-200"
					>Page {allPage}</button
				>
				<button
					class="btn join-item btn-sm"
					class:btn-disabled={allPage * allLimit >= count}
					on:click={() => loadPage(allPage + 1)}>»</button
				>
			</div>
			<div class="mt-1 md:flex-1">
				<div class="flex items-center justify-center space-x-3">
					<span>Limit: {allLimit}</span>
					<input
						type="range"
						min="1"
						max="20"
						bind:value={allLimit}
						class="range range-sm max-w-32"
						on:mouseup={() => loadPage(allPage)}
						on:touchend={() => loadPage(allPage)}
					/>
				</div>
			</div>
		</div>
	{/if}

	<div class="flex w-full flex-wrap justify-center">
		{#await data.projects}
			{#each [1, 2, 3, 4] as i}
				<div id={i} class="m-4 flex w-80 flex-col gap-4">
					<div class="skeleton h-28 w-full"></div>
					<div class="w-54 skeleton h-4"></div>
					<div class="skeleton h-4 w-full"></div>
					<div class="skeleton h-4 w-full"></div>
				</div>
			{/each}
		{:then projects}
			{#each projects as project}
				{@const status = currentStatus(project.startDate, project.endDate)}
				<div class="card m-4 w-80 border-2 border-neutral bg-base-100 shadow-xl">
					<div class="card-body">
						<h2 class="card-title">
							{project.name}
							{#if status === 'UPCOMING'}
								<span class="badge badge-warning">UPCOMING</span>
							{:else if status === 'ONGOING'}
								<span class="badge badge-primary">ONGOING</span>
							{:else}
								<span class="badge badge-secondary">FINISHED</span>
							{/if}
						</h2>
						<p class="whitespace-pre-line">{project.description}</p>
						<div class="card-actions justify-between">
							<div class="badge badge-outline">{project.startDate}</div>
							<div class="badge badge-outline">{project.endDate}</div>
						</div>
					</div>
				</div>
			{/each}
		{/await}
	</div>
</div>
