<script>
	export let type = 'button';
	export let variant = 'primary'; // primary, secondary, outline
	export let disabled = false;
	export let loading = false;

	$: classes = {
		primary: 'bg-brand-primary text-white hover:opacity-90',
		secondary: 'bg-gray-200 text-gray-900 hover:bg-gray-300',
		outline: 'border-2 border-brand-primary text-brand-primary hover:bg-brand-primary hover:text-white'
	}[variant];
</script>

<button
	{type}
	{disabled}
	on:click
	class="w-full font-semibold py-3 px-4 rounded-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed {classes}"
	class:opacity-50={loading}
>
	{#if loading}
		<span class="flex items-center justify-center">
			<svg class="animate-spin -ml-1 mr-3 h-5 w-5" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
			<slot name="loading">Loading...</slot>
		</span>
	{:else}
		<slot />
	{/if}
</button>
