<script>
	export let title = '';
	export let description = '';
</script>

<div class="bg-white shadow-lg rounded-lg p-8">
	{#if title}
		<h2 class="text-2xl font-bold text-center mb-2">{title}</h2>
	{/if}

	{#if description}
		<p class="text-center text-gray-600 mb-6">{description}</p>
	{/if}

	<slot />
</div>
