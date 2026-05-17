import { RiskBadge } from '../components/RiskBadge';
import { StatCard } from '../components/StatCard';
import type { PlatformData } from '../types';

type Props = {
  data: PlatformData;
};

export function Dashboard({ data }: Props) {
  const criticalServices = Object.values(data.risks).filter((risk) => risk.level === 'critical').length;
  const activeExperiments = data.experiments.filter((experiment) => experiment.lastOutcome !== 'passed').length;
  const recent = data.experiments.slice(0, 5);

  return (
    <section className="page">
      <header className="page-header">
        <div>
          <p className="eyebrow">Executive dashboard</p>
          <h1>Reliability signal, experiment safety and ownership in one view.</h1>
        </div>
        <div className="header-action">MTTD target: &lt; 15m</div>
      </header>

      <div className="stat-grid">
        <StatCard label="Total services" value={data.services.length} detail="tracked in catalog" />
        <StatCard label="Active findings" value={activeExperiments} detail="experiments needing review" />
        <StatCard label="Critical risk services" value={criticalServices} detail="need platform approval" />
        <StatCard label="MTTD trend" value="18m" detail="demo placeholder, -22% target" />
      </div>

      <div className="two-column">
        <section className="panel">
          <h2>Recent experiment outcomes</h2>
          <div className="table">
            {recent.map((experiment) => (
              <div className="table-row" key={experiment.id}>
                <div>
                  <strong>{experiment.name}</strong>
                  <span>{experiment.target}</span>
                </div>
                <RiskBadge level={experiment.blastRadius} />
                <span className={`status-pill status-${experiment.lastOutcome}`}>{experiment.lastOutcome}</span>
              </div>
            ))}
          </div>
        </section>

        <section className="panel">
          <h2>Highest service risk</h2>
          <div className="risk-list">
            {data.services
              .map((service) => ({ service, risk: data.risks[service.id] }))
              .filter(({ risk }) => Boolean(risk))
              .sort((a, b) => b.risk.score - a.risk.score)
              .slice(0, 4)
              .map(({ service, risk }) => (
                <article key={service.id} className="risk-item">
                  <div>
                    <strong>{service.name}</strong>
                    <span>{service.team || 'Owner metadata missing'}</span>
                  </div>
                  <RiskBadge level={risk.level} score={risk.score} />
                </article>
              ))}
          </div>
        </section>
      </div>
    </section>
  );
}
