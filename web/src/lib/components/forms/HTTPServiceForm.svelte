<script lang="ts">
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { type Service } from '$lib/gen/mantrae/v1/service_pb';
	import {
		type Cookie,
		type Service as HTTPService,
		type Server,
		type ServersLoadBalancer, ServerHealthCheckSchema
	} from '$lib/gen/zen/traefik-schemas';
	import { Plus, Trash, BrushCleaning } from '@lucide/svelte';
	import {marshalConfig, unmarshalConfig} from '$lib/types';
	import CustomSwitch from '../ui/custom-switch/custom-switch.svelte';
	import { serversTransportClient } from '$lib/api';
	import { profile } from '$lib/stores/profile';
	import {string} from "zod";
	import {Textarea} from "$lib/components/ui/textarea";

	const getInputType = (key: keyof HealthCheckType) => {
		const typeMap: Record<keyof HealthCheckType> = {
			scheme: 'text',
			mode: 'text',
			path: 'text',
			method: 'text',
			status: 'number',
			port: 'number',
			interval: 'text',
			unhealthyInterval: 'text',
			timeout: 'text',
			hostname: 'text',
			followRedirects: 'checkbox',
		};
		return typeMap[key];
	}

	const healthcheckPlaceholders = {
		scheme: "http",
		mode: "http",
		path: "/health",
		method: "GET",
		status: "200",
		port: "80",
		interval: "10s",
		unhealthyInterval: "3s",
		timeout: "2s",
		hostname: "example.internal",
		followRedirects: "true",
		headers: "X-Health-Check: true"
	};

	interface Props {
		service: Service;
	}
	let { service = $bindable() }: Props = $props();

	let servers = $state([{ url: '' }] as Server[]);
	let passHostHeader = $state(true);
	let serversTransport = $state('');
	let sticky = $state(false);
	let cookie = $state({} as Cookie);
	let healthCheck = $state({});
	let headers = $state({ key: '' });

	function updateConfig() {
		const config = {} as HTTPService;
		config.loadBalancer = {} as ServersLoadBalancer;
		config.loadBalancer.servers = servers;
		config.loadBalancer.passHostHeader = passHostHeader;
		config.loadBalancer.sticky = sticky ? { cookie } : undefined;
		config.loadBalancer.serversTransport = serversTransport;
		config.loadBalancer.healthCheck = healthCheck ?? {};
		config.loadBalancer.healthCheck.headers = headers ?? {};

		Object.entries(healthCheck).forEach(([key, value]) => {
			if (value === "") {
				delete healthCheck[key];
			}
		});

		service.config = marshalConfig(config);


	}

	$effect(() => {
		if (service.config) {
			let config = service.config as HTTPService;
			servers = config.loadBalancer?.servers || [];
			passHostHeader = config.loadBalancer?.passHostHeader ?? true;
			sticky = config.loadBalancer?.sticky?.cookie !== undefined;
			cookie = config.loadBalancer?.sticky?.cookie ?? {};
			serversTransport = config.loadBalancer?.serversTransport ?? '';
			healthCheck = config.loadBalancer?.healthCheck ?? {};
			headers = config.loadBalancer?.healthCheck?.headers ?? {};
		}
	});

	function capitalizeFirstLetter(string) {
		return string.replace(/^./, string[0].toUpperCase())
	}

	function getHeaderKeyNames() {
		return Object.keys(headers);
	}




</script>

