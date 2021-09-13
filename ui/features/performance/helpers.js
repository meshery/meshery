import dataFetch from "@/lib/data-fetch";
import fetchPerformanceResults from "./graphql/queries/PerformanceResultQuery";

/**
 *
 * @returns {Promise<{page: number, page_size: number, total_count: number, results:import("./performanceSlice").PerformanceProfileType[]>}} Profiles - Performance Profiles
 */
export const fetchTestProfiles = () =>
  new Promise((res, rej) => {
    dataFetch(
      "/api/user/performance/profiles",
      { credentials: "include" },
      (result) => {
        if (result) {
          res(result);
        }
      },
      (err) => rej(err)
    );
  });

/**
 *
 * @returns {Promise<{page: number, page_size: number, total_count: number, results:import("./performanceSlice").PerformanceResultType[]>}} Results - Performance Results
 */
/**
 *
 * @param {number | string} page - Page number
 * @param {number | string} pageSize - Size of the page
 * @param {string} search - Search keyword
 * @param {string} sortOrder - Order of sorting
 * @param {string} endpoint - Test Profile endpoint
 * @returns
 */
export const fetchTestProfileResults = (page, pageSize, search, sortOrder, endpoint) =>
  new Promise((res, rej) => {
    if (!search) search = "";
    if (!sortOrder) sortOrder = "";

    fetchPerformanceResults({
      selector: {
        pageSize: `${pageSize}`,
        page: `${page}`,
        search: `${encodeURIComponent(search)}`,
        order: `${encodeURIComponent(sortOrder)}`,
      },
      profileID: endpoint.split("/")[endpoint.split("/").length - 2],
    }).subscribe({
      next: (result) => {
        // @ts-ignore
        if (result.fetchResults) {
          res(result.fetchResults);
        }
      },
      error: (err) => rej(err),
    });
  });
