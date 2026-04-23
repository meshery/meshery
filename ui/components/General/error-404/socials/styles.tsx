import { styled, BrandColors } from '@/theme';
import Message from '@/assets/icons/general/message';
import Docker from '@/assets/icons/social/docker';
import Github from '@/assets/icons/social/github';
import Slack from '@/assets/icons/social/slack';
import Twitter from '@/assets/icons/social/twitter';
import Youtube from '@/assets/icons/social/youtube';

export const SocialMain = styled('div')(() => ({
  padding: '1.3rem 2rem',
  display: 'flex',
  justifyContent: 'center',
}));

export const SocialContainer = styled('div')(() => ({
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
}));

export const MessageIcon = styled(Message)(({ theme }) => ({
  fill: theme.palette.icon.default,
  color: theme.palette.icon.default,
  cursor: 'pointer',
  transition: 'all .3s',
  '&:hover': {
    fill: theme.palette.icon.default,
    color: theme.palette.icon.default,
  },
}));

export const SlackIcon = styled(Slack)(({ theme }) => ({
  fill: theme.palette.icon.default,
  color: theme.palette.icon.default,
  cursor: 'pointer',
  transition: 'all .3s',
  '&:hover': {
    fill: BrandColors.slack,
    color: BrandColors.slack,
  },
}));

export const TwitterIcon = styled(Twitter)(({ theme }) => ({
  fill: theme.palette.icon.default,
  color: theme.palette.icon.default,
  cursor: 'pointer',
  transition: 'all .3s',
  '&:hover': {
    fill: BrandColors.twitter,
    color: BrandColors.twitter,
  },
}));

export const TwitterHandleIcon = styled(Twitter)(({ theme }) => ({
  fill: theme.palette.keppelGreen,
  color: theme.palette.keppelGreen,
  cursor: 'pointer',
  transition: 'all .3s',
  '&:hover': {
    fill: BrandColors.twitter,
    color: BrandColors.twitter,
  },
}));

export const GithubHandleIcon = styled(Github)(({ theme }) => ({
  fill: theme.palette.keppelGreen,
  color: theme.palette.keppelGreen,
  cursor: 'pointer',
  transition: 'all .3s',
  '&:hover': {
    fill: BrandColors.github,
    color: BrandColors.github,
  },
}));

export const GithubIcon = styled(Github)(({ theme }) => ({
  fill: theme.palette.icon.default,
  color: theme.palette.icon.default,
  cursor: 'pointer',
  transition: 'all .3s',
  '&:hover': {
    fill: BrandColors.github,
    color: BrandColors.github,
  },
}));

export const YoutubeIcon = styled(Youtube)(({ theme }) => ({
  fill: theme.palette.icon.default,
  color: theme.palette.icon.default,
  cursor: 'pointer',
  transition: 'all .3s',
  '&:hover': {
    fill: BrandColors.youtube,
    color: BrandColors.youtube,
  },
}));

export const DockerIcon = styled(Docker)(({ theme }) => ({
  fill: theme.palette.icon.default,
  color: theme.palette.icon.default,
  cursor: 'pointer',
  transition: 'all .3s',
  '&:hover': {
    fill: BrandColors.docker,
    color: BrandColors.docker,
  },
}));
