import {
	create,
	type DescMessage,
	type DescMethodUnary,
	type MessageInitShape,
	type MessageShape
} from '@bufbuild/protobuf';
import { createConnectTransport } from '@connectrpc/connect-web';
import { BackendURL } from '$lib/config';
import { MutationCache, QueryCache, QueryClient } from '@tanstack/svelte-query';
import { browser } from '$app/environment';
import { mutationErrorHandler, queryErrorHandler } from './error';

export const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			enabled: browser,
			retry: false,
			refetchOnMount: true,
			refetchOnReconnect: true,
			refetchOnWindowFocus: true,
			refetchIntervalInBackground: false,
			refetchInterval: 300000 // 5min
		},
		mutations: {
			retry: false
		}
	},
	queryCache: new QueryCache({
		onError: queryErrorHandler
	}),
	mutationCache: new MutationCache({
		onError: mutationErrorHandler,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['connect-query'] });
		}
	})
});

export const transport = createConnectTransport({
	baseUrl: BackendURL,
	useHttpGet: true,
	fetch: (input, init) => {
		return fetch(input, { ...init, credentials: 'include' });
	},
	jsonOptions: {
		ignoreUnknownFields: true
	}
});

export type Result<T> = {
	success: boolean;
	response?: T;
	error?: string;
};

export interface CallUnaryOptions {
	signal?: AbortSignal;
}

export async function callUnary<I extends DescMessage, O extends DescMessage>(
	schema: DescMethodUnary<I, O>,
	input?: MessageInitShape<I>,
	customFetch?: typeof fetch,
	options?: CallUnaryOptions
): Promise<MessageShape<O> | undefined> {
	const wrappedFetch: typeof fetch = (input, init = {}) => {
		return (customFetch || fetch)(input, {
			...init,
			credentials: 'include'
		});
	};

	const transport = createConnectTransport({
		baseUrl: BackendURL,
		fetch: wrappedFetch
	});

	// Helper to actually perform the request
	const doRequest = async () => {
		return await transport.unary(
			schema,
			options?.signal,
			undefined,
			undefined,
			input ?? create(schema.input),
			undefined
		);
	};

	try {
		const result = await doRequest();
		if (result.message) return result.message;
	} catch (error) {
		return Promise.reject(error);
	}
}
