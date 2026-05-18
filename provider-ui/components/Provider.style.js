import {
  styled,
  MenuItem,
  Dialog,
  DialogActions,
  Typography,
  Popover,
  charcoal,
} from "@sistent/sistent";
export const CustomDiv = styled("div")(({ theme }) => ({
  width: "100%",
  maxWidth: theme.spacing(52),
  marginLeft: "auto",
  marginRight: "auto",
  marginTop: theme.spacing(4),
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  gap: theme.spacing(3),
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
  ".MuiPaper-root": {
    backgroundColor: theme.palette.background.brand.disabled,
    color: theme.palette.text.default,
  },
}));
export const LearnMore = styled("a")(({ theme }) => ({
  display: "inline-flex",
  alignItems: "center",
  justifyContent: "center",
  color: theme.palette.text.default,
  textDecoration: "none",
  fontSize: "1rem",
  lineHeight: 1.4,
  fontWeight: 500,
  cursor: "pointer",
  transition: "color 120ms ease, opacity 120ms ease",
  "&:hover": {
    color: theme.palette.text.inverse,
    opacity: 0.92,
  },
  "&:focus-visible": {
    outline: `2px solid ${theme.palette.primary.main}`,
    outlineOffset: "4px",
    borderRadius: "4px",
  },
}));
