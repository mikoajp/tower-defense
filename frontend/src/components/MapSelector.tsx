import { useState, useEffect } from 'react';
import { API_URL } from '../config';
import './MapSelector.css';

interface Map {
  id: string;
  name: string;
  difficulty: string;
  description: string;
  pathLength: number;
}

interface MapSelectorProps {
  onMapChange?: (mapId: string, mapName: string) => void;
}

export default function MapSelector({ onMapChange }: MapSelectorProps) {
  const [maps, setMaps] = useState<Map[]>([]);
  const [selectedMap, setSelectedMap] = useState<string>('classic');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchMaps();
  }, []);

  const fetchMaps = async () => {
    try {
      console.log('Fetching maps from:', `${API_URL}/maps`);
      const response = await fetch(`${API_URL}/maps`);
      console.log('Maps response status:', response.status);
      const data = await response.json();
      console.log('Maps data received:', data);
      setMaps(data.maps || []);
      console.log('Maps set to state:', data.maps?.length || 0);
    } catch (error) {
      console.error('Error fetching maps:', error);
    }
  };

  const handleMapChange = async (mapId: string) => {
    setLoading(true);
    try {
      const response = await fetch(`${API_URL}/map`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ mapId }),
      });

      if (response.ok) {
        setSelectedMap(mapId);
        const map = maps.find((m) => m.id === mapId);
        onMapChange?.(mapId, map?.name || mapId);
      }
    } catch (error) {
      console.error('Error changing map:', error);
    } finally {
      setLoading(false);
    }
  };

  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty.toLowerCase()) {
      case 'easy':
        return '#32cd32';
      case 'medium':
        return '#ffa500';
      case 'hard':
        return '#dc143c';
      case 'expert':
        return '#8b008b';
      default:
        return '#4169e1';
    }
  };

  return (
    <div className="map-selector">
      <h3 className="map-selector-title">
        <span className="map-icon">🗺️</span>
        Select Map
      </h3>

      <div className="map-grid">
        {maps.length === 0 ? (
          <div className="map-loading">
            <span>Loading maps...</span>
          </div>
        ) : (
          maps.map((map) => (
            <div
              key={map.id}
              className={`map-card ${selectedMap === map.id ? 'selected' : ''} ${
                loading ? 'disabled' : ''
              }`}
              onClick={() => !loading && handleMapChange(map.id)}
            >
              {selectedMap === map.id && (
                <div className="map-selected-indicator">✓</div>
              )}

              <div className="map-card-header">
                <div className="map-card-name">{map.name}</div>
                <div
                  className="map-card-difficulty"
                  style={{ color: getDifficultyColor(map.difficulty) }}
                >
                  {map.difficulty}
                </div>
              </div>

              <div className="map-card-description">{map.description}</div>

              <div className="map-card-stats">
                <div className="map-stat">
                  <span className="map-stat-icon">📍</span>
                  <span className="map-stat-value">{map.pathLength} points</span>
                </div>
              </div>
            </div>
          ))
        )}
      </div>

      {loading && (
        <div className="map-loading">
          <div className="loading-spinner"></div>
          <span>Loading map...</span>
        </div>
      )}
    </div>
  );
}
