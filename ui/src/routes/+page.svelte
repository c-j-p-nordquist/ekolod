<script lang="ts">
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';
	import {
		targets,
		metrics,
		fetchTargets,
		fetchMetrics,
		type Metrics,
		type MetricData
	} from '$lib/stores/targets';
	import { writable } from 'svelte/store';
	import { Grid, List, BarChart, RefreshCw, AlertCircle } from 'lucide-svelte';
	import ServiceGridItem from '$lib/components/ServiceGridItem.svelte';
	import ServiceTable from '$lib/components/ServiceTable.svelte';

	let loading = true;
	let error = '';
	let viewMode = writable('grid');
	let refreshing = false;
	let notification: { message: string; type: 'success' | 'error' } | null = null;

	function generateMockTimeSeriesData(count: number) {
		return Array.from({ length: count }, () => Math.random() * 100);
	}

	$: timeSeriesData = Object.fromEntries(
		Object.keys($metrics).map((service) => [service, generateMockTimeSeriesData(10)])
	);

	async function checkHealth() {
		try {
			const response = await fetch('http://localhost:8080/health');
			return response.ok;
		} catch {
			return false;
		}
	}

	async function fetchDataWithRetry(retries = 5, interval = 5000) {
		for (let i = 0; i < retries; i++) {
			if (await checkHealth()) {
				try {
					await fetchTargets();
					await fetchMetrics();
					loading = false;
					error = '';
					return true;
				} catch (e) {
					console.error('Error fetching data:', e);
				}
			}
			await new Promise((resolve) => setTimeout(resolve, interval));
		}
		loading = false;
		error = 'Failed to fetch data after multiple attempts. Please try refreshing the page.';
		return false;
	}

	async function handleRefreshClick() {
		refreshing = true;
		const success = await fetchDataWithRetry();
		refreshing = false;
		if (success) {
			showNotification('Data refreshed successfully', 'success');
		} else {
			showNotification('Failed to refresh data', 'error');
		}
	}

	function showNotification(message: string, type: 'success' | 'error') {
		notification = { message, type };
		setTimeout(() => {
			notification = null;
		}, 3000);
	}

	onMount(() => {
		fetchDataWithRetry();
	});

	function getMaxDuration(metrics: Metrics): number {
		return Math.max(
			...Object.values(metrics).map((endpoints) =>
				Math.max(...Object.values(endpoints).map((e) => e.Duration ?? 0))
			),
			0.1
		);
	}

	$: maxDuration = getMaxDuration($metrics);
</script>

<svelte:head>
	<title>Ekolod Dashboard</title>
</svelte:head>

<div class="mx-auto max-w-7xl space-y-4">
	<header class="mb-8">
		<div class="flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
			<h1 class="text-base-content text-3xl font-bold">Dashboard</h1>
			<div class="flex gap-2">
				<div class="btn-group">
					<button
						class="btn btn-sm {$viewMode === 'grid' ? 'btn-active' : ''}"
						on:click={() => ($viewMode = 'grid')}
						title="Grid view"
					>
						<Grid class="h-4 w-4" />
					</button>
					<button
						class="btn btn-sm {$viewMode === 'list' ? 'btn-active' : ''}"
						on:click={() => ($viewMode = 'list')}
						title="List view"
					>
						<List class="h-4 w-4" />
					</button>
					<button
						class="btn btn-sm {$viewMode === 'chart' ? 'btn-active' : ''}"
						on:click={() => ($viewMode = 'chart')}
						title="Chart view"
					>
						<BarChart class="h-4 w-4" />
					</button>
				</div>
				<button
					class="btn btn-primary btn-sm {refreshing ? 'loading' : ''}"
					on:click={handleRefreshClick}
					disabled={refreshing}
				>
					{#if !refreshing}
						<RefreshCw class="mr-2 h-4 w-4" />
					{/if}
					Refresh
				</button>
			</div>
		</div>
		<p class="text-base-content/70 mt-2 text-sm">
			Monitoring {Object.keys($metrics).length} services
		</p>
	</header>

	{#if notification}
		<div class="fixed bottom-4 right-4 z-50" transition:fly={{ y: 20, duration: 300 }}>
			<div class="alert alert-{notification.type} shadow-lg">
				<AlertCircle class="h-6 w-6" />
				<span>{notification.message}</span>
			</div>
		</div>
	{/if}

	{#if loading}
		<div class="flex h-64 items-center justify-center">
			<span class="loading loading-spinner loading-lg"></span>
		</div>
	{:else if error}
		<div class="alert alert-error">
			<AlertCircle class="h-6 w-6" />
			<span>{error}</span>
		</div>
	{:else if $viewMode === 'grid'}
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
			{#each Object.entries($metrics) as [service, endpoints]}
				<ServiceGridItem {service} {endpoints} timeSeriesData={timeSeriesData[service]} />
			{/each}
		</div>
	{:else if $viewMode === 'list'}
		<ServiceTable metrics={$metrics} />
	{:else if $viewMode === 'chart'}
		<div class="bg-base-100 border-base-300 rounded-lg border p-6 shadow-sm">
			<h2 class="mb-4 text-xl font-semibold">Response Times by Service</h2>
			<div class="overflow-x-auto">
				<svg width={Object.keys($metrics).length * 60 + 50} height="400" class="chart">
					{#each Object.entries($metrics) as [service, endpoints], i}
						{@const value = Object.values(endpoints)[0]?.Duration ?? 0}
						{@const x = i * 60 + 25}
						{@const y = 350 - (value / maxDuration) * 300}
						<g transform="translate({x}, 0)">
							<rect
								width="40"
								height={350 - y}
								{y}
								fill="hsl(var(--p) / 50%)"
								stroke="hsl(var(--p))"
							/>
							<text x="20" y="370" text-anchor="middle" font-size="12">{service}</text>
							<text x="20" y={y - 5} text-anchor="middle" font-size="12" fill="hsl(var(--p))"
								>{value.toFixed(3)}</text
							>
						</g>
					{/each}
					<line
						x1="20"
						y1="350"
						x2={Object.keys($metrics).length * 60 + 30}
						y2="350"
						stroke="currentColor"
					/>
					<text x="10" y="10" font-size="14" fill="currentColor">Response Time (s)</text>
				</svg>
			</div>
		</div>
	{/if}
</div>
