<script lang="ts">
	import { onMount } from 'svelte';
	import { targets, fetchTargets, fetchMetrics } from '../lib/stores/targets';
	import { writable } from 'svelte/store';

	const metrics = writable<{ [key: string]: number }>({});

	onMount(async () => {
		await fetchTargets();
		const metricsData = await fetchMetrics();
		metrics.set(metricsData);
	});
</script>

<main class="container mx-auto p-4">
	<h1 class="mb-4 text-4xl font-bold">Dashboard</h1>
	<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
		{#each $targets as target}
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<h2 class="card-title">{target}</h2>
					<p>
						Status: {#if $metrics[target]}<span class="text-green-500">Up</span>{:else}<span
								class="text-red-500">Down</span
							>{/if}
					</p>
					<p>Response Time: {$metrics[target] ? $metrics[target].toFixed(2) : 'N/A'} seconds</p>
				</div>
			</div>
		{/each}
	</div>
</main>
