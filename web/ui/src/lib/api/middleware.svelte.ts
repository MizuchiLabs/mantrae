import { toast } from 'svelte-sonner';
import { useMutation, useQuery } from '$lib/query';
import { profileID } from '$lib/store.svelte';
import { MiddlewareService } from '$lib/gen/mantrae/v1/middleware_pb';

export const middleware = {
	// Queries
	get: (id: string) =>
		useQuery(
			MiddlewareService.method.getMiddleware,
			{ id },
			{ enabled: !!profileID.current, select: (res) => res.middleware }
		),
	list: (pid?: bigint) =>
		useQuery(
			MiddlewareService.method.listMiddlewares,
			{ profileId: pid ?? profileID.current },
			{ enabled: !!profileID.current, select: (res) => res.middlewares }
		),
	plugins: () =>
		useQuery(
			MiddlewareService.method.getMiddlewarePlugins,
			{},
			{ enabled: !!profileID.current, select: (res) => res.plugins }
		),

	// Mutations
	create: () =>
		useMutation(MiddlewareService.method.createMiddleware, {
			onSuccess: () => toast.success('Middleware created!'),
			transform: (variables) => ({ ...variables, profileId: profileID.current })
		}),
	update: () =>
		useMutation(MiddlewareService.method.updateMiddleware, {
			onSuccess: () => toast.success('Middleware updated!')
		}),
	delete: () =>
		useMutation(MiddlewareService.method.deleteMiddleware, {
			onSuccess: () => toast.success('Middleware deleted!')
		})
};
