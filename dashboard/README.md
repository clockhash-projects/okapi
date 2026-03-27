# Okapi Dashboard

The Okapi dashboard is a React-based frontend providing a streamlined, **Uptime Kuma-inspired** visual interface for the universal service health proxy.

## Core Goals

1. **The "Universal Status Page" (Proxy Only)**  
   A high-density, real-time grid in the style of Uptime Kuma showing **only the health of your proxy instances**.  
   - Status (UP / DOWN / DEGRADED)  
   - Response time  
   - Last checked  
   - Uptime %  
   - Clean, zero-noise “everything is green” view  

2. **Service Discovery & Pinning**  
   A dedicated section to search across 80+ upstream services.  
   - View current upstream status  
   - Pin / unpin services  
   - Build a custom monitored set  

3. **Consolidated Incident Center**  
   A separate unified feed aggregating incidents across providers.  
   - Filter by service, severity, time  
   - Show status (Investigating / Identified / Resolved)  
   - Centralized outage visibility  

4. **Unified Maintenance Timeline**  
   A forward-looking section for scheduled maintenance.  
   - Upcoming maintenance windows  
   - Affected services  
   - Expected impact  

## Development

The dashboard is built using React, TypeScript, and Vite.

### Setup

```bash
cd dashboard
pnpm install
pnpm dev