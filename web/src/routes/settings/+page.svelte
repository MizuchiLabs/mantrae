<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import { Input } from '$lib/components/ui/input';
	import { Separator } from '$lib/components/ui/separator';
	import {
		DatabaseBackup,
		Download,
		FileCode,
		List,
		SaveIcon,
		Settings,
		Trash2,
		Upload
	} from 'lucide-svelte';
	import { settings, api, backups } from '$lib/api';
	import { onMount } from 'svelte';
	import type { Setting } from '$lib/types';
	import { toast } from 'svelte-sonner';
	import { DateFormat } from '$lib/stores/common';
	import PasswordInput from '$lib/components/ui/password-input/password-input.svelte';
	import { user } from '$lib/stores/user';

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
	const getInputType = (key: string, value: Setting['value']) => {
		if (typeof value === 'boolean') return 'boolean';
		if (key.toLowerCase().includes('password')) return 'password';
		if (key.toLowerCase().includes('interval')) return 'duration';
		if (key.toLowerCase().includes('port')) return 'port';
		if (typeof value === 'number') return 'number';
		if (typeof value === 'string' && value.includes('://')) return 'url';
		if (typeof value === 'string' && value.includes('@')) return 'email';
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
	let restoreDBFile: HTMLInputElement | null = $state(null);
	let restoreDynamicFile: HTMLInputElement | null = $state(null);
	let showBackupList = $state(false);

	onMount(async () => {
		if (user.isAdmin) {
			await api.listSettings();
			await api.listBackups();
		}
	});
</script>

<svelte:head>
	<title>Settings</title>
</svelte:head>

{#if user.isAdmin}
	<div class="container flex flex-col gap-6">
		<Card.Root>
			<Card.Header>
				<Card.Title>Backup Management</Card.Title>
				<Card.Description>Download or restore database backups</Card.Description>
			</Card.Header>
			<Card.Content>
				<div class="flex items-center justify-between gap-4">
					<input
						type="file"
						accept=".db"
						class="hidden"
						bind:this={restoreDBFile}
						onchange={(e) => api.restoreBackup(e.currentTarget.files)}
					/>

					<input
						type="file"
						accept=".yaml,.yml"
						class="hidden"
						bind:this={restoreDynamicFile}
						onchange={(e) => api.restoreDynamicConfig(e.currentTarget.files)}
					/>

					<div class="flex items-center gap-4">
						<Button onclick={() => api.downloadBackup()} variant="outline">
							<Download class="mr-2 size-4" />
							Download Backup
						</Button>

						<Button variant="outline" onclick={() => restoreDBFile?.click()}>
							<Upload class="mr-2 size-4" />
							Upload Backup
						</Button>

						<Button variant="outline" onclick={() => (showBackupList = true)}>
							<List class="mr-2 size-4" />
							View Backups
						</Button>
					</div>
					<Tooltip.Provider>
						<Tooltip.Root delayDuration={100}>
							<Tooltip.Trigger>
								<Button variant="outline" onclick={() => restoreDynamicFile?.click()}>
									<FileCode class="mr-2 size-4" />
									Import Configuration
								</Button>
							</Tooltip.Trigger>
							<Tooltip.Content side="bottom" align="end" class="w-80">
								<p>
									Restore using the dynamic Traefik config in yaml or json format. It will merge
									current routers/middlewares with the provided dynamic config.
								</p>
							</Tooltip.Content>
						</Tooltip.Root>
					</Tooltip.Provider>
				</div>
			</Card.Content>
		</Card.Root>

		<Dialog.Root bind:open={showBackupList}>
			<Dialog.Content class="max-w-[600px]">
				<Dialog.Header>
					<Dialog.Title>Latest Backups</Dialog.Title>
					<Dialog.Description class="flex items-start justify-between gap-2">
						Click on a backup to download it or use the buttons to either quickly restore a backup
						or delete it.
						<Button variant="default" onclick={() => api.createBackup()}>Create Backup</Button>
					</Dialog.Description>
				</Dialog.Header>
				<Separator />
				<div class="flex flex-col">
					{#each $backups || [] as backup (backup.name)}
						<div class="flex items-center justify-between font-mono text-sm">
							<Button
								variant="link"
								onclick={() => api.downloadBackupByName(backup.name)}
								class="flex items-center"
							>
								{DateFormat.format(new Date(backup.timestamp))}
								<Download />
							</Button>
							<span class="flex items-center">
								<span class="mr-2">
									{Intl.NumberFormat('en-US', {
										notation: 'compact',
										style: 'unit',
										unit: 'byte'
									}).format(backup.size)}
								</span>
								<Button
									variant="ghost"
									size="icon"
									class="rounded-full hover:bg-green-300"
									onclick={() => {
										api.restoreBackupByName(backup.name);
										showBackupList = false;
									}}
								>
									<DatabaseBackup />
								</Button>
								<Button
									variant="ghost"
									size="icon"
									class="rounded-full hover:bg-red-300"
									onclick={() => api.deleteBackup(backup.name)}
								>
									<Trash2 />
								</Button>
							</span>
						</div>
					{/each}
					{#if !$backups || $backups.length === 0}
						<p class="text-muted-foreground text-center text-sm">No backups available</p>
					{/if}
				</div>
			</Dialog.Content>
		</Dialog.Root>

		<!-- Settings -->
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
				{#each Object.entries($settings) as [key, setting] (key)}
					<div class="flex flex-col justify-start gap-4 sm:flex-row sm:justify-between">
						<Label>
							{formatSettingName(key)}
							{#if setting.description}
								<p class="text-muted-foreground text-sm">{setting.description}</p>
							{/if}
						</Label>

						<div class="flex w-full items-center justify-end gap-4 sm:w-auto md:w-[380px]">
							{#if getInputType(key, setting.value) === 'boolean'}
								<Switch
									id={key}
									checked={setting.value as boolean}
									onCheckedChange={(checked) => saveSetting(key, checked)}
								/>
							{:else if getInputType(key, setting.value) === 'password'}
								<PasswordInput
									class="sm:w-auto md:w-[380px]"
									bind:value={setting.value}
									autocomplete="new-password"
									onchange={(e) => handleChange(key, e.currentTarget.value)}
									onkeydown={(e) => handleKeydown(e, key, e.currentTarget.value)}
								/>
							{:else if getInputType(key, setting.value) === 'duration'}
								<Input
									type="text"
									id={key}
									value={setting.value}
									onchange={(e) => handleChange(key, parseDuration(e.currentTarget.value))}
									onkeydown={(e) => handleKeydown(e, key, parseDuration(e.currentTarget.value))}
								/>
							{:else if getInputType(key, setting.value) === 'port'}
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
									autocomplete="off"
									onchange={(e) => handleChange(key, e.currentTarget.value)}
									onkeydown={(e) => handleKeydown(e, key, e.currentTarget.value)}
								/>
							{/if}
						</div>
					</div>

					<Separator />
				{/each}

				<div class="flex justify-end">
					<Button
						variant={hasChanges ? 'default' : 'outline'}
						onclick={saveAllChanges}
						disabled={!hasChanges}
						size="icon"
					>
						<SaveIcon />
					</Button>
				</div>
			</Card.Content>
		</Card.Root>
	</div>
{/if}
