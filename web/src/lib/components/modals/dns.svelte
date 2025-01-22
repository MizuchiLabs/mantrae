<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import type { DNSProvider } from '$lib/types';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { api, loading } from '$lib/api';
	import { toast } from 'svelte-sonner';

	interface Props {
		dns: DNSProvider | undefined;
		open?: boolean;
	}

	let { dns: dns = $bindable({} as DNSProvider), open = $bindable(false) }: Props = $props();

	let providerTypes = [
		{ value: 'cloudflare', label: 'Cloudflare' },
		{ value: 'powerdns', label: 'PowerDNS' },
		{ value: 'technitium', label: 'Technitium' }
	];
	const handleSubmit = async () => {
		try {
			if (dns.id) {
				await api.updateDNSProvider(dns);
				toast.success('DNS Provider updated successfully');
			} else {
				await api.createDNSProvider(dns);
				toast.success('DNS Provider created successfully');
			}
			open = false;
		} catch (err: unknown) {
			const e = err as Error;
			toast.error('Failed to save dnsProvider', {
				description: e.message
			});
		}
	};

	const handleDelete = async () => {
		if (!dns.id) return;

		try {
			await api.deleteProfile(dns.id);
			toast.success('DNS Provider deleted successfully');
			open = false;
		} catch (err: unknown) {
			const e = err as Error;
			toast.error('Failed to delete dnsProvider', {
				description: e.message
			});
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-screen overflow-y-auto sm:max-w-[550px]">
		<Dialog.Header>
			<Dialog.Title>{dns.id ? 'Edit' : 'Add'} Provider</Dialog.Title>
			<!-- <Dialog.Description>Configure your Traefik instance connection details.</Dialog.Description> -->
		</Dialog.Header>

		<form onsubmit={handleSubmit} class="space-y-4">
			<div class="mb-4 flex items-center justify-end gap-2">
				<Tooltip.Root>
					<Tooltip.Trigger>
						<Label for="is_active" class="text-right">Default</Label>
						<Switch name="is_active" bind:checked={dns.isActive} required />
					</Tooltip.Trigger>
					<Tooltip.Content class="max-w-sm">
						<p>Sets this provider as the default, any new router created will use this provider</p>
					</Tooltip.Content>
				</Tooltip.Root>
			</div>

			<div class="grid grid-cols-4 items-center gap-2 space-y-2">
				<Label for="current" class="text-right">Type</Label>
				<Select.Root type="single" value={dns.type} onValueChange={(value) => (dns.type = value)}>
					<Select.Trigger class="col-span-3">
						{dns.type ? dns.type : 'Select type'}
					</Select.Trigger>
					<Select.Content class="no-scrollbar max-h-[300px] overflow-y-auto">
						{#each providerTypes as type}
							<Select.Item value={type.value} label={type.label}>
								{type.label}
							</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>

			<div class="grid grid-cols-4 items-center gap-2">
				<Label for="name" class="text-right">Name</Label>
				<Input
					name="name"
					type="text"
					bind:value={dns.name}
					placeholder="Name of the provider"
					required
				/>
			</div>

			<Dialog.Footer>
				{#if dns.id}
					<Button type="button" variant="destructive" onclick={handleDelete} disabled={$loading}
						>Delete</Button
					>
				{/if}
				<Button type="submit" disabled={$loading}>{dns.id ? 'Update' : 'Save'}</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
