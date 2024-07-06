<script>
	import { onMount } from 'svelte';
	import '../app.css';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import TopBar from '$lib/components/TopBar.svelte';

	const { children } = $props();
	let isDrawerOpen = $state(false);

	function toggleDrawer() {
		isDrawerOpen = !isDrawerOpen;
	}

	onMount(() => {
		const savedTheme = localStorage.getItem('theme') || 'light';
		document.documentElement.setAttribute('data-theme', savedTheme);
	});
</script>

<div class="drawer lg:drawer-open">
	<input id="main-drawer" type="checkbox" class="drawer-toggle" bind:checked={isDrawerOpen} />

	<div class="drawer-content flex flex-col">
		<TopBar {toggleDrawer} />
		<main class="flex-1 overflow-y-auto bg-base-100 p-4">
			{@render children()}
		</main>
	</div>

	<div class="drawer-side">
		<label for="main-drawer" class="drawer-overlay" aria-label="Close menu"></label>
		<Sidebar />
	</div>
</div>
