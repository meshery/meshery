import * as React from "react";
import TextField from "@mui/material/TextField";
import Dialog from "@mui/material/Dialog";
import DialogTitle from "@mui/material/DialogTitle";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogActions from "@mui/material/DialogActions";
import Button from "@material-ui/core/Button";
import Autocomplete, { createFilterOptions } from "@mui/material/Autocomplete";
import AddIcon from "@mui/icons-material/Add";
import { withStyles } from "@material-ui/core/styles";
import { Grid } from "@mui/material";

const filter = createFilterOptions();

const styles = (theme) => ({
  heading : {
    color : theme.palette.secondary.iconMain,
    "& .MuiInputLabel-root" : {
      color : theme.palette.secondary.iconMain,
    },
    "& .MuiOutlinedInput-root" : {
      color : theme.palette.secondary.iconMain,
    },
    "& .MuiSvgIcon-root" : {
      color : theme.palette.secondary.iconMain,
    }
  },
  wrapper : {
    background : theme.palette.type === 'dark' ? theme.palette.secondary.toolbarBg2 : theme.palette.secondary.toolbarBg1,
    padding : '8px'
  },
  autocompleteOptions : {
    "& .MuiAutocomplete-option" : {
      background : theme.palette.secondary.mainBackground,
    },
  },
  dialogTitle : {
    textAlign : 'center',
    minWidth : 400,
    padding : '10px',
    color : "#fff",
    backgroundColor : theme.palette.type === 'dark' ? theme.palette.secondary.headerColor : theme.palette.secondary.mainBackground,
  },
  subtitle : {
    minWidth : 400,
    overflowWrap : 'anywhere',
    textAlign : 'center',
    padding : '5px'
  },
  button : {
    backgroundColor : theme.palette.secondary.focused,
    "&:hover" : {
      backgroundColor : theme.palette.secondary.focused,
    },
    color : "#fff"
  },
  icon : {
    color : theme.palette.secondary.iconMain,
  },
});

function MySelectComponent({ classes }) {
  const [value, setValue] = React.useState(null);
  const [open, toggleOpen] = React.useState(false);
  const [labelVisible, setLabelVisible] = React.useState(true);

  const handleClose = () => {
    setDialogValue({
      title : "",
      year : "",
    });
    toggleOpen(false);
  };

  const [dialogValue, setDialogValue] = React.useState({
    title : "",
    year : "",
  });

  const handleSubmit = (event) => {
    event.preventDefault();
    setValue({
      title : dialogValue.title,
      year : parseInt(dialogValue.year, 10),
    });
    handleClose();
  };

  React.useEffect(() => {
    if (value) {
      setLabelVisible(false);
    }
  }, [value]);

  return (
    <React.Fragment>
      <Autocomplete
        className={classes.autocompleteOptions}
        value={value}
        onChange={(event, newValue) => {
          if (typeof newValue === "string") {
            // timeout to avoid instant validation of the dialog's form.
            setTimeout(() => {
              toggleOpen(true);
              setDialogValue({
                title : newValue,
              });
            });
          } else if (newValue && newValue.inputValue) {
            toggleOpen(true);
            setDialogValue({
              title : newValue.inputValue,
            });
          } else {
            setValue(newValue);
          }
        }}
        filterOptions={(options, params) => {
          const filtered = filter(options, params);

          if (params.inputValue !== "") {
            filtered.push({
              inputValue : params.inputValue,
              title : `Create:  "${params.inputValue}"`,
            });
          }

          return filtered;
        }}
        id="free-solo-dialog-demo"
        options={environmentOptions}
        getOptionLabel={(option) => {
          // e.g. value selected with enter, right from the input
          if (typeof option === "string") {
            return option;
          }
          if (option.inputValue) {
            return option.inputValue;
          }
          return option.title;
        }}
        selectOnFocus
        clearOnBlur
        handleHomeEndKeys
        renderOption={(props, option) => (
          <div className={classes.wrapper}>
            <li className={classes.wrapper} {...props} style={{ display : "flex", justifyContent : "space-between" }}>
              <div className={classes.heading}>{option.inputValue ? <>{option.title}</> : option.title}</div>
              <div>{option.inputValue && <AddIcon className={classes.icon} />}</div>
            </li>
          </div>
        )}
        sx={{ "& fieldset" : { border : "none" }, width : 250 }}
        renderInput={(params) => (
          <TextField
            className={classes.heading}
            {...params}
            label={labelVisible ? "Choose an environment" : ""}
            placeholder="Environment"
          />
        )}
      />
      <Dialog open={open} onClose={handleClose}>
        <form onSubmit={handleSubmit}>
          <DialogTitle className={classes.dialogTitle}>Add new environment</DialogTitle>
          <DialogContent>
            <DialogContentText className={classes.subtitle}>Create new environment</DialogContentText>
            <Grid container spacing={1} alignItems="center">
              <Grid item xs={3}>
                <h4 className={classes.heading} >Name</h4>
              </Grid>
              <Grid item xs={9}>
                <TextField
                  autoFocus
                  margin="dense"
                  id="name"
                  value={dialogValue.title}
                  onChange={(event) =>
                    setDialogValue({
                      ...dialogValue,
                      title : event.target.value,
                    })
                  }
                  label="Environment name"
                  type="text"
                  variant="outlined"
                  size="small"
                  style={{ width : "100%" }}
                />
              </Grid>
            </Grid>
          </DialogContent>
          <DialogActions>
            <Button size="large" variant="outlined"  onClick={handleClose}>Cancel</Button>
            <Button size="large" variant="contained" color="primary" className={classes.button}>Create</Button>
          </DialogActions>
        </form>
      </Dialog>
    </React.Fragment>
  );
}

const environmentOptions = [
  { title : "Environment 1" },
  { title : "Environment 2" },
  { title : "Environment 3" },
  { title : "Environment 4" },
  { title : "Environment 5" },
];

export default withStyles(styles)(MySelectComponent);
