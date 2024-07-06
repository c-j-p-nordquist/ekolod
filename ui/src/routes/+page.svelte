<script>
	import IconServices from '~icons/mdi/application';
	import IconAlerts from '~icons/mdi/alert';
	import IconCheck from '~icons/mdi/check';
	import IconRobot from '~icons/mdi/robot';
	import SecondaryNav from '$lib/components/SecondaryNav.svelte';

	let overallHealth = $state(98);
	let totalServices = $state(25);
	let activeAlerts = $state(3);
	let uptime = $state(99.99);

	let recentAlerts = $state([
		{
			service: 'Web API',
			message: 'High latency detected',
			severity: 'warning',
			time: '5 minutes ago'
		},
		{
			service: 'Database',
			message: 'Connection timeout',
			severity: 'error',
			time: '15 minutes ago'
		},
		{
			service: 'Auth Service',
			message: 'Increased error rate',
			severity: 'warning',
			time: '1 hour ago'
		}
	]);

	let topServices = $state([
		{ name: 'Web API', status: 'Healthy', responseTime: '120ms' },
		{ name: 'Database', status: 'Degraded', responseTime: '350ms' },
		{ name: 'Auth Service', status: 'Healthy', responseTime: '80ms' },
		{ name: 'Payment Gateway', status: 'Healthy', responseTime: '200ms' },
		{ name: 'Notification Service', status: 'Healthy', responseTime: '90ms' }
	]);

	let aiInsights = $state([
		"The Web API's latency has increased by 15% over the past 24 hours. Consider investigating potential bottlenecks.",
		'Database connections are showing intermittent timeouts. Recommend scaling up the database instance.',
		'Auth Service error rate spike correlates with a recent deployment. Suggest rolling back to the previous version.',
		'Probe in Asia-Pacific region is reporting higher latencies. Consider adding more probes in this area for better coverage.',
		'Overall system health has improved by 2% since implementing the last set of recommended optimizations.'
	]);

	let sections = [
		{ id: 'overview', title: 'Overview' },
		{ id: 'recent-alerts', title: 'Recent Alerts' },
		{ id: 'top-services', title: 'Top Services' },
		{ id: 'ai-insights', title: 'AI Insights' }
	];

	function getStatusColor(status) {
		switch (status.toLowerCase()) {
			case 'healthy':
				return 'text-success';
			case 'degraded':
				return 'text-warning';
			case 'error':
				return 'text-error';
			default:
				return 'text-info';
		}
	}

	function getSeverityColor(severity) {
		switch (severity.toLowerCase()) {
			case 'error':
				return 'badge-error';
			case 'warning':
				return 'badge-warning';
			default:
				return 'badge-info';
		}
	}
</script>

<div class="flex flex-col xl:flex-row">
	<div class="flex-grow space-y-6 xl:pr-4">
		<h1 class="text-3xl font-bold">Dashboard</h1>

		<div id="overview" class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4">
			<div class="stat bg-base-100 shadow">
				<div class="stat-figure text-primary">
					<IconServices class="w-8 h-8 sm:w-10 sm:h-10" />
				</div>
				<div class="stat-title text-xs sm:text-sm">Overall Health</div>
				<div class="stat-value text-primary text-2xl sm:text-3xl">{overallHealth}%</div>
				<div class="stat-desc text-xs">{totalServices} Services Monitored</div>
			</div>

			<div class="stat bg-base-100 shadow">
				<div class="stat-figure text-secondary">
					<IconAlerts class="w-8 h-8 sm:w-10 sm:h-10" />
				</div>
				<div class="stat-title text-xs sm:text-sm">Active Alerts</div>
				<div class="stat-value text-secondary text-2xl sm:text-3xl">{activeAlerts}</div>
				<div class="stat-desc text-xs">Across all services</div>
			</div>

			<div class="stat bg-base-100 shadow">
				<div class="stat-figure text-accent">
					<IconCheck class="w-8 h-8 sm:w-10 sm:h-10" />
				</div>
				<div class="stat-title text-xs sm:text-sm">Uptime</div>
				<div class="stat-value text-2xl sm:text-3xl">{uptime}%</div>
				<div class="stat-desc text-xs">Last 30 days</div>
			</div>
		</div>

		<div class="grid grid-cols-1 xl:grid-cols-2 gap-6">
			<div id="recent-alerts" class="card bg-base-100 shadow-xl">
				<div class="card-body p-4">
					<h2 class="card-title text-lg">Recent Alerts</h2>
					<div class="overflow-x-auto">
						<table class="table table-compact w-full">
							<thead>
								<tr>
									<th class="text-xs">Service</th>
									<th class="text-xs">Message</th>
									<th class="text-xs">Severity</th>
									<th class="text-xs">Time</th>
								</tr>
							</thead>
							<tbody>
								{#each recentAlerts as alert}
									<tr>
										<td class="text-xs whitespace-normal">{alert.service}</td>
										<td class="text-xs whitespace-normal break-words max-w-[150px]"
											>{alert.message}</td
										>
										<td>
											<div class="badge {getSeverityColor(alert.severity)} text-xs">
												{alert.severity}
											</div>
										</td>
										<td class="text-xs whitespace-nowrap">{alert.time}</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				</div>
			</div>

			<div id="top-services" class="card bg-base-100 shadow-xl">
				<div class="card-body p-4">
					<h2 class="card-title text-lg">Top Services</h2>
					<div class="overflow-x-auto">
						<table class="table table-compact w-full">
							<thead>
								<tr>
									<th class="text-xs">Service</th>
									<th class="text-xs">Status</th>
									<th class="text-xs">Response Time</th>
								</tr>
							</thead>
							<tbody>
								{#each topServices as service}
									<tr>
										<td class="text-xs whitespace-normal">{service.name}</td>
										<td>
											<div class="badge {getStatusColor(service.status)} text-xs">
												{service.status}
											</div>
										</td>
										<td class="text-xs whitespace-nowrap">{service.responseTime}</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				</div>
			</div>
		</div>

		<div id="ai-insights" class="card bg-base-100 shadow-xl">
			<div class="card-body p-4">
				<h2 class="card-title text-lg flex items-center">
					<IconRobot class="w-6 h-6 mr-2" />
					AI Insights
				</h2>
				<ul class="list-disc pl-5 space-y-2">
					{#each aiInsights as insight}
						<li class="text-xs">{insight}</li>
					{/each}
				</ul>
			</div>
		</div>
	</div>

	<div class="hidden xl:block xl:w-64">
		<SecondaryNav {sections} />
	</div>
</div>
