import type { Middleware } from '$lib/types/middlewares';
import type { SvelteComponent } from 'svelte';

export const MiddlewareForms = {
	addPrefix: (await import('$lib/components/forms/addPrefix.svelte')).default,
	stripPrefix: (await import('$lib/components/forms/stripPrefix.svelte')).default,
	stripPrefixRegex: (await import('$lib/components/forms/stripPrefixRegex.svelte')).default,
	replacePath: (await import('$lib/components/forms/replacePath.svelte')).default,
	replacePathRegex: (await import('$lib/components/forms/replacePathRegex.svelte')).default,
	chain: (await import('$lib/components/forms/chain.svelte')).default,
	ipAllowList: (await import('$lib/components/forms/ipAllowList.svelte')).default,
	headers: (await import('$lib/components/forms/headers.svelte')).default,
	errors: (await import('$lib/components/forms/errorPage.svelte')).default,
	rateLimit: (await import('$lib/components/forms/rateLimit.svelte')).default,
	redirectRegex: (await import('$lib/components/forms/redirectRegex.svelte')).default,
	redirectScheme: (await import('$lib/components/forms/redirectScheme.svelte')).default,
	basicAuth: (await import('$lib/components/forms/basicAuth.svelte')).default,
	digestAuth: (await import('$lib/components/forms/digestAuth.svelte')).default,
	forwardAuth: (await import('$lib/components/forms/forwardAuth.svelte')).default,
	inFlightReq: (await import('$lib/components/forms/inFlightReq.svelte')).default,
	buffering: (await import('$lib/components/forms/buffering.svelte')).default,
	circuitBreaker: (await import('$lib/components/forms/circuitBreaker.svelte')).default,
	compress: (await import('$lib/components/forms/compress.svelte')).default,
	// passTLSClientCert: (await import('$lib/components/forms/passTLSClientCert.svelte')).default,
	retry: (await import('$lib/components/forms/retry.svelte')).default,

	// TCP-specific
	inFlightConn: (await import('$lib/components/forms/inFlightConn.svelte')).default,
	tcpIpAllowList: (await import('$lib/components/forms/tcpIpAllowList.svelte')).default
};

export const LoadMiddlewareForm = async (
	mw: Middleware
): Promise<typeof SvelteComponent | null> => {
	const moduleKey = Object.keys(MiddlewareForms).find(
		(key) => key.toLowerCase() === mw.type?.toLowerCase()
	) as keyof typeof MiddlewareForms | undefined;

	return moduleKey ? MiddlewareForms[moduleKey] : null;
};
