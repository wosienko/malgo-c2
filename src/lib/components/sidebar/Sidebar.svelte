<div class="drawer h-full lg:drawer-open">
	<input id="my-drawer-2" type="checkbox" class="drawer-toggle" />
	<div class="drawer-content flex flex-col items-center">
		<slot name="content" />
		<label
			for="my-drawer-2"
			class="btn btn-primary drawer-button fixed bottom-3 left-3 mb-1 ml-3 lg:hidden">&gt;</label
		>
	</div>
	<div class="drawer-side lg:h-full">
		<label for="my-drawer-2" aria-label="close sidebar" class="drawer-overlay"></label>
		<ul class="menu min-h-full w-80 bg-base-200 p-4 text-base-content">
			<slot name="menu" />
		</ul>
	</div>
</div>
