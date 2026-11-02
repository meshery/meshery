// Immutable read/write helpers for the sparse configuration documents edited by
// ControllersConfigForm. "Absent" is meaningful in those documents - a field
// that is not present inherits from the next precedence layer - so writing
// `undefined` deletes the leaf and prunes any parent left empty rather than
// leaving an empty object behind.

export type FieldPath = (string | number)[];

export const getPath = (doc: unknown, path: FieldPath): unknown => {
  let node: unknown = doc;
  for (const key of path) {
    if (node == null || typeof node !== 'object') return undefined;
    node = (node as Record<string | number, unknown>)[key];
  }
  return node;
};

// deleteAtPath removes the leaf at path and prunes any parents left empty,
// so an all-inherit section disappears from the document entirely.
const deleteAtPath = (node: unknown, path: FieldPath): void => {
  if (node == null || typeof node !== 'object' || path.length === 0) return;
  const obj = node as Record<string | number, unknown>;
  const [head, ...rest] = path;
  if (rest.length === 0) {
    delete obj[head];
    return;
  }
  deleteAtPath(obj[head], rest);
  const child = obj[head];
  if (child && typeof child === 'object' && Object.keys(child as object).length === 0) {
    delete obj[head];
  }
};

/** Returns a deep copy of `doc` with `path` set to `value`, or deleted when `value` is undefined. */
export const setPath = <T>(doc: T, path: FieldPath, value: unknown): T => {
  // An empty path names no leaf to write. Without this guard
  // `path[path.length - 1]` is `undefined` and the assignment lands on an
  // "undefined" key, silently corrupting the draft's shape instead of failing.
  // `deleteAtPath` already refuses an empty path, so the delete branch is safe.
  if (path.length === 0) return doc;
  const next: T = JSON.parse(JSON.stringify(doc ?? {}));
  if (value === undefined) {
    deleteAtPath(next, path);
    return next;
  }
  let node: Record<string | number, unknown> = next as Record<string | number, unknown>;
  for (const key of path.slice(0, -1)) {
    if (node[key] == null || typeof node[key] !== 'object') {
      node[key] = {};
    }
    node = node[key] as Record<string | number, unknown>;
  }
  node[path[path.length - 1]] = value;
  return next;
};
