<script>
	import '../app.css';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import TopBar from '$lib/components/TopBar.svelte';

	const { children } = $props();
	let isDrawerOpen = $state(false);

	function toggleDrawer() {
		isDrawerOpen = !isDrawerOpen;
	}
</script>

<TopBar {toggleDrawer} />

<div class="drawer lg:drawer-open pt-16">
	<input id="main-drawer" type="checkbox" class="drawer-toggle" bind:checked={isDrawerOpen} />

	<div class="drawer-content flex flex-col">
		<main class="flex-1 overflow-y-auto bg-base-100 p-4">
			{@render children()}
		</main>
	</div>

	<div class="drawer-side">
		<label for="main-drawer" class="drawer-overlay" aria-label="Close menu"></label>
		<Sidebar />
	</div>
</div>
