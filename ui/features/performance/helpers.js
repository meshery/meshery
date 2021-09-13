import dataFetch from "@/lib/data-fetch";

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
export const fetchTestResults = () =>
  new Promise((res, rej) => {
    dataFetch(
      "/api/user/performance/profiles/results",
      { credentials: "include" },
      (result) => {
        if (result) {
          res(result);
        }
      },
      (err) => rej(err)
    );
  });
