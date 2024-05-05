<script lang="ts">
	import type { ApiError } from '$lib';
	import type { ZodIssue } from 'zod';
	import { dateSchema, fieldSchema } from '$lib/validationSchemas';
	import ModalRunCancel from '$lib/components/modals/ModalRunCancel.svelte';
	import ValidatedInputWithLabel from '$lib/components/inputs/ValidatedInputWithHorizontalLabel.svelte';

	type InputProps = {
		successMessage: string | undefined;
		zodIssues: ZodIssue[];
		apiError: ApiError | undefined;
		reloadCurrentPage: () => Promise<void>;
		showModal: () => void;
	};

	let {
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		successMessage = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		zodIssues = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		apiError = $bindable(),
		reloadCurrentPage,
		showModal = $bindable()
	}: InputProps = $props();

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
					return new Date(newProject.startDate) < new Date(newProject.endDate);
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
			newProject = {
				name: '',
				startDate: '',
				endDate: '',
				description: ''
			};
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
</script>

<ModalRunCancel
	id="create-new-project"
	title="Create new project"
	btnClass="btn-success"
	btnText="Create new project"
	btnDisabledCondition={!newProjectValid}
	onclickCallback={createNewProject}
	bind:showModal
	onHideModal={() => {
		newProject = {
			name: '',
			startDate: '',
			endDate: '',
			description: ''
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
			<div class="flex w-[45%] flex-col items-center">
				<textarea
					id="new-project-description"
					bind:value={newProject.description}
					class="textarea textarea-bordered textarea-xs mt-1.5 w-full max-w-xs resize-none"
				></textarea>
			</div>
		</div>
	</div>
</ModalRunCancel>
