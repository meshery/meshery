/**
 * Utilities for analyzing the pre-change impact of modified or removed components
 * and relationships within a Meshery Design.
 */

export interface ImpactedComponent {
  id: string;
  name: string;
  kind: string;
  model: string;
  explanation: string;
}

export interface ImpactedRelationship {
  id: string;
  kind?: string;
  type?: string;
  subType?: string;
  selectors?: any[];
  [key: string]: any;
}

export interface DesignImpactAnalysisResult {
  targetType: 'component' | 'relationship' | 'unknown';
  targetId: string;
  targetName: string;
  targetKind?: string;
  modificationType: 'update' | 'delete' | 'create';
  directlyAffected: ImpactedComponent[];
  indirectlyAffected: ImpactedComponent[];
  affectedRelationships: ImpactedRelationship[];
  totalAffectedCount: number;
  summary: string;
}

export interface AnalyzeDesignImpactOptions {
  design?: any;
  target?: any;
  targetType?: 'component' | 'relationship';
  modificationType?: 'update' | 'delete' | 'create';
}

function extractSelectorItemIds(items: any[] = []): string[] {
  const result: string[] = [];
  if (!Array.isArray(items)) return result;

  for (const item of items) {
    if (!item) continue;
    if (typeof item === 'string') {
      result.push(item);
    } else {
      if (item.id) result.push(item.id);
      if (item.name) result.push(item.name);
    }
  }
  return result;
}

/**
 * Analyzes the pre-change impact on components and relationships when a component
 * or relationship is modified or removed in a Meshery Design.
 */
