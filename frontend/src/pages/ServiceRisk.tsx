import { RiskBadge } from '../components/RiskBadge';
import type { PlatformData } from '../types';

type Props = {
  data: PlatformData;
};

export function ServiceRisk({ data }: Props) {
  return (
    <section className="page">
      <header className="page-header compact">
        <div>
          <p className="eyebrow">Service risk</p>
          <h1>Prioritized action list for safer experiments and faster detection.</h1>
        </div>
      </header>

      <div className="service-risk-list">
        {data.services.map((service) => {
          const risk = data.risks[service.id];
          return (
            <article className="service-risk-card" key={service.id}>
              <div className="service-risk-header">
                <div>
                  <h2>{service.name}</h2>
                  <span>{service.team} / {service.tier}</span>
                </div>
                {risk ? <RiskBadge level={risk.level} score={risk.score} /> : <span className="muted">No score</span>}
              </div>
              <p>{service.description}</p>
              <dl className="meta-grid four">
                <div>
                  <dt>Owner</dt>
                  <dd>{service.owner || 'Missing'}</dd>
                </div>
                <div>
                  <dt>SLO</dt>
                  <dd>{service.slo}</dd>
                </div>
                <div>
                  <dt>Observability</dt>
                  <dd>{service.observabilityMaturity}%</dd>
                </div>
                <div>
                  <dt>Error rate</dt>
                  <dd>{(service.recentErrorRate * 100).toFixed(2)}%</dd>
                </div>
              </dl>
              {risk ? (
                <div className="factor-layout">
                  <div>
                    <h3>Top factors</h3>
                    {risk.topContributingFactors.map((factor) => (
                      <div className="factor-row" key={factor.name}>
                        <span>{factor.name}</span>
                        <strong>{factor.contribution}</strong>
                      </div>
                    ))}
                  </div>
                  <div>
                    <h3>Recommended actions</h3>
                    <ol>
                      {risk.recommendedActions.map((action) => (
                        <li key={action}>{action}</li>
                      ))}
                    </ol>
                  </div>
                </div>
              ) : null}
            </article>
          );
        })}
      </div>
    </section>
  );
}
