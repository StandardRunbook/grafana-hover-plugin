# Screenshots Guide

Before submitting the plugin to Grafana, you need to add screenshots showing the plugin in action.

## Required Screenshots

Place these images in `src/img/`:

### 1. screenshot-main.png
**Main view showing logs on hover**
- Show a dashboard with multiple panels
- One panel should have a visible hover/crosshair
- The Hover panel should be visible with logs displayed
- Recommended size: 1920x1080 or similar
- Format: PNG

### 2. screenshot-config.png
**Panel configuration**
- Show the panel editor with the Hover panel options visible
- Display the configuration fields:
  - API Endpoint
  - Time Window (ms)
  - Max Logs
  - Max Log Length
  - Log Truncate Length
- Recommended size: 1920x1080 or similar
- Format: PNG

## How to Create Screenshots

1. Start Grafana with the Hover plugin installed
2. Create a dashboard with:
   - At least 2-3 time-series panels (e.g., CPU, Memory, Network)
   - One Hover panel configured with your API
3. For screenshot-main.png:
   - Hover over a data point in one of the panels
   - Wait for logs to appear in the Hover panel
   - Take a full browser screenshot
4. For screenshot-config.png:
   - Click Edit on the Hover panel
   - Scroll to show all configuration options
   - Take a screenshot of the panel editor

## Tips for Good Screenshots

- Use a clean, professional dashboard theme
- Make sure the logs shown are realistic (not too many error messages)
- Avoid showing sensitive data in logs or metrics
- Use high resolution (at least 1920x1080)
- Ensure text is readable
- Show the full context (dashboard name, time range, etc.)

## Placeholder Images

Currently, the plugin.json references these screenshots, but they don't exist yet. 
You need to create them before submitting to Grafana.

Alternatively, you can remove the screenshots section from plugin.json temporarily
if you want to test the submission process first.
