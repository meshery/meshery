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

const styles = (theme) => ({
    main: {
        paddingBottom: theme.spacing(4),
    },
    config: {
        margin: theme.spacing(5),
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
                    <Typography variant="h5" className={classes.main}>
                        Configure your Meshery deployment
                    </Typography>

                    <Location />
                    <ServiceMesh />

                    <Container>
                        <Button className={classes.mesheryConfig} color="primary" variant="contained">DEPLOY MESHERY</Button>
                        <Button className={classes.mesheryConfig} variant="contained">OPEN MESHERY</Button>
                    </Container>

                </div>
            </NoSsr>
        );
    }
}

export default withStyles(styles)(Configuration);


