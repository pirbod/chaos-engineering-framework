import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { RiskBadge } from '../src/components/RiskBadge';

describe('RiskBadge', () => {
  it('renders level and score', () => {
    render(<RiskBadge level="critical" score={91} />);

    expect(screen.getByText('critical')).toBeInTheDocument();
    expect(screen.getByText('91')).toBeInTheDocument();
  });
});
