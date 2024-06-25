<script lang="ts">
	import { writable } from 'svelte/store';

	const probes = writable([
		{ name: 'Example', status: 'UP', latency: '50ms' },
		{ name: 'Google', status: 'UP', latency: '100ms' }
	]);

	function addProbe() {
		probes.update((p) => [...p, { name: 'New Probe', status: 'UNKNOWN', latency: '0ms' }]);
	}
</script>

<h1 class="mb-4 text-2xl font-bold">Ekolod Dashboard</h1>

<button class="btn btn-primary mb-4" on:click={addProbe}>Add Probe</button>

<div class="overflow-x-auto">
	<table class="table w-full">
		<thead>
			<tr>
				<th>Name</th>
				<th>Status</th>
				<th>Latency</th>
			</tr>
		</thead>
		<tbody>
			{#each $probes as probe}
				<tr>
					<td>{probe.name}</td>
					<td
						><span class="badge {probe.status === 'UP' ? 'badge-success' : 'badge-error'}"
							>{probe.status}</span
						></td
					>
					<td>{probe.latency}</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
