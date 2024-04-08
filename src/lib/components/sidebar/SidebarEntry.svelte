<script lang="ts">
	type Props = {
		href: string;
	};

	let { href }: Props = $props();
</script>

<li>
	<a {href}>
		<slot />
	</a>
</li>
