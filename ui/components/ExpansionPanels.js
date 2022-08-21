import { Accordion, AccordionSummary } from "@mui/material";
import { withStyles } from "@mui/styles";

export const ExpansionPanel = withStyles({
  root: { border: "1px solid rgba(0,0,0,.125)" },
  expanded: { margin: "auto" },
})(Accordion);

export const ExpansionPanelSummary = withStyles({
  root: { borderBottom: "1px solid rgba(0,0,0,.125)" },
  content: { "&$expanded": { margin: "12px 0" } },
  expanded: {},
})((props) => <AccordionSummary {...props} />);
