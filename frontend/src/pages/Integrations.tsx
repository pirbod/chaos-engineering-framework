import type { IntegrationStatus } from '../types';

type Props = {
  integrations: IntegrationStatus[];
};

export function Integrations({ integrations }: Props) {
  return (
    <section className="page">
      <header className="page-header compact">
        <div>
          <p className="eyebrow">Integrations</p>
          <h1>Optional platform signals without blocking the local demo.</h1>
        </div>
      </header>

      <div className="card-grid integration-grid">
        {integrations.map((integration) => (
          <article className="integration-card" key={integration.id}>
            <div className="card-title-row">
              <div>
                <h2>{integration.name}</h2>
                <span>{integration.category}</span>
              </div>
              <span className={`status-pill status-${integration.status}`}>{integration.status}</span>
            </div>
            <p>{integration.details}</p>
            <div className="tag-row">
              {integration.signals.map((signal) => (
                <span key={signal}>{signal}</span>
              ))}
            </div>
            <a href={integration.docsUrl}>{integration.docsUrl}</a>
          </article>
        ))}
      </div>
    </section>
  );
}
