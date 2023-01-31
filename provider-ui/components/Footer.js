import React, { Fragment } from 'react'
import { styled } from '@mui/material/styles'
import Paper from '@mui/material/Paper'
import NoSsr from '@mui/material/NoSsr'
import Typography from '@mui/material/Typography'
import FavoriteIcon from '@mui/icons-material/Favorite'

const Item = styled(Paper)(({ theme }) => ({
  backgroundColor : theme.palette.background.paper,
  padding : theme.spacing(2),
  color : '#737373',
  border : '0'
}))

export default function Footer() {
  const handleL5CommunityClick = () => {
    if (typeof window !== 'undefined') {
      const w = window.open('https://layer5.io', '_blank')
      w.focus()
    }
  }

  return (
    <Fragment>
      <NoSsr>
        <Item component="footer" square variant="outlined">
          <Typography variant="body2" align="center" color="textSecondary" component="p">
            <span onClick={handleL5CommunityClick} style={{
              cursor : 'pointer',
              display : 'inline',
              verticalAlign : 'middle'
            }}>
              Built with
              {' '}
              <FavoriteIcon sx={{
                display : 'inline',
                verticalAlign : 'top'
              }} />
              {' '}
              by the Layer5 Community
            </span>
          </Typography>
        </Item>
      </NoSsr>
    </Fragment>
  )
}
