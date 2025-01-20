<script lang="ts">
	import ArrayInput from '../ui/array-input/array-input.svelte';
	import RuleEditor from '../utils/ruleEditor.svelte';
	import autoAnimate from '@formkit/auto-animate';
	import logo from '$lib/images/logo.svg';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { CircleCheck, Lock } from 'lucide-svelte';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Toggle } from '$lib/components/ui/toggle';
	import { api, dnsProviders, entrypoints, routers, middlewares } from '$lib/api';
	import { onMount } from 'svelte';
	import { type Router } from '$lib/types/router';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';

	interface Props {
		router: Router | undefined;
		disabled?: boolean;
	}

	let { router = $bindable({} as Router), disabled = false }: Props = $props();
	let routerProvider = router.name ? router.name.split('@')[1].toLowerCase() : 'http';

	const getCertResolver = () => {
		const certResolvers = $routers
			.filter((item) => item.tls && item.tls.certResolver)
			.map((item) => item.tls?.certResolver);

		// Remove duplicates
		return [...new Set(certResolvers)];
	};

	onMount(async () => {
		await api.listDNSProviders();
	});
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Router</Card.Title>
		<Card.Description>Configure your router settings</Card.Description>
	</Card.Header>
	<Card.Content class="flex flex-col gap-2">
		<!-- Provider -->
		{#if routerProvider === 'http'}
			<div class="flex items-center justify-end gap-1 font-mono text-sm" use:autoAnimate>
				<Toggle
					size="sm"
					pressed={router.type === 'http'}
					onPressedChange={() => (router.type = 'http')}
					class="font-bold data-[state=on]:bg-green-300  dark:data-[state=on]:text-black"
				>
					HTTP
				</Toggle>
				<Toggle
					size="sm"
					pressed={router.type === 'tcp'}
					onPressedChange={() => (router.type = 'tcp')}
					class="font-bold data-[state=on]:bg-blue-300 dark:data-[state=on]:text-black"
				>
					TCP
				</Toggle>
				<Toggle
					size="sm"
					pressed={router.type === 'udp'}
					onPressedChange={() => (router.type = 'udp')}
					class="font-bold data-[state=on]:bg-red-300 dark:data-[state=on]:text-black"
				>
					UDP
				</Toggle>
			</div>
		{/if}

		<!-- Name -->
		<div class="grid grid-cols-4 items-center gap-1">
			<Label for="name" class="mr-2 text-right">Name</Label>
			<div class="relative col-span-3">
				<Input
					id="name"
					name="name"
					type="text"
					bind:value={router.name}
					placeholder="Name of the router"
					required
					{disabled}
				/>
				<!-- Icon based on provider -->
				{#if routerProvider !== ''}
					<span
						class="pointer-events-none absolute inset-y-0 right-3 flex items-center text-gray-400"
					>
						{#if routerProvider === 'http'}
							<img src={logo} alt="HTTP" width="20" />
						{/if}
						{#if routerProvider === 'internal' || routerProvider === 'file'}
							<iconify-icon icon="devicon:traefikproxy" height="20"></iconify-icon>
						{/if}
						{#if routerProvider?.includes('docker')}
							<iconify-icon icon="logos:docker-icon" height="20"></iconify-icon>
						{/if}
						{#if routerProvider?.includes('kubernetes')}
							<iconify-icon icon="logos:kubernetes" height="20"></iconify-icon>
						{/if}
						{#if routerProvider === 'consul'}
							<iconify-icon icon="logos:consul" height="20"></iconify-icon>
						{/if}
						{#if routerProvider === 'nomad'}
							<iconify-icon icon="logos:nomad-icon" height="20"></iconify-icon>
						{/if}
						{#if routerProvider === 'kv'}
							<iconify-icon icon="logos:redis" height="20"></iconify-icon>
						{/if}
					</span>
				{/if}
			</div>
		</div>

		<!-- Entrypoints -->
		{#if routerProvider === 'http'}
			<div class="grid grid-cols-4 items-center gap-1">
				<Label for="entrypoints" class="mr-2 text-right">Entrypoints</Label>
				<Select.Root
					type="multiple"
					bind:value={router.entryPoints}
					onValueChange={(value) => (router.entryPoints = value)}
				>
					<Select.Trigger class="col-span-3">
						{#if router.entryPoints && router.entryPoints.length > 0}
							{router.entryPoints.join(', ')}
						{:else}
							Select an entrypoint
						{/if}
					</Select.Trigger>
					<Select.Content>
						{#each $entrypoints || [] as e}
							<Select.Item value={e.name}>
								<div class="flex flex-row items-center gap-2">
									{e.name}
									{#if e.http}
										{#if 'tls' in e.http}
											<Lock size="1rem" class="text-green-400" />
										{/if}
									{/if}
								</div>
							</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		{:else}
			<ArrayInput
				bind:items={router.entryPoints}
				label="Entrypoints"
				placeholder=""
				disabled={true}
			/>
		{/if}

		<!-- Middlewares -->
		{#if routerProvider === 'http' && router.type !== 'udp'}
			<div class="grid grid-cols-4 items-center gap-1">
				<Label for="middlewares" class="mr-2 text-right">Middlewares</Label>
				<Select.Root
					type="multiple"
					value={router.middlewares}
					onValueChange={(value) => (router.middlewares = value)}
				>
					<Select.Trigger class="col-span-3">
						{#if router.middlewares && router.middlewares.length > 0}
							{router.middlewares.join(', ')}
						{:else}
							Select middlewares
						{/if}
					</Select.Trigger>
					<Select.Content>
						{#each $middlewares as m}
							{#if router.type === m.type}
								<Select.Item value={m.name}>
									{m.name}
								</Select.Item>
							{/if}
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		{:else if router.type !== 'udp'}
			<ArrayInput
				bind:items={router.middlewares}
				label="Middlewares"
				placeholder="None"
				disabled={true}
			/>
		{/if}

		<!-- DNS Provider -->
		{#if $dnsProviders}
			<div class="grid grid-cols-4 items-center gap-1">
				<Label for="provider" class="mr-2 text-right">DNS Provider</Label>
				<Select.Root type="single">
					<Select.Trigger class="col-span-3">Select a dns provider</Select.Trigger>
					<Select.Content>
						<Select.Item value="" label="">None</Select.Item>
						{#each $dnsProviders as dns}
							<Select.Item value={dns.id.toString()} class="flex items-center gap-2">
								{dns.name} ({dns.type})
								{#if dns.is_active}
									<CircleCheck size="1rem" class="text-green-400" />
								{/if}
							</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		{/if}

		<!-- CertResolver -->
		{#if routerProvider === 'http' && router.type !== 'udp' && router.tls}
			<div class="space-y-0.5">
				<div class="grid grid-cols-4 items-center gap-1">
					<Label for="certresolver" class="mr-2 text-right">CertResolver</Label>
					<Input
						id="certresolver"
						name="certresolver"
						type="text"
						class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
						bind:value={router.tls.certResolver}
						placeholder="Certificate resolver"
					/>
				</div>
				<div class="grid grid-cols-4 flex-wrap items-center gap-1">
					<div></div>
					<div class="col-span-3">
						{#each getCertResolver() || [] as r}
							{#if router.tls.certResolver !== r}
								<div onclick={() => (router.tls.certResolver = r)} aria-hidden="true">
									<Badge>{r}</Badge>
								</div>
							{/if}
						{/each}
					</div>
				</div>
			</div>
		{/if}

		<!-- Rule -->
		{#if router.type === 'http' || router.type === 'tcp'}
			<RuleEditor bind:rule={router.rule} bind:type={router.type} {disabled} />
		{/if}
	</Card.Content>
</Card.Root>
