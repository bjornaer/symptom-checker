import React, { useState } from 'react';
import logo from './logo.svg';
import './App.css';
import { SymptomBox } from './SymptomBox';
import { AilmentBox } from './AilmentBox';
import { HPO } from './types';

function App() {
  const api: string = process.env.REACT_APP_API || "/api"
  const [showAilments, setShowAilments] = useState(false);
  const [selectedSymptoms, setselectedSymptoms] = useState<HPO[]>([]);
  const propagateSelectedList = (selected: HPO[]) => {
    setselectedSymptoms(selected)
  }
  const triggerSelection = () => {
    if (selectedSymptoms.length) {
      setShowAilments(true)
    } else {
      alert("Please select what symptoms you have!")
    }
  }
  const triggerSelectionReset = () => {
    setShowAilments(false)
    const btn = document.getElementById('invisibutton')
    btn && btn.click()
  }
  return (
    <div className="App">
      <header className="App-header">
        <h1>Symptom Checker</h1>
        {!showAilments && <button className="button-64 confirm" role="button" onClick={()=>{triggerSelection()}}><span className="text">What Do I Have?</span></button>}
        {showAilments && <button className="button-64 clear" role="button" onClick={()=>{triggerSelectionReset()}}><span className="text">Clear</span></button>}
        <div id='wrapper'>
          <div className='symptoms'>{<SymptomBox api={`${api}`} propagateSelectedList={propagateSelectedList}/>}</div>
          <div className='ailments'>
            {showAilments && <AilmentBox api={`${api}`} hpoList={selectedSymptoms}/>}
          </div>
        </div>
      </header>
    </div>
  );
}

export default App;
