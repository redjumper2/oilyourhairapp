<script>
	import { page } from '$app/stores';
</script>

<svelte:head>
	<title>Auth Portal</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center px-4 bg-gradient-to-br from-gray-50 to-gray-100">
	<div class="max-w-2xl w-full text-center">
		<h1 class="text-4xl font-bold mb-4 text-gray-900">Authentication Portal</h1>
		<p class="text-lg text-gray-600 mb-8">
			Secure multi-tenant authentication for your applications
		</p>

		<div class="bg-white shadow-lg rounded-lg p-8 text-left">
			<h2 class="text-xl font-semibold mb-4">Available Endpoints</h2>
			<ul class="space-y-3 text-gray-700">
				<li class="flex items-start">
					<span class="text-green-600 mr-2">→</span>
					<div>
						<code class="bg-gray-100 px-2 py-1 rounded text-sm">/invite?token=...</code>
						<p class="text-sm text-gray-600 mt-1">Accept invitation</p>
					</div>
				</li>
				<li class="flex items-start">
					<span class="text-green-600 mr-2">→</span>
					<div>
						<code class="bg-gray-100 px-2 py-1 rounded text-sm">/login?domain=...</code>
						<p class="text-sm text-gray-600 mt-1">Sign in with magic link or Google</p>
					</div>
				</li>
				<li class="flex items-start">
					<span class="text-green-600 mr-2">→</span>
					<div>
						<code class="bg-gray-100 px-2 py-1 rounded text-sm">/verify?token=...</code>
						<p class="text-sm text-gray-600 mt-1">Verify magic link</p>
					</div>
				</li>
				<li class="flex items-start">
					<span class="text-green-600 mr-2">→</span>
					<div>
						<code class="bg-gray-100 px-2 py-1 rounded text-sm">/oauth/google?domain=...</code>
						<p class="text-sm text-gray-600 mt-1">Google OAuth flow</p>
					</div>
				</li>
			</ul>
		</div>

		<p class="mt-8 text-sm text-gray-500">
			This is a centralized authentication portal. Users will be redirected here from your application.
		</p>
	</div>
</div>
