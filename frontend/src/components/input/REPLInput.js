import React from 'react';
import './REPLInput.css';

const REPLInput = ({ input, setInput, isRunning, handleInputSubmit }) => {
    return (
        <div>
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
                className="repl-button"
            >
                Send
            </button>
            
        </form>
        <div className="instruction">
                Please enter your input and click "Send Input" to submit.
            </div>
        </div>
    );
};

export default REPLInput;