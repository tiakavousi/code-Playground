import React from 'react';
import { HiSun, HiMoon } from 'react-icons/hi';
import './Header.css';

const Header = ({ isDarkMode, toggleTheme, accentColor, handleColorChange }) => {
    return (
        <div className="row repl-header">
            <div className='column repl-title'>
                <h1 className='title'>Code Playground</h1>
            </div>
            <div className="column repl-controls">
                <button onClick={toggleTheme} className="theme-toggle-switch" aria-label="Toggle dark mode">
                    <div className={`switch-track ${isDarkMode ? 'dark' : ''}`}>
                        <HiSun className="sun-icon" size={20} />
                        <HiMoon className="moon-icon" size={20} />
                        <div className="switch-thumb" />
                    </div>
                </button>
            </div>
        </div>
    );
};

export default Header;
