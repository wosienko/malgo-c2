<script lang="ts">
	import type { UsersWithRoles, UserWithRoles } from '$lib/db/schema/users';
	import type { ApiError } from '$lib';
	import { emailSchema, fieldSchema } from '$lib/validationSchemas';
	import type { ZodIssue } from 'zod';
	import ZodIssues from '$lib/components/ZodIssues.svelte';
	import ApiIssues from '$lib/components/ApiIssues.svelte';

	let { data } = $props();

	let users = $state(
		data.users.users.map((user) => ({
			...user,
			editing: false
		}))
	);

	let lastUserValues: (typeof users)[0] = $state((() => users[0])());

	let page = $state(data.page);
	let pageSize = $state(data.pageSize);

	let zodIssues: ZodIssue[] = $state([]);
	let apiError: ApiError | undefined = $state();

	const loadPage = (nextPage: number) => {
		return async () => {
			const res = await fetch(`/api/user?page=${nextPage}&pageSize=${pageSize}`);
			const usersPage: UsersWithRoles = await res.json();
			users = usersPage.users.map((user) => ({
				...user,
				editing: false
			}));
			data.users.count = usersPage.count;
			page = nextPage;
		};
	};

	const loadNextPage = () => {
		loadPage(page + 1)();
	};
	const reloadCurrentPage = () => {
		loadPage(page)();
	};
	const loadPreviousPage = () => {
		loadPage(page - 1)();
	};

	const updateUser = async (user: (typeof users)[0]) => {
		const res = await fetch(`/api/user/${user.id}`, {
			method: 'PATCH',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				name: user.name,
				surname: user.surname,
				email: user.email,
				admin: user.admin,
				operator: user.operator
			})
		});

		if (res.ok) {
			user.editing = false;
		} else {
			try {
				let body: ApiError | ZodIssue[] = await res.json();

				if (Array.isArray(body)) {
					zodIssues = body;
				} else {
					apiError = body;
				}
			} catch (e) {
				console.error(e);
			}
		}
	};

	let userToDelete: (typeof users)[0] = $state((() => users[0])());

	const deleteUser = async (user: (typeof users)[0]) => {
		const res = await fetch(`/api/user/${user.id}`, {
			method: 'DELETE'
		});

		if (res.ok) {
			let newUsers = await fetch(`/api/user?page=${page}&pageSize=${pageSize}`).then((res) =>
				res.json()
			);
			users = newUsers.users.map((user: UserWithRoles) => ({
				...user,
				editing: false,
				isEmailValid: emailSchema.safeParse(user.email),
				isNameValid: fieldSchema.safeParse(user.name),
				isSurnameValid: fieldSchema.safeParse(user.surname)
			}));
			data.users.count = newUsers.count;
		} else {
			try {
				let body: ApiError | ZodIssue[] = await res.json();

				if (Array.isArray(body)) {
					zodIssues = body;
				} else {
					apiError = body;
				}
			} catch (e) {
				console.error(e);
			}
		}
	};

	const prepareForDeletion = (user: typeof userToDelete) => {
		return () => {
			userToDelete = user;
			const dialog = document.getElementById('confirm-deletion') as HTMLDialogElement;
			dialog.showModal();
		};
	};
</script>

<svelte:head>
	<title>MALGO - Admin</title>
</svelte:head>

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

<dialog id="confirm-deletion" class="modal modal-bottom sm:modal-middle">
	<div class="modal-box">
		<h3 class="text-lg font-bold">Are you sure you want to delete?</h3>
		<p class="py-4">
			Account to be deleted: <span class="font-bold"
				>{userToDelete.name} {userToDelete.surname}</span
			>
		</p>
		<div class="modal-action">
			<form method="dialog">
				<button
					class="btn btn-error"
					on:click={() => {
						deleteUser(userToDelete);
					}}>Delete</button
				>
				<button class="btn">Cancel</button>
			</form>
		</div>
	</div>
</dialog>

<div class="mt-3 flex w-full justify-center space-x-3">
	<div class="hidden flex-1 text-center md:block">
		<a href="/admin/users/register" class="btn btn-sm -mt-1">Register new user</a>
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
			class:btn-disabled={page * pageSize >= data.users.count}
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

<p class="my-3 w-full text-center md:hidden">Users may be edited on larger screens!</p>

