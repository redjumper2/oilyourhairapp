/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				// These will be overridden dynamically by domain branding
				brand: {
					primary: 'var(--brand-primary, #000000)',
					secondary: 'var(--brand-secondary, #666666)'
				}
			}
		}
	},
	plugins: []
};
