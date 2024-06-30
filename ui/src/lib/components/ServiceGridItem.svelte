<script lang="ts">
	import type { MetricData } from '$lib/stores/targets';

	export let service: string;
	export let endpoints: { [key: string]: MetricData };
	export let timeSeriesData: number[];

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

	$: maxValue = Math.max(...timeSeriesData, 1);
</script>

<div
	class="bg-base-100 border-base-300 rounded-lg border p-4 shadow-sm transition-shadow duration-200 hover:shadow-md"
>
	<div class="mb-2 flex items-center justify-between">
		<h2 class="text-lg font-semibold">{service}</h2>
		<div class="flex items-center">
			<div class={`mr-2 h-3 w-3 rounded-full ${getStatusColor(getServiceStatus(endpoints))}`}></div>
			<span class="text-sm font-medium">
				{getServiceStatus(endpoints) ? 'Up' : 'Down'}
			</span>
		</div>
	</div>
	<div class="space-y-2 text-sm">
		<p>
			<span class="font-medium">Response Time:</span>
			{formatResponseTime(getServiceDuration(endpoints))}s
		</p>
		<p>
			<span class="font-medium">Content Length:</span>
			{formatContentLength(Object.values(endpoints)[0]?.ContentLength)}
		</p>
		<p>
			<span class="font-medium">TLS Version:</span>
			{Object.values(endpoints)[0]?.TLSVersion || 'N/A'}
		</p>
		<p>
			<span class="font-medium">Cert Expiry:</span>
			<span class="text-{getCertExpiryStatus(Object.values(endpoints)[0]?.CertExpiryDays || 0)}">
				{Object.values(endpoints)[0]?.CertExpiryDays || 'N/A'} days
			</span>
		</p>
	</div>
	<div class="mt-4 flex h-20 items-end">
		{#each timeSeriesData as value, i}
			<div
				class="bg-primary w-full"
				style="height: {(value / maxValue) * 100}%; margin-right: 2px;"
			></div>
		{/each}
	</div>
</div>
