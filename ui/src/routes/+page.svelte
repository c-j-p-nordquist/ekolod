<script>
	import IconServices from '~icons/mdi/application';
	import IconAlerts from '~icons/mdi/alert';
	import IconCheck from '~icons/mdi/check-circle';

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
</script>

<div class="p-4 space-y-6">
	<h1 class="text-3xl font-bold">Dashboard</h1>

	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
		<div class="stat bg-base-100 shadow">
			<div class="stat-figure text-primary">
				<IconServices class="w-8 h-8" />
			</div>
			<div class="stat-title">Overall Health</div>
			<div class="stat-value text-primary">{overallHealth}%</div>
			<div class="stat-desc">{totalServices} Services Monitored</div>
		</div>

		<div class="stat bg-base-100 shadow">
			<div class="stat-figure text-secondary">
				<IconAlerts class="w-8 h-8" />
			</div>
			<div class="stat-title">Active Alerts</div>
			<div class="stat-value text-secondary">{activeAlerts}</div>
			<div class="stat-desc">Across all services</div>
		</div>

		<div class="stat bg-base-100 shadow">
			<div class="stat-figure text-accent">
				<IconCheck class="w-8 h-8" />
			</div>
			<div class="stat-title">Uptime</div>
			<div class="stat-value">{uptime}%</div>
			<div class="stat-desc">Last 30 days</div>
		</div>
	</div>

	<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<h2 class="card-title">Recent Alerts</h2>
				<div class="overflow-x-auto">
					<table class="table w-full">
						<thead>
							<tr>
								<th>Service</th>
								<th>Message</th>
								<th>Severity</th>
								<th>Time</th>
							</tr>
						</thead>
						<tbody>
							{#each recentAlerts as alert}
								<tr>
									<td>{alert.service}</td>
									<td>{alert.message}</td>
									<td>
										<div class="badge badge-{alert.severity === 'error' ? 'error' : 'warning'}">
											{alert.severity}
										</div>
									</td>
									<td>{alert.time}</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>
		</div>

		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<h2 class="card-title">Top Services</h2>
				<div class="overflow-x-auto">
					<table class="table w-full">
						<thead>
							<tr>
								<th>Service</th>
								<th>Status</th>
								<th>Response Time</th>
							</tr>
						</thead>
						<tbody>
							{#each topServices as service}
								<tr>
									<td>{service.name}</td>
									<td>
										<div class="badge badge-{service.status === 'Healthy' ? 'success' : 'warning'}">
											{service.status}
										</div>
									</td>
									<td>{service.responseTime}</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>
		</div>
	</div>
</div>
