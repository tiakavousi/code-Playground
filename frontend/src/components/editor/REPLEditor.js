import React from 'react';
import Editor from "@monaco-editor/react";

const REPLEditor = ({
    language,
    setLanguage,
    code,
    handleEditorChange,
    handleEditorWillMount,
    handleEditorDidMount,
    handleEditorValidation,
    isDarkMode,
    isRunning,
    handleExecute,
    handleSaveAndShare,
    shareLink}) => {
    return(
        <>
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
            )}
        </>

    )
   
}

export default REPLEditor;