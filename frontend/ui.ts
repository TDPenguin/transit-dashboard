// HTML rendering functions for train predictions

import type { TrainPrediction } from './types.js';
import { DEST_EXPANSIONS } from './utils.js';

// Build HTML for train predictions with styling
export function buildTrainPredictionsHTML(trains: TrainPrediction[]): string {
    if (trains.length === 0) {
        return `<div style="margin-top: 16px; padding: 12px; background: #f5f5f5; border-radius: 4px;">
            <p style="margin: 0; color: #666;"><em>No trains currently predicted</em></p>
        </div>`;
    }
    
    // Use exact match instead of includes to avoid partial replacements
    const expandDest = (dest: string) => DEST_EXPANSIONS[dest] || dest;
    
    const formatTime = (min: string): { text: string; color: string } => {
        if (min === "BRD") return { text: "Boarding", color: 'color: #d32f2f;' };
        if (min === "ARR") return { text: "Arriving", color: 'color: #f57c00;' };
        if (!min || min.trim() === '' || min === '---') return { text: "—", color: 'color: #999;' };
        return { text: `${min} min`, color: '' };
    };
    
    const buildRow = (train: TrainPrediction): string => {
        const dest = expandDest(train.DestinationName);
        const { text: time, color: timeColor } = formatTime(train.Min);
        
        // Use X.svg for non-passenger trains
        const isRegular = train.Line && train.Line !== 'No' && train.Line !== '--';
        const icon = isRegular 
            ? `<img src="assets/${train.Line}.svg" alt="${train.Line}" style="height: 16px; width: 16px;">`
            : `<img src="assets/X.svg" alt="No Passenger" style="height: 16px; width: 16px;">`;
        
        return `<div style="display: flex; justify-content: space-between; align-items: center; padding: 6px 0; border-bottom: 1px solid #ddd;">
            <div style="display: flex; align-items: center; gap: 6px;">${icon}<span style="font-size: 1em;">${dest}</span></div>
            <span style="font-size: 0.9em; font-weight: 600; ${timeColor}">${time}</span>
        </div>`;
    };
    
    // Group by track, then platform
    const tracks: { [trackNum: string]: { [platform: string]: TrainPrediction[] } } = { '1': {}, '2': {} };
    for (const train of trains) {
        const platform = train.LocationCode;
        const track = train.Group || '1';
        if (!tracks[track]) tracks[track] = {};
        if (!tracks[track][platform]) tracks[track][platform] = [];
        tracks[track][platform].push(train);
    }
    
    const buildColumn = (trackNum: string): string => {
        let html = `<div><h3 style="font-size: 0.9em; font-weight: bold; margin: 0 0 8px 0;">Track ${trackNum}</h3>`;
        const platforms = tracks[trackNum] || {};
        const platformKeys = Object.keys(platforms).sort();
        
        if (platformKeys.length === 0) {
            html += `<div style="padding: 10px; background: #f9f9f9; border-radius: 4px; color: #999; font-size: 0.85em; font-style: italic;">No trains predicted</div>`;
        } else {
            for (const platform of platformKeys) {
                const platformTrains = platforms[platform];
                if (!platformTrains) continue;
                html += `<div style="margin-bottom: 12px; padding: 10px; background: #f9f9f9; border-radius: 4px; max-height: 150px; overflow-y: auto;">`;
                platformTrains.forEach(train => html += buildRow(train));
                html += '</div>';
            }
        }
        html += '</div>';
        return html;
    };
    
    return `<div style="margin-top: 16px;">
        <p style="font-weight: bold; margin-bottom: 8px;">Next Trains:</p>
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px;">
            ${buildColumn('1')}
            ${buildColumn('2')}
        </div>
    </div>`;
}
