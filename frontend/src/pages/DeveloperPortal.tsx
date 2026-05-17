const workflows = [
  {
    title: 'Request a chaos experiment',
    items: ['Pick an approved experiment from the catalog.', 'Confirm owner, SLO and rollback threshold.', 'Open a pull request with the manifest and runbook link.'],
  },
  {
    title: 'Onboard a service',
    items: ['Add Backstage catalog metadata.', 'Attach owner, dashboard, SLO and runbook links.', 'Start with pod-kill or latency in canary scope.'],
  },
  {
    title: 'Read the risk score',
    items: ['Scores combine blast radius, criticality, error rate and maturity.', 'Critical scores need platform approval.', 'Top factors explain the fastest risk reduction path.'],
  },
  {
    title: 'Reduce MTTD',
    items: ['Tag chaos events by service and risk.', 'Use SLO burn alerts and runbook links.', 'Capture AI summaries for incident handoff.'],
  },
];

export function DeveloperPortal() {
  return (
    <section className="page">
      <header className="page-header compact">
        <div>
          <p className="eyebrow">Developer portal</p>
          <h1>Self-service workflows for product teams using the platform.</h1>
        </div>
      </header>

      <div className="workflow-grid">
        {workflows.map((workflow) => (
          <article className="workflow-card" key={workflow.title}>
            <h2>{workflow.title}</h2>
            <ol>
              {workflow.items.map((item) => (
                <li key={item}>{item}</li>
              ))}
            </ol>
          </article>
        ))}
      </div>
    </section>
  );
}
