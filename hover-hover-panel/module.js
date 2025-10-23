define(["@grafana/data","react","@emotion/css","@grafana/ui","@grafana/runtime"],(e,t,o,a,i)=>(()=>{"use strict";var n={7:e=>{e.exports=a},89:e=>{e.exports=o},531:e=>{e.exports=i},781:t=>{t.exports=e},959:e=>{e.exports=t}},r={};function l(e){var t=r[e];if(void 0!==t)return t.exports;var o=r[e]={exports:{}};return n[e](o,o.exports,l),o.exports}l.n=e=>{var t=e&&e.__esModule?()=>e.default:()=>e;return l.d(t,{a:t}),t},l.d=(e,t)=>{for(var o in t)l.o(t,o)&&!l.o(e,o)&&Object.defineProperty(e,o,{enumerable:!0,get:t[o]})},l.o=(e,t)=>Object.prototype.hasOwnProperty.call(e,t),l.r=e=>{"undefined"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})};var s={};l.r(s),l.d(s,{plugin:()=>f});var d=l(781);function p(e,t,o,a,i,n,r){try{var l=e[n](r),s=l.value}catch(e){return void o(e)}l.done?t(s):Promise.resolve(s).then(a,i)}var c=l(959),g=l.n(c),m=l(89),u=l(7),x=l(531);const v=()=>({wrapper:m.css`
      position: relative;
      display: flex;
      flex-direction: column;
      height: 100%;
      overflow: hidden;
    `,header:m.css`
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 8px 16px;
      border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    `,title:m.css`
      margin: 0;
      font-size: 16px;
      font-weight: 500;
    `,count:m.css`
      font-size: 12px;
      opacity: 0.7;
    `,currentHoverWidget:m.css`
      display: flex;
      flex-direction: column;
      overflow: hidden;
      flex: 1;
      height: 100%;
    `,widgetHeader:m.css`
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;
      padding-bottom: 8px;
      border-bottom: 1px solid rgba(115, 191, 105, 0.3);
    `,widgetTitle:m.css`
      font-size: 11px;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 1px;
      color: rgba(115, 191, 105, 1);
    `,widgetLive:m.css`
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
    `,widgetContent:m.css`
      display: flex;
      flex-direction: column;
      gap: 8px;
      flex: 1;
      min-height: 0;
    `,widgetCompactHeader:m.css`
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 4px 0;
      font-size: 11px;
      flex-wrap: wrap;
      border-bottom: 1px solid rgba(255, 255, 255, 0.1);
      padding-bottom: 6px;
      margin-bottom: 4px;
    `,widgetCompactPanel:m.css`
      font-weight: 600;
      color: rgba(255, 255, 255, 0.9);
    `,widgetCompactSeparator:m.css`
      color: rgba(255, 255, 255, 0.3);
    `,widgetCompactSeries:m.css`
      color: rgba(255, 255, 255, 0.7);
    `,widgetCompactValue:m.css`
      font-weight: 700;
      color: rgba(115, 191, 105, 1);
      font-family: "Roboto Mono", monospace;
    `,widgetCompactTime:m.css`
      color: rgba(255, 255, 255, 0.6);
      font-family: "Roboto Mono", monospace;
      font-size: 11px;
    `,widgetLoading:m.css`
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 24px;
      color: rgba(255, 255, 255, 0.5);
      font-style: italic;
    `,widgetMainValue:m.css`
      display: flex;
      flex-direction: column;
      gap: 4px;
    `,widgetSeriesName:m.css`
      font-size: 13px;
      font-weight: 600;
      color: rgba(255, 255, 255, 0.9);
      letter-spacing: 0.3px;
    `,widgetValue:m.css`
      font-size: 32px;
      font-weight: 700;
      color: rgba(115, 191, 105, 1);
      line-height: 1.2;
      font-family: "Roboto Mono", monospace;
    `,widgetMetadata:m.css`
      display: flex;
      flex-direction: column;
      gap: 6px;
      padding-top: 8px;
      border-top: 1px solid rgba(255, 255, 255, 0.1);
    `,widgetMetaRow:m.css`
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 8px;
    `,widgetMetaLabel:m.css`
      font-size: 11px;
      font-weight: 600;
      color: rgba(255, 255, 255, 0.6);
      text-transform: uppercase;
      letter-spacing: 0.5px;
    `,widgetMetaValue:m.css`
      font-size: 12px;
      font-weight: 500;
      color: rgba(255, 255, 255, 0.9);
      text-align: right;
    `,noCurrentHover:m.css`
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
    `,noHoverIcon:m.css`
      font-size: 48px;
      margin-bottom: 12px;
      opacity: 0.5;
    `,noHoverText:m.css`
      font-size: 14px;
      color: rgba(255, 255, 255, 0.6);
      text-align: center;
    `,widgetLogsSection:m.css`
      display: flex;
      flex-direction: column;
      flex: 1;
      min-height: 0;
      padding-top: 8px;
    `,widgetLogsSectionHeader:m.css`
      font-size: 11px;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 0.5px;
      color: rgba(115, 191, 105, 0.8);
      margin-bottom: 8px;
    `,widgetLogsList:m.css`
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
    `,widgetLogItem:m.css`
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
    `,widgetLogToggle:m.css`
      color: rgba(115, 191, 105, 0.7);
      font-size: 9px;
      flex-shrink: 0;
      width: 12px;
      margin-top: 3px;
    `,widgetLogText:m.css`
      color: rgba(255, 255, 255, 0.85);
      word-break: break-word;
      flex: 1;
    `,widgetLogDivider:m.css`
      height: 1px;
      background: linear-gradient(
        to right,
        transparent,
        rgba(115, 191, 105, 0.3) 20%,
        rgba(115, 191, 105, 0.3) 80%,
        transparent
      );
      margin: 12px 0;
    `,eventList:m.css`
      flex: 1;
      overflow-y: auto;
      padding: 8px;
    `,event:m.css`
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
    `,eventHeader:m.css`
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 4px;
    `,panelName:m.css`
      font-weight: 500;
      font-size: 14px;
    `,timestamp:m.css`
      font-size: 11px;
      opacity: 0.6;
    `,eventDetails:m.css`
      display: flex;
      flex-direction: column;
      gap: 8px;
      font-size: 12px;
      opacity: 0.9;
    `,detail:m.css`
      display: inline-block;
    `,metricInfo:m.css`
      display: flex;
      gap: 8px;
      align-items: flex-start;
    `,metricLabel:m.css`
      font-weight: 600;
      color: rgba(255, 255, 255, 0.7);
      min-width: 60px;
    `,metricValue:m.css`
      color: rgba(255, 255, 255, 0.9);
      word-break: break-word;
      flex: 1;
    `,noEvents:m.css`
      text-align: center;
      padding: 32px;
      opacity: 0.5;
    `,debugInfo:m.css`
      margin-top: 8px;
      border-top: 1px solid rgba(255, 255, 255, 0.1);
      padding-top: 8px;
    `,debugSummary:m.css`
      cursor: pointer;
      font-size: 11px;
      color: rgba(255, 255, 255, 0.6);
      &:hover {
        color: rgba(255, 255, 255, 0.8);
      }
    `,debugPre:m.css`
      background: rgba(0, 0, 0, 0.3);
      padding: 8px;
      border-radius: 4px;
      font-size: 10px;
      overflow-x: auto;
      margin: 4px 0 0 0;
      max-height: 200px;
      overflow-y: auto;
    `,customTooltip:m.css`
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
    `,tooltipHeader:m.css`
      margin-bottom: 8px;
      padding-bottom: 6px;
      border-bottom: 1px solid rgba(115, 191, 105, 0.3);
      font-size: 14px;
      color: rgba(115, 191, 105, 1);
    `,tooltipContent:m.css`
      display: flex;
      flex-direction: column;
      gap: 4px;
    `,tooltipRow:m.css`
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      gap: 8px;
    `,tooltipLabel:m.css`
      font-weight: 600;
      color: rgba(255, 255, 255, 0.7);
      min-width: 50px;
      flex-shrink: 0;
    `,tooltipValue:m.css`
      color: rgba(255, 255, 255, 0.9);
      word-break: break-word;
      text-align: right;
      flex: 1;
    `,tooltipDivider:m.css`
      margin: 8px 0;
      height: 1px;
      background: rgba(115, 191, 105, 0.3);
    `,tooltipSection:m.css`
      margin-top: 8px;
    `,tooltipSectionHeader:m.css`
      font-weight: 600;
      color: rgba(115, 191, 105, 0.8);
      font-size: 11px;
      text-transform: uppercase;
      margin-bottom: 6px;
      letter-spacing: 0.5px;
    `,tooltipLoading:m.css`
      color: rgba(255, 255, 255, 0.6);
      font-style: italic;
      font-size: 11px;
    `,tooltipDataFrame:m.css`
      margin-bottom: 6px;
      padding-bottom: 4px;
      border-bottom: 1px solid rgba(255, 255, 255, 0.1);
      &:last-child {
        border-bottom: none;
        margin-bottom: 0;
      }
    `,tooltipFields:m.css`
      margin-top: 4px;
    `,tooltipMuted:m.css`
      color: rgba(255, 255, 255, 0.5);
      font-size: 11px;
      font-style: italic;
    `}),f=new d.PanelPlugin(e=>{var t,o,a,i;const{options:n,width:r,height:l,id:s,eventBus:f}=e,b=(0,u.useStyles2)(v),[w,h]=(0,c.useState)(null),[y,k]=(0,c.useState)([]),[L,S]=(0,c.useState)(new Set),[N,I]=(0,c.useState)(!1),E=(0,c.useCallback)((e,t)=>{return(o=function*(){try{var o,a,i;I(!0),k([]);const r=e.metricData,l=new Date,s=new Date(Date.now()-n.timeWindowMs),d=(null==r?void 0:r.seriesName)||(null==r?void 0:r.fieldName)||"Unknown Series",p=e.panelTitle||"Unknown Panel",c=null==t||null===(o=t._eventsOrigin)||void 0===o?void 0:o._state,g=(null==c?void 0:c.title)||"Unknown Dashboard",m={org:String((null===(i=x.config.bootData)||void 0===i||null===(a=i.user)||void 0===a?void 0:a.orgId)||1),dashboard:g,panel_title:p,metric_name:d,start_time:s.toISOString(),end_time:l.toISOString()},u=yield(0,x.getBackendSrv)().post("/api/plugins/hover-hover-panel/resources/query_logs",m);if(u&&Array.isArray(u.log_groups)){const e=[],t=n.maxLogs,o=n.maxLogLength;let a=0;u.log_groups.forEach((i,n)=>{if(a>=t)return;const r=i.relative_change,l=r>0?"↑":r<0?"↓":"→",s=r>50?"🔴":r>10?"🟠":r<-10?"🟢":"⚪";e.push(`${s} ${l} ${r.toFixed(1)}% change from baseline`),a++,Array.isArray(i.representative_logs)&&i.representative_logs.forEach(i=>{if(a>=t)return;const n=i.length>o?i.substring(0,o)+"... [truncated]":i;e.push(`  ${n}`),a++})}),a>=t&&e.push(`... [${u.log_groups.length-e.length} more logs truncated for performance]`),k(e),I(!1)}else k([]),I(!1)}catch(e){k([`Error: ${e}`]),I(!1)}},function(){var e=this,t=arguments;return new Promise(function(a,i){var n=o.apply(e,t);function r(e){p(n,a,i,r,l,"next",e)}function l(e){p(n,a,i,r,l,"throw",e)}r(void 0)})})();var o},[n.timeWindowMs,n.maxLogs,n.maxLogLength]);return(0,c.useEffect)(()=>{const e=(0,x.getAppEvents)(),t=null==f?void 0:f.getStream(d.DataHoverEvent).subscribe(e=>{var t,o,a,i,n,r,l,s,d;const p=null==e?void 0:e.origin,c=null===(t=e.payload)||void 0===t?void 0:t.origin,g=e.payload,m=(null==g?void 0:g.point)||{},u=null==g?void 0:g.data;let x,v,f,b="Unknown Series",w="";if(u&&u.fields){var y,k,L;const e=null!==(L=null==g?void 0:g.columnIndex)&&void 0!==L?L:null==g?void 0:g.fieldIndex;var S;const t=null!==(S=null==g?void 0:g.rowIndex)&&void 0!==S?S:null==g?void 0:g.dataIndex;let o=e;if(void 0===o&&(o=u.fields.findIndex(e=>"time"!==e.type)),void 0!==o&&u.fields[o]){var N;const e=u.fields[o];w=e.name||"",b=(null===(N=e.config)||void 0===N?void 0:N.displayName)||e.name||u.name||"Unknown Series",void 0!==t&&e.values&&e.values.length>t&&(x=e.values[t],v=e.display?e.display(x).text:String(x))}void 0!==t&&(null===(k=u.fields[0])||void 0===k||null===(y=k.values)||void 0===y?void 0:y.length)>t&&(f=u.fields[0].values[t])}void 0===x&&void 0!==(null==m?void 0:m.value)&&(x=m.value),"Unknown Series"===b&&(null==g?void 0:g.dataId)&&(b=g.dataId);const I=null==u?void 0:u.refId;let _="Unknown Panel",z="unknown";var C,T;(null==p||null===(a=p._eventsOrigin)||void 0===a||null===(o=a._state)||void 0===o?void 0:o.title)?_=p._eventsOrigin._state.title:(null==p||null===(i=p._state)||void 0===i?void 0:i.title)?_=p._state.title:(null==c||null===(r=c._eventsOrigin)||void 0===r||null===(n=r._state)||void 0===n?void 0:n.title)?_=c._eventsOrigin._state.title:(null==c||null===(l=c._state)||void 0===l?void 0:l.title)?_=c._state.title:(null==u||null===(d=u.meta)||void 0===d||null===(s=d.custom)||void 0===s?void 0:s.title)?_=u.meta.custom.title:(null==u?void 0:u.name)?_=`Panel: ${u.name}`:I&&(_=`Query ${I}`),I?z=I:(null==p?void 0:p.panelId)?z=p.panelId:(null==c?void 0:c.panelId)&&(z=c.panelId);const D={panelId:z,panelTitle:_,timestamp:Date.now(),x:(null==m?void 0:m.time)||(null==m?void 0:m.x)||0,y:x||(null==m?void 0:m.y)||0,elementType:"data-hover",metricData:{seriesName:b,fieldName:w,value:x,formattedValue:v,time:f,fieldIndex:null!==(C=null==g?void 0:g.columnIndex)&&void 0!==C?C:null==g?void 0:g.fieldIndex,dataIndex:null!==(T=null==g?void 0:g.rowIndex)&&void 0!==T?T:null==g?void 0:g.dataIndex,point:m,dataFrame:u}};h(D),E(D,p)}),o=e.subscribe(d.LegacyGraphHoverEvent,e=>{var t,o,a,i,n,r,l;const s={panelId:"legacy-graph-hover",panelTitle:"Legacy Graph Hover",timestamp:Date.now(),x:(null===(o=e.payload)||void 0===o||null===(t=o.pos)||void 0===t?void 0:t.x)||0,y:(null===(i=e.payload)||void 0===i||null===(a=i.pos)||void 0===a?void 0:a.y)||0,elementType:"legacy-graph-hover",metricData:{seriesName:(null===(r=e.payload)||void 0===r||null===(n=r.series)||void 0===n?void 0:n.name)||"Legacy Series",value:null===(l=e.payload)||void 0===l?void 0:l.value}};h(s),E(s)});return()=>{null==t||t.unsubscribe(),null==o||o.unsubscribe()}},[s,E,f]),g().createElement(g().Fragment,null,g().createElement("div",{className:(0,m.cx)(b.wrapper),style:{width:r,height:l}},w&&g().createElement("div",{className:b.currentHoverWidget},g().createElement("div",{className:b.widgetContent},g().createElement("div",{className:b.widgetCompactHeader},(null===(t=w.metricData)||void 0===t?void 0:t.time)&&g().createElement(g().Fragment,null,g().createElement("span",{className:b.widgetCompactTime},(0,d.dateTime)(w.metricData.time).format("YYYY-MM-DD HH:mm:ss")),g().createElement("span",{className:b.widgetCompactSeparator},"•")),g().createElement("span",{className:b.widgetCompactPanel},w.panelTitle),g().createElement("span",{className:b.widgetCompactSeparator},"•"),g().createElement("span",{className:b.widgetCompactSeries},(null===(o=w.metricData)||void 0===o?void 0:o.seriesName)||"No Series"),g().createElement("span",{className:b.widgetCompactSeparator},"•"),g().createElement("span",{className:b.widgetCompactValue},(null===(a=w.metricData)||void 0===a?void 0:a.formattedValue)||(null===(i=w.metricData)||void 0===i?void 0:i.value)||"—")),N?g().createElement("div",{className:b.widgetLoading},"Loading logs..."):y.length>0?g().createElement("div",{className:b.widgetLogsSection},g().createElement("div",{className:b.widgetLogsList},y.map((e,t)=>{const o=L.has(t),a=n.logTruncateLength,i=e.length>a;return g().createElement("div",{key:t,className:b.widgetLogItem,onClick:()=>i&&(e=>{S(t=>{const o=new Set(t);return o.has(e)?o.delete(e):o.add(e),o})})(t),style:{cursor:i?"pointer":"default"},title:i?"Click to expand/collapse":void 0},i&&g().createElement("span",{className:b.widgetLogToggle},o?"▼":"▶"),g().createElement("span",{className:b.widgetLogText},o||!i?e:e.substring(0,a)+"... (click to expand)"))}))):g().createElement("div",{className:b.widgetLoading},"No logs found"))),!w&&g().createElement("div",{className:b.noCurrentHover},g().createElement("div",{className:b.noHoverIcon},"👆"),g().createElement("div",{className:b.noHoverText},"Hover over a panel to see live data"))))}).setPanelOptions(e=>((e=>{e.addNumberInput({path:"timeWindowMs",name:"Time Window (ms)",description:"Time window for log queries in milliseconds (default: 1 hour)",defaultValue:36e5,settings:{min:1e3,step:1e3}}).addNumberInput({path:"maxLogs",name:"Max Logs",description:"Maximum number of log entries to display",defaultValue:500,settings:{min:1,step:1}}).addNumberInput({path:"maxLogLength",name:"Max Log Length",description:"Maximum length of individual log entry",defaultValue:1e4,settings:{min:100,step:100}}).addNumberInput({path:"logTruncateLength",name:"Log Truncate Length",description:"Character count threshold for expandable logs",defaultValue:120,settings:{min:50,step:10}})})(e),e));return s})());
//# sourceMappingURL=module.js.map