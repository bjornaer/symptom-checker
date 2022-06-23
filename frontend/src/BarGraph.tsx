import React, { useMemo } from 'react';
import { Group } from '@visx/group';
import { Bar } from '@visx/shape';
import { scaleLinear, scaleBand } from '@visx/scale';
import { Ailment, BarGraphProps } from './types';
import { GradientTealBlue } from '@visx/gradient';

// Define the graph dimensions and margins
const width = 500;
const height = 500;
// const verticalMargin = 120;
const margin = { top: 20, bottom: 20, left: 20, right: 20 };

// create some default bounds
const xDefaultMax = width - margin.left - margin.right;
const yDefaultMax = height - margin.top - margin.bottom;


// Finally we'll embed it all in an SVG
export const BarGraph = ({data, width = xDefaultMax, height = yDefaultMax, events = false}: BarGraphProps) => {
    const ailments = data

    const xMax = width
    const yMax = height // - verticalMargin

    // Accesors
    const getAilment = (d: Ailment) => d.name;
    const getFreq = (d: Ailment) => d.frequency * 100;

    // Scale the graph by our data
    const xScale = useMemo(() => scaleBand<string>({
        range: [0, xMax],
        round: true,
        domain: ailments.map(getAilment),
        padding: 0.4,
        }), [xMax])
    const yScale = useMemo(()=> scaleLinear<number>({
        range: [yMax, 0],
        round: true,
        domain: [0, Math.max(...ailments.map(getFreq))],
        }),[yMax]) 

    return (
    <svg width={width} height={height}>
        <GradientTealBlue id="teal" />
        <rect width={width} height={height} fill="url(#teal)" rx={14} />
        <Group>
            {ailments && ailments.map((d, idx) => {
                const ailment = getAilment(d)
                const barWidth = xScale.bandwidth();
                const barHeight = yMax - (yScale(getFreq(d)) ?? 0);
                const barX = xScale(ailment);
                const barY = yMax - barHeight;
                return (
                    <Bar
                        key={`bar-${ailment}-${idx}`}
                        x={barX}
                        y={barY}
                        width={barWidth}
                        height={barHeight}
                        fill="rgba(23, 233, 217, .5)"
                        onClick={() => {
                        if (events) alert(`clicked: ${JSON.stringify(Object.values(d))}`);
                        }}
                    />
                    );
            })}
        </Group>
    </svg>
    );
}
