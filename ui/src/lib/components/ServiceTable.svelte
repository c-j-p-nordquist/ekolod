<script lang="ts">
	import type { Metrics, MetricData } from '$lib/stores/targets';

	export let metrics: Metrics;

	function formatResponseTime(duration: number | undefined): string {
		if (duration === undefined || typeof duration !== 'number' || isNaN(duration)) return 'N/A';
		return duration.toFixed(3);
	}

	function formatContentLength(length: number | undefined): string {
		if (length === undefined || typeof length !== 'number' || isNaN(length)) return 'N/A';
		if (length < 1024) return `${length} B`;
		if (length < 1048576) return `${(length / 1024).toFixed(2)} KB`;
		return `${(length / 1048576).toFixed(2)} MB`;
	}

	function getServiceStatus(endpoints: { [key: string]: MetricData }): boolean {
		return Object.values(endpoints).every((endpoint) => endpoint.Success);
	}

	function getServiceDuration(endpoints: { [key: string]: MetricData }): number {
		return Object.values(endpoints)[0]?.Duration ?? 0;
	}

	function getCertExpiryStatus(days: number): 'success' | 'warning' | 'error' {
		if (days > 30) return 'success';
		if (days > 7) return 'warning';
		return 'error';
	}

	function getStatusColor(status: boolean): string {
		return status ? 'bg-success' : 'bg-error';
	}
</script>

<div class="bg-base-100 overflow-x-auto rounded-lg shadow-sm">
	<table class="table w-full">
		<thead>
			<tr>
				<th class="bg-base-200">Service</th>
				<th class="bg-base-200">Status</th>
				<th class="bg-base-200">Response Time</th>
				<th class="bg-base-200">Content Length</th>
				<th class="bg-base-200">TLS Version</th>
				<th class="bg-base-200">Cert Expiry</th>
			</tr>
		</thead>
		<tbody>
			{#each Object.entries(metrics) as [service, endpoints]}
				<tr class="hover:bg-base-100 border-base-300 border-b">
					<td class="font-medium">{service}</td>
					<td>
						<div class="flex items-center">
							<div
								class={`mr-2 h-3 w-3 rounded-full ${getStatusColor(getServiceStatus(endpoints))}`}
							></div>
							<span class="text-sm font-medium">
								{getServiceStatus(endpoints) ? 'Up' : 'Down'}
							</span>
						</div>
					</td>
					<td>{formatResponseTime(getServiceDuration(endpoints))}s</td>
					<td>{formatContentLength(Object.values(endpoints)[0]?.ContentLength)}</td>
					<td>{Object.values(endpoints)[0]?.TLSVersion || 'N/A'}</td>
					<td>
						<span
							class="text-{getCertExpiryStatus(Object.values(endpoints)[0]?.CertExpiryDays || 0)}"
						>
							{Object.values(endpoints)[0]?.CertExpiryDays || 'N/A'} days
						</span>
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
