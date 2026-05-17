import { useMemo, useState } from 'react';
import { EmptyState } from '../components/EmptyState';
import type { Runbook } from '../types';

type Props = {
  runbooks: Runbook[];
};

export function Runbooks({ runbooks }: Props) {
  const [query, setQuery] = useState('');
  const filtered = useMemo(
    () =>
      runbooks.filter((runbook) => {
        const haystack = `${runbook.title} ${runbook.summary} ${runbook.appliesTo.join(' ')}`.toLowerCase();
        return haystack.includes(query.toLowerCase());
      }),
    [query, runbooks],
  );

  return (
    <section className="page">
      <header className="page-header compact">
        <div>
          <p className="eyebrow">Runbooks</p>
          <h1>Actionable remediation paths tied to services and experiments.</h1>
        </div>
      </header>

      <div className="toolbar single">
        <label>
          Search runbooks
          <input value={query} onChange={(event) => setQuery(event.target.value)} placeholder="vault, runner, kubernetes" />
        </label>
      </div>

      {filtered.length === 0 ? (
        <EmptyState title="No runbooks found" detail="Try a service, integration or failure mode keyword." />
      ) : (
        <div className="runbook-list">
          {filtered.map((runbook) => (
            <article className="runbook-card" key={runbook.id}>
              <div className="card-title-row">
                <div>
                  <h2>{runbook.title}</h2>
                  <span>{runbook.appliesTo.join(', ')}</span>
                </div>
                <span className={`status-pill severity-${runbook.severity}`}>{runbook.severity}</span>
              </div>
              <p>{runbook.summary}</p>
              <div className="runbook-columns">
                <div>
                  <h3>Steps</h3>
                  <ol>
                    {runbook.steps.map((step) => (
                      <li key={step}>{step}</li>
                    ))}
                  </ol>
                </div>
                <div>
                  <h3>Checks</h3>
                  <ul>
                    {runbook.checks.map((check) => (
                      <li key={check}>{check}</li>
                    ))}
                  </ul>
                </div>
                <div>
                  <h3>Rollback</h3>
                  <ul>
                    {runbook.rollback.map((step) => (
                      <li key={step}>{step}</li>
                    ))}
                  </ul>
                </div>
              </div>
            </article>
          ))}
        </div>
      )}
    </section>
  );
}
