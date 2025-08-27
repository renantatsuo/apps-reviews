import React from "react";
import { type App as AppType } from "~/lib/server/apps";
import "./App.css";

interface AppProps {
  app: AppType;
  onSelect: (appId: string) => void;
  isSelected?: boolean;
}

export const App: React.FC<AppProps> = ({
  app,
  onSelect,
  isSelected = false,
}) => {
  const handleClick = () => {
    onSelect(app.id);
  };

  return (
    <div
      className={`app-item ${isSelected ? "app-item--selected" : ""}`}
      onClick={handleClick}
    >
      <div className="app-item__thumbnail">
        <img
          src={app.thumbnail_url}
          alt={app.name}
          className="app-item__image"
        />
      </div>
      <div className="app-item__content">
        <h3 className="app-item__name">{app.name}</h3>
        <p className="app-item__id">ID: {app.id}</p>
      </div>
    </div>
  );
};
