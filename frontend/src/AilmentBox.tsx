import React, { useEffect, useState } from "react"
import { BarGraph } from "./BarGraph";
import { AilmentList } from "./AilmentList";
import { AilmentHistogram } from "./types";

export const AilmentBox = ({api, hpoList}: {api: string, hpoList: string[]}) => {
    // Loading data
    const [data, setData] = useState<AilmentHistogram>();
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const endpoint = `${api}/symptoms`
    const getData = async () => {
        try {
            const response = await fetch(endpoint, {
                method: 'POST',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(hpoList)
            });
            if (!response.ok) {
            throw new Error(
                `This is an HTTP error: The status is ${response.status}`
            );
            }
            let actualData: AilmentHistogram = await response.json();
            setData(actualData);
            setError(null);
        } catch(err: any) {
            setError(err.message);
            setData(undefined);
        } finally {
            setLoading(false);
        }  
    }

    useEffect(() => {
        getData()
    }, [])

    return (
        <div className="App">
          {loading && <div>A moment please...</div>}
          {error && (
            <div>{`There is a problem fetching the post data - ${error}`}</div>
          )}
          {data && <div id="content_container">
          <h3>Most Likely Ailments</h3>
            <div id="list" className="item-container scrollable">
                <AilmentList data={data}/>
            </div>
            {
                // I was going to add a pretty histogram chart to show bars with the ailments and their frequency
                // but the bars aren't rendering and I ended up spending more than I wanted with it
            /* <div id="chart">
                <BarGraph data={data}/>
            </div> */}
            </div>}
        </div>
      );
}
