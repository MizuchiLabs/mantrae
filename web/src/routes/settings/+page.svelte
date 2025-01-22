<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import { Input } from '$lib/components/ui/input';
	import { Separator } from '$lib/components/ui/separator';
	import { Download, Eye, EyeOff, Save, Upload } from 'lucide-svelte';
	import { settings, api } from '$lib/api';
	import type { Setting } from '$lib/types';
	import { onMount } from 'svelte';

	// State management
	// let fileInput = $state<HTMLInputElement>();
	let showEmailPassword = $state(false);
	let settingsMap = $state<Record<string, string>>({});
	let changedSettings = $state<Record<string, string>>({});

	// Computed values
	// let agentCleanupEnabled = $derived(settingsMap['agent-cleanup-enabled'] === 'true');
	// let backupEnabled = $derived(settingsMap['backup-enabled'] === 'true');
	let hasChanges = $derived(Object.keys(changedSettings).length > 0);

	// async function handleFileUpload(event: Event) {
	// 	const file = (event.target as HTMLInputElement).files?.[0];
	// 	if (file) {
	// 		await uploadBackup(file);
	// 		fileInput.value = '';
	// 	}
	// }

	async function updateSetting(key: string) {
		await api.upsertSetting({ key, value: settingsMap[key] } as Setting);
		// const { [key]: _, ...rest } = changedSettings;
		// changedSettings = rest;
	}

	function markAsChanged(key: string) {
		const originalValue = $settings.find((s: Setting) => s.key === key)?.value;
		if (settingsMap[key] !== originalValue) {
			changedSettings = { ...changedSettings, [key]: settingsMap[key] };
			// } else {
			// const { [key]: _, ...rest } = changedSettings;
			// changedSettings = rest;
		}
	}

	async function saveAllChanges() {
		await Promise.all(Object.keys(changedSettings).map((key) => updateSetting(key)));
		changedSettings = {};
	}

	function handleKeydown(e: KeyboardEvent, key: string) {
		if (e.key === 'Enter') updateSetting(key);
	}

	onMount(async () => {
		await api.listSettings();
		// settingsMap = $settings.reduce(
		// 	(acc, setting) => ({ ...acc, [setting.key]: setting.value }),
		// 	{}
		// );
	});
</script>

<svelte:head>
	<title>Settings</title>
</svelte:head>

<div class="container mx-auto max-w-4xl p-6">
	<Card.Root>
		<Card.Header>
			<div class="flex items-center justify-between">
				<Card.Title class="text-2xl font-bold">Settings</Card.Title>
				{#if hasChanges}
					<Tooltip.Root>
						<Tooltip.Trigger>
							<Button variant="outline" size="icon" onclick={saveAllChanges}>
								<Save class="h-4 w-4 animate-pulse text-green-500" />
							</Button>
						</Tooltip.Trigger>
						<Tooltip.Content>Save all changes</Tooltip.Content>
					</Tooltip.Root>
				{/if}
			</div>
		</Card.Header>

		<Card.Content class="space-y-6">
			{#each $settings as setting}
				{setting.key}
				<div class="grid grid-cols-4 items-center gap-4">
					<div class="flex items-center gap-2">
						<Label for={setting.key} class="font-medium">
							{setting.key
								.split('_')
								.map((word) => word.charAt(0).toUpperCase() + word.slice(1))
								.join(' ')}
						</Label>
						<!-- {#if setting.description} -->
						<!-- 	<HoverInfo text={setting.description} /> -->
						<!-- {/if} -->
					</div>

					{#if setting.key.includes('enabled')}
						<Switch
							id={setting.key}
							checked={settingsMap[setting.key] === 'true'}
							onCheckedChange={(value) => {
								settingsMap[setting.key] = value.toString();
								updateSetting(setting.key);
							}}
							class="col-span-3 justify-self-end"
						/>
					{:else if setting.key === 'email-password'}
						<div class="relative col-span-3">
							<Input
								id={setting.key}
								type={showEmailPassword ? 'text' : 'password'}
								value={settingsMap[setting.key]}
								oninput={() => markAsChanged(setting.key)}
								onkeydown={(e) => handleKeydown(e, setting.key)}
								class="pr-10"
							/>
							<Button
								variant="ghost"
								size="icon"
								class="absolute right-2 top-1/2 -translate-y-1/2"
								onclick={() => (showEmailPassword = !showEmailPassword)}
							>
								{#if showEmailPassword}
									<Eye class="h-4 w-4" />
								{:else}
									<EyeOff class="h-4 w-4" />
								{/if}
							</Button>
						</div>
					{:else}
						<Input
							id={setting.key}
							type="text"
							value={settingsMap[setting.key]}
							oninput={() => markAsChanged(setting.key)}
							onkeydown={(e) => handleKeydown(e, setting.key)}
							class="col-span-3"
						/>
					{/if}
				</div>

				<Separator />
			{/each}

			<!-- <div class="flex gap-4 pt-4"> -->
			<!-- 	<input -->
			<!-- 		type="file" -->
			<!-- 		accept=".json" -->
			<!-- 		class="hidden" -->
			<!-- 		onchange={handleFileUpload} -->
			<!-- 		bind:this={fileInput} -->
			<!-- 	/> -->
			<!-- 	<Button  -->
			<!-- 		variant="outline"  -->
			<!-- 		class="flex-1"  -->
			<!-- 		onclick={() => fileInput.click()} -->
			<!-- 	> -->
			<!-- 		<Upload class="mr-2 h-4 w-4" /> -->
			<!-- 		Upload Backup -->
			<!-- 	</Button> -->
			<!-- 	<Button  -->
			<!-- 		variant="default"  -->
			<!-- 		class="flex-1"  -->
			<!-- 		onclick={() => downloadBackup()} -->
			<!-- 	> -->
			<!-- 		<Download class="mr-2 h-4 w-4" /> -->
			<!-- 		Download Backup -->
			<!-- 	</Button> -->
			<!-- </div> -->
		</Card.Content>
	</Card.Root>
</div>
