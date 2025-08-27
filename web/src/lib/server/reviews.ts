import * as HttpClient from "~/lib/http/typed-fetch";

type Review = {
  id: string;
  author: string;
  title: string;
  content: string;
  rating: number;
  sent_at: string;
};

type ReviewsResponse = {
  data: Review[];
};

export const fetchReviews = async (appId: string): Promise<ReviewsResponse> => {
  return await HttpClient.get<ReviewsResponse>(`/api/reviews/${appId}`);
};
