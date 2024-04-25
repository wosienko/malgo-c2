<script lang="ts">
	type InputProps = {
		name: string;
		createdAt: string;
		heartbeatAt: string;
	};

	let { name, createdAt, heartbeatAt }: InputProps = $props();

	const olderThanInSeconds = (date: string, seconds: number): boolean => {
		const d = new Date(date);
		const now = new Date();
		return (now.getTime() - d.getTime()) / 1000 > seconds;
	};

	const formatDateAndTime = (date: string): string => {
		const d = new Date(date);
		// format DD.MM.YYYY HH:MM:SS
		return `${d.getDate().toString().padStart(2, '0')}.${(d.getMonth() + 1).toString().padStart(2, '0')}.${d.getFullYear()} ${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}:${d.getSeconds().toString().padStart(2, '0')}`;
	};
</script>

<div class="flex flex-col">
	<span class="text-xl text-info">{name}</span>
	<table class="text-left">
		<tbody>
			<tr>
				<td>Created at:</td>
				<td>{formatDateAndTime(createdAt)}</td>
			</tr>
			<tr>
				<td>Last Heartbeat:</td>
				<td class={olderThanInSeconds(heartbeatAt, 3600) ? 'text-error' : 'text-success'}
					>{formatDateAndTime(heartbeatAt)}</td
				>
			</tr>
		</tbody>
	</table>
</div>
