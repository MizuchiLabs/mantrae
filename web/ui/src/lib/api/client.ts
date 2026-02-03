import {
	create,
	type DescMessage,
	type DescMethodServerStreaming,
	type DescMethodUnary,
	type DescService,
	type MessageInitShape,
	type MessageShape
} from '@bufbuild/protobuf';
import { createConnectTransport } from '@connectrpc/connect-web';
import { BackendURL } from '$lib/config';
import { MutationCache, QueryCache, QueryClient } from '@tanstack/svelte-query';
import { browser } from '$app/environment';
import { mutationErrorHandler, queryErrorHandler } from './error';
import { createClient, type Client } from '@connectrpc/connect';

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

export interface CallUnaryOptions {
	signal?: AbortSignal;
}

export interface CallServerStreamOptions {
	signal?: AbortSignal;
	onChunk?: (chunk: Uint8Array, received: number) => void;
}

export async function callUnary<I extends DescMessage, O extends DescMessage>(
	schema: DescMethodUnary<I, O>,
	input?: MessageInitShape<I>,
	options?: CallUnaryOptions
): Promise<MessageShape<O> | undefined> {
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

export function createServiceClient<S extends DescService>(service: S): Client<S> {
	return createClient(service, transport);
}

export async function collectServerStream<O>(
	stream: AsyncIterable<O>,
	extractChunk: (message: O) => Uint8Array,
	options?: CallServerStreamOptions
): Promise<Uint8Array> {
	const chunks: Uint8Array[] = [];
	let totalReceived = 0;

	for await (const response of stream) {
		const chunk = extractChunk(response);
		if (chunk.length > 0) {
			chunks.push(chunk);
			totalReceived += chunk.byteLength;
			options?.onChunk?.(chunk, totalReceived);
		}
	}

	// Concatenate all chunks into a single Uint8Array
	const result = new Uint8Array(totalReceived);
	let offset = 0;
	for (const chunk of chunks) {
		result.set(chunk, offset);
		offset += chunk.byteLength;
	}

	return result;
}

export function downloadBlob(
	data: Uint8Array,
	filename: string,
	mimeType = 'application/octet-stream'
) {
	const blob = new Blob([data.buffer as ArrayBuffer], { type: mimeType });
	const url = URL.createObjectURL(blob);

	const a = document.createElement('a');
	a.href = url;
	a.download = filename;
	a.click();

	URL.revokeObjectURL(url);
}

// import { collectServerStream, downloadBlob } from '$lib/api/client';

// async function downloadBackup(name?: string) {
// 	try {
// 		const stream = backupClient.downloadBackup({ name });
// 		const data = await collectServerStream(stream, (chunk) => chunk.data);
// 		downloadBlob(data, name || 'backup.db');
// 	} catch (err) {
// 		const e = ConnectError.from(err);
// 		toast.error('Failed to download backup', { description: e.message });
// 	}
// }
