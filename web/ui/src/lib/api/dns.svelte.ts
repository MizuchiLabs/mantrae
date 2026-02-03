import { toast } from 'svelte-sonner';
import { useMutation, useQuery } from '$lib/query';
import { DNSProviderService } from '$lib/gen/mantrae/v1/dns_provider_pb';

export const dns = {
	// Queries
	get: (id: string) =>
		useQuery(
			DNSProviderService.method.getDNSProvider,
			{ id },
			{ select: (res) => res.dnsProvider }
		),
	list: () =>
		useQuery(DNSProviderService.method.listDNSProviders, {}, { select: (res) => res.dnsProviders }),

	// Mutations
	create: () =>
		useMutation(DNSProviderService.method.createDNSProvider, {
			onSuccess: () => toast.success('DNS provider created!')
		}),
	update: () =>
		useMutation(DNSProviderService.method.updateDNSProvider, {
			onSuccess: () => toast.success('DNS provider updated!')
		}),
	delete: () =>
		useMutation(DNSProviderService.method.deleteDNSProvider, {
			onSuccess: () => toast.success('DNS provider deleted!')
		})
};
