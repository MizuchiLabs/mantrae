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

import type { Transport } from '@connectrpc/connect';
import { ConnectError } from '@connectrpc/connect';
import { getContext, setContext } from 'svelte';

const fallbackTransportError = new ConnectError(
	"To use Connect, you must provide a `Transport`: a simple object that handles `unary` and `stream` requests. `Transport` objects can easily be created by using `@connectrpc/connect-web`'s exports `createConnectTransport` and `createGrpcWebTransport`. see: https://connectrpc.com/docs/web/getting-started for more info."
);

// istanbul ignore next
export const fallbackTransport: Transport = {
	unary: () => {
		throw fallbackTransportError;
	},
	stream: () => {
		throw fallbackTransportError;
	}
};

const TRANSPORT_CONTEXT_KEY = Symbol.for('connect-query-transport');

/**
 * Use this helper to get the default transport that's currently attached to the Svelte context for the calling component.
 */
export const useTransport = (): Transport => {
	const transport = getContext<Transport>(TRANSPORT_CONTEXT_KEY);
	return transport ?? fallbackTransport;
};

/**
 * Sets the transport in the current Svelte context. This should be called from a parent component
 * to provide transport to all descendant components.
 *
 * @example
 * ```svelte
 * <script>
 *   import { createConnectTransport } from "@connectrpc/connect-web";
 *   import { setTransport } from "@connectrpc/connect-query-svelte";
 *
 *   const transport = createConnectTransport({
 *     baseUrl: "https://demo.connectrpc.com",
 *   });
 *
 *   setTransport(transport);
 * </script>
 *
 * <!-- Your app components -->
 * ```
 */
export const setTransport = (transport: Transport): void => {
	setContext(TRANSPORT_CONTEXT_KEY, transport);
};
