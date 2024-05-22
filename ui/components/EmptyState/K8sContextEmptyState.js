import { Button, Link, Typography, makeStyles } from '@material-ui/core';
import OperatorLight from '../../assets/img/OperatorLight';
import Operator from '../../assets/img/Operator';
import theme from '../../themes/app';
import AddIcon from '@material-ui/icons/Add';

const styles = makeStyles((theme) => ({
  textContent: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: '1rem',
    padding: '10px',
    borderRadius: '10px',
  },
  AddIcon: {
    width: theme.spacing(2.5),
    paddingRight: theme.spacing(0.5),
  },
}));

export const K8sEmptyState = ({ message }) => {
  const classes = styles();
  return (
    <div className={classes.textContent}>
      {theme.palette.type == 'dark' ? <OperatorLight /> : <Operator />}
      <Typography variant="h5">{message || 'No cluster connected yet'}</Typography>

      <Link href="/management/connections">
        <Button
          type="submit"
          variant="contained"
          color="primary"
          style={{ margin: '0.6rem 0.6rem', whiteSpace: 'nowrap' }}
        >
          <AddIcon className={classes.AddIcon} />
          Connect Clusters
        </Button>
      </Link>
    </div>
  );
};
