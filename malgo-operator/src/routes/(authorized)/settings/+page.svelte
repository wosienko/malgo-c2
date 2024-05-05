<script lang="ts">
	import ValidatedInputWithHorizontalLabel from '$lib/components/inputs/ValidatedInputWithHorizontalLabel.svelte';
	import { fieldSchema, emailSchema, passwordSchema } from '$lib/validationSchemas';
	import type { ZodIssue } from 'zod';
	import type { ApiError, User } from '$lib';
	import ZodIssues from '$lib/components/toasts/ZodIssues.svelte';
	import ApiIssues from '$lib/components/toasts/ApiIssues.svelte';
	import ChangeSuccessful from '$lib/components/toasts/ChangeSuccessful.svelte';
	import ModalRunCancel from '$lib/components/modals/ModalRunCancel.svelte';
	import ValidatedInputWithLabel from '$lib/components/inputs/ValidatedInputWithHorizontalLabel.svelte';

	let { data } = $props();
	let loading = $state(false);

	let user: User = $state({
		id: data.user.id,
		name: data.user.name,
		surname: data.user.surname,
		email: data.user.email
	});
	let editingValid = $derived(() => {
		return (
			fieldSchema.safeParse(user.name).success &&
			fieldSchema.safeParse(user.surname).success &&
			emailSchema.safeParse(user.email).success
		);
	});

	let zodIssues: ZodIssue[] = $state([]);
	let apiError: ApiError | undefined = $state();
	let successMessage: string | undefined = $state();

	const updateUser = async () => {
		loading = true;
		const res = await fetch(`/api/self`, {
			method: 'PATCH',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				name: user.name,
				surname: user.surname,
				email: user.email
			})
		});
		loading = false;

		if (res.ok) {
			successMessage = 'User updated successfully!';
			data.user = { ...user };
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

	let passwordChange = $state({
		currentPassword: '',
		newPassword: '',
		newPasswordConfirmation: ''
	});

	let passwordChangeVerification = $derived({
		currentPassword: passwordSchema.safeParse(passwordChange.currentPassword),
		newPassword: passwordSchema.safeParse(passwordChange.newPassword),
		newPasswordConfirmation: passwordSchema
			.refine((data) => data === passwordChange.newPassword, {
				message: "New passwords don't match",
				path: ['New Password Confirmation']
			})
			.safeParse(passwordChange.newPasswordConfirmation)
	});
	let passwordChangeValid = $derived(
		passwordChangeVerification.currentPassword.success &&
			passwordChangeVerification.newPassword.success &&
			passwordChangeVerification.newPasswordConfirmation.success
	);

	let showPasswordChangeModal = $state(() => {});

	const changePassword = async (): Promise<boolean> => {
		const res = await fetch(`/api/self/password`, {
			method: 'PATCH',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(passwordChange)
		});

		if (res.ok) {
			successMessage = 'Password changed successfully!';
			passwordChange = {
				currentPassword: '',
				newPassword: '',
				newPasswordConfirmation: ''
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
</script>

<svelte:head>
	<title>MALGO - Settings</title>
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
	id="password-change"
	title="Password change"
	btnClass="btn-warning"
	btnText="Change password"
	btnDisabledCondition={!passwordChangeValid}
	onclickCallback={changePassword}
	bind:showModal={showPasswordChangeModal}
>
	<div class="space-y-8">
		<ValidatedInputWithLabel
			label="Current password"
			type="password"
			id="current-password"
			name="currentPassword"
			bind:value={passwordChange.currentPassword}
			validation={passwordChangeVerification.currentPassword}
			classes="w-full max-w-xs"
		/>
		<ValidatedInputWithLabel
			label="New password"
			type="password"
			id="new-password"
			name="newPassword"
			bind:value={passwordChange.newPassword}
			validation={passwordChangeVerification.newPassword}
			classes="w-full max-w-xs"
		/>
		<ValidatedInputWithLabel
			label="Confirm new password"
			type="password"
			id="new-password-confirmation"
			name="newPasswordConfirmation"
			bind:value={passwordChange.newPasswordConfirmation}
			validation={passwordChangeVerification.newPasswordConfirmation}
			classes="w-full max-w-xs"
		/>
	</div>
</ModalRunCancel>

<div class="h-full w-full max-w-xs text-center md:h-auto">
	<h1 class="mb-6 text-4xl font-extrabold leading-none tracking-tight">Settings</h1>
	<form>
		<div class="space-y-6">
			<ValidatedInputWithHorizontalLabel
				id="name"
				label="Name"
				name="name"
				type="text"
				bind:value={user.name}
				validation={fieldSchema.safeParse(user.name)}
			/>
			<ValidatedInputWithHorizontalLabel
				id="surname"
				label="Surname"
				name="surname"
				type="text"
				bind:value={user.surname}
				validation={fieldSchema.safeParse(user.surname)}
			/>
			<ValidatedInputWithHorizontalLabel
				id="email"
				label="Email"
				name="email"
				type="email"
				bind:value={user.email}
				validation={emailSchema.safeParse(user.email)}
			/>
		</div>
		{#if loading}
			<div class="my-4 h-8">
				<span class="loading loading-spinner loading-lg text-info"></span>
			</div>
		{:else}
			<button
				class="btn btn-success mt-4 w-full"
				class:btn-disabled={JSON.stringify(data.user) === JSON.stringify(user) || !editingValid}
				on:click={updateUser}>Save</button
			>
		{/if}
	</form>
	<button class="btn btn-warning mt-4 w-full" on:click={showPasswordChangeModal}
		>Change password</button
	>
</div>
