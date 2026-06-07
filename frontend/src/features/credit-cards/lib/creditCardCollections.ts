export function findById<TItem extends { id: string }>(items: TItem[] | undefined, id: string) {
  return items?.find((item) => item.id === id)
}
