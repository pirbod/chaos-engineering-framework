import React from 'react';

type Props = {
  serviceName: string;
  owner: string;
  maturityScore: number;
  experiments: Array<{ id: string; name: string; blastRadius: string }>;
  runbooks: Array<{ id: string; title: string; href: string }>;
};

export function PlatformChaosHubPage({ serviceName, owner, maturityScore, experiments, runbooks }: Props) {
  return (
    <section>
      <header>
        <p>Platform Chaos Reliability Hub</p>
        <h1>{serviceName}</h1>
        <p>Owner: {owner || 'Missing owner metadata'}</p>
        <strong>Maturity score: {maturityScore}/100</strong>
      </header>
      <h2>Available experiments</h2>
      <ul>
        {experiments.map((experiment) => (
          <li key={experiment.id}>
            {experiment.name} / {experiment.blastRadius}
          </li>
        ))}
      </ul>
      <h2>Runbooks</h2>
      <ul>
        {runbooks.map((runbook) => (
          <li key={runbook.id}>
            <a href={runbook.href}>{runbook.title}</a>
          </li>
        ))}
      </ul>
    </section>
  );
}
