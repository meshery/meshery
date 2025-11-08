# Component Exclusion Filtering with Regex Support

## Overview

The Meshery REST API now supports exclusion filtering based on regular expressions for component queries. This allows clients to exclude components that match specific patterns, enabling more flexible and precise filtering scenarios.

## Dependencies

This feature requires changes in both:
- **meshkit** (v0.8.47+): Adds `ExcludeName` and `ExcludeDisplayName` fields to `ComponentFilter`
- **meshery**: Updates handlers to pass exclusion parameters from query strings

> **Note for Developers**: During local development, ensure the meshkit replace directive in go.mod points to your local meshkit repository with the required changes. For production builds, this should reference the published meshkit version.

## New Query Parameters

The following endpoints now support two new query parameters:

- `excludeName`: Exclude components whose kind matches the provided regex pattern
- `excludeDisplayName`: Exclude components whose display name matches the provided regex pattern

### Supported Endpoints

1. `/api/meshmodels/models/{model}/components`
2. `/api/meshmodels/components`
3. `/api/meshmodels/categories/{category}/components`
4. `/api/meshmodels/categories/{category}/models/{model}/components`

## Usage Examples

### Example 1: Exclude all components ending with "List"

```bash
# Exclude all Kubernetes components whose kind ends with "List"
curl "http://localhost:9081/api/meshmodels/models/kubernetes/components?excludeName=.*List$"
```

This will return all Kubernetes components except those like:
- PodList
- ServiceList
- DeploymentList
- etc.

### Example 2: Exclude components by display name pattern

```bash
# Exclude all components whose display name contains "Internal"
curl "http://localhost:9081/api/meshmodels/components?excludeDisplayName=.*Internal.*"
```

### Example 3: Combine with existing filters

```bash
# Get all Kubernetes components in v1, exclude those ending with "List"
curl "http://localhost:9081/api/meshmodels/models/kubernetes/components?apiVersion=v1&excludeName=.*List$"
```

### Example 4: Multiple exclusions

```bash
# Exclude multiple patterns using alternation
curl "http://localhost:9081/api/meshmodels/models/kubernetes/components?excludeName=(.*List$|.*Status$)"
```

This will exclude components ending with either "List" or "Status".

## Regex Pattern Syntax

The exclusion filters use PostgreSQL's regex pattern syntax (POSIX regular expressions):

- `.` - matches any character
- `*` - matches 0 or more of the preceding element
- `+` - matches 1 or more of the preceding element
- `^` - matches the start of the string
- `$` - matches the end of the string
- `|` - alternation (OR)
- `[]` - character class
- `()` - grouping

### Common Patterns

| Pattern | Description | Example Matches |
|---------|-------------|-----------------|
| `.*List$` | Ends with "List" | PodList, ServiceList |
| `^Internal.*` | Starts with "Internal" | InternalService, InternalConfig |
| `.*(Config\|Secret)$` | Ends with "Config" or "Secret" | AppConfig, DBSecret |
| `.*-v[0-9]+$` | Ends with version number | component-v1, service-v2 |

## Behavior and Precedence

1. **Exclusion filters are applied AFTER inclusion filters**: If both `search` and `excludeName` are provided, components are first filtered by the search term, then exclusions are applied.

2. **Case-sensitive matching**: The regex pattern matching is case-sensitive by default.

3. **Multiple exclusion conditions**: When both `excludeName` and `excludeDisplayName` are provided, a component is excluded if it matches ANY of the patterns (OR logic). The exclusion filters work independently - if a component's kind matches `excludeName` OR its display name matches `excludeDisplayName`, it will be excluded from results.

4. **Backward compatibility**: Existing queries without exclusion parameters continue to work as before.

## API Response

The API response format remains unchanged. The `count` field reflects the number of components after applying all filters including exclusions:

```json
{
  "page": 1,
  "page_size": 25,
  "count": 45,
  "components": [
    // ... array of component definitions
  ]
}
```

## Implementation Details

- The exclusion filtering is implemented at the database query level using PostgreSQL's `~` regex operator
- Filters are applied in the ComponentFilter.Get() method in meshkit
- The implementation uses parameterized queries to prevent SQL injection

## Error Handling

- Invalid regex patterns will result in a database error returned to the client
- Empty exclusion parameters are ignored (no filtering applied)

## Performance Considerations

- Regex pattern matching may have performance implications for large datasets
- Consider using indexes on component kind and display_name fields for better performance
- Use specific patterns rather than broad ones when possible

## Future Enhancements

Potential future improvements could include:

- Case-insensitive regex matching option
- Support for multiple exclusion patterns via comma-separated values
- Exclusion filtering for other entity types (models, relationships)
- Pre-compiled regex patterns for commonly used exclusions
