<script lang="ts">
	import { enhance } from '$app/forms';
	import { emailSchema, passwordSchema } from '$lib/validationSchemas';
	import ZodIssues from '$lib/components/ZodIssues.svelte';

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
	<title>Login</title>
</svelte:head>

{#if form?.issues && form.issues.length > 0}
	<ZodIssues
		issues={form.issues}
		on:close={() => {
			form.issues = [];
		}}
	/>
{/if}

<form action="/login" method="POST" class="flex flex-col space-y-5" use:enhance>
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
	<div class="flex flex-col items-center">
		<button class="btn btn-neutral w-full max-w-xs" type="submit" class:btn-disabled={!isFormValid}
			>Login</button
		>
	</div>
</form>
