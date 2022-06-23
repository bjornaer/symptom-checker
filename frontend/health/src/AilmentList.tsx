import React from "react"
import { AilmentHistogram } from "./types";

export const AilmentList = ({data}: {data: AilmentHistogram}) => (
    <div>
        <h3>Most Likely Ailments</h3>
        <ol>
        {data.map((item, idx) => {
            return <li key={`${item.id}-${idx}`}>{item.name}</li>;
        })}
        </ol>
    </div>
  );