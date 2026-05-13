/**
 * Canonical icon barrel for Meshery UI.
 *
 * Single home for typed SVG icons and the centralized re-export of
 * `@mui/icons-material` glyphs that have no in-house typed equivalent yet.
 * Phase 4 (#18744) consolidated three previously competing locations
 * (`ui/assets/icons/`, `ui/assets/new-icons/`, and `ui/components/icons/`)
 * into this directory; the MUI re-exports below preserve the centralization
 * established by Phase 2 #18736 so only one file changes when Sistent ships
 * a typed icon set.
 *
 * Conventions:
 *  - Each typed SVG icon is a React functional component accepting
 *    `IconProps` (standard SVG props plus optional `width`/`height`/`fill`
 *    overrides). `fill` defaults to `currentColor` so icons inherit color
 *    from CSS.
 *  - Each typed icon module also provides a `default` export for ergonomic
 *    `import EditIcon from '...'` consumption.
 *  - MUI glyphs are re-exported under their original `@mui/icons-material`
 *    names; prefer the typed SVG variants (`*Icon`) when both exist.
 */

export type { IconProps } from './types';

// ---------------------------------------------------------------------------
// Typed SVG icons (Phase 2 #18730).
// ---------------------------------------------------------------------------

export { ArrowDropDownIcon } from './ArrowDropDownIcon';
export { ChevronLeftIcon } from './ChevronLeftIcon';
export { ChevronRightIcon } from './ChevronRightIcon';
export { EditIcon } from './EditIcon';
export { FullscreenExitIcon } from './FullscreenExitIcon';
export { FullscreenIcon } from './FullscreenIcon';
export { GetAppIcon } from './GetAppIcon';
export { LockIcon } from './LockIcon';
export { SaveIcon } from './SaveIcon';
export { SettingsIcon } from './SettingsIcon';

// Pre-existing typed icons in `ui/assets/icons/` whose glyphs also appear
// >=3 times via `@mui/icons-material`. Re-exported here so that downstream
// migrations can import the entire canonical set from a single path.
//
// Note: these are intentionally re-exported as `default` since the source
// files only provide default exports today. Renaming their source exports
// is out of scope for this purely-additive issue.
export { default as DeleteIcon } from './DeleteIcon';
export { default as ExpandMoreIcon } from './ExpandMoreIcon';
export { default as InfoOutlinedIcon } from './InfoOutlined';

// ---------------------------------------------------------------------------
// Centralized MUI icon re-exports (Phase 2 #18736, relocated here in
// Phase 4 #18744 from the now-deleted `ui/components/icons/` barrel).
//
// All `@mui/icons-material` imports should go through this barrel. When
// Sistent provides a typed equivalent, only this file needs to change.
// ---------------------------------------------------------------------------

export {
  AccessTime,
  AddCircle,
  AddCircleOutlined,
  ArrowBack,
  ArrowDropDown,
  BarChart,
  BuildRounded,
  Cached,
  CheckCircle,
  Close,
  Code,
  Delete,
  DeleteForever,
  DirectionsCar,
  DoneAll,
  Edit,
  Error,
  ExpandLess,
  ExpandMore,
  Explore,
  FileCopy,
  FileUpload,
  Filter,
  FilterAlt,
  Fullscreen,
  FullscreenExit,
  GetApp,
  GroupAdd,
  Handyman,
  HelpOutlineOutlined,
  HelpOutlined,
  Info,
  InfoOutlined,
  InsertDriveFile,
  LanOutlined,
  Launch,
  Link,
  ListAlt,
  Lock,
  MergeOutlined,
  MoreHoriz,
  MoreVert,
  NotInterestedRounded,
  OpenInNewOutlined,
  Public,
  Publish,
  RemoveCircle,
  Reply,
  Save,
  SaveAlt,
  SaveAs,
  SaveOutlined,
  Search,
  Settings,
  SimCard,
  SupervisedUserCircle,
  Sync,
  TouchApp,
  Warning,
} from '@mui/icons-material';
