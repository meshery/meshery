import React from "react";
import { 
    NoSsr, 
    Container, 
    Typography, 
    Switch,
} from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";
import docker from "../../img/location/docker.png"
import k8s from "../../img/location/kubernetes.png"
import { CustomTypography } from "../CustomTypography";
import DockerIcon from "../../img/SVGs/dockerIcon";
import KubernetesIcon from "../../img/SVGs/kubernetesIcon";

const styles = (theme) => ({
    main: {
        width: "50%",
        float: "left",
        paddingBottom: theme.spacing(6),
        backgroundColor:"#393F49",
        borderBottomLeftRadius: "20px"
    },
    loc: {
        width: "50%",
        float: "left" 
    },
    offConfig: {
        filter: "grayscale(1) invert(0.35)"
    }
});

class Location extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            dockerChecked: true,
            k8sChecked: false
        }
    }

    handleToggle = (name) => () => {
        const self = this;
        if(name == 'docker'){
            self.setState((state) => ({ dockerChecked : !state.dockerChecked }), () => this.handleChange(name));
        } else {
            self.setState((state) => ({ k8sChecked : !state.k8sChecked }), () => this.handleChange(name));
        }
    }

    handleChange = (name) => {
        if(name == 'docker'){
            console.log("docker:", this.state.dockerChecked);
        } else {
            console.log("k8s:", this.state.k8sChecked);
        }
    }

    render() {

        const { classes } = this.props;
        const { dockerChecked, k8sChecked } = this.state;

        return (
            <NoSsr>
                <Container className={classes.main}>
                    <CustomTypography Style="margin-bottom:2rem">Location</CustomTypography>
                
                        <div className={classes.loc}>
                        <div className={dockerChecked  ? null : classes.offConfig}><DockerIcon width={57} height={40} /></div>                               
                         
                                  <Switch
                                 checked={dockerChecked}
                                 color="primary"
                                 onChange={this.handleToggle("docker")}
                                 inputProps={{ 'aria-label': 'controlled' }}
                                />
                                <CustomTypography>Out of Cluster</CustomTypography>
                            
                        </div>
                            
                        <div className={classes.loc}>
                        <div className={k8sChecked ? null : classes.offConfig}><KubernetesIcon width={40} height={40} /></div>                               
                                <Switch
                                checked={k8sChecked}
                                color="primary"
                                onChange={this.handleToggle("k8s")}
                                inputProps={{ 'aria-label': 'controlled' }}
                                />
                                <CustomTypography>In Cluster</CustomTypography>
                        </div>
                </Container>
            </NoSsr>
        );
    }
}

export default withStyles(styles)(Location);


