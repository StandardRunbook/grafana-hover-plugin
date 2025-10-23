import React, { useEffect, useState, useCallback } from "react";
import {
  PanelProps,
  DataHoverEvent,
  LegacyGraphHoverEvent,
  dateTime,
} from "@grafana/data";
import { SimpleOptions, HoverEvent } from "./types";
import { css, cx } from "@emotion/css";
import { useStyles2 } from "@grafana/ui";
import { getAppEvents, config as grafanaConfig, getBackendSrv } from "@grafana/runtime";

interface Props extends PanelProps<SimpleOptions> {}

export const HoverTrackerPanel: React.FC<Props> = (props) => {
  const { options, width, height, id, eventBus } = props;

  // Note: buildPanelRegistry removed - no longer needed
  // Note: buildRefIdToPanelMap removed - no longer needed
  // Panel title comes directly from event.origin._state.title

  const styles = useStyles2(getStyles);
  const [currentHover, setCurrentHover] = useState<HoverEvent | null>(null);
  const [apiLogs, setApiLogs] = useState<string[]>([]);
  const [expandedLogs, setExpandedLogs] = useState<Set<number>>(new Set());
  const [isLoadingLogs, setIsLoadingLogs] = useState(false);

  const toggleLogExpansion = (index: number) => {
    setExpandedLogs((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(index)) {
        newSet.delete(index);
      } else {
        newSet.add(index);
      }
      return newSet;
    });
  };

  // Function to send metric data to backend plugin
  const sendToAPI = useCallback(
    async (event: HoverEvent, eventOrigin?: any) => {
      try {
        // Clear old logs and show loading state
        setIsLoadingLogs(true);
        setApiLogs([]);

        const metricData = event.metricData;

        // Prepare the payload to match log analysis server spec
        // API expects: metric_name, start_time, end_time (ISO 8601 with Z suffix)
        const endTime = new Date();
        const startTime = new Date(Date.now() - options.timeWindowMs);

        // Extract metric name with better fallbacks

        // metric_name: the series/metric name from the data
        const metricName =
          metricData?.seriesName || metricData?.fieldName || "Unknown Series";

        // graph_name: the panel title
        const graphName = event.panelTitle || "Unknown Panel";

        // dashboard_name: from the dashboard state
        // Based on the structure: event.origin._eventsOrigin._state is the dashboard state
        const dashboardState = eventOrigin?._eventsOrigin?._state;
        const dashboardName = dashboardState?.title || "Unknown Dashboard";

        // org: Get the Grafana organization ID and convert to string
        const orgId = String(grafanaConfig.bootData?.user?.orgId || 1);

        const payload = {
          org: orgId,
          dashboard: dashboardName,
          panel_title: graphName,
          metric_name: metricName,
          start_time: startTime.toISOString(),
          end_time: endTime.toISOString(),
        };

        // Call backend plugin resource endpoint
        const result = await getBackendSrv().post(
          `/api/plugins/hover-hover-panel/resources/query_logs`,
          payload
        );

        // Parse log_groups response from log analysis server
        // Expected format: { log_groups: [{ representative_logs: [...], relative_change: number }] }
        if (result && Array.isArray(result.log_groups)) {
          const formattedLogs: string[] = [];
          const MAX_LOGS = options.maxLogs;
          const MAX_LOG_LENGTH = options.maxLogLength;
          let totalLogs = 0;

          result.log_groups.forEach((group: any, index: number) => {
            if (totalLogs >= MAX_LOGS) {
              return; // Stop processing if we hit the limit
            }

            const change = group.relative_change;
            const changeSymbol = change > 0 ? "↑" : change < 0 ? "↓" : "→";
            const changeColor =
              change > 50
                ? "🔴"
                : change > 10
                ? "🟠"
                : change < -10
                ? "🟢"
                : "⚪";

            formattedLogs.push(
              `${changeColor} ${changeSymbol} ${change.toFixed(
                1
              )}% change from baseline`
            );
            totalLogs++;

            if (Array.isArray(group.representative_logs)) {
              group.representative_logs.forEach((log: string) => {
                if (totalLogs >= MAX_LOGS) {
                  return;
                }
                // Truncate extremely long logs but allow expansion
                const truncatedLog =
                  log.length > MAX_LOG_LENGTH
                    ? log.substring(0, MAX_LOG_LENGTH) + "... [truncated]"
                    : log;
                formattedLogs.push(`  ${truncatedLog}`);
                totalLogs++;
              });
            }
          });

          if (totalLogs >= MAX_LOGS) {
            formattedLogs.push(
              `... [${
                result.log_groups.length - formattedLogs.length
              } more logs truncated for performance]`
            );
          }

          setApiLogs(formattedLogs);
          setIsLoadingLogs(false);
        } else {
          setApiLogs([]);
          setIsLoadingLogs(false);
        }
      } catch (error) {
        setApiLogs([`Error: ${error}`]);
        setIsLoadingLogs(false);
      }
    },
    [
      options.timeWindowMs,
      options.maxLogs,
      options.maxLogLength,
    ]
  );

  useEffect(() => {
    // Subscribe to Grafana hover events
    const appEvents = getAppEvents();

    // Method 1: Use eventBus from props (recommended approach)
    const dataHoverSub = eventBus
      ?.getStream(DataHoverEvent)
      .subscribe((event) => {

        // Extract event origin for panel info
        const eventOrigin = (event as any)?.origin;
        const payloadOrigin = (event.payload as any)?.origin;

        // Extract data from the payload
        const payload = event.payload as any;
        const point = payload?.point || {};
        const data = payload?.data; // This should be the DataFrame

        // Get series name from the data frame
        let seriesName = "Unknown Series";
        let fieldName = "";
        let value: any = undefined;
        let formattedValue: string | undefined = undefined;
        let time: number | undefined = undefined;

        // Extract from DataFrame if available
        if (data && data.fields) {
          // Find the field that was hovered
          const fieldIndex = payload?.columnIndex ?? payload?.fieldIndex;
          const rowIndex = payload?.rowIndex ?? payload?.dataIndex;

          // If columnIndex is undefined, find the first non-time field
          let targetFieldIndex = fieldIndex;
          if (targetFieldIndex === undefined) {
            targetFieldIndex = data.fields.findIndex(
              (f: any) => f.type !== "time"
            );
          }

          if (targetFieldIndex !== undefined && data.fields[targetFieldIndex]) {
            const field = data.fields[targetFieldIndex];
            fieldName = field.name || "";

            // Extract series name with priority: displayName > name > labels.__name__ > labels
            seriesName =
              field.config?.displayName ||
              field.name ||
              data.name ||
              "Unknown Series";

            // Get the value at the hovered row
            if (
              rowIndex !== undefined &&
              field.values &&
              field.values.length > rowIndex
            ) {
              value = field.values[rowIndex];
              formattedValue = field.display
                ? field.display(value).text
                : String(value);
            }
          }

          // Try to get time from the first field (usually time field)
          if (
            rowIndex !== undefined &&
            data.fields[0]?.values?.length > rowIndex
          ) {
            time = data.fields[0].values[rowIndex];
          }
        }

        // Fallback to point data
        if (value === undefined && point?.value !== undefined) {
          value = point.value;
        }
        if (seriesName === "Unknown Series" && payload?.dataId) {
          seriesName = payload.dataId;
        }

        // Extract panel info from available data
        const refId = data?.refId;

        // Try multiple sources for panel title and ID:
        // 1. Check if payload has panel info
        // 2. Use data frame name or refId
        // 3. Fallback to "Unknown Panel"

        let panelTitle = "Unknown Panel";
        let panelId = "unknown";

        // PRIMARY: Get title from event.origin._eventsOrigin._state.title
        if (eventOrigin?._eventsOrigin?._state?.title) {
          panelTitle = eventOrigin._eventsOrigin._state.title;
        }
        // FALLBACK 1: Try direct _state.title (in case structure varies)
        else if (eventOrigin?._state?.title) {
          panelTitle = eventOrigin._state.title;
        }
        // FALLBACK 2: Try payload origin
        else if (payloadOrigin?._eventsOrigin?._state?.title) {
          panelTitle = payloadOrigin._eventsOrigin._state.title;
        } else if (payloadOrigin?._state?.title) {
          panelTitle = payloadOrigin._state.title;
        }
        // FALLBACK 3: Try data frame metadata
        else if (data?.meta?.custom?.title) {
          panelTitle = data.meta.custom.title;
        }
        // FALLBACK 4: Use series name or refId
        else if (data?.name) {
          panelTitle = `Panel: ${data.name}`;
        } else if (refId) {
          panelTitle = `Query ${refId}`;
        }

        // Get panel ID from refId or origin
        if (refId) {
          panelId = refId;
        } else if (eventOrigin?.panelId) {
          panelId = eventOrigin.panelId;
        } else if (payloadOrigin?.panelId) {
          panelId = payloadOrigin.panelId;
        }

        const newEvent: HoverEvent = {
          panelId,
          panelTitle,
          timestamp: Date.now(),
          x: point?.time || point?.x || 0,
          y: value || point?.y || 0,
          elementType: "data-hover",
          metricData: {
            seriesName,
            fieldName,
            value,
            formattedValue,
            time,
            fieldIndex: payload?.columnIndex ?? payload?.fieldIndex,
            dataIndex: payload?.rowIndex ?? payload?.dataIndex,
            point: point,
            dataFrame: data,
          },
        };


        // Update current hover widget
        setCurrentHover(newEvent);

        // Send to API if configured (pass eventOrigin for dashboard/org info)
        sendToAPI(newEvent, eventOrigin);
      });

    // Intentionally NOT subscribing to DataHoverClearEvent
    // We want to persist the hover data until a new hover event occurs

    // Method 2: Fallback using getAppEvents() for legacy support
    const legacyHoverSub = appEvents.subscribe(
      LegacyGraphHoverEvent,
      (event) => {
        const newEvent: HoverEvent = {
          panelId: "legacy-graph-hover",
          panelTitle: "Legacy Graph Hover",
          timestamp: Date.now(),
          x: (event.payload as any)?.pos?.x || 0,
          y: (event.payload as any)?.pos?.y || 0,
          elementType: "legacy-graph-hover",
          metricData: {
            seriesName: (event.payload as any)?.series?.name || "Legacy Series",
            value: (event.payload as any)?.value,
          },
        };

        // Update current hover widget
        setCurrentHover(newEvent);

        // Send to API if configured
        sendToAPI(newEvent);
      }
    );

    return () => {
      dataHoverSub?.unsubscribe();
      legacyHoverSub?.unsubscribe();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id, sendToAPI, eventBus]);

  return (
    <>
      {/* Main Panel */}
      <div className={cx(styles.wrapper)} style={{ width, height }}>
        {/* Dynamic Current Hover Widget */}
        {currentHover && (
          <div className={styles.currentHoverWidget}>
            <div className={styles.widgetContent}>
              {/* Compact metadata at the top */}
              <div className={styles.widgetCompactHeader}>
                {currentHover.metricData?.time && (
                  <>
                    <span className={styles.widgetCompactTime}>
                      {dateTime(currentHover.metricData.time).format(
                        "YYYY-MM-DD HH:mm:ss"
                      )}
                    </span>
                    <span className={styles.widgetCompactSeparator}>•</span>
                  </>
                )}
                <span className={styles.widgetCompactPanel}>
                  {currentHover.panelTitle}
                </span>
                <span className={styles.widgetCompactSeparator}>•</span>
                <span className={styles.widgetCompactSeries}>
                  {currentHover.metricData?.seriesName || "No Series"}
                </span>
                <span className={styles.widgetCompactSeparator}>•</span>
                <span className={styles.widgetCompactValue}>
                  {currentHover.metricData?.formattedValue ||
                    currentHover.metricData?.value ||
                    "—"}
                </span>
              </div>

              {/* API Response Logs Section - takes up most of the space */}
              {isLoadingLogs ? (
                <div className={styles.widgetLoading}>Loading logs...</div>
              ) : apiLogs.length > 0 ? (
                <div className={styles.widgetLogsSection}>
                  <div className={styles.widgetLogsList}>
                    {apiLogs.map((log, index) => {
                      const isExpanded = expandedLogs.has(index);
                      const maxLength = options.logTruncateLength;
                      const isLong = log.length > maxLength;

                      return (
                        <div
                          key={index}
                          className={styles.widgetLogItem}
                          onClick={() => isLong && toggleLogExpansion(index)}
                          style={{
                            cursor: isLong ? "pointer" : "default",
                          }}
                          title={
                            isLong ? "Click to expand/collapse" : undefined
                          }
                        >
                          {isLong && (
                            <span className={styles.widgetLogToggle}>
                              {isExpanded ? "▼" : "▶"}
                            </span>
                          )}
                          <span className={styles.widgetLogText}>
                            {isExpanded || !isLong
                              ? log
                              : log.substring(0, maxLength) +
                                "... (click to expand)"}
                          </span>
                        </div>
                      );
                    })}
                  </div>
                </div>
              ) : (
                <div className={styles.widgetLoading}>No logs found</div>
              )}
            </div>
          </div>
        )}

        {!currentHover && (
          <div className={styles.noCurrentHover}>
            <div className={styles.noHoverIcon}>👆</div>
            <div className={styles.noHoverText}>
              Hover over a panel to see live data
            </div>
          </div>
        )}
      </div>
    </>
  );
};

