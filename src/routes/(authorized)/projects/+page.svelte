<script lang="ts">
	let { data } = $props();

	const currentStatus = (startDate: string, endDate: string) => {
		const start = new Date(startDate).setHours(0, 0, 0, 0);
		const end = new Date(endDate).setHours(23, 59, 59, 999);
		const now = new Date();
		if (start > now) return 'UPCOMING';
		if (start <= now && end >= now) return 'ONGOING';
		return 'FINISHED';
	};
</script>

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
