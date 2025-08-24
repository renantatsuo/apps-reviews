import { fetchReviews } from "~/lib/server/reviews";
import "./Review.css";

const renderStars = (rating: number) => {
  return "★".repeat(rating) + "☆".repeat(5 - rating);
};

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

interface ReviewProps {
  review: Awaited<ReturnType<typeof fetchReviews>>["data"][number];
}

export const Review = ({ review }: ReviewProps) => {
  return (
    <div className="review">
      <div className="review__header">
        <div className="review__rating">
          <span className="review__stars">{renderStars(review.rating)}</span>
          <span className="review__rating-number">({review.rating}/5)</span>
        </div>
        <div className="review__date">{formatDate(review.updated)}</div>
      </div>

      <h3 className="review__title">{review.title}</h3>

      <div className="review__content">
        <p>{review.content}</p>
      </div>

      <div className="review__author">
        <span>by {review.author.name}</span>
      </div>
    </div>
  );
};
