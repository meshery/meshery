import { test, expect } from '@playwright/test';
import { ENV } from './env';
import { RelationshipTestFixtures } from './fixtures/relationships';
import _ from 'lodash';

test.describe('Relationship Evaluation', { tag: '@relationship' }, () => {
    // test.use({ storageState: ENV.AUTHFILELOCALPROVIDER });

    for (const fixture of RelationshipTestFixtures) {
        test(`should identify relationships for ${fixture.name}`, async ({ request }) => {
            const response = await request.post(`${ENV.MESHERY_SERVER_URL}/api/meshmodels/relationships/evaluate`, {
                data: {
                    design: fixture.design,
                    options: {
                        returnDiffOnly: false,
                        trace: false
                    }
                }
            });

            expect(response.ok()).toBeTruthy();
            const responseBody = await response.json();

            console.log(responseBody);

            // Check if expected relationships are present in the response
            // The response structure contains a 'design' object which has 'relationships'
            // We need to check if the relationships in the response match the expected ones
            // The matching logic is based on kind, type, and subtype

            const actualRelationships = responseBody.design.relationships || [];

            for (const expectedRel of fixture.expected_relationships) {
                const found = actualRelationships.find((actualRel) => {

                    if (expectedRel.selectors) {
                        if (!_.isEqual(actualRel.selectors, expectedRel.selectors)) {
                            return false;
                        }
                    }

                    return (
                        actualRel.kind === expectedRel.kind &&
                        actualRel.type === expectedRel.type &&
                        actualRel.subType === expectedRel.subType
                    );
                });

                expect(found, `Expected relationship ${JSON.stringify(expectedRel)} not found`).toBeDefined();
            }
        });
    }
});
