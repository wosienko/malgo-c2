<script lang="ts">
	import ModalRunCancel from '$lib/components/modals/ModalRunCancel.svelte';
	import type { ApiError, UserWithRoles } from '$lib';
	import type { ZodIssue } from 'zod';

	type InputProps = {
		userToAlter: UserWithRoles;
		showModal: () => void;
		reloadCurrentPage: () => Promise<void>;
		successMessage: string | undefined;
		zodIssues: ZodIssue[];
		apiError: ApiError | undefined;
	};

	let {
		userToAlter,
		showModal = $bindable(),
		reloadCurrentPage,
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		successMessage = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		zodIssues = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		apiError = $bindable()
	}: InputProps = $props();

	const deleteUser = async (user: UserWithRoles): Promise<boolean> => {
		const res = await fetch(`/api/user/${user.id}`, {
			method: 'DELETE'
		});

		if (res.ok) {
			await reloadCurrentPage();

			successMessage = 'User deleted successfully!';
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
	message="Account to be deleted: "
	messageEmphasis={`${userToAlter.name} ${userToAlter.surname}`}
	btnClass="btn-error"
	btnText="Delete"
	onclickCallback={() => deleteUser(userToAlter)}
	bind:showModal
/>
