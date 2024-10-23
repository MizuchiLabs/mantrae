import type { Middleware } from '$lib/types/middlewares';
import type { SvelteComponent } from 'svelte';

export const MiddlewareForms = {
	addPrefix: import('$lib/components/forms/addPrefix.svelte'),
	stripPrefix: import('$lib/components/forms/stripPrefix.svelte'),
	stripPrefixRegex: import('$lib/components/forms/stripPrefixRegex.svelte'),
	replacePath: import('$lib/components/forms/replacePath.svelte'),
	replacePathRegex: import('$lib/components/forms/replacePathRegex.svelte'),
	chain: import('$lib/components/forms/chain.svelte'),
	ipAllowList: import('$lib/components/forms/ipAllowList.svelte'),
	headers: import('$lib/components/forms/headers.svelte'),
	errors: import('$lib/components/forms/errorPage.svelte'),
	rateLimit: import('$lib/components/forms/rateLimit.svelte'),
	redirectRegex: import('$lib/components/forms/redirectRegex.svelte'),
	redirectScheme: import('$lib/components/forms/redirectScheme.svelte'),
	basicAuth: import('$lib/components/forms/basicAuth.svelte'),
	digestAuth: import('$lib/components/forms/digestAuth.svelte'),
	forwardAuth: import('$lib/components/forms/forwardAuth.svelte'),
	inFlightReq: import('$lib/components/forms/inFlightReq.svelte'),
	buffering: import('$lib/components/forms/buffering.svelte'),
	circuitBreaker: import('$lib/components/forms/circuitBreaker.svelte'),
	compress: import('$lib/components/forms/compress.svelte'),
	passTLSClientCert: import('$lib/components/forms/passTLSClientCert.svelte'),
	retry: import('$lib/components/forms/retry.svelte'),
	plugin: import('$lib/components/forms/plugin.svelte'),

	// TCP-specific
	inFlightConn: import('$lib/components/forms/inFlightConn.svelte')
};

export const LoadMiddlewareForm = async (
	mw: Middleware
): Promise<typeof SvelteComponent | null> => {
	const moduleKey = Object.keys(MiddlewareForms).find(
		(key) => key.toLowerCase() === mw.type?.toLowerCase()
	) as keyof typeof MiddlewareForms | undefined;

	return moduleKey ? (await MiddlewareForms[moduleKey]).default : null;
};
