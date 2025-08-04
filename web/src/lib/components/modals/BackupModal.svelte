<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { Button } from '$lib/components/ui/button';
	import { Separator } from '$lib/components/ui/separator';
	import { CalendarDays, DatabaseBackup, Download, Trash2 } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { DateFormat } from '$lib/stores/common';
	import { backupClient } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Backup } from '$lib/gen/mantrae/v1/backup_pb';

	let { open = $bindable(false) } = $props();

	let sqliteBackups = $state<Backup[]>([]);
	let yamlBackups = $state<Backup[]>([]);

	async function deleteBackup(name: string) {
		await backupClient.deleteBackup({ name });
		const response = await backupClient.listBackups({});
		sqliteBackups = response.backups.filter((b) => b.name.endsWith('.db'));
		yamlBackups = response.backups.filter((b) => b.name.endsWith('.yaml'));
	}
	async function createBackup() {
		await backupClient.createBackup({});
		const response = await backupClient.listBackups({});
		sqliteBackups = response.backups.filter((b) => b.name.endsWith('.db'));
		yamlBackups = response.backups.filter((b) => b.name.endsWith('.yaml'));
	}
	async function downloadBackup(name?: string) {
		try {
			const stream = backupClient.downloadBackup({ name });

			const chunks: Uint8Array[] = [];
			for await (const chunk of stream) {
				if (chunk.data.length > 0) {
					chunks.push(chunk.data);
				}
			}

			const blob = new Blob(chunks, { type: 'application/octet-stream' });
			const url = URL.createObjectURL(blob);

			const a = document.createElement('a');
			a.href = url;
			a.download = name || 'backup.db';
			a.click();

			URL.revokeObjectURL(url);
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to download backup', { description: e.message });
		}
	}

	onMount(async () => {
		const response = await backupClient.listBackups({});
		sqliteBackups = response.backups.filter((b) => b.name.endsWith('.db'));
		yamlBackups = response.backups.filter(
			(b) => b.name.endsWith('.yaml') || b.name.endsWith('.yml')
		);
	});
</script>

<svelte:head>
	<title>Settings</title>
</svelte:head>

<Dialog.Root bind:open>
	<Dialog.Content class="flex max-w-[700px] flex-col gap-4">
		<Dialog.Header>
			<Dialog.Title>Latest Backups</Dialog.Title>
			<Dialog.Description class="flex items-start justify-between gap-2">
				Click on a backup to download it or use the buttons to either quickly restore a backup or
				delete it.
				<Button variant="default" onclick={createBackup}>Create Backup</Button>
			</Dialog.Description>
		</Dialog.Header>

		<Separator />

		<Tabs.Root value="sqlite">
			<Tabs.List class="grid w-full grid-cols-2">
				<Tabs.Trigger value="sqlite" class="px-2 py-0.5 font-bold">SQLite</Tabs.Trigger>
				<Tabs.Trigger value="yaml" class="px-2 py-0.5 font-bold">YAML</Tabs.Trigger>
			</Tabs.List>
			<Tabs.Content value="sqlite" class="space-y-2">
				{#each sqliteBackups || [] as b (b.name)}
					<div class="flex items-center justify-between font-mono text-sm">
						<Button variant="link" class="flex items-center" onclick={() => downloadBackup(b.name)}>
							<HoverCard.Root openDelay={400}>
								<HoverCard.Trigger class="max-w-[250px] truncate">
									{b.name}
								</HoverCard.Trigger>
								<HoverCard.Content class="w-full">
									<div class="flex items-center">
										<CalendarDays class="mr-2 size-4 opacity-70" />
										<span class="text-xs text-muted-foreground">
											Created
											{#if b.createdAt}
												{DateFormat.format(timestampDate(b.createdAt))}
											{/if}
										</span>
									</div>
								</HoverCard.Content>
							</HoverCard.Root>
							<Download />
						</Button>
						<span class="flex items-center">
							<span class="mr-2">
								{Intl.NumberFormat('en-US', {
									notation: 'compact',
									style: 'unit',
									unit: 'byte'
								}).format(b.size)}
							</span>
							<Button
								variant="ghost"
								size="icon"
								class="rounded-full hover:bg-green-300/50 dark:hover:bg-green-700/50"
								onclick={() => {
									backupClient.restoreBackup({ name: b.name });
									open = false;
								}}
							>
								<DatabaseBackup />
							</Button>
							<Button
								variant="ghost"
								size="icon"
								class="rounded-full hover:bg-red-300/50 dark:hover:bg-red-700/50"
								onclick={() => deleteBackup(b.name)}
							>
								<Trash2 />
							</Button>
						</span>
					</div>
				{/each}
				{#if sqliteBackups.length === 0}
					<p class="pt-2 text-center text-sm text-muted-foreground">No backups available</p>
				{/if}
			</Tabs.Content>
			<Tabs.Content value="yaml" class="space-y-2">
				{#each yamlBackups || [] as b (b.name)}
					<div class="flex items-center justify-between font-mono text-sm">
						<Button variant="link" class="flex items-center" onclick={() => downloadBackup(b.name)}>
							<HoverCard.Root openDelay={400}>
								<HoverCard.Trigger class="max-w-[250px] truncate">
									{b.name}
								</HoverCard.Trigger>
								<HoverCard.Content class="w-full">
									<div class="flex items-center">
										<CalendarDays class="mr-2 size-4 opacity-70" />
										<span class="text-xs text-muted-foreground">
											Created
											{#if b.createdAt}
												{DateFormat.format(timestampDate(b.createdAt))}
											{/if}
										</span>
									</div>
								</HoverCard.Content>
							</HoverCard.Root>
							<Download />
						</Button>
						<span class="flex items-center">
							<span class="mr-2">
								{Intl.NumberFormat('en-US', {
									notation: 'compact',
									style: 'unit',
									unit: 'byte'
								}).format(b.size)}
							</span>
							<Button
								variant="ghost"
								size="icon"
								class="rounded-full hover:bg-green-300/50 dark:hover:bg-green-700/50"
								onclick={() => {
									backupClient.restoreBackup({ name: b.name });
									open = false;
								}}
							>
								<DatabaseBackup />
							</Button>
							<Button
								variant="ghost"
								size="icon"
								class="rounded-full hover:bg-red-300/50 dark:hover:bg-red-700/50"
								onclick={() => deleteBackup(b.name)}
							>
								<Trash2 />
							</Button>
						</span>
					</div>
				{/each}
				{#if yamlBackups.length === 0}
					<p class="pt-2 text-center text-sm text-muted-foreground">No backups available</p>
				{/if}
			</Tabs.Content>
		</Tabs.Root>
	</Dialog.Content>
</Dialog.Root>
