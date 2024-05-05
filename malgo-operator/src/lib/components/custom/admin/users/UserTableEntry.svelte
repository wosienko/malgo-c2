<script lang="ts">
	import { emailSchema, fieldSchema } from '$lib/validationSchemas';
	import type { ApiError, UserWithRoles } from '$lib';
	import type { ZodIssue } from 'zod';
	import ValidatedInput from '$lib/components/inputs/ValidatedInput.svelte';

	type InputProps = {
		user: UserWithRoles;
		elementIndex: number;
		count: number;
		pageSize: number;
		prepareForPasswordChange: (user: UserWithRoles) => () => void;
		prepareForDeletion: (user: UserWithRoles) => () => void;
		successMessage: string | undefined;
		apiError: ApiError | undefined;
		zodIssues: ZodIssue[];
	};

	let {
		user = $bindable(),
		elementIndex,
		count,
		pageSize,
		prepareForPasswordChange,
		prepareForDeletion,
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		successMessage = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		apiError = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		zodIssues = $bindable()
	}: InputProps = $props();

	let editing = $state(false);
	let loading = $state(false);

	let lastUserValues = $state(JSON.parse(JSON.stringify(user)));

	const updateUser = async (user: UserWithRoles) => {
		loading = true;
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
		loading = false;

		if (res.ok) {
			editing = false;
			successMessage = 'User updated successfully!';
			lastUserValues = JSON.parse(JSON.stringify(user));
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
</script>

{#if !editing}
	<tr class="max-w-dvw hover">
		<td class="no-scrollbar hidden overflow-x-auto md:table-cell">{user.name}</td>
		<td class="no-scrollbar hidden overflow-x-auto md:table-cell">{user.surname}</td>
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
		<td>
			<!-- Dropdown accounting for hiding under the dropdown for last elements -->
			<div
				class="dropdown dropdown-end"
				class:dropdown-bottom={elementIndex <= 2 ||
					elementIndex + 1 < (count > pageSize ? pageSize : count) - 2}
				class:dropdown-top={elementIndex > 2 &&
					elementIndex + 1 >= (count > pageSize ? pageSize : count) - 2}
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
								editing = true;
								lastUserValues = JSON.parse(JSON.stringify(user));
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
	{@const editingValid = nameCheck.success && surnameCheck.success && emailCheck.success}
	<tr class="max-w-dvw hover">
		<td class="hidden md:table-cell">
			<div class="md:-mt-2">
				<ValidatedInput
					type="text"
					id="name"
					name="name"
					bind:value={user.name}
					validation={nameCheck}
					classes="input-sm mb-5 w-full"
				/>
			</div>
		</td>
		<td class="hidden md:table-cell">
			<div class="md:-mt-2">
				<ValidatedInput
					type="text"
					id="surname"
					name="surname"
					bind:value={user.surname}
					validation={surnameCheck}
					classes="input-sm mb-5 w-full"
				/>
			</div>
		</td>
		<td>
			<div class="md:-mt-2">
				<ValidatedInput
					type="email"
					id="email"
					name="email"
					bind:value={user.email}
					validation={emailCheck}
					classes="input-sm mb-5 w-full"
				/>
			</div>
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
			{#if loading}
				<div class="my-4 h-11">
					<span class="loading loading-spinner loading-lg text-info"></span>
				</div>
			{:else}
				<button
					class="btn btn-success btn-sm w-14"
					class:btn-disabled={JSON.stringify(user) === JSON.stringify(lastUserValues) ||
						!editingValid}
					on:click={() => {
						updateUser(user);
					}}>Save</button
				>
				<button
					class="btn btn-sm"
					on:click={() => {
						user = lastUserValues;
						editing = false;
					}}>Cancel</button
				>
			{/if}
		</td>
	</tr>
{/if}
