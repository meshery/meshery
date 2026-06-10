/**
 * Canonical icon barrel for Meshery UI.
 *
 * Single home for typed SVG icons and centralized design-system icon re-exports.
 * Phase 4 (#18744) consolidated three previously competing locations
 * (`ui/assets/icons/`, `ui/assets/new-icons/`, and `ui/components/icons/`)
 * into this directory so downstream callers can import icons from one place.
 *
 * Conventions:
 *  - Each typed SVG icon is a React functional component accepting
 *    `IconProps` (standard SVG props plus optional `width`/`height`/`fill`
 *    overrides). `fill` defaults to `currentColor` so icons inherit color
 *    from CSS.
 *  - Each typed icon module also provides a `default` export for ergonomic
 *    `import EditIcon from '...'` consumption.
 *  - Design-system glyphs are re-exported under their legacy names for
 *    compatibility; prefer the typed SVG variants (`*Icon`) when both exist.
 */

export type { IconProps } from './types';

// ---------------------------------------------------------------------------
// Typed SVG icons (Phase 2 #18730).
// ---------------------------------------------------------------------------

export { ArrowDropDownIcon } from './ArrowDropDownIcon';
export { EditIcon } from './EditIcon';
export { FullscreenExitIcon } from './FullscreenExitIcon';
export { FullscreenIcon } from './FullscreenIcon';

// Pre-existing typed icons in `ui/assets/icons/` whose glyphs also appear
// >=3 times via `@mui/icons-material`. Re-exported here so that downstream
// migrations can import the entire canonical set from a single path.
//
// Note: these are intentionally re-exported as `default` since the source
// files only provide default exports today. Renaming their source exports
// is out of scope for this purely-additive issue.

export { default as ExpandMoreIcon } from './ExpandMoreIcon';

// ---------------------------------------------------------------------------
// Centralized design-system icon re-exports (Phase 2 #18736, relocated here in
// Phase 4 #18744 from the now-deleted `ui/components/icons/` barrel).
//
// Keep the legacy export names stable for downstream imports while sourcing the
// actual glyphs from Sistent or local typed SVG components.
// ---------------------------------------------------------------------------

export {
  AccessTimeIcon,
  AddCircleIcon as AddCircle,
  ArrowBackIcon as ArrowBack,
  BarChartIcon as BarChart,
  BuildRoundedIcon as BuildRounded,
  CachedIcon as Cached,
  CheckCircleIcon as CheckCircle,
  CloseIcon as Close,
  CodeIcon as Code,
  DeleteIcon,
  DeleteIcon as Delete,
  DeleteForever,
  DirectionsCarIcon as DirectionsCar,
  DoneAllIcon as DoneAll,
  EditIcon as Edit,
  ErrorIcon as Error,
  ExpandLessIcon as ExpandLess,
  ExpandMoreIcon as ExpandMore,
  ExploreIcon as Explore,
  FileCopyIcon as FileCopy,
  FileUploadIcon as FileUpload,
  FilterIcon as Filter,
  GetAppIcon,
  GetAppIcon as GetApp,
  GroupAddIcon,
  HandymanIcon as Handyman,
  HelpOutlinedIcon as HelpOutlined,
  InfoIcon as Info,
  InfoOutlinedIcon,
  InsertDriveFileIcon as InsertDriveFile,
  LaunchIcon as Launch,
  Link,
  ListAltIcon as ListAlt,
  LockIcon,
  LockIcon as Lock,
  MergeOutlinedIcon as MergeOutlined,
  MoreHorizIcon as MoreHoriz,
  MoreVertIcon as MoreVert,
  OpenInNewIcon as OpenInNewOutlined,
  PublicIcon as Public,
  PublishIcon as Publish,
  RemoveIcon as RemoveCircle,
  ReplyIcon as Reply,
  SaveAsIcon as SaveAs,
  SaveIcon,
  SaveIcon as Save,
  SaveIcon as SaveOutlined,
  SearchIcon as Search,
  SettingsIcon,
  SettingsIcon as Settings,
  SimCardIcon as SimCard,
  SupervisedUserCircleIcon as SupervisedUserCircle,
  SyncIcon as Sync,
  TouchAppIcon as TouchApp,
  WarningIcon as Warning,
  DownloadIcon as SaveAlt,
  CancelIcon as NotInterestedRounded,
  ChevronRightIcon,
  ChevronLeftIcon
} from '@sistent/sistent';

export { default as AddCircleOutlined } from './AddIconCircleBorder';
export { ArrowDropDownIcon as ArrowDropDown } from './ArrowDropDownIcon';
export { default as FilterAlt } from './ContentFilterIcon';
export { default as Fullscreen } from './FullscreenIcon';
export { default as FullscreenExit } from './FullscreenExitIcon';
export { default as HelpOutlineOutlined } from './HelpOutlineIcon';
export { default as LanOutlined } from './Connection';
