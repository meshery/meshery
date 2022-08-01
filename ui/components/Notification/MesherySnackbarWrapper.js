import react, {useState} from "react";
import { styled } from "@mui/material/styles";
import {
    IconButton,
    Typography,
    Paper,
    Collapse,
    Card,
    CardActions,
    Grid,
    Snackbar
  } from "@mui/material";   

  function MesherySnackbarWrapper ({ message, onClose, details}) {

    const [expanded, setExpanded] = useState(false);
    const [cardHover, setCardHover] = useState(false)

    const handleExpandClick = () => {
        setExpanded(!expanded);
      };

    return(
        <Snackbar>
            <Card 
        aria-label="Show more"
        onMouseEnter={() => setCardHover(true)}
        onMouseLeave={() => setCardHover(false)}
      >
         <CardActions onClick={handleExpandClick}>
          <Grid container direction="row" justify="space-between" alignItems="center" wrap="nowrap">
            <Typography variant="subtitle2">
              <div style={{ display : "flex", alignItems : "center" }}>
                {/* <Icon /> */}
                <div>{message}</div>
              </div>
            </Typography>
            <Grid container item xs={4} justify="flex-end">
              <IconButton
                aria-label="Show more"
                
                onClick={handleExpandClick}
              >
                <ExpandMoreIcon />
              </IconButton>

              <IconButton  onClick={onClose}>
                <DoneIcon    
                  onMouseEnter={() => setCardHover(false)}
                  onMouseLeave={() => setCardHover(true)}
                />
              </IconButton>
            </Grid>
          </Grid>
        </CardActions>
        <Collapse in={expanded} timeout="auto" unmountOnExit >
          <Paper  square variant="outlined" elevation={0}>
            <Typography variant="subtitle2" gutterBottom>DETAILS</Typography>
            {details}
          </Paper>
        </Collapse>
      </Card>
        </Snackbar>
    )
}

export default MesherySnackbarWrapper