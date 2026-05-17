import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { fallbackExperiments } from '../src/data/fallback';
import { ExperimentCatalog } from '../src/pages/ExperimentCatalog';

describe('ExperimentCatalog', () => {
  it('filters by target service', () => {
    render(<ExperimentCatalog experiments={fallbackExperiments} />);

    fireEvent.change(screen.getByPlaceholderText('service id'), { target: { value: 'payments' } });

    expect(screen.getByText('Network Latency')).toBeInTheDocument();
    expect(screen.queryByText('Pod Kill')).not.toBeInTheDocument();
  });
});
