export interface Provider {
	id: number;
	name: string;
	type: string;
	external_ip: string;
	api_key?: string;
	api_url?: string;
}

export function newProvider(): Provider {
	return {
		id: 0,
		name: '',
		type: 'cloudflare',
		external_ip: '',
		api_key: '',
		api_url: ''
	};
}
