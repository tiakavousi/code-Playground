import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';

const wsUrl = process.env.REACT_APP_WS_URL || 'localhost:8080';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <App wsUrl={wsUrl} />
  </React.StrictMode>
);
