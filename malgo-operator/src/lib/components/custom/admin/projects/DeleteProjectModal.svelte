<script lang="ts">
	import ModalRunCancel from '$lib/components/modals/ModalRunCancel.svelte';
	import type { ProjectTableEntryType } from '$lib/components/custom/admin/projects/types';
	import type { ApiError } from '$lib';
	import type { ZodIssue } from 'zod';

	type InputProps = {
		projectToAlter: ProjectTableEntryType;
		showDeletionModal: () => void;
		reloadCurrentPage: () => Promise<void>;
		successMessage: string | undefined;
		zodIssues: ZodIssue[];
		apiError: ApiError | undefined;
	};

	let {
		projectToAlter,
		showDeletionModal = $bindable(),
		reloadCurrentPage,
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		successMessage = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		zodIssues = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		apiError = $bindable()
	}: InputProps = $props();

	const deleteProject = async (project: ProjectTableEntryType): Promise<boolean> => {
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
</script>

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
