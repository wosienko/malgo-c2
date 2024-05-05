<script lang="ts">
	import ModalRunCancel from '$lib/components/modals/ModalRunCancel.svelte';
	import type {
		CheckedUserType,
		ProjectTableEntryType
	} from '$lib/components/custom/admin/projects/types';
	import type { ZodIssue } from 'zod';
	import type { ApiError, UserWithRoles } from '$lib';

	type InputProps = {
		projectToAlter: ProjectTableEntryType;
		users: UserWithRoles[];
		showOperatorAssignmentModal: () => void;
		successMessage: string | undefined;
		zodIssues: ZodIssue[];
		apiError: ApiError | undefined;
	};

	let {
		projectToAlter,
		users,
		showOperatorAssignmentModal = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		successMessage = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		zodIssues = $bindable(),
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		apiError = $bindable()
	}: InputProps = $props();

	let checkedUsers: CheckedUserType[] = $state(
		users.map((user) => ({
			user,
			checked: false,
			changed: false
		}))
	);

	const isAnyUserChanged = $derived(checkedUsers.some((user) => user.changed));

	const getAssignedOperators = async (project: ProjectTableEntryType) => {
		const res = await fetch(`/api/project/${project.id}/operators`).then((res) => res.json());

		checkedUsers = users.map((user) => ({
			user,
			checked: res.some((u: string) => u === user.id),
			changed: false
		}));
	};

	let showModal = $state(() => {});
	// eslint-disable-next-line @typescript-eslint/no-unused-vars
	showOperatorAssignmentModal = () => {
		getAssignedOperators(projectToAlter);
		showModal();
	};

	const assignOperators = async () => {
		const res = await fetch(`/api/project/${projectToAlter.id}/operators`, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				users: checkedUsers.filter((user) => user.checked).map((user) => user.user.id)
			})
		});

		if (res.ok) {
			successMessage = 'Operators assigned successfully!';
			checkedUsers = checkedUsers.map((user) => ({
				...user,
				checked: false,
				changed: false
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
</script>

<ModalRunCancel
	id="operator-assignment"
	title="Assign operators"
	message="Assign operators to: "
	messageEmphasis={`${projectToAlter.name}`}
	btnClass="btn-warning"
	btnText="Assign operators"
	btnDisabledCondition={!isAnyUserChanged}
	onclickCallback={assignOperators}
	bind:showModal
	onHideModal={() => {
		for (const user of checkedUsers) {
			user.checked = false;
			user.changed = false;
		}
	}}
>
	<div class="grid w-full grid-cols-2">
		{#each checkedUsers as user}
			<div class="form-control">
				<label class="label cursor-pointer">
					<span class="label-text">{user.user.name} {user.user.surname}</span>
					<input
						type="checkbox"
						bind:checked={user.checked}
						on:click={() => {
							user.changed = !user.changed;
						}}
						class="checkbox-info checkbox"
					/>
				</label>
			</div>
		{/each}
	</div>
</ModalRunCancel>
