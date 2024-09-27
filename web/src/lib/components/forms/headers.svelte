<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';
	import ObjectInput from '../ui/object-input/object-input.svelte';
	import { z } from 'zod';
	import { CustomIPSchemaOptional } from '../utils/validation';

	export let middleware: Middleware;
	const emptyHeaders = {
		// SSL and Security Headers (commonly used)
		sslProxyHeaders: {},

		// Security and Privacy Policies (high importance)
		contentSecurityPolicy: '',
		contentTypeNosniff: false,
		browserXssFilter: false,
		frameDeny: false,
		customFrameOptionsValue: '',
		referrerPolicy: '',
		permissionsPolicy: '',

		// Access Control Headers (important for CORS and security)
		accessControlAllowOriginList: [],
		accessControlAllowOriginListRegex: [],
		accessControlAllowHeaders: [],
		accessControlAllowMethods: [],
		accessControlAllowCredentials: false,
		accessControlExposeHeaders: [],

		// STS (HTTP Strict Transport Security)
		stsIncludeSubdomains: false,
		stsPreload: false,
		forceSTSHeader: false,

		// Custom Headers (for custom configurations)
		customRequestHeaders: {},
		customResponseHeaders: {},

		// Less frequently used security options
		addVaryHeader: false,
		allowedHosts: [],
		hostsProxyHeaders: [],
		publicKey: '',

		// Miscellaneous
		customBrowserXSSValue: ''
	};

	const defaultTemplate = {
		// SSL and Security Headers
		sslProxyHeaders: {
			'X-Forwarded-Proto': 'https'
		},

		// Security and Privacy Policies
		contentSecurityPolicy:
			"default-src 'self'; script-src 'self'; object-src 'none'; style-src 'self' 'unsafe-inline'; frame-ancestors 'none';", // Mitigates XSS attacks
		contentTypeNosniff: true, // Prevents MIME-type sniffing
		browserXssFilter: true, // Helps prevent XSS attacks
		frameDeny: true, // Denies embedding in iframes
		customFrameOptionsValue: '', // Can be set for more granular control
		referrerPolicy: 'no-referrer', // Prevents referrer leakage
		permissionsPolicy: 'geolocation=(), microphone=(), camera=(), fullscreen=(self)', // Restricts access to sensitive APIs

		// Access Control Headers (CORS and Security)
		accessControlAllowHeaders: ['Authorization', 'Content-Type'],
		accessControlAllowMethods: ['GET', 'POST', 'OPTIONS'],
		accessControlAllowCredentials: true, // Allow sending credentials
		accessControlExposeHeaders: ['Authorization'],
		accessControlMaxAge: 86400, // Cache CORS preflight requests for 1 day

		// STS (HTTP Strict Transport Security)
		stsSeconds: 31536000, // Enforce HTTPS for 1 year
		stsIncludeSubdomains: true, // Apply STS to subdomains
		stsPreload: true, // Preload into browsers for STS
		forceSTSHeader: true, // Force the STS header

		// Custom Headers (for custom configurations)
		customResponseHeaders: {
			'X-Content-Type-Options': 'nosniff',
			'X-Frame-Options': 'DENY',
			'X-XSS-Protection': '1; mode=block',
			'X-Robots-Tag': 'none,noarchive,nosnippet,notranslate,noimageindex'
		},
		customRequestHeaders: {
			'X-Forwarded-Proto': 'https',
			'X-Permitted-Cross-Domain-Policies': 'none'
		},

		// Less frequently used security options
		addVaryHeader: true, // Useful for caching and negotiation
		hostsProxyHeaders: ['X-Forwarded-Host']
	};

	const headersSchema = z.object({
		sslProxyHeaders: z.record(z.string(), z.string()).optional(),
		contentSecurityPolicy: z.string().optional(),
		contentTypeNosniff: z.boolean().optional(),
		browserXssFilter: z.boolean().optional(),
		frameDeny: z.boolean().optional(),
		customFrameOptionsValue: z.string().optional(),
		referrerPolicy: z.string().optional(),
		permissionsPolicy: z.string().optional(),
		accessControlAllowOriginList: z.array(CustomIPSchemaOptional).default([]).optional(),
		accessControlAllowOriginListRegex: z.array(z.string()).default([]).optional(),
		accessControlAllowHeaders: z.array(z.string()).default([]).optional(),
		accessControlAllowMethods: z.array(z.string()).default([]).optional(),
		accessControlAllowCredentials: z.boolean().default(false).optional(),
		accessControlExposeHeaders: z.array(z.string()).default([]).optional(),
		accessControlMaxAge: z.coerce.number().int().nonnegative().optional(),
		stsSeconds: z.coerce.number().int().nonnegative().optional(),
		stsIncludeSubdomains: z.boolean().default(false).optional(),
		stsPreload: z.boolean().default(false).optional(),
		forceSTSHeader: z.boolean().default(false).optional(),
		customResponseHeaders: z.record(z.string(), z.string()).optional(),
		customRequestHeaders: z.record(z.string(), z.string()).optional(),
		addVaryHeader: z.boolean().default(true).optional(),
		hostsProxyHeaders: z.array(z.string()).default([]).optional()
	});
	middleware.headers = headersSchema.parse({ ...middleware.headers });

	let errors: Record<any, string[] | undefined> = {};
	const validate = () => {
		try {
			middleware.headers = headersSchema.parse(middleware.headers);
			console.log(middleware.headers.stsSeconds);
			errors = {};
		} catch (err) {
			if (err instanceof z.ZodError) {
				errors = err.flatten().fieldErrors;
			}
		}
	};

	//middleware.headers = { ...emptyHeaders, ...middleware.headers };
	let isTemplate = false;
	const toggleTemplate = () => {
		isTemplate = !isTemplate;
		if (isTemplate) {
			middleware.headers = { ...emptyHeaders, ...defaultTemplate };
		} else {
			middleware.headers = { ...emptyHeaders };
		}
		validate();
	};
