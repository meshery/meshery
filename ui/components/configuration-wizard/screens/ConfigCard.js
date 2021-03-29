import React from "react";
import {
  makeStyles,
  withStyles,
  Switch,
  FormControlLabel,
  Input,
  Card,
  CardContent,
  Typography,
} from "@material-ui/core/";
import TimerIcon from "@material-ui/icons/Timer";
import FiberManualRecordRoundedIcon from "@material-ui/icons/FiberManualRecordRounded";

const MeshySwitch = withStyles({
  switchBase: {
    color: "grey",
    "&$checked": {
      color: "#00B39F",
    },
    "&$checked + $track": {
      backgroundColor: "#00B39F",
    },
  },
  checked: {},
  track: {},
})(Switch);

const useStyles = makeStyles({
  card: {
    position: "relative",
    width: "10rem",
    minWidth: "10rem",
    border: "1px solid gray",
    borderRadius: "0.75rem",
    top: "2rem",
    padding: 0,
    margin: "0rem 2rem 5rem 2rem",
    ["@media (max-width:1024px)"]: { //eslint-disable-line no-useless-computed-key
      margin: "0rem 2rem 5rem 0",
    },
  },
  cardChecked: {
    height: "15rem",
    padding: 0,
    marginBottom: "0rem",
  },
  cardUnchecked: {
    height: "10rem",
    padding: 0,
  },
  cardContent: {
    background: "red",
    display: "flex",
    flexDirection: "column",
    justifyContent: "center",
    alignItems: "center",
    height: "100%",
    width: "100%",
    padding: 0,
  //  margin: "-1rem 0 0 -1rem",
  },
  contentTop: {
    background: "#434343",
    height: "12rem",
    width: "100%",
    display: "flex",
    alignItems: "center",
    padding: 0,
  },
  contentTopUnchecked: {
    background: "#434343",
    height: "100%",
    width: "100%",
    display: "flex",
    alignItems: "center",
    padding: 0,
  },
  contentTopSwitcher: {
    marginLeft: "0.5rem",
  },
  iconContainer: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
  },
  cardIcon: {
    width: "3rem",
  },
  cardIconText: {
    color: "white",
    fontSize: "0.85rem",
  },
  contentBottomChecked: {
    background: "white",
    height: "6rem",
    width: "100%",
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
  },
  contentBottomUnchecked: {
    display: "none",
  },
  contentBottomInput: {
    border: "1px solid lightgray",
    borderRadius: "5px",
    width: "9rem",
    height: "2rem",
    marginBottom: "0.15rem",
    fontSize: "0.75rem",
    padding: "0.50rem",
  },
  contentBottomIcon: {
    marginBottom: "-0.4rem",
    color: "#00B39F",
  },
});

const ConfigCard = ({
  Icon,
  name,
  topInputPlaceholder,
  bottomInputPlaceholder,
}) => {
  const [state, setState] = React.useState({
    checked: false,
  });
  const classes = useStyles();

  const handleChange = (event) => {
    setState({ ...state, [event.target.name]: event.target.checked });
  };
  return (
    <Card
      className={
        state.checked
          ? `${classes.card} ${classes.cardChecked}`
          : `${classes.card} ${classes.cardUnchecked}`
      }
      variant="outlined"
    >
      <CardContent className={classes.cardContent} style={{padding: 0}}>
        <div
          className={
            state.checked ? classes.contentTop : classes.contentTopUnchecked
          }
        >
          <FormControlLabel
            className={classes.contentTopSwitcher}
            control={<MeshySwitch checked={state.checked} name="checked" />}
            onChange={handleChange}
          />
          <div className={classes.iconContainer}>
            {Icon === "timer" ? (
              <TimerIcon />
            ) : (
              <Icon className={classes.cardIcon}
                alt={`${name} icon`}/>           
            )}
            <Typography className={classes.cardIconText} color="primary">
              {name}
            </Typography>
          </div>
        </div>
        <div
          className={
            state.checked
              ? classes.contentBottomChecked
              : classes.contentBottomUnchecked
          }
        >
          {name === "Open Service Mesh" ||
          name === "Consul" ||
          name === "Linkerd" ? (
              <>
                <Typography className={classes.contentBottomControlPlane}>
                Control Plane: 6{" "}
                  <FiberManualRecordRoundedIcon
                    className={classes.contentBottomIcon}
                  />
                </Typography>
                <Typography className={classes.contentBottomDataPlane}>
                Data Plane: 18
                  <FiberManualRecordRoundedIcon
                    className={classes.contentBottomIcon}
                  />
                </Typography>
              </>
            ) : (
              <>
                <Input
                  placeholder={topInputPlaceholder}
                  disableUnderline="false"
                  className={classes.contentBottomInput}
                ></Input>
                <Input
                  placeholder={bottomInputPlaceholder}
                  disableUnderline="false"
                  className={classes.contentBottomInput}
                ></Input>
              </>
            )}
        </div>
      </CardContent>
    </Card>
  );
};

export default ConfigCard;

{ /*<img
                className={classes.cardIcon}
                src={icon}
                alt={`${name} icon`}
              /> */ }