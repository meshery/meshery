import { api } from "./index"

/**
 * Convert an object with filters into a parsed object.
 *
 * @param {Object} filters - The input object containing filters.
 * @returns {Object} - The parsed object with filters.
 */
function parseFilters(filters) {
  return Object.entries(filters).reduce((parsedFilters, [key, value]) => {
    if (value || typeof value === 'string') {
      parsedFilters[key] =
        typeof value === 'string'
          ? value
          : JSON.stringify(value, (_key, val) =>
            val instanceof Set ? [...val] : val
          );
    }
    return parsedFilters;
  }, {});
}
export const PROVIDER_TAGS = {
  EVENT: "event"
}
export const notificationCenterApi = api
  .enhanceEndpoints({
    addTagTypes: Object.values(PROVIDER_TAGS),
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getEvents: builder.query({
        query: ({
          page = 1,
          filters = {}
        }) => {

          const parsedFilters = parseFilters(filters);
          // console.log("parsedFilters", parsedFilters)
          return {
            url: `v2/events`,
            params: {
              ...parsedFilters,
              page: page,
              sort: "created_at",
              order: "desc",
              page_size: 15
            }
          }
        },
        providesTags: [PROVIDER_TAGS.EVENT],
        // keepUnusedDataFor : "0.001"
      }),
      getEventsSummary: builder.query({
        query: () => {
          return {
            url: `v2/events?page=$1&page_size=1`
          }
        },
        transformResponse: (response) => {
          return {
            count_by_severity_level: response.count_by_severity_level,
            total_count: response.total_count
          }
        },
        providesTags: [PROVIDER_TAGS.EVENT],
      }),

      updateStatus: builder.mutation({
        query: ({ id, status }) => ({
          url: `events/status/${id}`,
          method: 'POST',
          body: {
            status: status
          }
        }),
        invalidatesTags: [PROVIDER_TAGS.EVENT],
      }),

      deleteEvent: builder.mutation({
        query: ({ id }) => ({
          url: `events/${id}`,
          method: 'DELETE',
        }),
        invalidatesTags: [PROVIDER_TAGS.EVENT],
      })
    }),
    overrideExisting: false,
  })

export const {
  useGetEventsSummaryQuery,
  useUpdateStatusMutation,
  useDeleteEventMutation,
  useLazyGetEventsQuery,
} = notificationCenterApi
