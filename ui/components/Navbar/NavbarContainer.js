/* eslint-disable react/prop-types */
import RemoveIcon from "@mui/icons-material/Remove";
import React from "react";
import { useState, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { fetchProviderDetailsThunk } from "@/features/provider";
import { compose } from "redux";
import { navigatorItemsTree } from "./constants";
import { fetchAvailableAdaptersThunk } from "@/features/mesheryComponents/mesheryComponentsSlice";
import extensionPointSchemaValidator from "@/utils/extensionPointSchemaValidator";

// -------------------Helper Functions-------------------------------------------------------------------------------------------------

/**
 *  Fetch the children of Lifecylce manageent category based on the adapters status
 * @param {Object} category
 * @param {import("@/features/mesheryComponents/mesheryComponentsSlice").AdaptersListType} meshAdapters - List of available adapters registered
 * @returns {Object[]} children
 */
const fetchLifecycleChildren = (category, meshAdapters) => {
  const children = [];
  category = category.toLowerCase();
  meshAdapters.forEach((adapter) => {
    let aName = adapter.name.toLowerCase();
    // Manually changing adapter name so that it matches the internal name
    if (aName === "osm") aName = "open service mesh";
    if (category !== aName) {
      return;
    }
    children.push({
      id: adapter.adapter_location,
      icon: <RemoveIcon />,
      href: `/management?adapter=${adapter.adapter_location}`,
      title: `Management - ${adapter.adapter_location}`,
      link: true,
      show: true,
    });
  });
  return children;
};

/**
 * Picks the service mesh icon based on the adapter name
 * @param {string} adapterName
 * @returns {import("react").ReactElement} Icon
 */

const pickIconForServiceMeshes = (adapterName) => {
  adapterName = adapterName.toLowerCase();
  let image = "/static/img/meshery-logo.png";
  let logoIcon = <img src={image} style={{ fontSize: 20 }} />;
  if (adapterName) {
    image = "/static/img/" + adapterName + "-light.svg";
    logoIcon = <img src={image} width="20rem" height="20rem" />;
  }
  return logoIcon;
};

/**
 * Generates a updater function which would update the categories based on the information from provider
 * @param {import("@/features/mesheryComponents/mesheryComponentsSlice").AdaptersListType} meshAdapters
 * @param {import("@/features/provider/providerSlice").Capability[]} providerCapabilities
 * @returns {(categories) => Object[]} updaterFunction - Function that updates the categories
 */
const updateCategoriesMenus = (providerCapabilities, meshAdapters) => (categories) => {
  categories.forEach((cat, ind) => {
    if (cat.id === "Lifecycle") {
      cat.children.forEach((catc, ind1) => {
        const cr = fetchLifecycleChildren(catc.id, meshAdapters);
        const icon = pickIconForServiceMeshes(catc.id);
        categories[ind].children[ind1].icon = icon;
        categories[ind].children[ind1].children = cr;
      });
    }

    if (cat.id === "Configuration") {
      let show = false;
      cat.children?.forEach((ch) => {
        if (ch.id === "Patterns") {
          const idx = providerCapabilities?.findIndex((cap) => cap.feature === "persist-meshery-patterns");
          if (idx != -1) {
            ch.show = true;
            show = true;
          }
        }
      });

      cat.show = show;
    }
  });
  return categories;
};

const updateAdaptersLink = (categories) => {
  categories.forEach((cat, ind) => {
    if (cat.id === "Lifecycle") {
      cat.children.forEach((catc, ind1) => {
        if (
          typeof categories[ind].children[ind1].children[0] !== "undefined" &&
          typeof categories[ind].children[ind1].children[0].href !== "undefined"
        ) {
          const val = true;
          const newhref = `${categories[ind].children[ind1].children[0].href}`;
          categories[ind].children[ind1].link = val;
          categories[ind].children[ind1].href = newhref;
        }
      });
    }
  });
  return categories;
};

// -------------------***********************-------------------------------------------------------------------------------------------------

export const NavbarContainer = (props) => {
  const [isDrawerCollapsed, setIsDrawerCollapsed] = useState(false);
  const dispatch = useDispatch();

  const onDrawerCollapse = () => setIsDrawerCollapsed((prev) => !prev);

  const extensions = useSelector((state) => state.provider.extensions);
  /** @type {import("@/features/provider/providerSlice").Capability} */
  const capabilities = useSelector((state) => state.provider.capabilities);
  /** @type {import("@/features/mesheryComponents/mesheryComponentsSlice").AdaptersListType} */
  const meshAdapters = useSelector((state) => state.mesheryComponents.adapters);
  const [categories, setCategories] = useState([]);

  useEffect(() => {
    // has to be dispatched elsewhere since this data will be used more or have to make sure it is dispatched only once
    dispatch(fetchProviderDetailsThunk());
    dispatch(fetchAvailableAdaptersThunk());
  }, []);

  const updateCategories = compose(updateAdaptersLink, updateCategoriesMenus(capabilities, meshAdapters));

  useEffect(() => {
    setCategories(updateCategories(navigatorItemsTree));
  }, [capabilities, extensions]);

  const extensionsNavigator = extensionPointSchemaValidator("navigator")(extensions.navigator);

  return props.render({ isDrawerCollapsed, extensionsNavigator, categories, onDrawerCollapse });
};