</script>

<div class="flex items-center justify-end gap-2">
	<Button on:click={toggleTemplate}>
		{isTemplate ? 'Clear Config' : 'Use Secure Template'}
	</Button>
</div>

{#if middleware.headers}
	<!-- SSL and Security Headers (commonly used) -->
	<ObjectInput
		bind:items={middleware.headers.sslProxyHeaders}
		label="SSL Proxy Headers"
		keyPlaceholder="Header Name"
		valuePlaceholder="Header Value"
		on:update={validate}
		class="my-4"
	/>
	{#if errors.sslProxyHeaders}
		<span class="text-sm text-red-500">{errors.sslProxyHeaders}</span>
	{/if}

	<!-- Security and Privacy Policies (high importance) -->
	<span class="my-4 border-b border-gray-200 pb-2 font-bold">Security and Privacy Policies</span>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="content-security-policy" class="text-right">Content Security Policy</Label>
		<Input
			id="content-security-policy"
			name="content-security-policy"
			type="text"
			bind:value={middleware.headers.contentSecurityPolicy}
			on:input={validate}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="default-src 'self'; script-src 'self' 'unsafe-inline';"
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="custom-frame-options-value" class="text-right">Custom Frame Options Value</Label>
		<Input
			id="custom-frame-options-value"
			name="custom-frame-options-value"
			type="text"
			bind:value={middleware.headers.customFrameOptionsValue}
			on:input={validate}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="SAMEORIGIN"
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="referrer-policy" class="text-right">Referrer Policy</Label>
		<Input
			id="referrer-policy"
			name="referrer-policy"
			type="text"
			bind:value={middleware.headers.referrerPolicy}
			on:input={validate}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="no-referrer"
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="permissions-policy" class="text-right">Permissions Policy</Label>
		<Input
			id="permissions-policy"
			name="permissions-policy"
			type="text"
			bind:value={middleware.headers.permissionsPolicy}
			on:input={validate}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="geolocation 'none'; microphone 'none';"
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="content-type-no-sniff" class="text-right">Content Type No Sniff</Label>
		<Switch
			id="content-type-no-snuff"
			bind:checked={middleware.headers.contentTypeNosniff}
			class="col-span-3"
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="browser-xss-filter" class="text-right">Browser XSS Filter</Label>
		<Switch
			id="browser-xss-filter"
			bind:checked={middleware.headers.browserXssFilter}
			class="col-span-3"
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="frame-deny" class="text-right">Frame Deny</Label>
		<Switch id="frame-deny" bind:checked={middleware.headers.frameDeny} class="col-span-3" />
	</div>

	<!-- Access Control Headers -->
	<span class="my-4 border-b border-gray-200 pb-2 font-bold">Access Control Headers</span>
	<ArrayInput
		bind:items={middleware.headers.accessControlAllowOriginList}
		placeholder="*"
		label="Access Control Allow Origin List"
		on:update={validate}
		class="my-2"
	/>
	<ArrayInput
		bind:items={middleware.headers.accessControlAllowOriginListRegex}
		placeholder="example\\.com"
		label="Access Control Allow Origin List Regex"
		class="my-2"
	/>
	<ArrayInput
		bind:items={middleware.headers.accessControlAllowHeaders}
		placeholder="Authorization"
		label="Access Control Allow Headers"
		on:update={validate}
		class="my-2"
	/>
	<ArrayInput
		bind:items={middleware.headers.accessControlAllowMethods}
		placeholder="GET, POST, PUT, DELETE, OPTIONS"
		label="Access Control Allow Methods"
		on:update={validate}
		class="my-2"
	/>
	<ArrayInput
		bind:items={middleware.headers.accessControlExposeHeaders}
		placeholder="Authorization"
		label="Access Control Expose Headers"
		on:update={validate}
		class="my-2"
	/>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="access-control-max-age" class="text-right">Access Control Max Age</Label>
		<Input
			id="access-control-max-age"
			name="access-control-max-age"
			type="number"
			bind:value={middleware.headers.accessControlMaxAge}
			on:input={validate}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="0"
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="access-control-allow-credentials" class="text-right"
			>Access Control Allow Credentials</Label
		>
		<Switch
			id="access-control-allow-credentials"
			bind:checked={middleware.headers.accessControlAllowCredentials}
			class="col-span-3"
		/>
	</div>

	<!-- STS -->
	<span class="my-4 border-b border-gray-200 pb-2 font-bold">Strict Transport Security (STS)</span>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="sts-seconds" class="text-right">STS Seconds</Label>
		<Input
			id="sts-seconds"
			name="sts-seconds"
			type="number"
			bind:value={middleware.headers.stsSeconds}
			on:input={validate}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="86400"
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="sts-include-sub-domains" class="text-right">STS Include Sub Domains</Label>
		<Switch
			id="sts-include-sub-domains"
			bind:checked={middleware.headers.stsIncludeSubdomains}
			class="col-span-3"
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="sts-pre-load" class="text-right">STS Pre Load</Label>
		<Switch id="sts-pre-load" bind:checked={middleware.headers.stsPreload} class="col-span-3" />
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="force-sts-header" class="text-right">Force STS Header</Label>
		<Switch
			id="force-sts-header"
			bind:checked={middleware.headers.forceSTSHeader}
			class="col-span-3"
		/>
	</div>

	<!-- Custom Headers -->
	<span class="my-4 border-b border-gray-200 pb-2 font-bold">Custom Headers</span>
	<ObjectInput
		bind:items={middleware.headers.customResponseHeaders}
		label="Custom Response Headers"
		keyPlaceholder="Header Name"
		valuePlaceholder="Header Value"
		on:update={validate}
		class="my-2"
	/>
	<ObjectInput
		bind:items={middleware.headers.customRequestHeaders}
		label="Custom Request Headers"
		keyPlaceholder="Header Name"
		valuePlaceholder="Header Value"
		on:update={validate}
		class="my-2"
	/>

	<!-- Less frequently used headers -->
	<span class="my-4 border-b border-gray-200 pb-2 font-bold">Various extra headers</span>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="add-vary-header" class="text-right">Add Vary Header</Label>
		<Switch
			id="add-vary-header"
			bind:checked={middleware.headers.addVaryHeader}
			class="col-span-3"
		/>
	</div>
	<ArrayInput
		bind:items={middleware.headers.allowedHosts}
		placeholder="example.com"
		label="Allowed Hosts"
		on:update={validate}
		class="my-2"
	/>
	<ArrayInput
		bind:items={middleware.headers.hostsProxyHeaders}
		placeholder="X-Forwarded-Host"
		label="Hosts Proxy Headers"
		on:update={validate}
		class="my-2"
	/>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="public-key" class="text-right">Public Key</Label>
		<Input
			id="public-key"
			name="public-key"
			type="text"
			bind:value={middleware.headers.publicKey}
			on:input={validate}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="MIIBIjANBgkqhkiG9w0BAQEFAA..."
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="custom-browser-xss-value" class="text-right">Custom Browser XSS Value</Label>
		<Input
			id="custom-browser-xss-value"
			name="custom-browser-xss-value"
			type="text"
			bind:value={middleware.headers.customBrowserXSSValue}
			on:input={validate}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="1; mode=block"
		/>
	</div>
{/if}
