<script lang="ts">
	import { setting } from '$lib/api/settings.svelte';
	import { backup } from '$lib/api/util.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import CustomSwitch from '$lib/components/ui/custom-switch/custom-switch.svelte';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import PasswordInput from '$lib/components/ui/password-input/password-input.svelte';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Separator } from '$lib/components/ui/separator';
	import * as Table from '$lib/components/ui/table';
	import * as Tabs from '$lib/components/ui/tabs';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { BackendURL } from '$lib/config';
	import { profileID } from '$lib/store.svelte';
	import { formatTs } from '$lib/utils';
	import { ConnectError } from '@connectrpc/connect';
	import { Download, Loader, RefreshCw, RotateCcw, Trash2, Upload } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';
	import { settingGroups, storageTypes } from './settings';

	// Settings Data
	const settings = $derived(setting.list());
	const updateSettings = setting.update();

	let settingsMap = $state<Record<string, string>>({});
	let initialized = $state(false);

	$effect(() => {
		if (settings.data && !initialized) {
			settingsMap = Object.fromEntries(settings.data.map((s) => [s.key, s.value]));
			initialized = true;
		}
	});

	// Backups Data
	const backupList = $derived(backup.list());
	const createBackupMutation = backup.create();
	const deleteBackupMutation = backup.delete();
	const restoreBackupMutation = backup.restore();

	let backups = $derived(backupList.data || []);
	let uploadInput: HTMLInputElement | null = $state(null);
	let isUploading = $state(false);

	// Visibility Helpers
	function shouldShowSetting(settingKey: string): boolean {
		if (settingKey === 'oidc_client_secret' && settingsMap['oidc_pkce'] === 'true') return false;
		if (
			(settingKey.startsWith('oidc_') || settingKey === 'password_login_enabled') &&
			settingKey !== 'oidc_enabled' &&
			settingsMap['oidc_enabled'] === 'false'
		) {
			return false;
		}
		return true;
	}

	// Settings Logic
	async function handleSave(key: string, value: string) {
		// Optimistic update
		settingsMap = { ...settingsMap, [key]: value };

		try {
			await updateSettings.mutateAsync({ key, value });
		} catch (err) {
			// Revert is handled by list refetch or could be done here if needed.
			// Since we use $state and init once, we might want to refetch on error to reset.
			const e = ConnectError.from(err);
			toast.error('Failed to save setting', { description: e.message });
		}
	}

	function handleKeydown(e: KeyboardEvent, key: string) {
		if (e.key === 'Enter') {
			handleSave(key, (e.currentTarget as HTMLInputElement).value);
		}
	}

	// Duration Parsing
	function parseDuration(str: string): string {
		const cleanStr = str.trim();
		try {
			const patterns = /^(\d+h)?(\d+m)?(\d+s)?$/;
			if (!patterns.test(cleanStr)) throw new Error('Invalid duration format');
			return cleanStr;
		} catch (err) {
			toast.error('Invalid duration format. Use format like "24h0m0s"');
			return str;
		}
	}

	// Backup Logic
	async function downloadBackup(filename: string) {
		try {
			const params = filename ? `?name=${encodeURIComponent(filename)}` : '';
			const response = await fetch(`${BackendURL}/backups/download${params}`, {
				credentials: 'include'
			});
			if (!response.ok) throw new Error(await response.text());

			const disposition = response.headers.get('Content-Disposition');
			const name = disposition?.match(/filename="?([^"]+)"?/i)?.[1] || filename;
			const blob = await response.blob();
			const url = URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = name;
			a.click();
			URL.revokeObjectURL(url);
		} catch (err) {
			toast.error('Download failed', { description: (err as Error).message });
		}
	}

	async function uploadBackup(input: HTMLInputElement | null) {
		if (!input?.files?.length) return;
		isUploading = true;
		try {
			const body = new FormData();
			body.append('file', input.files[0]);
			const response = await fetch(`${BackendURL}/backups/upload/${profileID.current}`, {
				method: 'POST',
				body,
				credentials: 'include'
			});
			if (!response.ok) throw new Error('Failed to upload');
			toast.success('Backup uploaded successfully');
			backupList.refetch(); // Refresh list
		} catch (err) {
			toast.error('Upload failed', { description: (err as Error).message });
		} finally {
			isUploading = false;
			if (input) input.value = '';
		}
	}

	// Group Rendering Component
	function getGroup(key: string) {
		return Object.entries(settingGroups).find(([k]) => k === key)?.[1];
	}
</script>

<svelte:head>
	<title>Settings - Mantrae</title>
</svelte:head>

