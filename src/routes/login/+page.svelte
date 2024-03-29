<script lang="ts">
	import { enhance } from '$app/forms';
	import { emailSchema, passwordSchema } from '$lib/validationSchemas';

	let email = $state('');
	let password = $state('');
	let isEmailValid = $derived(emailSchema.safeParse(email));
	let isPasswordValid = $derived(passwordSchema.safeParse(password));
	let isFormValid = $derived(isEmailValid && isPasswordValid);
</script>

<svelte:head>
	<title>Login</title>
</svelte:head>

<form action="/login" method="POST" class="flex flex-col space-y-5" use:enhance>
	<div class="flex flex-col items-center space-y-2">
		<label for="email">Email</label>
		{#if !isEmailValid.success}
			<p class="text-xs text-error">{isEmailValid.error.errors[0].message.replace('String', '')}</p>
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
			<p class="text-xs text-error">
				{isPasswordValid.error.errors[0].message.replace('String', '')}
			</p>
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
