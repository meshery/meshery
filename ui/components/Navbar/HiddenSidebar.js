import { styled } from "@mui/material/styles";

export const HiddenscrollbarStyle = styled("div")({
  overflow: "hidden auto",
  "scrollbar-width": "none",
  "-ms-overflow-style": "none",
  "&::-webkit-scrollbar": {
    display: "none",
  },
});
