import type { Config } from 'tailwindcss'
import daisyui from 'daisyui'

const config: Config = {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				primary: '#3B82F6',
				secondary: '#10B981',
				accent: '#F59E0B',
			},
			fontFamily: {
				sans: ['Inter', 'ui-sans-serif', 'system-ui', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Helvetica Neue', 'Arial', 'Noto Sans', 'sans-serif', 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji'],
			},
		},
	},
	plugins: [daisyui],
	daisyui: {
		themes: [
			{
				light: {
					"primary": "#3B82F6",
					"secondary": "#10B981",
					"accent": "#F59E0B",
					"neutral": "#3D4451",
					"base-100": "#FFFFFF",
				},
			},
		],
	},
}

export default config