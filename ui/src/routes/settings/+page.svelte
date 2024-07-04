<script>
	import { onMount } from 'svelte';
	import IconPalette from '~icons/mdi/palette';
	import IconCheck from '~icons/mdi/check';

	let currentTheme = $state('light'); // Default theme
	let isThemeDropdownOpen = $state(false);

	const themes = [
		{ name: 'Light', value: 'light' },
		{ name: 'Dark', value: 'dark' },
		{ name: 'Retro', value: 'retro' },
		{ name: 'Cyberpunk', value: 'cyberpunk' },
		{ name: 'Valentine', value: 'valentine' },
		{ name: 'Aqua', value: 'aqua' }
	];

	onMount(() => {
		// Access localStorage only after the component has mounted (client-side)
		currentTheme = localStorage.getItem('theme') || 'light';
	});

	function changeTheme(theme) {
		currentTheme = theme;
		document.documentElement.setAttribute('data-theme', theme);
		localStorage.setItem('theme', theme);
		isThemeDropdownOpen = false;
	}

	function toggleThemeDropdown() {
		isThemeDropdownOpen = !isThemeDropdownOpen;
	}
</script>

<div class="p-4 space-y-6">
	<h1 class="text-3xl font-bold">Settings</h1>

	<div class="card bg-base-100 shadow-xl">
		<div class="card-body">
			<h2 class="card-title">
				<IconPalette class="w-6 h-6 mr-2" />
				Theme Settings
			</h2>
			<div class="form-control">
				<label for="theme-select" class="label">
					<span class="label-text">Select Theme</span>
				</label>
				<div class="dropdown relative">
					<button
						id="theme-select"
						type="button"
						class="btn btn-outline w-full max-w-xs justify-between"
						onclick={toggleThemeDropdown}
						aria-haspopup="listbox"
						aria-expanded={isThemeDropdownOpen}
					>
						{themes.find((t) => t.value === currentTheme)?.name}
						<svg
							class="fill-current"
							xmlns="http://www.w3.org/2000/svg"
							width="24"
							height="24"
							viewBox="0 0 24 24"
							><path d="M7.41,8.58L12,13.17L16.59,8.58L18,10L12,16L6,10L7.41,8.58Z" /></svg
						>
					</button>
					{#if isThemeDropdownOpen}
						<ul
							class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-full max-w-xs absolute z-50 mt-1"
							role="listbox"
							aria-labelledby="theme-select"
						>
							{#each themes as theme}
								<li>
									<button
										type="button"
										class="flex justify-between"
										role="option"
										aria-selected={currentTheme === theme.value}
										onclick={() => changeTheme(theme.value)}
									>
										{theme.name}
										{#if currentTheme === theme.value}
											<IconCheck class="w-5 h-5" />
										{/if}
									</button>
								</li>
							{/each}
						</ul>
					{/if}
				</div>
			</div>
		</div>
	</div>

	<!-- Placeholder for other settings -->
	<div class="card bg-base-100 shadow-xl">
		<div class="card-body">
			<h2 class="card-title">Other Settings</h2>
			<p>More settings will be added here.</p>
		</div>
	</div>
</div>
