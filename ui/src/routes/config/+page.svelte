<script lang="ts">
	import { writable, derived } from 'svelte/store';

	const config = writable({
		targets: [
			{ name: 'Example', url: 'https://example.com', interval: '10s', timeout: '5s' },
			{ name: 'Google', url: 'https://www.google.com', interval: '15s', timeout: '5s' }
		]
	});

	const totalTargets = derived(config, ($config) => $config.targets.length);

	function addTarget() {
		config.update((c) => ({
			...c,
			targets: [
				...c.targets,
				{ name: 'New Target', url: 'https://example.org', interval: '30s', timeout: '10s' }
			]
		}));
	}
</script>

<h1 class="mb-4 text-2xl font-bold">Configuration</h1>

<p class="mb-4">Total Targets: {$totalTargets}</p>

<button class="btn btn-primary mb-4" on:click={addTarget}>Add Target</button>

<div class="overflow-x-auto">
	<table class="table w-full">
		<thead>
			<tr>
				<th>Name</th>
				<th>URL</th>
				<th>Interval</th>
				<th>Timeout</th>
			</tr>
		</thead>
		<tbody>
			{#each $config.targets as target}
				<tr>
					<td>{target.name}</td>
					<td>{target.url}</td>
					<td>{target.interval}</td>
					<td>{target.timeout}</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
