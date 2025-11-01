define(["@grafana/data","react","@emotion/css","@grafana/ui","@grafana/runtime"],(e,t,o,a,i)=>(()=>{"use strict";var n={7:e=>{e.exports=a},89:e=>{e.exports=o},531:e=>{e.exports=i},781:t=>{t.exports=e},959:e=>{e.exports=t}},r={};function l(e){var t=r[e];if(void 0!==t)return t.exports;var o=r[e]={exports:{}};return n[e](o,o.exports,l),o.exports}l.n=e=>{var t=e&&e.__esModule?()=>e.default:()=>e;return l.d(t,{a:t}),t},l.d=(e,t)=>{for(var o in t)l.o(t,o)&&!l.o(e,o)&&Object.defineProperty(e,o,{enumerable:!0,get:t[o]})},l.o=(e,t)=>Object.prototype.hasOwnProperty.call(e,t),l.r=e=>{"undefined"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})};var s={};l.r(s),l.d(s,{plugin:()=>v});var d=l(781);Object.create,Object.create,"function"==typeof SuppressedError&&SuppressedError;var p=l(959),c=l.n(p),g=l(89),m=l(7),u=l(531);const x=()=>({wrapper:g.css`
      position: relative;
      display: flex;
      flex-direction: column;
      height: 100%;
      overflow: hidden;
    `,header:g.css`
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 8px 16px;
      border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    `,title:g.css`
      margin: 0;
      font-size: 16px;
      font-weight: 500;
    `,count:g.css`
      font-size: 12px;
      opacity: 0.7;
    `,currentHoverWidget:g.css`
      display: flex;
      flex-direction: column;
      overflow: hidden;
      flex: 1;
      height: 100%;
    `,widgetHeader:g.css`
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;
      padding-bottom: 8px;
      border-bottom: 1px solid rgba(115, 191, 105, 0.3);
    `,widgetTitle:g.css`
      font-size: 11px;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 1px;
      color: rgba(115, 191, 105, 1);
    `,widgetLive:g.css`
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
    `,widgetContent:g.css`
      display: flex;
      flex-direction: column;
      gap: 8px;
      flex: 1;
      min-height: 0;
    `,widgetCompactHeader:g.css`
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 4px 0;
      font-size: 11px;
      flex-wrap: wrap;
      border-bottom: 1px solid rgba(255, 255, 255, 0.1);
      padding-bottom: 6px;
      margin-bottom: 4px;
    `,widgetCompactPanel:g.css`
      font-weight: 600;
      color: rgba(255, 255, 255, 0.9);
    `,widgetCompactSeparator:g.css`
      color: rgba(255, 255, 255, 0.3);
    `,widgetCompactSeries:g.css`
      color: rgba(255, 255, 255, 0.7);
    `,widgetCompactValue:g.css`
      font-weight: 700;
      color: rgba(115, 191, 105, 1);
      font-family: "Roboto Mono", monospace;
    `,widgetCompactTime:g.css`
      color: rgba(255, 255, 255, 0.6);
      font-family: "Roboto Mono", monospace;
      font-size: 11px;
    `,widgetLoading:g.css`
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 24px;
      color: rgba(255, 255, 255, 0.5);
      font-style: italic;
    `,widgetMainValue:g.css`
      display: flex;
      flex-direction: column;
      gap: 4px;
    `,widgetSeriesName:g.css`
      font-size: 13px;
      font-weight: 600;
      color: rgba(255, 255, 255, 0.9);
      letter-spacing: 0.3px;
    `,widgetValue:g.css`
      font-size: 32px;
      font-weight: 700;
      color: rgba(115, 191, 105, 1);
      line-height: 1.2;
      font-family: "Roboto Mono", monospace;
    `,widgetMetadata:g.css`
      display: flex;
      flex-direction: column;
      gap: 6px;
      padding-top: 8px;
      border-top: 1px solid rgba(255, 255, 255, 0.1);
    `,widgetMetaRow:g.css`
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 8px;
    `,widgetMetaLabel:g.css`
      font-size: 11px;
      font-weight: 600;
      color: rgba(255, 255, 255, 0.6);
      text-transform: uppercase;
      letter-spacing: 0.5px;
    `,widgetMetaValue:g.css`
      font-size: 12px;
      font-weight: 500;
      color: rgba(255, 255, 255, 0.9);
      text-align: right;
    `,noCurrentHover:g.css`
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
    `,noHoverIcon:g.css`
      font-size: 48px;
      margin-bottom: 12px;
      opacity: 0.5;
    `,noHoverText:g.css`
      font-size: 14px;
      color: rgba(255, 255, 255, 0.6);
      text-align: center;
    `,widgetLogsSection:g.css`
      display: flex;
      flex-direction: column;
      flex: 1;
      min-height: 0;
      padding-top: 8px;
    `,widgetLogsSectionHeader:g.css`
      font-size: 11px;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 0.5px;
      color: rgba(115, 191, 105, 0.8);
      margin-bottom: 8px;
    `,widgetLogsList:g.css`
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
    `,widgetLogItem:g.css`
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
    `,widgetLogToggle:g.css`
      color: rgba(115, 191, 105, 0.7);
      font-size: 9px;
      flex-shrink: 0;
      width: 12px;
      margin-top: 3px;
    `,widgetLogText:g.css`
      color: rgba(255, 255, 255, 0.85);
      word-break: break-word;
      flex: 1;
    `,widgetLogDivider:g.css`
      height: 1px;
      background: linear-gradient(
        to right,
        transparent,
        rgba(115, 191, 105, 0.3) 20%,
        rgba(115, 191, 105, 0.3) 80%,
        transparent
      );
      margin: 12px 0;
    `,eventList:g.css`
      flex: 1;
      overflow-y: auto;
      padding: 8px;
    `,event:g.css`
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
    `,eventHeader:g.css`
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 4px;
    `,panelName:g.css`
      font-weight: 500;
      font-size: 14px;
    `,timestamp:g.css`
      font-size: 11px;
      opacity: 0.6;
    `,eventDetails:g.css`
      display: flex;
      flex-direction: column;
      gap: 8px;
      font-size: 12px;
      opacity: 0.9;
    `,detail:g.css`
      display: inline-block;
    `,metricInfo:g.css`
      display: flex;
      gap: 8px;
      align-items: flex-start;
    `,metricLabel:g.css`
      font-weight: 600;
      color: rgba(255, 255, 255, 0.7);
      min-width: 60px;
    `,metricValue:g.css`
      color: rgba(255, 255, 255, 0.9);
      word-break: break-word;
      flex: 1;
    `,noEvents:g.css`
      text-align: center;
      padding: 32px;
      opacity: 0.5;
    `,debugInfo:g.css`
      margin-top: 8px;
      border-top: 1px solid rgba(255, 255, 255, 0.1);
      padding-top: 8px;
    `,debugSummary:g.css`
      cursor: pointer;
      font-size: 11px;
      color: rgba(255, 255, 255, 0.6);
      &:hover {
        color: rgba(255, 255, 255, 0.8);
      }
    `,debugPre:g.css`
      background: rgba(0, 0, 0, 0.3);
      padding: 8px;
      border-radius: 4px;
      font-size: 10px;
      overflow-x: auto;
      margin: 4px 0 0 0;
      max-height: 200px;
      overflow-y: auto;
    `,customTooltip:g.css`
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
    `,tooltipHeader:g.css`
      margin-bottom: 8px;
      padding-bottom: 6px;
      border-bottom: 1px solid rgba(115, 191, 105, 0.3);
      font-size: 14px;
      color: rgba(115, 191, 105, 1);
    `,tooltipContent:g.css`
      display: flex;
      flex-direction: column;
      gap: 4px;
    `,tooltipRow:g.css`
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      gap: 8px;
    `,tooltipLabel:g.css`
      font-weight: 600;
      color: rgba(255, 255, 255, 0.7);
      min-width: 50px;
      flex-shrink: 0;
    `,tooltipValue:g.css`
      color: rgba(255, 255, 255, 0.9);
      word-break: break-word;
      text-align: right;
      flex: 1;
    `,tooltipDivider:g.css`
      margin: 8px 0;
      height: 1px;
      background: rgba(115, 191, 105, 0.3);
    `,tooltipSection:g.css`
      margin-top: 8px;
    `,tooltipSectionHeader:g.css`
      font-weight: 600;
      color: rgba(115, 191, 105, 0.8);
      font-size: 11px;
      text-transform: uppercase;
      margin-bottom: 6px;
      letter-spacing: 0.5px;
    `,tooltipLoading:g.css`
      color: rgba(255, 255, 255, 0.6);
      font-style: italic;
      font-size: 11px;
    `,tooltipDataFrame:g.css`
      margin-bottom: 6px;
      padding-bottom: 4px;
      border-bottom: 1px solid rgba(255, 255, 255, 0.1);
      &:last-child {
        border-bottom: none;
        margin-bottom: 0;
      }
    `,tooltipFields:g.css`
      margin-top: 4px;
    `,tooltipMuted:g.css`
      color: rgba(255, 255, 255, 0.5);
      font-size: 11px;
      font-style: italic;
    `}),v=new d.PanelPlugin(e=>{var t,o,a,i;const{options:n,width:r,height:l,id:s,eventBus:v}=e,f=(0,m.useStyles2)(x),[b,w]=(0,p.useState)(null),[h,y]=(0,p.useState)([]),[k,L]=(0,p.useState)(new Set),[S,N]=(0,p.useState)(!1),E=(0,p.useCallback)((e,t)=>{return o=void 0,a=void 0,r=function*(){var o,a,i;try{N(!0),y([]);const r=e.metricData,l=new Date,s=new Date(Date.now()-n.timeWindowMs),d=(null==r?void 0:r.seriesName)||(null==r?void 0:r.fieldName)||"Unknown Series",p=e.panelTitle||"Unknown Panel",c=null===(o=null==t?void 0:t._eventsOrigin)||void 0===o?void 0:o._state,g=(null==c?void 0:c.title)||"Unknown Dashboard",m={org:String((null===(i=null===(a=u.config.bootData)||void 0===a?void 0:a.user)||void 0===i?void 0:i.orgId)||1),dashboard:g,panel_title:p,metric_name:d,start_time:s.toISOString(),end_time:l.toISOString()},x=yield(0,u.getBackendSrv)().post("/api/plugins/hover-hover-panel/resources/query_logs",m);if(x&&Array.isArray(x.log_groups)){const e=[],t=n.maxLogs,o=n.maxLogLength;let a=0;x.log_groups.forEach((i,n)=>{if(a>=t)return;const r=i.relative_change,l=r>0?"↑":r<0?"↓":"→",s=r>50?"🔴":r>10?"🟠":r<-10?"🟢":"⚪";e.push(`${s} ${l} ${r.toFixed(1)}% change from baseline`),a++,Array.isArray(i.representative_logs)&&i.representative_logs.forEach(i=>{if(a>=t)return;const n=i.length>o?i.substring(0,o)+"... [truncated]":i;e.push(`  ${n}`),a++})}),a>=t&&e.push(`... [${x.log_groups.length-e.length} more logs truncated for performance]`),y(e),N(!1)}else y([]),N(!1)}catch(e){y([`Error: ${e}`]),N(!1)}},new((i=void 0)||(i=Promise))(function(e,t){function n(e){try{s(r.next(e))}catch(e){t(e)}}function l(e){try{s(r.throw(e))}catch(e){t(e)}}function s(t){var o;t.done?e(t.value):(o=t.value,o instanceof i?o:new i(function(e){e(o)})).then(n,l)}s((r=r.apply(o,a||[])).next())});var o,a,i,r},[n.timeWindowMs,n.maxLogs,n.maxLogLength]);return(0,p.useEffect)(()=>{const e=(0,u.getAppEvents)(),t=null==v?void 0:v.getStream(d.DataHoverEvent).subscribe(e=>{var t,o,a,i,n,r,l,s,d,p,c,g,m,u,x,v;const f=null==e?void 0:e.origin,b=null===(t=e.payload)||void 0===t?void 0:t.origin,h=e.payload,y=(null==h?void 0:h.point)||{},k=null==h?void 0:h.data;let L,S,N,I="Unknown Series",_="";if(k&&k.fields){const e=null!==(o=null==h?void 0:h.columnIndex)&&void 0!==o?o:null==h?void 0:h.fieldIndex,t=null!==(a=null==h?void 0:h.rowIndex)&&void 0!==a?a:null==h?void 0:h.dataIndex;let l=e;if(void 0===l&&(l=k.fields.findIndex(e=>"time"!==e.type)),void 0!==l&&k.fields[l]){const e=k.fields[l];_=e.name||"",I=(null===(i=e.config)||void 0===i?void 0:i.displayName)||e.name||k.name||"Unknown Series",void 0!==t&&e.values&&e.values.length>t&&(L=e.values[t],S=e.display?e.display(L).text:String(L))}void 0!==t&&(null===(r=null===(n=k.fields[0])||void 0===n?void 0:n.values)||void 0===r?void 0:r.length)>t&&(N=k.fields[0].values[t])}void 0===L&&void 0!==(null==y?void 0:y.value)&&(L=y.value),"Unknown Series"===I&&(null==h?void 0:h.dataId)&&(I=h.dataId);const z=null==k?void 0:k.refId;let C="Unknown Panel",T="unknown";(null===(s=null===(l=null==f?void 0:f._eventsOrigin)||void 0===l?void 0:l._state)||void 0===s?void 0:s.title)?C=f._eventsOrigin._state.title:(null===(d=null==f?void 0:f._state)||void 0===d?void 0:d.title)?C=f._state.title:(null===(c=null===(p=null==b?void 0:b._eventsOrigin)||void 0===p?void 0:p._state)||void 0===c?void 0:c.title)?C=b._eventsOrigin._state.title:(null===(g=null==b?void 0:b._state)||void 0===g?void 0:g.title)?C=b._state.title:(null===(u=null===(m=null==k?void 0:k.meta)||void 0===m?void 0:m.custom)||void 0===u?void 0:u.title)?C=k.meta.custom.title:(null==k?void 0:k.name)?C=`Panel: ${k.name}`:z&&(C=`Query ${z}`),z?T=z:(null==f?void 0:f.panelId)?T=f.panelId:(null==b?void 0:b.panelId)&&(T=b.panelId);const D={panelId:T,panelTitle:C,timestamp:Date.now(),x:(null==y?void 0:y.time)||(null==y?void 0:y.x)||0,y:L||(null==y?void 0:y.y)||0,elementType:"data-hover",metricData:{seriesName:I,fieldName:_,value:L,formattedValue:S,time:N,fieldIndex:null!==(x=null==h?void 0:h.columnIndex)&&void 0!==x?x:null==h?void 0:h.fieldIndex,dataIndex:null!==(v=null==h?void 0:h.rowIndex)&&void 0!==v?v:null==h?void 0:h.dataIndex,point:y,dataFrame:k}};w(D),E(D,f)}),o=e.subscribe(d.LegacyGraphHoverEvent,e=>{var t,o,a,i,n,r,l;const s={panelId:"legacy-graph-hover",panelTitle:"Legacy Graph Hover",timestamp:Date.now(),x:(null===(o=null===(t=e.payload)||void 0===t?void 0:t.pos)||void 0===o?void 0:o.x)||0,y:(null===(i=null===(a=e.payload)||void 0===a?void 0:a.pos)||void 0===i?void 0:i.y)||0,elementType:"legacy-graph-hover",metricData:{seriesName:(null===(r=null===(n=e.payload)||void 0===n?void 0:n.series)||void 0===r?void 0:r.name)||"Legacy Series",value:null===(l=e.payload)||void 0===l?void 0:l.value}};w(s),E(s)});return()=>{null==t||t.unsubscribe(),null==o||o.unsubscribe()}},[s,E,v]),c().createElement(c().Fragment,null,c().createElement("div",{className:(0,g.cx)(f.wrapper),style:{width:r,height:l}},b&&c().createElement("div",{className:f.currentHoverWidget},c().createElement("div",{className:f.widgetContent},c().createElement("div",{className:f.widgetCompactHeader},(null===(t=b.metricData)||void 0===t?void 0:t.time)&&c().createElement(c().Fragment,null,c().createElement("span",{className:f.widgetCompactTime},(0,d.dateTime)(b.metricData.time).format("YYYY-MM-DD HH:mm:ss")),c().createElement("span",{className:f.widgetCompactSeparator},"•")),c().createElement("span",{className:f.widgetCompactPanel},b.panelTitle),c().createElement("span",{className:f.widgetCompactSeparator},"•"),c().createElement("span",{className:f.widgetCompactSeries},(null===(o=b.metricData)||void 0===o?void 0:o.seriesName)||"No Series"),c().createElement("span",{className:f.widgetCompactSeparator},"•"),c().createElement("span",{className:f.widgetCompactValue},(null===(a=b.metricData)||void 0===a?void 0:a.formattedValue)||(null===(i=b.metricData)||void 0===i?void 0:i.value)||"—")),S?c().createElement("div",{className:f.widgetLoading},"Loading logs..."):h.length>0?c().createElement("div",{className:f.widgetLogsSection},c().createElement("div",{className:f.widgetLogsList},h.map((e,t)=>{const o=k.has(t),a=n.logTruncateLength,i=e.length>a;return c().createElement("div",{key:t,className:f.widgetLogItem,onClick:()=>i&&(e=>{L(t=>{const o=new Set(t);return o.has(e)?o.delete(e):o.add(e),o})})(t),style:{cursor:i?"pointer":"default"},title:i?"Click to expand/collapse":void 0},i&&c().createElement("span",{className:f.widgetLogToggle},o?"▼":"▶"),c().createElement("span",{className:f.widgetLogText},o||!i?e:e.substring(0,a)+"... (click to expand)"))}))):c().createElement("div",{className:f.widgetLoading},"No logs found"))),!b&&c().createElement("div",{className:f.noCurrentHover},c().createElement("div",{className:f.noHoverIcon},"👆"),c().createElement("div",{className:f.noHoverText},"Hover over a panel to see live data"))))}).setPanelOptions(e=>((e=>{e.addNumberInput({path:"timeWindowMs",name:"Time Window (ms)",description:"Time window for log queries in milliseconds (default: 1 hour)",defaultValue:36e5,settings:{min:1e3,step:1e3}}).addNumberInput({path:"maxLogs",name:"Max Logs",description:"Maximum number of log entries to display",defaultValue:500,settings:{min:1,step:1}}).addNumberInput({path:"maxLogLength",name:"Max Log Length",description:"Maximum length of individual log entry",defaultValue:1e4,settings:{min:100,step:100}}).addNumberInput({path:"logTruncateLength",name:"Log Truncate Length",description:"Character count threshold for expandable logs",defaultValue:120,settings:{min:50,step:10}})})(e),e));return s})());
//# sourceMappingURL=module.js.map