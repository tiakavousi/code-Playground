import React from 'react';
import './REPLOutput.css';

const REPLOutput = ({ output, outputRef, accentColor }) => {
    return (
        <div className="repl-output-container">
            <h2 style={{ color: accentColor }} className="repl-output-title">
                Output:
            </h2>
            <pre 
                ref={outputRef} 
                className="repl-output-content"
            >
                {output || 'No output yet. Run your code to see results.'}
            </pre>
        </div>
    );
};

export default REPLOutput;