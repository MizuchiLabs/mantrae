<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
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
	} from '@lucide/svelte';
	import { settings, api, backups } from '$lib/api';
	import { onMount } from 'svelte';
	import type { Setting } from '$lib/types';
	import { toast } from 'svelte-sonner';
	import { DateFormat } from '$lib/stores/common';
	import PasswordInput from '$lib/components/ui/password-input/password-input.svelte';
	import { user } from '$lib/stores/user';

	let hasChanges = $state(false);
	let changedValues = $state<Record<string, Setting['value']>>({});

	// Define setting groups
	const settingGroups = {
		general: {
			title: 'General Settings',
			description: 'Basic application configuration',
			keys: [{ key: 'server_url', label: 'Server URL' }]
		},
		backup: {
			title: 'Backup Settings',
			description: 'Database backup configuration',
			keys: [
				{ key: 'backup_enabled', label: 'Enable Backups' },
				{ key: 'backup_interval', label: 'Backup Interval' },
				{ key: 'backup_keep', label: 'Number of Backups to Keep' },
				{ key: 'backup_path', label: 'Backup Path' },
				{ key: 'backup_storage_select', label: 'Backup Storage Type' }
			]
		},
		s3: {
			title: 'S3 Storage Settings',
			description: 'Amazon S3 or compatible storage configuration',
			keys: [
				{ key: 's3_endpoint', label: 'Endpoint' },
				{ key: 's3_bucket', label: 'Bucket Name' },
				{ key: 's3_region', label: 'Region' },
				{ key: 's3_access_key', label: 'Access Key' },
				{ key: 's3_secret_key', label: 'Secret Key' },
				{ key: 's3_use_path_style', label: 'Use Path Style' }
			]
		},
		email: {
			title: 'Email Settings',
			description: 'SMTP server configuration for sending emails',
			keys: [
				{ key: 'email_host', label: 'Host' },
				{ key: 'email_port', label: 'Port' },
				{ key: 'email_user', label: 'Username' },
				{ key: 'email_password', label: 'Password' },
				{ key: 'email_from', label: 'From Email Address' }
			]
		},
		oauth: {
			title: 'OIDC Settings',
			description: 'OIDC provider configuration',
			keys: [
				{ key: 'oidc_enabled', label: 'Enable OIDC' },
				{ key: 'oidc_client_id', label: 'Client ID' },
				{ key: 'oidc_client_secret', label: 'Client Secret' },
				{ key: 'oidc_issuer_url', label: 'Issuer URL' },
				{ key: 'oidc_provider_name', label: 'Provider Name' },
				{ key: 'oidc_scopes', label: 'Scopes' },
				{ key: 'oidc_pkce', label: 'Use PKCE' }
			]
		},
		agents: {
			title: 'Agent Settings',
			description: 'Agent management configuration',
			keys: [
				{ key: 'agent_cleanup_enabled', label: 'Enable Cleanup' },
				{ key: 'agent_cleanup_interval', label: 'Cleanup Interval' }
			]
		}
	};

	// Storage types for the select dropdown
	const storageTypes = [
		{ value: 'local', label: 'Local Storage' },
		{ value: 's3', label: 'S3 Storage' }
	];

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

	// Determine the input type based on the setting value
	const getInputType = (key: string, value: Setting['value']) => {
		if (typeof value === 'boolean') return 'boolean';
		if (key.toLowerCase().includes('select')) return 'select';
		if (key.toLowerCase().includes('password')) return 'password';
		if (key.toLowerCase().includes('secret')) return 'password';
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

	async function handleBackupList() {
		await api.listBackups();
		showBackupList = !showBackupList;
	}

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

						<Button variant="outline" onclick={handleBackupList}>
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
				<!-- Loop through each settings group -->
				{#each Object.entries(settingGroups) as [groupKey, group] (groupKey)}
					<div class="mt-4 first:mt-0">
						<h2 class="mb-0.5 text-xl font-semibold">{group.title}</h2>
						<p class="text-muted-foreground mb-2 text-sm">{group.description}</p>
						<Separator class="mb-4" />

						<!-- Loop through settings in this group -->
						{#each group.keys as { key, label } (key)}
							{#if $settings[key]}
								<div class="mb-4 flex flex-col justify-start gap-4 sm:flex-row sm:justify-between">
									<div class="border-muted-foreground border-l-2 pl-4">
										<Label>
											{label}
											{#if $settings[key].description}
												<p class="text-muted-foreground text-sm">{$settings[key].description}</p>
											{/if}
										</Label>
									</div>

									<div class="flex w-full items-center justify-end gap-4 sm:w-auto md:w-[380px]">
										{#if getInputType(key, $settings[key].value) === 'boolean'}
											<Switch
												id={key}
												checked={$settings[key].value as boolean}
												onCheckedChange={(checked) => saveSetting(key, checked)}
											/>
										{:else if getInputType(key, $settings[key].value) === 'select'}
											<Select.Root
												type="single"
												value={(changedValues[key] as string) || ($settings[key].value as string)}
												onValueChange={(value) => handleChange(key, value)}
											>
												<Select.Trigger>
													{changedValues[key] || $settings[key].value || 'Select...'}
												</Select.Trigger>
												<Select.Content>
													{#if key === 'backup_storage_select'}
														{#each storageTypes as option (option.value)}
															<Select.Item value={option.value}>{option.label}</Select.Item>
														{/each}
													{/if}
												</Select.Content>
											</Select.Root>
										{:else if getInputType(key, $settings[key].value) === 'password'}
											<PasswordInput
												class="sm:w-auto md:w-[380px]"
												value={changedValues[key] !== undefined
													? changedValues[key]
													: $settings[key].value}
												autocomplete="new-password"
												onchange={(e) => handleChange(key, e.currentTarget.value)}
												onkeydown={(e) => handleKeydown(e, key, e.currentTarget.value)}
											/>
										{:else if getInputType(key, $settings[key].value) === 'duration'}
											<Input
												type="text"
												id={key}
												value={changedValues[key] !== undefined
													? changedValues[key]
													: $settings[key].value}
												onchange={(e) => handleChange(key, parseDuration(e.currentTarget.value))}
												onkeydown={(e) =>
													handleKeydown(e, key, parseDuration(e.currentTarget.value))}
											/>
										{:else if getInputType(key, $settings[key].value) === 'port'}
											<Input
												type="number"
												id={key}
												value={changedValues[key] !== undefined
													? changedValues[key]
													: $settings[key].value}
												min="1"
												max="65535"
												onchange={(e) => handleChange(key, parseInt(e.currentTarget.value))}
												onkeydown={(e) => handleKeydown(e, key, parseInt(e.currentTarget.value))}
											/>
										{:else}
											<Input
												type="text"
												value={changedValues[key] !== undefined
													? changedValues[key]
													: $settings[key].value}
												autocomplete="off"
												onchange={(e) => handleChange(key, e.currentTarget.value)}
												onkeydown={(e) => handleKeydown(e, key, e.currentTarget.value)}
											/>
										{/if}
									</div>
								</div>

								<Separator class="mb-4" />
							{/if}
						{/each}
					</div>
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
