<script lang="ts">
	import type { ZodIssue } from 'zod';
	import type { ApiError, UserWithRoles } from '$lib';
	import ZodIssues from '$lib/components/toasts/ZodIssues.svelte';
	import ApiIssues from '$lib/components/toasts/ApiIssues.svelte';
	import ChangeSuccessful from '$lib/components/toasts/ChangeSuccessful.svelte';
	import { pushState } from '$app/navigation';
	import { onMount } from 'svelte';
	import ProjectTableEntryLoading from '$lib/components/custom/admin/projects/ProjectTableEntryLoading.svelte';
	import ProjectTableEntry from '$lib/components/custom/admin/projects/ProjectTableEntry.svelte';
	import type { ProjectTableEntryType } from '$lib/components/custom/admin/projects/types';
	import NewProjectModal from '$lib/components/custom/admin/projects/NewProjectModal.svelte';
	import DeleteProjectModal from '$lib/components/custom/admin/projects/DeleteProjectModal.svelte';
	import AssignOperatorsModal from '$lib/components/custom/admin/projects/AssignOperatorsModal.svelte';

	let zodIssues: ZodIssue[] = $state([]);
	let apiError: ApiError | undefined = $state();
	let successMessage: string | undefined = $state();

	let { data } = $props();

	let page = $state(data.page);
	let pageSize = $state(data.pageSize);

	let count = $state(0);

	let projectToAlter: ProjectTableEntryType = $state({
		id: '',
		name: '',
		startDate: '',
		endDate: '',
		description: ''
	});

	const loadPage = (nextPage: number) => {
		return async () => {
			data.projects = fetch(`/api/project?page=${nextPage}&pageSize=${pageSize}`).then((res) =>
				res.json()
			);
			pushState(`/admin/projects?page=${nextPage}&pageSize=${pageSize}`, {});

			count = (await data.projects).count;
			page = nextPage;
		};
	};

	const loadNextPage = async () => {
		await loadPage(page + 1)();
	};
	const reloadCurrentPage = async () => {
		await loadPage(page)();
	};
	const loadPreviousPage = async () => {
		await loadPage(page - 1)();
	};

	let showProjectCreationModal = $state(() => {});

	let showDeletionModal = $state(() => {});
	const prepareForDeletion = (project: typeof projectToAlter) => {
		return () => {
			projectToAlter = project;
			showDeletionModal();
		};
	};

	let allOperators = $state<UserWithRoles[]>([]);
	let showOperatorAssignmentModal = $state(() => {});
	const prepareForOperatorAssignment = (project: ProjectTableEntryType) => {
		return async () => {
			projectToAlter = project;
			showOperatorAssignmentModal();
		};
	};

	onMount(async () => {
		count = (await data.projects).count;
		allOperators = await data.users;
	});
</script>

{#if zodIssues.length > 0}
	<ZodIssues
		issues={zodIssues}
		on:close={() => {
			zodIssues = [];
		}}
	/>
{/if}

{#if apiError}
	<ApiIssues
		issue={apiError}
		on:close={() => {
			apiError = undefined;
		}}
	/>
{/if}

{#if successMessage}
	<ChangeSuccessful
		message={successMessage}
		on:close={() => {
			successMessage = undefined;
		}}
	/>
{/if}

<DeleteProjectModal
	{projectToAlter}
	bind:showModal={showDeletionModal}
	reloadCurrentPage={async () => {
		await loadPage(page)();
	}}
	bind:successMessage
	bind:apiError
	bind:zodIssues
/>

<AssignOperatorsModal
	{projectToAlter}
	users={allOperators}
	bind:showModal={showOperatorAssignmentModal}
	bind:successMessage
	bind:apiError
	bind:zodIssues
/>

<NewProjectModal
	bind:successMessage
	bind:apiError
	bind:zodIssues
	reloadCurrentPage={async () => {
		await loadPage(page)();
	}}
	bind:showModal={showProjectCreationModal}
/>

<div class="mt-3 flex w-full justify-center space-x-3">
	<div class="hidden flex-1 text-center md:block">
		<button class="btn btn-sm -mt-1" on:click={showProjectCreationModal}>Create new project</button>
	</div>
	<div class="join -mt-1">
		<button class="btn join-item btn-sm" class:btn-disabled={page === 1} on:click={loadPreviousPage}
			>«</button
		>
		<button class="btn join-item btn-sm hover:cursor-auto hover:border-base-200 hover:bg-base-200"
			>Page {page}</button
		>
		<button
			class="btn join-item btn-sm"
			class:btn-disabled={page * pageSize >= count}
			on:click={loadNextPage}>»</button
		>
	</div>
	<div class="mt-1 md:flex-1">
		<div class="flex items-center justify-center space-x-3">
			<span>Limit: {pageSize}</span>
			<input
				type="range"
				min="1"
				max="20"
				bind:value={pageSize}
				class="range range-sm max-w-32"
				on:mouseup={reloadCurrentPage}
				on:touchend={reloadCurrentPage}
			/>
		</div>
	</div>
</div>

<table class="table table-fixed overflow-x-auto text-center">
	<thead>
		<tr>
			<th>Name</th>
			<th class="hidden md:table-cell">Start date</th>
			<th class="hidden md:table-cell">End date</th>
			<th class="hidden md:table-cell">Description</th>
			<th class="w-28 md:w-auto"></th>
		</tr>
	</thead>
	<tbody>
		{#await data.projects}
			<ProjectTableEntryLoading />
		{:then projects}
			<!--eslint-disable-next-line @typescript-eslint/no-unused-vars-->
			{#each projects.projects as _, i}
				<ProjectTableEntry
					bind:project={projects.projects[i]}
					elementIndex={i}
					{prepareForDeletion}
					{prepareForOperatorAssignment}
					{count}
					{pageSize}
					bind:successMessage
					bind:apiError
					bind:zodIssues
				/>
			{/each}
		{/await}
	</tbody>
</table>
