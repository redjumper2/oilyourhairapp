<script>
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api } from '$lib/api';

	let loading = true;
	let error = null;

	// Get token and redirect URL from query params
	$: token = $page.url.searchParams.get('token');
	$: redirectUrl = $page.url.searchParams.get('redirect');

	onMount(async () => {
		if (!token) {
			error = 'Missing verification token';
			loading = false;
			return;
		}

		if (!redirectUrl) {
			error = 'Missing redirect URL';
			loading = false;
			return;
		}

		try {
			// Verify magic link token
			const result = await api.verifyMagicLink(token);

			// Redirect to domain with JWT in hash
			const redirectWithToken = `${redirectUrl}#token=${result.token}`;
			window.location.href = redirectWithToken;
		} catch (err) {
			error = err.message || 'Invalid or expired magic link';
			loading = false;
		}
	});
</script>

<svelte:head>
	<title>Verifying...</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center px-4">
	<div class="max-w-md w-full text-center">
		{#if loading}
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900 mx-auto mb-4"></div>
			<h2 class="text-xl font-semibold mb-2">Verifying your magic link...</h2>
			<p class="text-gray-600">You'll be redirected in a moment.</p>
		{:else if error}
			<div class="bg-red-50 border border-red-200 rounded-lg p-6">
				<h2 class="text-xl font-semibold text-red-900 mb-2">Verification Failed</h2>
				<p class="text-red-700 mb-4">{error}</p>
				<p class="text-sm text-gray-600">The link may have expired or already been used.</p>
			</div>
		{/if}
	</div>
</div>
