import React from 'react'
import PropTypes from 'prop-types'
import Container from '@mui/material/Container'
import NoSsr from '@mui/material/NoSsr'

export default function ProviderLayout ({ children }) {
  return (
    <NoSsr>
      <Container maxWidth="false">
        <div data-cy="root" style={{
          padding: '170px 0px',
          textAlign: 'center'
        }}>
          {children}
        </div>
      </Container>
    </NoSsr>
  )
}

ProviderLayout.propTypes = {
  children: PropTypes.node.isRequired
}
