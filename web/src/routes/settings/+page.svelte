<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch';
	import { downloadBackup, getSettings, uploadBackup, settings, updateSetting } from '$lib/api';
	import { onMount } from 'svelte';
	import { Input } from '$lib/components/ui/input';
	import type { Setting } from '$lib/types/base';
	import HoverInfo from '$lib/components/utils/hoverInfo.svelte';
	import { Download, Eye, EyeOff, Save, Upload } from 'lucide-svelte';

	let fileInput: HTMLInputElement;
	const handleFileUpload = (event: Event) => {
		const file = (event.target as HTMLInputElement).files?.[0];
		if (file) {
			uploadBackup(file);
		}
		fileInput.value = '';
	};

	// Settings
	let agentCleanupEnabled: boolean;
	let backupEnabled: boolean;

	let settingsMap: Record<string, string> = {};
	let changedSettings: Record<string, string> = {};
	let showEmailPassword = false;

	const update = async (key: string) => {
		const setting = { key, value: settingsMap[key] };
		await updateSetting(setting);
	};
	const markAsChanged = (key: string) => {
		const changed = $settings.find((s: Setting) => s.key === key)?.value;
		if (settingsMap[key] !== changed) {
			changedSettings = { ...changedSettings, [key]: settingsMap[key] };
		} else {
			// Remove it from changedSettings if the value hasn't changed
			const { [key]: _, ...rest } = changedSettings;
			changedSettings = rest;
		}
	};
	const saveChanges = async (): Promise<void> => {
		for (const key in changedSettings) {
			if (changedSettings.hasOwnProperty(key)) {
				await update(key);
			}
		}

		// Reset
		changedSettings = {};
	};

	const onKeydown = (e: KeyboardEvent, key: string) => {
		if (e.key === 'Enter') {
			update(key);
		}
	};

	onMount(async () => {
		await getSettings();
		settingsMap = $settings.reduce(
			(acc, setting) => ({ ...acc, [setting.key]: setting.value }),
			{}
		);
		agentCleanupEnabled = settingsMap['agent-cleanup-enabled'] === 'true';
		backupEnabled = settingsMap['backup-enabled'] === 'true';
	});
</script>

<svelte:head>
	<title>Settings</title>
</svelte:head>

