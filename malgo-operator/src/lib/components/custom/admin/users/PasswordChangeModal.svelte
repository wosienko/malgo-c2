<script lang="ts">
	import type { ApiError, UserWithRoles } from '$lib';
	import type { ZodIssue } from 'zod';
	import { passwordSchema } from '$lib/validationSchemas';
	import ModalRunCancel from '$lib/components/modals/ModalRunCancel.svelte';
	import ValidatedInputWithLabel from '$lib/components/inputs/ValidatedInputWithHorizontalLabel.svelte';

	type InputProps = {
		userToAlter: UserWithRoles;
		showModal: () => void;
		reloadCurrentPage: () => Promise<void>;
		successMessage: string | undefined;
		zodIssues: ZodIssue[];
		apiError: ApiError | undefined;
	};

	let {
		userToAlter,
		showModal = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		successMessage = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		zodIssues = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		apiError = $bindable()
	}: InputProps = $props();

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
	let passwordChangeValid = $derived(
		passwordChangeVerification.password.success &&
			passwordChangeVerification.passwordConfirmation.success
	);

	const changePassword = async (user: UserWithRoles): Promise<boolean> => {
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
</script>

<ModalRunCancel
	id="password-change"
	title="Password change"
	message="Password change for: "
	messageEmphasis={`${userToAlter.name} ${userToAlter.surname}`}
	btnClass="btn-warning"
	btnText="Change password"
	btnDisabledCondition={!passwordChangeValid}
	onclickCallback={() => changePassword(userToAlter)}
	bind:showModal
	onHideModal={() => {
		passwordChange = {
			password: '',
			passwordConfirmation: ''
		};
	}}
>
	<div class="space-y-8">
		<ValidatedInputWithLabel
			label="Password"
			type="password"
			id="new-password"
			name="password"
			bind:value={passwordChange.password}
			validation={passwordChangeVerification.password}
			classes="w-full max-w-xs"
		/>
		<ValidatedInputWithLabel
			label="Confirm Password"
			type="password"
			id="new-passwordConfirmation"
			name="passwordConfirmation"
			bind:value={passwordChange.passwordConfirmation}
			validation={passwordChangeVerification.passwordConfirmation}
			classes="w-full max-w-xs"
		/>
	</div>
</ModalRunCancel>
