import React from 'react';
import { HiSun, HiMoon } from 'react-icons/hi';
import './Header.css';

const Header = ({ isDarkMode, toggleTheme }) => {
    return (
        <div className="row repl-header">
            <div className='column repl-title'>
                <h3 className='title'>Code Playground</h3>
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
