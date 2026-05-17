type Props = {
  title: string;
  detail: string;
};

export function EmptyState({ title, detail }: Props) {
  return (
    <div className="empty-state" role="status">
      <strong>{title}</strong>
      <p>{detail}</p>
    </div>
  );
}
