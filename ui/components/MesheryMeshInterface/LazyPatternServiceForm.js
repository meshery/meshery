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
import { getPatternServiceName as getItemName, getPatternServiceID as getItemID, getPatternServiceType } from "./helpers"

const useStyles = makeStyles((theme) => ({
  root: {
    width: "100%",
  },
  heading: {
    fontSize: theme.typography.pxToRem(15),
    fontWeight: theme.typography.fontWeightRegular,
  },
}));

async function fetchJSONSchema(name, type, id) {
  if (type !== "workload" && type !== "trait") throw Error("type can be either \"workload\" or \"trait\"");
  const url = `/api/oam/${type}/${name}/${id}`;

  const res = await promisifiedDataFetch(url);
  return JSON.parse(res?.oam_ref_schema) || {}
}

export async function getWorkloadTraitAndType(schemaSet) {
  // Get the schema sets for the workload
  const workloadSchema = await fetchJSONSchema(getItemName(schemaSet?.workload), "workload", getItemID(schemaSet?.workload));
  workloadSchema._internal = { patternAttributeName: getItemName(schemaSet?.workload, false) };

  // Get the schema sets for the traits
  const traitsSchemas = await Promise.all(schemaSet?.traits?.map(async t => {
    const schema = await fetchJSONSchema(getItemName(t, false), "trait", getItemID(t));

    schema._internal = { patternAttributeName: getItemName(t, false) };

    return schema;
  }));

  const type = getPatternServiceType(schemaSet?.workload)
  return { workload: workloadSchema, traits: traitsSchemas, type }
}

export default function LazyPatternServiceForm(props) {
  const [expanded, setExpanded] = React.useState(false);
  const [isLoading, setIsLoading] = React.useState(false);
  const [schemaSet, setSchemaSet] = React.useState({});

  const classes = useStyles();

  async function expand(state) {
    if (!state) {
      setExpanded(false);
      return;
    }

    setExpanded(true);
    setIsLoading(true);

    try {
      // Get the schema sets for the workload
      // const workloadSchema = await fetchJSONSchema(getItemName(props?.schemaSet?.workload), "workload", getItemID(props?.schemaSet?.workload));
      // workloadSchema._internal = { patternAttributeName : getItemName(props?.schemaSet?.workload, false) };

      // // Get the schema sets for the traits
      // const traitsSchemas = await Promise.all(props?.schemaSet?.traits?.map(async t => {
      //   const schema = await fetchJSONSchema(getItemName(t, false), "trait", getItemID(t));

      //   schema._internal = { patternAttributeName : getItemName(t, false) };

      //   return schema;
      // }));

      // console.log({ workloadSchema, traitsSchemas })

      const { workload, traits, type } = await getWorkloadTraitAndType(props?.schemaSet)

      setSchemaSet({
        workload,
        traits,
        type,
      });
    } catch (error) {
      console.error(error)
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <div className={classes.root}>
      <Accordion expanded={expanded} onChange={() => expand(!expanded)}>
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Typography className={classes.heading}>{getItemName(props?.schemaSet?.workload)}</Typography>
        </AccordionSummary>
        <LazyAccordionDetails expanded={expanded}>
          {isLoading ? <CircularProgress /> : <PatternServiceForm {...props} schemaSet={schemaSet} />}
        </LazyAccordionDetails>
      </Accordion>
    </div>
  );
}

function LazyAccordionDetails(props) {
  if (!props.expanded) return <AccordionDetails />

  // LEE: This behavior is more like what we need - https://codesandbox.io/s/upbeat-tesla-uchsb?file=/src/MyAccordion.js
  return <AccordionDetails style={{ height: "50rem", overflow: "auto" }}>{props.children}</AccordionDetails>
}
