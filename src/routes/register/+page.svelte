<script lang="ts">
	import { enhance } from '$app/forms';
	import { fieldSchema, emailSchema, passwordSchema } from '$lib/validationSchemas';
	import ZodIssues from '$lib/components/ZodIssues.svelte';
	import {goto} from '$app/navigation';

	let { form } = $props();

	let name = $state('');
	let surname = $state('');
	let email = $state('');
	let password = $state('');
	let passwordConfirmation = $state('');
	let isNameValid = $derived(fieldSchema.safeParse(name));
	let isSurnameValid = $derived(fieldSchema.safeParse(surname));
	let isEmailValid = $derived(emailSchema.safeParse(email));
	let isPasswordValid = $derived(passwordSchema.safeParse(password));
	let isPasswordConfirmationValid = $derived(passwordSchema.safeParse(passwordConfirmation));
	let isFormValid = $derived(
		isNameValid.success &&
			isSurnameValid.success &&
			isEmailValid.success &&
			isPasswordValid.success &&
			isPasswordConfirmationValid.success
	);

	// Reset form upon sending it
	$effect(() => {
		if (form) {
			password = '';
			passwordConfirmation = '';
		}
	});

	// Show registration successful
	let countdown = $state(3);
	let modalTimeout: ReturnType<typeof setTimeout>;
	$effect(() => {
		if (form?.message) {
			const dialog = document.getElementById('registration-successful') as HTMLDialogElement;
			dialog.showModal();
			let countdownInterval = setInterval(() => {
				countdown -= 1;
			}, 1000);
			modalTimeout = setTimeout(() => {
				dialog.close();
				clearInterval(countdownInterval);
				goto('/login');
			}, 4000); // 3 seconds + 1 second for the last countdown
		}
	});
</script>

<svelte:head>
	<title>MALGO - Register</title>
</svelte:head>

{#if form?.issues && form.issues.length > 0}
	<ZodIssues
		issues={form.issues}
		on:close={() => {
			form.issues = [];
		}}
	/>
{/if}

<dialog id="registration-successful" class="modal modal-bottom sm:modal-middle">
	<div class="modal-box">
		<h3 class="font-bold text-lg">User registered successfully!</h3>
		<p class="py-4">You will be redirected to login in <span class="countdown"><span style="--value:{countdown};"></span></span></p>
		<div class="modal-action">
			<form method="dialog">
				<button class="btn" on:click={()=>{clearTimeout(modalTimeout); goto('/login')}}>Go right now!</button>
			</form>
		</div>
	</div>
</dialog>

<form method="POST" class="flex flex-col space-y-4" use:enhance>
	<div class="flex flex-col items-center space-y-2">
		<label for="name">Name</label>
		{#if !isNameValid.success}
			<p class="text-xs text-error" style="margin: 0">
				{isNameValid.error.errors[0].message.replace('String', '')}
			</p>
		{:else}
			<p class="text-xs text-transparent" style="margin: 0">For formatting sake</p>
		{/if}
		<input
			type="text"
			id="name"
			name="name"
			autocomplete="off"
			required
			class="input input-bordered w-full max-w-xs"
			bind:value={name}
			class:input-error={!isNameValid.success}
		/>
	</div>
	<div class="flex flex-col items-center space-y-2">
		<label for="surname">Surname</label>
		{#if !isSurnameValid.success}
			<p class="text-xs text-error" style="margin: 0">
				{isSurnameValid.error.errors[0].message.replace('String', '')}
			</p>
		{:else}
			<p class="text-xs text-transparent" style="margin: 0">For formatting sake</p>
		{/if}
		<input
			type="text"
			id="surname"
			name="surname"
			autocomplete="off"
			required
			class="input input-bordered w-full max-w-xs"
			bind:value={surname}
			class:input-error={!isSurnameValid.success}
		/>
	</div>
	<div class="flex flex-col items-center space-y-2">
		<label for="email">Email</label>
		{#if !isEmailValid.success}
			<p class="text-xs text-error" style="margin: 0">
				{isEmailValid.error.errors[0].message.replace('String', '')}
			</p>
		{:else}
			<p class="text-xs text-transparent" style="margin: 0">For formatting sake</p>
		{/if}
		<input
			type="email"
			id="email"
			name="email"
			autocomplete="off"
			required
			class="input input-bordered w-full max-w-xs"
			bind:value={email}
			class:input-error={!isEmailValid.success}
		/>
	</div>
	<div class="flex flex-col items-center space-y-2">
		<label for="password">Password</label>
		{#if !isPasswordValid.success}
			<p class="text-xs text-error" style="margin: 0">
				{isPasswordValid.error.errors[0].message.replace('String', '')}
			</p>
		{:else}
			<p class="text-xs text-transparent" style="margin: 0">For formatting sake</p>
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
	<div class="flex flex-col items-center space-y-2">
		<label for="passwordConfirmation">Confirm Password</label>
		{#if !isPasswordConfirmationValid.success}
			<p class="text-xs text-error" style="margin: 0">
				{isPasswordConfirmationValid.error.errors[0].message.replace('String', '')}
			</p>
		{:else}
			<p class="text-xs text-transparent" style="margin: 0">For formatting sake</p>
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
	<div class="flex flex-col items-center">
		<button class="btn btn-neutral w-full max-w-xs" type="submit" class:btn-disabled={!isFormValid}
			>Register</button
		>
	</div>
</form>
