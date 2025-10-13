import { PanelPlugin } from "@grafana/data";
import { SimpleOptions } from "./types";
import { HoverTrackerPanel } from "./HoverTrackerPanel";

export const plugin = new PanelPlugin<SimpleOptions>(
  HoverTrackerPanel
).setPanelOptions((builder) => {
  return builder
    .addNumberInput({
      path: "historySize",
      name: "History Size",
      description: "Number of hover events to keep in history",
      defaultValue: 100,
    })
    .addBooleanSwitch({
      path: "showTimestamp",
      name: "Show Timestamp",
      description: "Display timestamp for each hover event",
      defaultValue: true,
    })
    .addBooleanSwitch({
      path: "trackOwnPanel",
      name: "Track Own Panel",
      description: "Include hover events from this panel",
      defaultValue: false,
    })
    .addTextInput({
      path: "eventChannel",
      name: "Event Channel",
      description: "Custom event channel name for cross-panel communication",
      defaultValue: "hover-tracker",
    })
    .addTextInput({
      path: "apiEndpoint",
      name: "API Endpoint",
      description:
        "Optional API endpoint to send hover data (metric name, time range, value)",
      defaultValue: "",
      settings: {
        placeholder: "http://127.0.0.1:3001/query_logs",
      },
    });
});
