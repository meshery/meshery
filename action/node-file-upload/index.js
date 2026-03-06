const fs = require("fs");
const path = require("path");
const axios = require("axios");

const themes = {
  DARK: "dark",
  LIGHT: "light"
}

function convertFileToBase64(filePath) {
  try {
    // Read file as binary data
    const fileData = fs.readFileSync(filePath);

    // Convert binary data to Base64
    const base64Data = fileData.toString("base64");

    return base64Data;
  } catch (error) {
    console.error("Error:", error);
    return null;
  }
}

let location = [];
let snapshotlocationToThemeMap = {}
if (process.env.assetLocation != "" && process.env.assetLocation !== undefined) {
  location = process.env.assetLocation.split(",")
}

location?.forEach(val => {
  if (val.endsWith("dark.png")) {
    snapshotlocationToThemeMap[themes.DARK] = val;
  } else {
    snapshotlocationToThemeMap[themes.LIGHT] = val;
  }
})

const dirPath = path.join(
  __dirname,
  "..",
  "cypress-action",
  "cypress",
  "screenshots",
  "loadDesign.js"
);

fs.readdirSync(dirPath).forEach((fileName, index) => {
  setTimeout(() => {
    const filePath = path.join(dirPath, fileName);
    const base64Data = convertFileToBase64(filePath);
    const formData = new FormData();
    formData.append("image", base64Data);
    if (location.length > 0) {
      let theme = themes.LIGHT;
      if (fileName.endsWith("dark.png")) {
        theme = theme.DARK;
      }
      formData.append(
        "assetLocation",
        snapshotlocationToThemeMap[theme]
      );
    }
    const url =
      "https://cloud.layer5.io/api/integrations/github/meta/artifacts";

    const headers = {
      "Content-Type": "multipart/form-data",
      Authorization: `Bearer ${process.env.MESHERY_TOKEN}`,
    };

    if (formData) {
      axios
        .post(url, formData, {
          headers,
        })
        .then((response) => {
          console.log(response.data);
        })
        .catch((e) => {
          console.log(e.response.data);
        });
    } else {
      console.log(null);
    }
  }, 5000 * (index + 1));
});
