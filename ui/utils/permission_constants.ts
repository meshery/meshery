/**
 * Constants for Identity & Access Management Keys
 *
 * This file acts as a bridge between the legacy SCREAMING_SNAKE_CASE monikers
 * and the canonical, schema-aligned PascalCase permission definitions exported
 * by @meshery/schemas.
 *
 * All UUIDs and subject strings are resolved dynamically at runtime from
 * the schemas package. Do not hardcode UUIDs here.
 */

import { Keys } from '@meshery/schemas/permissions';

export type { PermissionKey } from '@meshery/schemas/permissions';

const legacyToPascal = {
  EDIT_ORGANIZATION: 'IdentityAccessManagementEditOrganization',
  VIEW_ALL_ORGANIZATIONS: 'IdentityAccessManagementViewAllOrganizations',
  VIEW_PROFILE: 'AccountManagementViewProfile',
  EDIT_ACCOUNT: 'AccountManagementEditAccount',
  VIEW_CREDENTIALS: 'SecurityManagementViewCredentials',
  CREATE_CREDENTIAL: 'SecurityManagementCreateCredential',
  EDIT_CREDENTIAL: 'SecurityManagementEditCredential',
  DELETE_CREDENTIAL: 'SecurityManagementDeleteCredential',
  VIEW_SESSIONS: 'SecurityManagementViewSessions',
  LOGOUT_SESSION_FROM_A_SESSION: 'SecurityManagementLogoutFromASession',
  VIEW_TOKENS: 'SecurityManagementViewTokens',
  DOWNLOAD_TOKEN: 'SecurityManagementDownloadToken',
  CREATE_TOKEN: 'SecurityManagementCreateToken',
  DELETE_TOKEN: 'SecurityManagementDeleteToken',
  VIEW_WORKSPACE: 'WorkspaceManagementViewWorkspace',
  CREATE_WORKSPACE: 'WorkspaceManagementCreateWorkspace',
  DELETE_WORKSPACE: 'WorkspaceManagementDeleteWorkspace',
  EDIT_WORKSPACE: 'WorkspaceManagementEditWorkspace',
  CONNECT_GOOGLE_ACCOUNT_TO_WORKSPACE: 'WorkspaceManagementConnectGoogleAccountToWorkspace',
  CONNECT_GITHUB_ACCOUNT_TO_WORKSPACE: 'WorkspaceManagementConnectGithubAccountToWorkspace',
  ASSIGN_TEAM_TO_WORKSPACE: 'WorkspaceManagementAssignTeamToWorkspace',
  REMOVE_TEAM_FROM_WORKSPACE: 'WorkspaceManagementRemoveTeamFromWorkspace',
  ASSIGN_DESIGNS_TO_WORKSPACE: 'WorkspaceManagementAssignDesignsToWorkspaces',
  REMOVE_DESIGNS_FROM_WORKSPACE: 'WorkspaceManagementRemoveDesignsFromWorkspaces',
  ASSIGN_VIEWS_TO_WORKSPACE: 'KanvasAssignViewsToWorkspace',
  REMOVE_VIEWS_FROM_WORKSPACE: 'KanvasUnassignViewsFromWorkspace',
  VIEW_TEAMS: 'IdentityAccessManagementViewTeams',
  DELETE_TEAM: 'IdentityAccessManagementDeleteTeam',
  CREATE_TEAM: 'IdentityAccessManagementCreateTeam',
  EDIT_TEAM: 'IdentityAccessManagementEditTeam',
  LEAVE_TEAM: 'IdentityAccessManagementLeaveTeam',
  ASSIGN_ENVIRONMENT_TO_WORKSPACE: 'WorkspaceManagementAssignEnvironmentToWorkspace',
  REMOVE_ENVIRONMENT_FROM_WORKSPACE: 'WorkspaceManagementRemoveEnvironmentFromWorkspace',
  VIEW_AUDIT: 'EventsManagementViewAudit',
  VIEW_PROJECTS: 'WorkspaceManagementViewProjects',
  VIEW_CONNECTIONS: 'WorkspaceManagementViewConnections',
  VIEW_ENVIRONMENTS: 'WorkspaceManagementViewEnvironment',
  CREATE_ENVIRONMENT: 'WorkspaceManagementCreateEnvironment',
  EDIT_ENVIRONMENT: 'WorkspaceManagementEditEnvironment',
  DELETE_ENVIRONMENT: 'WorkspaceManagementDeleteEnvironment',
  ASSIGN_CONNECTIONS_TO_ENVIRONMENT: 'WorkspaceManagementAssignConnectionsToEnvironment',
  REMOVE_CONNECTIONS_FROM_ENVIRONMENT: 'WorkspaceManagementRemoveConnectionsFromEnvironments',
  VIEW_CATALOG: 'CatalogManagementViewCatalog',
  VIEW_APPLICATIONS: 'CatalogManagementViewApplications',
  VIEW_DESIGNS: 'CatalogManagementViewDesigns',
  SHARE_DESIGN: 'CatalogManagementShareDesign',
  CLONE_DESIGN: 'CatalogManagementCloneDesign',
  DOWNLOAD_DESIGN: 'CatalogManagementDownloadADesign',
  OPEN_IN_PLAYGROUND: 'CatalogManagementOpenInPlayground',
  VIEW_FILTERS: 'CatalogManagementViewFilters',
  CREATE_NEW_DESIGN: 'CatalogManagementCreateNewDesign',
  IMPORT_DESIGN: 'CatalogManagementImportDesign',
  PUBLISH_DESIGN: 'CatalogManagementPublishDesign',
  UNPUBLISH_DESIGN: 'CatalogManagementUnpublishDesign',
  VALIDATE_DESIGN: 'CatalogManagementValidateDesign',
  DEPLOY_DESIGN: 'CatalogManagementDeployDesign',
  UNDEPLOY_DESIGN: 'CatalogManagementUndeployDesign',
  DETAILS_OF_DESIGN: 'CatalogManagementDetailsOfDesign',
  EDIT_DESIGN: 'CatalogManagementEditDesign',
  DELETE_A_DESIGN: 'CatalogManagementDeleteADesign',
  EVALUATE_RELATIONSHIPS: 'CatalogManagementEvaluateRelationships',
  EXPORT_DESIGN: 'CatalogManagementExportDesign',
  DOWNLOAD_A_DESIGN: 'CatalogManagementDownloadADesign',
  IMPORT_FILTER: 'CatalogManagementImportFilter',
  PUBLISH_WASM_FILTER: 'CatalogManagementPublishWasmFilter',
  UNPUBLISH_WASM_FILTER: 'CatalogManagementUnpublishWasmFilter',
  DOWNLOAD_A_WASM_FILTER: 'CatalogManagementDownloadAWasmFilter',
  DETAILS_OF_WASM_FILTER: 'CatalogManagementDetailsOfWasmFilter',
  EDIT_WASM_FILTER: 'CatalogManagementEditWasmFilter',
  DELETE_WASM_FILTER: 'CatalogManagementDeleteWasmFilter',
  CLONE_WASM_FILTER: 'CatalogManagementCloneWasmFilter',
  ADD_CLUSTER: 'LifecycleManagementAddCluster',
  CHANGE_CONNECTION_STATE: 'LifecycleManagementChangeConnectionState',
  FLUSH_MESHSYNC_DATA: 'LifecycleManagementFlushMeshsyncData',
  REGISTER_DISCOVERED_MESHSYNC_RESOURCE: 'LifecycleManagementRegisterDiscoveredMeshsyncResource',
  DELETE_A_CONNECTION: 'LifecycleManagementDeleteAConnection',
  INSTALL_EXTENSION: 'ExtensibilityInstallExtension',
  // No separate "uninstall" permission exists in schemas;
  // ExtensibilityInstallExtension covers install, enable, and disable.
  UNINSTALL_EXTENSION: 'ExtensibilityInstallExtension',
  VIEW_EXTENSIONS: 'ExtensibilityViewExtensions',
  VIEW_MESHERY_USER_PREFERENCES: 'ExtensibilityViewMesheryUserPreferences',
  ASSIGN_USER_ROLES: 'IdentityAccessManagementAssignUserRoles',
  VIEW_SETTINGS: 'MesherySystemViewSettings',
  CONNECT_ADAPTER: 'MesherySystemConnectAdapter',
  CONNECT_METRICS: 'MesherySystemConnectMetrics',
  VIEW_METRICS: 'MesherySystemViewMetrics',
  VIEW_REGISTRY: 'MesherySystemViewRegistry',
  ADD_PERFORMANCE_PROFILE: 'PerformanceManagementAddPerformaceProfile',
  RUN_TEST: 'PerformanceManagementRunTest',
  VIEW_RESULTS: 'PerformanceManagementViewResults',
  EDIT_PERFORMANCE_TEST: 'PerformanceManagementEditPerformanceTest',
  DELETE_PERFORMANCE_TEST: 'PerformanceManagementDeletePerformanceTest',
  VIEW_VIEWS: 'KanvasViewViews',
  CREATE_VIEW: 'KanvasCreateView',
  DELETE_VIEW: 'KanvasDeleteView',
  EDIT_VIEW: 'KanvasEditView',
  INVITE_COLLABORATORS_PUBLIC_DESIGNS: 'CollaborationInviteAnyMesheryCloudUserOrAllMesheryUsers',
  INVITE_COLLABORATORS_PRIVATE_DESIGNS: 'CollaborationInviteAnyMesheryCloudUserToOnAPrivateDesign',
  MESSAGE_IN_REAL_TIME: 'CollaborationMessageInRealTime',
  DISCUSS_ANY_DESIGN_BY_LEAVING_REVIEW_COMMENTS:
    'CollaborationDiscussAnyDesignByLeavingReviewComments',
  MANAGE_ACCESS_TO_DESIGN: 'CollaborationManageAccessToDesigns',
  CREATE_AND_COLLABORATE_IN_ONLINE_DESIGNS_IN_REAL_TIME:
    'CollaborationCreateAndCollaborateInOnlineDesignsInRealTime',
  MANAGE_CLOUD_NATIVE_INFRASTRUCTURE_LIFE_CYCLE:
    'InfrastructureManagementManageCloudNativeInfrastructureLifeCycle',
  MANAGE_CLOUD_NATIVE_INFRASTRUCTURE_CONFIGURATION:
    'InfrastructureManagementManageCloudNativeInfrastructureConfiguration',
  APPLY_CLOUD_NATIVE_INFRASTRUCTURE_CONFIGURATION:
    'InfrastructureManagementApplyCloudNativeInfrastructureConfiguration',
  VALIDATE_CLOUD_NATIVE_INFRASTRUCTURE_CONFIGURATION:
    'InfrastructureManagementValidateCloudNativeInfrastructureConfiguration',
  APPLY_CUSTOM_CLOUD_NATIVE_CONFIGURATION:
    'InfrastructureManagementApplyCustomCloudNativeConfiguration',
  DEPLOY_CLOUD_NATIVE_INFRASTRUCTURE: 'InfrastructureManagementDeployCloudNativeInfrastructure',
  UNDEPLOY_CLOUD_NATIVE_INFRASTRUCTURE: 'InfrastructureManagementUndeployCloudNativeInfrastructure',
  VIEW_CLOUD_NATIVE_INFRASTRUCTURE: 'InfrastructureManagementViewCloudNativeInfrastructure',
  VIEW_ALL_KUBERNETES_CLUSTERS: 'IdentityAccessManagementViewAllKubernetesClusters',
  VIEW_PERFORMANCE_PROFILES: 'PerformanceManagementViewPerformanceProfiles',
  RESET_DATABASE: 'MesherySystemResetDatabase',
} as const satisfies Record<string, keyof typeof Keys>;

