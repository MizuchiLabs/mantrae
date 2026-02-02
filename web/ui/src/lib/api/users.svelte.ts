import { toast } from 'svelte-sonner';
import { useMutation, useQuery } from '$lib/query';
import { UserService } from '$lib/gen/mantrae/v1/user_pb';

export const user = {
	// Queries
	self: () => useQuery(UserService.method.getUser, {}, { select: (res) => res.user }),

	// Mutations
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
