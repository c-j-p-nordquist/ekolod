<script lang="ts">
	import { fade } from 'svelte/transition';
	import { onMount } from 'svelte';

	export let message: string;
	export let type: 'success' | 'error' = 'success';
	export let duration: number = 3000;

	let visible = true;

	onMount(() => {
		setTimeout(() => {
			visible = false;
		}, duration);
	});
</script>

{#if visible}
	<div
		transition:fade={{ duration: 300 }}
		class="fixed right-4 top-4 rounded-lg p-4 shadow-lg {type === 'success'
			? 'bg-success text-success-content'
			: 'bg-error text-error-content'}"
	>
		{message}
	</div>
{/if}
