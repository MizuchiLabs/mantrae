// Copyright 2021-2023 The Connect Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import type {
	DescMessage,
	DescMethodUnary,
	MessageInitShape,
	MessageShape
} from '@bufbuild/protobuf';
import type { ConnectError, Transport } from '@connectrpc/connect';
import type { ConnectInfiniteQueryOptions, ConnectQueryKey } from '@connectrpc/connect-query-core';
import { createInfiniteQueryOptions } from '@connectrpc/connect-query-core';
import type {
	InfiniteData,
	SkipToken,
	CreateInfiniteQueryOptions as TanStackCreateInfiniteQueryOptions,
	CreateInfiniteQueryResult
} from '@tanstack/svelte-query';
import { createInfiniteQuery as tsCreateInfiniteQuery } from '@tanstack/svelte-query';

import { useTransport } from './use-transport.js';

/**
 * Options for useInfiniteQuery
 */
export type UseInfiniteQueryOptions<
	I extends DescMessage,
	O extends DescMessage,
	ParamKey extends keyof MessageInitShape<I>
> = Omit<
	TanStackCreateInfiniteQueryOptions<
		MessageShape<O>,
		ConnectError,
		InfiniteData<MessageShape<O>>,
		ConnectQueryKey<O>,
		MessageInitShape<I>[ParamKey]
	>,
	'getNextPageParam' | 'initialPageParam' | 'queryFn' | 'queryKey'
> &
	ConnectInfiniteQueryOptions<I, O, ParamKey> & {
		/** The transport to be used for the fetching. */
		transport?: Transport;
	};

/**
 * Query the method provided. Maps to createInfiniteQuery on tanstack/svelte-query
 */
export function useInfiniteQuery<
	I extends DescMessage,
	O extends DescMessage,
	ParamKey extends keyof MessageInitShape<I>
>(
	schema: DescMethodUnary<I, O>,
	input: SkipToken | (MessageInitShape<I> & Required<Pick<MessageInitShape<I>, ParamKey>>),
	{
		transport,
		pageParamKey,
		getNextPageParam,
		...queryOptions
	}: UseInfiniteQueryOptions<I, O, ParamKey>
): CreateInfiniteQueryResult<InfiniteData<MessageShape<O>>, ConnectError> {
	const transportFromCtx = useTransport();
	const baseOptions = createInfiniteQueryOptions(schema, input, {
		transport: transport ?? transportFromCtx,
		getNextPageParam,
		pageParamKey
	});
	return tsCreateInfiniteQuery(() => ({
		...baseOptions,
		...queryOptions
	}));
}
