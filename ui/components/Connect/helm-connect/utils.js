import axios from "axios";

export const postHelmInstall = async (body, installationId) => {
  const res = await axios.post(
    process.env.API_ENDPOINT_PREFIX +
      `/api/integrations/helm/install` +
      `?id=${installationId}`,
    body
  );
  return res.data;
};

export const VerifyHelmRepoLink = async (url) => {
  return await axios.post("/api/integrations/helm/status", url);
};

export const modifiedArrayOfObjects = (source, res) => {
  const newArray = source.map(obj => ({ ...obj }));

  for (const resObj of res) {
    const resKey = Object.keys(resObj)[0];
    const resContent = resObj[resKey];

    const objIndex = newArray.findIndex(obj => Object.keys(obj)[0] === resKey);
    if (objIndex !== -1) {
      newArray[objIndex] = { [resKey]: resContent };
    }
  }

  return newArray;
};