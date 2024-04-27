<script lang="ts">
	import type { Command } from '$lib/components/custom/command/command';

	let isResultVisible = $state(true);
	let isCommandVisible = $state(true);

	type InputProps = {
		command: Command;
	};

	let { command }: InputProps = $props();
</script>

<div class="mx-3 rounded-2xl border-2 border-neutral bg-base-100 p-4 shadow-xl">
	<div>
		<div class="mb-3 flex justify-between">
			<h2 class="card-title">
				<span class="badge badge-primary">{command.type.toUpperCase()}</span>
				<span class="badge badge-secondary">{command.status.toUpperCase()}</span>
			</h2>
			<div class="flex flex-col justify-around space-y-1.5 md:flex-row md:space-x-1 md:space-y-0">
				<button class="btn btn-sm" on:click={() => (isCommandVisible = !isCommandVisible)}
					>Toggle command</button
				>
				<button
					class="btn btn-sm"
					class:btn-disabled={!command.result}
					on:click={() => (isResultVisible = !isResultVisible)}>Toggle result</button
				>
				<button class="btn btn-sm">Cancel</button>
			</div>
		</div>
		<p>Operator: {command.operator}</p>
		<table class="w-full">
			<tbody>
				<tr>
					<td class="text-left">Created:</td>
					<td class="text-right">{command.created_at}</td>
				</tr>
				<tr>
					<td class="text-left">Last result update:</td>
					<td class="text-right">08.04.2024 12:15:36</td>
				</tr>
			</tbody>
		</table>
		{#if isCommandVisible}
			<div class="divider m-0"></div>
			<p class="whitespace-pre-line">{command.command}</p>
		{/if}
		{#if isResultVisible && command.result}
			<div class="divider m-0"></div>
			<p class="whitespace-pre-line">
				{command.result}
			</p>
		{/if}
	</div>
</div>
