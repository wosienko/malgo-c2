<script lang="ts">
	import type { UsersWithRoles, UserWithRoles } from '$lib/db/schema/users';
	import type { ApiError } from '$lib';
	import { emailSchema, fieldSchema, passwordSchema } from '$lib/validationSchemas';
	import type { ZodIssue } from 'zod';
	import { pushState } from '$app/navigation';
	import ZodIssues from '$lib/components/toasts/ZodIssues.svelte';
	import ApiIssues from '$lib/components/toasts/ApiIssues.svelte';
	import ChangeSuccessful from '$lib/components/toasts/ChangeSuccessful.svelte';
	import ModalRunCancel from '$lib/components/modals/ModalRunCancel.svelte';
	import ValidatedInput from '$lib/components/inputs/ValidatedInput.svelte';
	import ValidatedInputWithLabel from '$lib/components/inputs/ValidatedInputWithLabel.svelte';

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
	const createEmptyUser = (): User => ({
		id: '',
		name: '',
		surname: '',
		email: '',
		admin: false,
		operator: false,
		editing: false
	});

	let users: User[] = $state(
		data.users.users.map((user) => ({
			...user,
			editing: false
		}))
	);

	let lastUserValues: User = $state(createEmptyUser());

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

	let userToAlter: User = $state(createEmptyUser());

	let showDeletionModal = $state(() => {});
	const prepareForDeletion = (user: typeof userToAlter) => {
		return () => {
			userToAlter = user;
			showDeletionModal();
		};
	};
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
				editing: false
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

	let passwordChange = $state({
		password: '',
		passwordConfirmation: ''
	});

	let passwordChangeVerification = $derived({
		password: passwordSchema.safeParse(passwordChange.password),
		passwordConfirmation: passwordSchema
			.refine((data) => data === passwordChange.password, { message: "Passwords don't match" })
			.safeParse(passwordChange.passwordConfirmation)
	});

	let showPasswordChangeModal = $state(() => {});
	const prepareForPasswordChange = (user: typeof userToAlter) => {
		return () => {
			userToAlter = user;
			showPasswordChangeModal();
		};
	};
	const changePassword = async (user: User): Promise<boolean> => {
		const res = await fetch(`/api/user/${user.id}/password`, {
			method: 'PATCH',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(passwordChange)
		});

		if (res.ok) {
			successMessage = 'Password changed successfully!';
			passwordChange = {
				password: '',
				passwordConfirmation: ''
			};
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

	let showRegistrationModal = $state(() => {});
	const prepareForRegistration = () => {
		showRegistrationModal();
	};
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

<ModalRunCancel
	id="confirm-deletion"
	title="Are you sure you want to delete?"
	message="Account to be deleted: "
	messageEmphasis={`${userToAlter.name} ${userToAlter.surname}`}
	btnClass="btn-error"
	btnText="Delete"
	onclickCallback={() => deleteUser(userToAlter)}
	bind:showModal={showDeletionModal}
/>

<ModalRunCancel
	id="password-change"
	title="Password change"
	message="Password change for: "
	messageEmphasis={`${userToAlter.name} ${userToAlter.surname}`}
	btnClass="btn-warning"
	btnText="Change password"
	btnDisabledCondition={!passwordChangeVerification.password.success ||
		!passwordChangeVerification.passwordConfirmation.success}
	onclickCallback={() => changePassword(userToAlter)}
	bind:showModal={showPasswordChangeModal}
>
	<div class="space-y-8">
		<ValidatedInputWithLabel
			label="Password"
			type="password"
			id="password"
			name="password"
			bind:value={passwordChange.password}
			validation={passwordChangeVerification.password}
			classes="w-full max-w-xs"
		/>
		<ValidatedInputWithLabel
			label="Confirm Password"
			type="password"
			id="passwordConfirmation"
			name="passwordConfirmation"
			bind:value={passwordChange.passwordConfirmation}
			validation={passwordChangeVerification.passwordConfirmation}
			classes="w-full max-w-xs"
		/>
	</div>
</ModalRunCancel>

<ModalRunCancel
	id="register-new-user"
	title="Register new user"
	btnClass="btn-success"
	btnText="Register new user"
	btnDisabledCondition={!newUserVerification.name.success ||
		!newUserVerification.surname.success ||
		!newUserVerification.email.success ||
		!newUserVerification.password.success ||
		!newUserVerification.passwordConfirmation.success}
	onclickCallback={registerNewUser}
	bind:showModal={showRegistrationModal}
>
	<div class="space-y-8">
		<ValidatedInputWithLabel
			label="Name"
			type="text"
			id="name"
			name="name"
			bind:value={newUser.name}
			validation={newUserVerification.name}
			classes="w-full max-w-xs"
		/>
		<ValidatedInputWithLabel
			label="Surname"
			type="text"
			id="surname"
			name="surname"
			bind:value={newUser.surname}
			validation={newUserVerification.surname}
			classes="w-full max-w-xs"
		/>
		<ValidatedInputWithLabel
			label="Email"
			type="email"
			id="email"
			name="email"
			bind:value={newUser.email}
			validation={newUserVerification.email}
			classes="w-full max-w-xs"
		/>
		<ValidatedInputWithLabel
			label="Password"
			type="password"
			id="password"
			name="password"
			bind:value={newUser.password}
			validation={newUserVerification.password}
			classes="w-full max-w-xs"
		/>
		<ValidatedInputWithLabel
			label="Confirm Password"
			type="password"
			id="passwordConfirmation"
			name="passwordConfirmation"
			bind:value={newUser.passwordConfirmation}
			validation={newUserVerification.passwordConfirmation}
			classes="w-full max-w-xs"
		/>
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
</ModalRunCancel>

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
			{#if !user.editing}
				<tr class="max-w-dvw hover">
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
						<!-- Dropdown accounting for hiding under the dropdown for last elements -->
						<div
							class="dropdown dropdown-end"
							class:dropdown-bottom={i <= 2 ||
								i + 1 < (data.users.count > pageSize ? pageSize : data.users.count) - 2}
							class:dropdown-top={i > 2 &&
								i + 1 >= (data.users.count > pageSize ? pageSize : data.users.count) - 2}
						>
							<div tabindex="-1" role="button" class="btn btn-neutral btn-sm mb-1">Options</div>
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
				</tr>
			{:else}
				{@const nameCheck = fieldSchema.safeParse(user.name)}
				{@const surnameCheck = fieldSchema.safeParse(user.surname)}
				{@const emailCheck = emailSchema.safeParse(user.email)}
				<tr class="max-w-dvw hover">
					<td>
						<ValidatedInput
							type="text"
							id="name"
							name="name"
							bind:value={user.name}
							validation={nameCheck}
							classes="input-sm mb-5 w-full"
						/>
					</td>
					<td>
						<ValidatedInput
							type="text"
							id="surname"
							name="surname"
							bind:value={user.surname}
							validation={surnameCheck}
							classes="input-sm mb-5 w-full"
						/>
					</td>
					<td class="w-1/5">
						<ValidatedInput
							type="email"
							id="email"
							name="email"
							bind:value={user.email}
							validation={emailCheck}
							classes="input-sm mb-5 w-full"
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
				</tr>
			{/if}
		{/each}
	</tbody>
</table>
