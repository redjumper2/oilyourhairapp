<script>
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api } from '$lib/api';
	import { applyBranding } from '$lib/stores/branding';

	let email = '';
	let loading = false;
	let error = null;
	let success = false;

	// Get domain and redirect URL from query params
	$: domain = $page.url.searchParams.get('domain');
	$: redirectUrl = $page.url.searchParams.get('redirect') || `https://${domain}`;

	onMount(async () => {
		if (!domain) {
			error = 'Missing domain parameter';
			return;
		}

		try {
			// Fetch and apply domain branding
			const brandingData = await api.getDomainBranding(domain);
			applyBranding(brandingData);
		} catch (err) {
			console.error('Failed to load branding:', err);
		}
	});

	async function requestMagicLink(event) {
		event.preventDefault();
		loading = true;
		error = null;
		success = false;

		try {
			await api.requestMagicLink(email, domain);
			success = true;
		} catch (err) {
			error = err.message || 'Failed to send magic link';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Login</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center px-4">
	<div class="max-w-md w-full">
		<div class="bg-white shadow-lg rounded-lg p-8">
			<!-- Logo placeholder -->
			<div class="text-center mb-8">
				<h1 class="text-3xl font-bold" style="color: var(--brand-primary)">
					{domain || 'Welcome'}
				</h1>
			</div>

			{#if success}
				<!-- Success State -->
				<div class="text-center">
					<div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-green-100 mb-4">
						<svg
							class="h-6 w-6 text-green-600"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M5 13l4 4L19 7"
							/>
						</svg>
					</div>
					<h2 class="text-xl font-semibold mb-2">Check your email</h2>
					<p class="text-gray-600 mb-4">
						We've sent a magic link to <span class="font-medium">{email}</span>
					</p>
					<p class="text-sm text-gray-500">Click the link in your email to sign in.</p>
				</div>
			{:else}
				<!-- Login Form -->
				<form on:submit={requestMagicLink} class="space-y-6">
					<div>
						<label for="email" class="block text-sm font-medium text-gray-700 mb-2">
							Email address
						</label>
						<input
							id="email"
							type="email"
							required
							bind:value={email}
							disabled={loading}
							class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-primary focus:border-transparent disabled:opacity-50"
							placeholder="you@example.com"
						/>
					</div>

					{#if error}
						<div class="bg-red-50 border border-red-200 rounded-lg p-3">
							<p class="text-sm text-red-700">{error}</p>
						</div>
					{/if}

					<button
						type="submit"
						disabled={loading}
						class="w-full bg-brand-primary text-white font-semibold py-3 px-4 rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50 disabled:cursor-not-allowed"
						style="background-color: var(--brand-primary)"
					>
						{loading ? 'Sending...' : 'Send Magic Link'}
					</button>
				</form>

				<!-- Divider -->
				<div class="relative my-6">
					<div class="absolute inset-0 flex items-center">
						<div class="w-full border-t border-gray-300"></div>
					</div>
					<div class="relative flex justify-center text-sm">
						<span class="px-2 bg-white text-gray-500">Or continue with</span>
					</div>
				</div>

				<!-- Google OAuth Button -->
				<a
					href="/oauth/google?domain={domain}&redirect={redirectUrl}"
					class="flex items-center justify-center w-full bg-white border border-gray-300 text-gray-700 font-semibold py-3 px-4 rounded-lg hover:bg-gray-50 transition-colors"
				>
					<svg class="w-5 h-5 mr-2" viewBox="0 0 24 24">
						<path
							fill="currentColor"
							d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
						/>
						<path
							fill="currentColor"
							d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
						/>
						<path
							fill="currentColor"
							d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
						/>
						<path
							fill="currentColor"
							d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
						/>
					</svg>
					Sign in with Google
				</a>
			{/if}
		</div>
	</div>
</div>
