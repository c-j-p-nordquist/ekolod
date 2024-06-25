<script>
	import '../app.css';
</script>

<div class="bg-base-200 min-h-screen">
	<nav class="navbar bg-base-100">
		<div class="flex-1">
			<a href="/" class="btn btn-ghost text-xl normal-case">Ekolod</a>
		</div>
		<div class="flex-none">
			<ul class="menu menu-horizontal px-1">
				<li><a href="/">Dashboard</a></li>
				<li><a href="/config">Configuration</a></li>
			</ul>
		</div>
	</nav>

	<main class="container mx-auto p-4">
		<slot />
	</main>
</div>
