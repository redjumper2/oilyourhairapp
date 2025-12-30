/**
 * Simple API client for auth backend
 */
class AuthAPI {
    constructor() {
        this.baseURL = 'https://api.oilyourhair.com/api/v1';
    }

    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;

        try {
            const response = await fetch(url, {
                ...options,
                headers: {
                    'Content-Type': 'application/json',
                    ...options.headers
                }
            });

            if (!response.ok) {
                const error = await response.json().catch(() => ({ error: 'Request failed' }));
                throw new Error(error.error || `HTTP ${response.status}`);
            }

            return await response.json();
        } catch (error) {
            console.error('API request failed:', error);
            throw error;
        }
    }

    async getDomainBranding(domain) {
        // Return default branding for now
        // TODO: Add actual API endpoint if needed
        return {
            company_name: domain,
            primary_color: '#2E7D32',
            logo_url: '',
            support_email: `support@${domain}`
        };
    }

    async requestMagicLink(email, domain) {
        console.log('🔧 auth-api.js requestMagicLink called:');
        console.log('   email param:', email);
        console.log('   domain param:', domain);
        const payload = { email, domain };
        console.log('   payload object:', payload);
        console.log('   stringified:', JSON.stringify(payload));

        return this.request('/auth/magic-link/request', {
            method: 'POST',
            body: JSON.stringify(payload)
        });
    }

    async verifyInvitation(token) {
        return this.request(`/auth/invitation/verify?token=${token}`);
    }

    async acceptInvitation(token, email) {
        return this.request('/auth/invitation/accept', {
            method: 'POST',
            body: JSON.stringify({
                token,
                email,
                auth_provider: 'magic_link',
                provider_id: ''
            })
        });
    }
}

// Create global API instance
const api = new AuthAPI();
