<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import { Input } from '$lib/components/ui/input';
	import { Separator } from '$lib/components/ui/separator';
	import {
		Braces,
		CalendarDays,
		DatabaseBackup,
		Download,
		List,
		SaveIcon,
		Settings,
		Trash2,
		Upload
	} from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { DateFormat } from '$lib/stores/common';
	import PasswordInput from '$lib/components/ui/password-input/password-input.svelte';
	import { backupClient, settingClient, upload } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Backup } from '$lib/gen/mantrae/v1/backup_pb';
	import { profile } from '$lib/stores/profile';
	import { settingGroups, storageTypes } from './settings';

	let settingsMap = $state<Record<string, string>>({});
	let originalSettings: Record<string, string> = {};
	let changedSettings = $state<Record<string, string>>({});

	function handleKeydown(e: KeyboardEvent, key: string) {
		if (e.key === 'Enter') {
			let input = e.currentTarget as HTMLInputElement;
			updateSetting(key, input.value);
			saveSettings();
		}
	}

	async function updateSetting(key: string, value: string) {
		settingsMap = { ...settingsMap, [key]: value };

		if (value !== originalSettings[key]) {
			changedSettings = { ...changedSettings, [key]: value };
		} else {
			// remove it if it was reverted to original
			const { [key]: _, ...rest } = changedSettings;
			changedSettings = rest;
		}
	}

	async function saveSettings() {
		if (Object.keys(changedSettings).length === 0) return;

		try {
			await Promise.all(
				Object.entries(changedSettings).map(([key, value]) =>
					settingClient.updateSetting({ key, value })
				)
			);

			// merge deltas into original and clear
			originalSettings = { ...originalSettings, ...changedSettings };
			changedSettings = {};
			toast.success('Settings saved');
		} catch (err) {
			// Revert UI changes back to original values
			const revertedSettings = { ...settingsMap };
			Object.keys(changedSettings).forEach((key) => {
				revertedSettings[key] = originalSettings[key] || '';
			});
			settingsMap = revertedSettings;

			// Clear changed settings since we reverted
			changedSettings = {};

			const e = ConnectError.from(err);
			toast.error('Failed to save setting', { description: e.message });
		}
	}

	// Validate the go duration string format
	function parseDuration(str: string): string {
		const cleanStr = str.trim();
		try {
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

	// Backup handling
	let backups = $state<Backup[]>([]);
	let showBackupList = $state(false);
	let uploadBackupFile: HTMLInputElement | null = $state(null);
	let uploadDynamicFile: HTMLInputElement | null = $state(null);

	async function deleteBackup(name: string) {
		await backupClient.deleteBackup({ name });
		const response = await backupClient.listBackups({});
		backups = response.backups;
	}
	async function createBackup() {
		await backupClient.createBackup({});
		const response = await backupClient.listBackups({});
		backups = response.backups;
	}
	async function downloadBackup(name?: string) {
		if (!name) name = backups[0].name;
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
	}

	onMount(async () => {
		const response = await settingClient.listSettings({});
		settingsMap = Object.fromEntries(response.settings.map((s) => [s.key, s.value]));
		originalSettings = { ...settingsMap };
		changedSettings = {};

		const response2 = await backupClient.listBackups({});
		backups = response2.backups;
	});
</script>

<svelte:head>
	<title>Settings</title>
</svelte:head>

<div class="flex min-h-full w-full items-start justify-center">
	<div class="w-full max-w-[100rem] space-y-6">
		<Card.Root>
			<Card.Header>
				<Card.Title>Backup Management</Card.Title>
				<Card.Description>Download or restore database backups</Card.Description>
			</Card.Header>
			<Card.Content>
				<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
					<input
						type="file"
						accept=".db"
						class="hidden"
						bind:this={uploadBackupFile}
						onchange={() => upload(uploadBackupFile, 'backup')}
					/>
					<input
						type="file"
						accept=".yaml,.yml,.json"
						class="hidden"
						bind:this={uploadDynamicFile}
						onchange={() => upload(uploadDynamicFile, `backup?profile_id=${profile.id}`)}
					/>

					<div class="flex flex-wrap items-center gap-2 sm:gap-4">
						<Button onclick={() => downloadBackup()} variant="outline" class="flex-1 sm:flex-none">
							<Download class="mr-2 size-4" />
							Download Backup
						</Button>

						<Button
							variant="outline"
							onclick={() => uploadBackupFile?.click()}
							class="flex-1 sm:flex-none"
						>
							<Upload class="mr-2 size-4" />
							Upload Backup
						</Button>

						<Button
							variant="outline"
							onclick={() => (showBackupList = true)}
							class="flex-1 sm:flex-none"
						>
							<List class="mr-2 size-4" />
							View Backups
						</Button>

						<Tooltip.Provider>
							<Tooltip.Root delayDuration={100}>
								<Tooltip.Trigger class="flex-1 sm:flex-none">
									<Button
										variant="outline"
										onclick={() => uploadDynamicFile?.click()}
										class="flex-1 sm:flex-none"
									>
										<Braces class="mr-2 size-4" />
										Upload Configuration
									</Button>
								</Tooltip.Trigger>
								<Tooltip.Content side="bottom" align="end" class="max-w-md">
									<p>
										Restore using the dynamic Traefik config in yaml or json format. It will merge
										current routers/middlewares with the provided dynamic config.
									</p>
								</Tooltip.Content>
							</Tooltip.Root>
						</Tooltip.Provider>
					</div>
				</div>
			</Card.Content>
		</Card.Root>

		<Dialog.Root bind:open={showBackupList}>
			<Dialog.Content class="flex max-w-[600px] flex-col gap-4">
				<Dialog.Header>
					<Dialog.Title>Latest Backups</Dialog.Title>
					<Dialog.Description class="flex items-start justify-between gap-2">
						Click on a backup to download it or use the buttons to either quickly restore a backup
						or delete it.
						<Button variant="default" onclick={createBackup}>Create Backup</Button>
					</Dialog.Description>
				</Dialog.Header>
				<Separator />
				<div class="flex flex-col">
					{#each backups || [] as b (b.name)}
						<div class="flex items-center justify-between font-mono text-sm">
							<Button
								variant="link"
								class="flex items-center"
								onclick={() => downloadBackup(b.name)}
							>
								<HoverCard.Root openDelay={400}>
									<HoverCard.Trigger>
										{b.name}
									</HoverCard.Trigger>
									<HoverCard.Content class="w-full">
										<div class="flex items-center">
											<CalendarDays class="mr-2 size-4 opacity-70" />
											<span class="text-muted-foreground text-xs">
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
										showBackupList = false;
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
					{#if backups.length === 0}
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
			<Card.Content class="flex flex-col gap-6">
				<!-- Loop through each settings group -->
				{#each Object.entries(settingGroups) as [groupKey, group] (groupKey)}
					<div class="mt-4 first:mt-0">
						<h2 class="mb-0.5 text-xl font-semibold" id={groupKey}>{group.title}</h2>
						<p class="text-muted-foreground mb-2 text-sm">{group.description}</p>
						<Separator class="mb-4" />

						<!-- Loop through settings in this group -->
						{#each group.keys as setting (setting.key)}
							<div class="mb-4 flex flex-col justify-start gap-4 sm:flex-row sm:justify-between">
								<div class="border-muted-foreground border-l-2 pl-4">
									<Label class="flex flex-col items-start gap-1">
										<span class="text-sm font-medium">{setting.label}</span>
										<p class="text-muted-foreground text-xs">{setting.description}</p>
									</Label>
								</div>

								<div class="flex w-full items-center justify-end gap-4 sm:w-auto md:w-[380px]">
									{#if setting.type === 'text'}
										<Input
											type="text"
											value={settingsMap[setting.key]}
											autocomplete="off"
											onchange={(e) => updateSetting(setting.key, e.currentTarget.value)}
											onblur={saveSettings}
											onkeydown={(e) => handleKeydown(e, setting.key)}
										/>
									{/if}
									{#if setting.type === 'number'}
										<Input
											type="number"
											value={settingsMap[setting.key]}
											autocomplete="off"
											onchange={(e) => updateSetting(setting.key, e.currentTarget.value)}
											onblur={saveSettings}
											onkeydown={(e) => handleKeydown(e, setting.key)}
										/>
									{/if}
									{#if setting.type === 'boolean'}
										<Switch
											id={setting.key}
											checked={settingsMap[setting.key] === 'true'}
											onCheckedChange={(checked) => {
												updateSetting(setting.key, checked ? 'true' : 'false');
												saveSettings();
											}}
										/>
									{/if}
									{#if setting.type === 'password'}
										<PasswordInput
											class="sm:w-auto md:w-[380px]"
											value={settingsMap[setting.key]}
											autocomplete="new-password"
											onblur={saveSettings}
											onchange={(e) => updateSetting(setting.key, e.currentTarget.value)}
											onkeydown={(e) => handleKeydown(e, setting.key)}
										/>
									{/if}
									{#if setting.type === 'duration'}
										<Input
											type="text"
											value={settingsMap[setting.key]}
											autocomplete="off"
											onblur={saveSettings}
											onchange={(e) =>
												updateSetting(setting.key, parseDuration(e.currentTarget.value))}
											onkeydown={(e) => handleKeydown(e, setting.key)}
										/>
									{/if}
									{#if setting.type === 'select'}
										<Select.Root
											type="single"
											value={settingsMap[setting.key]}
											onValueChange={(value) => {
												updateSetting(setting.key, value);
												saveSettings();
											}}
										>
											<Select.Trigger class="w-full">
												{settingsMap[setting.key] || 'Select...'}
											</Select.Trigger>
											<Select.Content>
												{#if setting.key === 'storage_select'}
													{#each storageTypes as option (option.value)}
														<Select.Item value={option.value}>{option.label}</Select.Item>
													{/each}
												{/if}
											</Select.Content>
										</Select.Root>
									{/if}
								</div>
							</div>
							<Separator class="mb-4" />
						{/each}
					</div>
				{/each}

				<div class="flex justify-end">
					<Button
						variant={Object.keys(changedSettings).length === 0 ? 'outline' : 'default'}
						onclick={saveSettings}
						size="icon"
						disabled={Object.keys(changedSettings).length === 0}
					>
						<SaveIcon />
					</Button>
				</div>
			</Card.Content>
		</Card.Root>
	</div>
</div>
