import * as HttpClient from "~/lib/http/typed-fetch";

type ReviewAuthor = {
  name: string;
  uri: string;
};

type Review = {
  id: string;
  author: ReviewAuthor;
  title: string;
  content: string;
  rating: number;
  updated: string; // ISO date string
};

type ReviewsResponse = {
  data: Review[];
};

export const fetchReviews = async (appId: string): Promise<ReviewsResponse> => {
  return await HttpClient.get<ReviewsResponse>(`/api/reviews/${appId}`);
};
