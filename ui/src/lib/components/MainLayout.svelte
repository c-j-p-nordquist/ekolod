<script>
	import Sidebar from './Sidebar.svelte';
	import TopBar from './TopBar.svelte';

	let isDrawerOpen = $state(false);

	function toggleDrawer() {
		isDrawerOpen = !isDrawerOpen;
	}

	const { children } = $props();
</script>

<div class="drawer lg:drawer-open">
	<input id="main-drawer" type="checkbox" class="drawer-toggle" bind:checked={isDrawerOpen} />

	<div class="drawer-content flex flex-col">
		<TopBar {toggleDrawer} />

		<main class="flex-1 overflow-y-auto bg-base-100 p-4 lg:ml-64">
			{@render children()}
		</main>
	</div>

	<div class="drawer-side">
		<label for="main-drawer" class="drawer-overlay" aria-label="Close menu"></label>
		<Sidebar />
	</div>
</div>