<table class="table table-fixed overflow-x-auto text-center">
	<thead>
		<tr>
			<th class="md:w-1/5 lg:w-[23%]">Name</th>
			<th class="md:w-1/5 lg:w-[23%]">Surname</th>
			<th class="md:w-1/5 lg:w-[23%]">Email</th>
			<th class="w-20 md:w-24">Admin</th>
			<th class="w-20 md:w-24">Operator</th>
			<th class="hidden md:block"></th>
		</tr>
	</thead>
	<tbody>
		{#each users as user, i}
			<tr class="max-w-dvw hover">
				{#if !user.editing}
					<td class="no-scrollbar overflow-x-auto">{user.name}</td>
					<td class="no-scrollbar overflow-x-auto">{user.surname}</td>
					<td class="no-scrollbar overflow-x-auto">{user.email}</td>
					<td>
						<input
							type="checkbox"
							checked={user.admin}
							disabled
							class="checkbox-info checkbox checkbox-sm"
						/>
					</td>
					<td>
						<input
							type="checkbox"
							checked={user.operator}
							disabled
							class="checkbox-info checkbox checkbox-sm"
						/>
					</td>
					<td class="hidden md:block">
						<div
							class="dropdown dropdown-end"
							class:dropdown-bottom={i <= 2 ||
								i + 1 < (data.users.count > pageSize ? pageSize : data.users.count) - 2}
							class:dropdown-top={i > 2 &&
								i + 1 >= (data.users.count > pageSize ? pageSize : data.users.count) - 2}
						>
							<div tabindex="-1" role="button" class="btn btn-neutral btn-sm mb-1">Hover</div>
							<ul
								tabindex="-1"
								class="menu dropdown-content z-[1] w-52 space-y-1.5 rounded-box bg-base-100 p-2 shadow"
							>
								<li>
									<button
										class="btn btn-sm"
										on:click={() => {
											user.editing = true;
											lastUserValues = { ...user };
										}}>Edit</button
									>
								</li>
								<li><button class="btn btn-sm">Change password</button></li>
								<li>
									<button
										class="btn btn-error btn-sm"
										on:click={() => {
											prepareForDeletion(user)();
										}}>Delete</button
									>
								</li>
							</ul>
						</div>
					</td>
				{:else}
					{@const nameCheck = fieldSchema.safeParse(user.name)}
					{@const surnameCheck = fieldSchema.safeParse(user.surname)}
					{@const emailCheck = emailSchema.safeParse(user.email)}
					<td>
						{#if !nameCheck.success}
							<p class="mb-1.5 text-xs text-error">
								{nameCheck.error.errors[0].message.replace('String', '')}
							</p>
						{:else}
							<p class="mb-1.5 text-xs text-transparent">For formatting sake</p>
						{/if}
						<input
							type="text"
							id="name"
							name="name"
							autocomplete="off"
							required
							class="input input-sm input-bordered mb-5 w-full"
							bind:value={user.name}
							class:input-error={!nameCheck.success}
						/>
					</td>
					<td>
						{#if !surnameCheck.success}
							<p class="mb-1.5 text-xs text-error">
								{surnameCheck.error.errors[0].message.replace('String', '')}
							</p>
						{:else}
							<p class="mb-1.5 text-xs text-transparent">For formatting sake</p>
						{/if}
						<input
							type="text"
							id="surname"
							name="surname"
							autocomplete="off"
							required
							class="input input-sm input-bordered mb-5 w-full"
							bind:value={user.surname}
							class:input-error={!surnameCheck.success}
						/>
					</td>
					<td class="w-1/5">
						{#if !emailCheck.success}
							<p class="mb-1.5 text-xs text-error">
								{emailCheck.error.errors[0].message.replace('String', '')}
							</p>
						{:else}
							<p class="mb-1.5 text-xs text-transparent">For formatting sake</p>
						{/if}
						<input
							type="email"
							id="email"
							name="email"
							autocomplete="off"
							required
							class="input input-sm input-bordered mb-5 w-full"
							bind:value={user.email}
							class:input-error={!emailCheck.success}
						/>
					</td>
					<td>
						<input
							type="checkbox"
							bind:checked={user.admin}
							class="checkbox-info checkbox checkbox-sm"
						/>
					</td>
					<td>
						<input
							type="checkbox"
							bind:checked={user.operator}
							class="checkbox-info checkbox checkbox-sm"
						/>
					</td>
					<td class="my-1.5 flex flex-col items-center justify-center space-y-3">
						<button
							class="btn btn-success btn-sm w-14"
							class:btn-disabled={JSON.stringify(user) === JSON.stringify(lastUserValues) ||
								!nameCheck.success ||
								!surnameCheck.success ||
								!emailCheck.success}
							on:click={() => {
								updateUser(user);
							}}>Save</button
						>
						<button
							class="btn btn-sm"
							on:click={() => {
								// so that the input fields are reset. In runes mode, cannot use `user = lastUserValues`
								users[i] = lastUserValues;
								users[i].editing = false;
							}}>Cancel</button
						>
					</td>
				{/if}
			</tr>
		{/each}
	</tbody>
</table>
