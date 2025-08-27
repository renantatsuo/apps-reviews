import * as HttpClient from "~/lib/http/typed-fetch";

export type App = {
  id: string;
  name: string;
  thumbnail_url: string;
  created_at: string;
  updated_at: string;
};

export type AppsResponse = {
  data: App[];
};

export type AppResponse = {
  data: App;
};

export const fetchApps = async (): Promise<AppsResponse> => {
  return await HttpClient.get<AppsResponse>("/api/apps");
};

export const fetchApp = async (appId: string): Promise<AppResponse> => {
  return await HttpClient.get<AppResponse>(`/api/apps/${appId}`);
};

export const addApp = async (appId: string): Promise<string> => {
  return await HttpClient.execute<string>(
    new Request(`/api/apps/${appId}`, { method: "POST" })
  );
};
