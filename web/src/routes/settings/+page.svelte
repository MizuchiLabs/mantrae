<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import { Input } from '$lib/components/ui/input';
	import { Separator } from '$lib/components/ui/separator';
	import { SaveIcon, Settings } from 'lucide-svelte';
	import { settings, api } from '$lib/api';
	import { onMount } from 'svelte';
	import type { Setting } from '$lib/types';
	import { toast } from 'svelte-sonner';

	// State management
	// let fileInput = $state<HTMLInputElement>();

	// async function handleFileUpload(event: Event) {
	// 	const file = (event.target as HTMLInputElement).files?.[0];
	// 	if (file) {
	// 		await uploadBackup(file);
	// 		fileInput.value = '';
	// 	}
	// }

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
	const getInputType = (value: Setting) => {
		if (typeof value === 'boolean') return 'boolean';
		if (typeof value === 'number') return 'number';
		if (value?.toString().includes('://')) return 'url';
		if (value?.toString().includes('@')) return 'email';
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

	onMount(async () => {
		await api.listSettings();
	});
</script>

<svelte:head>
	<title>Settings</title>
</svelte:head>

<div class="container">
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
			{#each Object.entries($settings) as [key, setting]}
				<div class="flex flex-col justify-start gap-4 sm:flex-row sm:justify-between">
					<Label>
						{formatSettingName(key)}
						{#if setting.description}
							<p class="text-sm text-muted-foreground">{setting.description}</p>
						{/if}
					</Label>

					<div class="flex w-full items-center justify-end gap-4 sm:w-auto md:w-[380px]">
						{#if getInputType(setting.value) === 'boolean'}
							<Switch
								id={key}
								checked={setting.value}
								onCheckedChange={(checked) => saveSetting(key, checked)}
							/>
						{:else if key.includes('interval')}
							<Input
								type="text"
								id={key}
								value={setting.value}
								onchange={(e) => handleChange(key, parseDuration(e.currentTarget.value))}
								onkeydown={(e) => handleKeydown(e, key, parseDuration(e.currentTarget.value))}
							/>
						{:else if key.includes('port')}
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
								onchange={(e) => handleChange(key, e.currentTarget.value)}
								onkeydown={(e) => handleKeydown(e, key, e.currentTarget.value)}
							/>
						{/if}
					</div>
				</div>

				<Separator />
			{/each}

			<div class="flex justify-end">
				<Button variant={hasChanges ? 'default' : 'outline'} onclick={saveAllChanges} size="icon">
					<SaveIcon />
				</Button>
			</div>
		</Card.Content>
	</Card.Root>
</div>
