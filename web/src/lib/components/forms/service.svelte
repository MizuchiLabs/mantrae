<script lang="ts">
	import { Switch } from '$lib/components/ui/switch/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { type Service } from '$lib/types/config';
	import ArrayInput from '../ui/array-input/array-input.svelte';
	import { z } from 'zod';
	import { onMount } from 'svelte';

	export let service: Service;
	export let disabled = false;
	let passHostHeader = service?.loadBalancer?.passHostHeader ?? true;
	let servers: string[] = [];

	let errors: Record<any, string[] | undefined> = {};
	const formSchema = z
		.object({
			serviceType: z
				.string()
				.toLowerCase()
				.regex(/^(http|tcp|udp)$/),
			loadBalancer: z
				.object({
					servers: z.array(z.object({ url: z.string().trim() })).optional()
				})
				.nullable()
				.optional(),
			tcpLoadBalancer: z
				.object({
					servers: z.array(z.object({ address: z.string().trim() })).optional()
				})
				.nullable()
				.optional(),
			udpLoadBalancer: z
				.object({
					servers: z.array(z.object({ address: z.string().trim() })).optional()
				})
				.nullable()
				.optional()
		})
		.refine(
			(data) =>
				!!data.loadBalancer?.servers?.length ||
				!!data.tcpLoadBalancer?.servers?.length ||
				!!data.udpLoadBalancer?.servers?.length,
			{
				message:
					"Exactly one of 'loadBalancer', 'tcpLoadBalancer', or 'udpLoadBalancer' must be provided.",
				path: ['loadBalancer', 'tcpLoadBalancer', 'udpLoadBalancer']
			}
		);

	export const validate = () => {
		try {
			formSchema.parse({
				serviceType: service.serviceType,
				loadBalancer: service.loadBalancer,
				tcpLoadBalancer: service.tcpLoadBalancer,
				udpLoadBalancer: service.udpLoadBalancer
			});

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
		if (service.tcpLoadBalancer === undefined) service.tcpLoadBalancer = { servers: [] };
		if (service.udpLoadBalancer === undefined) service.udpLoadBalancer = { servers: [] };

		switch (service.serviceType) {
			case 'http':
				service.loadBalancer.servers = servers.map((s) => {
					return { url: s };
				});
				break;
			case 'tcp':
				service.tcpLoadBalancer.servers = servers.map((s) => {
					return { address: s };
				});
				break;
			case 'udp':
				service.udpLoadBalancer.servers = servers.map((s) => {
					return { address: s };
				});
				break;
		}
	};

	const onSwitchChange = () => {
		passHostHeader = !passHostHeader;
		if (service.loadBalancer === undefined) service.loadBalancer = { servers: [] };
		service.loadBalancer.passHostHeader = passHostHeader;
	};

	onMount(() => {
		switch (service.serviceType) {
			case 'http':
				servers = service.loadBalancer?.servers?.map((s) => s.url ?? '') ?? [''];
				break;
			case 'tcp':
				servers = service.tcpLoadBalancer?.servers?.map((s) => s.address ?? '') ?? [''];
				break;
			case 'udp':
				servers = service.udpLoadBalancer?.servers?.map((s) => s.address ?? '') ?? [''];
				break;
		}
	});
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Service</Card.Title>
		<Card.Description>
			Make changes to your Service here. Click save when you're done.
		</Card.Description>
	</Card.Header>
	<Card.Content class="space-y-2">
		{#if service.serviceType === 'http'}
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
	</Card.Content>
</Card.Root>
