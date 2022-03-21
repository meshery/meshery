import React from "react";
import { 
    NoSsr, 
    Container, 
    Typography, 
    Button
} from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";
import Location from "./Location";
import ServiceMesh from "./ServiceMesh";
import { CustomTypography } from "../CustomTypography";

const styles = (theme) => ({
    main: {
        padding: theme.spacing(3),
        color: "#AAAAAA"
    },
    config: {
        margin: theme.spacing(5),
        backgroundColor: "#393F49",
       borderRadius: "20px 20px 0 0"
    },
    mesheryConfig: {
        margin: theme.spacing(6)
    }
});

class Configuration extends React.Component {
    render() {

        const { classes } = this.props;

        return (
            <NoSsr>
                <div className={classes.config}>
                    <CustomTypography variant="h6" className={classes.main}>
                        CONFIGURE YOUR MESHERY DEPLOYMENT
                    </CustomTypography>

                    <Location />
                    <ServiceMesh />
{/* 
                    <Container>
                        <Button className={classes.mesheryConfig} color="primary" variant="contained">DEPLOY MESHERY</Button>
                        <Button className={classes.mesheryConfig} variant="contained">OPEN MESHERY</Button>
                    </Container> */}

                </div>
            </NoSsr>
        );
    }
}

export default withStyles(styles)(Configuration);


