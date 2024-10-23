import React from 'react';
import { HiSun, HiMoon } from 'react-icons/hi';
import './Header.css';

const Header = ({ isDarkMode, toggleTheme, accentColor, handleColorChange }) => {
    return (
        <div className="repl-header">
            <h1 className="repl-title">Playground</h1>
            <div className="repl-controls">
                <button 
                    onClick={toggleTheme}
                    className="theme-toggle-switch"
                    aria-label="Toggle dark mode"
                >
                    <div className={`switch-track ${isDarkMode ? 'dark' : ''}`}>
                    <HiSun className="sun-icon" size={16} />
                    <HiMoon className="moon-icon" size={16} />
                    <div className="switch-thumb" />
                    </div>
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
