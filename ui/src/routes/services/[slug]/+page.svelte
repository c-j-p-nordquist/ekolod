<script>
	import { page } from '$app/stores';
	import IconRefresh from '~icons/mdi/refresh';
	import IconEdit from '~icons/mdi/pencil';
	import IconDelete from '~icons/mdi/delete';
	import IconMapMarker from '~icons/mdi/map-marker';
	import { Chart, registerables } from 'chart.js';
	import { onMount } from 'svelte';

	Chart.register(...registerables);

	let service = $state({
		id: $page.params.slug,
		name: 'Web API',
		status: 'Healthy',
		uptime: '99.99%',
		responseTime: '120ms',
		description: 'Main API service for handling client requests.',
		lastChecked: '2 minutes ago',
		checks: [
			{ id: 1, name: 'HTTP Check', status: 'Passed', lastRun: '1 minute ago' },
			{ id: 2, name: 'SSL Certificate', status: 'Passed', lastRun: '5 minutes ago' },
			{ id: 3, name: 'Database Connection', status: 'Passed', lastRun: '2 minutes ago' }
		],
		probes: [
			{
				id: 1,
				location: 'New York',
				status: 'Active',
				latency: '50ms',
				lat: 40.7128,
				lon: -74.006
			},
			{ id: 2, location: 'London', status: 'Active', latency: '100ms', lat: 51.5074, lon: -0.1278 },
			{ id: 3, location: 'Tokyo', status: 'Active', latency: '200ms', lat: 35.6762, lon: 139.6503 }
		]
	});

	let activeTab = $state('overview');
	let responseTimeChart;
	let map;

	onMount(() => {
		// Response Time Chart
		const ctx = document.getElementById('responseTimeChart').getContext('2d');
		responseTimeChart = new Chart(ctx, {
			type: 'line',
			data: {
				labels: ['1h ago', '45m ago', '30m ago', '15m ago', 'Now'],
				datasets: [
					{
						label: 'Response Time (ms)',
						data: [100, 120, 115, 130, 120],
						borderColor: 'rgb(75, 192, 192)',
						tension: 0.1
					}
				]
			},
			options: {
				responsive: true,
				scales: {
					y: {
						beginAtZero: true
					}
				}
			}
		});

		// Initialize map
		map = L.map('probeMap').setView([0, 0], 2);
		L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
			attribution: '© OpenStreetMap contributors'
		}).addTo(map);

		service.probes.forEach((probe) => {
			L.marker([probe.lat, probe.lon])
				.addTo(map)
				.bindPopup(`${probe.location}<br>Latency: ${probe.latency}`);
		});
	});

	function refreshService() {
		console.log('Refreshing service data');
	}

	function getStatusColor(status) {
		return status.toLowerCase() === 'healthy' ? 'text-success' : 'text-error';
	}

	function openEditModal() {
		document.getElementById('edit_modal').showModal();
	}

	function openDeleteModal() {
		document.getElementById('delete_modal').showModal();
	}

	function closeModals() {
		document.getElementById('edit_modal').close();
		document.getElementById('delete_modal').close();
	}
</script>

<svelte:head>
	<link rel="stylesheet" href="https://unpkg.com/leaflet@1.7.1/dist/leaflet.css" />
	<script src="https://unpkg.com/leaflet@1.7.1/dist/leaflet.js"></script>
</svelte:head>

