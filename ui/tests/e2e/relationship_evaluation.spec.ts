import { test, expect } from '@playwright/test';
import { ENV } from './env';
import _ from 'lodash';

interface Design {
  id: string;
  name: string;
}

const DESIGNS_TO_TEST: Design[] = [
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

interface SelectorItem {
  kind?: string;
}

interface SelectorAllow {
  from?: SelectorItem[];
  to?: SelectorItem[];
}

interface Selector {
  from?: Record<string, unknown>;
  to?: Record<string, unknown>;
  allow?: SelectorAllow;
}

interface RelationshipMetadata {
  isAnnotation?: boolean;
}

interface RelationshipModel {
  name?: string;
}

interface Relationship {
  kind: string;
  type: string;
  subType: string;
  status?: string;
  selectors?: Selector[];
  metadata?: RelationshipMetadata;
  model?: RelationshipModel;
}

interface DesignFile {
  name: string;
  relationships: Relationship[];
}

interface DesignFetchResponse {
  patternFile: string;
}

interface EvaluateResponseBody {
  design: {
    relationships: Relationship[];
  };
}

interface RelationshipAnnotation {
  kind: string;
  type: string;
  subType: string;
  from: string;
  to: string;
  model: string;
  designName: string;
  status: 'pass' | 'fail';
}

test.describe('Relationship Evaluation', { tag: '@relationship' }, () => {
  for (const { id, name } of DESIGNS_TO_TEST) {
    test(`should identify relationships for ${name}`, async ({ request }, testInfo) => {
      const designResponse = await request.get(
        `${ENV.REMOTE_PROVIDER_URL}/api/content/patterns/${id}`,
      );

      expect(
        designResponse.ok(),
        `Failed to fetch design ${name}. Status ${designResponse.status()}`,
      ).toBeTruthy();

      const responseJson = (await designResponse.json()) as DesignFetchResponse;
      const design = JSON.parse(responseJson.patternFile) as DesignFile;

      const designToTest: DesignFile = { ...design, relationships: [] };

      const response = await request.post(
        `${ENV.MESHERY_SERVER_URL}/api/registry/relationships/evaluate`,
        {
          data: {
            design: designToTest,
            options: {
              returnDiffOnly: false,
              enableTrace: false,
            },
          },
        },
      );

      expect(response.ok()).toBeTruthy();
      const responseBody = (await response.json()) as EvaluateResponseBody;

      const actualRelationships: Relationship[] = responseBody.design.relationships || [];

      let failures = 0;

      for (const expectedRel of design.relationships) {
        if (
          expectedRel.status !== 'approved' ||
          !expectedRel.selectors ||
          expectedRel?.metadata?.isAnnotation ||
          expectedRel.subType === 'annotation'
        ) {
          continue;
        }

        const found = actualRelationships.find((actualRel: Relationship) => {
          const expectedSelector = expectedRel.selectors?.[0];
          const actualSelector = actualRel.selectors?.[0];

          return (
            !!expectedSelector &&
            !!actualSelector &&
            actualRel.kind === expectedRel.kind &&
            actualRel.type === expectedRel.type &&
            actualRel.subType === expectedRel.subType &&
            _.isEqual(actualSelector.from, expectedSelector.from) &&
            _.isEqual(actualSelector.to, expectedSelector.to)
          );
        });

        const selector = expectedRel.selectors[0];
        const fromKind = selector?.allow?.from?.[0]?.kind || '-';
        const toKind = selector?.allow?.to?.[0]?.kind || '-';
        const modelName = expectedRel.model?.name || '-';

        const relationshipData: RelationshipAnnotation = {
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
      }

      console.log(`Total relationship failures for ${design.name}: ${failures}`);
    });
  }
});
