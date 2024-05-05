<script lang="ts">
	import '../main.css';
	import { onMount } from 'svelte';
	import { navigating } from '$app/stores';
	import Footer from '$lib/components/Footer.svelte';
	import Header from '$lib/components/Header.svelte';

	let { data } = $props();

	// main tag overflow check
	let isOverflowing = $state(false);
	let contentElement: Element | null = null;
	const checkOverflow = () => {
		isOverflowing = contentElement!.scrollHeight > contentElement!.clientHeight;
	};

	let logoutTimeout: ReturnType<typeof setTimeout> | null = null;
	const logoutTimeoutTime = 29 * 60 * 1000;
	const logoutAutomatically = () => {
		// get form element
		const form = document.getElementById('logout-form') as HTMLFormElement;
		// submit form
		form.submit();
	};

	onMount(() => {
		// check if content is overflowing
		contentElement = document.querySelector('main');
		checkOverflow();
		window.addEventListener('resize', checkOverflow);
		let unsubscribe = navigating.subscribe(checkOverflow);

		// logout after 29 (to not have problems with 30min session) minutes of inactivity
		if (data.loggedIn) {
			logoutTimeout = setTimeout(logoutAutomatically, logoutTimeoutTime);
		}
		window.addEventListener('mousemove', () => {
			if (logoutTimeout) clearTimeout(logoutTimeout);
			logoutTimeout = setTimeout(logoutAutomatically, logoutTimeoutTime);
		});

		// cleanup
		return () => {
			window.removeEventListener('resize', checkOverflow);
			if (logoutTimeout) clearTimeout(logoutTimeout);
			unsubscribe();
		};
	});
</script>

<div class="flex h-dvh flex-col">
	<Header isLoggedIn={data.loggedIn} isAdmin={data.isAdmin} isOperator={data.isOperator} />

	<main class="grow overflow-y-auto">
		<div class="flex w-full flex-col items-center justify-center" class:h-full={!isOverflowing}>
			<slot />
		</div>
	</main>

	<Footer />
</div>
