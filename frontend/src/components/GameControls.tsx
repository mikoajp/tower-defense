import { useState } from 'react';
import { API_URL } from '../config';
import './GameControls.css';

interface GameControlsProps {
  onSaveSuccess?: () => void;
  onLoadSuccess?: () => void;
}

export default function GameControls({ onSaveSuccess, onLoadSuccess }: GameControlsProps) {
  const [savedState, setSavedState] = useState<string | null>(null);
  const [message, setMessage] = useState<{ text: string; type: 'success' | 'error' } | null>(null);

  const showMessage = (text: string, type: 'success' | 'error') => {
    setMessage({ text, type });
    setTimeout(() => setMessage(null), 3000);
  };

  const handleSave = async () => {
    try {
      const response = await fetch(`${API_URL}/state`);
      const state = await response.json();
      const stateStr = JSON.stringify(state);
      setSavedState(stateStr);
      
      // Also save to localStorage
      localStorage.setItem('towerDefenseSave', stateStr);
      
      showMessage('Game saved successfully!', 'success');
      onSaveSuccess?.();
    } catch (error) {
      console.error('Error saving game:', error);
      showMessage('Failed to save game', 'error');
    }
  };

  const handleLoad = async () => {
    try {
      let stateToLoad = savedState;
      
      // Try localStorage if no saved state in memory
      if (!stateToLoad) {
        const localSave = localStorage.getItem('towerDefenseSave');
        if (localSave) {
          stateToLoad = localSave;
          setSavedState(localSave);
        }
      }
      
      if (!stateToLoad) {
        showMessage('No saved game found', 'error');
        return;
      }

      const response = await fetch(`${API_URL}/load`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: stateToLoad,
      });

      if (response.ok) {
        showMessage('Game loaded successfully!', 'success');
        onLoadSuccess?.();
      } else {
        showMessage('Failed to load game', 'error');
      }
    } catch (error) {
      console.error('Error loading game:', error);
      showMessage('Failed to load game', 'error');
    }
  };

  const handleReset = async () => {
    if (!confirm('Are you sure you want to reset the game? All progress will be lost.')) {
      return;
    }

    try {
      await fetch(`${API_URL}/reset`, { method: 'POST' });
      showMessage('Game reset successfully!', 'success');
    } catch (error) {
      console.error('Error resetting game:', error);
      showMessage('Failed to reset game', 'error');
    }
  };

  const handleDownload = () => {
    if (!savedState) {
      showMessage('No saved game to download', 'error');
      return;
    }

    const blob = new Blob([savedState], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `tower-defense-save-${Date.now()}.json`;
    a.click();
    URL.revokeObjectURL(url);
    showMessage('Save file downloaded!', 'success');
  };

  const handleUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = async (e) => {
      try {
        const content = e.target?.result as string;
        JSON.parse(content); // Validate JSON
        setSavedState(content);
        
        // Auto-load the uploaded save
        const response = await fetch(`${API_URL}/load`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: content,
        });

        if (response.ok) {
          showMessage('Save file loaded successfully!', 'success');
          onLoadSuccess?.();
        } else {
          showMessage('Failed to load save file', 'error');
        }
      } catch (error) {
        console.error('Error loading save file:', error);
        showMessage('Invalid save file', 'error');
      }
    };
    reader.readAsText(file);
  };

  return (
    <div className="game-controls">
      <h3 className="game-controls-title">
        <span className="control-icon">🎮</span>
        Game Controls
      </h3>
      
      {message && (
        <div className={`control-message ${message.type}`}>
          {message.type === 'success' ? '✓' : '✗'} {message.text}
        </div>
      )}
      
      <div className="control-buttons">
        <button className="control-btn save" onClick={handleSave}>
          <span className="btn-icon">💾</span>
          Save Game
        </button>
        
        <button 
          className="control-btn load" 
          onClick={handleLoad}
          disabled={!savedState && !localStorage.getItem('towerDefenseSave')}
        >
          <span className="btn-icon">📂</span>
          Load Game
        </button>
        
        <button className="control-btn reset" onClick={handleReset}>
          <span className="btn-icon">🔄</span>
          Reset
        </button>
      </div>
      
      <div className="control-file-actions">
        <button className="file-btn download" onClick={handleDownload} disabled={!savedState}>
          <span className="btn-icon">⬇️</span>
          Download Save
        </button>
        
        <label className="file-btn upload">
          <span className="btn-icon">⬆️</span>
          Upload Save
          <input
            type="file"
            accept=".json"
            onChange={handleUpload}
            style={{ display: 'none' }}
          />
        </label>
      </div>
    </div>
  );
}
