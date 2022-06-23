import React, { useState } from "react"
import { AilmentHistogram } from "./types";

export const AilmentList = ({data}: {data: AilmentHistogram}) => {
    const toggleHidden = (event: React.MouseEvent<HTMLLIElement>) => {
        // event.preventDefault();
        const box: HTMLLIElement = event.currentTarget;
        const inner = box.children[0]
        inner.toggleAttribute("hidden")
    }
    return (
        <div>
            <ol>
            {data.map((item, idx) => {
                const symptoms = new Set(item.symptoms.map(s => {return s.name}))
                const otherSymptoms = Array.from(symptoms)
                return (
                    <li key={`${item.id}-${idx}`} className="nice-list item item-container" onClick={toggleHidden}> {item.name}
                        <div className="detail-box scrollable" hidden>
                            <a href={item.expert} target="_blank" rel="noopener noreferrer">Consult an expert on the matter</a>
                            <ul>
                                <p className="detail">Other symptoms may include:</p>
                                {otherSymptoms.map((symp, ix) => {
                                    return <li key={ix}>{symp}</li>
                                })}
                            </ul>
                        </div>
                    </li>
                )
            })}
            </ol>
        </div>
    )};