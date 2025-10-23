export interface SimpleOptions {
    timeWindowMs: number;
    maxLogs: number;
    maxLogLength: number;
    logTruncateLength: number;
}
export interface HoverEvent {
    panelId: string;
    panelTitle: string;
    timestamp: number;
    x: number;
    y: number;
    elementType?: string;
    elementData?: any;
    metricData?: {
        seriesName?: string;
        fieldName?: string;
        value?: number | string;
        formattedValue?: string;
        labels?: Record<string, string>;
        time?: number;
        timeRange?: any;
        dataFrame?: any;
        fieldIndex?: number;
        dataIndex?: number;
        point?: any;
        series?: any;
    };
}
