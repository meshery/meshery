import { Grid, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import clsx from 'clsx';
import React from 'react';
// import { ValueLabel } from '../Label';

const useStyles = makeStyles((theme) => ({
  metadataNameCell: {
    fontSize: '1rem',
    textAlign: 'left',
    maxWidth: '100%',
    minWidth: '10rem',
    verticalAlign: 'top',
    color: theme.palette.text.secondary,
    borderBottom: `1px solid ${theme.palette.divider}`,
    padding: '7px 12px',
    [theme.breakpoints.down('sm')]: {
      color: theme.palette.text.primary,
      fontSize: '1.5rem',
      minWidth: '100%',
      width: '100%',
      maxWidth: '100%',
      display: 'block',
      borderTop: `1px solid ${theme.palette.divider}`,
      borderBottom: `none`,
    },
  },
  metadataCell: {
    width: '100%',
    verticalAlign: 'top',
    fontSize: '1rem',
    overflowWrap: 'anywhere',
    padding: '7px 12px',
    borderBottom: `1px solid ${theme.palette.divider}`,
    [theme.breakpoints.down('sm')]: {
      color: theme.palette.text.secondary,
      minWidth: '100%',
      width: '100%',
      maxWidth: '100%',
      display: 'block',
      marginBottom: '2rem',
      borderBottom: `none`,
    },
  },
  metadataRow: {
    borderBottom: `1px solid ${theme.palette.divider}`,
  },
  metadataLast: {
    borderBottom: 'none',
  },
  table: {
    border: '1px solid #e7e7e7',
    borderRadius: theme.shape.borderRadius,
  },
  highlightRow: {
    color: theme.palette.tables.head.color,
    fontWeight: 'bold',
    background: theme.palette.tables.head.background,
  },
  valueLabel: {
    color: theme.palette.text.primary,
    fontSize: theme.typography.pxToRem(16),
    wordBreak: 'break-word',
  },
}));

// export interface NameValueTableRow {
//   /** The name (key) for this row */
//   name: string | JSX.Element;
//   /** The value for this row */
//   value?: string | JSX.Element | JSX.Element[];
//   /** Whether this row should be hidden (can be a boolean or a function that will take the
//    * @param value and return a boolean) */
//   hide?: boolean | ((value: NameValueTableRow['value']) => boolean);
//   /** Extra properties to pass to the value cell */
//   valueCellProps?: GridProps;
//   /** Whether to highlight the row (used for titles, separators, etc.). */
//   withHighlightStyle?: boolean;
// // }

// export interface NameValueTableProps {
//   rows: NameValueTableRow[];
//   valueCellProps?: GridProps;
// }

function ValueLabel(props) {
  const classes = useStyles();
  return (
    <Typography className={classes.valueLabel} component="span">
      {props.children}
    </Typography>
  );
}

function Value({ value }) {
  if (typeof value === 'undefined') {
    return null;
  } else if (typeof value === 'string') {
    return <ValueLabel>{value}</ValueLabel>;
  } else if (Array.isArray(value)) {
    return (
      <>
        {value.map((val, i) => (
          <Value value={val} key={i} />
        ))}
      </>
    );
  } else {
    return value;
  }
}

export default function NameValueTable(props) {
  const classes = useStyles();
  const { rows, valueCellProps: globalValueCellProps } = props;

  const visibleRows = React.useMemo(
    () =>
      rows.filter(({ value, hide = false }) => {
        let shouldHide = false;
        if (typeof hide === 'function') {
          shouldHide = hide(value);
        } else {
          shouldHide = hide;
        }

        return !shouldHide;
      }),
    [rows],
  );

  return (
    <Grid
      container
      component="dl" // mount a Definition List
      className={classes.table}
    >
      {visibleRows.map(
        ({ name, value, hide = false, withHighlightStyle = false, valueCellProps = {} }, i) => {
          let shouldHide = false;
          if (typeof hide === 'function') {
            shouldHide = hide(value);
          } else {
            shouldHide = hide;
          }

          if (shouldHide) {
            return null;
          }

          const last = visibleRows.length === i + 1;
          const { className, ...otherValueCellProps } = globalValueCellProps || {};

          const hideValueGridItem = withHighlightStyle && !value;

          const parts = name.split('.');
          const lastPart = parts[parts.length - 1];

          return (
            <>
              <Grid
                item
                key={i}
                xs={12}
                spacing={1}
                sm={hideValueGridItem ? 12 : 4}
                component="dt"
                className={clsx(
                  last ? classes.metadataLast : '',
                  classes.metadataNameCell,
                  withHighlightStyle ? classes.highlightRow : '',
                )}
              >
                {lastPart}
              </Grid>
              {!hideValueGridItem && (
                <Grid
                  item
                  key={i + 10000}
                  xs={12}
                  sm={8}
                  spacing={1}
                  component="dd"
                  className={clsx(
                    last ? classes.metadataLast : '',
                    classes.metadataCell,
                    className ? className : '',
                    withHighlightStyle ? classes.highlightRow : '',
                  )}
                  {...otherValueCellProps}
                  {...valueCellProps}
                >
                  <Value value={value} />
                </Grid>
              )}
            </>
          );
        },
      )}
    </Grid>
  );
}
