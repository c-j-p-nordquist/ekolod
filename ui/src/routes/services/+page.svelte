<script>
	import IconSearch from '~icons/mdi/magnify';
	import IconPlus from '~icons/mdi/plus';
	import IconEdit from '~icons/mdi/pencil';
	import IconDelete from '~icons/mdi/delete';

	let services = $state([
		{ id: 1, name: 'Web API', status: 'Healthy', uptime: '99.99%', responseTime: '120ms' },
		{ id: 2, name: 'Database', status: 'Degraded', uptime: '99.95%', responseTime: '350ms' },
		{ id: 3, name: 'Auth Service', status: 'Healthy', uptime: '99.98%', responseTime: '80ms' },
		{ id: 4, name: 'Payment Gateway', status: 'Healthy', uptime: '99.99%', responseTime: '200ms' },
		{
			id: 5,
			name: 'Notification Service',
			status: 'Healthy',
			uptime: '99.97%',
			responseTime: '90ms'
		},
		{ id: 6, name: 'User Management', status: 'Healthy', uptime: '99.99%', responseTime: '110ms' },
		{ id: 7, name: 'Logging Service', status: 'Degraded', uptime: '99.90%', responseTime: '180ms' },
		{ id: 8, name: 'CDN', status: 'Healthy', uptime: '99.99%', responseTime: '30ms' }
	]);

	let searchQuery = $state('');

	let filteredServices = $derived(
		services.filter((service) => service.name.toLowerCase().includes(searchQuery.toLowerCase()))
	);

	function getStatusColor(status) {
		switch (status.toLowerCase()) {
			case 'healthy':
				return 'text-success';
			case 'degraded':
				return 'text-warning';
			case 'down':
				return 'text-error';
			default:
				return 'text-info';
		}
	}

	function editService(id) {
		// Placeholder for edit functionality
		console.log('Edit service:', id);
	}

	function deleteService(id) {
		// Placeholder for delete functionality
		console.log('Delete service:', id);
	}
</script>

<div class="p-4 space-y-6">
	<h1 class="text-3xl font-bold">Services</h1>

	<div class="flex justify-between items-center">
		<div class="form-control">
			<div class="input-group">
				<input
					type="text"
					placeholder="Search services..."
					class="input input-bordered"
					bind:value={searchQuery}
				/>
				<button class="btn btn-square">
					<IconSearch />
				</button>
			</div>
		</div>
		<button class="btn btn-primary">
			<IconPlus />
			Add New Service
		</button>
	</div>

	<div class="overflow-x-auto">
		<table class="table w-full">
			<thead>
				<tr>
					<th>Name</th>
					<th>Status</th>
					<th>Uptime</th>
					<th>Response Time</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredServices as service}
					<tr>
						<td>
							<a href="/services/{service.id}" class="link link-hover">
								{service.name}
							</a>
						</td>
						<td>
							<div class="badge {getStatusColor(service.status)}">
								{service.status}
							</div>
						</td>
						<td>{service.uptime}</td>
						<td>{service.responseTime}</td>
						<td>
							<button class="btn btn-ghost btn-xs" onclick={() => editService(service.id)}>
								<IconEdit />
							</button>
							<button class="btn btn-ghost btn-xs" onclick={() => deleteService(service.id)}>
								<IconDelete />
							</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
