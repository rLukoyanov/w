/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				base: '#1e1e2e',
				surface: '#181825',
				overlay: '#313244',
				text: '#cdd6f4',
				subtext: '#a6adc8',
				blue: '#89b4fa',
				sky: '#74c7ec',
				red: '#f38ba8',
				pink: '#eba0ac'
			}
		}
	},
	plugins: []
};
