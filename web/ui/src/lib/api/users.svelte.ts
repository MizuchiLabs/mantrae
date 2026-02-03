import { toast } from 'svelte-sonner';
import { useMutation, useQuery } from '$lib/query';
import { UserService } from '$lib/gen/mantrae/v1/user_pb';
import { goto } from '$app/navigation';
import { queryClient } from './client';

export const user = {
	// Queries
	self: () => useQuery(UserService.method.getUser, {}, { select: (res) => res.user }),
	// get: (id: string) => useQuery(UserService.method.getUser, { id }, { select: (res) => res.user }),
	list: () => useQuery(UserService.method.listUsers, {}, { select: (res) => res.users }),
	oidc: () => useQuery(UserService.method.getOIDCStatus, {}, {}),

	// Mutations
	login: () =>
		useMutation(UserService.method.loginUser, {
			onSuccess: () => {
				goto('/');
				toast.success('Welcome back 👋');
			}
		}),
	logout: () =>
		useMutation(UserService.method.logoutUser, {
			onSuccess: () => {
				queryClient.cancelQueries();
				queryClient.clear();
				goto('/login');
				toast.success('Logged out 👋');
			}
		}),
	create: () =>
		useMutation(UserService.method.createUser, {
			onSuccess: () => toast.success('User created')
		}),
	update: () =>
		useMutation(UserService.method.updateUser, {
			onSuccess: () => toast.success('User updated')
		}),
	delete: () =>
		useMutation(UserService.method.deleteUser, {
			onSuccess: () => toast.success('User deleted')
		})
};
