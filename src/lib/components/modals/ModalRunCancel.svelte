<script lang="ts">
	type ModalProps = {
		id: string;
		title: string;
		message?: string;
		messageEmphasis?: string;
		btnClass?: string;
		btnDisabledCondition?: boolean;
		btnText: string;
		onclickCallback: () => Promise<boolean>;
		showModal: () => void;
	};

	let {
		id,
		title,
		message = '',
		messageEmphasis = '',
		btnClass = '',
		btnDisabledCondition = false,
		btnText,
		onclickCallback,
		// eslint-disable-next-line no-undef
		showModal = $bindable()
	}: ModalProps = $props();

	// eslint-disable-next-line @typescript-eslint/no-unused-vars
	showModal = () => {
		const dialog = document.getElementById(id) as HTMLDialogElement;
		dialog.showModal();
	};

	const cleanupClose = () => {
		const dialog = document.getElementById(id) as HTMLDialogElement;
		dialog.close();
	};

	let waitingForResponse = $state(false);
</script>

<dialog {id} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box">
		<h3 class="text-lg font-bold" class:mb-3={!message && !messageEmphasis}>{title}</h3>
		{#if message || messageEmphasis}
			<p class="py-4">
				{#if message}
					{message} <span class="font-bold">{messageEmphasis}</span>
				{:else}
					<span class="font-bold">{messageEmphasis}</span>
				{/if}
			</p>
		{/if}
		<form method="dialog">
			<slot />
			{#if waitingForResponse}
				<div class="mt-6 h-12 w-full text-center">
					<span class="loading loading-spinner loading-lg text-info"></span>
				</div>
			{:else}
				<div class="modal-action space-x-3">
					<button
						class={`btn ${btnClass}`}
						class:btn-disabled={btnDisabledCondition}
						on:click|preventDefault={async () => {
							waitingForResponse = true;
							let result = await onclickCallback();
							waitingForResponse = false;
							if (result) cleanupClose();
						}}
					>
						{btnText}
					</button>
					<button class="btn" class:btn-disabled={waitingForResponse}>Cancel</button>
				</div>
			{/if}
		</form>
	</div>
</dialog>
