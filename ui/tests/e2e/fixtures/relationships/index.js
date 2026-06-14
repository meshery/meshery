import serviceDeploymentNetworkFixture from './service-deployment-edge-non-binding-network-fixture.json';
import configmapPodReferenceFixture from './configmap-pod-edge-non-binding-reference-fixture.json';
import namespaceHierarchicalParentInventoryFixture from './namespace-hierarchical-parent-inventory-fixture.json';
import containerHierarchicalParentAliasFixture from './container-hierarchical-parent-alias-fixture.json';
import mesheryDesignFixture from './meshery-design-fixture.json';
import pvPvcReferenceFixture from './pv-pvc-edge-non-binding-reference-fixture.json';
import secretDeploymentReferenceFixture from './secret-deployment-edge-non-binding-reference-fixture.json';
import roleRolebindingPermissionFixture from './role-rolebinding-serviceaccount-edge-binding-permission-fixture.json';
import namespaceNamespaceDenyFixture from './namespace-namespace-deny-hierarchical-parent-inventory-fixture.json';

export const RelationshipTestFixtures = [
  serviceDeploymentNetworkFixture,
  namespaceHierarchicalParentInventoryFixture,
  configmapPodReferenceFixture,
  containerHierarchicalParentAliasFixture,
  mesheryDesignFixture,
  pvPvcReferenceFixture,
  secretDeploymentReferenceFixture,
  roleRolebindingPermissionFixture,
  namespaceNamespaceDenyFixture,
];
