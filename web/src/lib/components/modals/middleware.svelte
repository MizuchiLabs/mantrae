<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import MiddlewareForm from '../forms/middleware.svelte';
	import type { Middleware, UpsertMiddlewareParams } from '$lib/types/middlewares';
	import { Button } from '$lib/components/ui/button/index.js';
	import { api, profile, loading } from '$lib/api';
	import { toast } from 'svelte-sonner';

	interface Props {
		middleware?: Middleware;
		open?: boolean;
		mode: 'create' | 'edit';
	}

	const defaultMiddleware: Middleware = {
		name: '',
		protocol: 'http'
	};

	let {
		middleware = $bindable(defaultMiddleware),
		open = $bindable(false),
		mode = 'create'
	}: Props = $props();

	let disabled = middleware.name.split('@')[1] !== 'http';

	const update = async () => {
		try {
			// Ensure proper name formatting and synchronization
			if (!middleware.name.includes('@')) {
				middleware.name = `${middleware.name}@http`;
			}

			let params: UpsertMiddlewareParams = {
				name: middleware.name,
				protocol: middleware.protocol,
				type: middleware.type
			};
			switch (middleware.type) {
				case 'http':
					params.middleware = middleware;
					break;
				case 'tcp':
					params.tcpMiddleware = middleware;
					break;
			}

			await api.upsertMiddleware($profile.id, params);
			open = false;
			toast.success(`Middleware ${mode === 'create' ? 'created' : 'updated'} successfully`);
		} catch (err: unknown) {
			const e = err as Error;
			toast.error(`Failed to ${mode} router`, {
				description: e.message
			});
		}
	};

	const handleDelete = async () => {
		if (!middleware.name) return;

		try {
			let params: Middleware = {
				name: middleware.name,
				protocol: middleware.protocol
			};
			await api.deleteMiddleware($profile.id, params);
			toast.success('Middleware deleted successfully');
			open = false;
		} catch (err: unknown) {
			const e = err as Error;
			toast.error('Failed to delete middleware', {
				description: e.message
			});
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[80vh] max-w-2xl overflow-y-auto">
		<MiddlewareForm bind:middleware {mode} {disabled} />
		<Dialog.Footer>
			{#if middleware.name}
				<Button type="button" variant="destructive" onclick={handleDelete} disabled={$loading}
					>Delete</Button
				>
			{/if}
			<Button type="submit" onclick={() => update()} disabled={$loading}
				>{middleware.name ? 'Update' : 'Save'}</Button
			>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
