import { Activity, Zap, Moon, Sun } from 'lucide-react';
import { getConnectionIndicator, getModeBadgeClass } from '../utils/Format';
import { useTheme } from '../Context/Theme';
const Header = ({ connectionStatus, systemMode }) => {
  const indicator = getConnectionIndicator(connectionStatus);
  const { theme, toggleTheme } = useTheme();

  return (
    <header className="card p-4 mb-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="w-10 h-10 bg-gradient-to-br from-blue-600 to-blue-700 rounded-lg flex items-center justify-center">
            <Zap className="w-6 h-6 text-white" />
          </div>
          <div>
            <h1 className="text-2xl font-bold" style={{ color: 'var(--text-primary)' }}>
              NanoPulse Core
            </h1>
            <p className="text-sm" style={{ color: 'var(--text-secondary)' }}>
              Low-Latency Trading Engine
            </p>
          </div>
        </div>

        <div className="flex items-center space-x-6">
          <div className="flex items-center space-x-2">
            <Activity className="w-4 h-4" style={{ color: 'var(--text-secondary)' }} />
            <div className={indicator.class}></div>
            <span className="text-sm" style={{ color: 'var(--text-secondary)' }}>
              {indicator.text}
            </span>
          </div>

          <div className={getModeBadgeClass(systemMode)}>
            {systemMode}
          </div>

          <button
            onClick={toggleTheme}
            className="p-2 rounded-lg transition-all duration-200 hover:scale-110"
            style={{
              backgroundColor: 'var(--bg-secondary)',
              border: '1px solid var(--border-color)'
            }}
            title={`Switch to ${theme === 'dark' ? 'light' : 'dark'} mode`}
          >
            {theme === 'dark' ? (
              <Sun className="w-5 h-5 text-yellow-400" />
            ) : (
              <Moon className="w-5 h-5 text-blue-600" />
            )}
          </button>
        </div>
      </div>
    </header>
  );
};

export default Header;