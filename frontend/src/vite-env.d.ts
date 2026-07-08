/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

declare module 'tabulator-tables' {
  export class TabulatorFull {
    constructor(element: HTMLElement, options: any)
    destroy(): void
    setData(data: any[]): void
    getData(): any[]
    addRow(data: any): void
    deleteRow(row: any): void
    getSelectedRows(): any[]
    setSort(field: string, dir: string): void
    setFilter(field: string, type: string, value: any): void
    clearFilter(): void
    updateData(data: any[]): void
    redraw(force?: boolean): void
  }
}

declare module 'sql-formatter' {
  export function format(sql: string, options?: any): string
}

declare module '@vue-flow/core' {
  export const VueFlow: any
  export const useVueFlow: any
}

declare module '@vue-flow/background' {
  export const Background: any
}

declare module '@vue-flow/controls' {
  export const Controls: any
}
