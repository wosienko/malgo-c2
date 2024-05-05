<script lang="ts">
	import type { ApiError } from '$lib';
	import type { ZodIssue } from 'zod';
	import { emailSchema, fieldSchema, passwordSchema } from '$lib/validationSchemas';
	import ModalRunCancel from '$lib/components/modals/ModalRunCancel.svelte';
	import ValidatedInputWithLabel from '$lib/components/inputs/ValidatedInputWithHorizontalLabel.svelte';

	type InputProps = {
		successMessage: string | undefined;
		zodIssues: ZodIssue[];
		apiError: ApiError | undefined;
		reloadCurrentPage: () => Promise<void>;
		showModal: () => void;
	};

	let {
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		successMessage = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		zodIssues = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		apiError = $bindable(),
		reloadCurrentPage,
		showModal = $bindable()
	}: InputProps = $props();

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

	let newUserValid = $derived(
		newUserVerification.name.success &&
			newUserVerification.surname.success &&
			newUserVerification.email.success &&
			newUserVerification.password.success &&
			newUserVerification.passwordConfirmation.success
	);

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

<ModalRunCancel
	id="register-new-user"
	title="Register new user"
	btnClass="btn-success"
	btnText="Register new user"
	btnDisabledCondition={!newUserValid}
	onclickCallback={registerNewUser}
	bind:showModal
	onHideModal={() => {
		newUser = {
			name: '',
			surname: '',
			email: '',
			password: '',
			passwordConfirmation: '',
			admin: false,
			operator: false
		};
	}}
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
