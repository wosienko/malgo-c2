<script lang="ts">
	import '../main.css';
	import { onMount } from 'svelte';
	import { navigating } from '$app/stores';
	import { enhance } from '$app/forms';
	import { version } from '$lib';

	let { data } = $props();

	// inline dark mode store, no need for separate file
	let lightMode = $state(false);

	// main tag overflow check
	let isOverflowing = $state(false);
	let contentElement: Element | null = null;
	const checkOverflow = () => {
		isOverflowing = contentElement!.scrollHeight > contentElement!.clientHeight;
	};

	onMount(() => {
		// load light mode from local storage
		lightMode = localStorage.getItem('lightMode') === 'true';
		$effect(() => {
			localStorage.setItem('lightMode', lightMode.toString());
		});

		// check if content is overflowing
		contentElement = document.querySelector('main');
		checkOverflow();
		window.addEventListener('resize', checkOverflow);
		let unsubscribe = navigating.subscribe(checkOverflow);

		// cleanup. We can't use onDestroy because onDestroy runs on server side
		return () => {
			window.removeEventListener('resize', checkOverflow);
			unsubscribe();
		};
	});
</script>

<div class="flex h-dvh flex-col">
	<header id="header" class="flex-none">
		<div class="navbar bg-base-300">
			<div class="flex-1">
				{#if data.loggedIn}
					<a class="kbd kbd-lg text-xl" href="/home">MALGO</a>
					<ul class="menu menu-horizontal px-1">
						{#if data.isOperator}
							<li><a href="/projects">Projects</a></li>
						{/if}
					</ul>
				{:else}
					<a class="kbd kbd-lg text-xl" href="/">MALGO</a>
				{/if}
			</div>
			<div class="flex-none">
				<ul class="menu menu-horizontal px-1">
					{#if data.loggedIn}
						{#if data.isAdmin}
							<li><a href="/admin">Admin</a></li>
						{/if}
						<li>
							<div class="dropdown dropdown-end dropdown-bottom dropdown-hover">
								<span tabindex="-1" role="button">Hello!</span>
								<ul
									tabindex="-1"
									class="menu dropdown-content z-[1] w-52 rounded-box bg-base-200 p-2 shadow"
								>
									<li><a href="/settings">Settings</a></li>
									<li>
										<form method="post" action="/logout" use:enhance>
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

	<main class="grow overflow-y-auto pt-3">
		<div class="flex w-full flex-col items-center justify-center" class:h-full={!isOverflowing}>
			<slot />
		</div>
	</main>

	<footer id="footer" class="footer footer-center flex-none bg-base-300 p-4 text-base-content">
		<aside>
			<p>PW © {new Date().getFullYear()} - All rights reserved</p>
			<p>Version: {version}</p>
			<label class="swap swap-rotate fixed bottom-3 right-3 z-50 m-2">
				<!-- this hidden checkbox controls the state -->
				<input type="checkbox" class="theme-controller" value="emerald" bind:checked={lightMode} />

				<!-- sun icon -->
				<svg
					class="swap-on h-10 w-10 fill-current"
					xmlns="http://www.w3.org/2000/svg"
					viewBox="0 0 24 24"
					><path
						d="M5.64,17l-.71.71a1,1,0,0,0,0,1.41,1,1,0,0,0,1.41,0l.71-.71A1,1,0,0,0,5.64,17ZM5,12a1,1,0,0,0-1-1H3a1,1,0,0,0,0,2H4A1,1,0,0,0,5,12Zm7-7a1,1,0,0,0,1-1V3a1,1,0,0,0-2,0V4A1,1,0,0,0,12,5ZM5.64,7.05a1,1,0,0,0,.7.29,1,1,0,0,0,.71-.29,1,1,0,0,0,0-1.41l-.71-.71A1,1,0,0,0,4.93,6.34Zm12,.29a1,1,0,0,0,.7-.29l.71-.71a1,1,0,1,0-1.41-1.41L17,5.64a1,1,0,0,0,0,1.41A1,1,0,0,0,17.66,7.34ZM21,11H20a1,1,0,0,0,0,2h1a1,1,0,0,0,0-2Zm-9,8a1,1,0,0,0-1,1v1a1,1,0,0,0,2,0V20A1,1,0,0,0,12,19ZM18.36,17A1,1,0,0,0,17,18.36l.71.71a1,1,0,0,0,1.41,0,1,1,0,0,0,0-1.41ZM12,6.5A5.5,5.5,0,1,0,17.5,12,5.51,5.51,0,0,0,12,6.5Zm0,9A3.5,3.5,0,1,1,15.5,12,3.5,3.5,0,0,1,12,15.5Z"
					/></svg
				>

				<!-- moon icon -->
				<svg
					class="swap-off h-10 w-10 fill-current"
					xmlns="http://www.w3.org/2000/svg"
					viewBox="0 0 24 24"
					><path
						d="M21.64,13a1,1,0,0,0-1.05-.14,8.05,8.05,0,0,1-3.37.73A8.15,8.15,0,0,1,9.08,5.49a8.59,8.59,0,0,1,.25-2A1,1,0,0,0,8,2.36,10.14,10.14,0,1,0,22,14.05,1,1,0,0,0,21.64,13Zm-9.5,6.69A8.14,8.14,0,0,1,7.08,5.22v.27A10.15,10.15,0,0,0,17.22,15.63a9.79,9.79,0,0,0,2.1-.22A8.11,8.11,0,0,1,12.14,19.73Z"
					/></svg
				>
			</label>
		</aside>
	</footer>
</div>