<div class="container mx-auto max-w-7xl">
	<div class="mb-8 flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold tracking-tight">Settings</h1>
			<p class="text-muted-foreground">Manage your application configuration and backups.</p>
		</div>
	</div>

	<Tabs.Root value="general" class="space-y-6">
		<Tabs.List class="w-full">
			<Tabs.Trigger value="general">General</Tabs.Trigger>
			<Tabs.Trigger value="backups">Backups</Tabs.Trigger>
			{#if settingsMap['storage_select'] === 's3'}
				<Tabs.Trigger value="s3">S3 Storage</Tabs.Trigger>
			{/if}
			<Tabs.Trigger value="auth">Auth</Tabs.Trigger>
			<Tabs.Trigger value="email">Email</Tabs.Trigger>
			<Tabs.Trigger value="agents">Agents</Tabs.Trigger>
		</Tabs.List>

		<!-- General Tab -->
		<Tabs.Content value="general">
			<Card.Root>
				<Card.Header>
					<Card.Title>General Configuration</Card.Title>
					<Card.Description>Core system settings.</Card.Description>
				</Card.Header>
				<Card.Content class="space-y-6">
					{@render settingsGroup('general')}
				</Card.Content>
			</Card.Root>
		</Tabs.Content>

		<!-- Backups Tab -->
		<Tabs.Content value="backups" class="space-y-6">
			<!-- Backup Configuration -->
			<Card.Root>
				<Card.Header>
					<Card.Title>Configuration</Card.Title>
					<Card.Description>Configure how automatic backups are handled.</Card.Description>
				</Card.Header>
				<Card.Content class="space-y-6">
					{@render settingsGroup('backup')}
				</Card.Content>
			</Card.Root>

			<!-- Backup Management -->
			<Card.Root>
				<Card.Header>
					<div class="flex items-center justify-between">
						<div>
							<Card.Title>Manage Backups</Card.Title>
							<Card.Description>Create, download, or restore from backups.</Card.Description>
						</div>
						<div class="flex gap-2">
							<input
								type="file"
								accept=".db,.yaml,.yml,.json"
								class="hidden"
								bind:this={uploadInput}
								onchange={() => uploadBackup(uploadInput)}
							/>
							<Button variant="outline" onclick={() => uploadInput?.click()} disabled={isUploading}>
								{#if isUploading}
									<Loader class="mr-2 size-4 animate-spin" />
								{:else}
									<Upload class="mr-2 size-4" />
								{/if}
								Upload
							</Button>
							<Button onclick={() => createBackupMutation.mutate({})}>
								<RefreshCw class="mr-2 size-4" />
								Create Backup
							</Button>
						</div>
					</div>
				</Card.Header>
				<Card.Content>
					<Table.Root>
						<Table.Header>
							<Table.Row>
								<Table.Head>Name</Table.Head>
								<Table.Head>Type</Table.Head>
								<Table.Head>Size</Table.Head>
								<Table.Head>Created</Table.Head>
								<Table.Head class="text-right">Actions</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each backups as backup (backup.name)}
								<Table.Row>
									<Table.Cell class="font-medium">{backup.name}</Table.Cell>
									<Table.Cell>
										{#if backup.name.endsWith('.db')}
											<span
												class="inline-flex items-center rounded-full bg-blue-100 px-2.5 py-0.5 text-xs font-medium text-blue-800 dark:bg-blue-900 dark:text-blue-300"
											>
												SQLite
											</span>
										{:else}
											<span
												class="inline-flex items-center rounded-full bg-yellow-100 px-2.5 py-0.5 text-xs font-medium text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300"
											>
												YAML
											</span>
										{/if}
									</Table.Cell>
									<Table.Cell>
										{Intl.NumberFormat('en-US', {
											notation: 'compact',
											style: 'unit',
											unit: 'byte'
										}).format(backup.size)}
									</Table.Cell>
									<Table.Cell>{formatTs(backup.createdAt)}</Table.Cell>
									<Table.Cell class="text-right">
										<div class="flex justify-end gap-2">
											<Tooltip.Provider>
												<Tooltip.Root>
													<Tooltip.Trigger>
														<Button
															variant="ghost"
															size="icon"
															onclick={() => downloadBackup(backup.name)}
														>
															<Download class="size-4" />
														</Button>
													</Tooltip.Trigger>
													<Tooltip.Content>Download</Tooltip.Content>
												</Tooltip.Root>

												<Tooltip.Root>
													<Tooltip.Trigger>
														<Button
															variant="ghost"
															size="icon"
															class="text-green-600 hover:bg-green-50 hover:text-green-700"
															onclick={() => restoreBackupMutation.mutate({ name: backup.name })}
														>
															<RotateCcw class="size-4" />
														</Button>
													</Tooltip.Trigger>
													<Tooltip.Content>Restore</Tooltip.Content>
												</Tooltip.Root>

												<Tooltip.Root>
													<Tooltip.Trigger>
														<Button
															variant="ghost"
															size="icon"
															class="text-red-600 hover:bg-red-50 hover:text-red-700"
															onclick={() => deleteBackupMutation.mutate({ name: backup.name })}
														>
															<Trash2 class="size-4" />
														</Button>
													</Tooltip.Trigger>
													<Tooltip.Content>Delete</Tooltip.Content>
												</Tooltip.Root>
											</Tooltip.Provider>
										</div>
									</Table.Cell>
								</Table.Row>
							{:else}
								<Table.Row>
									<Table.Cell colspan={5} class="text-center text-muted-foreground">
										No backups found.
									</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</Card.Content>
			</Card.Root>
		</Tabs.Content>

		<!-- S3 Tab -->
		<Tabs.Content value="s3">
			<Card.Root>
				<Card.Header>
					<Card.Title>S3 Storage</Card.Title>
					<Card.Description>Configure external object storage.</Card.Description>
				</Card.Header>
				<Card.Content class="space-y-6">
					{@render settingsGroup('s3')}
				</Card.Content>
			</Card.Root>
		</Tabs.Content>

		<!-- Auth Tab -->
		<Tabs.Content value="auth">
			<Card.Root>
				<Card.Header>
					<Card.Title>Authentication</Card.Title>
					<Card.Description>
						Manage OIDC and login methods. Callback endpoint is <span class="font-mono">
							/oidc/callback
						</span>.
					</Card.Description>
				</Card.Header>
				<Card.Content class="space-y-6">
					{@render settingsGroup('oauth')}
				</Card.Content>
			</Card.Root>
		</Tabs.Content>

		<!-- Email Tab -->
		<Tabs.Content value="email">
			<Card.Root>
				<Card.Header>
					<Card.Title>Email Settings</Card.Title>
					<Card.Description>SMTP configuration for notifications.</Card.Description>
				</Card.Header>
				<Card.Content class="space-y-6">
					{@render settingsGroup('email')}
				</Card.Content>
			</Card.Root>
		</Tabs.Content>

		<!-- Agents Tab -->
		<Tabs.Content value="agents">
			<Card.Root>
				<Card.Header>
					<Card.Title>Agent Configuration</Card.Title>
					<Card.Description>Manage connected agents.</Card.Description>
				</Card.Header>
				<Card.Content class="space-y-6">
					{@render settingsGroup('agents')}
				</Card.Content>
			</Card.Root>
		</Tabs.Content>
	</Tabs.Root>
</div>

{#snippet settingsGroup(groupKey: string)}
	{#if getGroup(groupKey)}
		{#each getGroup(groupKey)!.keys as setting (setting.key)}
			{#if shouldShowSetting(setting.key)}
				{#if setting.type === 'boolean'}
					<div class="flex flex-row items-center justify-between rounded-lg border p-4 shadow-sm">
						<div class="space-y-0.5">
							<Label class="text-base">{setting.label}</Label>
							<p class="text-sm text-muted-foreground">{setting.description}</p>
						</div>
						<CustomSwitch
							checked={settingsMap[setting.key] === 'true'}
							onCheckedChange={(checked) => handleSave(setting.key, checked ? 'true' : 'false')}
						/>
					</div>
				{:else}
					<Separator />
					<div class="grid gap-4 sm:grid-cols-12 sm:gap-6">
						<div class="sm:col-span-5 md:col-span-7">
							<Label class="text-base">{setting.label}</Label>
							<p class="text-sm text-muted-foreground">{setting.description}</p>
						</div>
						<div class="sm:col-span-7 md:col-span-5">
							{#if setting.type === 'text'}
								<Input
									type="text"
									value={settingsMap[setting.key]}
									autocomplete="off"
									onchange={(e) => handleSave(setting.key, e.currentTarget.value)}
									onkeydown={(e) => handleKeydown(e, setting.key)}
								/>
							{:else if setting.type === 'number'}
								<Input
									type="number"
									value={settingsMap[setting.key]}
									autocomplete="off"
									onchange={(e) => handleSave(setting.key, e.currentTarget.value)}
									onkeydown={(e) => handleKeydown(e, setting.key)}
								/>
							{:else if setting.type === 'password'}
								<PasswordInput
									value={settingsMap[setting.key]}
									autocomplete="new-password"
									onchange={(e) => handleSave(setting.key, e.currentTarget.value)}
									onkeydown={(e) => handleKeydown(e, setting.key)}
								/>
							{:else if setting.type === 'duration'}
								<Input
									type="text"
									value={settingsMap[setting.key]}
									autocomplete="off"
									onchange={(e) => handleSave(setting.key, parseDuration(e.currentTarget.value))}
									onkeydown={(e) => handleKeydown(e, setting.key)}
								/>
							{:else if setting.type === 'select'}
								<Select.Root
									type="single"
									value={settingsMap[setting.key]}
									onValueChange={(value) => handleSave(setting.key, value)}
								>
									<Select.Trigger class="w-full">
										{settingsMap[setting.key] === 'local'
											? 'Local Storage'
											: settingsMap[setting.key] === 's3'
												? 'S3 Storage'
												: settingsMap[setting.key] || 'Select...'}
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
				{/if}
			{/if}
		{/each}
	{/if}
{/snippet}
