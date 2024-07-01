import { writable } from 'svelte/store';
import { API_URL } from '../config';

export interface MetricData {
    Duration: number;
    Success: boolean;
    Message: string;
    StatusCode: number;
    ContentLength: number;
    TLSVersion: string;
    CertExpiryDays: number;
}

export interface Metrics {
    [service: string]: {
        [endpoint: string]: MetricData;
    };
}

export const targets = writable<string[]>([]);
export const metrics = writable<Metrics>({});

export async function fetchTargets() {
    const response = await fetch(`${API_URL}/probe-metrics`);
    const data: Metrics = await response.json();
    const serviceNames = Object.keys(data);
    targets.set(serviceNames);
    metrics.set(data);
}

export async function fetchMetrics() {
    const response = await fetch(`${API_URL}/probe-metrics`);
    const data: Metrics = await response.json();
    metrics.set(data);
    return data;
}