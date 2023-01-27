import React from 'react'
import PropTypes from 'prop-types'
import NoSsr from '@mui/material/NoSsr'

export default function ProviderLayout({ children }) {
  return (
    <>
      <NoSsr>
        <div data-cy="root" style={{ padding : '170px 0px', textAlign : 'center' }}>
          {children}
        </div>
      </NoSsr>
    </>
  )
}

ProviderLayout.propTypes = {
  children : PropTypes.node.isRequired
}
