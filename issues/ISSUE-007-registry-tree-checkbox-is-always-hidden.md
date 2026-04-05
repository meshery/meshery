### Current Behavior
The bulk-selection checkbox rendered by `StyledTreeItem` in the registry tree is always hidden, even when the row is hovered or already checked.

The visibility expression returns `'hidden'` in both branches, so the control never becomes visible.

### Expected Behavior
When checkbox support is enabled for a tree item, the checkbox should become visible on hover or while checked.

Users should not have a permanently hidden interactive control inside the registry tree.

### Screenshots/Logs
Relevant code path:

```tsx
// ui/components/Settings/Registry/StyledTreeItem.tsx
<Checkbox
  onClick={(e) => {
    e.stopPropagation();
    setChecked((prevcheck) => !prevcheck);
  }}
  size="small"
  checked={checked}
  style={{
    visibility: hover || checked ? 'hidden' : 'hidden',
  }}
/>
```

This component is used by multiple registry tree views, including:

```tsx
// ui/components/Settings/Registry/MesheryTreeViewItem.tsx
<StyledTreeItem
  ...
  check={true}
>
```

### Environment

- **Host OS:** Mac Linux Windows
- **Platform:** Docker or Kubernetes
- **Meshery Server Version:** stable-v
- **Meshery Client Version:** stable-v

<!-- Optional 
### To Reproduce
1. Open the Mesh Model / Registry tree UI.
2. Navigate to a node rendered with `check={true}`.
3. Hover the tree row or toggle selection.
4. See that the checkbox never becomes visible.
-->

---
### Contributor [Guides](https://docs.meshery.io/project/contributing) and [Handbook](https://meshery.io/community#handbook)
- 🛠 [Meshery Build & Release Strategy](https://docs.meshery.io/project/contributing/build-and-release)
- 📚 [Instructions for contributing to documentation](https://docs.meshery.io/project/contributing/contributing-docs)
   - Meshery documentation [site](https://docs.meshery.io/) and [source](https://github.com/meshery/meshery/tree/master/docs)
- 🎨 Wireframes and [designs for Meshery UI](https://www.figma.com/file/SMP3zxOjZztdOLtgN4dS2W/Meshery-UI) in Figma [(open invite)](https://www.figma.com/team_invite/redeem/GvB8SudhEOoq3JOvoLaoMs)
- 📺 [Self-paced Contributor Trainings](https://meshery.io/talks-and-trainings#trainings)
- 🙋🏾🙋🏼 Questions: [Discussion Forum](https://meshery.io/community#community-forums) and [Community Slack](https://slack.meshery.io)