type UiPermissionKey = { action: string; subject: string };
type UiPermissionKeyMap = Record<string, UiPermissionKey>;

// Map the new Keys to the legacy action/subject shape for the UI
const mappedKeys: UiPermissionKeyMap = Object.fromEntries(
  Object.entries(Keys).map(([name, keyObj]) => [
    name,
    {
      action: keyObj.id,
      subject: keyObj.function,
    },
  ]),
);

// Create legacy keys mapped to schema details
const legacyKeys = (Object.entries(legacyToPascal) as Array<[string, keyof typeof Keys]>).reduce(
  (acc, [legacy, pascal]) => {
    const keyObj = Keys[pascal];
    if (!keyObj) {
      const msg = `Unknown permission key mapping: ${legacy} -> ${pascal}`;
      if (process.env.NODE_ENV !== 'production') {
        throw new Error(msg);
      }
      // Keep a visible, non-empty sentinel in production to avoid silent failures.
      acc[legacy] = { action: `MISSING:${pascal}`, subject: msg };
      return acc;
    }
    acc[legacy] = {
      action: keyObj.id,
      subject: keyObj.function,
    };
    return acc;
  },
  {} as UiPermissionKeyMap,
);

export const keys = {
  ...mappedKeys,
  ...legacyKeys,
};
