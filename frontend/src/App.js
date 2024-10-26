import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Routes, useParams } from 'react-router-dom';
import Main from './components/Main';
import './styles/App.css';

function SharedCodeLoader({ wsUrl }) {
  const [sharedCode, setSharedCode] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const { id } = useParams();

  useEffect(() => {
    if (id) {
      setIsLoading(true);
      fetch(`http://${wsUrl}/share/${id}`)
        .then(response => {
          if (!response.ok) {
            throw new Error('Failed to fetch shared code');
          }
          return response.json();
        })
        .then(data => {
          setSharedCode(data);
          setIsLoading(false);
        })
        .catch(error => {
          console.error('Error fetching shared code:', error);
          setError(error.message);
          setIsLoading(false);
        });
    }
  }, [id, wsUrl]);

  if (isLoading) {
    return <div>Loading shared code...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return <Main wsUrl={wsUrl} initialCode={sharedCode} />;
}

function App({ wsUrl }) {
  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/" element={<Main wsUrl={wsUrl} />} />
          <Route path="/share/:id" element={<SharedCodeLoader wsUrl={wsUrl} />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;