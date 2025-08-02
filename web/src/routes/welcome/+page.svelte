<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import { checkHealth } from '$lib/api';
	import { baseURL } from '$lib/stores/common';
	import { Loader } from '@lucide/svelte';
	import { onDestroy, onMount } from 'svelte';

	let url = $state(baseURL.value ?? '');
	let isChecking = $state(false);
	let isHealthy = $state(false);
	let reconnectTimer = $state<ReturnType<typeof setInterval> | null>(null);

	const silentCheckConnection = async () => {
		if (isChecking) return; // Don't run if already checking

		try {
			const healthy = await checkHealth();
			isHealthy = healthy;

			if (healthy && reconnectTimer) {
				toast.success('Backend reconnected!');
				clearInterval(reconnectTimer);
				reconnectTimer = null;
				goto('/login');
			}
		} catch (_) {
			isHealthy = false;
		}
	};

	const checkConnection = async (testUrl?: string) => {
		isChecking = true;
		const urlToTest = testUrl || url;

		try {
			// Temporarily update baseURL for health check
			const originalUrl = baseURL.value;
			baseURL.value = urlToTest;

			const healthy = await checkHealth();
			isHealthy = healthy;

			if (healthy) {
				toast.success('Backend connected successfully!');
				goto('/login');
			} else {
				baseURL.value = originalUrl; // Restore if failed
				toast.error('Unable to connect to backend');
			}
		} catch (_) {
			toast.error('Connection failed');
			isHealthy = false;
		} finally {
			isChecking = false;
		}
	};

	const handleSave = () => {
		if (!url.trim()) {
			toast.error('Please enter a backend URL');
			return;
		}
		if (!url.startsWith('http://') && !url.startsWith('https://')) {
			url = `http://${url}`;
		}

		baseURL.value = url;
		checkConnection(url);
	};

	const handleKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') {
			handleSave();
		}
	};

	// Check initial connection status
	$effect(() => {
		checkConnection(baseURL.value);
	});

	onMount(() => {
		reconnectTimer = setInterval(() => {
			silentCheckConnection();
		}, 5000);
	});
	onDestroy(() => {
		if (reconnectTimer) {
			clearInterval(reconnectTimer);
			reconnectTimer = null;
		}
	});
</script>

<Card.Root class="max-w-md sm:min-w-[350px]">
	<Card.Header class="flex flex-col items-center gap-3 text-center">
		<div class="flex items-center gap-2">
			<Card.Title class="text-2xl font-bold">Server Configuration</Card.Title>
		</div>

		<Badge variant={isHealthy ? 'default' : 'destructive'}>
			{#if isChecking}
				Checking connection...
			{:else if isHealthy}
				Connected
			{:else}
				Disconnected
			{/if}
		</Badge>

		<Card.Description class="text-center">
			Enter the URL where your backend server is running.
		</Card.Description>
	</Card.Header>

	<Card.Content class="space-y-4">
		<div>
			<div class="space-y-2">
				<Label for="backend-url">Backend URL</Label>
				<Input
					id="backend-url"
					type="url"
					placeholder="http://localhost:3000"
					bind:value={url}
					onkeydown={handleKeydown}
					disabled={isChecking}
				/>
			</div>
			{#if !isHealthy && reconnectTimer}
				<small class="text-muted-foreground animate-pulse">Attempting to reconnect...</small>
			{/if}
		</div>

		<Button class="w-full" onclick={handleSave} disabled={isChecking || !url.trim()}>
			{#if isChecking}
				<Loader class="mr-2 h-4 w-4 animate-spin" />
				Connecting...
			{:else}
				Save & Connect
			{/if}
		</Button>
	</Card.Content>
</Card.Root>