<div class="p-4 space-y-6">
	<div class="text-sm breadcrumbs">
		<ul>
			<li><a href="/services">Services</a></li>
			<li>{service.name}</li>
		</ul>
	</div>

	<div class="flex justify-between items-center">
		<h1 class="text-3xl font-bold">{service.name}</h1>
		<div class="space-x-2">
			<button class="btn btn-primary" onclick={openEditModal}>
				<IconEdit class="w-5 h-5 mr-2" /> Edit
			</button>
			<button class="btn btn-error" onclick={openDeleteModal}>
				<IconDelete class="w-5 h-5 mr-2" /> Delete
			</button>
			<button class="btn btn-ghost" onclick={refreshService}>
				<IconRefresh class="w-5 h-5 mr-2" /> Refresh
			</button>
		</div>
	</div>

	<div class="tabs tabs-boxed bg-base-200 p-2 rounded-lg">
		<button
			class="tab {activeTab === 'overview' ? 'tab-active bg-primary text-primary-content' : ''}"
			onclick={() => (activeTab = 'overview')}
			onkeydown={(e) => e.key === 'Enter' && (activeTab = 'overview')}
		>
			Overview
		</button>
		<button
			class="tab {activeTab === 'checks' ? 'tab-active bg-primary text-primary-content' : ''}"
			onclick={() => (activeTab = 'checks')}
			onkeydown={(e) => e.key === 'Enter' && (activeTab = 'checks')}
		>
			Checks
		</button>
		<button
			class="tab {activeTab === 'probes' ? 'tab-active bg-primary text-primary-content' : ''}"
			onclick={() => (activeTab = 'probes')}
			onkeydown={(e) => e.key === 'Enter' && (activeTab = 'probes')}
		>
			Probes
		</button>
	</div>

	{#if activeTab === 'overview'}
		<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
			<div class="stat bg-base-100 shadow">
				<div class="stat-title">Status</div>
				<div class="stat-value {getStatusColor(service.status)}">{service.status}</div>
			</div>
			<div class="stat bg-base-100 shadow">
				<div class="stat-title">Uptime</div>
				<div class="stat-value">{service.uptime}</div>
			</div>
			<div class="stat bg-base-100 shadow">
				<div class="stat-title">Response Time</div>
				<div class="stat-value">{service.responseTime}</div>
			</div>
			<div class="stat bg-base-100 shadow">
				<div class="stat-title">Last Checked</div>
				<div class="stat-value">{service.lastChecked}</div>
			</div>
		</div>
		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<h2 class="card-title">Description</h2>
				<p>{service.description}</p>
			</div>
		</div>
		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<h2 class="card-title">Response Time Trend</h2>
				<canvas id="responseTimeChart"></canvas>
			</div>
		</div>
		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<h2 class="card-title">Probe Locations</h2>
				<div id="probeMap" style="height: 400px;"></div>
			</div>
		</div>
	{:else if activeTab === 'checks'}
		<div class="overflow-x-auto">
			<table class="table w-full">
				<thead>
					<tr>
						<th>Check Name</th>
						<th>Status</th>
						<th>Last Run</th>
					</tr>
				</thead>
				<tbody>
					{#each service.checks as check}
						<tr>
							<td>{check.name}</td>
							<td>
								<div class="badge {check.status === 'Passed' ? 'badge-success' : 'badge-error'}">
									{check.status}
								</div>
							</td>
							<td>{check.lastRun}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{:else if activeTab === 'probes'}
		<div class="overflow-x-auto">
			<table class="table w-full">
				<thead>
					<tr>
						<th>Location</th>
						<th>Status</th>
						<th>Latency</th>
					</tr>
				</thead>
				<tbody>
					{#each service.probes as probe}
						<tr>
							<td>
								<div class="flex items-center">
									<IconMapMarker class="w-5 h-5 mr-2" />
									{probe.location}
								</div>
							</td>
							<td>{probe.status}</td>
							<td>{probe.latency}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

<dialog id="edit_modal" class="modal modal-bottom sm:modal-middle">
	<div class="modal-box">
		<h3 class="font-bold text-lg">Edit Service</h3>
		<p class="py-4">Edit service form would go here.</p>
		<div class="modal-action">
			<form method="dialog">
				<button class="btn" onclick={closeModals}>Close</button>
			</form>
		</div>
	</div>
</dialog>

<dialog id="delete_modal" class="modal modal-bottom sm:modal-middle">
	<div class="modal-box">
		<h3 class="font-bold text-lg">Delete Service</h3>
		<p class="py-4">Are you sure you want to delete this service? This action cannot be undone.</p>
		<div class="modal-action">
			<form method="dialog">
				<button class="btn btn-error" onclick={closeModals}>Delete</button>
				<button class="btn" onclick={closeModals}>Cancel</button>
			</form>
		</div>
	</div>
</dialog>
