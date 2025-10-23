import { PanelPlugin } from "@grafana/data";
import { SimpleOptions } from "./types";
import { HoverTrackerPanel } from "./HoverTrackerPanel";
import { addStandardOptions } from "./HoverTrackerEditor";

export const plugin = new PanelPlugin<SimpleOptions>(
  HoverTrackerPanel
).setPanelOptions((builder) => {
  addStandardOptions(builder);
  return builder;
});


