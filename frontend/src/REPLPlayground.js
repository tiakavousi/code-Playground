import React, { useState, useEffect, useRef } from 'react';

const REPLPlayground = ({ wsUrl }) => {
    const [code, setCode] = useState('');
    const [language, setLanguage] = useState('Python');
    const [output, setOutput] = useState('');
    const [isRunning, setIsRunning] = useState(false);
    const [input, setInput] = useState('');
    const [error, setError] = useState(null);
    const [isDarkMode, setIsDarkMode] = useState(false);
    const [accentColor, setAccentColor] = useState('#007bff');
    const outputRef = useRef(null);
    const ws = useRef(null);

    // useEffect(() => {
    //     return () => {
    //         if (ws.current) {
    //             ws.current.close();
    //         }
    //     };
    // }, []);

    useEffect(() => {
        document.documentElement.style.setProperty('--accent-color', accentColor);
    }, [accentColor]);

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

    const toggleTheme = () => {
        setIsDarkMode(!isDarkMode);
    };

    const handleColorChange = (e) => {
        setAccentColor(e.target.value);
    };

    return (
        <div className={`repl-container ${isDarkMode ? 'dark-mode' : ''}`}>
            <div className="repl-header">
                <h1 className="repl-title">Interactive REPL Playground</h1>
                <div className="repl-controls">
                    <button onClick={toggleTheme} className="repl-theme-toggle">
                        {isDarkMode ? '☀️ Light' : '🌙 Dark'}
                    </button>
                    <div className="repl-color-picker">
                        <label htmlFor="colorPicker" className="repl-color-picker-label">Accent:</label>
                        <input
                            type="color"
                            id="colorPicker"
                            value={accentColor}
                            onChange={handleColorChange}
                        />
                    </div>
                </div>
            </div>
            {error && (
                <div style={{ color: 'red', marginBottom: '20px' }}>
                    {error}
                </div>
            )}
            <select
                value={language}
                onChange={(e) => setLanguage(e.target.value)}
                className="repl-select"
            >
                <option value="c">C</option>
                <option value="cpp">C++</option>
                <option value="python3">Python</option>
                <option value="javascript">JavaScript</option>
                <option value="java">Java</option>
                <option value="bash">Bash</option>
            </select>
            <textarea
                value={code}
                onChange={(e) => setCode(e.target.value)}
                placeholder="Enter your code here..."
                rows={10}
                className="repl-textarea"
            />
            <button
                onClick={handleExecute}
                disabled={isRunning}
                className="repl-button"
            >
                {isRunning ? 'Running...' : 'Run'}
            </button>
            <div>
                <h2 style={{ color: accentColor }}>Output:</h2>
                <pre ref={outputRef} className="repl-output">
                    {output}
                </pre>
            </div>
            <form onSubmit={handleInputSubmit} className="repl-input-form">
                <input
                    type="text"
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    placeholder="Enter input here..."
                    className="repl-input"
                />
                <button
                    type="submit"
                    disabled={!isRunning}
                    className="repl-submit-button"
                >
                    Send Input
                </button>
            </form>
        </div>
    );
};

export default REPLPlayground;