<div class="flex flex-col gap-3">
	<div class="flex flex-col gap-2">
		<Label for="servers">Server Endpoints</Label>
		{#each servers || [] as server, i (i)}
			<div class="flex gap-2">
				<Input
						type="text"
						bind:value={server.url}
						oninput={updateConfig}
						placeholder="http://127.0.0.1:8080"
				/>
				<Button
						variant="ghost"
						size="icon"
						type="button"
						class="text-red-500"
						onclick={() => {
						if (i === 0) return;
						servers = servers.filter((_, j) => j !== i);
					}}
				>
					<Trash />
				</Button>
			</div>
		{/each}
	</div>
	<Button
			type="button"
			variant="outline"
			class="w-full"
			onclick={() => (servers = [...servers, { url: '' }])}
	>
		<Plus />
		Add Server
	</Button>

	<div class="flex items-center justify-between rounded-lg border p-3">
		<div class="space-y-1">
			<Label class="flex items-center gap-1 text-sm">Pass Host Header</Label>
			<p class="text-xs text-muted-foreground">Forward client host header to server</p>
		</div>
		<CustomSwitch bind:checked={passHostHeader} onCheckedChange={updateConfig} size="md" />
	</div>

	<div class=" items-center justify-between rounded-lg border p-3">
		<Label class="flex items-center text-1xl">Healthcheck</Label>
		<p class="text-xs text-muted-foreground pb-3">Monitors service health</p>

		<div class="flex items-center justify-between rounded-lg border p-3">
				<div class="space-y-1">
					<Label class="flex items-center gap-1 text-xs">Clear Healthcheck</Label>
				</div>
				<Button
						variant="ghost"
						size="default"
						type="button"
						class="text-orange-300"

						onclick={()=> {healthCheck = {}; headers = {}; updateConfig();}}
				>
					<BrushCleaning />
				</Button>
			</div>

		{#each Object.keys(ServerHealthCheckSchema.shape) as key}
			{@const inputType = getInputType(key)}

			{#if inputType !== 'checkbox' && key !== 'headers'}
				<div class="flex flex-col gap-2 pt-2">
					<Label for="cookie-name">{capitalizeFirstLetter(key)}</Label>
					<Input
							id="cookie-name"
							type={inputType}
							bind:value={healthCheck[key]}
							oninput={updateConfig}
							placeholder="{healthcheckPlaceholders[key]}"
					/>
				</div>
			{:else if inputType === 'checkbox' && key !== 'headers'}
				<div class="flex flex-col gap-2 pt-2 pb-2">
					<Label for="cookie-name">{capitalizeFirstLetter(key)}</Label>
					<CustomSwitch bind:checked={healthCheck[key]} onCheckedChange={updateConfig} size="md" />
				</div>
			{:else}
				<Label class="pt-2 pb-2" for="cookie-name">{capitalizeFirstLetter(key)}</Label>
				{#each Object.entries(headers) || {} as [keyName, value], i (i)}
					{@const nameofKey = getHeaderKeyNames()}
					<div class="flex gap-4 pb-2">
						<div class="col-span-3 space-y-1">

							<Input
									type="text"
									bind:value={nameofKey[i]}
									oninput={(e) => {
										const newKey = e.target.value;
										if (newKey === keyName) return;

										headers = {
											...Object.fromEntries(
												Object.entries(headers).filter(([k]) => k !== keyName)
											),
											[newKey]: value
										};
										updateConfig();
									}}

							/>
						</div>
						<div class="col-span-3 space-y-1">
							<Input
									type="text"
									bind:value={headers[keyName]}
									oninput={updateConfig}
									}}
							/>
						</div>
						<Button
								variant="ghost"
								size="icon"
								type="button"
								class="text-red-500"
								onclick={() => { delete headers[keyName]; updateConfig();}}
						>
							<Trash />
						</Button>
					</div>

				{/each}
				<Button
						type="button"
						variant="outline"
						class="w-full"
						onclick={

						() => {
							const randomString = crypto.randomUUID();
							headers = {...headers, [randomString]: 'value' };
							updateConfig();
						}

						}
				>
					<Plus />
					Add Header
				</Button>
			{/if}
		{/each}
	</div>
	<div class="grid grid-cols-4 items-center gap-3 rounded-lg border p-3">
		<div class="col-span-3 space-y-1">
			<Label class="flex items-center gap-1 text-sm">Sticky Cookie</Label>
			<p class="text-xs text-muted-foreground">
				Defines a Set-Cookie header is set on the initial response to let the client know which
				server handles the first response.
			</p>
		</div>
		<div class="justify-self-end">
			<CustomSwitch bind:checked={sticky} onCheckedChange={updateConfig} size="md" />
		</div>
	</div>

	{#if sticky}
		<div class="flex flex-col gap-2">
			<Label for="cookie-name">Cookie Name</Label>
			<Input
					id="cookie-name"
					type="text"
					bind:value={cookie.name}
					oninput={updateConfig}
					placeholder="session"
			/>
		</div>
	{/if}

	{#await serversTransportClient.listServersTransports({ profileId: profile.id }) then value}
		<div class="flex flex-col gap-2">
			<Label for="servers-transport">Servers Transport</Label>
			<Select.Root type="single" bind:value={serversTransport} onValueChange={updateConfig}>
				<Select.Trigger class="w-full">
                    <span class="truncate text-left">
                        {serversTransport || 'Select servers transport'}
                    </span>
				</Select.Trigger>
				<Select.Content>
					<Select.Item value="">
						<span class="truncate">Select servers transport</span>
					</Select.Item>
					{#each value.serversTransports || [] as transport (transport.id)}
						<Select.Item value={transport.name}>
							<span class="truncate">{transport.name}</span>
						</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</div>
	{/await}
</div>