// @ts-check
import React from "react";
import { makeStyles } from "@material-ui/core/styles";
import Accordion from "@material-ui/core/Accordion";
import AccordionSummary from "@material-ui/core/AccordionSummary";
import AccordionDetails from "@material-ui/core/AccordionDetails";
import Typography from "@material-ui/core/Typography";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import { promisifiedDataFetch } from "../../lib/data-fetch";
import { CircularProgress } from "@material-ui/core";
import PatternServiceForm from "./PatternServiceForm";
import {
  getPatternServiceName as getItemName,
  getPatternServiceID as getItemID,
  getPatternServiceType,
  getHumanReadablePatternServiceName as getReadableItemName
} from "./helpers";
import { isEmptyObj } from "../../utils/utils";
import { useSnackbar } from 'notistack';

const useStyles = makeStyles((theme) => ({
  accordionRoot : {
    width : "100%",
  },
  heading : {
    fontSize : theme.typography.pxToRem(15),
    fontWeight : theme.typography.fontWeightRegular,
  },
}));

// Question: So this function only returns the response coming from "/api/oam/(workload/trait)/--"
// So, it is not present in the current schemaSet because we had passed trim=true
// Discussion: Can we make network request without trim so that everything is fetched at once?
// Probable-Problem: It may be too large response.
async function fetchJSONSchema(name, type, id) {
  if (type !== "workload" && type !== "trait") throw Error("type can be either \"workload\" or \"trait\"");
  const url = `/api/oam/${type}/${name}/${id}`;

  const res = await promisifiedDataFetch(url);
  return JSON.parse(res?.oam_ref_schema) || {};
}

export async function getWorkloadTraitAndType(schemaSet) {
  // Get the schema sets for the workload
  const workloadSchema = await fetchJSONSchema(getItemName(schemaSet?.workload, false), "workload", getItemID(schemaSet?.workload));
  workloadSchema._internal = { patternAttributeName : getItemName(schemaSet?.workload, false) };

  // Get the schema sets for the traits
  const traitsSchemas = await Promise.all(schemaSet?.traits?.map(async t => {
    const schema = await fetchJSONSchema(getItemName(t, false), "trait", getItemID(t));

    schema._internal = { patternAttributeName : getItemName(t, false) };

    return schema;
  }));

  const type = getPatternServiceType(schemaSet?.workload);
  return { workload : workloadSchema, traits : traitsSchemas, type };
}

/**
 * LazyPatternServiceForm renders a accordion with the workload details &
 * traits schema lazily and does
 *
 * @param {{
 *  schemaSet: { workload: any, traits: any[], type: string };
 *  onSubmit: Function;
 *  onDelete: Function;
 *  namespace: string;
 *  onChange?: Function
 *  onSettingsChange?: Function;
 *  onTraitsChange?: Function;
 *  formData?: Record<String, unknown>
 *  reference?: Record<any, any>;
 * }} props
 *
 * @returns
 */
export default function LazyPatternServiceForm(props) {
  const [expanded, setExpanded] = React.useState(false);
  const [schemaSet, setSchemaSet] = React.useState({});
  const { enqueueSnackbar } = useSnackbar();
  const classes = useStyles();

  async function expand(state) {
    if (!state) {
      setExpanded(false);
      return;
    }

    setExpanded(true);
    try {
      // schemasets are not un-changing property,
      // thus only trigger stateChange when required
      if (isEmptyObj(schemaSet)) {
        // Get the schema sets consisting the workloads, traits and type
        const { workload, traits, type } = await getWorkloadTraitAndType(props?.schemaSet);
        setSchemaSet({
          workload,
          traits,
          type,
        });
      }
    } catch (error) {
      console.error("error getting schema:", { error })
      enqueueSnackbar(`error getting schema: ${error?.message}`, { variant : "error" })
    }
  }

  return (
    <div className={classes.accordionRoot}>
      <Accordion elevation={0} expanded={expanded} onChange={() => expand(!expanded)}>
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Typography className={classes.heading}>{getReadableItemName(props?.schemaSet?.workload)}</Typography>
        </AccordionSummary>
        <LazyAccordionDetails expanded={expanded}>
          {isEmptyObj(schemaSet) ? <CircularProgress /> : <PatternServiceForm {...props}
            // @ts-ignore
            schemaSet={schemaSet} />}
        </LazyAccordionDetails>
      </Accordion>
    </div>
  );
}

function LazyAccordionDetails(props) {
  if (!props.expanded) return <AccordionDetails />;

  // @ts-ignore // LEE: This behavior is more like what we need - https://codesandbox.io/s/upbeat-tesla-uchsb?file=/src/MyAccordion.js
  return <AccordionDetails>{props.children}</AccordionDetails>
}

