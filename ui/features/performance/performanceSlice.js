import { createAsyncThunk, createSlice, nanoid } from "@reduxjs/toolkit";
import { fetchTestProfiles, fetchTestResults } from "./helpers";
/**
 * this feature/module is reponsible for all the actions related to performance such as,
 * performance tests, profiles and running tests
 */

/**
 * @typedef {{id: string, name: string, user_id: string, load_generators: string[], endpoints: string[], service_mesh: import("../serviceMeshes/serviceMeshesSlice").ServiceMeshType, concurrent_request: number, qps: number, duration: string, last_run: string, total_results: number, created_at: string, updated_at: string }} PerformanceProfileType
 */

/**
 * @typedef {{meshery_id: string, name: string, mesh: import("../serviceMeshes/serviceMeshesSlice").ServiceMeshType, test_start_time: string, user_id: string, performance_profile: string, created_at: string, updated_at: string, runner_results: {}}} PerformanceResultType
 */

const initialState = {
  /** @type {PerformanceProfileType[]} */
  profiles: [],
  /** @type {PerformanceResultType[]} */
  results: [],
  loading: false,
  /** @type {import("../provider/providerSlice").Error[]} */
  errors: [],
};

const fetchPerformanceProfilesThunk = createAsyncThunk("performance/fetchPerformanceProfiles", async () => {
  const response = await fetchTestProfiles();
  return response;
});

const fetchPerformanceResultsThunk = createAsyncThunk("performance/fetchPerformanceResults", async () => {
  const response = await fetchTestResults();
  return response;
});

const performanceSlice = createSlice({
  name: "performance",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(fetchPerformanceProfilesThunk.pending, (state) => {
      state.loading = true;
      return state;
    });
    builder.addCase(fetchPerformanceProfilesThunk.fulfilled, (state, action) => {
      state.loading = false;
      state.profiles = action.payload.results;
      return state;
    });
    builder.addCase(fetchPerformanceProfilesThunk.rejected, (state) => {
      state.loading = false;
      state.errors.push({
        code: "",
        id: nanoid(),
        description: "Unable to fetch performance profiles",
        remediation: "try again after sometime",
      });
      return state;
    });

    builder.addCase(fetchPerformanceResultsThunk.pending, (state) => {
      state.loading = true;
      return state;
    });
    builder.addCase(fetchPerformanceResultsThunk.fulfilled, (state, action) => {
      state.loading = false;
      state.results = action.payload.results;
      return state;
    });
    builder.addCase(fetchPerformanceResultsThunk.rejected, (state) => {
      state.loading = false;
      state.errors.push({
        code: "",
        id: nanoid(),
        description: "Unable to fetch performance results",
        remediation: "try again after sometime",
      });
      return state;
    });
  },
});

export default performanceSlice.reducer;
