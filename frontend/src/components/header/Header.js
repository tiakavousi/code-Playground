import React from 'react';

const Header = ({ isDarkMode, toggleTheme, accentColor, handleColorChange }) => {
    return (
        <div className="repl-header">
            <h1 className="repl-title">Playground</h1>
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
    );
};

export default Header;
