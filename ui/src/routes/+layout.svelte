<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';
	import {
		Menu,
		X,
		Home,
		Activity,
		Settings,
		HelpCircle,
		ChevronLeft,
		ChevronRight
	} from 'lucide-svelte';
	import { writable } from 'svelte/store';

	const isNavbarCollapsed = writable(false);

	function toggleNavbar() {
		isNavbarCollapsed.update((n) => !n);
	}
</script>

<div class="bg-base-200 flex h-screen">
	<!-- Sidebar -->
	<aside
		class="bg-base-100 text-base-content transition-all duration-300 ease-in-out {$isNavbarCollapsed
			? 'w-16'
			: 'w-64'} flex flex-col"
	>
		<div class="flex items-center justify-between p-4">
			{#if !$isNavbarCollapsed}
				<span class="text-xl font-bold">EKOLOD</span>
			{/if}
			<button on:click={toggleNavbar} class="btn btn-square btn-ghost">
				{#if $isNavbarCollapsed}
					<ChevronRight />
				{:else}
					<ChevronLeft />
				{/if}
			</button>
		</div>
		<nav class="flex-1 overflow-y-auto">
			<ul class="menu p-2">
				<li>
					<a href="/" class="flex items-center gap-2 {$page.url.pathname === '/' ? 'active' : ''}">
						<Home />
						{#if !$isNavbarCollapsed}
							<span>Dashboard</span>
						{/if}
					</a>
				</li>
				<li>
					<a
						href="/services"
						class="flex items-center gap-2 {$page.url.pathname === '/services' ? 'active' : ''}"
					>
						<Activity />
						{#if !$isNavbarCollapsed}
							<span>Services</span>
						{/if}
					</a>
				</li>
				<li>
					<a
						href="/settings"
						class="flex items-center gap-2 {$page.url.pathname === '/settings' ? 'active' : ''}"
					>
						<Settings />
						{#if !$isNavbarCollapsed}
							<span>Settings</span>
						{/if}
					</a>
				</li>
				<li>
					<a
						href="/help"
						class="flex items-center gap-2 {$page.url.pathname === '/help' ? 'active' : ''}"
					>
						<HelpCircle />
						{#if !$isNavbarCollapsed}
							<span>Help</span>
						{/if}
					</a>
				</li>
			</ul>
		</nav>
	</aside>

	<!-- Main content -->
	<div class="flex flex-1 flex-col overflow-hidden">
		<!-- Top navbar for mobile -->
		<nav class="bg-base-100 p-4 shadow-md lg:hidden">
			<div class="flex items-center justify-between">
				<span class="text-xl font-bold">EKOLOD</span>
				<button on:click={toggleNavbar} class="btn btn-square btn-ghost">
					<Menu />
				</button>
			</div>
		</nav>

		<!-- Page content -->
		<main class="bg-base-200 flex-1 overflow-y-auto p-4">
			<slot />
		</main>

		<!-- Footer -->
		<footer class="bg-base-300 text-base-content p-4 text-center">
			<p>© {new Date().getFullYear()} Ekolod. All rights reserved.</p>
		</footer>
	</div>
</div>
