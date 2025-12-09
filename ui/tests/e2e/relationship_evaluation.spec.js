import { test, expect } from '@playwright/test';
import { ENV } from './env';
import { RelationshipTestFixtures } from './fixtures/relationships';
import _ from 'lodash';

test.describe('Relationship Evaluation', { tag: '@relationship' }, () => {

    for (const design of RelationshipTestFixtures) {
        test(`should identify relationships for ${design.name}`, async ({ request }) => {
            const designToTest = { ...design, relationships: [] }
            const response = await request.post(`${ENV.MESHERY_SERVER_URL}/api/meshmodels/relationships/evaluate`, {

                data: {
                    design: designToTest,
                    options: {
                        returnDiffOnly: false,
                        trace: false
                    }
                }
            });

            expect(response.ok()).toBeTruthy();
            const responseBody = await response.json();


            const actualRelationships = responseBody.design.relationships || [];

            // assert number of relationships
            expect(actualRelationships.length).toBeGreaterThanOrEqual(design.relationships.length);

            for (const expectedRel of design.relationships) {
                if (!expectedRel.selectors || expectedRel?.metadata?.isAnnotation) {
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

                expect(found, `Expected relationship ${JSON.stringify(expectedRel)} not found`).toBeDefined();
            }
        });
    }
});
