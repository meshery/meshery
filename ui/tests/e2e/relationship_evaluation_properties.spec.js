import { test, expect } from '@playwright/test';
import { ENV } from './env';
import { RelationshipTestFixtures } from './fixtures/relationships/index';

// --- Helpers ---

function getAtPath(obj, pathSegments) {
  let current = obj;
  for (const key of pathSegments) {
    if (current == null) return undefined;
    current = current[key];
  }
  return current;
}

function findComponent(design, id) {
  return design.components?.find((c) => c.id === id);
}

function fromSelector(rel) {
  return rel.selectors?.[0]?.allow?.from?.[0];
}

function toSelector(rel) {
  return rel.selectors?.[0]?.allow?.to?.[0];
}

// Extract mutation targets from a relationship's patch fields.
function extractPatchMutations(rel) {
  const from = fromSelector(rel);
  const to = toSelector(rel);
  if (!from || !to) return [];

  const fromPatch = from.patch;
  const toPatch = to.patch;

  let mutatorRef, mutatedRef, mutatorId, mutatedId;

  if (fromPatch?.mutatorRef?.[0]) {
    mutatorRef = fromPatch.mutatorRef[0];
    mutatorId = from.id;
  } else if (toPatch?.mutatorRef?.[0]) {
    mutatorRef = toPatch.mutatorRef[0];
    mutatorId = to.id;
  }

  if (toPatch?.mutatedRef?.[0]) {
    mutatedRef = toPatch.mutatedRef[0];
    mutatedId = to.id;
  } else if (fromPatch?.mutatedRef?.[0]) {
    mutatedRef = fromPatch.mutatedRef[0];
    mutatedId = from.id;
  }

  if (!mutatorRef || !mutatedRef || !mutatorId || !mutatedId) return [];
  return [{ mutatorRef, mutatedRef, mutatorId, mutatedId }];
}

// Extract mutation targets from a relationship's match fields (used by binding relationships).
function extractMatchMutations(rel) {
  const results = [];
  const selectorItems = [fromSelector(rel), toSelector(rel)];

  for (const sel of selectorItems) {
    if (!sel?.match) continue;
    const matchFrom = sel.match.from ?? [];
    const matchTo = sel.match.to ?? [];

    // Find the side with mutatorRef and the side with mutatedRef
    for (const mSel of [...matchFrom, ...matchTo]) {
      if (!mSel?.mutatorRef?.length || !mSel.id) continue;
      for (const oSel of [...matchFrom, ...matchTo]) {
        if (oSel === mSel || !oSel?.mutatedRef?.length || !oSel.id) continue;
        const count = Math.min(mSel.mutatorRef.length, oSel.mutatedRef.length);
        for (let i = 0; i < count; i++) {
          // Skip wildcard paths (contain "_")
          if (oSel.mutatedRef[i].includes('_')) continue;
          results.push({
            mutatorRef: mSel.mutatorRef[i],
            mutatedRef: oSel.mutatedRef[i],
            mutatorId: mSel.id,
            mutatedId: oSel.id,
          });
        }
      }
    }
  }
  return results;
}

// Extract all mutation targets from a relationship (patch fields first, then match fields).
function extractAllMutations(rel) {
  const patchMutations = extractPatchMutations(rel);
  if (patchMutations.length > 0) return patchMutations;
  return extractMatchMutations(rel);
}

// Deep clone a JSON-serializable object
function deepClone(obj) {
  return JSON.parse(JSON.stringify(obj));
}

// Delete the value at the given path in an object.
// Returns true if something was deleted.
function deleteAtPath(obj, pathSegments) {
  let current = obj;
  for (let i = 0; i < pathSegments.length - 1; i++) {
    if (current == null) return false;
    current = current[pathSegments[i]];
  }
  if (current == null) return false;
  const lastKey = pathSegments[pathSegments.length - 1];
  if (!(lastKey in current)) return false;
  delete current[lastKey];
  return true;
}

async function evaluateDesign(request, design) {
  const resp = await request.post(
    `${ENV.MESHERY_SERVER_URL}/api/meshmodels/relationships/evaluate`,
    {
      data: {
        design: { ...design, relationships: [] },
        options: { returnDiffOnly: false, trace: false },
      },
    },
  );
  expect(resp.ok()).toBeTruthy();
  return resp.json();
}

// Cache evaluated designs to avoid redundant API calls across tests for the same fixture.
const evaluationCache = new Map();

async function getEvaluatedDesign(request, fixture) {
  const key = fixture.id;
  if (evaluationCache.has(key)) return evaluationCache.get(key);
  const response = await evaluateDesign(request, fixture);
  evaluationCache.set(key, response);
  return response;
}

// --- Property-based Tests ---