const getStyles = () => {
  return {
    wrapper: css`
      position: relative;
      display: flex;
      flex-direction: column;
      height: 100%;
      overflow: hidden;
    `,
    header: css`
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 8px 16px;
      border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    `,
    title: css`
      margin: 0;
      font-size: 16px;
      font-weight: 500;
    `,
    count: css`
      font-size: 12px;
      opacity: 0.7;
    `,
    currentHoverWidget: css`
      display: flex;
      flex-direction: column;
      overflow: hidden;
      flex: 1;
      height: 100%;
    `,
    widgetHeader: css`
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;
      padding-bottom: 8px;
      border-bottom: 1px solid rgba(115, 191, 105, 0.3);
    `,
    widgetTitle: css`
      font-size: 11px;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 1px;
      color: rgba(115, 191, 105, 1);
    `,
    widgetLive: css`
      font-size: 10px;
      font-weight: 600;
      color: rgba(115, 191, 105, 1);
      animation: blink 1.5s ease-in-out infinite;

      @keyframes blink {
        0%,
        100% {
          opacity: 1;
        }
        50% {
          opacity: 0.4;
        }
      }
    `,
    widgetContent: css`
      display: flex;
      flex-direction: column;
      gap: 8px;
      flex: 1;
      min-height: 0;
    `,
    widgetCompactHeader: css`
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 4px 0;
      font-size: 11px;
      flex-wrap: wrap;
      border-bottom: 1px solid rgba(255, 255, 255, 0.1);
      padding-bottom: 6px;
      margin-bottom: 4px;
    `,
    widgetCompactPanel: css`
      font-weight: 600;
      color: rgba(255, 255, 255, 0.9);
    `,
    widgetCompactSeparator: css`
      color: rgba(255, 255, 255, 0.3);
    `,
    widgetCompactSeries: css`
      color: rgba(255, 255, 255, 0.7);
    `,
    widgetCompactValue: css`
      font-weight: 700;
      color: rgba(115, 191, 105, 1);
      font-family: "Roboto Mono", monospace;
    `,
    widgetCompactTime: css`
      color: rgba(255, 255, 255, 0.6);
      font-family: "Roboto Mono", monospace;
      font-size: 11px;
    `,
    widgetLoading: css`
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 24px;
      color: rgba(255, 255, 255, 0.5);
      font-style: italic;
    `,
    widgetMainValue: css`
      display: flex;
      flex-direction: column;
      gap: 4px;
    `,
    widgetSeriesName: css`
      font-size: 13px;
      font-weight: 600;
      color: rgba(255, 255, 255, 0.9);
      letter-spacing: 0.3px;
    `,
    widgetValue: css`
      font-size: 32px;
      font-weight: 700;
      color: rgba(115, 191, 105, 1);
      line-height: 1.2;
      font-family: "Roboto Mono", monospace;
    `,
    widgetMetadata: css`
      display: flex;
      flex-direction: column;
      gap: 6px;
      padding-top: 8px;
      border-top: 1px solid rgba(255, 255, 255, 0.1);
    `,
    widgetMetaRow: css`
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 8px;
    `,
    widgetMetaLabel: css`
      font-size: 11px;
      font-weight: 600;
      color: rgba(255, 255, 255, 0.6);
      text-transform: uppercase;
      letter-spacing: 0.5px;
    `,
    widgetMetaValue: css`
      font-size: 12px;
      font-weight: 500;
      color: rgba(255, 255, 255, 0.9);
      text-align: right;
    `,
    noCurrentHover: css`
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 32px 16px;
      margin: 12px;
      background: rgba(255, 255, 255, 0.03);
      border: 2px dashed rgba(255, 255, 255, 0.1);
      border-radius: 8px;
      opacity: 0.6;
    `,
    noHoverIcon: css`
      font-size: 48px;
      margin-bottom: 12px;
      opacity: 0.5;
    `,
    noHoverText: css`
      font-size: 14px;
      color: rgba(255, 255, 255, 0.6);
      text-align: center;
    `,
    widgetLogsSection: css`
      display: flex;
      flex-direction: column;
      flex: 1;
      min-height: 0;
      padding-top: 8px;
    `,
    widgetLogsSectionHeader: css`
      font-size: 11px;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 0.5px;
      color: rgba(115, 191, 105, 0.8);
      margin-bottom: 8px;
    `,
    widgetLogsList: css`
      display: flex;
      flex-direction: column;
      gap: 6px;
      flex: 1;
      overflow-y: auto;
      padding-right: 4px;

      /* Custom scrollbar */
      &::-webkit-scrollbar {
        width: 6px;
      }
      &::-webkit-scrollbar-track {
        background: rgba(0, 0, 0, 0.2);
        border-radius: 3px;
      }
      &::-webkit-scrollbar-thumb {
        background: rgba(115, 191, 105, 0.4);
        border-radius: 3px;
      }
      &::-webkit-scrollbar-thumb:hover {
        background: rgba(115, 191, 105, 0.6);
      }
    `,
    widgetLogItem: css`
      display: flex;
      gap: 6px;
      align-items: flex-start;
      padding: 4px 6px;
      margin: 0;
      font-family: "Roboto Mono", "Courier New", monospace;
      font-size: 11px;
      line-height: 1.4;
      transition: all 0.15s ease;
      user-select: text;
      border-left: 2px solid transparent;

      &:hover {
        background: rgba(255, 255, 255, 0.03);
      }

      &[style*="cursor: pointer"] {
        &:hover {
          background: rgba(115, 191, 105, 0.08);
          border-left-color: rgba(115, 191, 105, 0.6);
        }
      }
    `,
    widgetLogToggle: css`
      color: rgba(115, 191, 105, 0.7);
      font-size: 9px;
      flex-shrink: 0;
      width: 12px;
      margin-top: 3px;
    `,
    widgetLogText: css`
      color: rgba(255, 255, 255, 0.85);
      word-break: break-word;
      flex: 1;
    `,
    widgetLogDivider: css`
      height: 1px;
      background: linear-gradient(
        to right,
        transparent,
        rgba(115, 191, 105, 0.3) 20%,
        rgba(115, 191, 105, 0.3) 80%,
        transparent
      );
      margin: 12px 0;
    `,
    eventList: css`
      flex: 1;
      overflow-y: auto;
      padding: 8px;
    `,
    event: css`
      margin-bottom: 8px;
      padding: 8px 12px;
      background: rgba(255, 255, 255, 0.05);
      border-radius: 4px;
      border-left: 3px solid rgba(115, 191, 105, 0.7);
      transition: all 0.2s ease;

      &:hover {
        background: rgba(255, 255, 255, 0.08);
        border-left-color: rgba(115, 191, 105, 1);
      }
    `,
    eventHeader: css`
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 4px;
    `,
    panelName: css`
      font-weight: 500;
      font-size: 14px;
    `,
    timestamp: css`
      font-size: 11px;
      opacity: 0.6;
    `,
    eventDetails: css`
      display: flex;
      flex-direction: column;
      gap: 8px;
      font-size: 12px;
      opacity: 0.9;
    `,
    detail: css`
      display: inline-block;
    `,
    metricInfo: css`
      display: flex;
      gap: 8px;
      align-items: flex-start;
    `,
    metricLabel: css`
      font-weight: 600;
      color: rgba(255, 255, 255, 0.7);
      min-width: 60px;
    `,
    metricValue: css`
      color: rgba(255, 255, 255, 0.9);
      word-break: break-word;
      flex: 1;
    `,
    noEvents: css`
      text-align: center;
      padding: 32px;
      opacity: 0.5;
    `,
    debugInfo: css`
      margin-top: 8px;
      border-top: 1px solid rgba(255, 255, 255, 0.1);
      padding-top: 8px;
    `,
    debugSummary: css`
      cursor: pointer;
      font-size: 11px;
      color: rgba(255, 255, 255, 0.6);
      &:hover {
        color: rgba(255, 255, 255, 0.8);
      }
    `,
    debugPre: css`
      background: rgba(0, 0, 0, 0.3);
      padding: 8px;
      border-radius: 4px;
      font-size: 10px;
      overflow-x: auto;
      margin: 4px 0 0 0;
      max-height: 200px;
      overflow-y: auto;
    `,
    customTooltip: css`
      position: fixed;
      z-index: 9999;
      background: rgba(0, 0, 0, 0.9);
      color: white;
      padding: 12px;
      border-radius: 6px;
      border: 1px solid rgba(115, 191, 105, 0.7);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
      max-width: 300px;
      font-size: 12px;
      pointer-events: none;
      backdrop-filter: blur(4px);
    `,
    tooltipHeader: css`
      margin-bottom: 8px;
      padding-bottom: 6px;
      border-bottom: 1px solid rgba(115, 191, 105, 0.3);
      font-size: 14px;
      color: rgba(115, 191, 105, 1);
    `,
    tooltipContent: css`
      display: flex;
      flex-direction: column;
      gap: 4px;
    `,
    tooltipRow: css`
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      gap: 8px;
    `,
    tooltipLabel: css`
      font-weight: 600;
      color: rgba(255, 255, 255, 0.7);
      min-width: 50px;
      flex-shrink: 0;
    `,
    tooltipValue: css`
      color: rgba(255, 255, 255, 0.9);
      word-break: break-word;
      text-align: right;
      flex: 1;
    `,
    tooltipDivider: css`
      margin: 8px 0;
      height: 1px;
      background: rgba(115, 191, 105, 0.3);
    `,
    tooltipSection: css`
      margin-top: 8px;
    `,
    tooltipSectionHeader: css`
      font-weight: 600;
      color: rgba(115, 191, 105, 0.8);
      font-size: 11px;
      text-transform: uppercase;
      margin-bottom: 6px;
      letter-spacing: 0.5px;
    `,
    tooltipLoading: css`
      color: rgba(255, 255, 255, 0.6);
      font-style: italic;
      font-size: 11px;
    `,
    tooltipDataFrame: css`
      margin-bottom: 6px;
      padding-bottom: 4px;
      border-bottom: 1px solid rgba(255, 255, 255, 0.1);
      &:last-child {
        border-bottom: none;
        margin-bottom: 0;
      }
    `,
    tooltipFields: css`
      margin-top: 4px;
    `,
    tooltipMuted: css`
      color: rgba(255, 255, 255, 0.5);
      font-size: 11px;
      font-style: italic;
    `,
  };
};