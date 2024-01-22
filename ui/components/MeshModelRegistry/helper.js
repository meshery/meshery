import { MODELS, REGISTRANTS } from '@/constants/navigator';

/**
 * Retrieves filtered data for the details component based on the selected item ID.
 *
 * @param {Array} data - An array of data representing the tree.
 * @param {Array} selectedItemUUID - Node ID of selected element in tree.
 * @param {string} view - Selected tab in registry.
 * @returns {Object} - An object containing the selected component, model, and relationship.
 */
export const getFilteredDataForDetailsComponent = (data, selectedItemUUID, view) => {
  let selectedIdArr = selectedItemUUID.split('.');
  const selectedParent = data.find((parentItem) => parentItem.id === selectedIdArr[0]);
  let selectedComponent = [];
  let selectedModel = {};
  let selectedRelationship = [];

  selectedModel = view === REGISTRANTS ? selectedParent?.models : selectedParent;
  const componentKey = view === MODELS ? 'components' : 'relationships';
  selectedComponent = selectedModel?.[componentKey]?.find(
    (component) => component && `${selectedModel.id}.1.${component.id}` === selectedItemUUID,
  );
  selectedRelationship = selectedModel?.relationships?.find(
    (relationship) =>
      relationship && `${selectedModel.id}.2.${relationship.id}` === selectedItemUUID,
  );

  return { selectedComponent, selectedModel, selectedRelationship };
};

export const RelationInitiationType = {
  FROM: "FROM",
  TO: "TO"
};

// Utility function to check if a component can initiate a relation
const canInitiateRelation = (relation, fromComponent) => {
  const fromSelector = relation.selectors?.allow?.from;
  return predicate(fromComponent, fromSelector);
};

// Utility function to check if a component can accept a relation
const canAcceptRelation = (relation, toComponent) => {
  const toSelector = relation.selectors?.allow?.to;
  return predicate(toComponent, toSelector);
};

// Utility function to get all possible components to connect
const getAllPossibleComponentsToConnect = (relation, fromComponent) => {
  if (!canInitiateRelation(relation, fromComponent)) return [];
  return relation.selectors?.allow?.to || [];
};

// Utility function to get all possible components to get connected from
const getAllPossibleComponentsToGetConnectedFrom = (relation, toComponent) => {
  if (!canAcceptRelation(relation, toComponent)) return [];
  return relation.selectors?.allow?.from || [];
};

// Utility function for the common predicate logic
const predicate = (component, selectors) => {
  return !!selectors?.find?.(s => {
    if (s.kind) {
      if (s.kind === "*") {
        //kind wild card matching
        if (s.model === "*" || s.model === component.model.name) {
          return true;
        }
      }
      if (!(s.kind === component.kind)) return false;
    }
    if (s.model) {
      if (s.model === "*") return true; // wildcard matching
      if (!(s.model === component.model.name)) return false;
      if (s.version) {
        if (!(s.version === component.model.version)) return false;
      }
    }
    return true;
  });
};


export const groupPossibleRelationsByKind = (fromComponent, relationships) => {
  const groupedByKind = relationships.reduce((grouped, rel) => {
    grouped[rel.kind] = grouped[rel.kind] || {
      [RelationInitiationType.FROM]: [],
      [RelationInitiationType.TO]: []
    };

    const componentsToConnect = getAllPossibleComponentsToGetConnectedFrom(rel, fromComponent);
    grouped[rel.kind][RelationInitiationType.FROM] = grouped[rel.kind][RelationInitiationType.FROM].concat(
      componentsToConnect.map(from => ({ rel, ...from }))
    );

    const componentsToGetConnectedFrom = getAllPossibleComponentsToConnect(rel, fromComponent);
    grouped[rel.kind][RelationInitiationType.TO] = grouped[rel.kind][RelationInitiationType.TO].concat(
      componentsToGetConnectedFrom.map(to => ({ rel, ...to }))
    );

    return grouped;
  }, {});

  return groupedByKind;
};