<div class="mt-4 flex flex-col gap-4 p-4">
	<div class="container flex flex-col items-center justify-center gap-4 py-4">
		<Card.Root class="w-full sm:w-3/4 md:w-2/3">
			<Card.Header>
				<Card.Title class="flex flex-row items-center justify-between gap-2 text-xl font-bold">
					Settings
					{#if Object.keys(changedSettings)?.length > 0}
						<Tooltip.Root>
							<Tooltip.Trigger>
								<Button variant="ghost" size="icon" on:click={saveChanges}>
									<Save size="1rem" class="animate-pulse" />
								</Button>
							</Tooltip.Trigger>
							<Tooltip.Content>
								<p>Save changes</p>
							</Tooltip.Content>
						</Tooltip.Root>
					{/if}
				</Card.Title>
			</Card.Header>
			<Card.Content class="mt-4 flex flex-col gap-4">
				<h2 class="border-b border-gray-200 pb-2 text-lg">Server</h2>
				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="server-url" class="col-span-1 flex items-center gap-0.5">
						Server URL
						<HoverInfo text="The URL of the server which agents should use." />
					</Label>
					<Input
						name="server-url"
						type="text"
						bind:value={settingsMap['server-url']}
						on:keydown={(e) => onKeydown(e, 'server-url')}
						on:input={() => markAsChanged('server-url')}
						class="col-span-3 text-right"
					/>
				</div>

				<h2 class="border-b border-gray-200 pb-2 text-lg">Agents</h2>
				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="server-url" class="col-span-1 flex items-center gap-0.5">
						Cleanup
						<HoverInfo
							text="Automatically cleanup disconnected agents after a certain amount of time."
						/>
					</Label>
					<Switch
						name="agent-cleanup-enabled"
						class="col-span-3 justify-self-end"
						bind:checked={agentCleanupEnabled}
						onCheckedChange={(value) =>
							updateSetting({ key: 'agent-cleanup-enabled', value: value.toString() })}
					/>
				</div>
				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="server-url" class="col-span-1 flex items-center gap-0.5">
						Timout
						<HoverInfo
							text="The amount of time after which disconnected agents should be cleaned up. Valid time units are ns, us, ms, s, m, h."
						/>
					</Label>
					<Input
						name="agent-cleanup-timeout"
						type="text"
						bind:value={settingsMap['agent-cleanup-timeout']}
						on:keydown={(e) => onKeydown(e, 'agent-cleanup-timeout')}
						on:input={() => markAsChanged('agent-cleanup-timeout')}
						class="col-span-3 text-right"
					/>
				</div>

				<!-- Email -->
				<h2 class="border-b border-gray-200 pb-2 text-lg">Email</h2>
				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="ehost" class="col-span-1 flex items-center gap-0.5">
						Host
						<HoverInfo text="The IP/Domain of the smtp server." />
					</Label>
					<Input
						name="ehost"
						type="text"
						bind:value={settingsMap['email-host']}
						on:keydown={(e) => onKeydown(e, 'email-host')}
						on:input={() => markAsChanged('email-host')}
						class="col-span-3 text-right"
					/>
				</div>
				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="eport" class="col-span-1 flex items-center gap-0.5">
						Port
						<HoverInfo text="The port of the smtp server." />
					</Label>
					<Input
						name="eport"
						type="text"
						bind:value={settingsMap['email-port']}
						on:keydown={(e) => onKeydown(e, 'email-port')}
						on:input={() => markAsChanged('email-port')}
						class="col-span-3 text-right"
					/>
				</div>
				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="mail-user" class="col-span-1 flex items-center gap-0.5">
						Username
						<HoverInfo text="The username of the smtp server." />
					</Label>
					<Input
						name="mail-user"
						type="text"
						bind:value={settingsMap['email-username']}
						on:keydown={(e) => onKeydown(e, 'email-username')}
						on:input={() => markAsChanged('email-username')}
						class="col-span-3 text-right"
					/>
				</div>
				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="mail-password" class="col-span-1 flex items-center gap-0.5">
						Password
						<HoverInfo text="The password of the smtp server." />
					</Label>
					<div class="col-span-3 flex flex-row items-center justify-end gap-1">
						{#if showEmailPassword}
							<Input
								name="mail-password"
								type="text"
								bind:value={settingsMap['email-password']}
								on:keydown={(e) => onKeydown(e, 'email-password')}
								on:input={() => markAsChanged('email-password')}
								class="pr-10 text-right"
							/>
						{:else}
							<Input
								name="mail-password"
								type="password"
								bind:value={settingsMap['email-password']}
								on:keydown={(e) => onKeydown(e, 'email-password')}
								on:input={() => markAsChanged('email-password')}
								class="pr-10 text-right"
							/>
						{/if}
						<Button
							variant="ghost"
							size="icon"
							class="absolute hover:bg-transparent hover:text-red-400"
							on:click={() => (showEmailPassword = !showEmailPassword)}
						>
							{#if showEmailPassword}
								<Eye size="1rem" />
							{:else}
								<EyeOff size="1rem" />
							{/if}
						</Button>
					</div>
				</div>
				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="from" class="col-span-1 flex items-center gap-0.5">
						Sender
						<HoverInfo text="The from address of the email account." />
					</Label>
					<Input
						name="from"
						type="text"
						bind:value={settingsMap['email-from']}
						on:keydown={(e) => onKeydown(e, 'email-from')}
						on:input={() => markAsChanged('email-from')}
						class="col-span-3 text-right"
					/>
				</div>

				<h2 class="border-b border-gray-200 pb-2 text-lg">Backups</h2>
				<div class="mt-4 grid grid-cols-4 items-center justify-between gap-2">
					<Label for="backup-enabled" class="col-span-1">Enabled</Label>
					<Switch
						name="backup-enabled"
						class="col-span-3 justify-self-end"
						bind:checked={backupEnabled}
						onCheckedChange={(value) =>
							updateSetting({ key: 'backup-enabled', value: value.toString() })}
					/>
				</div>
				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="backup-keep" class="col-span-1 flex items-center gap-0.5">
						Retention
						<HoverInfo text="How many backups to keep. Set to 0 to keep all backups." />
					</Label>
					<Input
						name="backup-keep"
						type="text"
						bind:value={settingsMap['backup-keep']}
						on:keydown={(e) => onKeydown(e, 'backup-keep')}
						on:input={() => markAsChanged('backup-keep')}
						class="col-span-3 text-right"
						placeholder="3"
					/>
				</div>
				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="backup-schedule" class="col-span-1 flex items-center gap-0.5">
						Schedule
						<HoverInfo
							text="Cron expression for the backup schedule (e.g., * * * * *, or special keywords: @yearly, @annually, @monthly, @weekly, @daily"
						/>
					</Label>
					<Input
						name="backup-schedule"
						type="text"
						bind:value={settingsMap['backup-schedule']}
						on:keydown={(e) => onKeydown(e, 'backup-schedule')}
						class="col-span-3 text-right"
						placeholder="0 0 * * *"
					/>
				</div>

				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="backup" class="col-span-1 flex items-center gap-0.5">
						Backup & Restore
						<HoverInfo text="Manually backup and restore the database." />
					</Label>
					<div class="col-span-3 flex w-full gap-2">
						<input
							type="file"
							accept=".json"
							class="hidden"
							on:change={handleFileUpload}
							bind:this={fileInput}
							required
						/>
						<Button
							variant="ghost"
							class="w-full bg-orange-400"
							on:click={() => fileInput.click()}
							size="icon"
						>
							<Upload />
						</Button>
						<Button variant="default" class="w-full" on:click={() => downloadBackup()} size="icon">
							<Download />
						</Button>
					</div>
				</div>
			</Card.Content>
		</Card.Root>
	</div>
</div>
