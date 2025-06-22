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
	import { backupClient, settingClient } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Backup } from '$lib/gen/mantrae/v1/backup_pb';
	import { uploadBackup } from '$lib/api';

	let settingsMap = $state<Record<string, string>>({});

	// Define setting groups
	const settingGroups = {
		general: {
			title: 'General Settings',
			description: 'Core configuration for connecting to backend services.',
			keys: [
				{
					key: 'server_url',
					label: 'Server URL',
					type: 'text',
					description:
						'The base URL of your backend server, including protocol (e.g., https://example.com).'
				}
			]
		},
		backup: {
			title: 'Backup Settings',
			description: 'Control how and where the system creates database backups.',
			keys: [
				{
					key: 'backup_enabled',
					label: 'Enable Backups',
					type: 'boolean',
					description: 'Toggle automatic backups of the database.'
				},
				{
					key: 'backup_interval',
					label: 'Backup Interval',
					type: 'duration',
					description: 'How often backups should run (e.g., 24h, 1h30m).'
				},
				{
					key: 'backup_keep',
					label: 'Backups to Keep',
					type: 'number',
					description: 'The number of recent backups to retain.'
				},
				{
					key: 'backup_path',
					label: 'Backup Path',
					type: 'text',
					description: 'Local filesystem path where backups will be stored.'
				},
				{
					key: 'backup_storage_select',
					label: 'Storage Type',
					type: 'select',
					description: 'Select the storage backend for backups (e.g., local, S3).'
				}
			]
		},
		s3: {
			title: 'S3 Storage Settings',
			description: 'Configure connection to Amazon S3 or any compatible object storage.',
			keys: [
				{
					key: 's3_endpoint',
					label: 'S3 Endpoint',
					type: 'text',
					description:
						'Custom endpoint for your S3-compatible storage (e.g., https://s3.amazonaws.com).'
				},
				{
					key: 's3_bucket',
					label: 'Bucket Name',
					type: 'text',
					description: 'The name of the S3 bucket used for storing backups.'
				},
				{
					key: 's3_region',
					label: 'Region',
					type: 'text',
					description: 'AWS region of your bucket (e.g., us-east-1).'
				},
				{
					key: 's3_access_key',
					label: 'Access Key',
					type: 'text',
					description: 'Access key ID for your S3 storage credentials.'
				},
				{
					key: 's3_secret_key',
					label: 'Secret Key',
					type: 'text',
					description: 'Secret access key for S3 storage. Keep it safe.'
				},
				{
					key: 's3_use_path_style',
					label: 'Use Path Style',
					type: 'boolean',
					description:
						'Enable if your S3 provider requires path-style URLs (instead of virtual-host style).'
				}
			]
		},
		email: {
			title: 'Email Settings',
			description: 'Set up SMTP for sending system and user notification emails.',
			keys: [
				{
					key: 'email_host',
					label: 'SMTP Host',
					type: 'text',
					description: 'Hostname of the SMTP server (e.g., smtp.mailgun.org).'
				},
				{
					key: 'email_port',
					label: 'SMTP Port',
					type: 'number',
					description: 'Port used to connect to the SMTP server (e.g., 587 or 465).'
				},
				{
					key: 'email_user',
					label: 'Username',
					type: 'text',
					description: 'Login username for the SMTP server.'
				},
				{
					key: 'email_password',
					label: 'Password',
					type: 'password',
					description: 'SMTP password or app-specific token.'
				},
				{
					key: 'email_from',
					label: 'From Address',
					type: 'text',
					description: 'Email address to use as the sender in outgoing emails.'
				}
			]
		},
		oauth: {
			title: 'OIDC Authentication',
			description:
				'Configure OpenID Connect for secure single sign-on (SSO) and identity management.',
			keys: [
				{
					key: 'oidc_enabled',
					label: 'Enable OIDC',
					type: 'boolean',
					description: 'Turn on OpenID Connect authentication.'
				},
				{
					key: 'password_login_disabled',
					label: 'Disable Password Login',
					type: 'boolean',
					description: 'Force users to log in only via OIDC (no local passwords).'
				},
				{
					key: 'oidc_client_id',
					label: 'Client ID',
					type: 'text',
					description: 'OIDC client ID issued by your identity provider.'
				},
				{
					key: 'oidc_client_secret',
					label: 'Client Secret',
					type: 'password',
					description: 'OIDC client secret issued by your identity provider.'
				},
				{
					key: 'oidc_issuer_url',
					label: 'Issuer URL',
					type: 'text',
					description: 'URL of the OIDC provider (e.g., https://auth.example.com).'
				},
				{
					key: 'oidc_provider_name',
					label: 'Provider Name',
					type: 'text',
					description: 'Friendly name shown in the login UI for the OIDC provider.'
				},
				{
					key: 'oidc_scopes',
					label: 'OIDC Scopes',
					type: 'text',
					description:
						'Requested scopes (space-separated) for the OIDC flow (e.g., openid email profile).'
				},
				{
					key: 'oidc_pkce',
					label: 'Use PKCE',
					type: 'boolean',
					description: 'Enable PKCE (Proof Key for Code Exchange) for better security.'
				}
			]
		},
		agents: {
			title: 'Agent Configuration',
			description: 'Manage automated cleanup tasks for connected agents.',
			keys: [
				{
					key: 'agent_cleanup_enabled',
					label: 'Enable Cleanup',
					type: 'boolean',
					description: 'Automatically remove stale or offline agents on a schedule.'
				},
				{
					key: 'agent_cleanup_interval',
					label: 'Cleanup Interval',
					type: 'duration',
					description: 'Frequency of cleanup jobs (e.g., 1h, 24h).'
				}
			]
		}
	};

	// Storage types for the select dropdown
	const storageTypes = [
		{ value: 'local', label: 'Local Storage' },
		{ value: 's3', label: 'S3 Storage' }
	];

	async function updateSetting(key: string, value: string) {
		settingsMap = { ...settingsMap, [key]: value };
	}
	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') saveSettings();
	}
	async function saveSettings() {
		try {
			await Promise.all(
				Object.entries(settingsMap).map(([key, value]) =>
					settingClient.updateSetting({ key, value })
				)
			);
			toast.success('Settings saved');
		} catch (err) {
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
	// let restoreDBFile: HTMLInputElement | null = $state(null);
	// let restoreDynamicFile: HTMLInputElement | null = $state(null);
	let backups = $state<Backup[]>([]);
	let showBackupList = $state(false);
	let uploadBackupFile: HTMLInputElement | null = $state(null);

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
		const response2 = await backupClient.listBackups({});
		backups = response2.backups;
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
			<div class="flex items-center justify-between gap-4">
				<input
					type="file"
					accept=".db"
					class="hidden"
					bind:this={uploadBackupFile}
					onchange={uploadBackup}
				/>
				<!---->
				<!-- <input -->
				<!-- 	type="file" -->
				<!-- 	accept=".yaml,.yml" -->
				<!-- 	class="hidden" -->
				<!-- 	bind:this={restoreDynamicFile} -->
				<!-- 	onchange={(e) => api.restoreDynamicConfig(e.currentTarget.files)} -->
				<!-- /> -->

				<div class="flex items-center gap-4">
					<Button onclick={() => downloadBackup()} variant="outline">
						<Download class="mr-2 size-4" />
						Download Backup
					</Button>

					<Button variant="outline" onclick={() => uploadBackupFile?.click()}>
						<Upload class="mr-2 size-4" />
						Upload Backup
					</Button>

					<Button variant="outline" onclick={() => (showBackupList = true)}>
						<List class="mr-2 size-4" />
						View Backups
					</Button>
				</div>
				<!-- <Tooltip.Provider> -->
				<!-- 	<Tooltip.Root delayDuration={100}> -->
				<!-- 		<Tooltip.Trigger> -->
				<!-- 			<Button variant="outline" onclick={() => restoreDynamicFile?.click()}> -->
				<!-- 				<FileCode class="mr-2 size-4" /> -->
				<!-- 				Import Configuration -->
				<!-- 			</Button> -->
				<!-- 		</Tooltip.Trigger> -->
				<!-- 		<Tooltip.Content side="bottom" align="end" class="w-80"> -->
				<!-- 			<p> -->
				<!-- 				Restore using the dynamic Traefik config in yaml or json format. It will merge -->
				<!-- 				current routers/middlewares with the provided dynamic config. -->
				<!-- 			</p> -->
				<!-- 		</Tooltip.Content> -->
				<!-- 	</Tooltip.Root> -->
				<!-- </Tooltip.Provider> -->
			</div>
		</Card.Content>
	</Card.Root>

	<Dialog.Root bind:open={showBackupList}>
		<Dialog.Content class="flex max-w-[600px] flex-col gap-4">
			<Dialog.Header>
				<Dialog.Title>Latest Backups</Dialog.Title>
				<Dialog.Description class="flex items-start justify-between gap-2">
					Click on a backup to download it or use the buttons to either quickly restore a backup or
					delete it.
					<Button variant="default" onclick={createBackup}>Create Backup</Button>
				</Dialog.Description>
			</Dialog.Header>
			<Separator />
			<div class="flex flex-col">
				{#each backups || [] as b (b.name)}
					<div class="flex items-center justify-between font-mono text-sm">
						<Button variant="link" class="flex items-center" onclick={() => downloadBackup(b.name)}>
							{#if b.createdAt}
								{DateFormat.format(timestampDate(b.createdAt))}
							{/if}
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
								class="rounded-full hover:bg-green-300"
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
								class="rounded-full hover:bg-red-300"
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
		<Card.Content class="flex  flex-col gap-6">
			<!-- Loop through each settings group -->
			{#each Object.entries(settingGroups) as [groupKey, group] (groupKey)}
				<div class="mt-4 first:mt-0">
					<h2 class="mb-0.5 text-xl font-semibold" id={groupKey}>{group.title}</h2>
					<p class="text-muted-foreground mb-2 text-sm">{group.description}</p>
					<Separator class="mb-4" />

					<!-- Loop through settings in this group -->
					{#each group.keys as setting (setting.key)}
						<div class="mb-4 flex flex-col justify-start gap-4 sm:flex-row sm:justify-between">
							<div class="border-muted-foreground flex items-center border-l-2 pl-4">
								<Label class="flex flex-col items-start justify-start gap-1">
									{setting.label}
									<p class="text-muted-foreground text-sm">{setting.description}</p>
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
										onkeydown={handleKeydown}
									/>
								{/if}
								{#if setting.type === 'number'}
									<Input
										type="number"
										value={settingsMap[setting.key]}
										autocomplete="off"
										onchange={(e) => updateSetting(setting.key, e.currentTarget.value)}
										onblur={saveSettings}
										onkeydown={handleKeydown}
									/>
								{/if}
								{#if setting.type === 'boolean'}
									<Switch
										id={setting.key}
										checked={settingsMap[setting.key] === 'true'}
										onCheckedChange={(checked) =>
											updateSetting(setting.key, checked ? 'true' : 'false')}
									/>
								{/if}
								{#if setting.type === 'password'}
									<PasswordInput
										class="sm:w-auto md:w-[380px]"
										value={settingsMap[setting.key]}
										autocomplete="new-password"
										onchange={(e) => updateSetting(setting.key, e.currentTarget.value)}
										onblur={saveSettings}
										onkeydown={handleKeydown}
									/>
								{/if}
								{#if setting.type === 'duration'}
									<Input
										type="text"
										value={settingsMap[setting.key]}
										autocomplete="off"
										onchange={(e) =>
											updateSetting(setting.key, parseDuration(e.currentTarget.value))}
										onblur={saveSettings}
										onkeydown={handleKeydown}
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
											{#if setting.key === 'backup_storage_select'}
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
				<Button variant="default" onclick={saveSettings} size="icon">
					<SaveIcon />
				</Button>
			</div>
		</Card.Content>
	</Card.Root>
</div>
