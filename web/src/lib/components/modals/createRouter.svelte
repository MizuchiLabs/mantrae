<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import type { Selected } from 'bits-ui';
	import { routers, entrypoints, middlewares, upsertRouter } from '$lib/api';
	import { newRouter, newService, type Router } from '$lib/types/config';
	import RuleEditor from '../utils/ruleEditor.svelte';
	import Service from '../forms/service.svelte';

	let router = newRouter();
	let service = newService();

	const create = async () => {
		if (router.name === '' || isNameTaken) return;
		await upsertRouter(router.name, router, service);

		router = newRouter();
		service = newService();
	};

	let routerType: Selected<string> | undefined = { label: 'HTTP', value: 'http' };
	const changeType = (serviceType: Selected<string> | undefined) => {
		if (serviceType === undefined) return;
		router = newRouter();
		service = newService();
		router.routerType = serviceType.value;
		service.serviceType = serviceType.value;
		routerType = { label: serviceType.label || '', value: serviceType.value };
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

	// Check if router name is taken
	let isNameTaken = false;
	$: isNameTaken = $routers.some(
		(r) => r.name.split('@')[0].toLowerCase() === router.name.split('@')[0].toLowerCase()
	);
</script>

<Dialog.Root>
	<Dialog.Trigger>
		<div class="flex w-full flex-row items-center justify-between">
			<Button class="flex items-center gap-2 bg-red-400 text-black">
				<span>New Router</span>
				<iconify-icon icon="fa6-solid:plus" />
			</Button>
		</div>
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
							<Label for="current" class="text-right">Type</Label>
							<Select.Root onSelectedChange={changeType} selected={routerType}>
								<Select.Trigger class="col-span-3">
									<Select.Value placeholder="Select a type" />
								</Select.Trigger>
								<Select.Content>
									<Select.Item value="http" label="HTTP">HTTP</Select.Item>
									<Select.Item value="tcp" label="TCP">TCP</Select.Item>
									<Select.Item value="udp" label="UDP">UDP</Select.Item>
								</Select.Content>
							</Select.Root>
						</div>
						<div class="grid grid-cols-4 items-center gap-4">
							<Label for="name" class="text-right">Name</Label>
							<Input
								id="name"
								name="name"
								type="text"
								bind:value={router.name}
								class={isNameTaken
									? 'col-span-3 border-red-400 focus-visible:ring-0 focus-visible:ring-offset-0'
									: 'col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0'}
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
									{#each $entrypoints as entrypoint}
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
									{#each $middlewares as middleware}
										{#if router.routerType === middleware.middlewareType}
											<Select.Item value={middleware.name}>
												{middleware.name}
											</Select.Item>
										{/if}
									{/each}
								</Select.Content>
							</Select.Root>
						</div>
						<!-- Insane hacky editor for traefik rules -->
						<div class:hidden={router.routerType === 'udp'}>
							<RuleEditor bind:rule={router.rule} />
						</div>
					</Card.Content>
				</Card.Root>
			</Tabs.Content>
			<Tabs.Content value="service">
				<Service bind:service />
			</Tabs.Content>
		</Tabs.Root>
		<Dialog.Close class="w-full">
			<Button type="submit" class="w-full" on:click={() => create()}>Save</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
