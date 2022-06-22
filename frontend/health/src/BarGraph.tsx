import React from 'react';
import { Group } from '@visx/group';
import { Bar } from '@visx/shape';
import { scaleLinear, scaleBand, StringLike } from '@visx/scale';
import { AilmentHistogram } from './types';
import { ScaleBand, ScaleLinear } from 'd3-scale';

// Define the graph dimensions and margins
const width = 500;
const height = 500;
const margin = { top: 20, bottom: 20, left: 20, right: 20 };

// Then we'll create some bounds
const xMax = width - margin.left - margin.right;
const yMax = height - margin.top - margin.bottom;


// Finally we'll embed it all in an SVG
export const BarGraph = ({data}: {data: AilmentHistogram}) => {
    const ailments = Object.keys(data)
    // We'll make some helpers to get at the data we want
    const x = (d: string) => d;
    const y = (d: string) => +data[d] * 100;

    // And then scale the graph by our data
    const xScale = scaleBand({
    range: [0, xMax],
    round: true,
    domain: ailments.map(x),
    padding: 0.4,
    });
    const yScale = scaleLinear({
    range: [yMax, 0],
    round: true,
    domain: [0, Math.max(...ailments.map(y))],
    });

    // Compose together the scale and accessor functions to get point functions
    const compose = (
        scale: ScaleBand<string>|ScaleLinear<number,number,number>, accessor: any
        ) => (data: string) => scale(accessor(data));
    const xPoint = compose(xScale, x);
    const yPoint = compose(yScale, y);

    return (
    <svg width={width} height={height}>
        {ailments.map((d, i) => {
        const barHeight = yMax - yPoint(d)!;
        return (
            <Group key={`bar-${i}`}>
            <Bar
                x={xPoint(d)}
                y={yMax - barHeight}
                height={barHeight}
                width={xScale.bandwidth()}
                fill="#fc2e1c"
            />
            </Group>
        );
        })}
    </svg>
    );
}
