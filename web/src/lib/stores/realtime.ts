import { writable } from 'svelte/store';
import {
	ResourceType,
	EventType,
	type ProfileEvent,
	type ProfileEventsResponse,
	EventService
} from '$lib/gen/mantrae/v1/event_pb';
import { profile } from './profile';
import { ConnectError, createClient } from '@connectrpc/connect';
import type { Router } from '$lib/gen/mantrae/v1/router_pb';
import type { Service } from '$lib/gen/mantrae/v1/service_pb';
import type { Middleware } from '$lib/gen/mantrae/v1/middleware_pb';
import { createConnectTransport } from '@connectrpc/connect-web';
import { BASE_URL } from '$lib/api';

type ResourceUnion = Router | Service | Middleware;

interface RealtimeState<T extends ResourceUnion> {
	data: T[];
	rowCount: number;
	loading: boolean;
	error: string | null;
}

interface FetchResult<T> {
	data: T[];
	rowCount: number;
}

const eventClient = createClient(
	EventService,
	createConnectTransport({
		baseUrl: BASE_URL,
		useBinaryFormat: true,
		interceptors: [
			(next) => async (req) => {
				// Debug the request
				console.log('Request headers:', req.header);
				console.log('Request method:', req.method);
				return await next(req);
			}
		]
	})
);

export function createRealtimeStore<T extends ResourceUnion & { id: bigint }>(
	resourceType: ResourceType,
	fetchFn: (pageSize: number, pageIndex: number) => Promise<FetchResult<T>>
) {
	const { subscribe, set, update } = writable<RealtimeState<T>>({
		data: [],
		rowCount: 0,
		loading: false,
		error: null
	});

	let eventStream: AsyncIterable<ProfileEventsResponse> | null = null;
	let isStreaming = false;

	function getResourceFromEvent(event: ProfileEvent): T | null {
		if (!event.resource) return null;

		switch (resourceType) {
			case ResourceType.ROUTER:
				return event.resource.case === 'router' ? (event.resource.value as T) : null;
			case ResourceType.SERVICE:
				return event.resource.case === 'service' ? (event.resource.value as T) : null;
			case ResourceType.MIDDLEWARE:
				return event.resource.case === 'middleware' ? (event.resource.value as T) : null;
			default:
				return null;
		}
	}

	function isRelevantEvent(event: ProfileEvent): boolean {
		if (!event.resource || event.resourceType !== resourceType) {
			return false;
		}

		switch (resourceType) {
			case ResourceType.ROUTER:
				return event.resource.case === 'router';
			case ResourceType.SERVICE:
				return event.resource.case === 'service';
			case ResourceType.MIDDLEWARE:
				return event.resource.case === 'middleware';
			default:
				return false;
		}
	}

	async function startEventStream() {
		if (isStreaming || !profile?.id) return;

		isStreaming = true;
		try {
			eventStream = eventClient.profileEvents({
				profileId: profile.id,
				resourceTypes: [resourceType]
			});
			if (!eventStream) return;

			for await (const response of eventStream) {
				const event = response.event;
				if (!event || !isRelevantEvent(event)) continue;

				const resource = getResourceFromEvent(event);
				if (!resource) continue;

				update((state) => {
					const newData = [...state.data];

					switch (event.eventType) {
						case EventType.CREATED: {
							// Add to beginning of list, avoid duplicates
							const exists = newData.some((item) => item.id === resource.id);
							if (!exists) {
								newData.unshift(resource);
								return {
									...state,
									data: newData,
									rowCount: state.rowCount + 1
								};
							}
							return state;
						}

						case EventType.UPDATED: {
							const updateIndex = newData.findIndex((item) => item.id === resource.id);
							if (updateIndex !== -1) {
								newData[updateIndex] = resource;
								return { ...state, data: newData };
							}
							return state;
						}

						case EventType.DELETED: {
							const filteredData = newData.filter((item) => item.id !== resource.id);
							const rowCountDelta = newData.length - filteredData.length;
							return {
								...state,
								data: filteredData,
								rowCount: Math.max(0, state.rowCount - rowCountDelta)
							};
						}

						default:
							return state;
					}
				});
			}
		} catch (error) {
			const e = ConnectError.from(error);
			console.error('Event stream error:', e.message);
			update((state) => ({ ...state, error: e.message }));
		} finally {
			isStreaming = false;
		}
	}

	async function loadData(pageSize: number, pageIndex: number) {
		update((state) => ({ ...state, loading: true, error: null }));

		try {
			const result = await fetchFn(pageSize, pageIndex);
			set({
				data: result.data,
				rowCount: result.rowCount,
				loading: false,
				error: null
			});
		} catch (error) {
			const e = ConnectError.from(error);
			update((state) => ({
				...state,
				loading: false,
				error: e.message
			}));
		}
	}

	function stopEventStream() {
		isStreaming = false;
		eventStream = null;
	}

	return {
		subscribe,
		loadData,
		startEventStream,
		stopEventStream,
		// Optimistic updates for better UX
		optimisticUpdate: (id: bigint, updates: Partial<T>) => {
			update((state) => ({
				...state,
				data: state.data.map((item) => (item.id === id ? { ...item, ...updates } : item))
			}));
		}
	};
}

// Type-safe factory functions for each resource type
export function createRouterStore(
	fetchFn: (pageSize: number, pageIndex: number) => Promise<FetchResult<Router>>
) {
	return createRealtimeStore<Router>(ResourceType.ROUTER, fetchFn);
}

export function createServiceStore(
	fetchFn: (pageSize: number, pageIndex: number) => Promise<FetchResult<Service>>
) {
	return createRealtimeStore<Service>(ResourceType.SERVICE, fetchFn);
}

export function createMiddlewareStore(
	fetchFn: (pageSize: number, pageIndex: number) => Promise<FetchResult<Middleware>>
) {
	return createRealtimeStore<Middleware>(ResourceType.MIDDLEWARE, fetchFn);
}
