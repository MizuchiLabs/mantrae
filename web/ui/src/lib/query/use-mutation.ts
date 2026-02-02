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
import { callUnaryMethod } from '@connectrpc/connect-query-core';
import type {
	CreateMutationResult,
	CreateMutationOptions as TSCreateMutationOptions
} from '@tanstack/svelte-query';
import { createMutation as tsCreateMutation } from '@tanstack/svelte-query';

import { useTransport } from './use-transport.js';

/**
 * Options for useMutation
 */
export type UseMutationOptions<
	I extends DescMessage,
	O extends DescMessage,
	Ctx = unknown
> = TSCreateMutationOptions<
	MessageShape<O>,
	ConnectError,
	Omit<MessageInitShape<I>, '$typeName'>,
	Ctx
> & {
	/** The transport to be used for the fetching. */
	transport?: Transport;
	/** Transform the input before sending the request. */
	transform?: (
		input: Omit<MessageInitShape<I>, '$typeName'>
	) => Omit<MessageInitShape<I>, '$typeName'>;
};

/**
 * Query the method provided. Maps to createMutation on tanstack/svelte-query
 */
export function useMutation<I extends DescMessage, O extends DescMessage, Ctx = unknown>(
	schema: DescMethodUnary<I, O>,
	{ transport, transform, ...queryOptions }: UseMutationOptions<I, O, Ctx> = {}
): CreateMutationResult<
	MessageShape<O>,
	ConnectError,
	Omit<MessageInitShape<I>, '$typeName'>,
	Ctx
> {
	const transportFromCtx = useTransport();
	const transportToUse = transport ?? transportFromCtx;

	const mutationFn = async (input: Omit<MessageInitShape<I>, '$typeName'>) =>
		callUnaryMethod(
			transportToUse,
			schema,
			transform ? (transform(input) as MessageInitShape<I>) : (input as MessageInitShape<I>)
		);

	return tsCreateMutation(() => ({
		...queryOptions,
		mutationFn
	}));
}
