import './ConnectionStatus.css';

interface ConnectionStatusProps {
  connected: boolean;
}

export default function ConnectionStatus({ connected }: ConnectionStatusProps) {
  if (connected) return null;

  return (
    <div className="connection-status">
      <div className="connection-spinner"></div>
      <div className="connection-content">
        <div className="connection-title">⚠️ Connecting to server...</div>
        <div className="connection-subtitle">
          Make sure backend is running on port 8080
        </div>
      </div>
    </div>
  );
}
