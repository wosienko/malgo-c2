<script lang="ts">
	import { page } from '$app/stores';
	import { get } from 'svelte/store';
	import { afterNavigate } from '$app/navigation';
	const ALL_TABS = ['/admin/users'] as const;

	let currentUrl = $state(get(page).url.pathname);

	afterNavigate(() => {
		currentUrl = get(page).url.pathname;
	});

	const capitalize = (s: string) => {
		return s.charAt(0).toUpperCase() + s.slice(1);
	};
</script>

<div class="flex h-full w-full flex-col">
	<div role="tablist" class="tabs tabs-bordered tabs-md">
		{#each ALL_TABS as tab}
			<a
				role="tab"
				class="tab pb-9 hover:opacity-50"
				class:tab-active={currentUrl === tab}
				href={tab}>{capitalize(tab.replace('/admin/', ''))}</a
			>
		{/each}
	</div>
	<div class="grow">
		<slot />
	</div>
</div>
