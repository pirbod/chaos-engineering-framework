import type React from 'react';

type PageKey = 'dashboard' | 'experiments' | 'risk' | 'runbooks' | 'integrations' | 'portal';

type Props = {
  activePage: PageKey;
  onNavigate: (page: PageKey) => void;
  children: React.ReactNode;
};

const navItems: Array<{ key: PageKey; label: string }> = [
  { key: 'dashboard', label: 'Dashboard' },
  { key: 'experiments', label: 'Experiments' },
  { key: 'risk', label: 'Service Risk' },
  { key: 'runbooks', label: 'Runbooks' },
  { key: 'integrations', label: 'Integrations' },
  { key: 'portal', label: 'Developer Portal' },
];

export function Shell({ activePage, onNavigate, children }: Props) {
  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <span className="brand-mark">PCH</span>
          <div>
            <strong>Platform Chaos</strong>
            <span>Reliability Hub</span>
          </div>
        </div>
        <nav aria-label="Primary navigation">
          {navItems.map((item) => (
            <button
              key={item.key}
              className={item.key === activePage ? 'active' : ''}
              type="button"
              onClick={() => onNavigate(item.key)}
            >
              {item.label}
            </button>
          ))}
        </nav>
      </aside>
      <main className="content">{children}</main>
    </div>
  );
}

export type { PageKey };
