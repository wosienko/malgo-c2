<script lang="ts">
	type Props = {
		href: string;
		active: boolean;
	};

	let { href, active = false }: Props = $props();
</script>

<li class={active ? 'rounded-lg bg-neutral text-neutral-content' : ''}>
	<a {href}>
		<slot />
	</a>
</li>
