/**
 * Constants for Identity & Access Management Keys
 * When add new key make sure that subject property matches with key function property
 */

export const keys = {
  VIEW_ALL_ORGANIZATIONS: {
    subject: 'View All Organizations',
    action: 'e996c998-a50f-4cb8-ae7b-f2f1b523c971',
  },
  VIEW_PROFILE: {
    subject: 'View Profile',
    action: 'fa7de118-2d08-4b07-b9d7-3e0baead6d04',
  },
  EDIT_ACCOUNT: {
    subject: 'Edit Account',
    action: 'f5e6bb39-c89a-4172-86f9-14a4a59792c1',
  },
  VIEW_CREDENTIALS: {
    subject: 'View Credentials',
    action: '96759f76-4add-45f8-b4ef-d4ace5ab1bc4',
  },
  CREATE_CREDENTIAL: {
    subject: 'Create Credential',
    action: '30023b1b-01a7-4613-8364-38d3487d1789',
  },
  EDIT_CREDENTIAL: {
    subject: 'Edit Credential',
    action: 'e4cd5bb0-8afb-4b35-8716-0e2ead13c9b7',
  },
  DELETE_CREDENTIAL: {
    subject: 'Delete Credential',
    action: 'cb09f530-aa87-4a18-b3d3-bbcc2d6ca1a6',
  },
  VIEW_SESSIONS: {
    subject: 'View Sessions',
    action: '26cf042a-91db-4237-8644-4d617a0d49e1',
  },
  LOGOUT_SESSION_FROM_A_SESSION: {
    subject: 'Logout from a Session',
    action: '177b928b-71ee-4ecd-a30b-3154ff4ba0d9',
  },
  VIEW_TOKENS: {
    subject: 'View Tokens',
    action: '46d914bc-18c1-438f-aa74-fb78823aa25c',
  },
  DOWNLOAD_TOKEN: {
    subject: 'Download Token',
    action: 'ee5fc23e-d629-4c7b-8169-27e526394e8b',
  },
  CREATE_TOKEN: {
    subject: 'Create Token',
    action: '8aa0df56-57e8-44b7-9d6e-7df413048ed5',
  },
  DELETE_TOKEN: {
    subject: 'Delete Token',
    action: 'ddba5064-ac3c-470e-b405-d2a0e99db477',
  },
  VIEW_WORKSPACE: {
    subject: 'View Workspace',
    action: 'bc9379e8-dc18-4655-b53c-c641271c4ba3',
  },
  CREATE_WORKSPACE: {
    subject: 'Create Workspace',
    action: 'eb42ac41-a883-465e-843c-d64e962a3a0e',
  },
  DELETE_WORKSPACE: {
    subject: 'Delete Workspace',
    action: '09eb0507-2f14-4bc4-92c5-9e26a4efbd5e',
  },
  EDIT_WORKSPACE: {
    subject: 'Edit Workspace',
    action: '4112230f-5d1e-4d30-9790-942ad5c1dc50',
  },
  CONNECT_GOOGLE_ACCOUNT_TO_WORKSPACE: {
    subject: 'Connect Google Account to Workspace',
    action: '69179641-6c41-40d8-87a0-81dd99bcb396',
  },
  CONNECT_GITHUB_ACCOUNT_TO_WORKSPACE: {
    subject: 'Connect Github Account to Workspace',
    action: '410b2d3c-8194-44d1-9f80-7b5fea689b4f',
  },
  ASSIGN_TEAM_TO_WORKSPACE: {
    subject: 'Assign team to workspace',
    action: '6ab4263b-0bb3-492e-9878-6936a5b6312f',
  },
  REMOVE_TEAM_FROM_WORKSPACE: {
    subject: 'Remove team from workspace',
    action: 'c4ed82f5-783d-4451-9b34-44f50cae71df',
  },
  ASSIGN_ENVIRONMENT_TO_WORKSPACE: {
    subject: 'Assign environment to workspace',
    action: 'f421fc20-c14a-4282-b526-776c6cacfd99',
  },
  REMOVE_ENVIRONMENT_FROM_WORKSPACE: {
    subject: 'Remove environment from workspace',
    action: 'd0657715-80fb-4b00-af27-b78bb0fa56df',
  },
  VIEW_AUDIT: {
    subject: 'View Audit',
    action: '80bb9c66-0657-49ff-a064-667e9875bb3f',
  },
  VIEW_PROJECTS: {
    subject: 'View Projects',
    action: '141a5f3d-b5e2-4f36-8f83-df7f73744ee1',
  },
  VIEW_CONNECTIONS: {
    subject: 'View Connections',
    action: 'b35c9ce0-e787-4de6-8560-631007b0b947',
  },
  VIEW_ENVIRONMENTS: {
    subject: 'View Environment',
    action: 'e3656bbc-fba2-483d-9996-34f8614cd21b',
  },
  CREATE_ENVIRONMENT: {
    subject: 'Create Environment',
    action: 'a97b7f3b-3349-4a86-b917-2ce0b64a540b',
  },
  EDIT_ENVIRONMENT: {
    subject: 'Edit Environment',
    action: '145ab6ed-b4b6-4e34-ada5-78dada250f89',
  },
  DELETE_ENVIRONMENT: {
    subject: 'Delete Environment',
    action: '70747966-dfad-4523-93ce-bd7421258955',
  },
  ASSIGN_CONNECTIONS_TO_ENVIRONMENT: {
    subject: 'Assign Connections To Environment',
    action: '52cbe0b8-9aa7-4605-8eed-aa37e595adbb',
  },
  REMOVE_CONNECTIONS_FROM_ENVIRONMENT: {
    subject: 'Remove Connections From Environment',
    action: '65648682-e47f-43d7-a5ad-dc042803f951',
  },
  VIEW_CATALOG: {
    subject: 'View Catalog',
    action: '0cd05106-36b6-4393-a08e-4222fc10c8de',
  },
  VIEW_APPLICATIONS: {
    subject: 'View Applications',
    action: 'bfb200b6-0ba9-4783-95d4-eaf1c8fe004c',
  },
  VIEW_DESIGNS: {
    subject: 'View Designs',
    action: '3798736d-1f5d-41b3-876f-f3f01453dd15',
  },
  SHARE_DESIGN: {
    subject: 'Share Design',
    action: 'd9ae2b08-762f-418f-916f-43de736b53e2',
  },
  CLONE_DESIGN: {
    subject: 'Clone Design',
    action: '94a12f80-3c45-4a1f-afb2-a68b909d0d7f',
  },
  OPEN_IN_PLAYGROUND: {
    subject: 'Open in Playground', // not seeded
    action: 'c4d6c676-6e26-4b0c-9fdd-5eea1b780e98',
  },
  VIEW_FILTERS: {
    subject: 'View Filters',
    action: 'df41c45f-7c73-49c2-a055-0584fdcec1c1',
  },
  CREATE_NEW_DESIGN: {
    subject: 'Create new design',
    action: '14bd933e-83b7-464d-9a4d-d8c8eb9682ab',
  },
  IMPORT_DESIGN: {
    subject: 'Import Design',
    action: 'cc040d21-3160-4a96-8efa-833487a234cd',
  },
  PUBLISH_DESIGN: {
    subject: 'Publish Design',
    action: '9e66bdec-4177-42f9-8cec-d9eb52a12c38',
  },
  UNPUBLISH_DESIGN: {
    subject: 'Unpublish Design',
    action: 'c1595c90-b85b-4ac7-b921-f08959926db3',
  },
  VALIDATE_DESIGN: {
    subject: 'Validate Design',
    action: 'da5339dd-a4bc-4b91-8865-d8a703656516',
  },
  DEPLOY_DESIGN: {
    subject: 'Deploy Design',
    action: '595b921a-ea1e-4611-83f0-503db0eeb94d',
  },
  UNDEPLOY_DESIGN: {
    subject: 'Undeploy Design',
    action: '16b11ffa-7b92-4666-a0ff-191df9cd18b2',
  },
  DETAILS_OF_DESIGN: {
    subject: 'Details of design',
    action: '10a03036-53a0-40b3-9f69-6daab852e434',
  },
  EDIT_DESIGN: {
    subject: 'Edit design',
    action: '7f2b7084-4533-4824-b688-50cf35de7ef8',
  },
  DELETE_A_DESIGN: {
    subject: 'Delete a design',
    action: 'f024fcf7-3c3d-4521-b83e-6d659353ca0e',
  },
  EXPORT_DESIGN: {
    subject: 'Export Design',
    action: '9a783f51-3b4a-47a6-a02e-b0db9e78cd85',
  },
  DOWNLOAD_A_DESIGN: {
    subject: 'Download a design',
    action: '64de96b7-60db-4aab-b311-afc64066b2c4',
  },
  IMPORT_FILTER: {
    subject: 'Import Filter',
    action: 'cb79d7fb-19de-45fa-aaf5-0a0afc832bf8',
  },
  PUBLISH_WASM_FILTER: {
    subject: 'Publish WASM Filter',
    action: '173d99b7-3820-4c0c-88b2-a8455bd7a6b5',
  },
  UNPUBLISH_WASM_FILTER: {
    subject: 'Unpublish WASM Filter',
    action: '773f0a4d-ba04-40ed-9298-59ac8749804a',
  },
  DOWNLOAD_A_WASM_FILTER: {
    subject: 'Download a WASM filter',
    action: '24325b2c-5e08-4ba8-809f-8a4a1bf91084',
  },
  DETAILS_OF_WASM_FILTER: {
    subject: 'Details of WASM Filter',
    action: '86c457b5-b9ec-4223-af1f-30a5be67d69d',
  },
  EDIT_WASM_FILTER: {
    subject: 'Edit WASM filter',
    action: '88cd144e-806e-472a-a31a-ef6d64643291',
  },
  DELETE_WASM_FILTER: {
    subject: 'Delete WASM Filter',
    action: '9225d5a7-7255-49be-9233-daeabefae306',
  },
  CLONE_WASM_FILTER: {
    subject: 'Clone WASM Filter',
    action: 'c84718ca-7479-4ad9-a2b7-a5784baa51fb',
  },
  ADD_CLUSTER: {
    subject: 'Add cluster',
    action: 'fce15b20-78ac-42af-b79c-b8f19bdb0802',
  },
  CHANGE_CONNECTION_STATE: {
    subject: 'Change connection state',
    action: '14ac9622-3170-4580-8403-ed7a584f90ef',
  },
  FLUSH_MESHSYNC_DATA: {
    subject: 'Flush Meshsync data',
    action: '8dd4c54a-bccd-4fb3-a18c-269195653a91',
  },
  REGISTER_DISCOVERED_MESHSYNC_RESOURCE: {
    subject: 'Register discovered Meshsync resource',
    action: '214ad6b1-df4d-44a6-8872-8ad1f751ef68',
  },
  DELETE_A_CONNECTION: {
    subject: 'Delete a connection',
    action: '61afb8c2-cda6-4175-aad9-74ff87fed323',
  },
  INSTALL_EXTENSION: {
    subject: 'Install extension',
    action: '24f41e98-7ce1-40c4-a82d-4ae0294d237d',
  },
  UNINSTALL_EXTENSION: {
    // not seeded
    subject: 'Uninstall extension from Meshery',
    action: 'd25e0950-bcd5-4f31-9d19-26c91d17d89b',
  },
  VIEW_EXTENSIONS: {
    subject: 'View Extensions',
    action: 'c1330df4-1bbe-4d5d-8828-f4bd9ee989e5',
  },
  VIEW_MESHERY_USER_PREFERENCES: {
    subject: 'View Meshery User Preferences',
    action: 'eb1d8db3-e110-41b9-9f5f-cd271f7fc4a9',
  },
  VIEW_SETTINGS: {
    subject: 'View Settings',
    action: 'fdc038e3-1fdf-403a-af8a-53c0de8d7820',
  },
  CONNECT_ADAPTER: {
    subject: 'Connect adapter',
    action: 'c93bd211-1dac-42cc-9086-859288826d1b',
  },
  CONNECT_METRICS: {
    subject: 'Connect Metrics',
    action: 'b0aee906-c549-445f-be0c-b98b04d47d09',
  },
  VIEW_METRICS: {
    subject: 'View Metrics',
    action: '7fe36f60-fd0a-4fda-84e5-c64a04c3ad06',
  },
  VIEW_REGISTRY: {
    subject: 'View Registry',
    action: 'cc069117-08cc-44e3-9c61-ae0eeca0bcf1',
  },
  ADD_PERFORMANCE_PROFILE: {
    subject: 'Add performace profile',
    action: 'b2861578-c573-45fe-a95e-0356d56e1d1b',
  },
  RUN_TEST: {
    subject: 'Run test',
    action: '06de2b07-b4f4-4701-b87f-d92ebb66ba42',
  },
  VIEW_RESULTS: {
    subject: 'View Results',
    action: '0c757cc7-4038-4d9b-9b60-fa8d9fc9d27e',
  },
  EDIT_PERFORMANCE_TEST: {
    subject: 'Edit performance test',
    action: '33aa5c47-a8aa-4ad5-9950-7c17042c001d',
  },
  DELETE_PERFORMANCE_TEST: {
    subject: 'Delete performance test',
    action: '84aa9d3c-3d4b-4587-947d-ae17b2dcd5f5',
  },
  VIEW_VIEWS: {
    subject: 'View Views',
    action: '4b3117e6-176c-4d9b-8e09-4278d2f48280',
  },
  CREATE_VIEW: {
    subject: 'Create View',
    action: 'a7392c29-34af-4ee9-af1c-71e83876a8aa',
  },
  DELETE_VIEW: {
    subject: 'Delete View',
    action: 'b2b9c027-0f39-4a26-b3ed-5b30aaf82060',
  },
  EDIT_VIEW: {
    subject: 'Edit View',
    action: 'f2e04b59-42ee-4af6-b488-7bc7047fc01b',
  },
  ASSIGN_VIEWS_TO_WORKSPACE: {
    subject: 'Assign views to workspace',
    action: '5a26b21b-bf15-4622-9681-d907fc5154fd',
  },
  REMOVE_VIEWS_FROM_WORKSPACE: {
    subject: 'Remove views from workspace',
    action: '5b48eded-15d8-4c51-9ac4-b11097dbdd75',
  },
  INVITE_COLLABORATORS_PUBLIC_DESIGNS: {
    subject: 'Invite any Meshery Cloud user, or all Meshery users',
    action: 'ccc4bc8d-f484-42b3-8a62-2667284605c3',
  },
  INVITE_COLLABORATORS_PRIVATE_DESIGNS: {
    subject: 'Invite any Meshery Cloud user to on a private design',
    action: 'e28b851f-9a49-4ecf-a86e-493db1a27540',
  },
  MESSAGE_IN_REAL_TIME: {
    subject: 'Message in real-time',
    action: 'c42d08b2-c3e0-47b6-9e47-cfb149c0a5af',
  },
  DISCUSS_ANY_DESIGN_BY_LEAVING_REVIEW_COMMENTS: {
    subject: 'Discuss any design by leaving review comments',
    action: 'da5adf96-9fb5-49b2-a55e-dec9c9c4acba',
  },
  MANAGE_ACCESS_TO_DESIGN: {
    subject: 'Manage access to designs',
    action: '7a17c8d3-bba2-474b-bb1e-be5b5eee5dad',
  },
  CREATE_AND_COLLABORATE_IN_ONLINE_DESIGNS_IN_REAL_TIME: {
    // seed new function name for this key
    subject: 'Create and collaborate in online designs in real-time',
    action: 'd5267c04-b3ee-43fe-8b97-2a3321eb7f8e',
  },
  MANAGE_SERVICE_MESH_LIFE_CYCLE: {
    subject: 'Manage service mesh life cycle',
    action: '255fd148-e3fd-4408-a48c-0d157a57d4d9',
  },
  MANAGE_SERVICE_MESH_CONFIGURATION: {
    subject: 'Manage service mesh configuration',
    action: '0eb0558d-9b21-4e50-b4c6-bd8e9e3414f5',
  },
  APPLY_SERVICE_MESH_CONFIGURATION: {
    subject: 'Apply service mesh configuration',
    action: '3f20a106-24f5-4da6-a8eb-6eddaad50944',
  },
  VALIDATE_SERVICE_MESH_CONFIGURATION: {
    subject: 'Validate service mesh configuration',
    action: '8bb93f97-fcfb-4827-9fed-f931fdca7b95',
  },
  APPLY_CUSTOM_SERVICE_MESH_CONFIGURATION: {
    subject: 'Apply custom service mesh configuration',
    action: '2f4e2300-4c7e-4d48-95aa-74614a4826fe',
  },
  DEPLOY_SERVICE_MESH: {
    subject: 'Deploy service mesh',
    action: 'f7e70ffb-333d-43b3-a76e-0e6c63b9fbfa',
  },
  UNDEPLOY_SERVICE_MESH: {
    subject: 'Undeploy service mesh',
    action: '6e7f6f4f-4321-4e42-9eff-6a8323f32e84',
  },
  VIEW_SERVICE_MESH: {
    subject: 'View service mesh',
    action: 'fdc485dc-f68b-405c-9e54-7b9a7254c282',
  },
  VIEW_ALL_KUBERNETES_CLUSTERS: {
    subject: 'View All Kubernetes Clusters',
    action: 'b99a9a0a-2cb9-4be7-8251-14a249e4038e',
  },
  VIEW_PERFORMANCE_PROFILES: {
    subject: 'View performance profiles',
    action: '6593ac26-820b-4e87-be32-64ee740ea204',
  },
};
