import Fuse from "fuse.js";
import React, { useEffect, useState } from "react"
import Item from "./Item";
import { HPO, Symptom } from "./types";


export const SymptomBox = ({api, propagateSelectedList}: {api: string, propagateSelectedList: any}) => {
    // Loading data
    const [data, setData] = useState<Symptom[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [searchData, setSearchData] = useState(data);
    const endpoint = `${api}/symptoms`
    const getData = async () => {
        try {
            const response = await fetch(endpoint);
            if (!response.ok) {
            throw new Error(
                `This is an HTTP error: The status is ${response.status}`
            );
            }
            let actualData: Symptom[] = await response.json();
            setData(actualData);
            setSearchData(actualData);
            setError(null);
        } catch(err: any) {
            setError(err.message);
            setData([]);
        } finally {
            setLoading(false);
        }  
    }
    // fuzzy search through symptoms
    const searchItem = (query: string) => {
        if (!query) {
        setSearchData(data);
        return;
        }
        query = query.toLowerCase();

        const fuse = new Fuse(data, { 
            keys: ["name"]    
        });
        const fuzzyResult = fuse.search(query);
        const finalResult: Symptom[] = [];
        if (fuzzyResult.length) {
            fuzzyResult.forEach((fuzzy) => {
              finalResult.push(fuzzy.item);
            });
            setSearchData(finalResult);
          } else {
            setSearchData([]);
          }
        setSearchData(finalResult);
    };
    // handle selected items
    const [selected, setSelected] = useState<HPO[]>([]);
    const handleSelection = (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
        const box: HTMLButtonElement = event.currentTarget;
        const hpo = box.id
        if (box.className.includes("selected")) {
            box.className = box.className.split(" ")[0]
            const index = selected.indexOf(hpo);
            if (index > -1) {
                selected.splice(index, 1); // 2nd parameter means remove one item only
            }
        } else {
            box.className = box.className + " selected"
            if (!(hpo in selected)) {
                selected.push(hpo)
            }
        }
        console.debug(selected)
        setSelected(selected)
        propagateSelectedList(selected)
    }

    const resetSelection = (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
        const allSelected = Array.from(
            document.querySelectorAll('div.item.selected')
          );
        
          allSelected.forEach(element=>{
            element.className = "item"
          })

        setSelected([])
        propagateSelectedList([])
    }

    useEffect(() => {
        getData()
    }, [])

    return (
        <div className="App">
          {loading && <div>A moment please...</div>}
          {error && (
            <div>{`There is a problem fetching the post data - ${error} - try reloading`}</div>
          )}
          <div>
            <button id="invisibutton" onClick={resetSelection}></button>
            <p className="title">Search your symptoms and then click to select</p>
            <div className="search-container">
                <input
                type="search"
                onChange={(e) => searchItem(e.target.value)}
                placeholder="What symptom do you have?"
                />
            </div>

            <div className="item-container scrollable">
                {data && searchData.map((symptom: Symptom) => (
                <Item {...symptom} onClick={handleSelection} key={symptom.hpo} />
                ))}
            </div>
         </div>
        </div>
      );
}
