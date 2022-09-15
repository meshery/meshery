import React from 'react';
import {Badge, Box} from "@mui/material";
import { styled } from "@mui/material/styles";


const StyledBadge = styled(Badge)(({ theme }) => ({
    boxShadow : `0 0 0 2px ${theme.palette.background.paper}`,
}));

const CustomBox = styled(Box)(({ theme }) => ({
  display : 'flex',
    '& > *' : {
      marginLeft : theme.spacing(0.5),
      marginRight : -theme.spacing(0.75),
    },

}));


export default function BadgeAvatars({ children }) {

  return (
    <CustomBox>
      <StyledBadge
        overlap="circular"
        anchorOrigin={{
          vertical : 'bottom',
          horizontal : 'right',
        }}
        variant="dot"
      >
        {children}
      </StyledBadge>
    </CustomBox>
  );
}
