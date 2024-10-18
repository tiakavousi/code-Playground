import React, { useState, useEffect, useRef } from 'react';

const REPLPlayground = ({ wsUrl }) => {
    const [code, setCode] = useState('');
    const [language, setLanguage] = useState('python');
    const [output, setOutput] = useState('');
    const [isRunning, setIsRunning] = useState(false);
    const [input, setInput] = useState('');
    const [error, setError] = useState(null);
    const outputRef = useRef(null);
    const ws = useRef(null);

    useEffect(() => {
        return () => {
            if (ws.current) {
                ws.current.close();
            }
        };
    }, []);

    const handleExecute = async () => {
        setIsRunning(true);
        setOutput('');
        setError(null);

        try {
            ws.current = new WebSocket(`ws://${wsUrl}/execute`);

            ws.current.onopen = () => {
                console.log('WebSocket connection established');
                ws.current.send(JSON.stringify({ language, code }));
            };

            ws.current.onmessage = (event) => {
                setOutput(prev => prev + event.data + '\n');
                if (outputRef.current) {
                    outputRef.current.scrollTop = outputRef.current.scrollHeight;
                }
            };

            ws.current.onclose = () => {
                console.log('WebSocket connection closed');
                setIsRunning(false);
            };

            ws.current.onerror = (event) => {
                console.error('WebSocket error:', event);
                setError('WebSocket error occurred. Check console for details.');
                setIsRunning(false);
            };
        } catch (err) {
            console.error('Error setting up WebSocket:', err);
            setError(`Error: ${err.message}`);
            setIsRunning(false);
        }
    };

    const handleInputSubmit = (e) => {
        e.preventDefault();
        if (ws.current && ws.current.readyState === WebSocket.OPEN) {
            ws.current.send(input);
            setInput('');
        } else {
            setError('WebSocket is not connected');
        }
    };

    return (
        <div style={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
            <h1 style={{ textAlign: 'center' }}>Interactive REPL Playground</h1>
            {error && (
                <div style={{ color: 'red', marginBottom: '20px' }}>
                    {error}
                </div>
            )}
            <div style={{ marginBottom: '20px' }}>
                <select
                    value={language}
                    onChange={(e) => setLanguage(e.target.value)}
                    style={{ width: '100%', padding: '10px' }}
                >
                    <option value="c">C</option>
                    <option value="cpp">C++</option>
                    <option value="python">Python</option>
                    <option value="javascript">JavaScript</option>
                    <option value="java">Java</option>
                    <option value="bash">Bash</option>
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
                disabled={isRunning}
                style={{
                    width: '100%',
                    padding: '10px',
                    backgroundColor: isRunning ? '#cccccc' : '#007bff',
                    color: 'white',
                    border: 'none',
                    cursor: isRunning ? 'not-allowed' : 'pointer'
                }}
            >
                {isRunning ? 'Running...' : 'Run'}
            </button>
            <div style={{ marginTop: '20px' }}>
                <h2>Output:</h2>
                <pre
                    ref={outputRef}
                    style={{
                        backgroundColor: '#f0f0f0',
                        padding: '10px',
                        whiteSpace: 'pre-wrap',
                        wordWrap: 'break-word',
                        height: '200px',
                        overflowY: 'auto'
                    }}
                >
                    {output}
                </pre>
            </div>
            <form onSubmit={handleInputSubmit} style={{ marginTop: '20px' }}>
                <input
                    type="text"
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    placeholder="Enter input here..."
                    style={{ width: '70%', padding: '10px' }}
                />
                <button
                    type="submit"
                    disabled={!isRunning}
                    style={{
                        width: '30%',
                        padding: '10px',
                        backgroundColor: isRunning ? '#28a745' : '#cccccc',
                        color: 'white',
                        border: 'none',
                        cursor: isRunning ? 'pointer' : 'not-allowed'
                    }}
                >
                    Send Input
                </button>
            </form>
        </div>
    );
};

export default REPLPlayground;