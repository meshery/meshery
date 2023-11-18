import { ToggleButton } from '@layer5/sistent-components';
import GridOnIcon from '@material-ui/icons/GridOn';
import TableChartIcon from '@material-ui/icons/TableChart';
import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles((theme) => ({
  toggleButton: {
    border: 'none',
  },
  icon: {
    color: theme.palette.secondary.iconMain,
  },
}));

function ViewSwitch({ view, changeView }) {
  const classes = useStyles();
  return (
    <ToggleButton
      className={classes.toggleButton}
      size="small"
      value={view}
      onChange={() => {
        changeView(view === 'grid' ? 'table' : 'grid');
      }}
      aria-label="Switch View"
      sx={{
        border: 'none',
      }}
    >
      {view === 'grid' ? (
        <TableChartIcon className={classes.icon} />
      ) : (
        <GridOnIcon className={classes.icon} />
      )}
    </ToggleButton>
  );
}

export default ViewSwitch;
