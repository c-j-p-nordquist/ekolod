<script lang="ts">
	import { onMount } from 'svelte';
	import {
		targets,
		metrics,
		fetchTargets,
		fetchMetrics,
		type Metrics
	} from '../lib/stores/targets';
	import { writable } from 'svelte/store';
	import { Clock, Grid, List, BarChart } from 'lucide-svelte';

	let loading = true;
	let viewMode = writable('grid');

	onMount(async () => {
		await fetchTargets();
		await fetchMetrics();
		loading = false;
	});

	function getStatusColor(success: boolean | undefined): string {
		if (success === undefined) return 'bg-gray-300';
		return success ? 'bg-green-500' : 'bg-red-500';
	}

	function formatResponseTime(duration: number | undefined): string {
		if (duration === undefined || typeof duration !== 'number' || isNaN(duration)) return 'N/A';
		return duration.toFixed(3);
	}

	$: getIconColor = (mode: string) => ($viewMode === mode ? 'text-blue-600' : 'text-gray-600');

	$: chartData = Object.entries($metrics).map(([service, endpoints]) => {
		const duration = Object.values(endpoints)[0]?.Duration ?? 0;
		return { label: service, value: duration };
	});

	$: maxValue = Math.max(...chartData.map((d) => d.value), 0.1);

	function getBarHeight(value: number): number {
		return (value / maxValue) * 300; // 300 is the maximum bar height
	}

	function getServiceStatus(endpoints: { [key: string]: { Success: boolean } }): boolean {
		return Object.values(endpoints).every((endpoint) => endpoint.Success);
	}

	function getServiceDuration(endpoints: { [key: string]: { Duration: number } }): number {
		return Object.values(endpoints)[0]?.Duration ?? 0;
	}

	// Debug function
	function logMetrics() {
		console.log('Current metrics:', $metrics);
		console.log('Chart data:', chartData);
	}
</script>

<main class="container mx-auto p-4">
	<div class="mb-6 flex items-center justify-between">
		<h1 class="text-3xl font-bold text-gray-800">Dashboard</h1>
		<div class="flex space-x-2">
			<button
				class="rounded-md p-2 transition-colors duration-200 hover:bg-gray-200"
				on:click={() => ($viewMode = 'grid')}
			>
				<Grid class="h-5 w-5 {getIconColor('grid')}" />
			</button>
			<button
				class="rounded-md p-2 transition-colors duration-200 hover:bg-gray-200"
				on:click={() => ($viewMode = 'list')}
			>
				<List class="h-5 w-5 {getIconColor('list')}" />
			</button>
			<button
				class="rounded-md p-2 transition-colors duration-200 hover:bg-gray-200"
				on:click={() => {
					$viewMode = 'chart';
					logMetrics();
				}}
			>
				<BarChart class="h-5 w-5 {getIconColor('chart')}" />
			</button>
		</div>
	</div>

	{#if loading}
		<div class="flex h-64 items-center justify-center">
			<div class="h-16 w-16 animate-spin rounded-full border-b-2 border-t-2 border-blue-600"></div>
		</div>
	{:else if $viewMode === 'grid'}
		<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
			{#each Object.entries($metrics) as [service, endpoints]}
				<div
					class="rounded-lg bg-white p-6 shadow-md transition-shadow duration-300 hover:shadow-lg"
				>
					<h2 class="mb-4 text-xl font-semibold text-gray-800">{service}</h2>
					<div class="mb-2 flex items-center">
						<div
							class={`mr-2 h-3 w-3 rounded-full ${getStatusColor(getServiceStatus(endpoints))}`}
						></div>
						<p class="text-sm font-medium">
							Status: {#if getServiceStatus(endpoints)}
								<span class="text-green-600">Up</span>
							{:else}
								<span class="text-red-600">Down</span>
							{/if}
						</p>
					</div>
					<div class="flex items-center">
						<Clock class="mr-2 h-4 w-4 text-gray-600" />
						<p class="text-sm">
							Response Time: {formatResponseTime(getServiceDuration(endpoints))} seconds
						</p>
					</div>
				</div>
			{/each}
		</div>
	{:else if $viewMode === 'list'}
		<div class="overflow-hidden rounded-lg bg-white shadow-md">
			<table class="min-w-full divide-y divide-gray-200">
				<thead class="bg-gray-50">
					<tr>
						<th
							class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
							>Service</th
						>
						<th
							class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
							>Status</th
						>
						<th
							class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
							>Response Time</th
						>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-200 bg-white">
					{#each Object.entries($metrics) as [service, endpoints]}
						<tr>
							<td class="whitespace-nowrap px-6 py-4 text-sm font-medium text-gray-900"
								>{service}</td
							>
							<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
								<div class="flex items-center">
									<div
										class={`mr-2 h-3 w-3 rounded-full ${getStatusColor(getServiceStatus(endpoints))}`}
									></div>
									{#if getServiceStatus(endpoints)}
										<span class="text-green-600">Up</span>
									{:else}
										<span class="text-red-600">Down</span>
									{/if}
								</div>
							</td>
							<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500"
								>{formatResponseTime(getServiceDuration(endpoints))} seconds</td
							>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{:else if $viewMode === 'chart'}
		<div class="rounded-lg bg-white p-6 shadow-md">
			<h2 class="mb-4 text-xl font-semibold text-gray-800">Response Times by Service</h2>
			<div class="overflow-x-auto">
				<svg width={chartData.length * 60 + 50} height="400" class="chart">
					{#each chartData as { label, value }, i}
						<g transform="translate({i * 60 + 25}, 0)">
							<rect
								width="40"
								height={getBarHeight(value)}
								y={350 - getBarHeight(value)}
								fill="rgba(59, 130, 246, 0.5)"
								stroke="rgb(59, 130, 246)"
							/>
							<text x="20" y="370" text-anchor="middle" font-size="12">{label}</text>
							<text
								x="20"
								y={345 - getBarHeight(value)}
								text-anchor="middle"
								font-size="12"
								fill="rgb(59, 130, 246)">{formatResponseTime(value)}</text
							>
						</g>
					{/each}
					<line x1="20" y1="350" x2={chartData.length * 60 + 30} y2="350" stroke="black" />
					<text x="10" y="10" font-size="12">Response Time (s)</text>
				</svg>
			</div>
		</div>
	{/if}
</main>
