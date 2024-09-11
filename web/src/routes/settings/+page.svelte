<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch';
	import { downloadBackup, getSettings, uploadBackup, settings, updateSetting } from '$lib/api';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { Input } from '$lib/components/ui/input';
	import type { Setting } from '$lib/types/base';

	let fileInput: HTMLInputElement;
	const handleFileUpload = (event: Event) => {
		const file = (event.target as HTMLInputElement).files?.[0];
		if (file) {
			uploadBackup(file);
		}
		fileInput.value = '';
	};

	// Settings
	let backupEnabled: boolean;
	let backupKeep: string;
	let backupSchedule: string;

	const update = async (s: Setting) => {
		await updateSetting(s);
		toast.success(`Setting ${s.key} updated`);
	};

	const onKeydown = (e: KeyboardEvent, s: any) => {
		if (e.key === 'Enter') {
			update(s);
		}
	};

	onMount(async () => {
		await getSettings();
		backupEnabled = $settings?.find((s) => s.key === 'backup-enabled')?.value === 'true';
		backupKeep = $settings?.find((s) => s.key === 'backup-keep')?.value ?? '';
		backupSchedule = $settings?.find((s) => s.key === 'backup-schedule')?.value ?? '';
	});
</script>

<svelte:head>
	<title>Settings</title>
</svelte:head>

<div class="mt-4 flex flex-col gap-4 p-4">
	<div class="container mx-auto flex flex-col items-center justify-center gap-4 py-4">
		<Card.Root class="w-[800px]">
			<Card.Header>
				<Card.Title class="flex flex-row items-center justify-between gap-2 text-xl font-bold">
					Settings
				</Card.Title>
			</Card.Header>
			<Card.Content class="mt-4 flex flex-col gap-4">
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
					<Label for="backup-keep" class="col-span-1">Retention</Label>
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
					<Label for="backup-schedule" class="col-span-1">Schedule</Label>
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
					<Label for="backup" class="col-span-1">Backup & Restore</Label>
					<div class="col-span-3 flex w-full gap-2">
						<input
							type="file"
							accept=".tar.gz"
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
