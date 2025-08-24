import { useQuery } from "@tanstack/react-query";
import React from "react";
import { Button } from "~/components/Button";
import { Input } from "~/components/Input";
import { Message } from "~/components/Message";
import { Review } from "~/components/Review";
import { isHttpError } from "~/lib/http/http-error";
import { fetchReviews } from "~/lib/server/reviews";
import "./App.css";

function App() {
  const [searchAppId, setSearchAppId] = React.useState("");

  const {
    data: reviews,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["reviews", searchAppId],
    queryFn: () => fetchReviews(searchAppId),
    enabled: !!searchAppId,
    staleTime: 5 * 60 * 1000, // 5 minutes
    retry: false,
  });

  const handleSearch = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const formData = new FormData(e.currentTarget);
    const appId = formData.get("appId") as string;

    if (appId.trim()) {
      setSearchAppId(appId.trim());
    }
  };

  const handleAppleIdInputChange = () => {
    if (searchAppId) {
      setSearchAppId("");
      refetch();
    }
  };

  const errorMessage = getErrorMessage(error);

  return (
    <div className="app">
      <div className="reviews-container">
        <div className="search">
          <h1>Apps Reviews</h1>

          <form onSubmit={handleSearch} className="search__form">
            <div className="input-group">
              <Input
                type="text"
                name="appId"
                onChange={handleAppleIdInputChange}
                placeholder="Enter App Store App ID (e.g., 1234567890)"
                required
              />
              <Button type="submit" disabled={isLoading}>
                {isLoading ? "Loading..." : "Get Reviews"}
              </Button>
            </div>
          </form>
          {error && (
            <div className="message-container">
              <Message type="error">
                <Message.Title>{errorMessage.title}</Message.Title>
                <span>{errorMessage.description}</span>
                <Button onClick={() => refetch()} variant="link">
                  Try Again
                </Button>
              </Message>
            </div>
          )}
        </div>

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
