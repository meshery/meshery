import { describe, expect, it } from 'vitest';
import {
  ACTIONS,
  APP_MODE,
  CONNECTION_KINDS,
  CONNECTION_KINDS_DEF,
  CONNECTION_STATES,
  CONNECTION_STATE_TO_TRANSITION_MAP,
  CON_OPS,
  CONTROLLERS,
  CONTROLLER_STATES,
  DEPLOYMENT_TYPE,
  EVENT_TYPES,
  EXTENSIONS,
  FILE_OPS,
  MESHSYNC_DEPLOYMENT_TYPE,
  MESHSYNC_STATES,
  MesheryFiltersCatalog,
  MesheryPatternsCatalog,
  REGISTRY_ITEM_STATES,
  REGISTRY_ITEM_STATES_TO_TRANSITION_MAP,
  RESOURCE_TYPE,
  TRANSFER_COMPONENT,
  VIEW_VISIBILITY,
  VISIBILITY,
} from '../Enum';

describe('Enum constants', () => {
  it('FILE_OPS exposes file operation labels', () => {
    expect(FILE_OPS).toEqual({
      FILE_UPLOAD: 'upload',
      URL_UPLOAD: 'url_upload',
      UPDATE: 'update',
      DELETE: 'delete',
      DOWNLOAD: 'download',
      CLONE: 'clone',
    });
  });

  it('CON_OPS exposes connection operation labels', () => {
    expect(CON_OPS.DELETE).toBe('delete');
    expect(CON_OPS.UPDATE).toBe('update');
    expect(CON_OPS.CREATE).toBe('create');
  });

  it('ACTIONS uses numeric codes for deploy/undeploy/verify', () => {
    expect(ACTIONS.DEPLOY).toBe(2);
    expect(ACTIONS.UNDEPLOY).toBe(1);
    expect(ACTIONS.VERIFY).toBe(0);
  });

  it('DEPLOYMENT_TYPE has expected string values', () => {
    expect(DEPLOYMENT_TYPE.IN_CLUSTER).toBe('in_cluster');
    expect(DEPLOYMENT_TYPE.OUT_CLUSTER).toBe('out_of_cluster');
  });

  it('VISIBILITY values are lowercase strings', () => {
    expect(VISIBILITY).toEqual({
      PRIVATE: 'private',
      PUBLIC: 'public',
      PUBLISHED: 'published',
    });
  });

  it('EVENT_TYPES values are uppercase', () => {
    expect(EVENT_TYPES.ADDED).toBe('ADDED');
    expect(EVENT_TYPES.DELETED).toBe('DELETED');
    expect(EVENT_TYPES.MODIFIED).toBe('MODIFIED');
  });

  it('exposes the /extension/meshmap page as an available extension', () => {
    const meshmapPages = Object.values(EXTENSIONS).filter(
      (ext) => 'signup_url' in ext && ext.signup_url === '/extension/meshmap',
    );
    expect(meshmapPages).toHaveLength(1);
  });

  it('REGISTRY_ITEM_STATES and transition map line up by key', () => {
    expect(REGISTRY_ITEM_STATES.ENABLED).toBe('enabled');
    expect(REGISTRY_ITEM_STATES.IGNORED).toBe('ignored');

    expect(REGISTRY_ITEM_STATES_TO_TRANSITION_MAP[REGISTRY_ITEM_STATES.ENABLED]).toBe('Enable');
    expect(REGISTRY_ITEM_STATES_TO_TRANSITION_MAP[REGISTRY_ITEM_STATES.IGNORED]).toBe('Ignore');
  });

  it('CONNECTION_STATES covers every documented lifecycle phase', () => {
    expect(CONNECTION_STATES).toMatchObject({
      DISCOVERED: 'discovered',
      REGISTERED: 'registered',
      CONNECTED: 'connected',
      IGNORED: 'ignored',
      MAINTENANCE: 'maintenance',
      DISCONNECTED: 'disconnected',
      DELETED: 'deleted',
      NOTFOUND: 'not found',
    });
  });

  it('CONNECTION_STATE_TO_TRANSITION_MAP maps each known state to a human label', () => {
    expect(CONNECTION_STATE_TO_TRANSITION_MAP[CONNECTION_STATES.IGNORED]).toBe('Ignore');
    expect(CONNECTION_STATE_TO_TRANSITION_MAP[CONNECTION_STATES.CONNECTED]).toBe('Connect');
    expect(CONNECTION_STATE_TO_TRANSITION_MAP[CONNECTION_STATES.REGISTERED]).toBe('Register');
    expect(CONNECTION_STATE_TO_TRANSITION_MAP[CONNECTION_STATES.DISCOVERED]).toBe('Discover');
    expect(CONNECTION_STATE_TO_TRANSITION_MAP[CONNECTION_STATES.DELETED]).toBe('Delete');
    expect(CONNECTION_STATE_TO_TRANSITION_MAP[CONNECTION_STATES.MAINTENANCE]).toBe('Maintenance');
    expect(CONNECTION_STATE_TO_TRANSITION_MAP[CONNECTION_STATES.DISCONNECTED]).toBe('Disconnect');
    expect(CONNECTION_STATE_TO_TRANSITION_MAP[CONNECTION_STATES.NOTFOUND]).toBe('Not Found');
  });

  it('CONTROLLERS / CONTROLLER_STATES expose the documented constants', () => {
    expect(CONTROLLERS.BROKER).toBe('BROKER');
    expect(CONTROLLERS.OPERATOR).toBe('OPERATOR');
    expect(CONTROLLERS.MESHSYNC).toBe('MESHSYNC');

    expect(CONTROLLER_STATES.DEPLOYED).toBe('DEPLOYED');
    expect(CONTROLLER_STATES.UNKNOWN).toBe('UNKNOWN');
    expect(CONTROLLER_STATES.UNKOWN).toBe('UNKOWN');
  });

  it('catalog identifier constants are wire-format strings', () => {
    expect(MesheryPatternsCatalog).toBe('meshery-patterns-catalog');
    expect(MesheryFiltersCatalog).toBe('meshery-filters-catalog');
  });

  it('CONNECTION_KINDS_DEF and CONNECTION_KINDS describe the same provider set', () => {
    expect(CONNECTION_KINDS_DEF).toEqual([
      'MESHERY',
      'KUBERNETES',
      'PROMETHEUS',
      'GRAFANA',
      'GITHUB',
    ]);
    expect(Object.keys(CONNECTION_KINDS).sort()).toEqual([...CONNECTION_KINDS_DEF].sort());
    expect(CONNECTION_KINDS.MESHERY).toBe('meshery');
    expect(CONNECTION_KINDS.KUBERNETES).toBe('kubernetes');
  });

  it('MESHSYNC enums expose deployment/state strings', () => {
    expect(MESHSYNC_DEPLOYMENT_TYPE.OPERATOR).toBe('operator');
    expect(MESHSYNC_DEPLOYMENT_TYPE.EMBEDDED).toBe('embedded');
    expect(MESHSYNC_STATES.DISCOVERED).toBe('discovered');
    expect(MESHSYNC_STATES.REGISTER).toBe('register');
  });

  it('TRANSFER_COMPONENT distinguishes chip vs other transfer payloads', () => {
    expect(TRANSFER_COMPONENT.CHIP).toBe('chip');
    expect(TRANSFER_COMPONENT.OTHER).toBe('other');
  });

  it('EXTENSIONS map contains the registered extension entries keyed by EXTENSION_NAMES', () => {
    expect(EXTENSIONS.Catalog.name).toBe('Meshery Catalog');
    expect(EXTENSIONS.PerformanceAnalysis.name).toBe('Performance Analysis');
  });

  it('RESOURCE_TYPE / APP_MODE / VIEW_VISIBILITY expose lowercase values', () => {
    expect(RESOURCE_TYPE).toEqual({
      FILTER: 'filter',
      DESIGN: 'design',
      CATALOG: 'catalog',
      VIEW: 'view',
    });

    expect(APP_MODE.DESIGN).toBe('design');
    expect(APP_MODE.OPERATOR).toBe('operator');

    expect(VIEW_VISIBILITY.PUBLIC).toBe('public');
    expect(VIEW_VISIBILITY.PRIVATE).toBe('private');
  });
});
