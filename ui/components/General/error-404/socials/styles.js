import { styled } from "@material-ui/core";
import Message from "@/assets/icons/general/message";
import Docker from "@/assets/icons/social/docker";
import Github from "@/assets/icons/social/github";
import Slack from "@/assets/icons/social/slack";
import Twitter from "@/assets/icons/social/twitter";
import Youtube from "@/assets/icons/social/youtube";

export const SocialMain = styled("div")(({ theme }) => ({
  padding: "1.3rem 2rem",
  display: "flex",
  justifyContent: "center"
}));

export const SocialContainer = styled("div")(({ theme }) => ({
  display: "flex",
  justifyContent: "center",
  alignItems: "center"
}));

export const MessageIcon = styled(Message)(({ theme }) => ({
  fill: theme.palette.limedSpruce,
  color: theme.palette.limedSpruce,
  cursor: "pointer",
  transition: "all .3s",
  "&:hover": {
    fill: theme.palette.limedSpruce,
    color: theme.palette.limedSpruce
  }
}));

export const SlackIcon = styled(Slack)(({ theme }) => ({
  fill: theme.palette.limedSpruce,
  color: theme.palette.limedSpruce,
  cursor: "pointer",
  transition: "all .3s",
  "&:hover": {
    fill: theme.palette.slack,
    color: theme.palette.slack
  }
}));

export const TwitterIcon = styled(Twitter)(({ theme }) => ({
  fill: theme.palette.limedSpruce,
  color: theme.palette.limedSpruce,
  cursor: "pointer",
  transition: "all .3s",
  "&:hover": {
    fill: theme.palette.twitter,
    color: theme.palette.twitter
  }
}));

export const TwitterHandleIcon = styled(Twitter)(({ theme }) => ({
  fill: theme.palette.keppelGreen,
  color: theme.palette.keppelGreen,
  cursor: "pointer",
  transition: "all .3s",
  "&:hover": {
    fill: theme.palette.twitter,
    color: theme.palette.twitter
  }
}));

export const GithubHandleIcon = styled(Github)(({ theme }) => ({
  fill: theme.palette.keppelGreen,
  color: theme.palette.keppelGreen,
  cursor: "pointer",
  transition: "all .3s",
  "&:hover": {
    fill: theme.palette.github,
    color: theme.palette.github
  }
}));

export const GithubIcon = styled(Github)(({ theme }) => ({
  fill: theme.palette.limedSpruce,
  color: theme.palette.limedSpruce,
  cursor: "pointer",
  transition: "all .3s",
  "&:hover": {
    fill: theme.palette.github,
    color: theme.palette.github
  }
}));

export const YoutubeIcon = styled(Youtube)(({ theme }) => ({
  fill: theme.palette.limedSpruce,
  color: theme.palette.limedSpruce,
  cursor: "pointer",
  transition: "all .3s",
  "&:hover": {
    fill: theme.palette.youtube,
    color: theme.palette.youtube
  }
}));

export const DockerIcon = styled(Docker)(({ theme }) => ({
  fill: theme.palette.limedSpruce,
  color: theme.palette.limedSpruce,
  cursor: "pointer",
  transition: "all .3s",
  "&:hover": {
    fill: theme.palette.docker,
    color: theme.palette.docker
  }
}));
