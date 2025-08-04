import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import tailwindcss from '@tailwindcss/vite';
import { compression } from 'vite-plugin-compression2';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit(), compression()],
	define: {},
	server: {
		host: '127.0.0.1'
	}
});
