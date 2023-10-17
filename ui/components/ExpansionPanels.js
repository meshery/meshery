import MuiExpansionPanel from '@material-ui/core/ExpansionPanel';
import MuiExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import { withStyles } from '@material-ui/core';

export const ExpansionPanel = withStyles({
  root: { border: '1px solid rgba(0,0,0,.125)' },
  expanded: { margin: 'auto' },
})(MuiExpansionPanel);

export const ExpansionPanelSummary = withStyles({
  root: { borderBottom: '1px solid rgba(0,0,0,.125)' },
  content: { '&$expanded': { margin: '12px 0' } },
  expanded: {},
})((props) => <MuiExpansionPanelSummary {...props} />);
