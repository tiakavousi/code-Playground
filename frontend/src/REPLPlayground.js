import React, { useState, useEffect, useRef } from 'react';

const REPLPlayground = ({ wsUrl, initialCode = null }) => {
    const [code, setCode] = useState('');
    const [language, setLanguage] = useState('Python');
    const [output, setOutput] = useState('');
    const [isRunning, setIsRunning] = useState(false);
    const [input, setInput] = useState('');
    const [error, setError] = useState(null);
    const [isDarkMode, setIsDarkMode] = useState(false);
    const [accentColor, setAccentColor] = useState('#007bff');
    const [shareLink, setShareLink] = useState('');
    const outputRef = useRef(null);
    const ws = useRef(null);


    useEffect(() => {
        document.documentElement.style.setProperty('--accent-color', accentColor);
    }, [accentColor]);

    useEffect(() => {
        if (initialCode) {
            setCode(initialCode.code || '');
            setLanguage(initialCode.language || 'Python');
        }
    }, [initialCode]);

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

    const handleSaveAndShare = async () => {
        try {
            const response = await fetch(`http://${wsUrl}/save`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ language, code }),
            });

            if (!response.ok) {
                throw new Error('Failed to save code');
            }

            const data = await response.json();
            setShareLink(`http://${window.location.host}/share/${data.id}`);
        } catch (err) {
            setError(`Error: ${err.message}`);
        }
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
                <div className="repl-error">
                    {error}
                </div>
            )}

            <div className="repl-editor">
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

                <div className="repl-button-group">
                    <button
                        onClick={handleExecute}
                        disabled={isRunning}
                        className="repl-button"
                    >
                        {isRunning ? 'Running...' : 'Run'}
                    </button>
                    <button
                        onClick={handleSaveAndShare}
                        className="repl-button"
                    >
                        Save & Share
                    </button>
                </div>
            </div>

            {shareLink && (
                <div className="share-link-container">
                    <p>Share your code with this link:</p>
                    <input
                        type="text"
                        value={shareLink}
                        readOnly
                        className="share-link-input"
                    />
                </div>
            )}

            <div className="repl-output-container">
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