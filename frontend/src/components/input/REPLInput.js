import React from 'react';
import './REPLInputStyles.css';

const REPLInput = ({ input, setInput, isRunning, handleInputSubmit }) => {
    return (
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
    );
};

export default REPLInput;