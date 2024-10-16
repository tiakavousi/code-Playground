import React, { useState } from 'react';

const REPLPlayground = () => {
    const [code, setCode] = useState('');
    const [language, setLanguage] = useState('python');
    const [output, setOutput] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const handleExecute = async () => {
        setIsLoading(true);
        setOutput('');
        try {
            console.log('Sending request to:', '/execute');
            console.log('Request body:', { language, code });

            const response = await fetch('/execute', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ language, code }),
            });

            console.log('Response status:', response.status);
            console.log('Response headers:', Object.fromEntries(response.headers.entries()));

            const text = await response.text();
            console.log('Response text:', text);

            let data;
            try {
                data = JSON.parse(text);
            } catch (error) {
                console.error('Error parsing JSON:', error);
                setOutput(`Error: Invalid JSON response: ${text.substring(0, 100)}...`);
                return;
            }

            if (!response.ok) {
                setOutput(`Error: ${data.error || 'An error occurred during execution'}`);
            } else {
                setOutput(data.output || 'No output received');
            }
        } catch (error) {
            console.error('Execution error:', error);
            setOutput(`Error: ${error.message}`);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div style={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
            <h1 style={{ textAlign: 'center' }}>REPL Playground</h1>
            <div style={{ marginBottom: '20px' }}>
                <select
                    value={language}
                    onChange={(e) => setLanguage(e.target.value)}
                    style={{ width: '100%', padding: '10px' }}
                >
                    <option value="python">Python</option>
                    <option value="javascript">JavaScript</option>
                    <option value="bash">Bash</option>
                    <option value="java">Java</option>
                    <option value="c">C</option>
                    <option value="cpp">C++</option>
                </select>
            </div>
            <textarea
                value={code}
                onChange={(e) => setCode(e.target.value)}
                placeholder="Enter your code here..."
                rows={10}
                style={{ width: '100%', marginBottom: '20px', padding: '10px' }}
            />
            <button
                onClick={handleExecute}
                disabled={isLoading}
                style={{
                    width: '100%',
                    padding: '10px',
                    backgroundColor: isLoading ? '#cccccc' : '#007bff',
                    color: 'white',
                    border: 'none',
                    cursor: isLoading ? 'not-allowed' : 'pointer'
                }}
            >
                {isLoading ? 'Executing...' : 'Execute'}
            </button>
            <div style={{ marginTop: '20px' }}>
                <h2>Output:</h2>
                <pre style={{ backgroundColor: '#f0f0f0', padding: '10px', whiteSpace: 'pre-wrap', wordWrap: 'break-word' }}>
                    {output || 'No output yet'}
                </pre>
            </div>
        </div>
    );
};

export default REPLPlayground;