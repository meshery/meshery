import { Box, Card, CardActions, CardHeader, IconButton } from "@mui/material";
import AddIcon from '@mui/icons-material/Add';
import DeleteIcon from '@mui/icons-material/Delete';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';

function CategoryCard({cat}) {
  if (typeof cat === "undefined") {
    cat = 0;
  }
  let content;
  let description;
  switch (cat) {
    case 0:
      content = "Manage Service Mesh Lifecycle";
      description = "Deploy a service mesh or SMI adapter on your cluster.";
      break;

    case 1:
      content = "Manage Sample Application Lifecycle";
      description = "Deploy sample applications on/off the service mesh.";
      break;

    case 2:
      content = "Apply Service Mesh Configuration";
      description = "Configure your service mesh using some pre-defined options.";
      break;

    case 3:
      content = "Validate Service Mesh Configuration";
      description = "Validate your service mesh configuration against best practices.";
      break;

    case 4:
      content = "Apply Custom Configuration";
      description = "Customize the configuration of your service mesh.";
      break;
  }
  return (
    <Card sx={
      {
        height:'100%',
        display:'flex',
        flexDirection:'column',
      }
    }>
      <CardHeader title={content} subheader={description} sx={{ flexGrow : 1 }} />
      <CardActions disableSpacing>
        <IconButton
          aria-label="install"
        >
          {cat !== 4
            ? <AddIcon />
            : <PlayArrowIcon />}
        </IconButton>
        {cat !== 3 && (
          <Box sx={{width:"100%"}}>
            <IconButton
              aria-label="delete"
              sx={{float:'right'}}
            >
              <DeleteIcon />
            </IconButton>
          </Box>
        )}
      </CardActions>
    </Card>
  );
}

export default CategoryCard;