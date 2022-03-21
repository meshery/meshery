import React from "react";
import {
    NoSsr, 
    Container, 
    Typography, 
    Switch,
} from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";
import { CustomTypography } from "../CustomTypography";
import ConsulIcon from "../../img/SVGs/consulIcon";
import IstioIcon from "../../img/SVGs/IstioIcon";
import KumaIcon from "../../img/SVGs/kumaIcon";
import LinkerdIcon from "../../img/SVGs/linkerdIcon";
import NginxIcon from "../../img/SVGs/nginxIcon";



const styles = (theme) => ({
    main: {
        width: "50%",
        float: "left",
        paddingBottom: theme.spacing(9),
        backgroundColor:"#393F49",
        borderBottomRightRadius: "20px"
    },
    sm: {
        width: "20%",
        float: "left",
        flexDirection: "row"
    },
    img: {
        width: "auto",
        height: "50px"
    },
    offConfig: {
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
                <div className={classes.main}>
                    <CustomTypography Style="margin-bottom:2rem">Service Mesh</CustomTypography>
                        <div className={classes.sm}>
                           <div className={consulChecked ? null : classes.offConfig}><ConsulIcon width={40} height={40} /></div>                               
                           <Switch
                                checked={consulChecked}
                                color="primary"
                                onChange={this.handleToggle("consul")}
                                inputProps={{ 'aria-label': 'controlled' }}
                                />
                           
                       </div>
                       <div className={classes.sm}>
                           <div className={istioChecked ? null : classes.offConfig}><IstioIcon width={40} height={40} /></div>                               
                          
                                <Switch
                                checked={istioChecked}
                                color="primary"
                                onChange={this.handleToggle("istio")}
                                inputProps={{ 'aria-label': 'controlled' }}
                                />
                          </div>
                          <div className={classes.sm}>
                          <div className={linkerdChecked ? null : classes.offConfig}><LinkerdIcon width={40} height={40} /></div>                               
                        
                                <Switch
                                checked={linkerdChecked}
                                color="primary"
                                onChange={this.handleToggle("linkerd")}
                                inputProps={{ 'aria-label': 'controlled' }}
                                />
                        </div>
                        <div className={classes.sm}>
                        <div className={nginxChecked ? null : classes.offConfig}><NginxIcon width={38} height={40} /></div>                               
                            <Switch
                             checked={nginxChecked}
                             color="primary"
                             onChange={this.handleToggle("nginx")}
                             inputProps={{ 'aria-label': 'controlled' }}
                            />
                            </div>
                     <div className={classes.sm}>
                    <div className={kumaChecked ? null : classes.offConfig}><KumaIcon width={40} height={40} /></div>                               

                         
                                <Switch
                                checked={kumaChecked}
                                color="primary"
                                onChange={this.handleToggle("kuma")}
                                inputProps={{ 'aria-label': 'controlled' }}
                                />
                          </div>
                </div>
            </NoSsr>
        );
    }
}

export default withStyles(styles)(Configuration);


