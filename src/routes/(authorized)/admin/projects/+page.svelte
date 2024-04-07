<script lang="ts">
	import type { ZodIssue } from 'zod';
	import type { ApiError, Projects } from '$lib';
	import ZodIssues from '$lib/components/toasts/ZodIssues.svelte';
	import ApiIssues from '$lib/components/toasts/ApiIssues.svelte';
	import ChangeSuccessful from '$lib/components/toasts/ChangeSuccessful.svelte';
	import { pushState } from '$app/navigation';
	import { dateSchema, fieldSchema } from '$lib/validationSchemas';
	import ValidatedInput from '$lib/components/inputs/ValidatedInput.svelte';
	import ModalRunCancel from '$lib/components/modals/ModalRunCancel.svelte';
	import ValidatedInputWithLabel from '$lib/components/inputs/ValidatedInputWithHorizontalLabel.svelte';

	let zodIssues: ZodIssue[] = $state([]);
	let apiError: ApiError | undefined = $state();
	let successMessage: string | undefined = $state();

	let { data } = $props();

	let loading = $state(false);

	let page = $state(data.page);
	let pageSize = $state(data.pageSize);

	type Project = {
		id: string;
		name: string;
		startDate: string;
		endDate: string;
		description: string;
		editing: boolean;
	};

	const createEmptyProject = (): Project => ({
		id: '',
		name: '',
		startDate: '',
		endDate: '',
		description: '',
		editing: true
	});

	let projects: Project[] = $state(
		data.projects.projects.map((project) => ({
			...project,
			editing: false
		}))
	);

	let lastProjectValues: Project = $state(createEmptyProject());
	let projectToAlter: Project = $state(createEmptyProject());

	const loadPage = (nextPage: number) => {
		return async () => {
			const res = await fetch(`/api/project?page=${nextPage}&pageSize=${pageSize}`);
			pushState(`/admin/projects?page=${nextPage}&pageSize=${pageSize}`, {});
			const projectsPage: Projects = await res.json();
			projects = projectsPage.projects.map((project) => ({
				...project,
				editing: false
			}));
			data.projects.count = projectsPage.count;
			page = nextPage;
		};
	};

	const loadNextPage = async () => {
		await loadPage(page + 1)();
	};
	const reloadCurrentPage = async () => {
		await loadPage(page)();
	};
	const loadPreviousPage = async () => {
		await loadPage(page - 1)();
	};

	let newProject = $state({
		name: '',
		startDate: '',
		endDate: '',
		description: ''
	});

	const newProjectVerification = $derived({
		name: fieldSchema.safeParse(newProject.name),
		startDate: dateSchema.safeParse(newProject.startDate),
		endDate: dateSchema
			.refine(
				() => {
					return newProject.startDate < newProject.endDate;
				},
				{ message: 'End date must be after start date' }
			)
			.safeParse(newProject.endDate)
	});

	const newProjectValid = $derived(
		newProjectVerification.name.success &&
			newProjectVerification.startDate.success &&
			newProjectVerification.endDate.success
	);

	let showProjectCreationModal = $state(() => {});
	const createNewProject = async (): Promise<boolean> => {
		const res = await fetch('/api/project', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(newProject)
		});

		if (res.ok) {
			successMessage = 'Project created successfully!';
			newProject = createEmptyProject();
			await reloadCurrentPage();
			return true;
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
			return false;
		}
	};

	const updateProject = async (project: Project) => {
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
			project.editing = false;
			successMessage = 'Project updated successfully!';
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

	let showDeletionModal = $state(() => {});
	const prepareForDeletion = (user: typeof projectToAlter) => {
		return () => {
			projectToAlter = user;
			showDeletionModal();
		};
	};
	const deleteProject = async (project: Project): Promise<boolean> => {
		const res = await fetch(`/api/project/${project.id}`, {
			method: 'DELETE'
		});

		if (res.ok) {
			await reloadCurrentPage();

			successMessage = 'Project deleted successfully!';
			return true;
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
			return false;
		}
	};

	let checkedUsers = $state(
		data.users.map((user) => ({
			...user,
			checked: false
		}))
	);

	let currentlyAssignedOperators: string[] = $state([]);

	const isAnythingDifferentInCheckedUsers = $derived.by(() => {
		return checkedUsers.filter((user) => user.checked).length !== currentlyAssignedOperators.length;
	});

	const getAssignedOperators = async (project: Project) => {
		const res = await fetch(`/api/project/${project.id}/operators`);
		const users = await res.json();
		checkedUsers = data.users.map((user) => ({
			...user,
			checked: users.some((u) => u === user.id)
		}));
		currentlyAssignedOperators = users.map((user) => user.id);
	};

	let showOperatorAssignmentModal = $state(() => {});
	const prepareForOperatorAssignment = (project: Project) => {
		return async () => {
			projectToAlter = project;
			await getAssignedOperators(project);
			showOperatorAssignmentModal();
		};
	};

	const assignOperators = async () => {
		const res = await fetch(`/api/project/${projectToAlter.id}/operators`, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				users: checkedUsers.filter((user) => user.checked).map((user) => user.id)
			})
		});

		if (res.ok) {
			successMessage = 'Operators assigned successfully!';
			checkedUsers = checkedUsers.map((user) => ({
				...user,
				checked: false
			}));
			return true;
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
			return false;
		}
	};

	const formatDate = (date: string): string => {
		const d = new Date(date);
		// format DD.MM.YYYY
		return `${d.getDate().toString().padStart(2, '0')}.${(d.getMonth() + 1).toString().padStart(2, '0')}.${d.getFullYear()}`;
	};
