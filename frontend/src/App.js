import React from 'react';
import REPLPlayground from './REPLPlayground';

function App({ wsUrl }) {
  return (
    <div className="App">
      <REPLPlayground wsUrl={wsUrl} />
    </div>
  );
}

export default App;