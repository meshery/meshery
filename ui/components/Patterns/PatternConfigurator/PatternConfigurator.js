import React, { useState, useEffect } from "react";
import {
  AppBar,
  FormControl,
  IconButton,
  MenuItem,
  TextField,
  Toolbar,
  Tooltip,
  Autocomplete,
  Select,
} from "@mui/material";
import { styled } from "@mui/system";
import DeleteIcon from "@mui/icons-material/Delete";
import FileCopyIcon from "@mui/icons-material/FileCopy";
import ListAltIcon from "@mui/icons-material/ListAlt";
import SaveIcon from "@mui/icons-material/Save";
import {
  getHumanReadablePatternServiceName,
  getPatternServiceName,
  isEmptyObj,
  NameToIcon,
  getMeshProperties,
} from "./helpers";

const PattrenBar = styled(AppBar)(({theme}) => ({
  marginBottom: theme.spacing(2),
  backgroundColor: theme.palette.white,
  borderRadius: theme.spacing(1),
}));

const PatternConfigurator = ({ pattern, onSubmit, show: setSelectedPattern }) => {
  // const { workloadTraitSet, meshWorkloads } = useContext(SchemaContext);
  // const [workloadTraitsSet, setWorkloadTraitsSet] = useState(workloadTraitSet);
  const [selectedMeshType, setSelectedMeshType] = useState("core");
  const [selectedVersion, setSelectedVersion] = useState("");
  const [selectedVersionMesh, setSelectedVersionMesh] = useState();
  const [briefCrsInformations, setBriefCrsInformations] = useState(null);
  const [activeForm, setActiveForm] = useState();
  const [activeCR, setActiveCR] = useState({});
  const [viewType, setViewType] = useState("list");

  useEffect(() => {
    // core is not versioned
    if (selectedMeshType == "core") {
      setSelectedVersionMesh(null);
    } else {
      const meshVersionsWithDetails = groupWlByVersion();
      setSelectedVersionMesh(meshVersionsWithDetails);
    }
    setViewType("list");
    setActiveForm(null);
  }, [selectedMeshType]);

  useEffect(() => {
    if (selectedVersionMesh) {
      setSelectedVersion(Object.keys(selectedVersionMesh).sort().reverse()[0]);
    }
  }, [selectedVersionMesh]);

  useEffect(() => {
    if (selectedVersion) {
      const crsBriefs = getFormBriefInformationKeys();
      setBriefCrsInformations(crsBriefs);
      setActivePatternWithRefinedSchema(
        selectedVersionMesh?.[selectedVersion]?.sort((a, b) =>
          getPatternServiceName(a.workload) < getPatternServiceName(b.workload) ? -1 : 1
        )[0]
      );
      setActiveCR(crsBriefs[0]);
    }
  }, [selectedVersion]);

  const groupWlByVersion = () => {
    // const mfw = meshWorkloads[selectedMeshType];
    // return mfw ? groupWorkloadByVersion(mfw) : {};
  };

  const resetSelectedPattern = () => {
    return { show: false, pattern: null };
  };

  const handleMeshSelection = (event) => setSelectedMeshType(event.target.value);
  
  const getMeshOptions = () => {
    return meshWorkloads ? Object.keys(meshWorkloads) : [];
  };

  const handleVersionChange = (_, value) => {
    setSelectedVersion(value);
  };

  const handleSubmitFinalPattern = (yaml, id, name, action) => {
    onSubmit({
      data: yaml,
      id: id,
      name: name,
      type: action,
    });
    setSelectedPattern(resetSelectedPattern());
  };

  const setActivePatternWithRefinedSchema = async (schema) => {
    // const refinedSchema = await getWorkloadTraitAndType(schema);
    // setActiveForm(refinedSchema);
  };

  /**
   * get keys and mapping to the correct icons and colors
   * for all the CRDs available in any SM
   * @returns {{name: String, icon: React.ReactElement, color: String}}
   */
  function getFormBriefInformationKeys() {
    if (selectedMeshType === "core" || selectedMeshType === "kubernetes") {
      return meshWorkloads[selectedMeshType].map((mwl) => {
        const name = mwl?.workload?.metadata?.["display.ui.meshery.io/name"];
        return {
          name,
          icon: <NameToIcon name={name.split(".")[0]} color={getMeshProperties(selectedMeshType).color} />,
          readableName: getHumanReadablePatternServiceName(mwl?.workload),
        };
      });
    }
    return selectedVersionMesh?.[selectedVersion]
      ?.sort((a, b) => (getPatternServiceName(a.workload) < getPatternServiceName(b.workload) ? -1 : 1))
      .map((item) => {
        const name = item.workload?.oam_definition?.metadata?.name;
        return {
          name,
          icon: <NameToIcon name={name.split(".")[0]} color={getMeshProperties(selectedMeshType).color} />,
          readableName: getHumanReadablePatternServiceName(item?.workload),
        };
      });
  }

  const toggleView = () => {
    if (viewType == "list") {
      if (isEmptyObj(activeForm)) {
        // core resources are handled sepaeratrly since they are not versioned
        setBriefCrsInformations(getFormBriefInformationKeys());
        if (selectedMeshType === "core" || selectedMeshType === "kubernetes") {
          setActivePatternWithRefinedSchema(
            meshWorkloads[selectedMeshType]?.sort((a, b) =>
              getPatternServiceName(a.workload) < getPatternServiceName(b.workload) ? -1 : 1
            )[0]
          );
        } else {
          setActivePatternWithRefinedSchema(
            selectedVersionMesh?.[selectedVersion]?.sort((a, b) =>
              getPatternServiceName(a.workload) < getPatternServiceName(b.workload) ? -1 : 1
            )[0]
          );
        }
      }
      setViewType("form");
    } else {
      setViewType("list");
    }
  };

  return (
    <>
      <PattrenBar position="static" elevation={0}>
        <Toolbar
          sx={{
            display: "flex",
            justifyContent: "space-between",
          }}
        >
          <FormControl
            sx={{
              width: "60px",
              minWidth: "60px",
              maxWidth: "60px",
              marginRight: 8,
            }}
          >
            <Select
              id="service-mesh-selector"
              labelId="service-mesh-selector"
              value={selectedMeshType}
              onChange={handleMeshSelection}
              fullWidth
            >
              {/* {getMeshOptions().map((meshItem) => {
                const meshDetails = getMeshProperties(meshItem);
                return (
                  <MenuItem
                    key={meshDetails.name}
                    value={meshDetails.name}
                    sx={{
                      padding: "5px 0",
                      justifyContent: "center",
                    }}
                  >
                    <img src={meshDetails.image} height="32px" />
                  </MenuItem>
                );
              })} */}
            </Select>
          </FormControl>
          {selectedVersion && selectedVersionMesh && (
            <Autocomplete
              options={Object.keys(selectedVersionMesh).sort().reverse()}
              renderInput={(params) => <TextField {...params} variant="outlined" label="Version" />}
              value={selectedVersion}
              onChange={handleVersionChange}
              disableClearable
              sx={{
                width: "120px",
                minWidth: "120px",
                maxWidth: 120,
              }}
            />
          )}
          {viewType === "form" && briefCrsInformations && briefCrsInformations.length > 0 && (
            <Autocomplete
              disableClearable
              value={activeCR}
              options={briefCrsInformations}
              getOptionLabel={(option) => option.readableName}
              onChange={(_, newVal) => {
                setActiveCR(newVal);
              }}
              renderOption={(option) => {
                return (
                  <>
                    <IconButton color="primary">{option.icon}</IconButton>
                    {option.name}
                  </>
                );
              }}
              renderInput={(params) => (
                <TextField {...params} variant="outlined" label="Configure" placeholder={selectedMeshType} />
              )}
              sx={{
                width: 250,
                marginLeft: 16,
                marginRight: "auto",
                padding: 0,
              }}
            />
          )}

          <div>
            <Tooltip title="Save Pattern as New File">
              <IconButton
                aria-label="Save"
                color="primary"
                // onClick={() =>
                //   handleSubmitFinalPattern(yaml, "", getRandomName(), "upload")
                // }
              >
                <FileCopyIcon />
              </IconButton>
            </Tooltip>
            <Tooltip title="Update Pattern">
              <IconButton
                aria-label="Update"
                color="primary"
                // onClick={() =>
                //   handleSubmitFinalPattern(
                //     yaml,
                //     pattern.id,
                //     pattern.name,
                //     "update"
                //   )
                // }
              >
                <SaveIcon />
              </IconButton>
            </Tooltip>
            <Tooltip title="Delete Pattern">
              <IconButton
                aria-label="Delete"
                color="primary"
                // onClick={() =>
                //   handleSubmitFinalPattern(
                //     yaml,
                //     pattern.id,
                //     pattern.name,
                //     "delete"
                //   )
                // }
              >
                <DeleteIcon />
              </IconButton>
            </Tooltip>
            <Tooltip title="Toggle View">
              <IconButton color="primary" onClick={toggleView}>
                <ListAltIcon />
              </IconButton>
            </Tooltip>
          </div>
        </Toolbar>
      </PattrenBar>
    </>
  );
};

export default PatternConfigurator;

