/* eslint-disable no-unused-vars */
import {
  withStyles,
  Typography,
  Grid,
  Chip,
  IconButton,
} from "@material-ui/core/";
import { withSnackbar } from "notistack";

const chipStyles = (theme) => ({ chipIcon : { width : theme.spacing(2.5) },
  chip : { marginRight : theme.spacing(1),
    marginBottom : theme.spacing(1), }, })

const AdapterChip = withStyles(chipStyles)(({
  classes, handleClick, handleDelete, label, image, isActive
}) => (

  <Chip
    label={label}
    onDelete={handleDelete}
    onClick={handleClick}
    icon={<img src={image} className={classes.chipIcon} />}
    className={classes.chip}
    key={label+"-key"}
    variant={isActive
      ? "outlined"
      : 'default'}
  />
))


export default AdapterChip
