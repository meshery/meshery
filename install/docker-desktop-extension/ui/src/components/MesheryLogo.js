import React from "react";
import { NoSsr, Typography } from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";
import logo from "../img/meshery-logo/meshery-logo.svg";
import meshery from "../img/meshery-logo/meshery-text-dark.png";

const styles = (theme) => ({
    header: {
        width: "fit-content",
        margin: "auto"
    },
    logoClass: {
      width: theme.spacing(10),
      maxWidth: "100%",
      marginBottom: "50px",
      height: "auto",
    },
    mesheryClass: {
      width: "auto",
      maxWidth: "100%",
      height: theme.spacing(25),
    },
    description: {
        maxWidth: "60%",
        margin: "auto",
    }
});

class MesheryLogo extends React.Component {

    render() {
        const { classes } = this.props;

        return (
            <NoSsr>
                <div className={classes.header}>
                    <img src={logo} className={classes.logoClass} alt="logo"/>
                    <img src={meshery} className={classes.mesheryClass} alt="meshery"/>
                    <Typography className={classes.description}>
                        Design and operate your cloud native deployments with the extensible management plane, Meshery
                    </Typography>
                </div>
            </NoSsr>
        );
    }
}

export default withStyles(styles)(MesheryLogo);
