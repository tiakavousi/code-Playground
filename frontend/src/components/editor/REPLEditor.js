import React, { useState, useEffect } from 'react';
import './REPLEditor.css';
import Editor from "@monaco-editor/react";

const REPLEditor = ({
    language,
    setLanguage,
    code,
    handleEditorChange,
    handleEditorDidMount,
    handleEditorValidation,
    isDarkMode,
    isRunning,
    handleExecute,
    handleSaveAndShare,
    shareLink}) => {

        const languageMap = {
            c: {
                id: "c",
                placeholder: "// write code here ..."
            },
            cpp: {
                id: "cpp",
                placeholder: "// write code here ..."
            },
            python: {
                id: "python",
                placeholder: "# write code here ..."
            },
            javascript: {
                id: "javascript",
                placeholder: "// write code here ..."
            },
            java: {
                id: "java",
                placeholder: "// write code here ..."
            },
            bash: {
                id: "bash",
                placeholder: "# write code here ..."
            }
        };
    
        const [currentCode, setCurrentCode] = useState(languageMap[language]?.placeholder || "");
    
        // Update placeholder when the language changes
        useEffect(() => {
            setCurrentCode(languageMap[language]?.placeholder || "");
        }, [language]);
    
    return(
        <>
            <div className="repl-editor">
                <div className="editor-controls">
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

                    <button
                        onClick={handleExecute}
                        disabled={isRunning}
                        className="repl-button"
                    >
                        {isRunning ? 'Running...' : 'Run'}
                    </button>
                </div>
                <div className="monaco-editor-container">
                    <Editor
                        height="100%"
                        // language={language}
                        language={languageMap[language]?.id}
                        defaultValue={languageMap[language]?.placeholder}
                        value={currentCode}
                        // value={code}
                        theme={isDarkMode ? "vs-dark" : "light"}
                        // onChange={handleEditorChange}
                        onChange={(newValue) => {
                            console.log("editor value changing: ", newValue);
                            setCurrentCode(newValue); // Update local code state
                            handleEditorChange(newValue); // Call external handler
                        }}
                        onMount={handleEditorDidMount}
                        onValidate={handleEditorValidation}
                        options={{
                            minimap: { enabled: false }
                        }}
                        loading={<div>Loading editor...</div>}
                    />
                </div>

                <div className="bottom-controls">
                    
                    <button
                        onClick={handleSaveAndShare}
                        className="repl-button"
                    >
                        Share
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
        </>
    ); 
};

export default REPLEditor;