test.describe('Relationship Evaluation Properties', { tag: '@relationship' }, () => {
  for (const fixture of RelationshipTestFixtures) {
    test.describe(fixture.name, () => {
      // 1. Structural integrity: every approved relationship's from/to components exist
      test('structural integrity', async ({ request }) => {
        const { design } = await getEvaluatedDesign(request, fixture);

        for (const rel of design.relationships ?? []) {
          if (rel.status !== 'approved') continue;

          const from = fromSelector(rel);
          const to = toSelector(rel);
          if (!from?.id || !to?.id) continue;

          expect(
            findComponent(design, from.id),
            `from component ${from.id} (${from.kind}) missing for ${rel.kind}/${rel.type}/${rel.subType}`,
          ).toBeDefined();
          expect(
            findComponent(design, to.id),
            `to component ${to.id} (${to.kind}) missing for ${rel.kind}/${rel.type}/${rel.subType}`,
          ).toBeDefined();
        }
      });

      // 2. Alias resolution: every resolvedAlias points to a valid value in the parent config
      test('alias resolution', async ({ request }) => {
        const { design } = await getEvaluatedDesign(request, fixture);
        const aliases = design.metadata?.resolvedAliases;
        if (!aliases || Object.keys(aliases).length === 0) return;

        for (const [aliasId, alias] of Object.entries(aliases)) {
          const parent = findComponent(design, alias.resolved_parent_id);
          expect(
            parent,
            `parent ${alias.resolved_parent_id} missing for alias ${aliasId}`,
          ).toBeDefined();

          const val = getAtPath(parent, alias.resolved_ref_field_path);
          const pathStr = alias.resolved_ref_field_path.join('.');
          expect(
            val,
            `alias ${aliasId} path ${pathStr} not found on parent ${parent.component?.kind}`,
          ).toBeDefined();
          expect(
            typeof val === 'object' && val !== null,
            `alias ${aliasId} path ${pathStr} should be an object, got ${typeof val}`,
          ).toBe(true);
        }
      });

      // 3. Config patching correctness: use the engine-evaluated design as baseline
      //    (with correctly discovered relationships), clear mutated values, re-evaluate
      //    WITH relationships, and verify the engine patches values back.
      test('config patching correctness', async ({ request }) => {
        // Step 1: Get baseline design with engine-discovered relationships
        const { design: baseline } = await getEvaluatedDesign(request, fixture);

        // Step 2: Extract mutation targets from the engine's relationships
        const patchTargets = [];
        for (const rel of baseline.relationships ?? []) {
          if (rel.status !== 'approved') continue;
          for (const info of extractAllMutations(rel)) {
            const mutatorComp = findComponent(baseline, info.mutatorId);
            if (!mutatorComp) continue;
            const expectedVal = getAtPath(mutatorComp, info.mutatorRef);
            if (expectedVal === undefined) continue;
            patchTargets.push({ ...info, expectedVal });
          }
        }

        if (patchTargets.length === 0) return;

        // Step 3: Clone baseline, keep relationships, clear mutated values
        const tampered = deepClone(baseline);
        const cleared = [];
        for (const target of patchTargets) {
          const comp = findComponent(tampered, target.mutatedId);
          if (!comp) continue;
          if (target.mutatedRef[0] !== 'configuration') continue;
          if (deleteAtPath(comp, target.mutatedRef)) {
            cleared.push(target);
          }
        }

        expect(cleared.length, 'no mutated values were cleared, nothing to test').toBeGreaterThan(
          0,
        );

        // Step 4: Re-evaluate WITH relationships so the engine patches values back
        const resp = await request.post(
          `${ENV.MESHERY_SERVER_URL}/api/meshmodels/relationships/evaluate`,
          {
            data: {
              design: tampered,
              options: { returnDiffOnly: false, trace: false },
            },
          },
        );
        expect(resp.ok()).toBeTruthy();
        const { design } = await resp.json();

        // Step 5: Verify engine patched the values back
        for (const target of cleared) {
          const mutatedComp = findComponent(design, target.mutatedId);
          if (!mutatedComp) continue;

          const patchedVal = getAtPath(mutatedComp, target.mutatedRef);
          expect(
            patchedVal,
            `engine failed to patch ${target.mutatedRef.join('.')} ` +
              `(expected ${JSON.stringify(target.expectedVal)})`,
          ).toEqual(target.expectedVal);
        }
      });

      // 4. Evaluation idempotency: re-evaluating produces no new meaningful actions
      test('evaluation idempotency', async ({ request }) => {
        const { design } = await getEvaluatedDesign(request, fixture);

        // Pass design WITH relationships intact for re-evaluation
        const resp = await request.post(
          `${ENV.MESHERY_SERVER_URL}/api/meshmodels/relationships/evaluate`,
          {
            data: {
              design,
              options: { returnDiffOnly: false, trace: false },
            },
          },
        );
        expect(resp.ok()).toBeTruthy();
        const secondResponse = await resp.json();

        const meaningfulActions = (secondResponse.actions ?? []).filter(
          (a) => a.op !== 'update_relationship',
        );
        expect(
          meaningfulActions,
          `re-evaluation produced ${meaningfulActions.length} unexpected actions`,
        ).toHaveLength(0);
      });
    });
  }
});
