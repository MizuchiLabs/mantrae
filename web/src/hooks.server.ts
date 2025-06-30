import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) =>
	await resolve(event, {
		filterSerializedResponseHeaders: (name) => name === 'content-type'
	});
