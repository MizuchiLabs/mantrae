<script lang="ts">
	import { Switch } from '$lib/components/ui/switch/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { type Router, type Service } from '$lib/types/config';
	import ArrayInput from '../ui/array-input/array-input.svelte';
	import { z } from 'zod';
	import { services } from '$lib/api';
	import { onMount } from 'svelte';

	export let router: Router;
	export let service: Service;
	export let disabled = false;
	$: router, routerChange();

	let passHostHeader = service?.loadBalancer?.passHostHeader ?? true;
	let servers: string[] = [];

	const serviceSchema = z.object({
		provider: z.string().trim().optional(),
		type: z.string().trim().optional(),
		status: z.string().trim().optional(),
		protocol: z
			.string()
			.trim()
			.toLowerCase()
			.regex(/^(http|tcp|udp)$/),
		serverStatus: z.record(z.string().trim()).optional(),
		loadBalancer: z.object({
			servers: z
				.array(
					z
						.object({
							url: z.string().trim().optional(),
							address: z.string().trim().optional()
						})
						.refine((data) => data.url || data.address, {
							message: 'At least one server is required',
							path: ['servers'] // Points to the 'servers' array in case of error
						})
				)
				.nonempty('At least one server is required'),
			passHostHeader: z.boolean().optional()
		})
	});

	let errors: Record<any, string[] | undefined> = {};
	export const validate = () => {
		try {
			serviceSchema.parse({ ...service });
			errors = {};
			return true;
		} catch (err) {
			if (err instanceof z.ZodError) {
				errors = err.flatten().fieldErrors;
			}
			return false;
		}
	};

	const update = () => {
		if (service.loadBalancer === undefined) service.loadBalancer = { servers: [] };

		service.loadBalancer.passHostHeader = passHostHeader;
		switch (service.protocol) {
			case 'http':
				service.loadBalancer.servers = servers.map((s) => {
					return { url: s };
				});
				break;
			case 'tcp':
			case 'udp':
				service.loadBalancer.servers = servers.map((s) => {
					return { address: s };
				});
				break;
		}
	};

	const routerChange = () => {
		service.name = router.name;
		service.provider = router.provider;
		service.profileId = router.profileId;
		service.protocol = router.protocol;
		service = $services.find((s) => s.name === router.name) ?? service;
		servers = service.loadBalancer?.servers?.map((s) => s.url ?? s.address ?? '') ?? [];
	};

	const onSwitchChange = () => {
		passHostHeader = !passHostHeader;
		if (service.loadBalancer === undefined) service.loadBalancer = { servers: [] };
		service.loadBalancer.passHostHeader = passHostHeader;
	};

	// onMount(() => {
	// 	servers = service.loadBalancer?.servers?.map((s) => s.url ?? s.address ?? '') ?? [];
	// 	update();
	// });
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Service</Card.Title>
		<Card.Description>Configure your {router.protocol} service</Card.Description>
	</Card.Header>
	<Card.Content class="flex flex-col gap-2">
		{#if service.protocol === 'http'}
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="passHostHeader" class="text-right">Pass Host Header</Label>
				<Switch
					id="passHostHeader"
					class="col-span-3"
					bind:checked={passHostHeader}
					onCheckedChange={onSwitchChange}
					{disabled}
				/>
			</div>
		{/if}
		<ArrayInput
			bind:items={servers}
			label="Servers"
			placeholder="http://192.168.1.1:8080"
			on:update={update}
			{disabled}
		/>
		{#if errors.loadBalancer}
			<div class="col-span-4 text-right text-sm text-red-500">{errors.loadBalancer}</div>
		{/if}
	</Card.Content>
</Card.Root>
