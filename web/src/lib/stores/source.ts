import { createLocalStorage } from '$lib/storage.svelte';
import { TraefikSource } from '$lib/types';

class TraefikSourceStore {
	private store = createLocalStorage<TraefikSource>('selected_source', TraefikSource.API);

	get value(): TraefikSource {
		return this.store.value ?? TraefikSource.API;
	}

	set value(source: TraefikSource) {
		if (this.isValid(source)) {
			this.store.value = source;
		}
	}

	isValid(source: string): source is TraefikSource {
		return Object.values(TraefikSource).includes(source as TraefikSource);
	}

	isApi(): boolean {
		return this.value === TraefikSource.API;
	}

	isLocal(): boolean {
		return this.value === TraefikSource.LOCAL;
	}

	isAgent(): boolean {
		return this.value === TraefikSource.AGENT;
	}

	// Reset to default
	reset(): void {
		this.store.value = TraefikSource.API;
	}
}

export const source = new TraefikSourceStore();
