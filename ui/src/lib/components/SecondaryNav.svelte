<script>
	const { sections } = $props();

	function scrollToSection(id) {
		const element = document.getElementById(id);
		if (element) {
			const yOffset = -70; // Adjust this value based on your TopBar height
			const y = element.getBoundingClientRect().top + window.pageYOffset + yOffset;
			window.scrollTo({ top: y, behavior: 'smooth' });
		}
	}

	function handleClick(event, id) {
		event.preventDefault();
		scrollToSection(id);
	}
</script>

<div class="hidden lg:block w-64 p-4 border-l bg-base-100">
	<div class="fixed top-16 w-64">
		<h3 class="text-lg font-semibold mb-4">On this page</h3>
		<ul class="space-y-2">
			{#each sections as section}
				<li>
					<a
						href="#{section.id}"
						class="text-base-content hover:text-primary transition-colors duration-200"
						onclick={(event) => handleClick(event, section.id)}
					>
						{section.title}
					</a>
				</li>
			{/each}
		</ul>
	</div>
</div>
