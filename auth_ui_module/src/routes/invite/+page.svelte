<script>
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { applyBranding } from '$lib/stores/branding';

	let loading = true;
	let error = null;
	let invitation = null;
	let accepting = false;

	// Get token and redirect URL from query params
	$: token = $page.url.searchParams.get('token');
	$: redirectUrl = $page.url.searchParams.get('redirect') || 'https://example.com';

	onMount(async () => {
		if (!token) {
			error = 'Missing invitation token';
			loading = false;
			return;
		}

		try {
			// Verify invitation
			invitation = await api.verifyInvitation(token);

			// Apply domain branding
			if (invitation.branding) {
				applyBranding(invitation.branding);
			}

			loading = false;
		} catch (err) {
			error = err.message || 'Failed to verify invitation';
			loading = false;
		}
	});

	async function acceptInvitation() {
		accepting = true;
		error = null;

		try {
			// Accept invitation (creates user and returns JWT)
			const result = await api.acceptInvitation(
				token,
				invitation.email,
				'magic_link',
				''
			);

			// Redirect to domain with JWT in hash
			const redirectWithToken = `${redirectUrl}#token=${result.token}`;
			window.location.href = redirectWithToken;
		} catch (err) {
			error = err.message || 'Failed to accept invitation';
			accepting = false;
		}
	}
</script>

<svelte:head>
	<title>Accept Invitation</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center px-4">
	<div class="max-w-md w-full">
		{#if loading}
			<div class="text-center">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900 mx-auto"></div>
				<p class="mt-4 text-gray-600">Verifying invitation...</p>
			</div>
		{:else if error}
			<div class="bg-red-50 border border-red-200 rounded-lg p-6">
				<h2 class="text-xl font-semibold text-red-900 mb-2">Invalid Invitation</h2>
				<p class="text-red-700">{error}</p>
			</div>
		{:else if invitation}
			<div class="bg-white shadow-lg rounded-lg p-8">
				<!-- Logo -->
				{#if invitation.branding?.logo_url}
					<img
						src={invitation.branding.logo_url}
						alt={invitation.branding.company_name}
						class="h-12 mx-auto mb-6"
					/>
				{/if}

				<!-- Invitation Details -->
				<h1 class="text-2xl font-bold text-center mb-2">You're Invited!</h1>
				<p class="text-center text-gray-600 mb-6">
					Join <span class="font-semibold">{invitation.branding?.company_name || invitation.domain}</span>
					as a <span class="font-semibold capitalize">{invitation.role}</span>
				</p>

				<!-- Details Card -->
				<div class="bg-gray-50 rounded-lg p-4 mb-6 space-y-2">
					<div class="flex justify-between">
						<span class="text-gray-600">Email:</span>
						<span class="font-medium">{invitation.email}</span>
					</div>
					<div class="flex justify-between">
						<span class="text-gray-600">Role:</span>
						<span class="font-medium capitalize">{invitation.role}</span>
					</div>
					{#if invitation.promo_code}
						<div class="flex justify-between">
							<span class="text-gray-600">Promo Code:</span>
							<span class="font-mono font-medium">{invitation.promo_code}</span>
						</div>
					{/if}
					{#if invitation.discount_percent}
						<div class="flex justify-between">
							<span class="text-gray-600">Discount:</span>
							<span class="font-medium">{invitation.discount_percent}% off</span>
						</div>
					{/if}
					{#if invitation.time_remaining}
						<div class="flex justify-between">
							<span class="text-gray-600">Expires in:</span>
							<span class="font-medium text-orange-600">{invitation.time_remaining}</span>
						</div>
					{/if}
				</div>

				<!-- Accept Button -->
				<button
					on:click={acceptInvitation}
					disabled={accepting}
					class="w-full bg-brand-primary text-white font-semibold py-3 px-4 rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50 disabled:cursor-not-allowed"
					style="background-color: var(--brand-primary)"
				>
					{accepting ? 'Accepting...' : 'Accept Invitation'}
				</button>

				{#if error}
					<p class="mt-4 text-red-600 text-sm text-center">{error}</p>
				{/if}

				<!-- Support -->
				{#if invitation.branding?.support_email}
					<p class="mt-6 text-center text-sm text-gray-500">
						Questions? Contact <a
							href="mailto:{invitation.branding.support_email}"
							class="text-brand-primary hover:underline"
							style="color: var(--brand-primary)"
						>
							{invitation.branding.support_email}
						</a>
					</p>
				{/if}
			</div>
		{/if}
	</div>
</div>
