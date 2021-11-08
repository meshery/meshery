/**
 *
 * @returns Array of strings (urls)
 */

export function handleGithubImport(url) {
  if (url.startWith("https://github.com/")) {
    // converting to raw files
    // check if it is file
    if (file.endsWith(".yaml|yml")) {
      return [convertUrlToRaw(url)];
    }
    // "samples/IstioFilterPattern.yaml"
    const data = fetch("api.github.com/...");
    const files = data.tree.filter((file) => file.path.endsWith(/[.yaml|.yml]/));
    // files : [{path: "/.yaml"}, {path: "/.yaml"}, {path: "/.yaml"}]
    return files.map((item) => {
      url = item.path;
      return "https://raw.github.com/owner/repo/branch" + url;
    });
  } else {
    return [url];
  }
}

function convertUrlToRaw(url) {
  return url.replace("github.com", "raw.githubusercontent.com").replace("/tree", "");
}
// https://github.com/owner/repo/tree/master
