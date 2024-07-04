<script>
	import IconDashboard from '~icons/mdi/view-dashboard';
	import IconServices from '~icons/mdi/application';
	import IconMonitoring from '~icons/mdi/monitor';
	import IconAlerts from '~icons/mdi/alert';
	import IconReports from '~icons/mdi/chart-box';
	import IconSettings from '~icons/mdi/cog';

	let menuItems = $state([
		{ icon: IconDashboard, label: 'Dashboard', href: '/' },
		{
			icon: IconServices,
			label: 'Services',
			submenu: [
				{ label: 'All Services', href: '/services' },
				{ label: 'Add New Service', href: '/services/new' }
			]
		},
		{
			icon: IconMonitoring,
			label: 'Monitoring',
			submenu: [
				{ label: 'Checks', href: '/monitoring/checks' },
				{ label: 'Probes', href: '/monitoring/probes' },
				{ label: 'Integrations', href: '/monitoring/integrations' }
			]
		},
		{
			icon: IconAlerts,
			label: 'Alerts',
			submenu: [
				{ label: 'Active Alerts', href: '/alerts' },
				{ label: 'Alert History', href: '/alerts/history' },
				{ label: 'Alert Rules', href: '/alerts/rules' }
			]
		},
		{
			icon: IconReports,
			label: 'Reports & Analytics',
			submenu: [
				{ label: 'Performance Reports', href: '/reports' },
				{ label: 'Custom Dashboards', href: '/reports/custom' },
				{ label: 'Data Export', href: '/reports/export' }
			]
		},
		{
			icon: IconSettings,
			label: 'System Config',
			submenu: [
				{ label: 'General Settings', href: '/config' },
				{ label: 'User Management', href: '/config/users' },
				{ label: 'API Configuration', href: '/config/api' },
				{ label: 'Notification Settings', href: '/config/notifications' }
			]
		}
	]);
</script>

<div class="bg-base-200 text-base-content w-64 min-h-full border-r border-base-300">
	<div class="p-4">
		<h1 class="text-2xl font-bold">Ekolod</h1>
	</div>
	<ul class="menu p-4 w-full">
		{#each menuItems as item}
			<li>
				{#if item.submenu}
					<details class="w-full">
						<summary class="flex items-center justify-between px-4 py-2 w-full">
							<div class="flex items-center">
								<svelte:component this={item.icon} class="w-5 h-5 mr-2" />
								{item.label}
							</div>
						</summary>
						<ul class="pl-4 mt-2">
							{#each item.submenu as subitem}
								<li><a href={subitem.href} class="py-2">{subitem.label}</a></li>
							{/each}
						</ul>
					</details>
				{:else}
					<a href={item.href} class="flex items-center justify-between px-4 py-2 w-full">
						<div class="flex items-center">
							<svelte:component this={item.icon} class="w-5 h-5 mr-2" />
							{item.label}
						</div>
					</a>
				{/if}
			</li>
		{/each}
	</ul>
</div>
