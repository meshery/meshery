import environmentsResponse from './mockData/environments.json';
import connectionsResponse from './mockData/connections.json';
import patternsResponse from './mockData/patterns.json';

export const mockEnvironmentsApi = async (route) => {
  await route.fulfill({
    status: 200,
    contentType: 'application/json',
    body: JSON.stringify(environmentsResponse),
  });
};

export const mockConnectionsApi = async (route) => {
  await route.fulfill({
    status: 200,
    contentType: 'application/json',
    body: JSON.stringify(connectionsResponse),
  });
};

export const mockPatternsApi = async (route) => {
  await route.fulfill({
    status: 200,
    contentType: 'application/json',
    body: JSON.stringify(patternsResponse),
  });
};
