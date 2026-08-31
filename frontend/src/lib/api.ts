const BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api';

export interface User {
	id: number;
	full_name: string;
	email: string;
	phone: string;
	role: 'admin' | 'staff' | 'customer';
	active: boolean;
	created_at: string;
}

export interface Resource {
	id: number;
	name: string;
	type: 'doctor' | 'room' | 'other';
	description: string;
	location: string;
	capacity: number;
	active: boolean;
}

export interface Reservation {
	id: number;
	status: 'pending' | 'confirmed' | 'cancelled' | 'completed';
	start_time: string;
	end_time: string;
	notes: string;
	user: User;
	resource: Resource;
	created_at: string;
}

type Params = Record<string, string>;

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
	const res = await fetch(`${BASE_URL}${path}`, {
		headers: { 'Content-Type': 'application/json' },
		...options
	});

	if (!res.ok) {
		let message = `Error ${res.status}`;
		try {
			const body = await res.json();
			if (body && typeof body.error === 'string') message = body.error;
		} catch {
			// ignore
		}
		throw new Error(message);
	}

	if (res.status === 204) return undefined as T;
	return (await res.json()) as T;
}

export const api = {
	get: <T>(path: string, params?: Params) => {
		const query = params && Object.keys(params).length
			? `?${new URLSearchParams(params).toString()}`
			: '';
		return request<T>(`${path}${query}`);
	},
	post: <T>(path: string, body: unknown) =>
		request<T>(path, { method: 'POST', body: JSON.stringify(body) }),
	put: <T>(path: string, body: unknown) =>
		request<T>(path, { method: 'PUT', body: JSON.stringify(body) }),
	patch: <T>(path: string, body: unknown) =>
		request<T>(path, { method: 'PATCH', body: JSON.stringify(body) }),
	del: (path: string) => request<void>(path, { method: 'DELETE' })
};
