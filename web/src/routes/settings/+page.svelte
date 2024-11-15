<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch';
	import { downloadBackup, getSettings, uploadBackup, settings, updateSetting } from '$lib/api';
	import { onMount } from 'svelte';
	import { Input } from '$lib/components/ui/input';
	import type { Setting } from '$lib/types/base';
	import HoverInfo from '$lib/components/utils/hoverInfo.svelte';

	let fileInput: HTMLInputElement;
	const handleFileUpload = (event: Event) => {
		const file = (event.target as HTMLInputElement).files?.[0];
		if (file) {
			uploadBackup(file);
		}
		fileInput.value = '';
	};

	// Settings
	let serverURL: string;
	let agentCleanupEnabled: boolean;
	let agentCleanupTimout: string;
	let emailHost: string;
	let emailPort: string;
	let emailUsername: string;
	let emailPassword: string;
	let emailFrom: string;
	let backupEnabled: boolean;
	let backupKeep: string;
	let backupSchedule: string;

	const update = async (s: Setting) => {
		await updateSetting(s);
	};

	const onKeydown = (e: KeyboardEvent, s: any) => {
		if (e.key === 'Enter') {
			update(s);
		}
	};

	onMount(async () => {
		await getSettings();
		serverURL = $settings?.find((s) => s.key === 'server-url')?.value ?? '';

		agentCleanupEnabled =
			$settings?.find((s) => s.key === 'agent-cleanup-enabled')?.value === 'true';
		agentCleanupTimout = $settings?.find((s) => s.key === 'agent-cleanup-timeout')?.value ?? '';

		backupEnabled = $settings?.find((s) => s.key === 'backup-enabled')?.value === 'true';
		backupKeep = $settings?.find((s) => s.key === 'backup-keep')?.value ?? '';
		backupSchedule = $settings?.find((s) => s.key === 'backup-schedule')?.value ?? '';

		emailHost = $settings?.find((s) => s.key === 'email-host')?.value ?? '';
		emailPort = $settings?.find((s) => s.key === 'email-port')?.value ?? '';
		emailUsername = $settings?.find((s) => s.key === 'email-username')?.value ?? '';
		emailPassword = $settings?.find((s) => s.key === 'email-password')?.value ?? '';
		emailFrom = $settings?.find((s) => s.key === 'email-from')?.value ?? '';
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
						on:keydown={(e) => onKeydown(e, { key: 'server-url', value: serverURL })}
						bind:value={serverURL}
						class="col-span-3 text-right"
					/>
				</div>

				<h2 class="border-b border-gray-200 pb-2 text-lg">Agents</h2>
				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="server-url" class="col-span-1 flex items-center gap-0.5">
						Cleanup Enabled
						<HoverInfo
							text="Automatically cleanup disconnected agents after a certain amount of time."
						/>
					</Label>
					<Switch
						name="agent-cleanup-enabled"
						class="col-span-3 justify-self-end"
						bind:checked={agentCleanupEnabled}
						onCheckedChange={(value) =>
							updateSetting({ id: 0, key: 'agent-cleanup-enabled', value: value.toString() })}
					/>
				</div>
				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="server-url" class="col-span-1 flex items-center gap-0.5">
						Cleanup Timout
						<HoverInfo
							text="The amount of time after which disconnected agents should be cleaned up. Valid time units are ns, us, ms, s, m, h."
						/>
					</Label>
					<Input
						name="agent-cleanup-timeout"
						type="text"
						on:keydown={(e) =>
							onKeydown(e, { key: 'agent-cleanup-timeout', value: agentCleanupTimout })}
						bind:value={agentCleanupTimout}
						class="col-span-3 text-right"
					/>
				</div>

				<h2 class="border-b border-gray-200 pb-2 text-lg">Email</h2>
				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="email-host" class="col-span-1 flex items-center gap-0.5">
						Host
						<HoverInfo text="The host of the email server." />
					</Label>
					<Input
						name="email-host"
						type="text"
						on:keydown={(e) => onKeydown(e, { key: 'email-host', value: emailHost })}
						bind:value={emailHost}
						class="col-span-3 text-right"
					/>
				</div>
				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="email-port" class="col-span-1 flex items-center gap-0.5">
						Port
						<HoverInfo text="The port of the email server." />
					</Label>
					<Input
						name="email-port"
						type="text"
						on:keydown={(e) => onKeydown(e, { key: 'email-port', value: emailPort })}
						bind:value={emailPort}
						class="col-span-3 text-right"
					/>
				</div>
				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="email-username" class="col-span-1 flex items-center gap-0.5">
						Username
						<HoverInfo text="The username of the email account." />
					</Label>
					<Input
						name="email-username"
						type="text"
						on:keydown={(e) => onKeydown(e, { key: 'email-username', value: emailUsername })}
						bind:value={emailUsername}
						class="col-span-3 text-right"
					/>
				</div>
				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="email-password" class="col-span-1 flex items-center gap-0.5">
						Password
						<HoverInfo text="The password of the email account." />
					</Label>
					<Input
						name="email-password"
						type="password"
						on:keydown={(e) => onKeydown(e, { key: 'email-password', value: emailPassword })}
						bind:value={emailPassword}
						class="col-span-3 text-right"
					/>
				</div>
				<div class="grid grid-cols-4 items-center justify-between gap-2">
					<Label for="email-from" class="col-span-1 flex items-center gap-0.5">
						Sender
						<HoverInfo text="The from address of the email account." />
					</Label>
					<Input
						name="email-from"
						type="text"
						on:keydown={(e) => onKeydown(e, { key: 'email-from', value: emailFrom })}
						bind:value={emailFrom}
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
							updateSetting({ id: 0, key: 'backup-enabled', value: value.toString() })}
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
						on:keydown={(e) => onKeydown(e, { key: 'backup-keep', value: backupKeep })}
						bind:value={backupKeep}
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
						on:keydown={(e) => onKeydown(e, { key: 'backup-schedule', value: backupSchedule })}
						bind:value={backupSchedule}
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
						<Button variant="ghost" class="w-full bg-orange-400" on:click={() => fileInput.click()}>
							<iconify-icon icon="fa6-solid:upload" width="16" height="16" />
						</Button>
						<Button variant="default" class="w-full" on:click={() => downloadBackup()}>
							<iconify-icon icon="fa6-solid:download" width="16" height="16" />
						</Button>
					</div>
				</div>
			</Card.Content>
		</Card.Root>
	</div>
</div>
