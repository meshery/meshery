/**
 * Constants for Identity & Access Management Keys
 */

export const keys = {
  VIEW_ALL_ORGANIZATIONS: {
    subject: "View All Organizations",
    action: "e996c998-a50f-4cb8-ae7b-f2f1b523c971"
  },
  VIEW_PROFILE: {
    subject: "View Profile",
    action: "fa7de118-2d08-4b07-b9d7-3e0baead6d04"
  },
  EDIT_ACCOUNT: {
    subject: "Edit Account",
    action: "f5e6bb39-c89a-4172-86f9-14a4a59792c1"
  },
  VIEW_CREDENTIALS: {
    subject: "View Credentials",
    action: "96759f76-4add-45f8-b4ef-d4ace5ab1bc4"
  },
  CREATE_CREDENTIAL: {
    subject: "Create Credential",
    action: "30023b1b-01a7-4613-8364-38d3487d1789"
  },
  EDIT_CREDENTIAL: {
    subject: "Edit Credential",
    action: "e4cd5bb0-8afb-4b35-8716-0e2ead13c9b7"
  },
  DELETE_CREDENTIAL: {
    subject: "Delete Credential",
    action: "cb09f530-aa87-4a18-b3d3-bbcc2d6ca1a6"
  },
  VIEW_SESSIONS: {
    subject: "View Sessions",
    action: "26cf042a-91db-4237-8644-4d617a0d49e1"
  },
  LOGOUT_SESSION: {
    subject: "Logout from a Session",
    action: "177b928b-71ee-4ecd-a30b-3154ff4ba0d9"
  },
  VIEW_TOKENS: {
    subject: "View Tokens",
    action: "46d914bc-18c1-438f-aa74-fb78823aa25c"
  },
  DOWNLOAD_TOKEN: {
    subject: "Download Token",
    action: "ee5fc23e-d629-4c7b-8169-27e526394e8b"
  },
  CREATE_TOKEN: {
    subject: "Create Token",
    action: "8aa0df56-57e8-44b7-9d6e-7df413048ed5"
  },
  DELETE_TOKEN: {
    subject: "Delete Token",
    action: "ddba5064-ac3c-470e-b405-d2a0e99db477"
  },
  VIEW_WORKSPACE: {
    subject: "View Workspace",
    action: "bc9379e8-dc18-4655-b53c-c641271c4ba3"
  },
  CREATE_WORKSPACE: {
    subject: "Create Workspace",
    action: "eb42ac41-a883-465e-843c-d64e962a3a0e"
  },
  DELETE_WORKSPACE: {
    subject: "Delete Workspace",
    action: "09eb0507-2f14-4bc4-92c5-9e26a4efbd5e"
  },
  EDIT_WORKSPACE: {
    subject: "Edit Workspace",
    action: "4112230f-5d1e-4d30-9790-942ad5c1dc50"
  },
  CONNECT_GOOGLE_ACCOUNT_TO_WORKSPACE: {
    subject: "Connect Google Account to Workspace",
    action: "69179641-6c41-40d8-87a0-81dd99bcb396"
  },
  CONNECT_GITHUB_ACCOUNT_TO_WORKSPACE: {
    subject: "Connect Github Account to Workspace",
    action: "410b2d3c-8194-44d1-9f80-7b5fea689b4f"
  },
  ASSIGN_TEAM_TO_WORKSPACE: {
    subject: "Assign team to workspace",
    action: "6ab4263b-0bb3-492e-9878-6936a5b6312f"
  },
  REMOVE_TEAM_FROM_WORKSPACE: {
    subject: "Remove team from workspace",
    action: "c4ed82f5-783d-4451-9b34-44f50cae71df"
  },
  ASSIGN_ENVIRONMENT_TO_WORKSPACE: {
    subject: "Assign environment to workspace",
    action: "f421fc20-c14a-4282-b526-776c6cacfd99"
  },
  REMOVE_ENVIRONMENT_FROM_WORKSPACE: {
    subject: "Remove environment from workspace",
    action: "d0657715-80fb-4b00-af27-b78bb0fa56df"
  },
  VIEW_AUDIT: {
    subject: "View Audit",
    action: "80bb9c66-0657-49ff-a064-667e9875bb3f"
  },
  VIEW_PROJECTS: {
    subject: "View Projects",
    action: "141a5f3d-b5e2-4f36-8f83-df7f73744ee1"
  },
  VIEW_CONNECTIONS: {
    subject: "View Connections",
    action: "b35c9ce0-e787-4de6-8560-631007b0b947"
  },
  VIEW_ENVIRONMENTS: {
    subject: "View Environment",
    action: "e3656bbc-fba2-483d-9996-34f8614cd21b"
  },
  CREATE_ENVIRONMENT: {
    subject: "Create Environment",
    action: "a97b7f3b-3349-4a86-b917-2ce0b64a540b"
  },
  EDIT_ENVIRONMENT: {
    subject: "Edit Environment",
    action: "145ab6ed-b4b6-4e34-ada5-78dada250f89"
  },
  DELETE_ENVIRONMENT: {
    subject: "Delete Environment",
    action: "70747966-dfad-4523-93ce-bd7421258955"
  },
  ASSIGN_CONNECTIONS_TO_ENVIRONMENT: {
    subject: "Assign Connections To Environment",
    action: "52cbe0b8-9aa7-4605-8eed-aa37e595adbb"
  },
  REMOVE_CONNECTIONS_FROM_ENVIRONMENT: {
    subject: "Remove Connections From Environment",
    action: "65648682-e47f-43d7-a5ad-dc042803f951"
  },
  VIEW_CATALOG: {
    subject: "View Catalog",
    action: "0cd05106-36b6-4393-a08e-4222fc10c8de"
  },
  VIEW_APPLICATIONS: {
    subject: "View Applications",
    action: "bfb200b6-0ba9-4783-95d4-eaf1c8fe004c"
  },
  VIEW_DESIGNS: {
    subject: "View Designs",
    action: "3798736d-1f5d-41b3-876f-f3f01453dd15"
  },
  SHARE_DESIGN: {
    subject: "Share Design", // not seeded
    action: "d9ae2b08-762f-418f-916f-43de736b53e2"
  },
  DOWNLOAD_DESIGN: {
    subject: "Download Design", // not seeded
    action: "7a7d3cad-4cf4-48fa-b69c-e6a5abc97a07"
  },
  CLONE_DESIGN: {
    subject: "Clone Design", // not seeded
    action: "94a12f80-3c45-4a1f-afb2-a68b909d0d7f"
  },
  OPEN_IN_PLAYGROUND: {
    subject: "Open in Playground", // not seeded
    action: "c4d6c676-6e26-4b0c-9fdd-5eea1b780e98"
  },
  VIEW_FILTER: {
    subject: "View Filters",
    action: "df41c45f-7c73-49c2-a055-0584fdcec1c1"
  },
};
