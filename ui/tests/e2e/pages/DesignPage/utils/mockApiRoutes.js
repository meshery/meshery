import environmentsResponse from './mockData/environments.json';
import connectionsResponse from './mockData/connections.json';

export const mockEnvironmentsApi = async (route) => {
  await route.fulfill({
    status: 200,
    body: JSON.stringify(environmentsResponse),
  });
};

export const mockConnectionsApi = async (route) => {
  await route.fulfill({
    status: 200,
    body: JSON.stringify(connectionsResponse),
  });
};
