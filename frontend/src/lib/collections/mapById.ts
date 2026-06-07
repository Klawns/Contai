export function mapById<TItem extends { id: string }>(items: TItem[] | undefined) {
  return new Map((items ?? []).map((item) => [item.id, item]))
}
