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

<div class="w-64 p-4 bg-base-100 border-l border-base-300">
	<div class="sticky top-20">
		<h3 class="text-lg font-semibold mb-4">On this page</h3>
		<ul class="space-y-2">
			{#each sections as section}
				<li>
					<a
						href="#{section.id}"
						class="text-sm text-base-content hover:text-primary transition-colors duration-200"
						onclick={(event) => handleClick(event, section.id)}
					>
						{section.title}
					</a>
				</li>
			{/each}
		</ul>
	</div>
</div>
