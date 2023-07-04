import { Box, Button, IconButton, Toolbar, Typography } from "@mui/material";
import { AppBar } from "./AppBar";
import MenuIcon from "@mui/icons-material/Menu";
import Settings from "@mui/icons-material/Settings";
import Link from "next/link";

export default {
  title : "Components/AppBar",
  component : AppBar,
  tags: ['autodocs'],
    parameters: {
    // More on how to position stories at: https://storybook.js.org/docs/react/configure/story-layout
    layout: 'fullscreen',
  },
};

export function Primary() {
  return (
    <AppBar position="fixed" color="primary">
      <Toolbar>
        <Typography variant="h6" noWrap component="div">
                    Dashboard
        </Typography>
      </Toolbar>
    </AppBar>
  )
}

export function BasicAppBar() {
  return (
    <Box sx={{ flexGrow : 1 }}>
      <AppBar position="static">
        <Toolbar>
          <IconButton size="large" edge="start" color="inherit" aria-label="menu" sx={{ mr : 2 }}>
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" component="div" sx={{ flexGrow : 1 }}>
                        Dashboard
          </Typography>
          <Button color="inherit">Login</Button>
        </Toolbar>
      </AppBar>
    </Box>
  )
}

export function MesheryHeader() {
  const title = "Dashboard";
  const isBeta = "";

  return (
    <Box sx={{ flexGrow : 1 }}>
      <AppBar color="primary" position="sticky" elevation={2}>
        <Toolbar>
          <Typography variant="h6" noWrap component="div" sx={{ display : { xs : "none", sm : "block" } }}>
            {title}{isBeta ? <sup>BETA</sup> : ""}
          </Typography>
          <Box sx={{ flexGrow : 1 }} />
          <Box data-test="settings-button" sx={{ display : { xs : "none", md : "flex" } }}>
            <Link href="/settings">
              <IconButton>
                <Settings />
              </IconButton>
            </Link>
          </Box>
        </Toolbar>
      </AppBar>
    </Box>
  )
}