import React, { useState } from 'react';
import logo from './logo.svg';
import './App.css';
import { SymptomBox } from './SymptomBox';
import { AilmentBox } from './AilmentBox';
import { HPO } from './types';

function App() {
  const api: string = process.env.REACT_APP_API || "http://localhost:8081"
  const [showAilments, setShowAilments] = useState(false);
  const [selectedSymptoms, setselectedSymptoms] = useState<HPO[]>([]);
  const propagateSelectedList = (selected: HPO[]) => {
    setselectedSymptoms(selected)
  }
  return (
    <div className="App">
      <header className="App-header">
        <h1>Symptom Checker</h1>
        <div id='wrapper'>
          <div className='symptoms'>{<SymptomBox api={`${api}`} propagateSelectedList={propagateSelectedList}/>}</div>
          <div className='ailments'>
            <button className='button' onClick={()=>{setShowAilments(true)}}></button>
            {showAilments && <AilmentBox api={`${api}`} hpoList={selectedSymptoms}/>}
          </div>
        </div>
      </header>
    </div>
  );
}

export default App;
