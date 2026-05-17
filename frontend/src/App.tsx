import { useEffect, useState } from 'react';
import { loadPlatformData } from './api/client';
import { Shell, type PageKey } from './components/Shell';
import { fallbackExperiments, fallbackIntegrations, fallbackRisks, fallbackRunbooks, fallbackServices } from './data/fallback';
import { Dashboard } from './pages/Dashboard';
import { DeveloperPortal } from './pages/DeveloperPortal';
import { ExperimentCatalog } from './pages/ExperimentCatalog';
import { Integrations } from './pages/Integrations';
import { Runbooks } from './pages/Runbooks';
import { ServiceRisk } from './pages/ServiceRisk';
import type { PlatformData } from './types';

const initialData: PlatformData = {
  experiments: fallbackExperiments,
  services: fallbackServices,
  runbooks: fallbackRunbooks,
  integrations: fallbackIntegrations,
  risks: fallbackRisks,
  usingFallback: true,
};

export default function App() {
  const [activePage, setActivePage] = useState<PageKey>('dashboard');
  const [data, setData] = useState<PlatformData>(initialData);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let mounted = true;
    loadPlatformData()
      .then((loaded) => {
        if (mounted) {
          setData(loaded);
        }
      })
      .finally(() => {
        if (mounted) {
          setLoading(false);
        }
      });
    return () => {
      mounted = false;
    };
  }, []);

  return (
    <Shell activePage={activePage} onNavigate={setActivePage}>
      <div className="topbar">
        <div>
          <strong>Platform partner mode</strong>
          <span>Risk-first chaos workflows for product teams</span>
        </div>
        <span className="environment-pill">{loading ? 'Loading API' : data.usingFallback ? 'Fallback demo data' : 'Live API'}</span>
      </div>

      {data.usingFallback && !loading ? (
        <div className="notice" role="status">
          Backend unavailable or partially unavailable. Showing local demo data so developers can still explore the platform.
        </div>
      ) : null}

      {activePage === 'dashboard' ? <Dashboard data={data} /> : null}
      {activePage === 'experiments' ? <ExperimentCatalog experiments={data.experiments} /> : null}
      {activePage === 'risk' ? <ServiceRisk data={data} /> : null}
      {activePage === 'runbooks' ? <Runbooks runbooks={data.runbooks} /> : null}
      {activePage === 'integrations' ? <Integrations integrations={data.integrations} /> : null}
      {activePage === 'portal' ? <DeveloperPortal /> : null}
    </Shell>
  );
}
