export type NoRefs<T> = T extends Record<string, unknown>
  ? Omit<T, " $refType" | " $fragmentRefs">
  : T;

export type ExtractRelayEdgeNode<
  T extends { edges: ReadonlyArray<{ node: any | null } | null> | null } | null
> = NoRefs<
  NonNullable<NonNullable<NonNullable<NonNullable<T>["edges"]>[0]>["node"]>
>;
