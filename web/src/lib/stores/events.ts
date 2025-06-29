import { eventClient } from "$lib/api";
import type {
	GlobalEvent,
	ProfileEvent,
	ResourceType,
	StreamGlobalEventsResponse,
	StreamProfileEventsResponse,
} from "$lib/gen/mantrae/v1/event_pb";
import { ConnectError } from "@connectrpc/connect";

class EventStore {
	private profileStream: AbortController | null = null;
	private globalStream: AbortController | null = null;
	private listeners = new Map<
		ResourceType,
		Set<(event: ProfileEvent | GlobalEvent) => void>
	>();
	private currentProfileId?: bigint;
	private currentResourceTypes?: ResourceType[];

	subscribe(
		resourceType: ResourceType,
		callback: (event: ProfileEvent | GlobalEvent) => void,
	) {
		if (!this.listeners.has(resourceType)) {
			this.listeners.set(resourceType, new Set());
		}
		this.listeners.get(resourceType)!.add(callback);

		// Return unsubscribe function
		return () => {
			this.listeners.get(resourceType)?.delete(callback);
		};
	}

	async startProfileEvents(profileId: bigint, resourceTypes?: ResourceType[]) {
		this.currentProfileId = profileId;
		this.currentResourceTypes = resourceTypes;

		if (this.profileStream) this.stopProfileEvents();

		await this.connectProfileStream();
	}

	private async connectProfileStream() {
		if (!this.currentProfileId) return;

		this.profileStream = new AbortController();

		try {
			const stream = eventClient.streamProfileEvents(
				{
					profileId: this.currentProfileId,
					resourceTypes: this.currentResourceTypes || [],
				},
				{ signal: this.profileStream.signal },
			);

			for await (const event of stream) {
				this.handleProfileEvent(event);
			}
		} catch (error) {
			const e = ConnectError.from(error);
			if (!this.profileStream?.signal.aborted) {
				console.error("Profile stream error:", e.message);
			}
		}
	}

	private handleProfileEvent(res: StreamProfileEventsResponse) {
		if (!res.event) return;
		const listeners = this.listeners.get(res.event.resourceType);
		if (listeners) {
			listeners.forEach((callback) => callback(res.event));
		}
	}

	stopProfileEvents() {
		this.profileStream?.abort();
		this.profileStream = null;
	}

	async startGlobalEvents(resourceTypes?: ResourceType[]) {
		if (this.globalStream) this.stopGlobalEvents();

		this.globalStream = new AbortController();

		try {
			const stream = eventClient.streamGlobalEvents(
				{ resourceTypes: resourceTypes || [] },
				{ signal: this.globalStream.signal },
			);

			for await (const event of stream) {
				this.handleGlobalEvent(event);
			}
		} catch (error) {
			const e = ConnectError.from(error);
			if (!this.globalStream?.signal.aborted) {
				console.error("Global stream error:", e.message);
			}
		}
	}

	private handleGlobalEvent(res: StreamGlobalEventsResponse) {
		if (!res.event) return;
		const listeners = this.listeners.get(res.event.resourceType);
		if (listeners) {
			listeners.forEach((callback) => callback(res.event));
		}
	}

	stopGlobalEvents() {
		this.globalStream?.abort();
		this.globalStream = null;
	}

	destroy() {
		this.stopProfileEvents();
		this.stopGlobalEvents();
		this.listeners.clear();
	}
}

export const eventStore = new EventStore();
