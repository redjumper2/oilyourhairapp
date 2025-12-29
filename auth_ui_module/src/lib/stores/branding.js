import { writable } from 'svelte/store';

/**
 * Store for domain branding
 */
export const branding = writable({
	company_name: '',
	primary_color: '#000000',
	secondary_color: '#666666',
	logo_url: '',
	support_email: ''
});

/**
 * Apply branding to CSS variables
 * @param {Object} brandingData - Branding configuration
 */
export function applyBranding(brandingData) {
	branding.set(brandingData);

	if (typeof document !== 'undefined') {
		const root = document.documentElement;
		root.style.setProperty('--brand-primary', brandingData.primary_color || '#000000');
		root.style.setProperty('--brand-secondary', brandingData.secondary_color || '#666666');
	}
}
