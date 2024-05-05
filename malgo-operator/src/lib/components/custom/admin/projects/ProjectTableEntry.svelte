<script lang="ts">
	import { dateSchema, fieldSchema } from '$lib/validationSchemas';
	import type { ProjectTableEntryType } from '$lib/components/custom/admin/projects/types';
	import ValidatedInput from '$lib/components/inputs/ValidatedInput.svelte';
	import type { ApiError } from '$lib';
	import type { ZodIssue } from 'zod';

	type InputProps = {
		project: ProjectTableEntryType;
		elementIndex: number;
		count: number;
		pageSize: number;
		successMessage: string | undefined;
		apiError: ApiError | undefined;
		zodIssues: ZodIssue[];
		prepareForOperatorAssignment: (project: ProjectTableEntryType) => () => Promise<void>;
		prepareForDeletion: (project: ProjectTableEntryType) => () => void;
	};

	let {
		project = $bindable(),
		elementIndex,
		count,
		pageSize,
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		successMessage = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		apiError = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		zodIssues = $bindable(),
		prepareForOperatorAssignment,
		prepareForDeletion
	}: InputProps = $props();

	let editing = $state(false);
	let loading = $state(false);

	// deep copy of the project object to store the last values
	let lastProjectValues: ProjectTableEntryType = $state(JSON.parse(JSON.stringify(project)));

	const updateProject = async () => {
		loading = true;
		const res = await fetch(`/api/project/${project.id}`, {
			method: 'PATCH',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				name: project.name,
				startDate: project.startDate,
				endDate: project.endDate,
				description: project.description
			})
		});
		loading = false;

		if (res.ok) {
			editing = false;
			successMessage = 'Project updated successfully!';
			lastProjectValues = JSON.parse(JSON.stringify(project));
		} else {
			try {
				let body: ApiError | ZodIssue[] = await res.json();

				if (Array.isArray(body)) {
					zodIssues = body;
				} else {
					apiError = body;
				}
			} catch (e) {
				console.error(e);
			}
		}
	};

	const formatDate = (date: string): string => {
		const d = new Date(date);
		// format DD.MM.YYYY
		return `${d.getDate().toString().padStart(2, '0')}.${(d.getMonth() + 1).toString().padStart(2, '0')}.${d.getFullYear()}`;
	};
</script>

{#if !editing}
	<tr class="max-w-dvw hover">
		<td class="no-scrollbar overflow-x-auto">{project.name}</td>
		<td class="no-scrollbar hidden overflow-x-auto md:table-cell"
			>{formatDate(project.startDate)}</td
		>
		<td class="no-scrollbar hidden overflow-x-auto md:table-cell">{formatDate(project.endDate)}</td>
		<td class="no-scrollbar hidden overflow-x-auto whitespace-pre-line md:table-cell"
			>{project.description}</td
		>
		<td>
			<!-- Dropdown accounting for hiding under the dropdown for last elements -->
			<div
				class="dropdown dropdown-end"
				class:dropdown-bottom={elementIndex <= 2 ||
					elementIndex + 1 < (count > pageSize ? pageSize : count) - 2}
				class:dropdown-top={elementIndex > 2 &&
					elementIndex + 1 >= (count > pageSize ? pageSize : count) - 2}
			>
				<div tabindex="-1" role="button" class="btn btn-neutral btn-sm mb-1">Options</div>
				<ul
					tabindex="-1"
					class="menu dropdown-content z-[1] w-52 space-y-1.5 rounded-box bg-base-100 p-2 shadow"
				>
					<li>
						<button
							class="btn btn-sm"
							on:click={() => {
								editing = true;
							}}>Edit</button
						>
					</li>
					<li>
						<button
							class="btn btn-warning btn-sm"
							on:click={() => {
								prepareForOperatorAssignment(project)();
							}}>Assign operators</button
						>
					</li>
					<li>
						<button class="btn btn-error btn-sm" on:click={() => prepareForDeletion(project)()}
							>Delete</button
						>
					</li>
				</ul>
			</div>
		</td>
	</tr>
{:else}
	{@const nameCheck = fieldSchema.safeParse(project.name)}
	{@const startDateCheck = dateSchema.safeParse(project.startDate)}
	{@const projectStartDate = new Date(project.startDate)}
	{@const projectEndDate = new Date(project.endDate)}
	{@const endDateCheck = dateSchema
		.refine(
			() => {
				return projectStartDate < projectEndDate;
			},
			{
				message: 'End date must be after start date'
			}
		)
		.safeParse(project.endDate)}
	{@const areValuesDifferent =
		project.name !== lastProjectValues.name ||
		project.startDate !== lastProjectValues.startDate ||
		project.endDate !== lastProjectValues.endDate ||
		project.description !== lastProjectValues.description}
	{@const editingValid = nameCheck.success && startDateCheck.success && endDateCheck.success}
	<tr class="max-w-dvw hover">
		<td>
			<div class="md:-mt-2">
				<ValidatedInput
					type="text"
					id="name"
					name="name"
					bind:value={project.name}
					validation={nameCheck}
					classes="input-sm mb-5 w-full"
				/>
			</div>
		</td>
		<td class="hidden md:table-cell">
			<div class="md:-mt-2">
				<ValidatedInput
					type="date"
					id="start-date"
					name="startDate"
					bind:value={project.startDate}
					validation={startDateCheck}
					classes="input-sm mb-5"
				/>
			</div>
		</td>
		<td class="hidden md:table-cell">
			<div class="md:-mt-2">
				<ValidatedInput
					type="date"
					id="end-date"
					name="endDate"
					bind:value={project.endDate}
					validation={endDateCheck}
					classes="input-sm mb-5"
				/>
			</div>
		</td>
		<td class="hidden md:table-cell">
			<div class="md:-mt-2">
				<textarea
					bind:value={project.description}
					class="textarea textarea-bordered textarea-xs mt-1.5 w-full max-w-xs resize-none"
				></textarea>
			</div>
		</td>
		<td class="my-1.5 flex flex-col items-center justify-center space-y-3">
			{#if loading}
				<div class="my-4 h-11">
					<span class="loading loading-spinner loading-lg text-info"></span>
				</div>
			{:else}
				<button
					class="btn btn-success btn-sm w-14"
					class:btn-disabled={!areValuesDifferent || !editingValid}
					on:click={updateProject}>Save</button
				>
				<button
					class="btn btn-sm"
					on:click={() => {
						editing = false;
						project = lastProjectValues;
					}}>Cancel</button
				>
			{/if}
		</td>
	</tr>
{/if}
