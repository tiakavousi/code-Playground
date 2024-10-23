import React, { useState, useEffect, useRef } from 'react';
import REPLOutput from './output/REPLOutput';
import REPLInput from './input/REPLInput';
import REPLEditor from './editor/REPLEditor';
import Header from './header/Header';


const Main = ({ wsUrl, initialCode = null }) => {
    const [code, setCode] = useState(''); // Holds the code entered by the user
    const [language, setLanguage] = useState('Python'); // Selected programming language
    const [output, setOutput] = useState(''); // Holds the result/output of the executed code
    const [isRunning, setIsRunning] = useState(false); // Flag indicating if code execution is in progress
    const [input, setInput] = useState('');  // Stores any input to be sent during execution
    const [error, setError] = useState(null); // Handles any error messages
    const [isDarkMode, setIsDarkMode] = useState(false); // Toggle for dark/light mode
    const [accentColor, setAccentColor] = useState('#007bff'); // Customizable accent color
    const [shareLink, setShareLink] = useState(''); // Holds the link to share code
    const outputRef = useRef(null); // Reference to the output area for scrolling
    const ws = useRef(null); // WebSocket connection reference
    const editorRef = useRef(null); // Reference to the Monaco editor instance
    const [isEditorReady, setIsEditorReady] = useState(false); // Tracks if the editor is fully mounted

    // Updates the accent color across the application
    useEffect(() => {
        document.documentElement.style.setProperty('--accent-color', accentColor);
    }, [accentColor]);


    // Sets the initial code and language when component mounts
    useEffect(() => {
        if (initialCode) {
            setCode(initialCode.code || '');
            setLanguage(initialCode.language || 'Python');
        }
    }, [initialCode]);


    // Toggles between dark mode and light mode for the editor and UI
    const toggleTheme = () => {
        setIsDarkMode(!isDarkMode);
    };

   // Triggered when the Monaco editor is mounted and ready to be used
    const handleEditorDidMount = (editor, monaco) => {
        editorRef.current = editor;
        setIsEditorReady(true);
        console.log("Editor mounted successfully");
    }

    // Logs the editor mounting process for debugging
    const handleEditorWillMount = (monaco) => {
        console.log("Editor will mount");
    }

    // Logs any validation markers (warnings/errors) from the Monaco editor
    const handleEditorValidation = (markers) => {
        markers.forEach((marker) => console.log('onValidate:', marker.message));
    }

    // Updates the code state when the user modifies the content in the editor
    const handleEditorChange = (value, event) => {
        setCode(value);
    }

    // Handles the execution of the code by sending it to the server via WebSocket
    const handleExecute = async () => {
        setIsRunning(true);
        setOutput('');
        setError(null);
        console.log('Executing code:', code); // Debug log

        try {
            // Establish WebSocket connection to send and execute the code on the backend
            ws.current = new WebSocket(`ws://${wsUrl}/execute`);
            ws.current.onopen = () => {
                console.log('WebSocket connection established');
                ws.current.send(JSON.stringify({ language, code })); // Sends code and language to the server
            };

            // Receive and append execution output to the output area       
            ws.current.onmessage = (event) => {
                console.log("Received message:", event.data); // Log incoming message
                setOutput(prev => prev + event.data + '\n');
                console.log ("output: " + event.data);
                // Auto-scroll the output to the bottom as new messages arrive
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

    // Handles user input submission to the running code during execution
    const handleInputSubmit = (e) => {
        e.preventDefault();
        if (ws.current && ws.current.readyState === WebSocket.OPEN) {
            ws.current.send(input);
            setInput('');
        } else {
            setError('WebSocket is not connected');
        }
    };

    // Updates the accent color for UI customization
    const handleColorChange = (e) => {
        setAccentColor(e.target.value);
    };

    // Handles saving the current code to the server and generating a shareable link
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
            // Create a shareable link for the saved code
            setShareLink(`http://${window.location.host}/share/${data.id}`);
        } catch (err) {
            setError(`Error: ${err.message}`);
        }
    };

    return (
        <div className={`repl-container ${isDarkMode ? 'dark-mode' : ''}`}>
            {/* <div className="repl-header">
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
            </div> */}
            <Header
                isDarkMode={isDarkMode}
                toggleTheme={toggleTheme}
                accentColor={accentColor}
                handleColorChange={handleColorChange}
            />

            {error && (
                <div className="repl-error">
                    {error}
                </div>
            )}

            {/* <div className="repl-editor">
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

                <Editor
                    height="400px"
                    language={language}
                    value={code}
                    theme={isDarkMode ? "vs-dark" : "light"}
                    onChange={handleEditorChange}
                    beforeMount={handleEditorWillMount}
                    onMount={handleEditorDidMount}
                    onValidate={handleEditorValidation}
                    options={{
                        minimap: { enabled: false }
                    }}
                    loading={<div>Loading editor...</div>}
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
            )} */}
            <REPLEditor 
                language={language}
                setLanguage={setLanguage}
                code={code}
                handleEditorChange={handleEditorChange}
                handleEditorWillMount={handleEditorWillMount}
                handleEditorDidMount={handleEditorDidMount}
                handleEditorValidation={handleEditorValidation}
                isDarkMode={isDarkMode}
                isRunning={isRunning}
                handleExecute={handleExecute}
                handleSaveAndShare={handleSaveAndShare}
                shareLink={shareLink}
            />

            <REPLOutput
                output={output}
                outputRef={outputRef}
                accentColor={accentColor}
            />
              <REPLInput
                input={input}
                setInput={setInput}
                isRunning={isRunning}
                handleInputSubmit={handleInputSubmit}
            />
        </div>
    );
};

export default Main;

