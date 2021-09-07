import { createAsyncThunk, createSlice, nanoid } from "@reduxjs/toolkit";
import { fetchProviderDetails } from "./helpers";

/**
 * This feature/module is responsible for all the data and operations
 * regarding provider component
 */

/**
 * @typedef {{title: string, href: {uri: string,external: boolean}, component: string, icon : string, link: boolean, show: boolean }} Navigator
 */

/**
 * @typedef {{component: string}} UserPref
 */

/**
 * @typedef {{component: string, path: string}} Graphql
 */

/**
 * @typedef {{feature: string, endpoint :string}} Capability
 */

/**
 * @typedef {{code: string, description: string, remediation: string, id:string }} Error
 */

/**
 * @typedef {{loading: boolean, errors: Array.<Error>, type: string, version: string, name: string, description: string, url: string, extensions:{navigator: Array.<Navigator>, userPrefs: Array.<UserPref>, graphql: Array.<Graphql> }, capabilities: Array.<Capability> }} Provider
 */

/** @type {Provider} */
const initialState = {
  type: "",
  version: "",
  description: "",
  url: "",
  extensions: {},
  loading: false,
  errors: null,
};

export const fetchProviderDetailsThunk = createAsyncThunk("provider/fetchProviderDetails", async () => {
  const response = await fetchProviderDetails();
  return response;
});

const providerSlice = createSlice({
  name: "provider",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(fetchProviderDetailsThunk.pending, (state) => {
      state.loading = true;
      return state;
    }),
      builder.addCase(fetchProviderDetailsThunk.fulfilled, (state, action) => {
        state.type = action.payload?.provider_type;
        state.name = action.payload?.provider_name;
        state.description = action.payload?.description;
        state.url = action.payload?.provider_url;
        (state.extensions = {
          navigator: action.payload?.extensions?.navigator,
          userPrefs: action.payload?.extensions?.user_prefs,
          graphql: action.payload?.extensions?.graphql,
        }),
          (state.version = action.payload?.package_version),
          (state.capabilities = action.payload?.capabilities);
        return state;
      }),
      builder.addCase(fetchProviderDetailsThunk.rejected, (state, action) => {
        state.loading = false;
        if (action.error) {
          state.errors.push({
            id: nanoid(),
            code: "NtwrkErr",
            description: "Internal Server Error",
            remediation: "Please try again after sometime",
          });
        }
        return;
      });
  },
});

export default providerSlice.reducer;

export const loadingSelector = (state) => state.provider.loading;
export const providerCapabilitiesSelector = (state) => state.provider.capabilities;
