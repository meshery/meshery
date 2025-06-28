import { v4 as uuidv4 } from 'uuid';
import perfMockData from './mockPerfData.json';

const TEST_PROFILE_ID = uuidv4();
const USER_ID = uuidv4();

function getMockResult(profileID, url) {
  const result = structuredClone(perfMockData.perfResults);
  result.performance_profile = profileID;
  result.runner_results.URL = url;
  return result;
}

export async function mockPerfApis(page, requestData) {
  const timestamp = new Date().toISOString();

  await page.route('**/api/system/graphql/query', async (route, request) => {
    const postData = await request.postDataJSON();
    const query = postData?.query || '';

    if (query.includes('getPerformanceProfiles')) {
      return route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: {
            getPerformanceProfiles: {
              page: 0,
              page_size: 10,
              total_count: 1,
              profiles: [
                {
                  id: TEST_PROFILE_ID,
                  name: requestData.profileName,
                  load_generators: [requestData.loadGenerator],
                  endpoints: [requestData.url],
                  service_mesh: requestData.serviceMesh,
                  concurrent_request: Number(requestData.concurrentRequest),
                  qps: Number(requestData.qps),
                  duration: requestData.duration,
                  created_at: timestamp,
                  updated_at: timestamp,
                  last_run: timestamp,
                  total_results: 1,
                  user_id: USER_ID,
                  metadata: { additional_options: [''], ca_certificate: {} },
                },
              ],
            },
          },
        }),
      });
    }

    if (query.includes('fetchResults')) {
      return route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: {
            fetchResults: {
              page: 0,
              page_size: 10,
              total_count: 1,
              results: [getMockResult(postData.variables.profileID, requestData.url)],
            },
          },
        }),
      });
    }

    return route.continue();
  });

  await page.route('**/api/**/performance/profiles', async (route, request) => {
    const postData = await request.postDataJSON();
    return route.fulfill({
      status: 201,
      contentType: 'application/json',
      body: JSON.stringify({
        ...postData,
        id: TEST_PROFILE_ID,
        created_at: timestamp,
        updated_at: timestamp,
        user_id: USER_ID,
      }),
    });
  });
}
