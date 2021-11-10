import { DataGrid } from '@mui/x-data-grid';
import { styled } from "@mui/material/styles";
import { useTheme } from "@mui/system";

const CustomDataTableWrapper = styled(DataGrid)(({theme})=>({
    autoHeight: true,
}))

export const CustomDataTable = ({rows, columns, options}) =>  { const theme = useTheme();
    return (
      <CustomDataTableWrapper>
        <DataGrid
            rows={rows}
            columns={columns}
            options={options}
            pageSize={25}
            rowsPerPageOptions={[6]}
            checkboxSelection
        />
      </CustomDataTableWrapper>
    );
  }

export default CustomDataTableWrapper;