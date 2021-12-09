import * as React from "react";
import { Grid, Box, Typography, Button, CardContent, CardActions, Card } from "@mui/material/";

// eslint-disable-next-line react/prop-types
export default function BigNumber({ title, description, buttonText, handleClick }) {
  return (
    <Card sx={{ minWidth: 275, p: 1 }}>
      <CardContent>
        <Grid container direction="row" justifyContent="center" alignItems="baseline" spacing={2}>
          <Grid item>
            <Typography variant="h2" component="h2" color="text.secondary" gutterBottom>
              {title}
            </Typography>
          </Grid>
          <Grid item>
            <Typography color="text.secondary">{description}</Typography>
          </Grid>
        </Grid>
      </CardContent>
      <CardActions>
        <Box sx={{ flexGrow: 1 }} />
        <Button variant="contained" onClick={handleClick}>
          {buttonText}
        </Button>
      </CardActions>
    </Card>
  );
}
