<script lang="ts">
	import type { UserWithRoles } from '$lib';
	import type { ApiError } from '$lib';
	import type { ZodIssue } from 'zod';
	import { pushState } from '$app/navigation';
	import ZodIssues from '$lib/components/toasts/ZodIssues.svelte';
	import ApiIssues from '$lib/components/toasts/ApiIssues.svelte';
	import ChangeSuccessful from '$lib/components/toasts/ChangeSuccessful.svelte';
	import { onMount } from 'svelte';
	import UserTableEntryLoading from '$lib/components/custom/admin/projects/UserTableEntryLoading.svelte';
	import UserTableEntry from '$lib/components/custom/admin/users/UserTableEntry.svelte';
	import NewUserModal from '$lib/components/custom/admin/users/NewUserModal.svelte';
	import DeleteUserModal from '$lib/components/custom/admin/users/DeleteUserModal.svelte';
	import PasswordChangeModal from '$lib/components/custom/admin/users/PasswordChangeModal.svelte';

	let { data } = $props();

	let count = $state(0);

	let page = $state(data.page);
	let pageSize = $state(data.pageSize);

	let zodIssues: ZodIssue[] = $state([]);
	let apiError: ApiError | undefined = $state();
	let successMessage: string | undefined = $state();

	const loadPage = (nextPage: number) => {
		return async () => {
			data.users = fetch(`/api/user?page=${nextPage}&pageSize=${pageSize}`).then((res) =>
				res.json()
			);
			pushState(`/admin/users?page=${nextPage}&pageSize=${pageSize}`, {});

			count = (await data.users).count;
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

	let userToAlter: UserWithRoles = $state({
		id: '',
		name: '',
		surname: '',
		email: '',
		admin: false,
		operator: false
	});

	let showDeletionModal = $state(() => {});
	const prepareForDeletion = (user: typeof userToAlter) => {
		return () => {
			userToAlter = user;
			showDeletionModal();
		};
	};

	let showPasswordChangeModal = $state(() => {});
	const prepareForPasswordChange = (user: typeof userToAlter) => {
		return () => {
			userToAlter = user;
			showPasswordChangeModal();
		};
	};

	let showRegistrationModal = $state(() => {});
	const prepareForRegistration = () => {
		showRegistrationModal();
	};

	onMount(async () => {
		count = (await data.users).count;
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

<PasswordChangeModal
	{userToAlter}
	bind:showModal={showPasswordChangeModal}
	reloadCurrentPage={async () => {
		await loadPage(page)();
	}}
	bind:successMessage
	bind:apiError
	bind:zodIssues
/>

<DeleteUserModal
	{userToAlter}
	bind:showModal={showDeletionModal}
	reloadCurrentPage={async () => {
		await loadPage(page)();
	}}
	bind:successMessage
	bind:apiError
	bind:zodIssues
/>

<NewUserModal
	bind:successMessage
	bind:apiError
	bind:zodIssues
	reloadCurrentPage={async () => {
		await loadPage(page)();
	}}
	bind:showModal={showRegistrationModal}
/>

<div class="mt-3 flex w-full justify-center space-x-3">
	<div class="hidden flex-1 text-center md:block">
		<button class="btn btn-sm -mt-1" on:click={prepareForRegistration}>Register new user</button>
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
			<th class="hidden md:table-cell md:w-1/5 lg:w-[23%]">Name</th>
			<th class="hidden md:table-cell md:w-1/5 lg:w-[23%]">Surname</th>
			<th class="md:w-1/5 lg:w-[23%]">Email</th>
			<th class="w-20 md:w-24">Admin</th>
			<th class="w-20 md:w-24">Operator</th>
			<th class="w-28 md:w-auto"></th>
		</tr>
	</thead>
	<tbody>
		{#await data.users}
			<UserTableEntryLoading />
		{:then users}
			<!--eslint-disable-next-line @typescript-eslint/no-unused-vars-->
			{#each users.users as _, i}
				<UserTableEntry
					bind:user={users.users[i]}
					elementIndex={i}
					{prepareForPasswordChange}
					{prepareForDeletion}
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
