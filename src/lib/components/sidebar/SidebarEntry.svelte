<script lang="ts">
	type Props = {
		href: string;
		active: boolean;
	};

	let { href, active = false }: Props = $props();
</script>

<li class={active ? 'bg-neutral text-neutral-content rounded-md': ''}>
	<a {href}>
		<slot />
	</a>
</li>
