import React from 'react';
import './REPLOutput.css';

const REPLOutput = ({ output, outputRef, isDarkMode }) => {
    return (
        <div className={`repl-output-container ${isDarkMode ? 'dark-mode' : ''}`}>
                <pre 
                    ref={outputRef} 
                    className="repl-output-content"
                >
                    {output}
                </pre>
        </div>
    );
};

export default REPLOutput;