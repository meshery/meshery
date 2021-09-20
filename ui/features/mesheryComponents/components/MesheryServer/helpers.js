/**
 * getMesheryVersionText returs a well formatted version text
 *
 * If the meshery is running latest version then and is using "edge" channel
 * then it will just show "edge-latest". However, if the meshery is on edge and
 * is running an outdated version then it will return "edge-$version".
 *
 * If on stable channel, then it will always show "stable-$version"
 */
export const getMesheryVersionText = (serverVersion) => {
  // If the version is outdated then no matter what the
  // release channel is, specify the build
  if (serverVersion.outdated) return `${serverVersion.release_channel}-${serverVersion.build}`;

  if (serverVersion.release_channel === "edge") return `${serverVersion.release_channel}-latest`;
  if (serverVersion.release_channel === "stable") return `${serverVersion.release_channel}-${serverVersion.build}`;

  return `${serverVersion.build}`;
};
