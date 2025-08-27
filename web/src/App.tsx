import { useQuery } from "@tanstack/react-query";
import React from "react";
import { App as AppComponent } from "~/components/App";
import { Button } from "~/components/Button";
import { Input } from "~/components/Input";
import { Message } from "~/components/Message";
import { Review } from "~/components/Review";
import { isHttpError } from "~/lib/http/http-error";
import { addApp, fetchApps } from "~/lib/server/apps";
import { fetchReviews } from "~/lib/server/reviews";
import "./App.css";

function App() {
  const [selectedAppId, setSelectedAppId] = React.useState("");
  const [newAppId, setNewAppId] = React.useState("");
  const [isAddingApp, setIsAddingApp] = React.useState(false);
  const [addError, setAddError] = React.useState<string | null>(null);

  const {
    data: appsResponse,
    isLoading: appsLoading,
    error: appsError,
    refetch: refetchApps,
  } = useQuery({
    queryKey: ["apps"],
    queryFn: fetchApps,
    staleTime: 2 * 60 * 1000, // 2 minutes
    retry: false,
  });

  const {
    data: reviews,
    error: reviewsError,
    refetch: refetchReviews,
  } = useQuery({
    queryKey: ["reviews", selectedAppId],
    queryFn: () => fetchReviews(selectedAppId),
    enabled: !!selectedAppId,
    staleTime: 5 * 60 * 1000, // 5 minutes
    retry: false,
  });

  const handleAppSelect = (appId: string) => {
    setSelectedAppId(appId);
  };

  const handleAddApp = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setAddError(null);

    if (!newAppId.trim()) return;

    setIsAddingApp(true);
    try {
      await addApp(newAppId.trim());
      setNewAppId("");
      refetchApps(); // Refresh the apps list
    } catch (error) {
      if (isHttpError(error)) {
        if (error.status === 400) {
          setAddError(
            "Invalid App ID. Please enter a valid numeric App Store ID."
          );
        } else if (error.status === 404) {
          setAddError("App not found. Please check the App ID and try again.");
        } else {
          setAddError("Failed to add app. Please try again.");
        }
      } else {
        setAddError("An unexpected error occurred. Please try again.");
      }
    } finally {
      setIsAddingApp(false);
    }
  };

  const apps = appsResponse?.data || [];
  const reviewsErrorMessage = getErrorMessage(reviewsError);

  return (
    <div className="app">
      <div className="reviews-container">
        <div className="header">
          <h1>Apps Reviews</h1>
          <p className="subtitle">Select an app to view its recent reviews</p>
        </div>

        <div className="apps-list">
          <div className="apps-list__header">
            <h2>Select an App {apps.length > 0 && `(${apps.length})`}</h2>
            <p className="apps-list__subtitle">
              Choose an app to view its recent reviews
            </p>
          </div>

          {appsLoading ? (
            <div className="apps-list__loading">Loading apps...</div>
          ) : appsError ? (
            <div className="apps-list__error">
              <Message type="error">
                <Message.Title>Error</Message.Title>
                <span>
                  {isHttpError(appsError) && appsError.status === 404
                    ? "No apps found. Add your first app below."
                    : "Failed to load apps. Please try again."}
                </span>
                <Button onClick={() => refetchApps()} variant="link">
                  Try Again
                </Button>
              </Message>
            </div>
          ) : apps.length === 0 ? (
            <div className="apps-list__empty">
              <p>No apps available. Add your first app below.</p>
            </div>
          ) : (
            <div className="apps-list__grid">
              {apps.map((app) => (
                <AppComponent
                  key={app.id}
                  app={app}
                  onSelect={handleAppSelect}
                  isSelected={selectedAppId === app.id}
                />
              ))}
            </div>
          )}

          <div className="apps-list__add-app">
            <h3>Add New App</h3>
            <form onSubmit={handleAddApp} className="add-app-form">
              <div className="input-group">
                <Input
                  type="text"
                  value={newAppId}
                  onChange={(e) => setNewAppId(e.target.value)}
                  placeholder="Enter App Store App ID (e.g., 1234567890)"
                  disabled={isAddingApp}
                  required
                />
                <Button
                  type="submit"
                  disabled={isAddingApp || !newAppId.trim()}
                >
                  {isAddingApp ? "Adding..." : "Add App"}
                </Button>
              </div>
            </form>
            {addError && (
              <div className="add-app-form__error">
                <Message type="error">
                  <span>{addError}</span>
                </Message>
              </div>
            )}
          </div>
        </div>

        {reviewsError && (
          <div className="message-container">
            <Message type="error">
              <Message.Title>{reviewsErrorMessage.title}</Message.Title>
              <span>{reviewsErrorMessage.description}</span>
              <Button onClick={() => refetchReviews()} variant="link">
                Try Again
              </Button>
            </Message>
          </div>
        )}

        {reviews && (
          <div className="reviews">
            <div className="reviews__header">
              <h2>Recent Reviews ({reviews.data.length})</h2>
              <p className="reviews__subtitle">
                Reviews from the last 48 hours
              </p>
            </div>

            {reviews.data.length === 0 ? (
              <div className="reviews__no-reviews">
                <p>No recent reviews found for this app.</p>
              </div>
            ) : (
              <div className="reviews__list">
                {reviews.data.map((review) => (
                  <Review key={review.id} review={review} />
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

const getErrorMessage = (error: unknown) => {
  if (isHttpError(error)) {
    if (error.status === 404) {
      return {
        title: "App not found",
        description: "The app you are looking for does not exist.",
      };
    }
  }

  return {
    title: "An unexpected error occurred",
    description: "Please try again later.",
  };
};

export default App;
