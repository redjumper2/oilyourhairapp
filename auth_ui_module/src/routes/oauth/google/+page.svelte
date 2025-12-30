<script>
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';

	// Get domain and redirect URL from query params
	$: domain = $page.url.searchParams.get('domain');
	$: redirectUrl = $page.url.searchParams.get('redirect');

	onMount(() => {
		if (!domain) {
			return;
		}

		// Build OAuth URL with domain and redirect params
		const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';
		const params = new URLSearchParams();
		params.set('domain', domain);
		if (redirectUrl) params.set('redirect', redirectUrl);

		// Redirect to backend Google OAuth endpoint
		// The backend will handle the OAuth flow and redirect back
		window.location.href = `${apiUrl}/auth/google?${params.toString()}`;
	});
</script>

<svelte:head>
	<title>Signing in with Google...</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center px-4">
	<div class="max-w-md w-full text-center">
		{#if !domain}
			<div class="bg-red-50 border border-red-200 rounded-lg p-6">
				<h2 class="text-xl font-semibold text-red-900 mb-2">Missing Domain</h2>
				<p class="text-red-700">Domain parameter is required for Google OAuth.</p>
			</div>
		{:else}
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900 mx-auto mb-4"></div>
			<h2 class="text-xl font-semibold mb-2">Redirecting to Google...</h2>
			<p class="text-gray-600">Please wait while we redirect you to sign in.</p>
		{/if}
	</div>
</div>
