export interface DNSProvider {
	id: number;
	name: string;
	type: string;
	external_ip: string;
	api_key?: string;
	api_url?: string;
	is_active: boolean;
}

export function newProvider(): DNSProvider {
	return {
		id: 0,
		name: '',
		type: 'cloudflare',
		external_ip: '',
		api_key: '',
		api_url: '',
		is_active: false
	};
}