</script>

{#if zodIssues.length > 0}
	<ZodIssues
		issues={zodIssues}
		on:close={() => {
			zodIssues = [];
		}}
	/>
{/if}

{#if apiError}
	<ApiIssues
		issue={apiError}
		on:close={() => {
			apiError = undefined;
		}}
	/>
{/if}

{#if successMessage}
	<ChangeSuccessful
		message={successMessage}
		on:close={() => {
			successMessage = undefined;
		}}
	/>
{/if}

<ModalRunCancel
	id="confirm-deletion"
	title="Are you sure you want to delete?"
	message="Project to be deleted: "
	messageEmphasis={`${projectToAlter.name}`}
	btnClass="btn-error"
	btnText="Delete"
	onclickCallback={() => deleteProject(projectToAlter)}
	bind:showModal={showDeletionModal}
/>

<ModalRunCancel
	id="operator-assignment"
	title="Assign operators"
	message="Assign operators to: "
	messageEmphasis={`${projectToAlter.name}`}
	btnClass="btn-warning"
	btnText="Assign operators"
	btnDisabledCondition={!isAnythingDifferentInCheckedUsers}
	onclickCallback={assignOperators}
	bind:showModal={showOperatorAssignmentModal}
	onHideModal={() => {
		checkedUsers = checkedUsers.map((user) => ({
			...user,
			checked: false
		}));
	}}
>
	<div class="grid w-full grid-cols-2">
		{#each checkedUsers as user}
			<div class="form-control">
				<label class="label cursor-pointer">
					<span class="label-text">{user.name} {user.surname}</span>
					<input type="checkbox" bind:checked={user.checked} class="checkbox-info checkbox" />
				</label>
			</div>
		{/each}
	</div>
</ModalRunCancel>

<ModalRunCancel
	id="create-new-project"
	title="Create new project"
	btnClass="btn-success"
	btnText="Create new project"
	btnDisabledCondition={!newProjectValid}
	onclickCallback={createNewProject}
	bind:showModal={showProjectCreationModal}
	onHideModal={() => {
		newProject = {
			name: '',
			startDate: '',
			endDate: ''
		};
	}}
>
	<div class="space-y-8">
		<ValidatedInputWithLabel
			label="Name"
			type="text"
			id="name"
			name="name"
			bind:value={newProject.name}
			validation={newProjectVerification.name}
			classes="w-full max-w-xs"
		/>
		<ValidatedInputWithLabel
			label="Start date"
			type="date"
			id="start-date"
			name="startDate"
			bind:value={newProject.startDate}
			validation={newProjectVerification.startDate}
			classes="w-48"
		/>
		<ValidatedInputWithLabel
			label="End date"
			type="date"
			id="end-date"
			name="endDate"
			bind:value={newProject.endDate}
			validation={newProjectVerification.endDate}
			classes="w-48"
		/>
		<div class="flex items-center justify-between">
			<label for="new-project-description">Project description</label>
			<div class="flex flex-col items-center w-[45%]">
				<textarea id="new-project-description" bind:value={newProject.description} class="mt-1.5 textarea resize-none textarea-bordered textarea-xs w-full max-w-xs" ></textarea>
			</div>
		</div>
	</div>
</ModalRunCancel>

<div class="mt-3 flex w-full justify-center space-x-3">
	<div class="hidden flex-1 text-center md:block">
		<button class="btn btn-sm -mt-1" on:click={showProjectCreationModal}>Create new project</button>
	</div>
	<div class="join -mt-1">
		<button class="btn join-item btn-sm" class:btn-disabled={page === 1} on:click={loadPreviousPage}
			>«</button
		>
		<button class="btn join-item btn-sm hover:cursor-auto hover:border-base-200 hover:bg-base-200"
			>Page {page}</button
		>
		<button
			class="btn join-item btn-sm"
			class:btn-disabled={page * pageSize >= data.projects.count}
			on:click={loadNextPage}>»</button
		>
	</div>
	<div class="mt-1 md:flex-1">
		<div class="flex items-center justify-center space-x-3">
			<span>Limit: {pageSize}</span>
			<input
				type="range"
				min="1"
				max="20"
				bind:value={pageSize}
				class="range range-sm max-w-32"
				on:mouseup={reloadCurrentPage}
				on:touchend={reloadCurrentPage}
			/>
		</div>
	</div>
</div>

<table class="table table-fixed overflow-x-auto text-center">
	<thead>
		<tr>
			<th>Name</th>
			<th class="hidden md:table-cell">Start date</th>
			<th class="hidden md:table-cell">End date</th>
			<th class="hidden md:table-cell">Description</th>
			<th class="w-28 md:w-auto"></th>
		</tr>
	</thead>
	<tbody>
		{#each projects as project, i}
			{#if !project.editing}
				<tr class="max-w-dvw hover">
					<td class="no-scrollbar overflow-x-auto">{project.name}</td>
					<td class="no-scrollbar hidden overflow-x-auto md:table-cell"
						>{formatDate(project.startDate)}</td
					>
					<td class="no-scrollbar hidden overflow-x-auto md:table-cell"
						>{formatDate(project.endDate)}</td
					>
					<td class="no-scrollbar hidden overflow-x-auto md:table-cell whitespace-pre-line"
					>{project.description}</td
					>
					<td>
						<!-- Dropdown accounting for hiding under the dropdown for last elements -->
						<div
							class="dropdown dropdown-end"
							class:dropdown-bottom={i <= 2 ||
								i + 1 < (data.projects.count > pageSize ? pageSize : data.projects.count) - 2}
							class:dropdown-top={i > 2 &&
								i + 1 >= (data.projects.count > pageSize ? pageSize : data.projects.count) - 2}
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
											project.editing = true;
											lastProjectValues = { ...project };
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
									<button
										class="btn btn-error btn-sm"
										on:click={() => prepareForDeletion(project)()}>Delete</button
									>
								</li>
							</ul>
						</div>
					</td>
				</tr>
			{:else}
				{@const nameCheck = fieldSchema.safeParse(project.name)}
				{@const startDateCheck = dateSchema.safeParse(project.startDate)}
				{@const endDateCheck = dateSchema
					.refine(
						() => {
							return project.startDate < project.endDate;
						},
						{
							message: 'End date must be after start date'
						}
					)
					.safeParse(project.endDate)}
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
								label="Start date"
								type="date"
								id="start-date"
								name="startDate"
								bind:value={project.startDate}
								validation={startDateCheck}
								classes="w-48 input-sm mb-5"
							/>
						</div>
					</td>
					<td class="hidden md:table-cell">
						<div class="md:-mt-2">
							<ValidatedInput
								label="End date"
								type="date"
								id="end-date"
								name="endDate"
								bind:value={project.endDate}
								validation={endDateCheck}
								classes="w-48 input-sm mb-5"
							/>
						</div>
					</td>
					<td class="hidden md:table-cell">
						<div class="md:-mt-2">
							<textarea bind:value={project.description} class="mt-1.5 textarea resize-none textarea-bordered textarea-xs w-full max-w-xs" ></textarea>
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
								class:btn-disabled={JSON.stringify(project) === JSON.stringify(lastProjectValues) ||
									!editingValid}
								on:click={() => updateProject(project)}>Save</button
							>
							<button
								class="btn btn-sm"
								on:click={() => {
									// so that the input fields are reset. In runes mode, cannot use `project = lastProjectValues`
									projects[i] = lastProjectValues;
									projects[i].editing = false;
								}}>Cancel</button
							>
						{/if}
					</td>
				</tr>
			{/if}
		{/each}
	</tbody>
</table>
