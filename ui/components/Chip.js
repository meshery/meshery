import { Chip } from "@mui/material";
import { styled } from "@mui/material/styles";

export default styled(Chip)(({ theme }) => ({
  padding: theme.spacing(1),
  cursor: "pointer",
}));
