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

const styles = (theme) => ({
    main: {
        width: "50%",
        float: "left"
    },
    loc: {
        width: "50%",
        float: "left" 
    },
    img: {
        width: "auto",
        height: "50px",
    },
    offConfig: {
        width: "auto",
        height: "50px",
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
                    <Typography Style="margin-bottom:2rem">Location</Typography>
                    <Container>
                        <Container className={classes.loc}>
                            {
                                dockerChecked &&
                                <img
                                src={docker}
                                className={classes.img}
                                alt="docker"
                                />
                            }{
                                !dockerChecked && 
                                <img
                                src={docker}
                                className={classes.offConfig}
                                alt="docker"
                                /> 
                            }
                            <Container>
                                <Switch
                                 checked={dockerChecked}
                                 color="primary"
                                 onChange={this.handleToggle("docker")}
                                 inputProps={{ 'aria-label': 'controlled' }}
                                />
                                <Typography>Out of Cluster</Typography>
                            </Container>
                        </Container>
                            
                        <Container className={classes.loc}>
                            {
                                k8sChecked && 
                                <img
                                src={k8s}
                                className={classes.img}
                                alt="k8s"
                                />
                            }
                            {
                                !k8sChecked &&
                                <img
                                src={k8s}
                                className={classes.offConfig}
                                alt="k8s"
                                />
                            }
                            <Container>
                                <Switch
                                checked={k8sChecked}
                                color="primary"
                                onChange={this.handleToggle("k8s")}
                                inputProps={{ 'aria-label': 'controlled' }}
                                />
                                <Typography>In Cluster</Typography>
                            </Container>
                        </Container>
                    </Container>
                </Container>
            </NoSsr>
        );
    }
}

export default withStyles(styles)(Location);


