import { test, expect } from '@playwright/test';
import { ENV } from './env';
import _ from 'lodash';

const DESIGNS_TO_TEST = [
  {
    id: '13e803b7-596c-4620-bdc4-4d3a28a027a2',
    name: 'Container-Hierarchical-Parent-Alias-Relationship',
  },
  {
    id: '83191f0e-001a-430e-8620-158a05323ac5',
    name: 'deployment-configmap-reference-relationship',
  },
  {
    id: '28e09dc6-de75-4c61-8531-e6e3744f7e98',
    name: 'Hierarchical-Parent-Namespace-Relationship',
  },
  {
    id: '2aa761d3-cd30-419a-a441-2f042934116f',
    name: 'Service-To-Deployment-Network',
  },
  {
    name: 'pv-pvc-edge-non-binding-reference-relationship',
    id: '9f249a4b-13da-4ea1-860e-98c4060a444a',
  },
  {
    id: '34fe3846-4b90-4558-914c-8e57caefd52f',
    name: 'meshery-design',
  },
  {
    id: 'a674263f-2e62-49e5-a986-2585b13c6591',
    name: 'All Relationships',
  },
];

test.describe('Relationship Evaluation', { tag: '@relationship' }, () => {
  for (const { id, name } of DESIGNS_TO_TEST) {
    test(`should identify relationships for ${name}`, async ({ request }, testInfo) => {
      const designResponse = await request.get(
        `${ENV.REMOTE_PROVIDER_URL}/api/content/patterns/${id}`,
      );
      const responseJson = await designResponse.json();
      const design = JSON.parse(responseJson.pattern_file);

      const designToTest = { ...design, relationships: [] };

      const response = await request.post(
        `${ENV.MESHERY_SERVER_URL}/api/meshmodels/relationships/evaluate`,
        {
          data: {
            design: designToTest,
            options: {
              returnDiffOnly: false,
              trace: false,
            },
          },
        },
      );

      expect(response.ok()).toBeTruthy();
      const responseBody = await response.json();

      const actualRelationships = responseBody.design.relationships || [];

      var failures = 0;

      for (const expectedRel of design.relationships) {
        if (
          expectedRel.status !== 'approved' ||
          !expectedRel.selectors ||
          expectedRel?.metadata?.isAnnotation ||
          expectedRel.subType === 'annotation'
        ) {
          continue;
        }
        const found = actualRelationships.find((actualRel) => {
          const expectedSelector = expectedRel.selectors[0];
          const actualSelector = actualRel.selectors[0];

          return (
            actualRel.kind === expectedRel.kind &&
            actualRel.type === expectedRel.type &&
            actualRel.subType === expectedRel.subType &&
            _.isEqual(actualSelector.from, expectedSelector.from) &&
            _.isEqual(actualSelector.to, expectedSelector.to)
          );
        });

        // Add annotation for this relationship test result
        const selector = expectedRel.selectors[0];
        const fromKind = selector?.allow?.from?.[0]?.kind || '-';
        const toKind = selector?.allow?.to?.[0]?.kind || '-';
        const modelName = expectedRel.model?.name || '-';

        const relationshipData = {
          kind: expectedRel.kind,
          type: expectedRel.type,
          subType: expectedRel.subType,
          from: fromKind,
          to: toKind,
          model: modelName,
          designName: design.name,
          status: found ? 'pass' : 'fail',
        };
        testInfo.annotations.push({
          type: 'relationship',
          description: JSON.stringify(relationshipData),
        });

        if (!found) {
          failures++;
        }

        // expect(found, `Expected relationship ${JSON.stringify(expectedRel)} not found`).toBeDefined();
      }

      // assert number of relationships
      expect(failures).toBe(0);
    });
  }
});
