import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import {
  Card, CardContent, Container, Typography
} from '@material-ui/core';


// const MeshySwitch = withStyles({
//   switchBase: {
//     color: "grey",
//     "&$checked": {
//       color: "#00B39F",
//     },
//     "&$checked + $track": {
//       backgroundColor: "#00B39F",
//     },
//   },
//   checked: {},
//   track: {},
// })(Switch);

const styles = () => ({

  cardContainer : {
    position : "relative",
    display : "flex",
    flexDirection : "column",
    justifyContent : "space-between",
    // padding: "2rem 6rem",
  },

  // Card
  card : {
    position : "relative",
    width : "16rem",
    minWidth : "10rem",
    border : "1px solid gray",
    borderRadius : "0.75rem",
    //top: "2rem",
    margin : "0rem 0rem 6rem 0rem",
    ["@media (max-width:1024px)"] : {
      //eslint-disable-line no-useless-computed-key
      margin : "0rem 0rem 6rem 0", },
  },
  cardChecked : { height : "15rem",
    marginBottom : "1rem", },
  cardUnchecked : { height : "10rem", },
  cardContent : { height : "10rem", //change this to increase the card size for adding configuration
    width : "100%",
    padding : "0 !important", },
  cardContentChecked : { height : "100%", //change this to increase the card size for adding configuration
    width : "100%",
    padding : "0 !important", },
  contentTop : {
    background : "#434343",
    height : "10rem",
    width : "100%",
    display : "flex",
    justifyContent : "center",
    alignItems : "center",
    transition : "backgroudn 1s ease"
  },
  contentTopConnected : {
    background : "white",
    height : "10rem",
    width : "100%",
    display : "flex",
    justifyContent : "center",
    alignItems : "center",
    transition : "backgroudn 1s ease"
  },
  contentTopSwitcher : { paddingLeft : "2rem", },
  iconContainer : {
    display : "flex",
    flexDirection : "column",
    alignItems : "center",
    justifyContent : "center",
    // marginLeft: "1rem",
  },
  cardIcon : { width : "3rem", },
  cardIconText : {
    color : "white",
    fontSize : "0.85rem",
    textAlign : "center",
    "&:first-letter" : { textTransform : "capitalize", },
    marginTop : "0.8rem"
  },
  cardIconTextConnected : {
    color : "black",
    fontSize : "0.85rem",
    marginTop : "0.8rem",
    textAlign : "center",
    "&:first-letter" : { textTransform : "capitalize", },
  },
});




const ServiceSwitch = ({ serviceInfo, classes,isConnected }) => {

  const ServiceIcon = serviceInfo.logoComponent
  const ConfigComponent = serviceInfo.configComp

  const showConfigComponent = () => !isConnected

  return (
    <Container className={ classes.cardContainer} key={ serviceInfo.name }>
      <Card className={`${classes.card} `} variant="outlined">
        <CardContent className={ showConfigComponent()
          ? classes.cardContentChecked
          : classes.cardContent}>
          <div className={isConnected
            ? classes.contentTopConnected
            : classes.contentTop}>
            <div className={classes.iconContainer}>
              <ServiceIcon isActive={isConnected}/>
              <Typography className={isConnected
                ? classes.cardIconTextConnected
                : classes.cardIconText}  color="primary">
                {serviceInfo.name}
              </Typography>
            </div>
            {/* <FormControlLabel
              className={classes.contentTopSwitcher}
              control={<MeshySwitch checked={isChecked} name={`{serviceInfo.name}`} />}
              onChange={handleSwitch}
            />*/}
          </div>
          {showConfigComponent() &&
          <div style={{ paddingLeft : "0.6rem", paddingRight : "0.6rem" }}>
            {ConfigComponent}
          </div>
          }

        </CardContent>
      </Card>
    </Container>
  )
}

ServiceSwitch.propTypes = { serviceInfo : PropTypes.object.isRequired, };

export default withStyles(styles)(ServiceSwitch);
