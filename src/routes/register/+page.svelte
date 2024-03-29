<form action="/login" method="POST" class="flex flex-col space-y-2 sm:space-y-5">
	<div class="flex flex-col items-center space-y-2">
		<label for="name">Name</label>
		<input
			type="text"
			id="name"
			name="name"
			autocomplete="off"
			required
			class="input input-bordered w-full max-w-xs"
		/>
	</div>
	<div class="flex flex-col items-center space-y-2">
		<label for="surname">Surname</label>
		<input
			type="text"
			id="surname"
			name="surname"
			autocomplete="off"
			required
			class="input input-bordered w-full max-w-xs"
		/>
	</div>
	<div class="flex flex-col items-center space-y-2">
		<label for="email">Email</label>
		<input
			type="email"
			id="email"
			name="email"
			autocomplete="off"
			required
			class="input input-bordered w-full max-w-xs"
		/>
	</div>
	<div class="flex flex-col items-center space-y-2">
		<label for="password">Password</label>
		<input
			type="password"
			id="password"
			name="password"
			autocomplete="off"
			required
			class="input input-bordered w-full max-w-xs"
		/>
	</div>
	<div class="flex flex-col items-center space-y-2">
		<label for="password-confirm">Confirm Password</label>
		<input
			type="password"
			id="password-confirm"
			name="password-confirm"
			autocomplete="off"
			required
			class="input input-bordered w-full max-w-xs"
		/>
	</div>
	<div class="flex flex-col items-center">
		<button class="btn btn-neutral w-full max-w-xs" type="submit">Register</button>
	</div>
</form>
