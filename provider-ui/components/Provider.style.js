import {
  styled,
  MenuItem,
  Dialog,
  DialogActions,
  Typography,
  Popover,
  charcoal,
  accentGrey,
  KEPPEL,
} from "@sistent/sistent";
// Wraps the dropdown button (and its anchored popovers). Equal vertical
// margin above/below makes the chooser sit predictably under the Learn
// more link without depending on its sibling's intrinsic height. The
// LearnMore link is now centered between the logo and this container via
// its own marginTop in the styled component below.
export const CustomDiv = styled("div")(({ theme }) => ({
  width: "100%",
  maxWidth: theme.spacing(52),
  marginLeft: "auto",
  marginRight: "auto",
  marginTop: theme.spacing(4),
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  gap: theme.spacing(3)
}));

export const CustomTypography = styled(Typography)(({ theme }) => ({
  fontWeight: 400,
  fontStyle: "normal",
  fontSize: "1.5rem",
  lineHeight: "2rem",
  letterSpacing: "0.01em",
  color: theme.palette.text.inverse,
}));

export const MesheryLogo = styled("img")(({ theme }) => ({
  width: theme.spacing(50),
  maxWidth: "100%",
  height: "auto",
}));

export const MenuProviderDisabled = styled(MenuItem)(({ theme }) => ({
  display: "flex",
  justifyContent: "space-between",
  backgroundColor: theme.palette.text.default,
  "> span": {
    fontStyle: "italic",
  },
  textOverflow: "ellipsis",
}));

export const CustomDialog = styled(Dialog)(({ theme }) => ({
  // overflow:hidden on the Paper clips the inner DialogContent /
  // DialogActions backgrounds against the dialog's border-radius, so the
  // rounded corners read as clean curves instead of showing the inner
  // surface bleeding past the radius (the "corner artifact" the chooser
  // showed in dark mode).
  "& .MuiDialog-paper": {
    borderRadius: theme.spacing(1),
    overflow: "hidden",
  },
  "& .MuiDialogContent-root": {
    padding: theme.spacing(2),
    background: theme.palette.background.elevatedComponents,
  },
  "& .MuiDialogActions-root": {
    padding: theme.spacing(1),
  },
  "& .MuiDialogContentText-root > a": {
    color: "#222",
  },
}));
export const CustomDialogActions = styled(DialogActions)(({ theme }) => ({
  display: "flex",
  justifyContent: "space-between",
  background: charcoal[20],
  "& div": {
    margin: ".8rem",
  },
  "& div > a": {
    color: theme.palette.text.inverse,
    paddingLeft: "1rem",
  },
}));
export const StyledPopover = styled(Popover)(({ theme }) => ({
  // borderRadius + overflow:hidden clip MenuItem hover backgrounds against
  // the rounded corners. The Paper background must match the MenuList's
  // own background (charcoal[20]) — when the two differ, the rectangular
  // MenuList shows through past the rounded Paper corners as four square
  // notches at each corner. Setting Paper to the same color as the
  // MenuList collapses the two surfaces into one continuous rounded fill.
  "& .MuiPaper-root": {
    backgroundColor: charcoal[20],
    color: theme.palette.text.inverse,
    borderRadius: theme.spacing(1),
    overflow: "hidden",
  },
  "& .MuiList-root": {
    backgroundColor: "transparent",
    borderRadius: "inherit",
    overflow: "hidden",
    padding: 0,
  },
}));
// Per-provider info popover. Anchored to the InfoOutlined icon next to each
// remote entry in the chooser. Theme tokens drive the surface so the card
// stays on-brand across Sistent palettes; the KEPPEL accent border matches
// the InfoOutlined icon's color in the menu, tying the trigger to its panel.
export const ProviderInfoPopover = styled(Popover)(({ theme }) => ({
  "& .MuiPaper-root": {
    background: charcoal[20],
    // Force inheritance for every nested Typography/Link so the dark
    // surface is not undone by MUI defaults (Typography defaults to
    // text.primary which is dark in light themes; Link to primary.main).
    color: theme.palette.text.inverse,
    maxWidth: 420,
    minWidth: 320,
    padding: theme.spacing(2),
    borderLeft: `3px solid ${KEPPEL}`,
    borderRadius: theme.spacing(1),
    overflow: "hidden",
    boxShadow: theme.shadows[6],
    "& .MuiTypography-root, & a, & p, & h1, & h2, & h3, & h4, & h5, & h6, & li":
      {
        color: "inherit",
      },
  },
}));

// Visual divider between sections inside the info popover. Lighter than
// Sistent's default Divider so it reads as a soft separator on the dark
// charcoal[20] surface without competing with the content.
export const ProviderInfoSectionRule = styled("div")(({ theme }) => ({
  height: 1,
  background: accentGrey[40],
  opacity: 0.4,
  margin: theme.spacing(1.25, 0),
}));

// Section heading used inside the info popover ("Capabilities", "What this
// provider offers"). Renders as a small, slightly muted label so it groups
// the section without pulling focus from the content beneath.
export const ProviderInfoSectionHeading = styled(Typography)(({ theme }) => ({
  fontSize: "0.7rem",
  fontWeight: 700,
  letterSpacing: "0.06em",
  textTransform: "uppercase",
  color: theme.palette.text.inverse,
  opacity: 0.7,
  marginBottom: theme.spacing(0.75),
}));

// Bulleted description list inside the info popover. Uses inherited
// Typography color so theme switches recolor automatically; only structural
// styles live here.
export const ProviderInfoDescriptionList = styled("ul")(({ theme }) => ({
  paddingLeft: theme.spacing(2.5),
  margin: 0,
  "& li": {
    marginBottom: theme.spacing(0.5),
    fontSize: "0.85rem",
    lineHeight: 1.45,
  },
  "& li:last-child": {
    marginBottom: 0,
  },
}));

// LearnMore is now a standalone block between the Meshery logo and the
// chooser dropdown. Equal vertical margin top/bottom places it visually
// centered in the gap, regardless of how tall the logo (image) renders at
// the current viewport width.
export const LearnMore = styled("a")(({ theme }) => ({
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  marginTop: theme.spacing(4),
  marginBottom: theme.spacing(4),
  color: theme.palette.text.inverse,
  textDecoration: "none",
  fontSize: "1rem",
  lineHeight: 1.4,
  fontWeight: 500,
  cursor: "pointer",
  transition: "color 120ms ease, opacity 120ms ease",
  "&:hover": {
    textDecoration: "underline",
    opacity: 0.92,
  },
  "&:focus-visible": {
    outline: `2px solid ${theme.palette.primary.main}`,
    outlineOffset: "4px",
    borderRadius: "4px",
  },
}));
