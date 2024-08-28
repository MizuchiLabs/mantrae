<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: Middleware;
	export let disabled = false;
	middleware.forwardAuth = {
		address: '',
		tls: { insecureSkipVerify: false, ca: '', cert: '', key: '' },
		trustForwardHeader: true,
		authResponseHeaders: [],
		authResponseHeadersRegex: '',
		authRequestHeaders: [],
		addAuthCookiesToResponse: [],
		...middleware.forwardAuth
	};
</script>

{#if middleware.forwardAuth}
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="address" class="text-right">Address</Label>
		<Input
			id="address"
			name="address"
			type="text"
			bind:value={middleware.forwardAuth.address}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="https://example.com/auth"
			{disabled}
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="trust-forward-header" class="text-right">Trust Forward Header</Label>
		<Switch
			id="trust-forward-header"
			bind:checked={middleware.forwardAuth.trustForwardHeader}
			class="col-span-3"
			{disabled}
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="auth-response-headers-regex" class="text-right">Auth Response Headers Regex</Label>
		<Input
			id="auth-response-headers-regex"
			name="auth-response-headers-regex"
			type="text"
			bind:value={middleware.forwardAuth.authResponseHeadersRegex}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="^X-"
			{disabled}
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="auth-request-headers" class="text-right">Auth Request Headers</Label>
		<Input
			id="auth-request-headers"
			name="auth-request-headers"
			type="text"
			bind:value={middleware.forwardAuth.authRequestHeaders}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="X-CustomHeader"
			{disabled}
		/>
	</div>
	<ArrayInput
		bind:items={middleware.forwardAuth.authResponseHeaders}
		label="Auth Response Headers"
		placeholder="X-Auth-User"
		{disabled}
	/>
	<ArrayInput
		bind:items={middleware.forwardAuth.authRequestHeaders}
		label="Auth Request Headers"
		placeholder="Accept"
		{disabled}
	/>
	<ArrayInput
		bind:items={middleware.forwardAuth.addAuthCookiesToResponse}
		label="Add Auth Cookies To Response"
		placeholder="Session-Cookie"
		{disabled}
	/>
{/if}
