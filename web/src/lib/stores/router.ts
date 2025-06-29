import { eventStore } from "./events.svelte";
import {
	EventType,
	ResourceType,
	type ProfileEvent,
} from "$lib/gen/mantrae/v1/event_pb";
import type { Router } from "$lib/gen/mantrae/v1/router_pb";
import { routerClient } from "$lib/api";

class RouterStore {
	private data = $state<Router[]>([]);
	private rowCount = $state<number>(0);
	private unsubscribe?: () => void;
	private isDestroyed = false;

	get routers() {
		return this.data;
	}

	get totalCount() {
		return this.rowCount;
	}

	constructor() {
		// Subscribe to router events
		this.unsubscribe = eventStore.subscribe(ResourceType.ROUTER, (event) => {
			if (this.isDestroyed) return; // Prevent updates after destruction

			if (event.resource?.case === "router") {
				this.handleRouterEvent(event as ProfileEvent);
			}
		});
	}

	async loadPage(pageSize: number, pageIndex: number, profileId: bigint) {
		if (this.isDestroyed) return;

		try {
			const response = await routerClient.listRouters({
				profileId,
				limit: BigInt(pageSize),
				offset: BigInt(pageIndex * pageSize),
			});

			this.data = response.routers;
			this.rowCount = Number(response.totalCount);
		} catch (error) {
			console.error("Failed to load routers:", error);
			throw error;
		}
	}

	private handleRouterEvent(event: ProfileEvent) {
		if (this.isDestroyed || event.resource.case !== "router") return;

		const router = event.resource.value;

		switch (event.eventType) {
			case EventType.CREATED:
				// Only add if not already present (avoid duplicates)
				if (!this.data.find((r) => r.id === router.id)) {
					this.data = [...this.data, router];
					this.rowCount++;
				}
				break;
			case EventType.UPDATED:
				const index = this.data.findIndex((r) => r.id === router.id);
				if (index >= 0) {
					this.data[index] = router;
					this.data = [...this.data]; // Trigger reactivity
				}
				break;
			case EventType.DELETED:
				const beforeLength = this.data.length;
				this.data = this.data.filter((r) => r.id !== router.id);
				if (this.data.length < beforeLength) {
					this.rowCount--;
				}
				break;
		}
	}

	destroy() {
		this.unsubscribe?.();
		this.isDestroyed = true;
		this.unsubscribe = undefined;
		this.data = [];
		this.rowCount = 0;
	}
}

export const routerStore = new RouterStore();
