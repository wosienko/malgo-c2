<script lang="ts">
	import { enhance } from '$app/forms';
	import { emailSchema, passwordSchema } from '$lib/validationSchemas';
	import ZodIssues from '$lib/components/toasts/ZodIssues.svelte';
	import ValidatedInputWithVerticalLabel from '$lib/components/inputs/ValidatedInputWithVerticalLabel.svelte';

	let { form } = $props();

	let email = $state('');
	let password = $state('');
	let isEmailValid = $derived(emailSchema.safeParse(email));
	let isPasswordValid = $derived(passwordSchema.safeParse(password));
	let isFormValid = $derived(isEmailValid.success && isPasswordValid.success);

	// Reset form upon sending it
	$effect(() => {
		if (form) {
			email = '';
			password = '';
		}
	});
</script>

<svelte:head>
	<title>MALGO - Login</title>
</svelte:head>

{#if form?.issues && form.issues.length > 0}
	<ZodIssues
		issues={form.issues}
		on:close={() => {
			form.issues = [];
		}}
	/>
{/if}

<form action="/login" method="POST" class="flex flex-col space-y-4" use:enhance>
	<ValidatedInputWithVerticalLabel
		label="Email"
		id="email"
		name="email"
		type="email"
		bind:value={email}
		validation={isEmailValid}
	/>
	<ValidatedInputWithVerticalLabel
		label="Password"
		id="password"
		name="password"
		type="password"
		bind:value={password}
		validation={isPasswordValid}
	/>
	<div class="flex flex-col items-center">
		<button class="btn btn-neutral w-full max-w-xs" type="submit" class:btn-disabled={!isFormValid}
			>Login</button
		>
	</div>
</form>
