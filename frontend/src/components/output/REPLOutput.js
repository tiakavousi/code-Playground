import React from 'react';
import './REPLOutput.css';

const REPLOutput = ({ output, outputRef, isDarkMode }) => {
    return (
        <div className={`repl-output-container ${isDarkMode ? 'dark-mode' : ''}`}>
                <div 
                    ref={outputRef} 
                    className="repl-output-content"
                >
                    {output || "No Output Yet..."}
                </div>
        </div>
    );
};

export default REPLOutput;