<script lang="ts">
	import Sidebar from '$lib/components/sidebar/Sidebar.svelte';
	import SidebarEntry from '$lib/components/sidebar/SidebarEntry.svelte';
</script>

<Sidebar>
	<div slot="content"><p>Here is content</p></div>
	<svelte:fragment slot="menu">
		<SidebarEntry href="/">Entry 1</SidebarEntry>
	</svelte:fragment>
</Sidebar>
