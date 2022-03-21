import React from "react";
import { NoSsr } from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";
// import MesheryIcon from "../img/meshery-logo/meshery-logo-light-white-text-side.svg"
import { CustomTypography } from "./CustomTypography";
import MesheryIcon from "../img/meshery-logo/CustomMesheryLogo";


const styles = (theme) => ({
    header: {
        width: "fit-content",
        margin: "auto"
    },
    logoClass: {
      marginBottom: "1rem"
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
                    <div className={classes.logoClass} >
                    <MesheryIcon className={classes.header} />
                    </div>
                    <CustomTypography className={classes.description}>
                        Design and operate your cloud native deployments with the extensible management plane, Meshery
                    </CustomTypography>
                </div>
            </NoSsr>
        );
    }
}

export default withStyles(styles)(MesheryLogo);
