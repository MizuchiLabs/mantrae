<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import { Input } from '$lib/components/ui/input';
	import { Separator } from '$lib/components/ui/separator';
	import { Download, List, SaveIcon, Settings, Trash, Upload } from 'lucide-svelte';
	import { settings, api, backups, loading } from '$lib/api';
	import { onMount } from 'svelte';
	import type { Setting } from '$lib/types';
	import { toast } from 'svelte-sonner';
	import { DateFormat } from '$lib/stores/common';

	let hasChanges = $state(false);
	let changedValues = $state<Record<string, Setting['value']>>({});

	function parseDuration(str: string): string {
		// Just clean up and validate the duration string
		const cleanStr = str.trim();
		try {
			// Validate the duration string format
			const patterns = /^(\d+h)?(\d+m)?(\d+s)?$/;
			if (!patterns.test(cleanStr)) {
				throw new Error('Invalid duration format');
			}
			return cleanStr;
		} catch (err) {
			const error = err as Error;
			toast.error('Invalid duration format. Use format like "24h0m0s"', {
				description: error.message
			});
			return str;
		}
	}

	// Helper to convert camelCase/snake_case to Title Case
	const formatSettingName = (key: string) => {
		return key
			.split(/[_\s]/)
			.map((word) => word.charAt(0).toUpperCase() + word.slice(1))
			.join(' ');
	};

	// Determine the input type based on the setting value
	const getInputType = (value: Setting) => {
		if (typeof value === 'boolean') return 'boolean';
		if (typeof value === 'number') return 'number';
		if (value?.toString().includes('://')) return 'url';
		if (value?.toString().includes('@')) return 'email';
		return 'text';
	};

	function handleChange(key: string, value: Setting['value']) {
		changedValues[key] = value;
		hasChanges = true;
	}

	async function saveSetting(key: string, value: Setting['value']) {
		try {
			await api.upsertSetting({
				key,
				value: value.toString(),
				description: $settings[key].description
			});
			delete changedValues[key];
			hasChanges = Object.keys(changedValues).length > 0;
			toast.success('Setting updated successfully');
		} catch (error) {
			toast.error('Failed to save setting', { description: (error as Error).message });
		}
	}

	async function saveAllChanges() {
		for (const [key, value] of Object.entries(changedValues)) {
			await saveSetting(key, value);
		}
		hasChanges = false;
	}

	function handleKeydown(e: KeyboardEvent, key: string, value: Setting['value']) {
		if (e.key === 'Enter') saveSetting(key, value);
	}

	// Backup handling
	let fileInput: HTMLInputElement;
	let showBackupList = $state(false);

	function humanFileSize(size: number) {
		var i = size == 0 ? 0 : Math.floor(Math.log(size) / Math.log(1024));
		return +(size / Math.pow(1024, i)).toFixed(2) * 1 + ' ' + ['B', 'kB', 'MB', 'GB', 'TB'][i];
	}

	onMount(async () => {
		await api.listSettings();
		await api.listBackups();
	});
</script>

<svelte:head>
	<title>Settings</title>
</svelte:head>

<div class="container flex flex-col gap-6">
	<Card.Root>
		<Card.Header>
			<Card.Title>Backup Management</Card.Title>
			<Card.Description>Download or restore database backups</Card.Description>
		</Card.Header>
		<Card.Content>
			<div class="flex items-center gap-4">
				<Button onclick={() => api.downloadBackup()} variant="outline">
					<Download class="mr-2 size-4" />
					Download Backup
				</Button>

				<Button variant="outline" onclick={() => fileInput?.click()} disabled={$loading}>
					<Upload class="mr-2 size-4" />
					{$loading ? 'Uploading...' : 'Restore Backup'}
				</Button>

				<input
					type="file"
					accept=".sql,.db"
					class="hidden"
					bind:this={fileInput}
					onchange={(e) => api.restoreBackup(e.currentTarget.files)}
				/>

				<Button variant="outline" onclick={() => (showBackupList = true)}>
					<List class="mr-2 size-4" />
					View Backups
				</Button>
			</div>
		</Card.Content>
	</Card.Root>
	<Card.Root>
		<Card.Header>
			<Card.Title class="mb-3">
				<div class="flex items-center gap-2">
					<Settings class="size-8" />
					<h1 class="text-3xl font-bold">Settings</h1>
				</div>
			</Card.Title>
			<Separator />
		</Card.Header>
		<Card.Content class="flex  flex-col gap-6">
			{#each Object.entries($settings) as [key, setting]}
				<div class="flex flex-col justify-start gap-4 sm:flex-row sm:justify-between">
					<Label>
						{formatSettingName(key)}
						{#if setting.description}
							<p class="text-sm text-muted-foreground">{setting.description}</p>
						{/if}
					</Label>

					<div class="flex w-full items-center justify-end gap-4 sm:w-auto md:w-[380px]">
						{#if getInputType(setting.value) === 'boolean'}
							<Switch
								id={key}
								checked={setting.value}
								onCheckedChange={(checked) => saveSetting(key, checked)}
							/>
						{:else if key.includes('interval')}
							<Input
								type="text"
								id={key}
								value={setting.value}
								onchange={(e) => handleChange(key, parseDuration(e.currentTarget.value))}
								onkeydown={(e) => handleKeydown(e, key, parseDuration(e.currentTarget.value))}
							/>
						{:else if key.includes('port')}
							<Input
								type="number"
								id={key}
								value={setting.value}
								min="1"
								max="65535"
								onchange={(e) => handleChange(key, parseInt(e.currentTarget.value))}
								onkeydown={(e) => handleKeydown(e, key, parseInt(e.currentTarget.value))}
							/>
						{:else}
							<Input
								type="text"
								value={setting.value}
								onchange={(e) => handleChange(key, e.currentTarget.value)}
								onkeydown={(e) => handleKeydown(e, key, e.currentTarget.value)}
							/>
						{/if}
					</div>
				</div>

				<Separator />
			{/each}

			<div class="flex justify-end">
				<Button variant={hasChanges ? 'default' : 'outline'} onclick={saveAllChanges} size="icon">
					<SaveIcon />
				</Button>
			</div>
		</Card.Content>
	</Card.Root>
</div>

<Dialog.Root bind:open={showBackupList}>
	<Dialog.Content class="max-w-[600px]">
		<Dialog.Header class="mb-4">
			<Dialog.Title>Available Backups</Dialog.Title>
			<Dialog.Description class="flex items-start justify-between gap-2">
				Click on a backup to download it or click the trash icon to delete it
				<Button variant="default" onclick={() => api.createBackup()}>Create Backup</Button>
			</Dialog.Description>
		</Dialog.Header>
		<div class="flex flex-col gap-2">
			{#each $backups || [] as backup}
				<Button
					variant="link"
					class="flex items-center justify-between p-3"
					onclick={() => api.downloadBackupByName(backup.name)}
				>
					<span class="font-mono text-sm">
						Backup:
						{DateFormat.format(new Date(backup.timestamp))}
					</span>
					<span class="flex items-center font-mono text-sm">
						{humanFileSize(backup.size)}
						<button
							class="ml-2 rounded-full p-2 hover:bg-red-300"
							onclick={() => api.deleteBackup(backup.name)}
						>
							<Trash />
						</button>
					</span>
				</Button>
			{/each}
			{#if !$backups || $backups.length === 0}
				<p class="text-center text-sm text-muted-foreground">No backups available</p>
			{/if}
		</div>
	</Dialog.Content>
</Dialog.Root>
