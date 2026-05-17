export type Experiment = {
  id: string;
  name: string;
  type: string;
  target: string;
  description: string;
  blastRadius: 'low' | 'medium' | 'high' | 'critical';
  defaultDuration: string;
  safetyNotes: string[];
  runbookId: string;
  lastOutcome: string;
  tags: string[];
};

export type Service = {
  id: string;
  name: string;
  description: string;
  owner: string;
  team: string;
  tier: string;
  criticality: string;
  recentErrorRate: number;
  ownershipCompleteness: number;
  runbookIds: string[];
  observabilityMaturity: number;
  lastExperimentOutcome: string;
  securitySensitivity: string;
  slo: string;
  experimentIds: string[];
  repository: string;
  dashboardUrl: string;
};

export type Runbook = {
  id: string;
  title: string;
  summary: string;
  appliesTo: string[];
  severity: string;
  steps: string[];
  checks: string[];
  rollback: string[];
  links: Array<{ label: string; url: string }>;
};

export type RiskFactor = {
  name: string;
  contribution: number;
  weight: number;
  recommendation: string;
};

export type RiskResult = {
  serviceId: string;
  score: number;
  level: 'low' | 'medium' | 'high' | 'critical';
  topContributingFactors: RiskFactor[];
  recommendedActions: string[];
};

export type IntegrationStatus = {
  id: string;
  name: string;
  status: string;
  category: string;
  details: string;
  signals: string[];
  docsUrl: string;
  lastChecked: string;
};

export type PlatformData = {
  experiments: Experiment[];
  services: Service[];
  runbooks: Runbook[];
  integrations: IntegrationStatus[];
  risks: Record<string, RiskResult>;
  usingFallback: boolean;
};
