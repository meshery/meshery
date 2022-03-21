import React from "react";
import {
    NoSsr, 
    Container, 
    Typography, 
    Switch,
} from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";
import consul from "../../img/service-meshes/consul.png";
import istio from "../../img/service-meshes/istio.png";
import kuma from "../../img/service-meshes/kuma.png";
import linkerd from "../../img/service-meshes/linkerd.png";
// import nginx from "../../img/service-meshes/nginx.svg";

const styles = (theme) => ({
    main: {
        width: "50%",
        float: "left",
        paddingBottom: theme.spacing(3)
    },
    sm: {
        width: "25%",
        float: "left"
    },
    img: {
        width: "auto",
        height: "50px"
    },
    offConfig: {
        width: "auto",
        height: "50px",
        filter: "grayscale(1) invert(0.35)"
    }
});

class Configuration extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            consulChecked: true,
            istioChecked: false,
            linkerdChecked: false,
            nginxChecked: false,
            kumaChecked: false
        }
    }

    handleToggle = (name) => () => {
        const self = this;
        if(name == 'consul'){
            self.setState((state) => ({ consulChecked : !state.consulChecked }), () => this.handleChange(name));
        } else if (name == 'istio') {
            self.setState((state) => ({ istioChecked : !state.istioChecked }), () => this.handleChange(name));
        } else if (name == 'linkerd') {
            self.setState((state) => ({ linkerdChecked : !state.linkerdChecked }), () => this.handleChange(name));
        } else if (name == "nginx") {
            self.setState((state) => ({ nginxChecked : !state.nginxChecked }), () => this.handleChange(name));
        } else {
            self.setState((state) => ({ kumaChecked : !state.kumaChecked }), () => this.handleChange(name));
        }
    }

    handleChange = (name) => {
        if(name == 'consul'){
            console.log("consul:", this.state.dockerChecked);
        } else if (name == 'istio') {
            console.log("istio:", this.state.istioChecked);
        } else if (name == 'linkerd') {
            console.log("linkerd:", this.state.linkerdChecked);
        } else if (name == "nginx") {
            console.log("nginx:", this.state.nginxChecked);
        } else {
            console.log("kuma:", this.state.kumaChecked);
        }
    }

    render() {
        const { classes } = this.props;
        const {
            consulChecked, 
            istioChecked, 
            linkerdChecked, 
            nginxChecked, 
            kumaChecked 
        } = this.state;

        return (
            <NoSsr>
                <Container className={classes.main}>
                    <Typography Style="margin-bottom:2rem">Service Mesh</Typography>
                    <Container>
                        <Container className={classes.sm}>
                            {
                                consulChecked &&
                                <img
                                 src={consul}
                                 className={classes.img}
                                 alt="consul"
                                />
                            }
                            {
                                !consulChecked &&
                                <img
                                 src={consul}
                                 className={classes.offConfig}
                                 alt="consul"
                                />
                            }
                            <Container>
                                <Switch
                                checked={consulChecked}
                                color="primary"
                                onChange={this.handleToggle("consul")}
                                inputProps={{ 'aria-label': 'controlled' }}
                                />
                            </Container>
                        </Container>
                        <Container className={classes.sm}>
                            { istioChecked &&
                                <img
                                 src={istio}
                                 className={classes.img}
                                 alt="istio"
                                />
                            }
                            {
                                !istioChecked && 
                                <img
                                 src={istio}
                                 className={classes.offConfig}
                                 alt="istio"
                                />
                            }
                            <Container>
                                <Switch
                                checked={istioChecked}
                                color="primary"
                                onChange={this.handleToggle("istio")}
                                inputProps={{ 'aria-label': 'controlled' }}
                                />
                            </Container>
                        </Container>
                        <Container className={classes.sm}>
                            {
                                linkerdChecked && 
                                <img
                                 src={linkerd}
                                 className={classes.img}
                                 alt="linkerd"
                                />
                            }
                            {
                                !linkerdChecked &&
                                <img
                                 src={linkerd}
                                 className={classes.offConfig}
                                 alt="linkerd"
                                />
                            }
                            <Container>
                                <Switch
                                checked={linkerdChecked}
                                color="primary"
                                onChange={this.handleToggle("linkerd")}
                                inputProps={{ 'aria-label': 'controlled' }}
                                />
                            </Container>

                        </Container>
                        {/* <Container className={classes.img}>
                            <img
                             src={nginx}
                             className={classes.img}
                             alt="nginx"
                            />
                            <Switch
                             checked={nginxChecked}
                             color="primary"
                             onChange={this.handleToggle("nginx")}
                             inputProps={{ 'aria-label': 'controlled' }}
                            />
                        </Container> */}
                        <Container className={classes.sm}>
                            {
                                kumaChecked &&
                                <img
                                 src={kuma}
                                 className={classes.img}
                                 alt="kuma"
                                />
                            }
                            {
                                !kumaChecked && 
                                <img
                                 src={kuma}
                                 className={classes.offConfig}
                                 alt="kuma"
                                />
                            }

                            <Container>
                                <Switch
                                checked={kumaChecked}
                                color="primary"
                                onChange={this.handleToggle("kuma")}
                                inputProps={{ 'aria-label': 'controlled' }}
                                />
                            </Container>
                        </Container>
                    </Container>
                </Container>
            </NoSsr>
        );
    }
}

export default withStyles(styles)(Configuration);


