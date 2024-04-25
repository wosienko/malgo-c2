<script>
	import Command from '$lib/components/custom/command/Command.svelte';
</script>

<div class="flex h-1/2 w-full justify-around border-b-2 border-neutral pb-6">
	<div class="m-3 flex h-full flex-1 flex-col">
		<textarea class="textarea textarea-warning w-full grow resize-none" placeholder="Command"
		></textarea>
		<button class="btn btn-warning mt-3 w-full">Run</button>
	</div>
	<div class="m-3 hidden h-full flex-1 overflow-auto md:block">
		<table class="table table-fixed">
			<thead>
				<tr>
					<th class="w-auto">Key</th>
					<th class="w-auto">Value</th>
					<th class="w-24"></th>
				</tr>
			</thead>
			<tbody>
				<tr>
					<td>key</td>
					<td>value</td>
					<td>
						<button class="btn btn-sm">Edit</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td>
						<button class="btn btn-sm">Edit</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td>
						<button class="btn btn-sm">Edit</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td>
						<button class="btn btn-sm">Edit</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td>
						<button class="btn btn-sm">Edit</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td>
						<button class="btn btn-sm">Edit</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td>
						<button class="btn btn-sm">Edit</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td>
						<button class="btn btn-sm">Edit</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td>
						<button class="btn btn-sm">Edit</button>
					</td>
				</tr>
				<tr>
					<td>key</td>
					<td>value</td>
					<td>
						<button class="btn btn-sm">Edit</button>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</div>
<div class="h-1/2 w-full overflow-y-auto py-3">
	<Command />
</div>
