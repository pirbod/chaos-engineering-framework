import { useMemo, useState } from 'react';
import { EmptyState } from '../components/EmptyState';
import { RiskBadge } from '../components/RiskBadge';
import type { Experiment } from '../types';

type Props = {
  experiments: Experiment[];
};

export function ExperimentCatalog({ experiments }: Props) {
  const [typeFilter, setTypeFilter] = useState('all');
  const [riskFilter, setRiskFilter] = useState('all');
  const [targetFilter, setTargetFilter] = useState('');

  const types = Array.from(new Set(experiments.map((experiment) => experiment.type))).sort();
  const filtered = useMemo(
    () =>
      experiments.filter((experiment) => {
        const matchesType = typeFilter === 'all' || experiment.type === typeFilter;
        const matchesRisk = riskFilter === 'all' || experiment.blastRadius === riskFilter;
        const matchesTarget = targetFilter === '' || experiment.target.toLowerCase().includes(targetFilter.toLowerCase());
        return matchesType && matchesRisk && matchesTarget;
      }),
    [experiments, riskFilter, targetFilter, typeFilter],
  );

  return (
    <section className="page">
      <header className="page-header compact">
        <div>
          <p className="eyebrow">Experiment catalog</p>
          <h1>Choose safe, scoped resilience tests with clear guardrails.</h1>
        </div>
      </header>

      <div className="toolbar" aria-label="Experiment filters">
        <label>
          Type
          <select value={typeFilter} onChange={(event) => setTypeFilter(event.target.value)}>
            <option value="all">All types</option>
            {types.map((type) => (
              <option key={type} value={type}>
                {type}
              </option>
            ))}
          </select>
        </label>
        <label>
          Blast radius
          <select value={riskFilter} onChange={(event) => setRiskFilter(event.target.value)}>
            <option value="all">All risk</option>
            <option value="low">Low</option>
            <option value="medium">Medium</option>
            <option value="high">High</option>
            <option value="critical">Critical</option>
          </select>
        </label>
        <label>
          Target
          <input value={targetFilter} onChange={(event) => setTargetFilter(event.target.value)} placeholder="service id" />
        </label>
      </div>

      {filtered.length === 0 ? (
        <EmptyState title="No experiments match" detail="Adjust the filters or add a new experiment manifest." />
      ) : (
        <div className="card-grid experiment-grid">
          {filtered.map((experiment) => (
            <article key={experiment.id} className="catalog-card">
              <div className="card-title-row">
                <div>
                  <h2>{experiment.name}</h2>
                  <span>{experiment.target}</span>
                </div>
                <RiskBadge level={experiment.blastRadius} />
              </div>
              <p>{experiment.description}</p>
              <dl className="meta-grid">
                <div>
                  <dt>Type</dt>
                  <dd>{experiment.type}</dd>
                </div>
                <div>
                  <dt>Duration</dt>
                  <dd>{experiment.defaultDuration}</dd>
                </div>
                <div>
                  <dt>Runbook</dt>
                  <dd>{experiment.runbookId}</dd>
                </div>
              </dl>
              <ul>
                {experiment.safetyNotes.map((note) => (
                  <li key={note}>{note}</li>
                ))}
              </ul>
            </article>
          ))}
        </div>
      )}
    </section>
  );
}
