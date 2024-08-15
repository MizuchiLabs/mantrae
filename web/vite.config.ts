import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	define: {
		'process.env': process.env
	}
});
process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0';
