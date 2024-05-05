<script lang="ts">
	import { enhance } from '$app/forms';

	type InputProps = {
		isLoggedIn: boolean;
		isAdmin: boolean;
		isOperator: boolean;
	};

	const { isLoggedIn, isAdmin, isOperator }: InputProps = $props();
</script>

<header id="header" class="flex-none">
	<div class="navbar bg-base-300">
		<div class="flex-1">
			{#if isLoggedIn}
				<a class="kbd kbd-lg text-xl" href="/home">MALGO</a>
				<ul class="menu menu-horizontal px-1">
					{#if isOperator}
						<li><a href="/projects" data-sveltekit-preload-data>Projects</a></li>
					{/if}
				</ul>
			{:else}
				<a class="kbd kbd-lg text-xl" href="/">MALGO</a>
			{/if}
		</div>
		<div class="flex-none">
			<ul class="menu menu-horizontal px-1">
				{#if isLoggedIn}
					{#if isAdmin}
						<li><a href="/admin">Admin</a></li>
					{/if}
					<li>
						<div class="dropdown dropdown-end dropdown-bottom dropdown-hover">
							<span tabindex="-1" role="button">Hello!</span>
							<ul
								tabindex="-1"
								class="menu dropdown-content z-[1] w-52 rounded-box bg-base-200 p-2 shadow"
							>
								<li><a href="/settings" data-sveltekit-preload-data>Settings</a></li>
								<li>
									<form id="logout-form" method="post" action="/logout" use:enhance>
										<button class="w-40 text-left">Sign out</button>
									</form>
								</li>
							</ul>
						</div>
					</li>
				{:else}
					<li><a href="/login">Login</a></li>
				{/if}
			</ul>
		</div>
	</div>
</header>
