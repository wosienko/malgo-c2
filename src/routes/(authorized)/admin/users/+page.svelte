<script lang="ts">
	import type { UsersWithRoles, UserWithRoles } from '$lib/db/schema/users';
	import type { ApiError } from '$lib';
	import { emailSchema, fieldSchema, passwordSchema } from '$lib/validationSchemas';
	import type { ZodIssue } from 'zod';
	import { pushState } from '$app/navigation';
	import ZodIssues from '$lib/components/ZodIssues.svelte';
	import ApiIssues from '$lib/components/ApiIssues.svelte';
	import ChangeSuccessful from '$lib/components/ChangeSuccessful.svelte';

	let { data } = $props();

	type User = {
		id: string;
		name: string;
		surname: string;
		email: string;
		admin: boolean;
		operator: boolean;
		editing: boolean;
	};

	let users: User[] = $state(
		data.users.users.map((user) => ({
			...user,
			editing: false
		}))
	);

	let lastUserValues: User = $state({
		id: '',
		name: '',
		surname: '',
		email: '',
		admin: false,
		operator: false,
		editing: false
	});

	let page = $state(data.page);
	let pageSize = $state(data.pageSize);

	let zodIssues: ZodIssue[] = $state([]);
	let apiError: ApiError | undefined = $state();
	let successMessage: string | undefined = $state();

	const loadPage = (nextPage: number) => {
		return async () => {
			const res = await fetch(`/api/user?page=${nextPage}&pageSize=${pageSize}`);
			pushState(`/admin/users?page=${nextPage}&pageSize=${pageSize}`, {});
			const usersPage: UsersWithRoles = await res.json();
			users = usersPage.users.map((user) => ({
				...user,
				editing: false
			}));
			data.users.count = usersPage.count;
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

	const updateUser = async (user: User) => {
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
			successMessage = 'User updated successfully!';
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

	let userToAlter: User = $state({
		id: '',
		name: '',
		surname: '',
		email: '',
		admin: false,
		operator: false,
		editing: false
	});

	const deleteUser = async (user: User): Promise<boolean> => {
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

			successMessage = 'User deleted successfully!';
			return true;
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
			return false;
		}
	};

	const prepareForDeletion = (user: typeof userToAlter) => {
		return () => {
			userToAlter = user;
			const dialog = document.getElementById('confirm-deletion') as HTMLDialogElement;
			dialog.showModal();
		};
	};

	const cleanupDeletion = () => {
		const dialog = document.getElementById('confirm-deletion') as HTMLDialogElement;
		dialog.close();
	};

	let password = $state('');
	let passwordConfirmation = $state('');

	const isPasswordValid = $derived(passwordSchema.safeParse(password));
	const isPasswordConfirmationValid = $derived(
		passwordSchema
			.refine((data) => data === password, { message: "Passwords don't match" })
			.safeParse(passwordConfirmation)
	);

	const changePassword = async (user: (typeof users)[0]): Promise<boolean> => {
		const res = await fetch(`/api/user/${user.id}/password`, {
			method: 'PATCH',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				password: password,
				passwordConfirmation: passwordConfirmation
			})
		});

		if (res.ok) {
			successMessage = 'Password changed successfully!';
			return true;
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
			return false;
		}
	};

	const prepareForPasswordChange = (user: typeof userToAlter) => {
		return () => {
			userToAlter = user;
			const dialog = document.getElementById('password-change') as HTMLDialogElement;
			dialog.showModal();
		};
	};

	const cleanupPasswordChange = () => {
		password = '';
		passwordConfirmation = '';
		const dialog = document.getElementById('password-change') as HTMLDialogElement;
		dialog.close();
	};

	let newUser = $state({
		name: '',
		surname: '',
		email: '',
		password: '',
		passwordConfirmation: '',
		admin: false,
		operator: false
	});

	let newUserVerification = $derived({
		name: fieldSchema.safeParse(newUser.name),
		surname: fieldSchema.safeParse(newUser.surname),
		email: emailSchema.safeParse(newUser.email),
		password: passwordSchema.safeParse(newUser.password),
		passwordConfirmation: passwordSchema
			.refine((data) => data === newUser.password, { message: "Passwords don't match" })
			.safeParse(newUser.passwordConfirmation)
	});

	const registerNewUser = async (): Promise<boolean> => {
		const res = await fetch('/api/user', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(newUser)
		});

		if (res.ok) {
			successMessage = 'User registered successfully!';
			return true;
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
			return false;
		}
	};

	const prepareForRegistration = () => {
		const dialog = document.getElementById('register-new-user') as HTMLDialogElement;
		dialog.showModal();
	};

	const cleanupRegistration = async () => {
		newUser = {
			name: '',
			surname: '',
			email: '',
			password: '',
			passwordConfirmation: '',
			admin: false,
			operator: false
		};
		await reloadCurrentPage();
		const dialog = document.getElementById('register-new-user') as HTMLDialogElement;
		dialog.close();
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

{#if successMessage}
	<ChangeSuccessful
		message={successMessage}
		on:close={() => {
			successMessage = undefined;
		}}
	/>
{/if}

<dialog id="confirm-deletion" class="modal modal-bottom sm:modal-middle">
	<div class="modal-box">
		<h3 class="text-lg font-bold">Are you sure you want to delete?</h3>
		<p class="py-4">
			Account to be deleted: <span class="font-bold">{userToAlter.name} {userToAlter.surname}</span>
		</p>
		<div class="modal-action">
			<form method="dialog" class="space-x-3">
				<button
					class="btn btn-error"
					on:click|preventDefault={async () => {
						let result = await deleteUser(userToAlter);
						if (result) cleanupDeletion();
					}}>Delete</button
				>
				<button class="btn">Cancel</button>
			</form>
		</div>
	</div>
</dialog>

<dialog id="password-change" class="modal modal-middle">
	<div class="modal-box">
		<h3 class="text-lg font-bold">Password change</h3>
		<p class="py-4">
			Password change for <span class="font-bold">{userToAlter.name} {userToAlter.surname}</span>
		</p>
		<div class="space-y-8">
			<div class="flex items-center justify-between">
				<label for="password">Password</label>
				<div class="-mt-4 flex flex-col items-center">
					{#if !isPasswordValid.success}
						<p class="mb-1.5 text-xs text-error">
							{isPasswordValid.error.errors[0].message.replace('String', '')}
						</p>
					{:else}
						<p class="mb-1.5 text-xs text-transparent">For formatting sake</p>
					{/if}
					<input
						type="password"
						id="password"
						name="password"
						autocomplete="off"
						required
						class="input input-bordered w-full max-w-xs"
						bind:value={password}
						class:input-error={!isPasswordValid.success}
					/>
				</div>
			</div>
			<div class="flex items-center justify-between">
				<label for="passwordConfirmation">Confirm Password</label>
				<div class="-mt-4 flex flex-col items-center">
					{#if !isPasswordConfirmationValid.success}
						<p class="mb-1.5 text-xs text-error">
							{isPasswordConfirmationValid.error.errors[0].message.replace('String', '')}
						</p>
					{:else}
						<p class="mb-1.5 text-xs text-transparent">For formatting sake</p>
					{/if}
					<input
						type="password"
						id="passwordConfirmation"
						name="passwordConfirmation"
						autocomplete="off"
						required
						class="input input-bordered w-full max-w-xs"
						bind:value={passwordConfirmation}
						class:input-error={!isPasswordConfirmationValid.success}
					/>
				</div>
			</div>
		</div>
		<div class="modal-action">
			<form method="dialog" class="space-x-3">
				<button
					class="btn btn-warning"
					class:btn-disabled={!isPasswordValid.success || !isPasswordConfirmationValid.success}
					on:click={async () => {
						let result = await changePassword(userToAlter);
						if (result) cleanupPasswordChange();
					}}>Change password</button
				>
				<button class="btn">Cancel</button>
			</form>
		</div>
	</div>
</dialog>

<dialog id="register-new-user" class="modal modal-middle">
	<div class="modal-box">
		<h3 class="mb-3 text-lg font-bold">Register new user</h3>
		<div class="space-y-8">
			<div class="flex items-center justify-between">
				<label for="name">Name</label>
				<div class="-mt-4 flex flex-col items-center">
					{#if !newUserVerification.name.success}
						<p class="mb-1.5 text-xs text-error">
							{newUserVerification.name.error.errors[0].message.replace('String', '')}
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
						class="input input-bordered w-full max-w-xs"
						bind:value={newUser.name}
						class:input-error={!newUserVerification.name.success}
					/>
				</div>
			</div>
			<div class="flex items-center justify-between">
				<label for="surname">Surname</label>
				<div class="-mt-4 flex flex-col items-center">
					{#if !newUserVerification.surname.success}
						<p class="mb-1.5 text-xs text-error">
							{newUserVerification.surname.error.errors[0].message.replace('String', '')}
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
						class="input input-bordered w-full max-w-xs"
						bind:value={newUser.surname}
						class:input-error={!newUserVerification.surname.success}
					/>
				</div>
			</div>
			<div class="flex items-center justify-between">
				<label for="email">Email</label>
				<div class="-mt-4 flex flex-col items-center">
					{#if !newUserVerification.email.success}
						<p class="mb-1.5 text-xs text-error">
							{newUserVerification.email.error.errors[0].message.replace('String', '')}
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
						class="input input-bordered w-full max-w-xs"
						bind:value={newUser.email}
						class:input-error={!newUserVerification.email.success}
					/>
				</div>
			</div>
			<div class="flex items-center justify-between">
				<label for="password">Password</label>
				<div class="-mt-4 flex flex-col items-center">
					{#if !newUserVerification.password.success}
						<p class="mb-1.5 text-xs text-error">
							{newUserVerification.password.error.errors[0].message.replace('String', '')}
						</p>
					{:else}
						<p class="mb-1.5 text-xs text-transparent">For formatting sake</p>
					{/if}
					<input
						type="password"
						id="password"
						name="password"
						autocomplete="off"
						required
						class="input input-bordered w-full max-w-xs"
						bind:value={newUser.password}
						class:input-error={!newUserVerification.password.success}
					/>
				</div>
			</div>
			<div class="flex items-center justify-between">
				<label for="passwordConfirmation">Confirm Password</label>
				<div class="-mt-4 flex flex-col items-center">
					{#if !newUserVerification.passwordConfirmation.success}
						<p class="mb-1.5 text-xs text-error">
							{newUserVerification.passwordConfirmation.error.errors[0].message.replace(
								'String',
								''
							)}
						</p>
					{:else}
						<p class="mb-1.5 text-xs text-transparent">For formatting sake</p>
					{/if}
					<input
						type="password"
						id="passwordConfirmation"
						name="passwordConfirmation"
						autocomplete="off"
						required
						class="input input-bordered w-full max-w-xs"
						bind:value={newUser.passwordConfirmation}
						class:input-error={!newUserVerification.passwordConfirmation.success}
					/>
				</div>
			</div>
			<div class="flex w-full justify-around">
				<div class="form-control max-w-40 flex-1">
					<label class="label cursor-pointer">
						<span class="label-text">Admin</span>
						<input type="checkbox" bind:checked={newUser.admin} class="checkbox-warning checkbox" />
					</label>
				</div>
				<div class="form-control max-w-40 flex-1">
					<label class="label cursor-pointer">
						<span class="label-text">Operator</span>
						<input type="checkbox" bind:checked={newUser.operator} class="checkbox-info checkbox" />
					</label>
				</div>
			</div>
		</div>
		<div class="modal-action">
			<form method="dialog" class="space-x-3">
				<button
					class="btn btn-success"
					class:btn-disabled={!newUserVerification.name.success ||
						!newUserVerification.surname.success ||
						!newUserVerification.email.success ||
						!newUserVerification.password.success ||
						!newUserVerification.passwordConfirmation.success}
					on:click|preventDefault={async () => {
						const result = await registerNewUser();
						if (result) await cleanupRegistration();
					}}>Register new user</button
				>
				<button class="btn">Cancel</button>
			</form>
		</div>
	</div>
</dialog>

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
							class="checkbox-warning checkbox checkbox-sm"
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
								<li>
									<button
										class="btn btn-warning btn-sm"
										on:click={() => {
											prepareForPasswordChange(user)();
										}}>Change password</button
									>
								</li>
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
							class="checkbox-warning checkbox checkbox-sm"
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
