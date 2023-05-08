import React from "react";
import { styled } from "@mui/material";
import { Avatar } from "./Avatar";

export default {
  title : "Components/Avatar",
  component : Avatar,
};

const CustomAvatar = styled(Avatar)(({ theme }) => ({
  marginRight : theme.spacing(1),
  marginBottom : theme.spacing(1),
}));

export function Default(props) {
  return (
    <CustomAvatar {...props} />
  )
}

export function Icons(props) {
  return (
    <CustomAvatar src="/static/img/kubernetes.svg" {...props} />
  )
}