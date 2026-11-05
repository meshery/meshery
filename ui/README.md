# <a name="contributing">Contributing Overview</a>

Please do! Thank you for your help in improving Meshery! :balloon:

Find the complete set of contributor guides at https://docs.meshery.io/project/contributing

---

## UI Development Patterns

### Table State Management

All data tables in Meshery UI persist their navigation state (page, page size, sort order, search, and custom filters) in the URL so that links are shareable and state survives navigation. Two hooks handle this:

#### `useTableUrlState` — URL-backed table state

Located at `utils/hooks/useTableUrlState.ts`. Each table gets a short, stable key (e.g. `"con"` for Connections, `"des"` for Designs, `"fil"` for Filters). State is serialised as prefixed query parameters:

| Param           | Description                                         |
| --------------- | --------------------------------------------------- |
| `<key>_page`    | Current page (omitted when 0)                       |
| `<key>_ps`      | Page size (omitted when equal to the default)       |
| `<key>_sort`    | Sort string, e.g. `"name asc"`                      |
| `<key>_q`       | Search query                                        |
| `<key>_<field>` | Custom filter fields declared in `defaults.filters` |

```tsx
const { tableState, updateTableState, copyRowDeepLink } = useTableUrlState({
  tableKey: 'con', // short stable prefix — never change once deployed
  rowParam: 'connectionId', // URL param for row deep-links (default: <key>_row)
  defaults: {
    page: 0,
    pageSize: 10,
    sortOrder: 'created_at desc',
    filters: { status: '', kind: '' },
  },
});

// tableState.page / .pageSize / .sortOrder / .search / .filters are read-only
// updateTableState({ page: 2 }) writes only the changed fields to the URL
// copyRowDeepLink(id) builds a URL pointing at that row and copies it to the clipboard
```

`updateTableState` is referentially stable across renders (uses internal refs for router and defaults), so it is safe to include in child `useCallback` dependency arrays without causing cascade re-renders.

#### `useColumnVisibilityPreference` — localStorage-backed column visibility

Located at `utils/hooks/useColumnVisibilityPreference.ts`. User-initiated column toggles are persisted to `localStorage` under a `meshery_cols_<tableKey>` key. Responsive-layout updates (e.g. hiding columns at narrow widths) are applied on top without overwriting saved preferences.

```tsx
const { columnVisibility, setColumnVisibilityByUser, setColumnVisibilityByResponsive } =
  useColumnVisibilityPreference('connections', responsiveColDefaults);

// In a width-change effect:
useEffect(() => {
  setColumnVisibilityByResponsive(responsiveColDefaults);
  // eslint-disable-next-line react-hooks/exhaustive-deps
}, [width, setColumnVisibilityByResponsive]);
```

Pass `setColumnVisibilityByUser` to the table's column-visibility toolbar so toggling a column is saved. Never call `setColumnVisibilityByResponsive` from user interaction — it does not save.

#### Row Deep-Links in `ResponsiveDataTable`

Use `getCopyDeepLinkAction` from `@sistent/sistent` to add a copy-link action to the row action menu:

```tsx
import { getCopyDeepLinkAction } from '@sistent/sistent';

const actions: TableAction[] = [
  getCopyDeepLinkAction(() => copyRowDeepLink(row.id)),
  // ... other actions
];
```

The helper accepts an optional `title` argument (default: `'Copy link'`) for i18n.
