import React from 'react';
import './REPLInput.css';

const REPLInput = ({ input, setInput, isRunning,isDarkMode, handleInputSubmit }) => {
    return (
        <div>
          <div className={`instruction ${isDarkMode ? 'dark-mode' : ''}`}>
                <p>If your code takes input, add it in the box below after running and click Send.</p>
            </div>
            <form onSubmit={handleInputSubmit} className="repl-input-form">
                <textarea
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    placeholder="Enter input here..."
                    className={`repl-input ${isDarkMode ? 'dark-mode' : ''}`}
                    rows="3"
                />
               <button
                    type="submit"
                    disabled={!isRunning}
                    className="repl-button"
                >
                    Send
                </button>    
            </form>
        </div>
    );
};

export default REPLInput;