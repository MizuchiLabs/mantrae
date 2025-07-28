<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { Separator } from '$lib/components/ui/separator';
	import { Download, List, SaveIcon, Settings, Upload } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import PasswordInput from '$lib/components/ui/password-input/password-input.svelte';
	import { backupClient, settingClient, upload } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import { settingGroups, storageTypes } from './settings';
	import BackupModal from '$lib/components/modals/BackupModal.svelte';
	import CustomSwitch from '$lib/components/ui/custom-switch/custom-switch.svelte';
	import { profile } from '$lib/stores/profile';

	let settingsMap = $state<Record<string, string>>({});
	let originalSettings: Record<string, string> = {};
	let changedSettings = $state<Record<string, string>>({});

	// Helper function to check if an entire group should be visible
	function shouldShowGroup(groupKey: string): boolean {
		// Hide S3 group when storage is local
		if (groupKey === 's3' && settingsMap['storage_select'] === 'local') {
			return false;
		}
		return true;
	}

	// Helper function to check if a setting should be visible
	function shouldShowSetting(settingKey: string): boolean {
		// Hide client secret when PKCE is enabled
		if (settingKey === 'oidc_client_secret' && settingsMap['oidc_pkce'] === 'true') {
			return false;
		}

		// Hide OIDC-related settings when OIDC is disabled (except the enable toggle itself)
		if (
			(settingKey.startsWith('oidc_') || settingKey === 'password_login_enabled') &&
			settingKey !== 'oidc_enabled' &&
			settingsMap['oidc_enabled'] === 'false'
		) {
			return false;
		}

		return true;
	}

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
	let showBackupList = $state(false);
	let uploadBackupFile: HTMLInputElement | null = $state(null);

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
		const response = await settingClient.listSettings({});
		settingsMap = Object.fromEntries(response.settings.map((s) => [s.key, s.value]));
		originalSettings = { ...settingsMap };
		changedSettings = {};
	});
</script>

<svelte:head>
	<title>Settings</title>
</svelte:head>

<BackupModal bind:open={showBackupList} />

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
						accept=".db,.yaml,.yml,.json"
						class="hidden"
						bind:this={uploadBackupFile}
						onchange={() => upload(uploadBackupFile, `backup/${profile.id}`)}
					/>

					<div class="grid grid-cols-1 gap-2 sm:grid-cols-6">
						<Button onclick={() => downloadBackup()} variant="outline" class="px-1">
							<Download />
							Download Backup
						</Button>

						<Tooltip.Provider>
							<Tooltip.Root delayDuration={500}>
								<Tooltip.Trigger class="w-full">
									<Button
										variant="outline"
										onclick={() => uploadBackupFile?.click()}
										class="w-full px-1"
									>
										<Upload />
										Upload Backup
									</Button>
								</Tooltip.Trigger>
								<Tooltip.Content side="bottom" align="end" class="max-w-md">
									<p>
										Upload a SQLite <code>.db</code> file to fully replace the current database, or
										upload a Traefik dynamic config file in <code>.yaml</code> or
										<code>.json</code>
										format. Dynamic configs will be merged with the existing routers, services, and middlewares—overwriting
										any entries with matching names.
									</p>
								</Tooltip.Content>
							</Tooltip.Root>
						</Tooltip.Provider>

						<Button variant="outline" onclick={() => (showBackupList = true)} class="px-1">
							<List />
							View Backups
						</Button>
					</div>
				</div>
			</Card.Content>
		</Card.Root>

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
					{#if shouldShowGroup(groupKey)}
						<div class="mt-4 first:mt-0">
							<h2 class="mb-0.5 text-xl font-semibold" id={groupKey}>{group.title}</h2>
							<p class="text-muted-foreground mb-2 text-sm">{group.description}</p>
							<Separator class="mb-4" />

							<!-- Loop through settings in this group -->
							{#each group.keys as setting (setting.key)}
								{#if shouldShowSetting(setting.key)}
									<div
										class="mb-4 flex flex-col justify-start gap-4 sm:flex-row sm:justify-between"
									>
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
												<CustomSwitch
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
								{/if}
							{/each}
						</div>
					{/if}
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
