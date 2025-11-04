// Configuration constants and utilities

import type { StationInfo } from './types.js';

// Line colors
export const LINE_COLORS: { [key: string]: string } = {
    red: '#e11738',
    blue: '#0575bf',
    orange: '#f99219',
    green: '#00a94f',
    yellow: '#fdd200',
    silver: '#a4a09c'
};

// Destination name expansions
export const DEST_EXPANSIONS: { [key: string]: string } = {
    'N Carrollton': 'New Carrollton',
    'NewCrlton': 'New Carrollton',
    'Branch Ave': 'Branch Avenue',
    'Branch Av': 'Branch Avenue',
    'Mt Vern Sq': 'Mt Vernon Sq',
    'MtVern Sq': 'Mt Vernon Sq',
    'Shady Grv': 'Shady Grove'
};

// Helper to get timestamp for logs
export function getTimestamp(): string {
    return new Date().toLocaleTimeString('en-US', { 
        hour12: false, 
        hour: '2-digit', 
        minute: '2-digit', 
        second: '2-digit', 
        fractionalSecondDigits: 3 
    });
}

// Get all station codes including connected platforms (like Metro Center has both A01 and C01)
export function getStationCodes(stationInfo: StationInfo): string[] {
    const codes = [stationInfo.Code];
    if (stationInfo.StationTogether1) codes.push(stationInfo.StationTogether1);
    if (stationInfo.StationTogether2) codes.push(stationInfo.StationTogether2);
    return codes;
}
