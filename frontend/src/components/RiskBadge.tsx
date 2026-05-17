import type { RiskResult } from '../types';

type Props = {
  level: RiskResult['level'] | string;
  score?: number;
};

export function RiskBadge({ level, score }: Props) {
  return (
    <span className={`risk-badge risk-${level}`}>
      <span>{level}</span>
      {typeof score === 'number' ? <strong>{score}</strong> : null}
    </span>
  );
}
