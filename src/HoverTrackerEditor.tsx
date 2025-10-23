import { PanelOptionsEditorBuilder } from "@grafana/data";
import { SimpleOptions } from "./types";

export const addStandardOptions = (
  builder: PanelOptionsEditorBuilder<SimpleOptions>
) => {
  builder
    .addNumberInput({
      path: "timeWindowMs",
      name: "Time Window (ms)",
      description:
        "Time window for log queries in milliseconds (default: 1 hour)",
      defaultValue: 3600000,
      settings: {
        min: 1000,
        step: 1000,
      },
    })
    .addNumberInput({
      path: "maxLogs",
      name: "Max Logs",
      description: "Maximum number of log entries to display",
      defaultValue: 500,
      settings: {
        min: 1,
        step: 1,
      },
    })
    .addNumberInput({
      path: "maxLogLength",
      name: "Max Log Length",
      description: "Maximum length of individual log entry",
      defaultValue: 10000,
      settings: {
        min: 100,
        step: 100,
      },
    })
    .addNumberInput({
      path: "logTruncateLength",
      name: "Log Truncate Length",
      description: "Character count threshold for expandable logs",
      defaultValue: 120,
      settings: {
        min: 50,
        step: 10,
      },
    });
};


