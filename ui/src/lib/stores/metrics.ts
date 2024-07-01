import { writable } from 'svelte/store';
import { API_URL } from '../config';

export const metrics = writable<{ [key: string]: number }>({});

export async function fetchMetrics() {
    const response = await fetch(`${API_URL}/probe-metrics`);
    const data = await response.json();
    metrics.set(data);
}
