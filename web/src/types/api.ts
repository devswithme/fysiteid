export type APIMeta = {
  message: string;
  request_id: string;
  success: boolean;
};

export type APIResponse<T> = {
  data: T;
  meta: APIMeta;
};
