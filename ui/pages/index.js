import { AdaptersChipList, AdaptersListContainer, ComponentsStatusContainer } from "@/features/mesheryComponents";
import { Button, useTheme } from "@material-ui/core";
import React from "react";

export default function Home() {
  // const theme = useTheme();
  // return <div style={{ color: theme.palette.newcolor }}>Hello</div>;
  return <AdaptersListContainer render={() => <Button>Hello</Button>} />;
}