export function analyzeDesignImpact({
  design,
  target,
  targetType,
  modificationType = 'update',
}: AnalyzeDesignImpactOptions): DesignImpactAnalysisResult {
  if (!design || !target) {
    return {
      targetType: 'unknown',
      targetId: '',
      targetName: '',
      modificationType,
      directlyAffected: [],
      indirectlyAffected: [],
      affectedRelationships: [],
      totalAffectedCount: 0,
      summary: 'No target or design provided for impact analysis.',
    };
  }

  const components = Array.isArray(design.components) ? design.components : [];
  const relationships = Array.isArray(design.relationships) ? design.relationships : [];

  const compMap = new Map<string, any>();
  for (const c of components) {
    if (!c) continue;
    if (c.id) compMap.set(c.id, c);
    if (c.name) compMap.set(c.name, c);
  }

  const relMap = new Map<string, any>();
  for (const r of relationships) {
    if (!r) continue;
    if (r.id) relMap.set(r.id, r);
  }

  const isRelExplicit = targetType === 'relationship';
  const isCompExplicit = targetType === 'component';

  const isTargetObject = typeof target === 'object' && target !== null;
  const targetId = typeof target === 'string' ? target : target.id || target.componentId || '';

  const isRel =
    isRelExplicit ||
    (!isCompExplicit &&
      (relMap.has(targetId) ||
        (isTargetObject &&
          (target.selectors !== undefined ||
            target.subType !== undefined ||
            target.relationshipType !== undefined))));

  if (!isRel) {
    // Target is a component
    const targetComp = compMap.get(targetId) || (isTargetObject ? target : null);
    const targetName = targetComp?.name || targetComp?.displayName || targetId;
    const targetKind =
      targetComp?.component?.kind || targetComp?.kind || (targetComp ? 'Component' : 'Resource');

    const directlyAffectedIds = new Set<string>();
    const affectedRelMap = new Map<string, ImpactedRelationship>();
    const explanationMap = new Map<string, string>();

    for (const rel of relationships) {
      if (!rel) continue;
      const selectors = Array.isArray(rel.selectors) ? rel.selectors : [];
      let relInvolvesTarget = false;
      const matchedFrom: string[] = [];
      const matchedTo: string[] = [];

      for (const sel of selectors) {
        if (!sel) continue;
        const fromList = extractSelectorItemIds(sel.allow?.from || sel.from);
        const toList = extractSelectorItemIds(sel.allow?.to || sel.to);

        const fromMatches = fromList.includes(targetId) || fromList.includes(targetName);
        const toMatches = toList.includes(targetId) || toList.includes(targetName);

        if (fromMatches || toMatches) {
          relInvolvesTarget = true;
          for (const f of fromList) {
            if (f !== targetId && f !== targetName) matchedFrom.push(f);
          }
          for (const t of toList) {
            if (t !== targetId && t !== targetName) matchedTo.push(t);
          }
        }
      }

      if (relInvolvesTarget) {
        const relId =
          rel.id || `${rel.kind || 'Relationship'}-${rel.type || ''}-${rel.subType || ''}`;
        affectedRelMap.set(relId, rel);

        const relDesc = `${rel.kind || 'Edge'} (${rel.type || rel.subType || 'connection'})`;

        for (const otherId of [...matchedFrom, ...matchedTo]) {
          if (otherId && otherId !== targetId && otherId !== targetName) {
            directlyAffectedIds.add(otherId);
            explanationMap.set(
              otherId,
              `Directly connected via ${relDesc} with ${targetKind}/${targetName}`,
            );
          }
        }
      }
    }

    // Traverse multi-level relationship chains (indirect dependencies)
    const indirectlyAffectedIds = new Set<string>();
    const queue = Array.from(directlyAffectedIds);
    const visited = new Set<string>([targetId, targetName, ...directlyAffectedIds]);

    while (queue.length > 0) {
      const currentId = queue.shift()!;
      const currentComp = compMap.get(currentId);
      const currentName = currentComp?.name;
      const currentKind = currentComp?.component?.kind || currentComp?.kind || 'Component';

      for (const rel of relationships) {
        if (!rel) continue;
        const selectors = Array.isArray(rel.selectors) ? rel.selectors : [];
        for (const sel of selectors) {
          if (!sel) continue;
          const fromList = extractSelectorItemIds(sel.allow?.from || sel.from);
          const toList = extractSelectorItemIds(sel.allow?.to || sel.to);

          const matches =
            fromList.includes(currentId) ||
            (currentName && fromList.includes(currentName)) ||
            toList.includes(currentId) ||
            (currentName && toList.includes(currentName));

          if (matches) {
            for (const neighbor of [...fromList, ...toList]) {
              if (neighbor && !visited.has(neighbor)) {
                visited.add(neighbor);
                indirectlyAffectedIds.add(neighbor);
                queue.push(neighbor);
                explanationMap.set(
                  neighbor,
                  `Transitively dependent through ${currentKind}/${currentName || currentId}`,
                );
              }
            }
          }
        }
      }
    }

    const formatComponent = (id: string): ImpactedComponent => {
      const comp = compMap.get(id);
      return {
        id,
        name: comp?.name || id,
        kind: comp?.component?.kind || comp?.kind || 'Component',
        model: comp?.model?.name || comp?.model || 'kubernetes',
        explanation: explanationMap.get(id) || 'Connected relationship dependency',
      };
    };

    const directList = Array.from(directlyAffectedIds).map(formatComponent);
    const indirectList = Array.from(indirectlyAffectedIds).map(formatComponent);
    const relList = Array.from(affectedRelMap.values());
    const totalCount = directList.length + indirectList.length;

    const actionVerb = modificationType === 'delete' ? 'Removing' : 'Modifying';
    const summary =
      totalCount === 0
        ? `${actionVerb} ${targetKind}/${targetName}: no dependent components or relationships found.`
        : `${actionVerb} ${targetKind}/${targetName} potentially affects ${directList.length} directly connected and ${indirectList.length} indirect component(s).`;

    return {
      targetType: 'component',
      targetId,
      targetName,
      targetKind,
      modificationType,
      directlyAffected: directList,
      indirectlyAffected: indirectList,
      affectedRelationships: relList,
      totalAffectedCount: totalCount,
      summary,
    };
  }

  // Target is a relationship
  const rel = relMap.get(targetId) || (isTargetObject ? target : null);
  const affectedCompIds = new Set<string>();

  if (rel) {
    const selectors = Array.isArray(rel.selectors) ? rel.selectors : [];
    for (const sel of selectors) {
      if (!sel) continue;
      const fromList = extractSelectorItemIds(sel.allow?.from || sel.from);
      const toList = extractSelectorItemIds(sel.allow?.to || sel.to);
      for (const id of [...fromList, ...toList]) {
        if (id) affectedCompIds.add(id);
      }
    }
  }

  const directList: ImpactedComponent[] = Array.from(affectedCompIds).map((id) => {
    const comp = compMap.get(id);
    return {
      id,
      name: comp?.name || id,
      kind: comp?.component?.kind || comp?.kind || 'Component',
      model: comp?.model?.name || comp?.model || 'kubernetes',
      explanation: `Endpoint of relationship ${rel?.kind || 'Relationship'}/${rel?.type || rel?.subType || targetId}`,
    };
  });

  const actionVerb = modificationType === 'delete' ? 'Removing' : 'Modifying';
  const relDisplayName = rel?.subType || rel?.type || rel?.kind || targetId;
  const summary =
    directList.length === 0
      ? `Relationship ${targetId} has no attached components.`
      : `${actionVerb} relationship "${relDisplayName}" affects ${directList.length} endpoint component(s).`;

  return {
    targetType: 'relationship',
    targetId,
    targetName: relDisplayName,
    modificationType,
    directlyAffected: directList,
    indirectlyAffected: [],
    affectedRelationships: rel ? [rel] : [],
    totalAffectedCount: directList.length,
    summary,
  };
}
