<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { activeProfile, updateRouter, updateService } from '$lib/api';
	import { newService, type Router, type Service } from '$lib/types/config';
	import RuleEditor from '../utils/ruleEditor.svelte';
	import type { Selected } from 'bits-ui';
	import { toast } from 'svelte-sonner';

	export let router: Router;
	let oldRouter = router.name;
	let oldService = router.service + '@' + router.provider;

	$: middlewares = Object.values($activeProfile?.instance?.dynamic?.middlewares || []);
	let service: Service | undefined =
		$activeProfile?.instance?.dynamic?.services?.[router.service + '@' + router.provider];
	$: servers = service?.loadBalancer?.servers?.length || 0;

	const update = async () => {
		if (service === undefined) return;
		router.name = router.service + '@' + router.provider;
		service.name = router.service + '@' + router.provider; // Extra check in case router name changed
		service.serviceType = router.routerType;
		try {
			await updateRouter($activeProfile.name, router, oldRouter);
			await updateService($activeProfile.name, service, oldService);
			toast.success(`Router ${router.name} updated`);
		} catch (e) {}
	};

	const toggleEntrypoint = (router: Router, item: Selected<unknown>[] | undefined) => {
		if (item === undefined) return;
		router.entrypoints = item.map((i) => i.value) as string[];
	};
	const toggleMiddleware = (router: Router, item: Selected<unknown>[] | undefined) => {
		if (item === undefined) return;
		router.middlewares = item.map((i) => i.value) as string[];
	};
	const getSelectedEntrypoints = (router: Router): Selected<unknown>[] => {
		let list = router?.entrypoints?.map((entrypoint) => {
			return { value: entrypoint, label: entrypoint };
		});
		return list ?? [];
	};
	const getSelectedMiddlewares = (router: Router): Selected<unknown>[] => {
		let list = router?.middlewares?.map((middleware) => {
			return { value: middleware, label: middleware };
		});
		return list ?? [];
	};

	const addServer = () => {
		if (service?.loadBalancer?.servers === undefined) {
			service = newService();
			return;
		}
		service.loadBalancer.servers = [...service.loadBalancer.servers, { url: '' }];
	};
	const removeServer = (index: number) => {
		if (service?.loadBalancer?.servers === undefined) return;
		if (service.loadBalancer.servers.length > 1) {
			service.loadBalancer.servers = service.loadBalancer.servers.filter((_, i) => i !== index);
		}
	};
</script>

<Dialog.Root>
	<Dialog.Trigger>
		<Button variant="ghost" class="h-8 w-4 rounded-full bg-orange-400">
			<iconify-icon icon="fa6-solid:pencil" />
		</Button>
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-[520px]">
		<Tabs.Root value="router" class="mt-4 w-[470px]">
			<Tabs.List class="grid w-full grid-cols-2">
				<Tabs.Trigger value="router">Router</Tabs.Trigger>
				<Tabs.Trigger value="service">Service</Tabs.Trigger>
			</Tabs.List>
			<Tabs.Content value="router">
				<Card.Root>
					<Card.Header>
						<Card.Title>Router</Card.Title>
						<Card.Description>
							Make changes to your Router here. Click save when you're done.
						</Card.Description>
					</Card.Header>
					<Card.Content class="space-y-2">
						<div class="grid grid-cols-4 items-center gap-4">
							<Label for="name" class="text-right">Name</Label>
							<Input
								id="name"
								name="name"
								type="text"
								class="col-span-3"
								bind:value={router.service}
								placeholder="Name of the router"
								required
							/>
						</div>
						<div class="grid grid-cols-4 items-center gap-4">
							<Label for="entrypoints" class="text-right">Entrypoints</Label>
							<Select.Root
								multiple={true}
								selected={getSelectedEntrypoints(router)}
								onSelectedChange={(value) => toggleEntrypoint(router, value)}
							>
								<Select.Trigger class="col-span-3">
									<Select.Value placeholder="Select an entrypoint" />
								</Select.Trigger>
								<Select.Content>
									{#each $activeProfile?.instance?.dynamic?.entrypoints || [] as entrypoint}
										<Select.Item value={entrypoint.name}>
											<div class="flex flex-row items-center gap-2">
												{entrypoint.name}
												{#if entrypoint.http}
													{#if 'tls' in entrypoint.http}
														<iconify-icon icon="fa6-solid:lock" class=" text-green-400" />
													{/if}
												{/if}
											</div>
										</Select.Item>
									{/each}
								</Select.Content>
							</Select.Root>
						</div>
						<div
							class="grid grid-cols-4 items-center gap-4"
							class:hidden={router.routerType === 'udp'}
						>
							<Label for="middlewares" class="text-right">Middlewares</Label>
							<Select.Root
								multiple={true}
								selected={getSelectedMiddlewares(router)}
								onSelectedChange={(value) => toggleMiddleware(router, value)}
							>
								<Select.Trigger class="col-span-3">
									<Select.Value placeholder="Select a middleware" />
								</Select.Trigger>
								<Select.Content>
									{#each middlewares as middleware}
										{#if router.routerType === middleware.type}
											<Select.Item value={middleware.name}>
												{middleware.name}
											</Select.Item>
										{/if}
									{/each}
								</Select.Content>
							</Select.Root>
						</div>
						<div class:hidden={router.routerType === 'udp'}>
							<RuleEditor bind:rule={router.rule} />
						</div>
					</Card.Content>
				</Card.Root>
			</Tabs.Content>
			<Tabs.Content value="service">
				<Card.Root>
					<Card.Header>
						<Card.Title>Service</Card.Title>
						<Card.Description>
							Make changes to your Service here. Click save when you're done.
						</Card.Description>
					</Card.Header>
					<Card.Content class="space-y-2">
						{#if router.routerType === 'http' && service?.loadBalancer !== undefined}
							<div class="flex items-center space-x-2">
								<Label for="passHostHeader">Pass Host Header</Label>
								<Switch id="passHostHeader" bind:checked={service.loadBalancer.passHostHeader} />
							</div>
						{/if}
						<div class="space-y-1">
							<div class="flex flex-row items-center justify-between">
								<Label for="url">Load Balancer URL</Label>
								<Button class="h-8 w-4 bg-red-400 text-black" on:click={() => addServer()}>
									<iconify-icon icon="fa6-solid:plus" />
								</Button>
							</div>
							{#each service?.loadBalancer?.servers || [] as server, idx}
								<div class="flex items-center gap-2">
									<Input
										id="url"
										type="text"
										bind:value={server.url}
										class="focus-visible:ring-0 focus-visible:ring-offset-0"
									/>
									{#if servers > 1 && idx >= 1}
										<Button on:click={() => removeServer(idx)}>-</Button>
									{/if}
								</div>
							{/each}
						</div>
					</Card.Content>
				</Card.Root>
			</Tabs.Content>
		</Tabs.Root>
		<Dialog.Close class="w-full">
			<Button class="w-full" on:click={() => update()}>Save</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
