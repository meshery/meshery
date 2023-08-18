import api from "./index"
import { ctxUrl } from "../utils/multi-ctx";
export const PATTERN_ENDPOINT = "/api/pattern"



export async function dryRunPattern(patternFileData, contexts) {
  return deployPatternWithData(patternFileData, contexts, {
    verify : false,
    dryRun : true
  });
}

/**
 * Deploys pattern with the content provided to it.
 * It is {pattern_file} property of the patterns meta
 *
 * @param {String} patternFileData
 * @param {String} contexts
 * @param {{verify: Boolean, dryRun: Boolean}} options
 */
export async function deployPatternWithData(patternFileData, contexts, options) {
  const { verify = false, dryRun = false } = options;
  const endpoint = `${ctxUrl(PATTERN_ENDPOINT+"/deploy",contexts)}${
    verify ? "&verify=true" : ""
  }${dryRun ? "&dryRun=true" : ""}`;
  return await api.post(endpoint, patternFileData) ;
}