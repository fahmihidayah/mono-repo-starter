import type { ColumnDef } from "@tanstack/react-table";
import { createDateColumn, type DateColumnConfig } from "./date-column";
import { createTextColumn, type TextColumnConfig } from "./text-column";
import { createActionColumn, type ActionColumnConfig } from "./action-column";

interface CreateColumnProps<T> {
    columnConfig : (string | TextColumnConfig<T> | DateColumnConfig<T>) []
    actionColumnConfig : ActionColumnConfig<T>

}

export default function createColumn<T>(props : CreateColumnProps<T>) : ColumnDef<T>[] {
    
    const columns : ColumnDef<T>[] = []

    props.columnConfig.forEach(e => {
        if(typeof e === "string") {
            // Capitalize first letter for header
            const header = e.charAt(0).toUpperCase() + e.slice(1);
            columns.push(createTextColumn<T>({
                accessorKey : e as any,
                header : header,
                type : "text"
            }))
        }
        else if (e.type === "text") {
           columns.push(createTextColumn(e))
        }
        else if(e.type === 'date') {
            columns.push(createDateColumn(e))
        }

    })

    

    return [
        ... columns,
        createActionColumn(props.actionColumnConfig)
    ]

}