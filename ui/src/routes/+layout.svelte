<script lang="ts">
	import '../app.css';
</script>

<div class="flex min-h-screen flex-col bg-gray-100">
	<nav class="bg-white shadow-sm">
		<div class="container mx-auto px-4">
			<div class="flex h-16 items-center justify-between">
				<a
					href="/"
					class="text-2xl font-extrabold tracking-wider text-gray-800 transition-colors duration-200 hover:text-blue-600"
					>EKOLOD</a
				>
				<div class="hidden md:block">
					<div class="ml-10 flex items-baseline space-x-4">
						<a
							href="/"
							class="rounded-md px-3 py-2 text-sm font-medium text-gray-800 transition-colors duration-200 hover:text-blue-600"
							>Dashboard</a
						>
						<!-- Add more navigation items here as needed -->
					</div>
				</div>
			</div>
		</div>
	</nav>

	<main class="container mx-auto flex-grow px-4 py-8">
		<slot />
	</main>

	<footer class="mt-auto bg-white shadow-sm">
		<div class="container mx-auto px-4 py-4 text-center text-gray-600">
			© {new Date().getFullYear()} Ekolod. All rights reserved.
		</div>
	</footer>
</div>
