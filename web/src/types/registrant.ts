export type RegistrantGetByTicketID = {
  username: string;
  picture: string;
  is_verified: boolean;
  created_at: string;
};

export type PaginatedRegistrants = {
  data: RegistrantGetByTicketID[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
};

export type RegistrantGetByUserID = {
  id: string;
  title: string;
  url: string;
  image: string;
  is_verified: string;
  created_at: string;
};
