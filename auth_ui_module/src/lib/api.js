import { browser } from '$app/environment';

// Get API base URL from environment variable
const API_BASE = browser
	? (import.meta.env.VITE_API_URL || 'http://localhost:9090/api/v1')
	: 'http://localhost:9090/api/v1';

/**
 * API client for auth backend
 */
class AuthAPI {
	/**
	 * Make API request
	 * @param {string} endpoint - API endpoint path
	 * @param {Object} options - Fetch options
	 * @param {string} domain - Domain for Host header
	 */
	async request(endpoint, options = {}, domain = null) {
		const url = `${API_BASE}${endpoint}`;
		const headers = {
			'Content-Type': 'application/json',
			...options.headers
		};

		// Add Host header if domain provided
		if (domain) {
			headers['Host'] = domain;
		}

		const response = await fetch(url, {
			...options,
			headers
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ error: 'Request failed' }));
			throw new Error(error.error || `HTTP ${response.status}`);
		}

		return response.json();
	}

	// Invitation endpoints
	async verifyInvitation(token) {
		return this.request(`/auth/invitation/verify?token=${token}`);
	}

	async acceptInvitation(token, email, authProvider = 'magic_link', providerId = '') {
		return this.request('/auth/invitation/accept', {
			method: 'POST',
			body: JSON.stringify({
				token,
				email,
				auth_provider: authProvider,
				provider_id: providerId
			})
		});
	}

	// Magic link endpoints
	async requestMagicLink(email, domain) {
		return this.request('/auth/magic-link/request', {
			method: 'POST',
			body: JSON.stringify({ email })
		}, domain);
	}

	async verifyMagicLink(token) {
		return this.request(`/auth/magic-link/verify?token=${token}`);
	}

	// Domain branding
	async getDomainBranding(domain) {
		// This would call admin endpoint or a public branding endpoint
		// For now, we'll extract from invitation or use defaults
		return {
			company_name: domain,
			primary_color: '#000000',
			logo_url: '',
			support_email: `support@${domain}`
		};
	}
}

export const api = new AuthAPI();
