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
    width: "12rem",
    minWidth: "10rem",
    border: "1px solid gray",
    borderRadius: "0.75rem",
    //top: "2rem",
    margin: "0rem 0rem 6rem 0rem",
    ["@media (max-width:1024px)"]: {
      //eslint-disable-line no-useless-computed-key
      margin: "0rem 0rem 6rem 0",
    },
  },
  cardChecked: {
    height: "15rem",
    marginBottom: "1rem",
  },
  cardUnchecked: {
    height: "10rem",
  },
  cardContent: {
    display: "flex",
    flexDirection: "column",
    justifyContent: "center",
    alignItems: "center",
    height: "100%",
    width: "100%",
    padding: "0",
  },
  contentTop: {
    background: "#434343",
    height: "12rem",
    width: "100%",
    display: "flex",
    alignItems: "center",
  },
  contentTopUnchecked: {
    background: "#434343",
    height: "100%",
    width: "100%",
    display: "flex",
    alignItems: "center",
  },
  contentTopSwitcher: {
    paddingLeft: "2rem",
  },
  iconContainer: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
    marginLeft: "1rem",
  },
  cardIcon: {
    width: "3rem",
  },
  cardIconText: {
    color: "white",
    fontSize: "0.85rem",
    textAlign: "center",
    "&:first-letter": {
      textTransform: "capitalize",
    },
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
  topInputIcon: {
    position: "absolute",
    fontSize: "1.25rem",
    color: "lightgray",
    bottom: "4.25rem",
    left: "9rem",
    cursor: "pointer",
    zIndex: "99999",
    "&:hover": {
      color: "grey",
    },
  },
  file: {
    display: "none",
  },
});

const ConfigCard = ({ handleSwitch, Icon, KubernetesInput, name, topInputPlaceholder, bottomInputPlaceholder }) => {
  const [isChecked, setIsChecked] = React.useState(false);
  const classes = useStyles();

  const handleChange = (e) => {
    setIsChecked(e.target.checked);
    if (handleSwitch) {
      handleSwitch(e.target.name, e.target.checked);
    }
  };
  return (
    <Card className={`${classes.card} ${classes.cardChecked}`} variant="outlined">
      <CardContent className={classes.cardContent}>
        <div className={classes.contentTop}>
          <div className={classes.iconContainer}>
            <Icon className={classes.cardIcon} alt={`${name} icon`} />
            {name === "openServiceMesh" ? (
              <Typography className={classes.cardIconText} color="primary">
                Open Service <br />
                Mesh
              </Typography>
            ) : (
              <Typography className={classes.cardIconText} color="primary">
                {name}
              </Typography>
            )}
          </div>
          <FormControlLabel
            className={classes.contentTopSwitcher}
            control={<MeshySwitch checked={isChecked} name={name} />}
            onChange={handleChange}
          />
        </div>
        <div className={classes.contentBottomChecked}>
          <>
            {name === "Kubernetes" ? (
              <KubernetesInput/>
            ) : (
              <>
                {" "}
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
          </>
        </div>
      </CardContent>
    </Card>
  );
};

export default ConfigCard;
