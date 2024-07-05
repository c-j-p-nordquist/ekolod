/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {}
  },
  plugins: [require('daisyui'),],
  daisyui: {
    themes: ['light', 'dark', 'corporate', 'cyberpunk', 'valentine', 'aqua', 'forest', 'nord'],
  },
};