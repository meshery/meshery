import { ENV } from '../env';

const getTokenFromAuthFile = () => {
  try {
    const authFile = require('../../../playwright/.auth/user.json');
    return authFile.cookies.find((cookie) => cookie.name === 'token').value;
  } catch (error) {
    console.error('Error while getting token from auth file:', error);
    return null;
  }
};

export async function fetchData(endpoint, method = 'GET', body = null) {
  const url = `${ENV.MESHERY_SERVER_URL}/api${endpoint}`;
  const yourToken = ENV.PROVIDER_TOKEN;
  const token = yourToken || getTokenFromAuthFile();
  const urlWithoutProtocol = new URL(ENV.REMOTE_PROVIDER_URL).host;

  const headers = {
    'meshery-token': token,
    Cookie: `meshery-provider=Meshery; ${urlWithoutProtocol}_ref=/;token=${token}`,
  };

  try {
    const response = await fetch(url, { headers, method, body });
    if (!response.ok) {
      throw new Error(`Error fetching data: ${response.status}`);
    }
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error:', error);
    return null; // Or handle the error differently
  }
}
