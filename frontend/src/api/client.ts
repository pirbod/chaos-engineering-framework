import { fallbackExperiments, fallbackIntegrations, fallbackRisks, fallbackRunbooks, fallbackServices } from '../data/fallback';
import type { Experiment, IntegrationStatus, PlatformData, RiskResult, Runbook, Service } from '../types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';

async function getJSON<T>(path: string): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`);
  if (!response.ok) {
    throw new Error(`${path} returned ${response.status}`);
  }
  return response.json() as Promise<T>;
}

async function safe<T>(operation: () => Promise<T>, fallback: T): Promise<{ data: T; fallback: boolean }> {
  try {
    return { data: await operation(), fallback: false };
  } catch {
    return { data: fallback, fallback: true };
  }
}

export async function loadPlatformData(): Promise<PlatformData> {
  const [experimentsResult, servicesResult, runbooksResult, integrationsResult] = await Promise.all([
    safe(() => getJSON<Experiment[]>('/api/experiments'), fallbackExperiments),
    safe(() => getJSON<Service[]>('/api/services'), fallbackServices),
    safe(() => getJSON<Runbook[]>('/api/runbooks'), fallbackRunbooks),
    safe(() => getJSON<IntegrationStatus[]>('/api/integrations'), fallbackIntegrations),
  ]);

  const riskEntries = await Promise.all(
    servicesResult.data.map(async (service) => {
      const result = await safe(() => getJSON<RiskResult>(`/api/services/${service.id}/risk`), fallbackRisks[service.id]);
      return [service.id, result.data] as const;
    }),
  );

  return {
    experiments: experimentsResult.data,
    services: servicesResult.data,
    runbooks: runbooksResult.data,
    integrations: integrationsResult.data,
    risks: Object.fromEntries(riskEntries.filter(([, risk]) => Boolean(risk))),
    usingFallback: [experimentsResult, servicesResult, runbooksResult, integrationsResult].some((result) => result.fallback),
  };
}
