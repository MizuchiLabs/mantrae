import { api, BASE_URL } from '$lib/api';
import { user } from '$lib/stores/user';
import { source } from '$lib/stores/source';

interface Event {
	type: string;
	category: string;
	message: string;
}

let eventSourceInstance: EventSource | null = null;

export function initializeSSE() {
	if (eventSourceInstance || !user.isLoggedIn()) {
		return;
	}

	eventSourceInstance = new EventSource(`${BASE_URL}/events`);
	eventSourceInstance.onmessage = async (event) => {
		// Safeguard check
		if (!user.isLoggedIn()) return;

		try {
			const data: Event = JSON.parse(event.data);

			switch (data.category) {
				case 'profile':
					await api.listProfiles();
					break;
				case 'traefik':
					// Check if current source is valid before refetching
					if (!source.isValid(source.value)) return;
					await api.getTraefikConfig(source.value);
					break;
				case 'user':
					await api.listUsers();
					break;
				case 'dns':
					await api.listDNSProviders();
					break;
				case 'agent':
					await api.listAgentsByProfile();
					break;
				case 'error':
					await api.listErrors();
					break;
				default:
					break;
			}
		} catch (error) {
			console.error('Failed to parse SSE message or fetch data:', error);
		}
	};

	eventSourceInstance.onerror = (error) => {
		console.error('SSE connection error:', error);
		closeSSE(); // Close on error to allow potential re-initialization
	};

	eventSourceInstance.onopen = () => {
		console.log('SSE connection opened.');
	};
}

export function closeSSE() {
	if (eventSourceInstance) {
		console.log('Closing SSE connection...');
		eventSourceInstance.close();
		eventSourceInstance = null;
	}
}
