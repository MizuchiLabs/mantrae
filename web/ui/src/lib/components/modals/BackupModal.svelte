<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { Button } from '$lib/components/ui/button';
	import { Separator } from '$lib/components/ui/separator';
	import { CalendarDays, DatabaseBackup, Download, Trash2 } from '@lucide/svelte';
	import { formatTs } from '$lib/utils';
	import { backup } from '$lib/api/util.svelte';
	import { BackendURL } from '$lib/config';

	let { open = $bindable(false) } = $props();

	const backupList = $derived(backup.list());
	const createMutation = backup.create();
	const deleteMutation = backup.delete();
	const restoreMutation = backup.restore();

	let sqliteBackups = $derived(backupList.data?.filter((b) => b.name.endsWith('.db')) || []);
	let yamlBackups = $derived(backupList.data?.filter((b) => b.name.endsWith('.yaml')) || []);

	async function downloadBackup(filename = '') {
		const params = filename ? `?name=${encodeURIComponent(filename)}` : '';
		const response = await fetch(`${BackendURL}/backups/download${params}`, {
			credentials: 'include'
		});
		if (!response.ok) {
			throw new Error(await response.text());
		}

		const disposition = response.headers.get('Content-Disposition');
		const name = disposition?.match(/filename="?([^"]+)"?/i)?.[1];
		if (!name) {
			throw new Error('Server did not provide filename');
		}

		const blob = await response.blob();
		const url = URL.createObjectURL(blob);

		const a = document.createElement('a');
		a.href = url;
		a.download = name;
		a.click();

		URL.revokeObjectURL(url);
	}
</script>

<svelte:head>
	<title>Settings</title>
</svelte:head>

<Dialog.Root bind:open>
	<Dialog.Content class="flex max-w-150 flex-col gap-4">
		<Dialog.Header>
			<Dialog.Title>Latest Backups</Dialog.Title>
			<Dialog.Description class="flex items-start justify-between gap-2">
				Click on a backup to download it or use the buttons to either quickly restore a backup or
				delete it.
				<Button variant="default" onclick={() => createMutation.mutate({})}>Create Backup</Button>
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
											{formatTs(b.createdAt)}
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
									restoreMutation.mutate({ name: b.name });
									open = false;
								}}
							>
								<DatabaseBackup />
							</Button>
							<Button
								variant="ghost"
								size="icon"
								class="rounded-full hover:bg-red-300/50 dark:hover:bg-red-700/50"
								onclick={() => deleteMutation.mutate({ name: b.name })}
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
											{formatTs(b.createdAt)}
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
									restoreMutation.mutate({ name: b.name });
									open = false;
								}}
							>
								<DatabaseBackup />
							</Button>
							<Button
								variant="ghost"
								size="icon"
								class="rounded-full hover:bg-red-300/50 dark:hover:bg-red-700/50"
								onclick={() => deleteMutation.mutate({ name: b.name })}
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
