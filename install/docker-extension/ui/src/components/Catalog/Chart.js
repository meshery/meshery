import BBChart from "../Chart/BBChart";
import { CHART_COLORS} from "../utils/constants";
import { bar } from "billboard.js";
import { ChartDiv } from "./style";

const capitalize = (str) => {
    return str.charAt(0).toUpperCase() + str.slice(1);
}

const groupItemsByTypeCount = data => {
  const dataByTypeCount = {};
  data.forEach(item => {
    const type = item?.catalog_data?.type;

    if (type) {
      const label = type.charAt(0).toUpperCase() + type.slice(1);
      const count = dataByTypeCount[label] || 0;
      dataByTypeCount[label] = count + 1;
    }
  });
  return dataByTypeCount;
};

const dataToColors = data => {
  const columns = data.map(item => item[0]);
  const colors = {};
  let colorIdx = 0;

  columns.forEach(col => {
    if (colorIdx >= CHART_COLORS.length) {
      colorIdx = 0;
    }
    colors[col] = CHART_COLORS[colorIdx];
    colorIdx += 1;
  });

  return colors;
};

const CatalogChart = ({filter, pattern, isTheme: isDarkTheme }) => {
  const groupedDataToArray = (groupedData, topics) => {
    return topics.map(topic => groupedData[topic] || 0);
  };

  const filtersByType = groupItemsByTypeCount(filter ? filter.filters : []);
  const patternsByType = {}
  pattern?.category_count?.forEach(e => {
    patternsByType[capitalize(e.type)] = e.count
  })

  const topicsList = pattern?.category_count?.map(e => ({label: capitalize(e.type), value: e.type }))

  let topics = new Set([
    ...topicsList.map(({ label }) => label),
    ...Object.keys(filtersByType),
    ...Object.keys(patternsByType)
  ]);
  topics = [...topics];

  const data = [
    ["Filters", ...groupedDataToArray(filtersByType, topics)],
    ["Designs", ...groupedDataToArray(patternsByType, topics)]
  ];

  const options = {
    data: {
      x: "x",
      columns: [["x", ...topics], ...data],
      groups: [["Filters", "Designs"]],
      type: bar(),
      colors: dataToColors(data)
    },
    axis: {
      x: {
        type: "category"
      }
    },
  };

  return <ChartDiv style={{ background:  isDarkTheme ? '#666A75' : '#D7DADE' }}>{<BBChart options={options} />}</ChartDiv>
}

export default CatalogChart;