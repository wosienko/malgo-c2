<script lang="ts">
	import { page } from '$app/stores';
	import { get } from 'svelte/store';
	import { afterNavigate } from '$app/navigation';

	const ALL_TABS = ['/commands', '/modules', '/history', '/settings'] as const;

	const { id, sessionid } = get(page).params;

	let currentUrl = $state(get(page).url.pathname);

	afterNavigate(() => {
		currentUrl = get(page).url.pathname;
	});

	const capitalize = (s: string) => {
		return s.charAt(0).toUpperCase() + s.slice(1);
	};
</script>

<div class="flex h-full w-full flex-col pt-3">
	<div role="tablist" class="tabs tabs-bordered tabs-md w-full">
		{#each ALL_TABS as tab}
			<a
				role="tab"
				class="tab pb-9 hover:opacity-50"
				class:tab-active={currentUrl.endsWith(tab)}
				href={`/projects/${id}/${sessionid}${tab}`}
				>{capitalize(tab.replace('/', '').replace('-', ' '))}</a
			>
		{/each}
	</div>
	<div class="max-w-full grow overflow-x-auto">
		<slot />
	</div>
</div>
