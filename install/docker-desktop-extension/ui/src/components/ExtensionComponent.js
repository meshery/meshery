import React from "react";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import { NoSsr, Typography } from "@material-ui/core";
import FavoriteIcon from '@material-ui/icons/Favorite';
import MesheryLogo from "./MesheryLogo";
import Configuration from "./Config/Configuration";
import { CustomTypography } from "./CustomTypography";
import {
Container, 
Button
} from "@material-ui/core";

const styles = (theme) => ({
  root: {
    textAlign: "center",
    backgroundColor: "#222C32",
    padding: "1rem",
    minHeight: "100vh"
  },
  container: {
    width: "60%",
    marginLeft: "auto",
    marginRight: "auto",
    marginTop: theme.spacing(2),
  },
  footerText: {
    cursor: 'pointer',
    display: 'inline',
    verticalAlign: 'middle',
    color: "#CCCCCC",
  },
  footerIcon: {
    display: 'inline',
    verticalAlign: 'top',
  },
  mesheryConfig1: {
    margin: theme.spacing(6),
    backgroundColor: "#5AA9F1"
}, mesheryConfig: {
  margin: theme.spacing(6),

}
});

class ExtensionComponent extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
    };
  }

  handleL5CommunityClick = () => {
    if (typeof window !== 'undefined') {
      const w = window.open('https://layer5.io', '_blank');
      w.focus();
    }
  }


  render() {
    const { classes } = this.props;

    return (
      <NoSsr>
        <div className={classes.root}>
          <MesheryLogo />
          <Configuration />
          <Container>
                        <Button className={classes.mesheryConfig1} variant="contained">DEPLOY MESHERY</Button>
                        <Button className={classes.mesheryConfig} variant="contained">OPEN MESHERY</Button>
                    </Container>
                    <div className={classes.footer} >
          <Typography variant="body2" align="center" component="p">
                <span onClick={this.handleL5CommunityClick} className={classes.footerText}>
                  Built with {' '} <FavoriteIcon className={classes.footerIcon} /> {' '} by the Layer5 Community
                </span>
          </Typography>
          </div>
        </div>
      </NoSsr>
    );
  }
}

ExtensionComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(ExtensionComponent);
