import { styled } from "@mui/material/styles";
import { Box } from "@mui/material";

export const ChartDiv = styled(Box)(() => ({
  padding: "1rem",
  borderRadius: "1rem",
  width: "50%",
  height: "20rem",
  marginInline: "1rem",
  marginBottom: "1rem",
  display: "block",
  ["@media (max-width:900px)"]: {
    height: "18rem",
    marginInline: "0",
    padding: "0.5rem"
  }
}));