import { goto } from '$app/navigation';
import { Code, ConnectError } from '@connectrpc/connect';
import { Mutation, Query, type QueryKey } from '@tanstack/svelte-query';
import { toast } from 'svelte-sonner';
import { formatConnectError } from './error-formatter';

const errorHandler = (error: unknown): void => {
	if (!(error instanceof ConnectError)) {
		toast.error((error as Error).message);
		return;
	}

	if (error.code === Code.Unauthenticated || error.code === Code.Aborted) {
		goto('/login');
		return;
	}

	const { title, description } = formatConnectError(error);
	toast.error(title, description ? { description } : undefined);
};

export const queryErrorHandler = (
	error: unknown,
	_query: Query<unknown, unknown, unknown, QueryKey>
) => {
	errorHandler(error);
};

export const mutationErrorHandler = (
	error: unknown,
	_variables: unknown,
	_context: unknown,
	_mutation: Mutation<unknown, unknown, unknown, unknown>
) => {
	errorHandler(error);
};
