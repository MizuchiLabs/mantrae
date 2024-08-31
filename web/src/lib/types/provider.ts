export interface Provider {
	name: string;
	type: string;
	externalIP: string;
	key?: string;
	url?: string;
}

export function newProvider(): Provider {
	return {
		name: '',
		type: 'cloudflare',
		externalIP: '',
		key: '',
		url: ''
	};
}
