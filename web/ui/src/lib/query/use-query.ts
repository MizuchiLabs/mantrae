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
import type { ConnectQueryKey, SkipToken } from '@connectrpc/connect-query-core';
import { createQueryOptions } from '@connectrpc/connect-query-core';
import type {
	CreateQueryResult,
	CreateQueryOptions as TanStackCreateQueryOptions
} from '@tanstack/svelte-query';
import { createQuery as tsCreateQuery } from '@tanstack/svelte-query';

import { useTransport } from './use-transport.js';

/**
 * Options for createQuery
 */
export type UseQueryOptions<O extends DescMessage, SelectOutData = MessageShape<O>> = Omit<
	TanStackCreateQueryOptions<MessageShape<O>, ConnectError, SelectOutData, ConnectQueryKey<O>>,
	'queryFn' | 'queryKey'
> & {
	/** The transport to be used for the fetching. */
	transport?: Transport;
};

/**
 * Query the method provided. Maps to createQuery on tanstack/svelte-query
 */
export function useQuery<
	I extends DescMessage,
	O extends DescMessage,
	SelectOutData = MessageShape<O>
>(
	schema: DescMethodUnary<I, O>,
	input?: SkipToken | MessageInitShape<I>,
	{ transport, ...queryOptions }: UseQueryOptions<O, SelectOutData> = {}
): CreateQueryResult<SelectOutData, ConnectError> {
	const transportFromCtx = useTransport();
	const baseOptions = createQueryOptions(schema, input, {
		transport: transport ?? transportFromCtx
	});
	return tsCreateQuery(() => ({
		...baseOptions,
		...queryOptions
	}));
}
