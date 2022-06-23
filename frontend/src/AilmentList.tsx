import React from "react"
import { AilmentHistogram } from "./types";

export const AilmentList = ({data}: {data: AilmentHistogram}) => (
    <div>
        <ol>
        {data.map((item, idx) => {
            return <li className="item" key={`${item.id}-${idx}`}>{item.name}</li>;
        })}
        </ol>
    </div>
  );